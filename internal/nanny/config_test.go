package nanny

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfig Config = Config{
	DailyTimeFrom:      "08:30",
	DailyTimeTo:        "21:20",
	DailyTimeAmountSec: 25,
	TickIntervalSec:    60,
	DbFilePath:         "./test_files/test_db.yaml",
}

func TestNewNannyNewDbFile(t *testing.T) {
	config := testConfig
	os.Remove(config.DbFilePath)
	_, err := NewNanny(&config)
	assert.NoError(t, err)
}

func TestNewNannyBadConfig(t *testing.T) {
	config := testConfig
	config.DailyTimeFrom = "222:30"
	os.Remove(config.DbFilePath)
	_, err := NewNanny(&config)
	assert.Error(t, err)
}

func TestNewNannyExistingDbFile(t *testing.T) {
	config := testConfig
	os.Remove(config.DbFilePath)
	n, err := NewNanny(&config)
	assert.NoError(t, err)
	// Close object (TODO: does this release the FD?)
	_ = n
	n = nil
	_, err = NewNanny(&config)
	assert.NoError(t, err)
}
