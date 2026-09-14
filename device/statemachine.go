package device

import (
	"fmt"

	modbus "github.com/goburrow/modbus"
)

type State struct {
	cfg Config
	u   [][]modbus.Client
	r   [][]Register
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

	return nil

}
