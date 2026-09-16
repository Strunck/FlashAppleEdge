package device

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gocron "github.com/go-co-op/gocron/v2"
	modbus "github.com/goburrow/modbus"
)

const (
	DataFolder = "Daten"
)

type State struct {
	cfg   Config
	sched gocron.Scheduler
	units [][]Unit
}

type Unit struct {
	m modbus.Client
	r Register
	f *os.File
}

func Run() (err error) {
	// Initialize Backend State
	se := State{}

	se.cfg, err = LoadConfig()
	if err != nil {
		fmt.Printf("Main LoadConfig error: %v", err)
		return err
	}
	// Initialisiert ein 2D-Array für die
	se.make2dArr()

	PrintConfig(se.cfg)

	if err := se.intiFiles(DataFolder); err != nil {
		return fmt.Errorf("Run intiFiles error: %v", err)
	}

	// Make Modbus  Clients and start
	if err := se.MakeClientsFromConfig(); err != nil {
		return fmt.Errorf("Run MakeClientsFromConfig error: %v", err)
	}

	se.sched, err = gocron.NewScheduler()
	defer func() { _ = se.sched.Shutdown() }()

	err = se.pollHourly()
	if err != nil {
		fmt.Printf("Main pollHourly error: %v", err)
		return err
	}
	se.sched.Start()

	bgctx := context.Background()

	ctx, userstop := signal.NotifyContext(bgctx, os.Interrupt, syscall.SIGTERM)

	select {
	case <-make(chan struct{}):
	// This will block forever, effectively keeping the program running

	case <-ctx.Done():
		fmt.Println("Received interrupt signal, shutting down...")
		_ = se.CloseAllFiles()
		userstop()
	}

	return nil

}

func (s *State) make2dArr() {
	s.units = make([][]Unit, s.cfg.Count)

	for l, line := range s.cfg.Lines {

		s.units[l] = make([]Unit, len(line.Id))
	}
}
