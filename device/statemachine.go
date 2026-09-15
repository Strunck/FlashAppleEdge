package device

import (
	"fmt"

	modbus "github.com/goburrow/modbus"
	"github.com/go-co-op/gocron"
)

type State struct {
	cfg Config
	u   [][]modbus.Client
	r   [][]Register
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
	defer func() { _= se.sched.Shutdown() }()

	err = se.pollHourly()
	if err != nil {
		fmt.Printf("Main pollHourly error: %v", err)
		return err
	}
	se.sched.Start()


	return nil

}
