# Examples

[简体中文](examples.zh-CN.md) | [Docs Index](README.md)

## Example 1: Add the Module to an Application

```bash
mkdir gnalloy-app && cd gnalloy-app
go mod init example.com/gnalloy-app
go get gnalloy.org/benchmarks@dev
go doc gnalloy.org/benchmarks
```

## Example 2: Inspect Current Packages

The current source tree exposes these package import paths:
- `gnalloy.org/benchmarks/benchdiff`
- `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff`
- `gnalloy.org/benchmarks/external/internal/benchh3`
- `gnalloy.org/benchmarks/external/internal/benchtls`
- `gnalloy.org/benchmarks/microbench`
- `gnalloy.org/benchmarks/parity`

Use `go doc` against the package that matches the behavior you need:

```bash
go doc gnalloy.org/benchmarks/benchdiff
go doc gnalloy.org/benchmarks/cmd/gnalloy-benchdiff
go doc gnalloy.org/benchmarks/external/internal/benchh3
go doc gnalloy.org/benchmarks/external/internal/benchtls
go doc gnalloy.org/benchmarks/microbench
go doc gnalloy.org/benchmarks/parity
```

Selected current exported entry points:
- `var ErrInvalidRunner = errors.New("gnalloy/benchmarks/benchdiff: invalid runner") ...`
- `func CompareSamples(baseSamples []Sample, candidateSamples []Sample) ([]Comparison, []Missing)`
- `func WriteMarkdown(w io.Writer, report Report) error`
- `type Comparison struct{ ... }`
- `type Missing struct{ ... }`
- `type Report struct{ ... }`
- `func ResponseBody(size int) []byte`
- `type Config struct{ ... }`
- `type LatencySummary struct{ ... }`
- `type Result struct{ ... }`
- `func ResponseBody(size int) []byte`
- `type Config struct{ ... }`
- `type LatencySummary struct{ ... }`
- `type Result struct{ ... }`

## Example 3: Use Executable Tests as Behavioral Examples

Repository tests are executable examples of supported behavior. Start with the selected names below, then read the matching `_test.go` files for complete setup and assertions. See [Testing and Performance](testing.md) for the complete discovered list.

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run Test -count=1
```

Selected current test, benchmark, fuzz, and example entry points:
- `TestALPNProtocolsTrimsEmptyItems`
- `TestAverageLatencyNanosKeepsPositiveFloor`
- `TestAverageLatencyNanosUsesWindowMean`
- `TestBaselineExternalHarnessesCanPassStrictGateWithRepoArtifacts`
- `TestBaselineGnalloyTCPEchoUsesSameLoadModel`
- `TestBaselineNettyTCPEchoUsesNativeEpoll`
- `TestBaselineSpecLoads`
- `TestCertificateAlgorithmsMatchCipherSuiteAuth`
- `TestCompareSamplesUsesMedianAndReportsMissing`
- `TestDefaultWorkerCountCapsLinuxEpoll`
- `TestDefaultWorkerCountCapsLinuxIOUring`
- `TestDefaultWorkerCountCapsWindowsIOCP`
- `TestDefaultWorkerCountKeepsNonIOCPParallelism`
- `TestDefaultWorkerCountNormalizesInvalidCPUCount`
- `TestElapsedLatencyNanosIsPositive`

## Example 4: Cross-Module Assembly

Direct Gnalloy dependencies for this module:
- `gnalloy.org/codec-http1`
- `gnalloy.org/codec-http3`
- `gnalloy.org/gnalloy`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

Assembly guidance:
- Use benchmark tooling to measure assembled Gnalloy stacks, not isolated documentation examples.
- Keep scenario definitions reproducible with host, OS, protocol, payload, concurrency, warmup, repetitions, throughput, and p99 latency.
- Do not compare results across hosts unless the hardware and software matrix is recorded.

## Example 5: Pressure-Test Harness

For sustained load, wire this module into a scenario under `gnalloy.org/benchmarks` or a runnable client under `gnalloy.org/examples` when the module participates in network traffic. Record host, OS, CPU, Go version, protocol, payload, concurrency, warmup, repetitions, throughput, and p99 latency in the report.
