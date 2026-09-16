package device

import (
	"fmt"

	gocron "github.com/go-co-op/gocron/v2"
)

func (s *State) pollClients() {

	for lnr, line := range s.units {
		for sid := 0; sid < len(line); sid++ {
			err := s.units[lnr][sid].r.Trig.Write(s.units[lnr][sid].m, F_Messung)
			if err != nil {
				fmt.Printf("Error writing trigger for line %d, slave %d: %v\n", lnr+1, sid+1, err)
			}

			err = s.units[lnr][sid].r.Mess.Read(s.units[lnr][sid].m)
			if err != nil {
				// Handle the error appropriately, e.g., log it
				fmt.Printf("Error reading Messung for line %d, slave %d: %v\n", lnr+1, sid+1, err)
			} else {
				s.units[lnr][sid].r.Print(lnr, sid)
				if err := s.units[lnr][sid].writeLineInCSV(); err != nil {
					fmt.Printf("Error writing line to CSV for line %d, slave %d: %v\n", lnr+1, sid+1, err)
				}
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
