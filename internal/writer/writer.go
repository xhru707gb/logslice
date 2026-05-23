// Package writer provides output writing utilities for logslice,
// supporting plain text and gzip-compressed output streams.
package writer

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// WriteCloser wraps an io.WriteCloser with an optional underlying file
// so both can be closed in the correct order.
type WriteCloser struct {
	w    io.WriteCloser
	file *os.File
}

// Write implements io.Writer.
func (wc *WriteCloser) Write(p []byte) (int, error) {
	return wc.w.Write(p)
}

// Close flushes and closes the writer, then closes the underlying file if present.
func (wc *WriteCloser) Close() error {
	if err := wc.w.Close(); err != nil {
		return fmt.Errorf("writer close: %w", err)
	}
	if wc.file != nil {
		if err := wc.file.Close(); err != nil {
			return fmt.Errorf("file close: %w", err)
		}
	}
	return nil
}

// Create opens path for writing. If path ends with ".gz" the output is
// gzip-compressed. Pass "-" or an empty string to write to stdout.
func Create(path string) (*WriteCloser, error) {
	if path == "" || path == "-" {
		return wrap(os.Stdout, nil), nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create output file: %w", err)
	}

	if strings.HasSuffix(path, ".gz") {
		gw := gzip.NewWriter(f)
		return &WriteCloser{w: gw, file: f}, nil
	}

	return wrap(f, f), nil
}

// wrap builds a WriteCloser; when w and file are the same object we still
// store file so Close can handle the two-step shutdown uniformly.
func wrap(w io.WriteCloser, file *os.File) *WriteCloser {
	return &WriteCloser{w: w, file: file}
}

// WriteLines writes each line followed by a newline to w.
func WriteLines(w io.Writer, lines []string) error {
	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return nil
}
