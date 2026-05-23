// Package reader provides utilities for reading log files with support
// for plain text and gzip-compressed archives.
package reader

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// LineReader wraps a buffered scanner over a (possibly compressed) log file.
type LineReader struct {
	scanner *bufio.Scanner
	closer  io.Closer
}

// Open opens the file at path for line-by-line reading.
// Files ending in ".gz" are transparently decompressed.
func Open(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reader: open %q: %w", path, err)
	}

	var rc io.ReadCloser = f

	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("reader: gzip %q: %w", path, err)
		}
		rc = &multiCloser{Reader: gr, closers: []io.Closer{gr, f}}
	}

	scanner := bufio.NewScanner(rc)
	scanner.Buffer(make([]byte, 1<<20), 1<<20) // 1 MiB max line

	return &LineReader{scanner: scanner, closer: rc}, nil
}

// Scan advances to the next line. Returns false when done or on error.
func (r *LineReader) Scan() bool {
	return r.scanner.Scan()
}

// Text returns the current line as a string.
func (r *LineReader) Text() string {
	return r.scanner.Text()
}

// Err returns the first non-EOF error encountered.
func (r *LineReader) Err() error {
	return r.scanner.Err()
}

// Close releases resources held by the reader.
func (r *LineReader) Close() error {
	return r.closer.Close()
}

// multiCloser wraps a reader and closes multiple closers in order.
type multiCloser struct {
	io.Reader
	closers []io.Closer
}

func (m *multiCloser) Close() error {
	var first error
	for _, c := range m.closers {
		if err := c.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
