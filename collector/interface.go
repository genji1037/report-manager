package collector

type Collector interface {
	Collect() error
	Render(string) string
}

// Ignorer is ignorable Collector
// Ignore() == true means we can ignore collected data.
type Ignorer interface {
	Collector
	Ignore() bool
}
