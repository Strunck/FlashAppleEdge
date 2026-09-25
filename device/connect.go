package device

import (
	"fmt"
	"time"

	modbus "github.com/goburrow/modbus"
)

// Erstelle Modbus Clients anhand der Eingelesenen Konfiguration
// 2d Array. [line][SlaveID]
func (s *State) MakeClientsFromConfig() {

	// Implement the logic to create Modbus clients based on the provided configuration
	for lnr, line := range s.Cfg.Lines {
		switch line.ConType {
		case "TCP":
			s.publishMsg(fmt.Sprintf("Creating TCP client for line %d", lnr+1), true)
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
				s.Units[lnr][sid].m = modbus.NewClient(handler)
			}
		case "RTU":
			s.publishMsg(fmt.Sprintf("Creating RTU client for line %d", lnr+1), true)
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
				s.Units[lnr][sid].m = modbus.NewClient(handler)
			}
		}
	}
}
