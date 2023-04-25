package factory

import (
	"context"
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
		cpuInfo: cpuInfo,
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
			fmt.Println(metric)
		}
		return nil
	})

	return eg.Wait()
}
