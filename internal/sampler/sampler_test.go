package sampler

import (
	"bytes"
	"strings"
	"testing"
)

func lines(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString("log line\n")
	}
	return sb.String()
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n")
}

func TestSampler_FullRate(t *testing.T) {
	s := New(Options{Rate: 1.0, Seed: 1})
	input := lines(100)
	var out bytes.Buffer
	n, err := s.Sample(strings.NewReader(input), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 100 {
		t.Errorf("expected 100 lines, got %d", n)
	}
	if countLines(out.String()) != 100 {
		t.Errorf("output line count mismatch")
	}
}

func TestSampler_ZeroRate_ClampsToFull(t *testing.T) {
	s := New(Options{Rate: 0, Seed: 1})
	if s.opts.Rate != 1.0 {
		t.Errorf("expected rate clamped to 1.0, got %f", s.opts.Rate)
	}
}

func TestSampler_HalfRate(t *testing.T) {
	s := New(Options{Rate: 0.5, Seed: 99})
	input := lines(1000)
	var out bytes.Buffer
	n, err := s.Sample(strings.NewReader(input), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With seed 99 and rate 0.5 over 1000 lines we expect roughly 500 ± 10%
	if n < 400 || n > 600 {
		t.Errorf("expected ~500 lines, got %d", n)
	}
}

func TestSampler_EmptyInput(t *testing.T) {
	s := New(Options{Rate: 1.0, Seed: 1})
	var out bytes.Buffer
	n, err := s.Sample(strings.NewReader(""), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestSampler_Deterministic(t *testing.T) {
	input := lines(200)
	s1 := New(Options{Rate: 0.5, Seed: 7})
	s2 := New(Options{Rate: 0.5, Seed: 7})
	var out1, out2 bytes.Buffer
	s1.Sample(strings.NewReader(input), &out1) //nolint:errcheck
	s2.Sample(strings.NewReader(input), &out2) //nolint:errcheck
	if out1.String() != out2.String() {
		t.Error("same seed should produce identical output")
	}
}
