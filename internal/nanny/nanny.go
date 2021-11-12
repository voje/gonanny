package nanny

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Interval          uint   `yaml:"interval"`
	TmpFile           string `yaml:"tmp_file"`
	DailyTimeLimitSec uint   `yaml:"daily_time_limit_sec"`
}

type State struct {
	LastSeen      time.Time     `yaml:"last_seen"`
	TimeAvailable time.Duration `yaml:"time_available"`
}

func ConfigFromFile(filePath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(bytes, c)
	return c, nil
}

type Nanny struct {
	conf  *Config
	state *State
}

func NewNanny(c *Config) *Nanny {
	n := &Nanny{}
	n.conf = c
	n.state = &State{}
	return n
}

func (n *Nanny) StoreState() error {
	// Create timer file if it doesn't exist
	yamlBytes, err := yaml.Marshal(n.state)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(
		n.conf.TmpFile,
		[]byte(yamlBytes),
		0755,
	)
	if err != nil {
		return err
	}
	return nil
}

func (n *Nanny) resetState() {
	n.state.LastSeen = time.Now()
	n.state.TimeAvailable = time.Duration(n.conf.DailyTimeLimitSec) * time.Second
}

func (n *Nanny) addDailyTime() {

}

func (n *Nanny) InitState() error {
	if _, err := os.Stat(n.conf.TmpFile); errors.Is(err, os.ErrNotExist) {
		n.resetState()
		n.addDailyTime()
		err := n.StoreState()
		if err != nil {
			return err
		}
	} else {
		_, err := n.GetPreviousState()
		if err != nil {
			return err
		}
	}
	return nil
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

func (n *Nanny) CheckTimeDiff(prevState *State) {
	timeDiff := n.state.Time.Sub(prevState.Time)
	fmt.Printf("Elapsed: %v", timeDiff)
}

func (n *Nanny) Run() error {
	n.InitState()

	// Start ticking
	ticker := time.NewTicker(time.Duration(n.conf.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:

			n.CheckTimeDiff()
		}
	}
}
