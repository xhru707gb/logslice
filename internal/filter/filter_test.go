package filter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
)

const layout = "2006-01-02T15:04:05"

func mustTime(s string) time.Time {
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestFilter_NoConstraints(t *testing.T) {
	cfg := &filter.Config{}
	ok, err := cfg.Filter([]byte("any log line"))
	if err != nil || !ok {
		t.Fatalf("expected line to pass with no constraints, got ok=%v err=%v", ok, err)
	}
}

func TestFilter_RegexMatch(t *testing.T) {
	cfg := &filter.Config{Pattern: regexp.MustCompile(`ERROR`)}

	ok, _ := cfg.Filter([]byte("2024-01-01T10:00:00 ERROR something broke"))
	if !ok {
		t.Fatal("expected ERROR line to match")
	}

	ok, _ = cfg.Filter([]byte("2024-01-01T10:00:00 INFO all good"))
	if ok {
		t.Fatal("expected INFO line to not match ERROR pattern")
	}
}

func TestFilter_TimeRange(t *testing.T) {
	cfg := &filter.Config{
		TimeLayout: layout,
		TimeField:  0,
		TimeFrom:   mustTime("2024-01-01T09:00:00"),
		TimeTo:     mustTime("2024-01-01T11:00:00"),
	}

	ok, _ := cfg.Filter([]byte("2024-01-01T10:00:00 INFO inside range"))
	if !ok {
		t.Fatal("expected line inside range to pass")
	}

	ok, _ = cfg.Filter([]byte("2024-01-01T08:00:00 INFO before range"))
	if ok {
		t.Fatal("expected line before range to be filtered")
	}

	ok, _ = cfg.Filter([]byte("2024-01-01T12:00:00 INFO after range"))
	if ok {
		t.Fatal("expected line after range to be filtered")
	}
}

func TestFilter_ShortLine(t *testing.T) {
	cfg := &filter.Config{
		TimeLayout: layout,
		TimeField:  0,
		TimeFrom:   mustTime("2024-01-01T09:00:00"),
	}
	ok, err := cfg.Filter([]byte("short"))
	if ok || err != nil {
		t.Fatalf("expected short line to be silently dropped, got ok=%v err=%v", ok, err)
	}
}
