package general

import (
	"context"
	"testing"
)

func TestServeStatsOfDisk(t *testing.T) {
	dataChannel := make(chan AggregateMetrics)
	go func() {
		err := ServeDiskRates(context.Background(), dataChannel)
		if err != nil {
			t.Log("failed to get disk-related messages", err)
			return
		}
	}()
	t.Log(<-dataChannel)
}

func TestServeStatsOfMemRates(t *testing.T) {
	dataChannel := make(chan AggregateMetrics)
	go func() {
		err := ServeMemRates(context.Background(), dataChannel)
		if err != nil {
			t.Log("failed to get memory-related messges", err)
			return
		}
	}()
	t.Log(<-dataChannel)
}
