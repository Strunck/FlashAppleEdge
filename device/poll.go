package device

import (
	"fmt"
	"time"
)

func (s *State) pollClients() {

	for lnr, line := range s.u {
		for sid, client := range line {

			err := s.r[lnr][sid].Mess.Read(client)
			if err != nil {
				// Handle the error appropriately, e.g., log it
				fmt.Printf("Error reading Messung for line %d, slave %d: %v\n", lnr+1, sid+1, err)
			}
			s.r[lnr][sid].Mess.Print()
		}
	}

	time.Sleep(10 * time.Second)
}
