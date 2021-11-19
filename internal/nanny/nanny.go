package nanny

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Nanny struct {
	// Values from config
	DailyTimeFrom      time.Time
	DailyTimeTo        time.Time
	DailyTimeAmountSec float64
	TickIntervalSec    time.Duration
	DbFilePath         string

	state *State
}

func NewNanny(c *Config) (*Nanny, error) {
	var err error
	n := &Nanny{
		state: &State{},
	}
	err = n.applyConfig(c)
	if err != nil {
		return nil, err
	}
	err = n.initState(time.Now())
	if err != nil {
		return nil, err
	}
	return n, nil
}

func daySeconds(t time.Time) int {
	year, month, day := t.Date()
	t2 := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return int(t.Sub(t2).Seconds())
}

func (n *Nanny) withinAllowedTimeInterval(currentTime time.Time) bool {
	fromSec := daySeconds(n.DailyTimeFrom)
	toSec := daySeconds(n.DailyTimeTo)
	currentSec := daySeconds(currentTime)
	if currentSec >= fromSec && currentSec <= toSec {
		return true
	}
	return false
}

func (n *Nanny) suspendUser() {
	n.storeState(time.Now())
	log.Info("Shutting down! TODO")
}

func (n *Nanny) Run() error {
	// Init nanny

	// Start ticking
	ticker := time.NewTicker(time.Duration(n.TickIntervalSec))
	for {
		select {
		case <-ticker.C:
			log.Info("Tick")

			n.addDailyTime(time.Now())

			if n.state.AvailableTimeSec <= 0 ||
				!n.withinAllowedTimeInterval(time.Now()) {
				n.suspendUser()
			}
		}
	}
}
