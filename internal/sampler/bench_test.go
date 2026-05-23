package sampler

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func BenchmarkSampler_Rate50(b *testing.B) {
	const numLines = 10_000
	var sb strings.Builder
	for i := 0; i < numLines; i++ {
		fmt.Fprintf(&sb, `{"ts":"2024-01-01T00:00:00Z","level":"info","msg":"request","id":%d}`+"\n", i)
	}
	input := sb.String()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s := New(Options{Rate: 0.5, Seed: int64(i + 1)})
		s.Sample(strings.NewReader(input), io.Discard) //nolint:errcheck
	}
}

func BenchmarkSampler_Rate100(b *testing.B) {
	const numLines = 10_000
	var sb strings.Builder
	for i := 0; i < numLines; i++ {
		fmt.Fprintf(&sb, `{"ts":"2024-01-01T00:00:00Z","level":"error","msg":"timeout","id":%d}`+"\n", i)
	}
	input := sb.String()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s := New(Options{Rate: 1.0, Seed: 1})
		var out bytes.Buffer
		out.Grow(len(input))
		s.Sample(strings.NewReader(input), &out) //nolint:errcheck
	}
}
