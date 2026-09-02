# API 参考

[English](api.md) | [文档索引](README.zh-CN.md)

本清单由本仓库 package 的 `go doc -short` 生成，用于快速查看公共面。精确语义以源码和测试为准。

## 包

### `gnalloy.org/benchmarks/benchdiff`

包名：`benchdiff`

```text
var ErrInvalidRunner = errors.New("gnalloy/benchmarks/benchdiff: invalid runner") ...
func CompareSamples(baseSamples []Sample, candidateSamples []Sample) ([]Comparison, []Missing)
func WriteMarkdown(w io.Writer, report Report) error
type Comparison struct{ ... }
type Missing struct{ ... }
type Report struct{ ... }
    func NewReport(baseLabel string, candidateLabel string, command []string, baseOutput string, ...) Report
type Runner struct{ ... }
type Sample struct{ ... }
    func ParseGoBenchOutput(output string) []Sample
type Summary struct{ ... }
```

### `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff`

包名：`main`

```text
(no exported API reported by go doc)
```

### `gnalloy.org/benchmarks/external/internal/benchh3`

包名：`benchh3`

```text
func ResponseBody(size int) []byte
type Config struct{ ... }
type LatencySummary struct{ ... }
type Result struct{ ... }
    func RunLoad(parent context.Context, cfg Config) (Result, error)
```

### `gnalloy.org/benchmarks/external/internal/benchtls`

包名：`benchtls`

```text
const DefaultServerName = "gnalloy.local"
func SelfSignedCertificate(serverName string) (tls.Certificate, error)
func SelfSignedCertificateWithAlgorithm(serverName string, algorithm CertificateKeyAlgorithm) (tls.Certificate, error)
func SelfSignedCertificates(serverName string, algorithms ...CertificateKeyAlgorithm) ([]tls.Certificate, error)
type CertificateKeyAlgorithm uint8
    const CertificateKeyECDSA CertificateKeyAlgorithm = iota + 1 ...
```

### `gnalloy.org/benchmarks/microbench`

包名：`microbench`

```text
type Scenario struct{ ... }
type Suite struct{ ... }
    func Lookup(name string) (Suite, bool)
    func Suites() []Suite
```

### `gnalloy.org/benchmarks/parity`

包名：`parity`

```text
var ErrInvalidSpec = errors.New("gnalloy/benchmarks/parity: invalid spec") ...
func ValidateExternalHarnesses(spec Spec, options ExternalHarnessOptions) error
func WriteReport(w io.Writer, report Report, format Format) error
type BenchmarkMetric struct{ ... }
    func ParseBenchmarkMetrics(output string) []BenchmarkMetric
type Duration time.Duration
type ExternalHarnessError struct{ ... }
type ExternalHarnessIssue struct{ ... }
type ExternalHarnessOptions struct{ ... }
type ExternalHarnessReport struct{ ... }
    func InspectExternalHarnesses(spec Spec, options ExternalHarnessOptions) (ExternalHarnessReport, error)
type Format string
    const FormatMarkdown Format = "markdown" ...
type Machine struct{ ... }
    func DetectMachine() Machine
type Report struct{ ... }
type Runner struct{ ... }
type Scenario struct{ ... }
type ScenarioResult struct{ ... }
type ScenarioSample struct{ ... }
type ScenarioStats struct{ ... }
    func ParseScenarioStats(output string) []ScenarioStats
type Spec struct{ ... }
    func LoadSpec(r io.Reader) (Spec, error)
type StatsSummary struct{ ... }
```
