package nanny

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Nanny struct {
	// Values from config
	DailyTimeFrom      time.Time
	DailyTimeTo        time.Time
	DailyTimeAmountSec int
	TickIntervalSec    int
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
	log.Infof("%s is outside allowed time interval.", currentTime)
	return false
}

// addDailyTime checks the amount of days since (var state.LastUpdated) and
// adds days * n.DailyAmountSec to n.AvailableTimeSec
func (n *Nanny) addDailyTime(currentTime time.Time) {
	daysSinceLastLogin := int(currentTime.Sub(n.state.LastUpdated).Hours() / 24)
	nSec := daysSinceLastLogin * n.DailyTimeAmountSec
	log.Infof("Last logged in %d days ago, adding %d seconds.",
		daysSinceLastLogin, nSec,
	)
	n.state.AvailableTimeSec += nSec
	n.state.LastUpdated = currentTime
	n.storeState(time.Now())
}

func (n *Nanny) subtractAvailableTime(nsec int) {
	n.state.AvailableTimeSec -= nsec
	if n.state.AvailableTimeSec < 0 {
		n.state.AvailableTimeSec = 0
	}
	n.storeState(time.Now())
}

func (n *Nanny) Run() error {
	// Init nanny
	n.addDailyTime(time.Now())

	// Run http server that displays app info to the user
	go n.runServer()

	// Start ticking
	ticker := time.NewTicker(time.Duration(n.TickIntervalSec) * time.Second)
	for {
		select {
		case <-ticker.C:
			n.subtractAvailableTime(n.TickIntervalSec)
			log.Infof("Available time in seconds: %d", n.state.AvailableTimeSec)
			n.storeState(time.Now())

			if n.state.AvailableTimeSec <= 0 ||
				!n.withinAllowedTimeInterval(time.Now()) {
				n.systemMessage()
				n.systemShutdown()
			}
		}
	}
}
