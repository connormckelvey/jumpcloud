package pwhash

import (
	"encoding/json"
	"time"

	"github.com/connormckelvey/jumpcloud/metrics"
)

type averageDuration struct {
	name  string
	count *metrics.Int64Collector
	sum   *metrics.Int64Collector
	units time.Duration
}

func NewAverageDuration(name string, units time.Duration) *averageDuration {
	return &averageDuration{
		name:  name,
		count: metrics.NewInt64Collector(""),
		sum:   metrics.NewInt64Collector(""),
		units: units,
	}
}

func (c *averageDuration) Name() string {
	return c.name
}

func (c *averageDuration) Observe(d time.Duration) {
	c.count.Add(1)
	c.sum.Add(int64(d / c.units))
}

func (c *averageDuration) Count() int64 {
	return c.count.Value()
}

func (c *averageDuration) Sum() int64 {
	return c.sum.Value()
}

func (c *averageDuration) Average() int64 {
	count := c.Count()
	if count > 0 {
		return c.Sum() / count
	}
	return 0
}

func (c *averageDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Count   int64 `json:"total"`
		Average int64 `json:"average"`
	}{
		c.Count(),
		c.Average(),
	})
}
