package slicer_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/slicer"
)

var sampleLog = `2024-01-01T10:00:00 INFO  service started
2024-01-01T10:01:00 ERROR disk full
2024-01-01T10:02:00 INFO  request handled
2024-01-01T10:03:00 ERROR timeout reached
2024-01-01T10:04:00 INFO  shutdown complete
`

func TestSlice_NoFilter(t *testing.T) {
	var out bytes.Buffer
	stats, err := slicer.Slice(strings.NewReader(sampleLog), &out, &filter.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if stats.LinesRead != 5 || stats.LinesMatched != 5 {
		t.Fatalf("expected 5/5, got %d/%d", stats.LinesMatched, stats.LinesRead)
	}
}

func TestSlice_RegexFilter(t *testing.T) {
	var out bytes.Buffer
	cfg := &filter.Config{Pattern: regexp.MustCompile(`ERROR`)}
	stats, err := slicer.Slice(strings.NewReader(sampleLog), &out, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if stats.LinesMatched != 2 {
		t.Fatalf("expected 2 ERROR lines, got %d", stats.LinesMatched)
	}
	if !strings.Contains(out.String(), "disk full") {
		t.Error("output missing 'disk full'")
	}
}

func TestSlice_EmptyInput(t *testing.T) {
	var out bytes.Buffer
	stats, err := slicer.Slice(strings.NewReader(""), &out, &filter.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if stats.LinesRead != 0 {
		t.Fatalf("expected 0 lines read, got %d", stats.LinesRead)
	}
}
