package metrics

import (
	"hash/fnv"
	"sort"
)

type Metric struct {
	Reporter
	labels map[string]string
	name   string
	value  int64
}

func (m *Metric) labelKeys() (keys []string) {
	for k := range m.labels {
		keys = append(keys, k)
	}
	return
}

func (m *Metric) hash(suffix string) uint64 {
	keys := m.labelKeys()
	sort.Strings(keys)

	hasher := fnv.New64()
	hasher.Write([]byte(m.name))
	for _, k := range keys {
		hasher.Write([]byte(k))
		hasher.Write([]byte(m.labels[k]))
	}
	hasher.Write([]byte(suffix))
	return hasher.Sum64()
}

func (m *Metric) Observe(value int64) {
	m.value = value
	m.Reporter.Observe(m)
}
