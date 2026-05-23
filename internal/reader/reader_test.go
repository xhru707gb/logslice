package reader_test

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func writePlain(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log*.log")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func writeGzip(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "log.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gw := gzip.NewWriter(f)
	for _, l := range lines {
		gw.Write([]byte(l + "\n"))
	}
	gw.Close()
	f.Close()
	return path
}

func collectLines(t *testing.T, path string) []string {
	t.Helper()
	r, err := reader.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer r.Close()

	var got []string
	for r.Scan() {
		got = append(got, r.Text())
	}
	if err := r.Err(); err != nil {
		t.Fatalf("Err: %v", err)
	}
	return got
}

func TestOpen_PlainFile(t *testing.T) {
	want := []string{"line one", "line two", "line three"}
	path := writePlain(t, want)
	got := collectLines(t, path)
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestOpen_GzipFile(t *testing.T) {
	want := []string{"compressed line 1", "compressed line 2"}
	path := writeGzip(t, want)
	got := collectLines(t, path)
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestOpen_MissingFile(t *testing.T) {
	_, err := reader.Open("/no/such/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
