package device

import (
	"fmt"

	gocron "github.com/go-co-op/gocron/v2"
)

func (s *State) pollClients() {

	for lnr, line := range s.u {
		for sid, client := range line {

			err := s.r[lnr][sid].Mess.Read(client)
			if err != nil {
				// Handle the error appropriately, e.g., log it
				fmt.Printf("Error reading Messung for line %d, slave %d: %v\n", lnr+1, sid+1, err)
			} else {
				s.r[lnr][sid].Print(lnr, sid)
			}
		}
	}
}

func (s *State) pollHourly() (err error) {

	_, err = s.sched.NewJob(
		gocron.CronJob("0 * * * *", false), // Every hour at minute 0
		gocron.NewTask(s.pollClients),
	)
	return err
}
