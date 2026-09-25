package device

import (
	"context"
	"fmt"
	"os"

	gocron "github.com/go-co-op/gocron/v2"
	modbus "github.com/goburrow/modbus"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

const (
	DataFolder = "Daten"
	Nats_port  = 8082
)

type State struct {
	Cfg        Config
	sched      gocron.Scheduler
	Units      [][]Unit
	Metrics    Metrics
	initDone   bool
	NatsServer *server.Server
	NatsClient *nats.Conn
}

type Unit struct {
	m modbus.Client
	r Register
	f *os.File
}

func Bootstrap() (se State, err error) {
	fmt.Println("statemachine starts")

	se.Cfg, err = LoadConfig()
	if err != nil {
		return se, fmt.Errorf("Main LoadConfig error: %v", err)
	}

	// NATS Message Bus Initialization
	if err := se.initNataSrv(); err != nil {
		return se, fmt.Errorf("Main initNataSrv error: %v", err)
	}

	se.initUnits()
	PrintConfig(se.Cfg)

	if err := se.intiFiles(DataFolder); err != nil {
		return se, fmt.Errorf("Run intiFiles error: %v", err)
	}

	se.MakeClientsFromConfig()

	se.initMetrics()

	se.initDone = true

	return se, nil
}

func (se *State) Run(bgctx context.Context) {

	fmt.Println("Game Loop started")
	var err error
	errCh := make(chan error, 1)
	httpStarted := make(chan bool, 1)
	go RunWebServer(*se, errCh, httpStarted)

	select {
	case httpServerStarted := <-httpStarted:
		if !httpServerStarted {
			errCh <- fmt.Errorf("HTTP server failed to start")
			return
		} else {
			se.PollClients()

			se.sched, err = gocron.NewScheduler()
			if err != nil {
				errCh <- fmt.Errorf("Main NewScheduler error: %v", err)

			}

			go func() {
				err = se.pollHourly()
				if err != nil {
					errCh <- fmt.Errorf("Main pollHourly error: %v", err)
				}
			}()

			se.sched.Start()
		}
	case err := <-errCh:
		se.publishMsg(fmt.Sprintf("Error received: %v\n", err), true)

	case <-bgctx.Done():
		se.NatsClient.ClosedHandler()
		if se.sched != nil {
			se.sched.Shutdown()
		}
		se.CloseAllFiles()
	}
}
