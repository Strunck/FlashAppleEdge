package device

import (
	"fmt"

	gocron "github.com/go-co-op/gocron/v2"
)

func (s *State) PollClients() {
	s.publishMsg("Polling startet...", true)

	for lnr, line := range s.Units {
		ort := s.Cfg.Lines[lnr].Ort
		for sid := 0; sid < len(line); sid++ {

			slaveId := s.Cfg.Lines[lnr].Id[sid]
			loc := fmt.Sprintf("[Line%d] %s▪️%d", lnr+1, ort, slaveId)

			unit := s.Units[lnr][sid]

			s.publishMsg(fmt.Sprintf("F_Messung %s", loc), true)
			err := unit.r.Trig.Write(unit.m, F_Messung)
			if err != nil {
				s.publishMsg(fmt.Sprintf("Error Trig(F_Messung) %s: %v\n", loc, err), true)
			}

			err = unit.r.Mess.Read(unit.m)
			if err != nil {
				s.publishMsg(fmt.Sprintf("Error reading Messung %s: %v\n", loc, err), true)
				continue
			} else {
				unit.r.Print(lnr, sid)
				if err := unit.writeLineInCSV(); err != nil {
					s.publishMsg(fmt.Sprintf("Error writing line to CSV %s: %v\n", loc, err), true)
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
