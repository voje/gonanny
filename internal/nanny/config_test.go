package nanny

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNanny(t *testing.T) {
	config := &Config{
		DailyTimeFrom:      "08:30",
		DailyTimeTo:        "21:20",
		DailyTimeAmountSec: 25,
		TickIntervalSec:    60,
		DbFilePath:         "./test_files/test_db.yaml",
	}
	_, err := NewNanny(config)
	assert.NoError(t, err)
}

func TestNewNannyFail(t *testing.T) {
	config := &Config{
		DailyTimeFrom:      "222:30",
		DailyTimeTo:        "21:20",
		DailyTimeAmountSec: 25,
		TickIntervalSec:    60,
		DbFilePath:         "./test_files/test_db.yaml",
	}
	_, err := NewNanny(config)
	assert.Error(t, err)
}
