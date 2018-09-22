package metrics

type CounterMetricGroup struct {
	*StandardMetricGroup
}

func (c *CounterMetricGroup) Inc() {
	c.StandardMetricGroup.Observe(1)
}

type CounterCollector struct {
	Name string
	*StandardCollector
}

func NewCounterCollector(name string) *CounterCollector {
	return &CounterCollector{
		Name:              name,
		StandardCollector: new(StandardCollector),
	}
}

func (c *CounterCollector) WithLabels(labels labelMap) *CounterMetricGroup {
	cm := new(CounterMetricGroup)
	cm.Collector = c
	cm.labels = labels
	cm.name = c.Name
	return cm
}
