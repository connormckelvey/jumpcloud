package pwhash

import (
	"encoding/json"
	"time"

	"github.com/connormckelvey/jumpcloud/metrics"
)

// AverageDuration is a metrics.Collector used for storing counts
// and accumulated values.
type AverageDuration struct {
	name  string
	count *metrics.Int64Collector
	sum   *metrics.Int64Collector
	units time.Duration
}

// NewAverageDuration returns an initialized AverageDuration.
func NewAverageDuration(name string, units time.Duration) *AverageDuration {
	return &AverageDuration{
		name:  name,
		count: metrics.NewInt64Collector(""),
		sum:   metrics.NewInt64Collector(""),
		units: units,
	}
}

// Name returns the name of the Collector. Name() is a requirement for the
// metrics.Collector interface.
func (c *AverageDuration) Name() string {
	return c.name
}

// Observe increments the AverageDuration.count value by 1 and adds the delta
// as an int64 (adjusted with AverageDuration.units) to AverageDuration.sum
func (c *AverageDuration) Observe(d time.Duration) {
	c.count.Add(1)
	c.sum.Add(int64(d / c.units))
}

// Count returns the value of AverageDuration.count
func (c *AverageDuration) Count() int64 {
	return c.count.Value()
}

// Sum returns the value of AverageDuration.sum
func (c *AverageDuration) Sum() int64 {
	return c.sum.Value()
}

// Average calculates and returns the average duration as an int64. To prevent
// a "Divide By Zero" panic, 0 is returned in the case that c.Count returns 0
func (c *AverageDuration) Average() int64 {
	count := c.Count()
	if count > 0 {
		return c.Sum() / count
	}
	return 0
}

// MarshalJSON implements the json.Marshaler interface allowing custom JSON formatting
// of AverageDuration for the /stats endpoint
func (c *AverageDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Count   int64 `json:"total"`
		Average int64 `json:"average"`
	}{
		c.Count(),
		c.Average(),
	})
}
