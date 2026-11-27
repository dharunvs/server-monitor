# Reproducible benchmark protocol

Run on an isolated machine from the commit under test:

```bash
go version
go test -race ./...
go test ./monitor -run '^$' -bench . -benchmem -count 5
```

Record commit SHA, Go version, OS, CPU, count, ns/op, B/op and allocs/op. Compare revisions on the same machine and toolchain. The benchmark uses fake probes and storage and does not measure PostgreSQL or network latency.

For a future end-to-end benchmark, run a disposable PostgreSQL instance and local TCP fixtures, vary target count and worker count, measure scheduler lag, completed checks per second, CPU, memory and shutdown time, and report dropped/deferred observations explicitly. Do not benchmark third-party hosts without permission.

No measured results are claimed until the commands have actually run.
