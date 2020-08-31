package report

import (
	"report-manager/collector"
	"report-manager/logger"
)

// Maker can make report by report template and specified collectors.
type Maker struct {
	ReportName string
	Template   string
	Collectors []collector.Collector
}

func NewMaker(name string, template string, collectors []collector.Collector) *Maker {
	return &Maker{ReportName: name, Template: template, Collectors: collectors}
}

// Make makes report.
func (m *Maker) Make() (string, error) {
	logger.Infof("[report] make %s report begin", m.ReportName)
	defer logger.Infof("[report] make %s report done", m.ReportName)
	collector.Collect(m.Collectors)

	ignore := m.evaluateIgnore()
	if ignore {
		return "", DoNotReport
	}

	// render
	for i := range m.Collectors {
		m.Template = m.Collectors[i].Render(m.Template)
	}
	return m.Template, nil
}

func (m *Maker) evaluateIgnore() bool {
	for _, c := range m.Collectors {
		ignorer, ok := c.(collector.Ignorer)
		if !ok { // un-ignorable
			return false
		}
		if !ignorer.Ignore() {
			return false
		}
	}
	return true
}
