# Testing and Performance

[简体中文](testing.zh-CN.md) | [Docs Index](README.md)

## Required Checks

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
gofmt -l .
git diff --check
```

## Focused Behavior Checks

Run focused tests while working on a small behavior change:

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run 'TestName' -count=1
```

## Discovered Test Entry Points

This inventory is generated from the current `_test.go` files in this repository. It is intentionally complete so documentation review can catch stale test, benchmark, fuzz, and example coverage when code changes.

Total discovered entry points: 116.

### Tests (116)
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
- `TestHTTP2MatrixSpecLoads`
- `TestHTTP3MatrixSpecLoads`
- `TestHTTPS1ALPNMatrixSpecLoads`
- `TestInspectExternalHarnessesChecksExpandedCommand`
- `TestLatencySamplingPredicate`
- `TestLinuxHTTP1MatrixExternalHarnessesCanPassStrictGateWithRepoArtifacts`
- `TestLinuxHTTP2MatrixSpecLoads`
- `TestLinuxHTTP3MatrixSpecLoads`
- `TestLinuxHTTPS1ALPNMatrixSpecLoads`
- `TestLinuxTLSVersionMatrixSpecLoads`
- `TestLinuxUDPEchoMatrixSpecLoads`
- `TestLoadSpecAllowsSkippedExternalScenarioWithoutCommand`
- `TestLoadSpecRejectsInvalidScenario`
- `TestLoadSpecRejectsNegativeSampling`
- `TestLookupHotPathSuiteBuildsStableBenchdiffInputs`
- `TestLookupRejectsUnknownSuite`
- `TestParseConfig`
- `TestParseConfigKeepsMinimumAutoReadBufferSize`
- `TestParseConfigRejectsCipherSuitesForTLS13`
- `TestParseConfigRejectsFixedBuffersWithoutMmap`
- `TestParseConfigRejectsHTTP1BenchmarkBypass`
- `TestParseConfigRejectsHTTP2TLS11`
- `TestParseConfigRejectsHTTP3TLS12`
- `TestParseConfigRejectsInsecureCipherSuiteByDefault`
- `TestParseConfigRejectsInvalidBackend`
- `TestParseConfigRejectsInvalidTLSVersion`
- `TestParseConfigRejectsMmapBlockSmallerThanReadBuffer`
- `TestParseConfigRejectsMmapSizeOverflow`
- `TestParseConfigRejectsNegativeLatencySampleRate`
- `TestParseConfigRejectsNegativeReadBufferSize`
- `TestParseConfigRejectsNegativeWarmupMessages`
- `TestParseConfigRejectsNegativeWorkers`
- `TestParseConfigRejectsUnsupportedProtocol`
- `TestParseConfigResolvesAutoReadBufferSize`
- `TestParseConfigResolvesAutoWorkers`
- `TestParseConfigResolvesNativePerformanceFlags`
- `TestParseConfigSupportsCipherSuites`
- `TestParseConfigSupportsHTTP2Family`
- `TestParseConfigSupportsHTTP3`
- `TestParseConfigSupportsHTTPS1ALPN`
- `TestParseConfigSupportsInsecureCipherSuiteOptIn`
- `TestParseConfigSupportsTLSVersions`
- `TestParseConfigSupportsUDPEcho`
- `TestParseGoBenchOutputTracksPackage`
- `TestParseScenarioStats`
- `TestParseScenarioStatsParsesJavaDuration`
- `TestPathCommandCandidatesAddsWindowsExecutableSuffix`
- `TestRequestHeaderBlockUsesStaticHPACKFields`
- `TestResolveBenchmarkSelectionKeepsExplicitOverrides`
- `TestResolveBenchmarkSelectionRejectsUnknownSuite`
- `TestResolveBenchmarkSelectionUsesSuiteDefaults`
- `TestRunBenchmarkHTTP1`
- `TestRunBenchmarkHTTP2`
- `TestRunBenchmarkHTTPS1`
- `TestRunBenchmarkHTTPS2ALPN`
- `TestRunBenchmarkHTTPS2SustainedLoad`
- `TestRunBenchmarkRejectsInvalidConfig`
- `TestRunBenchmarkReportsUnsupportedPlatform`
- `TestRunBenchmarkTCPEcho`
- `TestRunBenchmarkUDPEcho`
- `TestRunLoadHTTP1`
- `TestRunLoadHTTP2Cleartext`
- `TestRunLoadHTTP2TLSALPN`
- `TestRunLoadHTTP3QUIC`
- `TestRunLoadHTTPS1ALPN`
- `TestRunLoadTimeoutClosesBlockedClients`
- `TestRunnerCapturesCommandOutput`
- `TestRunnerDryRunProducesSkippedResults`
- `TestRunnerExpandsScenarioVariables`
- `TestRunnerRepeatsScenarioAndDropsWarmupOutput`
- `TestRunnerSkipsScenarioMarkedSkip`
- `TestSelfSignedCertificatesDeduplicatesAlgorithms`
- `TestSelfSignedCertificateSupportsIPName`
- `TestSelfSignedCertificateSupportsRSA`
- `TestSelfSignedCertificateUsesDefaultServerName`
- `TestServerTLSConfigUsesTLS13`
- `TestSummarizeLatencySamples`
- `TestTCPMatrixIncludesOptimizedIOUringScenario`
- `TestTCPMatrixSpecLoads`
- `TestTLSConfigUsesCipherSuites`
- `TestTLSConfigUsesSelectedVersion`
- `TestTLSVersionMatrixSpecLoads`
- `TestUDPEchoMatrixSpecLoads`
- `TestValidateExternalHarnessesAcceptsReadyCommand`
- `TestValidateExternalHarnessesChecksGnalloyParityHarness`
- `TestValidateExternalHarnessesChecksJavaJarArgument`
- `TestValidateExternalHarnessesRejectsSkippedScenario`
- `TestWindowsHTTP1MatrixExternalHarnessesCanPassStrictGateWithRepoArtifacts`
- `TestWindowsTCPMatrixExternalHarnessesCanPassStrictGateWithRepoArtifacts`
- `TestWindowsTCPMatrixSpecLoads`
- `TestWriteBenchmarkResult`
- `TestWriteHTTP3BenchmarkResultUsesRFC9000Backend`
- `TestWriteJSONReport`
- `TestWriteMarkdownIncludesComparisonRows`
- `TestWriteMarkdownReportIncludesMachineAndScenario`

### Benchmarks (0)
- No Benchmark functions are currently declared.

### Fuzz Targets (0)
- No Fuzz targets are currently declared.

### Examples (0)
- No Example functions are currently declared.

## Race Checks

```bash
GOWORK=off GOTOOLCHAIN=local go test -race ./... -count=1
```

Race checks are most valuable for core, transport, handler, resolver, observability, examples, and benchmark modules. Platform-specific transports may require native host capabilities.

## Benchmarks

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -benchtime=1s -count=5
```

Report `ns/op`, `B/op`, `allocs/op`, throughput, and p99 latency separately. Include host and OS details with every result.

## Pressure Testing

Pressure tests should run against a realistic assembled stack. Use `gnalloy.org/benchmarks` for repeatable matrices and `gnalloy.org/examples` for runnable clients. Keep warmup and measurement phases separate.

## CI

The repository validation workflow runs formatting, tests, and vet on Linux, macOS, and Windows for pushes and pull requests.
