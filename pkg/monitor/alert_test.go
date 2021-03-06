package monitor

import (
	"testing"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
	"github.com/stretchr/testify/assert"
)

func TestAboveThresholdTrafficAlert(t *testing.T) {
	alertMonitor := newAlertMonitor(2*time.Minute, 10)

	now := time.Now()

	for i := 0; i < 10; i++ {
		alert := alertMonitor.addLogEntry(newLogEntry(now))
		assert.Nil(t, alert)
	}

	alert, ok := alertMonitor.addLogEntry(newLogEntry(now)).(*TrafficAlert)
	assert.True(t, ok)
	assert.True(t, alert.AboveThreshold())
	assert.Equal(t, int64(11), alert.AverageTraffic)
}

func TestUnderThresholdTrafficAlert(t *testing.T) {
	alertMonitor := newAlertMonitor(2*time.Minute, 10)

	now := time.Now()

	for i := 0; i <= 10; i++ {
		alertMonitor.addLogEntry(newLogEntry(now))
	}

	entry := newLogEntry(now.Add(2*time.Minute + 1*time.Second))
	alert, ok := alertMonitor.addLogEntry(entry).(*TrafficAlert)
	assert.True(t, ok)
	assert.True(t, alert.UnderThreshold())
	assert.Equal(t, int64(1), alert.AverageTraffic)
}

func newLogEntry(at time.Time) commonformat.LogEntry {
	return commonformat.LogEntry{
		Time: at,
	}
}
