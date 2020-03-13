package collector

type Collector interface {
	Collect() error
	Render(string) string
}
