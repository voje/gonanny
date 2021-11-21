package nanny

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config corresponds to the config.yaml file
type Config struct {
	DailyTimeFrom      string `yaml:"daily_time_from"`
	DailyTimeTo        string `yaml:"daily_time_to"`
	DailyTimeAmountSec int    `yaml:"daily_time_amount_sec"`
	TickIntervalSec    int    `yaml:"tick_interval_sec"`
	DbFilePath         string `yaml:"db_file_path"`
	HttpPort           int    `yaml:"http_port"`
}

func ReadConfigFromFile(filePath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(bytes, c)
	return c, nil
}

// ParseConfig transforms config strings into valid types
func (n *Nanny) applyConfig(c *Config) error {
	hhmm := "15:04"
	var err error
	n.DailyTimeFrom, err = time.Parse(hhmm, c.DailyTimeFrom)
	if err != nil {
		return err
	}
	n.DailyTimeTo, err = time.Parse(hhmm, c.DailyTimeTo)
	if err != nil {
		return err
	}
	n.DailyTimeAmountSec = c.DailyTimeAmountSec
	n.TickIntervalSec = c.TickIntervalSec
	n.DbFilePath = c.DbFilePath
	if c.HttpPort == 0 {
		c.HttpPort = 8544
	}
	n.httpPort = c.HttpPort
	return nil
}
