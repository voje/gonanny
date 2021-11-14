package nanny

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// State contains all data that should be persisted
// between program restarts
type State struct {
	availableTimeSec time.Duration `yaml:"available_time_sec"`
	lastUpdated      time.Time     `yaml:"lastUpdated"`
}

// initState creates a new state and populates it with init values or reads
// state from an existing file
// Keeps an open file descriptor
func (n *Nanny) initState() error {
	var err error
	n.DbFile, err = os.OpenFile(n.DbFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// First time running the app, let's initalize DailyTimeAmountSec
			n.state.availableTimeSec = n.DailyTimeAmountSec
		} else {
			return err
		}
	}
	n.state.lastUpdated = time.Now()
	fmt.Printf("%+v\n", n)
	fmt.Printf("%+v\n", n.state)
	err = n.storeState()
	if err != nil {
		return err
	}
	return nil
}

func (n *Nanny) storeState() error {
	fmt.Printf("%v", n.state)
	yamlBytes, err := yaml.Marshal(n.state)
	if err != nil {
		return err
	}
	fmt.Println("AAA")
	_, err = n.DbFile.Write(yamlBytes)
	if err != nil {
		return err
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
