package writer_test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/writer"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(data)
}

func readGzip(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open gz: %v", err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer gr.Close()
	b, err := io.ReadAll(gr)
	if err != nil {
		t.Fatalf("read gz: %v", err)
	}
	return string(b)
}

func TestCreate_PlainFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	wc, err := writer.Create(path)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := writer.WriteLines(wc, []string{"line1", "line2"}); err != nil {
		t.Fatalf("WriteLines: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	got := readFile(t, path)
	if !strings.Contains(got, "line1") || !strings.Contains(got, "line2") {
		t.Errorf("unexpected content: %q", got)
	}
}

func TestCreate_GzipFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log.gz")

	wc, err := writer.Create(path)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := writer.WriteLines(wc, []string{"hello", "world"}); err != nil {
		t.Fatalf("WriteLines: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	got := readGzip(t, path)
	if !strings.Contains(got, "hello") || !strings.Contains(got, "world") {
		t.Errorf("unexpected gz content: %q", got)
	}
}

func TestCreate_Stdout(t *testing.T) {
	// Just verify Create("-") does not error and is closeable without
	// closing the real stdout.
	wc, err := writer.Create("-")
	if err != nil {
		t.Fatalf("Create('-'): %v", err)
	}
	// Do NOT call wc.Close() here — it would close os.Stdout.
	_ = wc
}
