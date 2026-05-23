package sampler

import (
	"bufio"
	"io"
	"math/rand"
)

// Options holds sampling configuration.
type Options struct {
	// Rate is the fraction of lines to keep (0.0 < Rate <= 1.0).
	Rate float64
	// Seed is the random seed; 0 means use default source.
	Seed int64
}

// Sampler selects a random subset of lines from a reader.
type Sampler struct {
	opts Options
	rng  *rand.Rand
}

// New creates a Sampler with the given options.
// If Rate is outside (0, 1] it is clamped to 1.0.
func New(opts Options) *Sampler {
	if opts.Rate <= 0 || opts.Rate > 1.0 {
		opts.Rate = 1.0
	}
	var src rand.Source
	if opts.Seed != 0 {
		src = rand.NewSource(opts.Seed)
	} else {
		src = rand.NewSource(42)
	}
	return &Sampler{
		opts: opts,
		rng:  rand.New(src),
	}
}

// Sample reads lines from r, writes a random sample to w, and returns
// the number of lines written and any error encountered.
func (s *Sampler) Sample(r io.Reader, w io.Writer) (int, error) {
	scanner := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)
	written := 0

	for scanner.Scan() {
		if s.rng.Float64() < s.opts.Rate {
			if _, err := bw.WriteString(scanner.Text() + "\n"); err != nil {
				return written, err
			}
			written++
		}
	}
	if err := scanner.Err(); err != nil {
		return written, err
	}
	return written, bw.Flush()
}
