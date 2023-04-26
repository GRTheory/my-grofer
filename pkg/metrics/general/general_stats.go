package general

import (
	"context"
	"sync"
)

// AggregateMetrics represents global metrics to be consumed.
type AggregateMetrics struct {
	NetStats  []NetMonitor
	FieldSet  string
	CPUStats  []CPUMonitor
	MemStats  MemMonitor
	DiskStats []DiskMonitor
	HostInfo  InfoMonitor
}

type serveFunc func(context.Context, chan AggregateMetrics) error

// GlobalStats gets stats about the mem and CPUs and prints it.
func GlobalStats(ctx context.Context, dataChannel chan AggregateMetrics, _ uint64) error {
	serveFuncs := []serveFunc{
		ServeCPURates,
		ServeMemRates,
		ServeDiskRates,
		ServeNetRates,
		ServeInfo,
	}

	return func(ctx context.Context) error {
		defer close(dataChannel)
		var wg sync.WaitGroup
		errCh := make(chan error, len(serveFuncs))

		for _, sf := range serveFuncs {
			wg.Add(1)
			go func(sf serveFunc, dc chan AggregateMetrics) {
				defer wg.Done()
				errCh <- sf(ctx, dc)
			}(sf, dataChannel)
		}
		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return err
			}
		}
		return nil
	}(ctx)
}
