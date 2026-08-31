# benchmarks

Benchmark, parity, and pressure-test tooling for the Gnalloy modular networking stack.

This repository is part of the Gnalloy modular networking stack. The default branch is `dev`; no release tag is created during bootstrap.

## Install

```bash
go get gnalloy.org/benchmarks@dev
```

## Module Boundary

- Module path: `gnalloy.org/benchmarks`
- Responsibility: Benchmark, parity, and pressure-test tooling for the Gnalloy modular networking stack
- Core dependency: `gnalloy.org/gnalloy` when this module uses Gnalloy buffers, channels, event loops, or bootstrap contracts.

## Gnalloy Dependencies

- `gnalloy.org/gnalloy`
- `gnalloy.org/codec-http3`
- `gnalloy.org/handler-tls`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

## Development

```bash
go test ./... -count=1
go vet ./...
go test ./... -run '^$' -bench . -benchmem -benchtime=100ms -count=1
```

For multi-repository development, use the workspace at `G:\opensource\gnalloy\go.work`. For standalone verification, set `GOWORK=off`.

## License

Apache-2.0.
