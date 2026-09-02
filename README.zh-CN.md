# benchmarks

[English](README.md) | [文档](docs/README.zh-CN.md)

Gnalloy 模块化网络栈的基准、性能对标与压测工具。

该模块包含 benchmark、parity 与 pressure-test 工具，用于测量和回归防护，不作为生产运行时依赖。

## 状态

- 导入路径：`gnalloy.org/benchmarks`
- 仓库：`github.com/gnalloy/benchmarks`
- 默认分支：`dev`
- 预览安装：`go get gnalloy.org/benchmarks@dev`
- 许可证：Apache-2.0

## 安装
```bash
go get gnalloy.org/benchmarks@dev
go doc gnalloy.org/benchmarks
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
```

## 文档
- [概览](docs/overview.zh-CN.md) ([English](docs/overview.md))
- [用法](docs/usage.zh-CN.md) ([English](docs/usage.md))
- [案例](docs/examples.zh-CN.md) ([English](docs/examples.md))
- [配置说明](docs/configuration.zh-CN.md) ([English](docs/configuration.md))
- [测试与性能](docs/testing.zh-CN.md) ([English](docs/testing.md))
- [API 参考](docs/api.zh-CN.md) ([English](docs/api.md))
- [注意事项与备注](docs/notes.zh-CN.md) ([English](docs/notes.md))
- [ADR-001 模块边界](docs/decisions/0001-module-boundary.zh-CN.md) ([English](docs/decisions/0001-module-boundary.md))

## 模块边界

本仓库负责：Gnalloy 模块化网络栈的基准、性能对标与压测工具。

它不吸收相邻模块职责。核心基础能力保留在 `gnalloy.org/gnalloy`；协议 codec、transport、handler、resolver、examples 与 benchmarks 分别由独立仓库负责。

## 包结构
- `gnalloy.org/benchmarks/benchdiff`（`benchdiff`）
- `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff`（`main`）
- `gnalloy.org/benchmarks/cmd/gnalloy-parity`（`main`）
- `gnalloy.org/benchmarks/external/internal/benchh2`（`benchh2`）
- `gnalloy.org/benchmarks/external/internal/benchh3`（`benchh3`）
- `gnalloy.org/benchmarks/external/internal/benchtls`（`benchtls`）
- `gnalloy.org/benchmarks/microbench`（`microbench`）
- `gnalloy.org/benchmarks/parity`（`parity`）

## Gnalloy 依赖

- `gnalloy.org/codec-http1`
- `gnalloy.org/codec-http3`
- `gnalloy.org/gnalloy`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

## 常见集成方式
- 比较结果前先配置 scenario、protocol、payload size、concurrency、warmup、measured repetitions、timeout 与 report 输出。
- throughput 与 latency 必须分开报告；不同主机采集的数据不能直接当成排名。
- 使用 `cmd/gnalloy-parity -matrix linux-full` 或 `-matrix windows-full` 生成完整 TCP、UDP、QUIC、HTTP/1、HTTP/2、HTTP/3 与 TLS 对标矩阵。

## 当前公共入口

生成的 API 参考列出了完整公共面。当前常用构造函数或 option 类型包括：
- `var ErrInvalidRunner = errors.New("gnalloy/benchmarks/benchdiff: invalid runner") ...`
- `type Config struct{ ... }`
- `type Config struct{ ... }`
- `const ProtocolHTTP1 = "http1"`
- `type Config struct{ ... }`
- `const DefaultServerName = "gnalloy.local"`
- `var ErrInvalidSpec = errors.New("gnalloy/benchmarks/parity: invalid spec") ...`
- `type ExternalHarnessOptions struct{ ... }`

## 验证

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -count=1
```

压测时，将该模块和相应 transport、codec、handler 栈装配后，使用 `gnalloy.org/benchmarks` 或 `gnalloy.org/examples` 中的场景运行。报告必须保留主机、操作系统、payload、并发度、warmup 和 repetition。

## 注意事项
- 本仓库保持窄边界。跨模块行为应在应用、recipes、examples 或 benchmark harness 中装配。
- 公共 API 必须保持 Go 原生和显式；热路径避免运行时扫描、隐藏全局注册表和重反射。
- 网络输入一律视为不可信。配置解析上限，返回类型化错误，不使用 panic 处理输入错误。
- 性能结论必须绑定具体主机、操作系统、协议、payload、并发度、warmup 和 repetition。
