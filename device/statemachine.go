package device

import (
	"context"
	"fmt"
	"os"

	gocron "github.com/go-co-op/gocron/v2"
	modbus "github.com/goburrow/modbus"
)

const (
	DataFolder = "Daten"
)

type State struct {
	Cfg       Config
	sched     gocron.Scheduler
	Units     [][]Unit
	Metrics   Metrics
	initState InitStateType
	initDone  bool
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
)

func Run(bgctx context.Context) {

	fmt.Println("statemachine starts")

	se := State{}
	errCh := make(chan error, 1)
	var err error

	initDoneCh := make(chan bool)

	se.setInitState(KONFIG_FILE)

	fmt.Println("statemachine select")
	for {
		select {
		case err := <-errCh:
			fmt.Printf("Error received: %v\n", err)
		default:
			switch se.initState {
			case KONFIG_FILE:

				se.Cfg, err = LoadConfig()
				if err != nil {
					errCh <- fmt.Errorf("Main LoadConfig error: %v", err)
					se.setInitState(SHUTDOWN)
				}
				se.setInitState(INIT_UNITS)

			case INIT_NATS: // NATS initialization
				se.setInitState(INIT_UNITS)

			case INIT_UNITS:
				se.initUnits()
				PrintConfig(se.Cfg)
				if err := se.intiFiles(DataFolder); err != nil {
					errCh <- fmt.Errorf("Run intiFiles error: %v", err)
				}
				se.setInitState(READY_TO_POLL)

			case READY_TO_POLL:
				// Start Polling
				go se.MakeClientsFromConfig(errCh, initDoneCh)
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
					err = se.pollHourly(errCh)
					if err != nil {
						errCh <- fmt.Errorf("Main pollHourly error: %v", err)
					}
				}()
				se.setInitState(RUNNING)

			case RUNNING:
				se.initDone = true
				se.sched.Start()
				se.PollClients(errCh)

			case SHUTDOWN:
				func() { _ = se.sched.Shutdown() }()

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
	case SHUTDOWN:
		stateTxt = "SHUTDOWN"
	default:
		stateTxt = "UNKNOWN"
	}

	s.initState = newState
	fmt.Println("State changed to:", stateTxt)
}
