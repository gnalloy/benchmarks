# Linux 跨机 HTTP/1 与 TLS 对比（2026-09-04）

## 目的

在 171 服务端、172 客户端上串行对比 Gnalloy、Hertz、fasthttp 与 Netty 的 HTTP/1.1 和 HTTPS/1.1。Gnalloy 服务端使用正式 TCP transport、Channel Pipeline、HTTP/1 codec、业务 handler；HTTPS 再经过正式 TLS handler。benchmark 不包含自写 HTTP 帧解析、固定响应旁路或复制 codec 逻辑。

gnet 和 CloudWeGo netpoll 不提供可直接用于本矩阵的 HTTP codec。按照禁止 benchmark 自行补协议的约束，两者明确标记为 N/A，不把 TCP benchmark 冒充 HTTP benchmark。

## 环境与方法

- 服务端：`172.16.8.171`，Debian 13，Linux `6.12.48+deb13-amd64`，Intel Xeon E3-1535M v5 @ 2.90GHz。
- 客户端：`172.16.8.172`，Debian 13，Linux `6.12.48+deb13-amd64`。
- 服务端固定 CPU `0,1,2,3`；客户端固定 CPU `0,1,2,4`；两端 `GOMAXPROCS=4`，测试期间使用 performance governor。
- 64 个连接，每连接 5,000 个测量请求和 500 个预热请求；每 64 个请求采样一次延迟。
- 每个案例执行 5 轮并取中位数；框架轮换先后顺序，任意时刻只运行一个服务端和一个客户端案例。
- 固定速率矩阵使用 70,000 ops/s。吞吐接近目标值只表示服务端能承载该速率，延迟结论使用实际发送后的 RTT，不使用包含客户端排队的 total latency。
- TLS 1.1 使用 `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA`，TLS 1.2 使用 `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`，TLS 1.3 使用 Go 标准套件；ALPN 为 `http/1.1`。
- 当前 Gnalloy 二进制 SHA-256：`acc4a5ead2a53ada9e94431ff1c281cd1f606d871597bdc5a74eaaf062ca94f0`。
- Hertz 二进制 SHA-256：`45a002d95efe434240a4e67dd7559bb36490037dd911b766f4eb2389eb9c3af6`。
- fasthttp 二进制 SHA-256：`539f5b1072956f5075ab7fbf2d9221a6041e759f360843876e82eaf957c0c297`。
- Netty JAR SHA-256：`4df099cb73622543e621e246cd79bb4b3c154f0b09f4bf2b8804df67a4b65aa5`。
- 公共客户端 SHA-256：`43cd55fa74eb3aca53d1d91c16eed8848c606ceaf95c55ee041034fc0d406106`。

## 明文 HTTP/1.1

### 饱和吞吐

| Payload | 服务端 | 吞吐 (ops/s) | RTT P50 | RTT P99 |
|---:|---|---:|---:|---:|
| 128B | Gnalloy | 150,439.64 | 0.390ms | 0.784ms |
| 128B | Hertz | 144,481.61 | 0.437ms | 0.842ms |
| 128B | fasthttp | 149,356.18 | 0.408ms | 0.817ms |
| 128B | Netty | 123,365.51 | 0.456ms | 4.389ms |
| 1KiB | Gnalloy | 103,542.15 | 0.603ms | 0.986ms |
| 1KiB | Hertz | 98,375.23 | 0.616ms | 1.013ms |
| 1KiB | fasthttp | 103,289.17 | 0.602ms | 1.046ms |
| 1KiB | Netty | 94,474.01 | 0.602ms | 4.493ms |

Gnalloy 在两个 payload 的饱和吞吐和 RTT P99 均为本矩阵最优。128B 吞吐相对 fasthttp 高约 `0.7%`，1KiB 高约 `0.2%`；这两个小差值应视为本机矩阵结果，不能外推为普遍优势。

### 70k offered-rate

| Payload | 服务端 | 实际吞吐 (ops/s) | RTT P50 | RTT P99 |
|---:|---|---:|---:|---:|
| 128B | Gnalloy | 69,995.07 | 0.358ms | 0.679ms |
| 128B | Hertz | 69,994.26 | 0.373ms | 0.813ms |
| 128B | fasthttp | 69,994.25 | 0.366ms | 0.772ms |
| 128B | Netty | 69,993.60 | 0.364ms | 1.236ms |
| 1KiB | Gnalloy | 69,993.22 | 0.424ms | 0.933ms |
| 1KiB | Hertz | 69,991.33 | 0.510ms | 0.901ms |
| 1KiB | fasthttp | 69,992.64 | 0.449ms | 0.899ms |
| 1KiB | Netty | 69,992.84 | 0.450ms | 2.932ms |

128B 时 Gnalloy RTT P50/P99 均最低。1KiB 时 Gnalloy P50 最低，但 P99 比 Hertz 高约 `3.5%`、比 fasthttp 高约 `3.8%`，因此明文固定速率场景尚不能声明所有尾延迟全面领先。

### 当前版本明文 1KiB 十轮复验

使用 benchmark `6e0d2b754350e63b2028a676494432f6e0d9b4e9` 和统一 `http1-load` 客户端，在 `70,000 ops/s` 下对四个服务端轮转执行 10 轮。Gnalloy 请求完整经过 TCP transport、Channel Pipeline、HTTP/1 codec 和业务 handler。共 40 个 case，全部 `errors=0`。

| 服务端 | 实际吞吐中位数 ops/s | RTT P50 中位数 ms | RTT P95 中位数 ms | RTT P99 中位数 ms | RTT P99 范围 ms |
| --- | ---: | ---: | ---: | ---: | ---: |
| Gnalloy | **69,992.83** | **0.432** | **0.719** | 0.903 | 0.819-1.063 |
| Hertz | 69,991.80 | 0.513 | 0.739 | 0.893 | 0.760-1.298 |
| fasthttp | 69,991.78 | 0.451 | 0.751 | **0.891** | 0.844-1.013 |
| Netty | 69,991.36 | 0.452 | 0.808 | 3.322 | 1.199-5.364 |

Gnalloy 的 RTT P50/P95 最低，但 RTT P99 比 fasthttp 高 `1.39%`，比 Hertz 高 `1.18%`，仍未在明文 1KiB 低负载尾延迟上全面领先。Netty 的客户端 schedule-delay 尖峰较大，排名继续使用实际发送后的 RTT 指标，不使用 total latency。

## HTTPS/1.1 完整 TLS 矩阵

完整 TLS 1.1/1.2/1.3 矩阵使用 Gnalloy 二进制 `74bf33a6d27f99213a8b70ecb65fd720075fc9101c8b25a59929551d7239f64c`。该二进制与当前二进制使用相同 Gnalloy core `13bb860`，当前二进制额外包含 TLS nonblocking 状态发布竞态修复 `7ca6328`。后文使用当前二进制重新验证了原先差距最大的 TLS 1.3/1KiB 场景。

### 饱和吞吐

| TLS | Payload | Gnalloy | Hertz | fasthttp | Netty |
|---|---:|---:|---:|---:|---:|
| 1.1 | 128B | 118.065k | 102.455k | 106.061k | 81.215k |
| 1.1 | 1KiB | 94.606k | 91.199k | 94.644k | 76.948k |
| 1.2 | 128B | 141.497k | 129.994k | 136.356k | 94.079k |
| 1.2 | 1KiB | 100.732k | 96.080k | 100.093k | 81.995k |
| 1.3 | 128B | 142.816k | 129.016k | 136.748k | 84.289k |
| 1.3 | 1KiB | 100.902k | 97.604k | 100.677k | 79.089k |

Gnalloy 在 6 个 TLS 场景均高于 Hertz 和 Netty；相对 fasthttp 为 5 个场景领先，TLS 1.1/1KiB 基本持平且低约 `0.04%`。

### 70k offered-rate RTT P99

| TLS | Payload | Gnalloy | Hertz | fasthttp | Netty |
|---|---:|---:|---:|---:|---:|
| 1.1 | 128B | 0.986ms | 1.076ms | 0.980ms | 5.075ms |
| 1.1 | 1KiB | 1.293ms | 1.192ms | 1.247ms | 5.481ms |
| 1.2 | 128B | 0.755ms | 0.807ms | 0.987ms | 4.749ms |
| 1.2 | 1KiB | 1.005ms | 0.999ms | 0.955ms | 5.471ms |
| 1.3 | 128B | 0.847ms | 0.855ms | 0.781ms | 5.392ms |
| 1.3 | 1KiB | 1.047ms | 0.974ms | 1.040ms | 5.547ms |

该完整矩阵中 Gnalloy 全部低于 Netty，128B 的三个 TLS 版本全部低于 Hertz；相对 Hertz/fasthttp 的若干 1KiB 场景仍有小幅差距。

## 当前 TLS 1.3/1KiB 复验

针对完整矩阵中原先落后的 TLS 1.3/1KiB，使用当前二进制重新执行四框架轮换五轮：

| 模式 | 服务端 | 实际吞吐 (ops/s) | RTT P50 | RTT P99 |
|---|---|---:|---:|---:|
| 饱和 | Gnalloy | 101,034.30 | 0.610ms | 1.253ms |
| 饱和 | Hertz | 96,863.31 | 0.622ms | 1.284ms |
| 饱和 | fasthttp | 100,309.19 | 0.614ms | 1.776ms |
| 饱和 | Netty | 72,085.12 | 0.644ms | 5.643ms |
| 70k/s | Gnalloy | 69,992.35 | 0.456ms | 0.970ms |
| 70k/s | Hertz | 69,989.88 | 0.522ms | 0.985ms |
| 70k/s | fasthttp | 69,990.96 | 0.482ms | 1.101ms |
| 70k/s | Netty | 69,989.38 | 0.569ms | 6.895ms |

当前 Gnalloy 在该关键场景的饱和吞吐、RTT P50/P99，以及等 offered-rate RTT P50/P99 均为本轮矩阵最低。与旧矩阵的 P99 差异包含运行间波动，不能仅凭这组结果把竞态修复解释为确定的性能提升。

## Poller A/B

使用同一当前二进制，仅切换 Gnalloy backend，未开启 fixed buffer、multishot accept 或 SQPOLL：

| Backend | 模式 | 吞吐中位数 (ops/s) | RTT P50 | RTT P99 |
|---|---|---:|---:|---:|
| epoll | 饱和 | 100,269.07 | 0.609ms | 5.413ms |
| io_uring | 饱和 | 99,339.03 | 0.612ms | 1.148ms |
| epoll | 70k/s | 69,990.49 | 0.454ms | 0.975ms |
| io_uring | 70k/s | 69,990.93 | 0.478ms | 1.136ms |

饱和 P99 容易被不同容量和偶发系统长尾放大，不用于选择 backend。io_uring 饱和吞吐低约 `0.93%`，70k/s RTT P99 高约 `16.5%`，因此当前 HTTP/1 默认继续使用 epoll，也不继续叠加 io_uring 专属开关。

## 已否决优化与 Profile

- 握手后 EventLoop owner 无锁密文队列在微基准中降低约 `34%-42%` 开销，但修复队列切换竞态后，70k/s RTT P99 仍从约 `1.085ms` 退化到约 `1.159ms`。该实现已删除，未提交。
- 明文 1KiB 将 flush 从 `read-complete` 改为 `immediate` 后，70k/s RTT P99 从 `0.933ms` 降到 `0.910ms`，但饱和吞吐从 `103.542k` 降到 `103.301k ops/s`，饱和 P99 也从 `0.986ms` 升到 `1.020ms`。该负载特定取舍不作为通用默认值。
- TLS 1.3/1KiB、70k/s CPU profile 中，`transport-tcp.writeNonblocking` 累计约 `37.5%`，Linux syscall 平坦占比约 `55.7%`；HTTP/1 encode 累计约 `2.6%`，TLS crypto 为个位数占比。
- 18 秒 runtime trace 的总 scheduler latency 约 `128ms`，未发现 GC 或 Go scheduler 主瓶颈。当前主要成本位于密文 socket 写出和内核 syscall，而不是 HTTP codec。

## 结论边界

- 当前明文矩阵中，Gnalloy 的饱和吞吐和饱和 P99 在 128B/1KiB 均领先；70k/s 的 128B P99 领先，1KiB P99 仍比 Hertz/fasthttp 高约 `3%-4%`。
- 完整 TLS 矩阵中，Gnalloy 吞吐全面领先 Hertz/Netty，并在 5/6 场景领先 fasthttp；固定速率 P99 全面领先 Netty，但没有在所有 TLS/payload 组合全面领先 Hertz/fasthttp。
- 当前二进制的 TLS 1.3/1KiB 复验已同时领先 Hertz、fasthttp 和 Netty，但不能由单一场景外推到所有 TLS 版本、连接数和硬件。
- gnet/netpoll 的 HTTP 项为 N/A。只有框架正式提供对应 HTTP codec 时，才能加入 HTTP 对比；TCP 对比结果不能替代 HTTP 对比。
- “任何场景都达到物理极限并全面领先所有框架”不是可证实的工程结论。本文只声明上述硬件、CPU 绑定、连接数、payload、TLS 套件和五轮中位数范围内的结果。

## 复现命令

```powershell
.\scripts\run-cross-host-http1-common-client.ps1 -ServerRepo /opt/test/gnalloy/benchmarks-http1-hertz-server-20260904 -ClientRepo /opt/test/gnalloy/benchmarks-http1-tls-client-20260904 -Protocols http1 -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 0 -Repetitions 5 -Frameworks gnalloy,gnet,netpoll,hertz,fasthttp,netty -GnalloyBackend epoll
.\scripts\run-cross-host-http1-common-client.ps1 -ServerRepo /opt/test/gnalloy/benchmarks-http1-hertz-server-20260904 -ClientRepo /opt/test/gnalloy/benchmarks-http1-tls-client-20260904 -Protocols http1 -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 70000 -Repetitions 5 -Frameworks gnalloy,gnet,netpoll,hertz,fasthttp,netty -GnalloyBackend epoll
.\scripts\run-cross-host-http1-common-client.ps1 -ServerRepo /opt/test/gnalloy/benchmarks-http1-hertz-server-20260904 -ClientRepo /opt/test/gnalloy/benchmarks-http1-tls-client-20260904 -Protocols https1 -TLSVersions 1.1,1.2,1.3 -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 0 -Repetitions 5 -Frameworks gnalloy,gnet,netpoll,hertz,fasthttp,netty -GnalloyBackend epoll
.\scripts\run-cross-host-http1-common-client.ps1 -ServerRepo /opt/test/gnalloy/benchmarks-http1-hertz-server-20260904 -ClientRepo /opt/test/gnalloy/benchmarks-http1-tls-client-20260904 -Protocols https1 -TLSVersions 1.1,1.2,1.3 -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 70000 -Repetitions 5 -Frameworks gnalloy,gnet,netpoll,hertz,fasthttp,netty -GnalloyBackend epoll
```
