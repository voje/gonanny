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
	DailyTimeAmountSec: 30,
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

func TestOutsieTimeLimits(t *testing.T) {
	config := testConfig
	config.DailyTimeFrom = "08:30"
	config.DailyTimeTo = "21:20"
	os.Remove(config.DbFilePath)
	n, _ := NewNanny(&config)

	testMap := make(map[string]bool)
	testMap["2014-11-12T05:45:26.371Z"] = false
	testMap["2014-11-12T11:45:26.371Z"] = true
	testMap["2014-11-12T21:30:26.371Z"] = false
	testMap["2014-11-12T00:30:26.371Z"] = false

	for key, el := range testMap {
		testTime, err := time.Parse(time.RFC3339, key)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Checking: %s", key)
		assert.Equal(t, n.withinAllowedTimeInterval(testTime), el)
	}
	_ = n
}