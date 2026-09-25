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
	initState  InitStateType
	initDone   bool
	NatsServer *server.Server
	NatsClient *nats.Conn
}

type Unit struct {
	m modbus.Client
	r Register
	f *os.File
}

type InitStateType int

const (
	KONFIG_FILE InitStateType = iota
	INIT_NATS
	INIT_WEBSERVER
	INIT_UNITS
	INIT_METRIC    // Fehlerfälle?
	READY_TO_POLL  // Fehlerfälle?
	INIT_SCHEDULER // Fehlerfälle?
	RUNNING        // msg to topic:Errors || Data
	SHUTDOWN
	FATAL
)

func Run(bgctx context.Context) {

	fmt.Println("statemachine starts")

	se := State{}
	errCh := make(chan error, 1)
	var err error

	se.setInitState(KONFIG_FILE)

	fmt.Println("statemachine select")
	for {
		select {
		case err := <-errCh:
			se.NatsClient.Publish("Error", fmt.Appendf(nil, "Error received: %v\n", err))
			fmt.Printf("Error received: %v\n", err)
		default:
			switch se.initState {
			case KONFIG_FILE:

				se.Cfg, err = LoadConfig()
				if err != nil {
					errCh <- fmt.Errorf("Main LoadConfig error: %v", err)
					se.setInitState(FATAL)
				}
				se.setInitState(INIT_UNITS)

			case INIT_NATS: // NATS initialization
				if err := se.initNataSrv(); err != nil {
					errCh <- fmt.Errorf("Main initNataSrv error: %v", err)
					se.setInitState(FATAL)
				} else {
					se.setInitState(INIT_UNITS)
				}

			case INIT_UNITS:
				se.initUnits()
				PrintConfig(se.Cfg)
				if err := se.intiFiles(DataFolder); err != nil {
					errCh <- fmt.Errorf("Run intiFiles error: %v", err)
				}
				se.setInitState(READY_TO_POLL)

			case READY_TO_POLL:
				// Start Polling
				se.MakeClientsFromConfig()
				se.setInitState(INIT_METRIC)

			case INIT_METRIC:
				se.initMetrics()
				se.setInitState(INIT_WEBSERVER)

			case INIT_WEBSERVER: // WEBSERVER initialization
				go RunWebServer(se, errCh)
				se.setInitState(INIT_SCHEDULER)

			case INIT_SCHEDULER:
				se.sched, err = gocron.NewScheduler()
				if err != nil {
					errCh <- fmt.Errorf("Main NewScheduler error: %v", err)
					se.setInitState(SHUTDOWN)
				}

				go func() {
					err = se.pollHourly()
					if err != nil {
						errCh <- fmt.Errorf("Main pollHourly error: %v", err)
					}
				}()
				se.setInitState(RUNNING)

			case RUNNING:
				if !se.initDone {
					se.initDone = true
					se.sched.Start()
					se.PollClients()
				}

			case SHUTDOWN:
				func() { _ = se.sched.Shutdown() }()

			case FATAL:
				fmt.Println("Fatal error received, shutting down...")
				os.Exit(1)
			}
		}
	}
}

func (s *State) setInitState(newState InitStateType) {

	var stateTxt string
	switch newState {
	case KONFIG_FILE:
		stateTxt = "KONFIG_FILE"
	case INIT_NATS:
		stateTxt = "INIT_NATS"
	case INIT_WEBSERVER:
		stateTxt = "INIT_WEBSERVER"
	case INIT_UNITS:
		stateTxt = "INIT_UNITS"
	case INIT_METRIC:
		stateTxt = "INIT_METRIC"
	case READY_TO_POLL:
		stateTxt = "READY_TO_POLL"
	case INIT_SCHEDULER:
		stateTxt = "INIT_SCHEDULER"
	case RUNNING:
		stateTxt = "RUNNING"
	case FATAL:
		stateTxt = "FATAL"
	case SHUTDOWN:
		stateTxt = "SHUTDOWN"
	default:
		stateTxt = "UNKNOWN"
	}

	s.initState = newState
	fmt.Println("State changed to:", stateTxt)
}
