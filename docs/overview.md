# Overview

[简体中文](overview.zh-CN.md) | [Docs Index](README.md)

## Purpose

Benchmark, parity, and pressure-test tooling for the Gnalloy modular networking stack.

This module contains benchmark, parity, and pressure-test tooling. It is for measurement and regression protection, not for production runtime dependencies.

## Repository Identity

- Module path: `gnalloy.org/benchmarks`
- GitHub repository: `github.com/gnalloy/benchmarks`
- Default branch: `dev`
- License: Apache-2.0

## Package Map
- `gnalloy.org/benchmarks/benchdiff` (`benchdiff`)
- `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff` (`main`)
- `gnalloy.org/benchmarks/external/internal/benchh2` (`benchh2`)
- `gnalloy.org/benchmarks/external/internal/benchh3` (`benchh3`)
- `gnalloy.org/benchmarks/external/internal/benchhttp` (`benchhttp`)
- `gnalloy.org/benchmarks/external/internal/benchtls` (`benchtls`)
- `gnalloy.org/benchmarks/microbench` (`microbench`)
- `gnalloy.org/benchmarks/parity` (`parity`)

## Direct Gnalloy Dependencies

- `gnalloy.org/codec-http3`
- `gnalloy.org/gnalloy`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

## Direct Dependents in the Current Repository Set

- `gnalloy.org/examples`

## Architecture Position

Gnalloy keeps the core small and dependency-light. This repository is a replaceable module around one responsibility, connected through explicit Go packages instead of runtime discovery.

## Compatibility

The public import path is `gnalloy.org/benchmarks`. Until the first stable tag is published, use `@dev` or an explicit pseudo-version selected by your dependency policy.
