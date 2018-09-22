package metrics

import (
	"fmt"
	"sync"
)

type Registry struct {
	collectors map[Collector]bool
	metrics    []*Metric
	mux        sync.RWMutex
}

func (r *Registry) register(col Collector) error {
	if _, exists := r.collectors[col]; exists {
		return fmt.Errorf("Collector %s is already registered", col)
	}
	r.collectors[col] = true
	return nil
}

func (r *Registry) Register(cols ...Collector) error {
	for _, col := range cols {
		if err := r.register(col); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) Retrieve() {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.metrics = []*Metric{}

	collectors := make(chan Collector)
	go func() {
		for c := range r.collectors {
			collectors <- c
		}
		close(collectors)
	}()

	// That was dumb, MetricGroup Makes more sense, or MetricFacts
	metrics := make(chan *Metric)

	go func() {
		wg := sync.WaitGroup{}

		for c := range collectors {
			wg.Add(1)

			colMetrics := make(chan *Metric)
			go c.Retrieve(colMetrics)

			go func(cm chan *Metric) {
				for m := range cm {
					metrics <- m
				}

				wg.Done()
			}(colMetrics)
		}

		wg.Wait()
		close(metrics)
	}()

	for m := range metrics {
		r.metrics = append(r.metrics, m)
	}
}

func foo() {
	c := NewCounterCollector("http_request_count")
	c.WithLabels(map[string]string{"path": "/hash"}).Inc()

	s := new(StandardCollector)
	s.WithLabels(map[string]string{"path": "/hash"}).Observe(400)

}
