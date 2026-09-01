# 案例

[English](examples.md) | [文档索引](README.zh-CN.md)

## 案例 1：将模块加入应用

```bash
mkdir gnalloy-app && cd gnalloy-app
go mod init example.com/gnalloy-app
go get gnalloy.org/benchmarks@dev
go doc gnalloy.org/benchmarks
```

## 案例 2：查看当前包

当前源码树暴露这些 package 导入路径：
- `gnalloy.org/benchmarks/benchdiff`
- `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff`
- `gnalloy.org/benchmarks/external/internal/benchh2`
- `gnalloy.org/benchmarks/external/internal/benchh3`
- `gnalloy.org/benchmarks/external/internal/benchhttp`
- `gnalloy.org/benchmarks/external/internal/benchtls`
- `gnalloy.org/benchmarks/microbench`
- `gnalloy.org/benchmarks/parity`

按需要的行为对对应 package 执行 `go doc`：

```bash
go doc gnalloy.org/benchmarks/benchdiff
go doc gnalloy.org/benchmarks/cmd/gnalloy-benchdiff
go doc gnalloy.org/benchmarks/external/internal/benchh2
go doc gnalloy.org/benchmarks/external/internal/benchh3
go doc gnalloy.org/benchmarks/external/internal/benchhttp
go doc gnalloy.org/benchmarks/external/internal/benchtls
go doc gnalloy.org/benchmarks/microbench
go doc gnalloy.org/benchmarks/parity
```

精选当前导出入口：
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
- `const ProtocolHTTP1 = "http1"`
- `func RequestBytes(host string) []byte`
- `func ResponseBody(size int) []byte`
- `func ResponseBytes(payload int) []byte`

## 案例 3：将可执行测试作为行为示例

仓库测试是受支持行为的可执行示例。先从下面的精选名称开始，再阅读对应 `_test.go` 文件中的完整 setup 和断言。完整发现列表见 [测试与性能](testing.zh-CN.md)。

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run Test -count=1
```

精选当前 test、benchmark、fuzz 与 example 入口：
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

## 案例 4：跨模块装配

本模块的直接 Gnalloy 依赖：
- `gnalloy.org/codec-http3`
- `gnalloy.org/gnalloy`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

装配说明：
- benchmark 工具用于测量装配后的 Gnalloy 栈，而不是孤立文档示例。
- 场景定义必须可复现，保留 host、OS、protocol、payload、concurrency、warmup、repetitions、throughput 与 p99 latency。
- 硬件和软件矩阵未记录时，不要跨主机比较结果。

## 案例 5：压测 Harness

持续负载测试时，如果该模块参与网络流量路径，将它接入 `gnalloy.org/benchmarks` 的场景，或接入 `gnalloy.org/examples` 的可运行客户端。报告中记录 host、OS、CPU、Go version、protocol、payload、concurrency、warmup、repetitions、throughput 和 p99 latency。
