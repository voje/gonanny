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
	AvailableTimeSec int       `yaml:"available_time_sec"`
	LastUpdated      time.Time `yaml:"lastUpdated"`
}

// initState creates a new state and populates it with init values or reads
// state from an existing file
// Keeps an open file descriptor
func (n *Nanny) initState(currentTime time.Time) (err error) {
	if _, err := os.Stat(n.DbFilePath); errors.Is(err, os.ErrNotExist) {
		// File doesn't exist; initialize
		log.Infof("First initialization, creating file: %s", n.DbFilePath)
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
		log.Infof("Setting inital AvailableTimeSec: %d", n.state.AvailableTimeSec)
		n.state.LastUpdated = currentTime
	} else {
		log.Infof("Initializing from existing state file: %s", n.DbFilePath)
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

// storeState sets LastUpdated to currentTime (should be time.Now() when not testing).
// Stores state to file.
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
