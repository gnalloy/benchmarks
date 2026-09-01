# benchmarks

[简体中文](README.zh-CN.md) | [Documentation](docs/README.md)

Benchmark, parity, and pressure-test tooling for the Gnalloy modular networking stack.

This module contains benchmark, parity, and pressure-test tooling. It is for measurement and regression protection, not for production runtime dependencies.

## Status

- Import path: `gnalloy.org/benchmarks`
- Repository: `github.com/gnalloy/benchmarks`
- Default branch: `dev`
- Preview install: `go get gnalloy.org/benchmarks@dev`
- License: Apache-2.0

## Install
```bash
go get gnalloy.org/benchmarks@dev
go doc gnalloy.org/benchmarks
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
```

## Documentation
- [Overview](docs/overview.md) ([中文](docs/overview.zh-CN.md))
- [Usage](docs/usage.md) ([中文](docs/usage.zh-CN.md))
- [Examples](docs/examples.md) ([中文](docs/examples.zh-CN.md))
- [Configuration](docs/configuration.md) ([中文](docs/configuration.zh-CN.md))
- [Testing and Performance](docs/testing.md) ([中文](docs/testing.zh-CN.md))
- [API Reference](docs/api.md) ([中文](docs/api.zh-CN.md))
- [Notes and Caveats](docs/notes.md) ([中文](docs/notes.zh-CN.md))
- [ADR-001 Module Boundary](docs/decisions/0001-module-boundary.md) ([中文](docs/decisions/0001-module-boundary.zh-CN.md))

## Module Boundary

This repository owns: Benchmark, parity, and pressure-test tooling for the Gnalloy modular networking stack.

It does not absorb neighboring module responsibilities. Core primitives stay in `gnalloy.org/gnalloy`; protocol codecs, transports, handlers, resolvers, examples, and benchmarks stay in their own repositories.

## Packages
- `gnalloy.org/benchmarks/benchdiff` (`benchdiff`)
- `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff` (`main`)
- `gnalloy.org/benchmarks/external/internal/benchh2` (`benchh2`)
- `gnalloy.org/benchmarks/external/internal/benchh3` (`benchh3`)
- `gnalloy.org/benchmarks/external/internal/benchhttp` (`benchhttp`)
- `gnalloy.org/benchmarks/external/internal/benchtls` (`benchtls`)
- `gnalloy.org/benchmarks/microbench` (`microbench`)
- `gnalloy.org/benchmarks/parity` (`parity`)

## Gnalloy Dependencies

- `gnalloy.org/codec-http3`
- `gnalloy.org/gnalloy`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

## Common Integration Pattern
- Configure scenario, protocol, payload size, concurrency, warmup count, measured repetitions, timeout, and report output before comparing results.
- Report throughput and latency separately. Do not compare numbers collected on different hosts as a ranking.

## Current Public Entry Points

The generated API reference lists the full public surface. Common constructors or option types currently include:
- `var ErrInvalidRunner = errors.New("gnalloy/benchmarks/benchdiff: invalid runner") ...`
- `type Config struct{ ... }`
- `type Config struct{ ... }`
- `const ProtocolHTTP1 = "http1"`
- `type Config struct{ ... }`
- `const DefaultServerName = "gnalloy.local"`
- `var ErrInvalidSpec = errors.New("gnalloy/benchmarks/parity: invalid spec") ...`
- `type ExternalHarnessOptions struct{ ... }`

## Verification

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -count=1
```

For pressure tests, assemble this module with the relevant transport, codec, and handler stack and run the scenario from `gnalloy.org/benchmarks` or `gnalloy.org/examples`. Keep host, operating system, payload, concurrency, warmup, and repetitions in the report.

## Caveats
- This repository is intentionally narrow. Cross-module behavior should be assembled in applications, recipes, examples, or benchmark harnesses.
- Public APIs should remain Go-native and explicit; avoid runtime scanning, hidden global registries, and reflection-heavy behavior in hot paths.
- Treat network input as untrusted. Configure parser limits and return typed errors instead of panics.
- Keep benchmark claims tied to a concrete host, operating system, protocol, payload, concurrency, warmup, and repetition count.
