package nanny

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testFilesDir = "test_files"

func setup() {
	os.Mkdir(testFilesDir, 0766)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	// shutdown()
	os.Exit(code)
}

var testConfig Config = Config{
	DailyTimeFrom:      "08:30",
	DailyTimeTo:        "21:20",
	DailyTimeAmountSec: 25,
	TickIntervalSec:    60,
	DbFilePath:         fmt.Sprintf("./%s/test_db.yaml", testFilesDir),
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

func TestAddDailytime(t *testing.T) {
	config := testConfig
	os.Remove(config.DbFilePath)
	n, err := NewNanny(&config)
	// Store state, simulate 2 days ago
	n.storeState(time.Now().Add(-time.Hour * 48))
	_ = n
	n = nil
	n, err = NewNanny(&config)
	n.addDailyTime(time.Now())
	t.Logf("AvailablePlayTime: %f", n.state.AvailableTimeSec)
	assert.NoError(t, err)
}
