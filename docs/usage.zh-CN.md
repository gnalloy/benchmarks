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

## 完整对标矩阵

先构建外部 harness，再从包含 `benchmarks` 目录的 workspace 根目录运行内置矩阵：

```bash
(cd benchmarks && GOWORK=off GOTOOLCHAIN=local go build -o external/bin/gnalloy-parity ./cmd/gnalloy-parity)
./benchmarks/external/bin/gnalloy-parity -matrix linux-full -strict-external -format markdown -out benchmarks/reports/linux-full-parity.md -timeout 6h
```

Linux 矩阵覆盖 TCP、UDP、QUIC stream、HTTP/1.1、HTTP/2、HTTP/3、HTTPS TLS 1.1/1.2/1.3、Gnalloy、Netty、gnet、fasthttp 与 netpoll。非法或不支持的组合会记录为 skipped scenario，并保留明确原因。可执行场景向 harness 传入 15 分钟 timeout，并使用 16 分钟外层场景保护。

## Linux TCP 专项测试

构建 `external/bin/gnalloy-bench` 和 `external/bin/netpoll-bench` 后，在同一台 Linux 主机执行：

```bash
./scripts/run-linux-tcp-optimized.sh
```

脚本任意时刻只运行一个进程，在五轮采样中交替 Gnalloy/netpoll 的先后顺序，并记录主机信息、参数、二进制哈希、吞吐和延迟。payload、负载、冷却时间、可执行文件、输出路径和进程优先级均可通过环境变量设置。脚本不会修改 CPU governor 或停止主机服务。已验证的 16KiB 结果和 syscall 分析见 `reports/linux-tcp-read-write-optimization.md`。

## Linux UDP 专项测试

构建 `external/bin/gnalloy-bench`、`external/bin/gnet-bench`、`external/bin/netty-bench.jar` 和 `external/bin/udp-load` 后，在同一台 Linux 主机用统一客户端测试全部 UDP 服务端：

```bash
./scripts/run-linux-udp-common-client.sh
TARGET_RATE=60000 ./scripts/run-linux-udp-common-client.sh
```

默认模式测量饱和吞吐。`TARGET_RATE` 为等负载延迟测试设置共享的聚合请求速率。固定速率模式下，`latency` 从每个请求的计划发送时刻开始统计，`scheduleDelay` 表示客户端实际开始发送的滞后时间，`roundTripLatency` 表示从实际发送尝试到收到回显的耗时。这组分解既保留调度积压，又能将它与网络和服务端处理耗时区分。脚本将服务端和客户端固定到互不重叠的物理核，通过 Linux sysfs 校验 package/core ID，轮换框架顺序，严格串行执行案例，并记录 CPU 拓扑、governor、参数和二进制哈希。仅当主机无法提供物理核隔离集合时才设置 `REQUIRE_PHYSICAL_CPU_ISOLATION=0`，并且不得把该结果与隔离测试直接比较。Netpoll 与 fasthttp 的 benchmark harness 没有可比的 UDP 服务端，因此保持排除。

如需通过内网固定主机角色验证 UDP，在 Debian `172.16.8.172` 上运行统一客户端，并在 Debian `172.16.8.171` 上严格串行启动各服务端。客户端 checkout 默认使用 `/opt/test/gnalloy/benchmarks-cross-host`，服务端 checkout 默认使用 `/opt/test/gnalloy/benchmarks-e481c9a`。在客户端构建 `external/bin/udp-load`，在服务端构建全部服务端二进制，然后从控制工作站执行：

```powershell
.\scripts\run-cross-host-udp-common-client.ps1
.\scripts\run-cross-host-udp-common-client.ps1 -TargetRate 60000
```

诊断时仅选择 Gnalloy，并且每轮只采集一种探针。文件先写入服务端 checkout，再下载到 `reports/raw/cross-host-profiles`；带探针的结果不得与正常性能结果直接比较：

```powershell
.\scripts\run-cross-host-udp-common-client.ps1 -Frameworks gnalloy -Payloads 128 -Repetitions 1 -CaptureCPUProfile
.\scripts\run-cross-host-udp-common-client.ps1 -Frameworks gnalloy -Payloads 128 -Repetitions 1 -TargetRate 60000 -CaptureRuntimeTrace
```

客户端默认绑定 `0,1,2,4`，避开 CPU 11 上的网卡 IRQ，并使用四个不同物理核。runner 会保存两端每个 CPU 的 governor，测试期间切换为 `performance`，清理阶段恢复全部原值。跨机结果包含两端网卡、交换链路和两端操作系统网络栈的成本，必须与同机物理核隔离结果分开报告。

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
