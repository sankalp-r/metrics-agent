package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSystemSource(t *testing.T) {
	expectedResponse := "# HELP algod_inoctets Algod inoctets (bytes)\n# TYPE algod_inoctets gauge\nalgod_inoctets 100\n# HELP algod_outoctets Algod outoctets (bytes)\n# TYPE algod_outoctets gauge\nalgod_outoctets 100\n# HELP algod_starttime Algod process start time\n# TYPE algod_starttime gauge\nalgod_starttime 1\n# HELP algod_stime Algod process stime\n# TYPE algod_stime gauge\nalgod_stime 1\n# HELP algod_utime Algod process utime\n# TYPE algod_utime gauge\nalgod_utime 1\n"
	testSystemSource := NewSystemSource("algod")
	findProcessID = mockFindPID
	findProcessStat = mockProcessStat

	b, err := testSystemSource.Collect()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResponse, string(b))
}

func mockFindPID(processName string) (string, error) {
	return "123", nil
}

func mockProcessStat(pid int) (*Stat, error) {
	return &Stat{
			Utime:     1,
			Stime:     1,
			StartTime: 1,
			InOctets:  100,
			OutOctets: 100,
		},
		nil
}
