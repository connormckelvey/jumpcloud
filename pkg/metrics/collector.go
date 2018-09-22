package metrics

import "sync"

type labelHashMap = map[hashedLabels]labelMap
type dataHashMap = map[hashedLabels]int64

type Collector interface {
	Observe(string, labelMap, int64)
	Retrieve(chan<- *Metric)
}

func NewCollector() Collector {
	return &StandardCollector{
		mux:    sync.RWMutex{},
		labels: make(labelHashMap),
		data:   make(dataHashMap),
	}
}

type StandardCollector struct {
	mux    sync.RWMutex
	labels labelHashMap
	data   dataHashMap
}

func (s *StandardCollector) WithLabels(labels labelMap) MetricGroup {
	return &StandardMetricGroup{
		Collector: s,
		labels:    labels,
	}
}

func (s *StandardCollector) Observe(name string, labels labelMap, value int64) {
	labels["__metric_name"] = name
	hash := s.findOrCreateLabelHash(labels)
	s.mux.Lock()
	defer s.mux.Unlock()
	s.data[hash] += value
}

type Metric struct {
	labels labelMap
	value  int64
}

func (s *StandardCollector) Retrieve(metrics chan<- *Metric) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	labelHashes := make(chan hashedLabels)

	go func() {
		for k := range s.labels {
			labelHashes <- k
		}
		close(labelHashes)
	}()

	for labelHash := range labelHashes {
		metrics <- &Metric{
			labels: s.labels[labelHash],
			value:  s.data[labelHash],
		}
	}
	close(metrics)
}

func (s *StandardCollector) findOrCreateLabelHash(ls labelMap) uint64 {
	hash := ls.hash()
	if _, exists := s.findLabelsByHash(hash); !exists {
		s.mux.Lock()
		s.labels[hash] = ls
		s.mux.Unlock()
	}
	return hash
}

func (s *StandardCollector) findLabelsByHash(hash uint64) (map[string]string, bool) {
	s.mux.RLock()
	defer s.mux.Unlock()
	if labels, exists := s.labels[hash]; exists {
		return labels, true
	}
	return nil, false
}
