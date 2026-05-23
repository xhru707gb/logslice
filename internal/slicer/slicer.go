// Package slicer reads a log source line-by-line and writes matching lines to a destination.
package slicer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/logslice/logslice/internal/filter"
)

// Stats holds counters collected during a slicing run.
type Stats struct {
	LinesRead    int64
	LinesMatched int64
	Errors       int64
}

// Slice reads lines from r, applies cfg, and writes matching lines to w.
// It returns aggregated Stats and the first non-EOF error encountered.
func Slice(r io.Reader, w io.Writer, cfg *filter.Config) (Stats, error) {
	var stats Stats
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1<<20), 1<<20) // 1 MiB line buffer

	for scanner.Scan() {
		line := scanner.Bytes()
		stats.LinesRead++

		ok, err := cfg.Filter(line)
		if err != nil {
			stats.Errors++
			continue
		}
		if !ok {
			continue
		}

		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return stats, err
		}
		stats.LinesMatched++
	}

	if err := scanner.Err(); err != nil {
		return stats, err
	}
	return stats, nil
}
