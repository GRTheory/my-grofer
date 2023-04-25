package factory

// MetricScraper scrapes metrics of some form and serves it based on
// the implementation.
type MetricScraper interface {
	// Serve serves the metrics to a 'sink', which can be a TUI or
	// logic that exports these served metrics to either a file or
	// maybe even served over an endpoint (some day).
	Serve(opts ...Option) error
}
