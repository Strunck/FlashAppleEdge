package device

import (
	"fmt"

	gocron "github.com/go-co-op/gocron/v2"
)

func (s *State) PollClients() {
	fmt.Println("Polling startet...")

	for lnr, line := range s.Units {
		ort := s.Cfg.Lines[lnr].Ort
		for sid := 0; sid < len(line); sid++ {

			slaveId := s.Cfg.Lines[lnr].Id[sid]
			loc := fmt.Sprintf("[Line%d] %s▪️%d", lnr+1, ort, slaveId)

			unit := s.Units[lnr][sid]

			err := unit.r.Trig.Write(unit.m, F_Messung)
			if err != nil {
				s.NatsClient.Publish("pool", fmt.Appendf(nil, "Error Trig(F_Messung) %s: %v\n", loc, err))
			}

			err = unit.r.Mess.Read(unit.m)
			if err != nil {
				// Handle the error appropriately, e.g., log it
				s.NatsClient.Publish("pool", fmt.Appendf(nil, "Error reading Messung %s: %v\n", loc, err))
				continue
			} else {
				unit.r.Print(lnr, sid)
				if err := unit.writeLineInCSV(); err != nil {
					s.NatsClient.Publish("pool", fmt.Appendf(nil, "Error writing line to CSV %s: %v\n", loc, err))
					continue
				}
			}
		}
	}

}

func (s *State) pollHourly() (err error) {

	_, err = s.sched.NewJob(
		gocron.CronJob("0 * * * *", false), // Every hour at minute 0
		gocron.NewTask(func() {
			s.PollClients()
		}),
	)
	return err
}
