package factory

// Option is used to inject command specific configuration.
type Option func(MetricScraper)

// WithAllAs sets the all flag value for the ContainerCommand.
// func WithAllAs(all bool) Option {
// 	return func(ms MetricScraper) {
// 		cms := ms.(*containerMetrics)
// 		cms.all = all
// 	}
// }

// // WithCPUInfoAs sets the cpuinfo flag value for the RootCommand.
func WithCPUInfoAs(cpuInfo bool) Option {
	return func(ms MetricScraper) {
		swm := ms.(*systemWideMetrics)
		swm.cpuInfo = cpuInfo
	}
}
