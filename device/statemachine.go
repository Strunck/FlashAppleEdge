package device

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	gocron "github.com/go-co-op/gocron/v2"
	modbus "github.com/goburrow/modbus"
)

type State struct {
	cfg   Config
	u     [][]modbus.Client
	r     [][]Register
	sched gocron.Scheduler
}

func Run() (err error) {
	// Initialize Backend State
	se := State{}

	se.cfg, err = LoadConfig()
	if err != nil {
		fmt.Printf("Main LoadConfig error: %v", err)
		return err
	}
	PrintConfig(se.cfg)

	err = se.MakeClientsFromConfig()
	if err != nil {
		fmt.Printf("Main MakeClientsFromConfig error: %v", err)
		return err
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

	ctx, userstop := signal.NotifyContext(bgctx, os.Interrupt)

	select {
	case <-make(chan struct{}):
	// This will block forever, effectively keeping the program running
	case <-ctx.Done():
		userstop()
	}

	return nil

}
