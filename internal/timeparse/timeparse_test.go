package timeparse_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
)

func mustUTC(s, layout string) time.Time {
	t, err := time.ParseInLocation(layout, s, time.UTC)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}

func TestParse_RFC3339(t *testing.T) {
	got, err := timeparse.Parse("2024-03-15T08:30:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := mustUTC("2024-03-15T08:30:00Z", time.RFC3339)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParse_DateOnly(t *testing.T) {
	got, err := timeparse.Parse("2024-03-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := mustUTC("2024-03-15", "2006-01-02")
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParse_SpaceSeparated(t *testing.T) {
	got, err := timeparse.Parse("2024-03-15 12:00:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Hour() != 12 {
		t.Errorf("expected hour 12, got %d", got.Hour())
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := timeparse.Parse("not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid input, got nil")
	}
}

func TestParseRange_BothSet(t *testing.T) {
	from, to, err := timeparse.ParseRange("2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if from.IsZero() || to.IsZero() {
		t.Fatal("expected non-zero times")
	}
	if !to.After(from) {
		t.Errorf("expected to > from")
	}
}

func TestParseRange_EmptyFrom(t *testing.T) {
	from, to, err := timeparse.ParseRange("", "2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.IsZero() {
		t.Errorf("expected zero from, got %v", from)
	}
	if to.IsZero() {
		t.Error("expected non-zero to")
	}
}

func TestParseRange_InvertedRange(t *testing.T) {
	_, _, err := timeparse.ParseRange("2024-12-31", "2024-01-01")
	if err == nil {
		t.Fatal("expected error for inverted range, got nil")
	}
}

func TestParseRange_BothEmpty(t *testing.T) {
	from, to, err := timeparse.ParseRange("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.IsZero() || !to.IsZero() {
		t.Error("expected both times to be zero")
	}
}
