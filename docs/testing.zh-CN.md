# 测试与性能

[English](testing.md) | [文档索引](README.zh-CN.md)

## 必跑检查

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
gofmt -l .
git diff --check
```

## 聚焦行为检查

处理小范围行为变更时先跑聚焦测试：

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run 'TestName' -count=1
```

## 已发现测试入口

本清单从当前仓库的 `_test.go` 文件生成。这里刻意保持完整，用于在代码变化时发现 test、benchmark、fuzz 与 example 覆盖说明是否过期。

已发现入口总数：116。

### Tests（116）
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
- `TestHTTP1RawHandlerBatchesResponsesFromOneRead`
- `TestHTTP1RawHandlerHandlesFragmentedRequest`
- `TestHTTP1RawHandlerWritesFixedResponse`
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
- `TestParseConfigRejectsHTTP1ModeForHTTP2`
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
- `TestParseConfigSupportsHTTP1RawMode`
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
- `TestResponseBytesUsesRequestedPayload`
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
- `TestServerStateCountsSplitRequests`
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

### Benchmarks（0）
- 当前没有声明 Benchmark 函数。

### Fuzz Targets（0）
- 当前没有声明 Fuzz 入口。

### Examples（0）
- 当前没有声明 Example 函数。

## Race 检查

```bash
GOWORK=off GOTOOLCHAIN=local go test -race ./... -count=1
```

Race 检查对 core、transport、handler、resolver、observability、examples 和 benchmark 模块最有价值。平台相关 transport 可能需要原生主机能力。

## 基准测试

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -benchtime=1s -count=5
```

`ns/op`、`B/op`、`allocs/op`、throughput 和 p99 latency 要分开报告。每份结果都要包含 host 与 OS 信息。

## 压测

压测应针对真实装配栈运行。使用 `gnalloy.org/benchmarks` 维护可重复矩阵，使用 `gnalloy.org/examples` 运行客户端。warmup 和 measurement 阶段必须分离。

## CI

仓库 validation workflow 会在 Linux、macOS 与 Windows 上为 push 和 pull request 运行格式检查、测试和 vet。
