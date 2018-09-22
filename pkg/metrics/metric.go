package metrics

type MetricGroup interface {
	Observe(int64)
}

type StandardMetricGroup struct {
	Collector
	MetricGroup
	name   string
	labels map[string]string
}

func (s *StandardMetricGroup) Observe(delta int64) {
	s.Collector.Observe(s.name, s.labels, delta)
}
