// Package timeparse provides helpers for parsing time range arguments
// supplied on the command line or in configuration files.
package timeparse

import (
	"fmt"
	"time"
)

// Common layouts tried in order when parsing a free-form timestamp string.
var layouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
	"2006-01-02",
}

// Parse attempts to parse s using a set of common timestamp layouts.
// It returns the parsed time in UTC or an error if no layout matches.
func Parse(s string) (time.Time, error) {
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, s, time.UTC)
		if err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("timeparse: cannot parse %q: no matching layout", s)
}

// ParseRange parses two optional timestamp strings into a (from, to) pair.
// Empty strings are returned as zero time.Time values (meaning "no bound").
func ParseRange(from, to string) (time.Time, time.Time, error) {
	var (
		t0, t1 time.Time
		err    error
	)
	if from != "" {
		t0, err = Parse(from)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("timeparse: from: %w", err)
		}
	}
	if to != "" {
		t1, err = Parse(to)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("timeparse: to: %w", err)
		}
	}
	if !t0.IsZero() && !t1.IsZero() && t1.Before(t0) {
		return time.Time{}, time.Time{}, fmt.Errorf("timeparse: 'to' (%s) is before 'from' (%s)", t1, t0)
	}
	return t0, t1, nil
}
