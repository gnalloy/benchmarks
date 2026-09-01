# API Reference

[简体中文](api.zh-CN.md) | [Docs Index](README.md)

This inventory is generated from `go doc -short` for the packages in this repository. It is a quick public-surface map; source files and tests remain the authority for exact semantics.

## Packages

### `gnalloy.org/benchmarks/benchdiff`

Package name: `benchdiff`

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

Package name: `main`

```text
(no exported API reported by go doc)
```

### `gnalloy.org/benchmarks/external/internal/benchh2`

Package name: `benchh2`

```text
func ResponseBody(size int) []byte
type Config struct{ ... }
type LatencySummary struct{ ... }
type Result struct{ ... }
    func RunLoad(parent context.Context, cfg Config) (Result, error)
```

### `gnalloy.org/benchmarks/external/internal/benchh3`

Package name: `benchh3`

```text
func ResponseBody(size int) []byte
type Config struct{ ... }
type LatencySummary struct{ ... }
type Result struct{ ... }
    func RunLoad(parent context.Context, cfg Config) (Result, error)
```

### `gnalloy.org/benchmarks/external/internal/benchhttp`

Package name: `benchhttp`

```text
const ProtocolHTTP1 = "http1"
func RequestBytes(host string) []byte
func ResponseBody(size int) []byte
func ResponseBytes(payload int) []byte
type Config struct{ ... }
type LatencySummary struct{ ... }
type Result struct{ ... }
    func RunLoad(parent context.Context, cfg Config) (Result, error)
type ServerState struct{ ... }
```

### `gnalloy.org/benchmarks/external/internal/benchtls`

Package name: `benchtls`

```text
const DefaultServerName = "gnalloy.local"
func SelfSignedCertificate(serverName string) (tls.Certificate, error)
func SelfSignedCertificateWithAlgorithm(serverName string, algorithm CertificateKeyAlgorithm) (tls.Certificate, error)
func SelfSignedCertificates(serverName string, algorithms ...CertificateKeyAlgorithm) ([]tls.Certificate, error)
type CertificateKeyAlgorithm uint8
    const CertificateKeyECDSA CertificateKeyAlgorithm = iota + 1 ...
```

### `gnalloy.org/benchmarks/microbench`

Package name: `microbench`

```text
type Scenario struct{ ... }
type Suite struct{ ... }
    func Lookup(name string) (Suite, bool)
    func Suites() []Suite
```

### `gnalloy.org/benchmarks/parity`

Package name: `parity`

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
