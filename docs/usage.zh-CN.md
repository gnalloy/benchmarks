# 用法

[English](usage.md) | [文档索引](README.zh-CN.md)

## 要求

- Go 1.25 或更新版本，并与 module 的 `go` 指令一致。
- 由 Gnalloy 应用、recipe、example 或 benchmark harness 负责生命周期与部署配置。
- 独立模块复验应设置 `GOWORK=off`，确保通过已发布依赖图测试。

## 安装
```bash
go get gnalloy.org/benchmarks@dev
```

## 导入
```go
import "gnalloy.org/benchmarks"
```

## 集成模式
- 比较结果前先配置 scenario、protocol、payload size、concurrency、warmup、measured repetitions、timeout 与 report 输出。
- throughput 与 latency 必须分开报告；不同主机采集的数据不能直接当成排名。

## API 选择

通过 API 清单选择当前协议路径需要的具体构造函数或 option 类型：

```bash
go doc gnalloy.org/benchmarks
```

当前常用入口：
- `var ErrInvalidRunner = errors.New("gnalloy/benchmarks/benchdiff: invalid runner") ...`
- `type Config struct{ ... }`
- `type Config struct{ ... }`
- `const ProtocolHTTP1 = "http1"`
- `type Config struct{ ... }`
- `const DefaultServerName = "gnalloy.local"`
- `var ErrInvalidSpec = errors.New("gnalloy/benchmarks/parity: invalid spec") ...`
- `type ExternalHarnessOptions struct{ ... }`

## 跨模块装配

多个 Gnalloy 仓库一起开发时，在自己选择的 workspace 中创建本地 `go.work` 文件。不要把应用本地 `replace` 指令提交到发布用 library module，除非它是明确的临时变更且不会进入提交。

## 错误处理

网络输入、对端行为、平台能力和超时失败都必须作为普通错误处理。不要用 panic 恢复协议正确性。返回或传播模块错误，并在所有权要求时关闭受影响的 Channel。
