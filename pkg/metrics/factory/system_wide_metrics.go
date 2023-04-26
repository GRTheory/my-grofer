package factory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GRTheory/my-grofer/pkg/metrics/general"
	"golang.org/x/sync/errgroup"
)

type systemWideMetrics struct {
	cpuInfo     bool
	refreshRate uint64
}

func NewSystemWideMetrics(cpuInfo bool, refreshRate uint64) *systemWideMetrics {
	return &systemWideMetrics{
		cpuInfo:     cpuInfo,
		refreshRate: refreshRate,
	}
}

func (swm *systemWideMetrics) Serve(opts ...Option) error {
	// apply command spcific options.
	for _, opt := range opts {
		opt(swm)
	}

	return swm.serveGenericMetrics()
}

func (swm *systemWideMetrics) serveGenericMetrics() error {
	eg, ctx := errgroup.WithContext(context.Background())
	metricBus := make(chan general.AggregateMetrics, 1)

	eg.Go(func() error {
		alteredRefreshRate := uint64(4 * swm.refreshRate / 5)
		return general.GlobalStats(ctx, metricBus, alteredRefreshRate)
	})

	eg.Go(func() error {
		for metric := range metricBus {
			var result []byte
			var err error
			switch metric.FieldSet {
			case "NET":
				result, err = json.Marshal(metric.NetStats)
			case "DISK":
				result, err = json.Marshal(metric.DiskStats)
			case "INFO":
				result, err = json.Marshal(metric.HostInfo)
			case "CPU":
				result, err = json.Marshal(metric.CPUStats)
			case "MEM":
				result, err = json.Marshal(metric.MemStats)
			}
			if err != nil {
				return err
			}
			fmt.Println(metric.FieldSet, string(result))
		}
		return nil
	})

	return eg.Wait()
}
