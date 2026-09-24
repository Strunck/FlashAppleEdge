package device

import (
	"fmt"

	gocron "github.com/go-co-op/gocron/v2"
)

func (s *State) PollClients(errCh chan error) {
	fmt.Println("Polling startet...")

	for lnr, line := range s.Units {
		ort := s.Cfg.Lines[lnr].Ort
		for sid := 0; sid < len(line); sid++ {

			slaveId := s.Cfg.Lines[lnr].Id[sid]
			loc := fmt.Sprintf("[Line%d] %s▪️%d", lnr+1, ort, slaveId)

			unit := s.Units[lnr][sid]

			err := unit.r.Trig.Write(unit.m, F_Messung)
			if err != nil {
				errCh <- fmt.Errorf("Error Trig(F_Messung) %s: %v\n", loc, err)
			}

			err = unit.r.Mess.Read(unit.m)
			if err != nil {
				// Handle the error appropriately, e.g., log it
				errCh <- fmt.Errorf("Error reading Messung %s: %v\n", loc, err)
			} else {
				unit.r.Print(lnr, sid)
				if err := unit.writeLineInCSV(); err != nil {
					errCh <- fmt.Errorf("Error writing line to CSV %s: %v\n", loc, err)
				}
			}
		}
	}
}

func (s *State) pollHourly(errCh chan error) (err error) {

	_, err = s.sched.NewJob(
		gocron.CronJob("0 * * * *", false), // Every hour at minute 0
		gocron.NewTask(func() { s.PollClients(errCh) }),
	)
	return err
}
