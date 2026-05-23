# logslice

Fast log file slicer and sampler for large structured log archives with regex and time-range filtering.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build ./...
```

---

## Usage

```bash
# Filter logs by time range
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log

# Filter using a regex pattern
logslice --pattern "ERROR|FATAL" app.log

# Sample 10% of matching lines from a compressed archive
logslice --pattern "timeout" --sample 0.1 logs/archive-2024-01.log.gz

# Combine time range, regex, and output to a file
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" \
         --pattern "user_id=42" --out slice.log app.log
```

### Flags

| Flag        | Description                              | Default |
|-------------|------------------------------------------|---------|
| `--from`    | Start of time range (RFC3339)            | —       |
| `--to`      | End of time range (RFC3339)              | —       |
| `--pattern` | Regex pattern to match log lines         | —       |
| `--sample`  | Fraction of matching lines to emit (0–1) | `1.0`   |
| `--out`     | Output file path (default: stdout)       | stdout  |

---

## Requirements

- Go 1.21+
- Supports `.log` and `.log.gz` input files

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

MIT © yourname