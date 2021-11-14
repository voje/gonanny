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
	AvailableTimeSec time.Duration `yaml:"available_time_sec"`
	LastUpdated      time.Time     `yaml:"lastUpdated"`
}

// initState creates a new state and populates it with init values or reads
// state from an existing file
// Keeps an open file descriptor
func (n *Nanny) initState() (err error) {
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
	n.state.LastUpdated = time.Now()
	err = n.storeState()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (n *Nanny) storeState() error {
	log.Printf("state: %+v\n", n.state)
	yamlBytes, err := yaml.Marshal(*n.state)
	log.Printf("marshalled: %v\n", yamlBytes)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(n.DbFilePath, yamlBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

/*
func (n *Nanny) resetState() {
	n.state.LastSeen = time.Now()
	n.state.TimeAvailable = time.Duration(n.conf.DailyTimeLimitSec) * time.Second
}

func (n *Nanny) addDailyTime() {

}


func (n *Nanny) GetPreviousState() (*State, error) {
	// Read data from existing timer
	prevState := &State{}
	bytes, err := ioutil.ReadFile(n.conf.TmpFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, prevState)
	if err != nil {
		return nil, err
	}
	return prevState, nil
}
*/
