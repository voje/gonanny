package nanny

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Nanny struct {
	// Values from config
	DailyTimeTo        time.Time
	DailyTimeFrom      time.Time
	DailyTimeAmountSec time.Duration
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
	err = n.initState()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Nanny) Run() error {
	// Init nanny

	// Start ticking
	ticker := time.NewTicker(time.Duration(n.TickIntervalSec))
	for {
		select {
		case <-ticker.C:
			log.Info("Tick")

			// n.CheckTimeDiff()
		}
	}
}
