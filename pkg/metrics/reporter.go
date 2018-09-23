package metrics

type Reporter interface {
	Observe(*Metric)
}

type StandardReporter struct {
	Name string
}

func (s *StandardReporter) WithLabels(labels map[string]string) *Metric {
	return &Metric{
		Reporter: s,
		name:     s.Name,
		labels:   labels,
	}
}

func (s *StandardReporter) incrementCount(metric *Metric) {
	defaultStore.Update(metric.hash("count"), func(v int64) int64 {
		return v + 1
	})
}

func (s *StandardReporter) addToSum(metric *Metric) {
	defaultStore.Update(metric.hash("sum"), func(v int64) int64 {
		return v + metric.value
	})
}

func (s *StandardReporter) Observe(metric *Metric) {
	s.incrementCount(metric)
	s.addToSum(metric)
}
