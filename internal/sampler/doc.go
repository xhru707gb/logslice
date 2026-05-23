// Package sampler provides probabilistic line-level sampling for log streams.
//
// It is designed to work alongside the filter and slicer packages:
// logs can first be sliced by time range and regex, then sampled down
// to a manageable volume for further analysis.
//
// Example usage:
//
//	s := sampler.New(sampler.Options{Rate: 0.1, Seed: 42})
//	n, err := s.Sample(inputReader, outputWriter)
package sampler
