// Package timeparse provides flexible timestamp parsing for logslice.
//
// It supports a variety of common date/time layouts so that users can
// specify time-range bounds without needing to know the exact RFC 3339
// format.  Layouts are tried from most specific to least specific;
// the first successful parse is returned.
//
// Typical usage:
//
//	from, to, err := timeparse.ParseRange(fromFlag, toFlag)
//	if err != nil {
//		log.Fatal(err)
//	}
package timeparse
