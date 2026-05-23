// Package filter provides log line filtering by time range and regex pattern.
package filter

import (
	"regexp"
	"time"
)

// Config holds the filtering configuration for a log slicing operation.
type Config struct {
	// TimeFrom filters out log lines before this time (zero value means no lower bound).
	TimeFrom time.Time
	// TimeTo filters out log lines after this time (zero value means no upper bound).
	TimeTo time.Time
	// Pattern is an optional compiled regex applied to each log line.
	Pattern *regexp.Regexp
	// TimeLayout is the Go time layout used to parse timestamps in log lines.
	TimeLayout string
	// TimeField is the byte offset in the line where the timestamp starts.
	TimeField int
}

// Filter evaluates a single raw log line against the configured criteria.
// It returns true if the line should be included in the output.
func (c *Config) Filter(line []byte) (bool, error) {
	if c.TimeLayout != "" && (!c.TimeFrom.IsZero() || !c.TimeTo.IsZero()) {
		if len(line) < c.TimeField+len(c.TimeLayout) {
			return false, nil
		}
		raw := string(line[c.TimeField : c.TimeField+len(c.TimeLayout)])
		t, err := time.Parse(c.TimeLayout, raw)
		if err != nil {
			return false, err
		}
		if !c.TimeFrom.IsZero() && t.Before(c.TimeFrom) {
			return false, nil
		}
		if !c.TimeTo.IsZero() && t.After(c.TimeTo) {
			return false, nil
		}
	}

	if c.Pattern != nil && !c.Pattern.Match(line) {
		return false, nil
	}

	return true, nil
}
