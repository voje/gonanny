package nanny

import (
	"errors"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// State contains all data that should be persisted
// between program restarts
type State struct {
	AvailableTimeSec float64   `yaml:"available_time_sec"`
	LastUpdated      time.Time `yaml:"lastUpdated"`
}

// initState creates a new state and populates it with init values or reads
// state from an existing file
// Keeps an open file descriptor
func (n *Nanny) initState(currentTime time.Time) (err error) {
	if _, err := os.Stat(n.DbFilePath); errors.Is(err, os.ErrNotExist) {
		// File doesn't exist; initialize
		log.Infof("Initializing file: %s", n.DbFilePath)
		file, err := os.Create(n.DbFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if file == nil {
			log.Fatalf(
				"Opened file descriptor for path %s is <nil>",
				n.DbFilePath,
			)
		}
		// On first init, fund the user with some playtime
		n.state.AvailableTimeSec = n.DailyTimeAmountSec
		log.Infof("Initializing AvailableTimeSec: %f", n.state.AvailableTimeSec)
		n.state.LastUpdated = currentTime
	} else {
		// Read previous state
		yamlBytes, err := os.ReadFile(n.DbFilePath)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(yamlBytes, n.state)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (n *Nanny) addDailyTime(currentTime time.Time) {
	daysSinceLastLogin := currentTime.Sub(n.state.LastUpdated).Hours() / 24
	nSec := daysSinceLastLogin * float64(n.DailyTimeAmountSec)
	log.Infof("Last logged in %f days ago, adding %f seconds.",
		daysSinceLastLogin, nSec,
	)
	n.state.AvailableTimeSec += nSec
	n.state.LastUpdated = currentTime
}

func (n *Nanny) storeState(currentTime time.Time) error {
	n.state.LastUpdated = currentTime
	yamlBytes, err := yaml.Marshal(*n.state)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(n.DbFilePath, yamlBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
