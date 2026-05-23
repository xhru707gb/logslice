// Package reader provides transparent log file reading for logslice.
//
// It supports plain-text log files as well as gzip-compressed archives
// (detected by a ".gz" suffix). Lines up to 1 MiB in length are handled
// correctly; longer lines are truncated by the underlying scanner.
//
// Typical usage:
//
//	r, err := reader.Open("app.log.gz")
//	if err != nil { ... }
//	defer r.Close()
//
//	for r.Scan() {
//		fmt.Println(r.Text())
//	}
//	if err := r.Err(); err != nil { ... }
package reader
