package nanny

import (
	"math"
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
	httpPort           int

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

// addDailyTime checks the amount of days since (var state.LastUpdated).
// In case of same calendar day, do nothing.
// In case of different days, calculate the difference and round up.
func (n *Nanny) addDailyTime(currentTime time.Time) {
	defer n.storeState(currentTime)
	lastUpdatedDay := n.state.LastUpdated.Format("2")
	currentDay := currentTime.Format("2")
	if lastUpdatedDay == currentDay {
		return
	}
	daysSinceLastLogin := int(math.Ceil(currentTime.Sub(n.state.LastUpdated).Hours() / 24))
	nSec := daysSinceLastLogin * n.DailyTimeAmountSec
	if nSec > 0 {
		log.Infof("Last logged in %d days ago, adding %d seconds.",
			daysSinceLastLogin, nSec,
		)
		n.state.AvailableTimeSec += nSec
	}
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
			n.addDailyTime(time.Now())
			n.subtractAvailableTime(n.TickIntervalSec)
			n.storeState(time.Now())
			log.Infof("Available time in seconds: %d", n.state.AvailableTimeSec)

			if n.state.AvailableTimeSec <= 0 ||
				!n.withinAllowedTimeInterval(time.Now()) {
				n.systemMessage("Shutting down!")
				n.systemShutdown()
			}
		}
	}
}
