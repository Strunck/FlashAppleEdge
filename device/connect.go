package device

import (
	"fmt"
	"time"

	modbus "github.com/goburrow/modbus"
)

// Erstelle Modbus Clients anhand der Eingelesenen Konfiguration
// 2d Array. [line][SlaveID]
func (s *State) MakeClientsFromConfig() (err error) {

	c := make([][]modbus.Client, s.cfg.Count)
	r := make([][]Register, s.cfg.Count)
	for i := range c {
		c[i] = make([]modbus.Client, s.cfg.Lines[i].IdCount)
		r[i] = make([]Register, s.cfg.Lines[i].IdCount)
	}

	// Implement the logic to create Modbus clients based on the provided configuration
	for lnr, line := range s.cfg.Lines {
		switch line.ConType {
		case "TCP":
			fmt.Println("Creating TCP client for line", lnr+1)
			// Create a Modbus client for each line based on its configuration
			handler := modbus.NewTCPClientHandler(line.Url)
			handler.Connect()
			defer handler.Close()

			for _, slaveID := range line.Id {
				handler.SlaveId = byte(slaveID)
			}

			// client creation moved inside the loop above
			for sid, slaveID := range line.Id {
				handler.SlaveId = byte(slaveID)
				handler.IdleTimeout = 3 * time.Hour
				handler.Timeout = 5 * time.Second
				c[lnr][sid] = modbus.NewClient(handler)
			}
		case "RTU":
			fmt.Println("Creating RTU client for line", lnr+1)
			// Create a Modbus client for each line based on its configuration
			handler := modbus.NewRTUClientHandler(line.Url)
			handler.Connect()
			defer handler.Close()

			for sid, slaveID := range line.Id {
				handler.SlaveId = byte(slaveID)
				handler.BaudRate = 19200
				handler.DataBits = 8
				handler.Parity = "N"
				handler.StopBits = 1
				handler.Timeout = 5 * time.Second
				c[lnr][sid] = modbus.NewClient(handler)
			}
		}
	}
	s.u = c // Modbus Clients
	s.r = r // Modbus Register and Methods to poll

	s.pollClients()

	return nil
}
