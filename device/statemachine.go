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
	Cfg     Config
	sched   gocron.Scheduler
	Units   [][]Unit
	Metrics Metrics
}

type Unit struct {
	m modbus.Client
	r Register
	f *os.File
}

func Run(bgctx context.Context, se *State, errCh chan error, initDoneCh chan bool) {
	// Initialize Backend State
	var err error
	se.Cfg, err = LoadConfig()

	if err != nil {
		errCh <- fmt.Errorf("Main LoadConfig error: %v", err)
	}

	// Initialisiert ein 2D-Array für die
	se.Units = make([][]Unit, se.Cfg.Count)
	for l, line := range se.Cfg.Lines {
		se.Units[l] = make([]Unit, len(line.Id))
	}

	PrintConfig(se.Cfg)

	if err := se.intiFiles(DataFolder); err != nil {
		errCh <- fmt.Errorf("Run intiFiles error: %v", err)
	}

	// Make Modbus  Clients and start
	go se.MakeClientsFromConfig(errCh, initDoneCh)

	se.sched, err = gocron.NewScheduler()
	defer func() { _ = se.sched.Shutdown() }()

	go func() {
		err = se.pollHourly(errCh)
		if err != nil {
			errCh <- fmt.Errorf("Main pollHourly error: %v", err)
		}
	}()
	se.sched.Start()
}
