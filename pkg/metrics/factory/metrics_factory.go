package factory

// // MetricScraperFactory constructs a MetricScaper for a command
// // and returns it.
// type MetricScraperFactory struct {
// 	// command is the command for which a MetricScraper
// 	// is created. This defaults to core.MainCommand.
// 	command core.Command
// 	// singualrEntityMetrics indicate whether metrics
// 	// that need to be scraped are for a singular entity
// 	// or not, for ex - metrics for some process ID or
// 	// some container ID.
// 	singularEntityMetrics bool
// 	// entity is an identifier that can be used to scrape
// 	// metrics for it, ex CID, PID. Conversion from string
// 	// to the appropriate form of this entity should be
// 	// handled by an implementation of the MetricScraper
// 	// interface.
// 	entity string
// 	// scrapeIntervalMillisecond is the frequency in ms at
// 	// which metrics will be scraped.
// 	scrapeIntervalMillisecond uint64
// }

// // NewMetricScraperFactory is a constructor for the MetricScraperFactory type.
// // By default, this will be for the core.MainCommand command.
// func NewMetricScraperFactory() *MetricScraperFactory {
// 	return &MetricScraperFactory{}
// }

// // ForCommand sets the command for which a MetricScraper needs to be constructed.
// func (msf *MetricScraperFactory) ForCommand(command core.Command) *MetricScraperFactory {
// 	msf.command = command
// 	return msf
// }

// // ForSingularEntity sets the factory to construct an entity specific MetricScraper.
// func (msf *MetricScraperFactory) ForSingularEntity(entity string) *MetricScraperFactory {
// 	msf.singularEntityMetrics = true
// 	msf.entity = entity
// 	return msf
// }

// // WithScrapeInterval sets teh scrape interval for the factory.
// func (msf *MetricScraperFactory) WithScrapeInterval(interval uint64) *MetricScraperFactory {
// 	msf.scrapeIntervalMillisecond = interval
// 	return msf
// }

// // Construct constructs the MetricScraper for a particular Command and returns it.
// func (msf *MetricScraperFactory) Construct() (MetricScraper, error) {
// 	switch msf.command {
// 	case core.RootCommand:

// 	}
// }

// func (msf *MetricScraperFactory) constructSystemWideMetricScraper() (MetricScraper, error) {
// 	return &systemWideMetrics{
// 		refreshRate: msf.scrapeIntervalMillisecond,
// 	}, nil
// }

// func (msf *MetricScraperFactory) constructContainerMetricScraper() (MetricScraper, error) {
// 	if msf.singularEntityMetrics {

// 	}
// }

// func (msf *MetricScraperFactory) newContainerMetrics() (*metrics)