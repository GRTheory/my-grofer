package general

import (
	"context"
	"fmt"
	"sync"
)

// AggregateMetrics represents global metrics to be consumed.
type AggregateMetrics struct {
	NetStats  map[string][]float64
	FieldSet  string
	CPUStats  []float64
	MemStats  []float64
	DiskStats [][]string
	HostInfo  [][]string
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
		fmt.Println("returned")
		return nil
	}(ctx)
}
