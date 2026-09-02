# Linux UDP 性能优化与对比报告

测试日期：2026-09-02

## 1. 范围与结论边界

本报告只比较各框架实际支持的 UDP echo 服务端：

| 框架 | UDP echo | 纳入原因 |
|---|---:|---|
| Gnalloy | 是 | 被测实现 |
| gnet | 是 | 可比实现 |
| Netty | 是 | 可比实现 |
| netpoll | 否 | 没有独立、可比的 UDP 服务端 API |
| fasthttp | 否 | 只提供 HTTP/1 服务端，不支持 UDP |

结论按测试场景分别成立，不把单一 payload、单一负载或不同主机上的数字扩展为通用排名。

## 2. 代码版本

| 仓库 | 分支 | 提交 | 作用 |
|---|---|---|---|
| benchmarks | dev | `e897219` | benchmark 消费正式批量接收实现 |
| transport-udp | dev | `2fd09e2` | Linux `recvmmsg` 批量接收 |
| transport-udp | dev | `26a032b` | 每 endpoint 复用 `sendmmsg` 描述符 |
| gnalloy | dev | `5357513` | EventLoop、Pipeline 与 epoll 核心 |

正式服务端二进制 SHA-256：

```text
e461f3d734ae3cfeff9abc1eea86852355a8f61dbe04728f18a456274f094e5e  gnalloy-bench
ffbdb5227259f1591383e615ce2ea5abbde65d4e5255aac3b668be88c22d30d0  gnet-bench
60d54dd6d132f99513cf43178189b3d75721b1510c8ebefc4ee0c4f34eb483da  netty-bench.jar
```

## 3. 环境与方法

服务端为 `172.16.8.171`，客户端为 `172.16.8.172`，测试目录均位于 `/opt/test/gnalloy/<独立目录>`。

| 项目 | 171 服务端 | 172 客户端 |
|---|---|---|
| 系统 | Debian 13, Linux 6.12.48 | Debian 13, Linux 6.12.48 |
| CPU | Xeon E3-1535M v5, 4C/8T | Intel 0000, 8C/16T |
| 网卡 | 1GbE `enp1s0` | 1GbE `eno1` |
| 工具链 | Go 1.27.0, Java 21 | Go 1.27.0 |

所有框架在同一服务端、同一客户端上串行运行。每个 case 都先停止上一个服务端，再启动下一个服务端，并在 case 之间冷却。测试期间 governor 临时切换为 `performance`，结束后恢复为 `powersave`。共同客户端、payload、连接数、消息数、预热和采样率保持一致。

主要参数：

- 64 个 connected UDP 客户端。
- 每连接预热 1,000 条消息。
- 饱和测试每连接 20,000 条消息。
- 固定速率测试每连接 10,000 条消息，总速率 60,000 ops/s。
- 每 64 条消息采集一次应用延迟。
- Linux 批量读上限为 2，兼顾吞吐和尾延迟。
- 最终低资源拓扑使用 2 个 worker，分别固定到两个物理核；boss 使用同核超线程。

## 4. 饱和结果

### 4.1 128B，相同 2-loop 资源，三轮中位数

| 框架 | 吞吐量 (ops/s) | P99 (ms) | 错误 |
|---|---:|---:|---:|
| Gnalloy | 287,330 | 0.340 | 0 |
| gnet | 206,397 | 0.545 | 0 |
| Netty | 161,899 | 0.643 | 0 |

Gnalloy 相对 gnet 吞吐高 `39.21%`、P99 低 `37.62%`；相对 Netty 吞吐高 `77.47%`、P99 低 `47.11%`。

### 4.2 跨机完整矩阵，4-loop，五轮中位数

| Payload | 框架 | 吞吐量 (ops/s) | P99 (ms) | 错误 |
|---:|---|---:|---:|---:|
| 128B | Gnalloy | 287,874 | 0.385 | 0 |
| 128B | gnet | 252,013 | 0.467 | 0 |
| 128B | Netty | 158,879 | 0.610 | 0 |
| 1KiB | Gnalloy | 105,707 | 0.930 | 0 |
| 1KiB | gnet | 106,975 | 1.089 | 0 |
| 1KiB | Netty | 104,582 | 1.242 | 0 |

1KiB 时两端 1GbE 已成为吞吐上限。按 1024B UDP payload 加以太网、IPv4、UDP、preamble 和 IFG 开销估算，理论上限约为 114,679 包/秒；Gnalloy 已达到约 `92.18%` 线速，gnet 约为 `93.29%`。该场景中 Gnalloy 吞吐比 gnet 低 `1.18%`，不能据此声称跨机 1KiB 吞吐领先。

### 4.3 171 本机 loopback，客户端与服务端物理核隔离，1KiB 五轮中位数

| 框架 | 吞吐量 (ops/s) | P99 (ms) | 错误 |
|---|---:|---:|---:|
| Gnalloy | 278,377 | 0.689 | 0 |
| gnet | 193,047 | 2.808 | 0 |
| Netty | 121,821 | 0.776 | 0 |

移除 1GbE 瓶颈后，Gnalloy 相对 gnet 吞吐高 `44.20%`、P99 低 `75.45%`；相对 Netty 吞吐高 `128.51%`、P99 低 `11.12%`。这组结果证明跨机 1KiB 的吞吐接近不是服务端 CPU 处理能力相同，而是物理链路封顶。

## 5. 固定 60k 结果

相同 2-loop 资源下的三轮中位数：

| Payload | 框架 | 总 P99 (ms) | RTT P99 (ms) | 错误 |
|---:|---|---:|---:|---:|
| 128B | Gnalloy | 0.562 | 0.314 | 0 |
| 128B | gnet | 0.225 | 0.191 | 0 |
| 128B | Netty | 0.223 | 0.186 | 0 |
| 1KiB | Gnalloy | 0.272 | 0.235 | 0 |
| 1KiB | gnet | 0.265 | 0.229 | 0 |
| 1KiB | Netty | 0.271 | 0.234 | 0 |

固定速率总延迟包含客户端在计划发送时间之后的调度延迟，不能直接作为服务端 RTT。128B 的 Gnalloy 三轮中有两轮同时出现客户端 schedule-delay 尖峰，因此中位数仍偏高；1KiB 三方已处于约 1% 到 3% 的同档范围。

为分离客户端 goroutine 唤醒，额外在 172 使用 `tcpdump` 内核时间戳进行短测。每个框架完整配对 134,400 个请求，未配对数为 0：

| 框架 | 内核 RTT P99 (ms) | 同轮应用 RTT P99 (ms) |
|---|---:|---:|
| Gnalloy | 0.165 | 0.181 |
| gnet | 0.422 | 0.401 |

该诊断轮证明 Gnalloy 的低负载尾延迟并非稳定落后，但单轮抓包也不足以替代多轮统计。因此固定 60k 结果应表述为“常态延迟同档、轮次尾部存在环境噪声”，不能声称在所有低负载运行中全面领先。

## 6. Profile 与根因

1KiB 饱和 CPU profile 的持续时间为 13.33 秒，累计 CPU 样本 16.78 秒，服务端只消耗约 1.26 个 CPU 核：

| 热点 | 平坦占比 | 累计占比 |
|---|---:|---:|
| Linux syscall | 79.14% | 79.14% |
| UDP `flushOutbound` | 1.55% | 47.97% |
| UDP batch reader prepare | 0.48% | 1.67% |
| Pipeline `FireChannelRead` | 0.42% | 1.37% |
| HeapAllocator Acquire | 0.36% | 1.19% |

runtime trace 没有显示持续 GC 压力；EventLoop 调度延迟累计约 154ms。Pipeline、handler、分配器都不是当前主要 CPU 瓶颈。跨机 1KiB 的主要限制是 1GbE，128B 的主要成本是收发系统调用。

## 7. 被否决的候选

| 候选 | 结果 | 决策 |
|---|---|---|
| 每个 datagram `WriteAndFlush` | 128B 饱和吞吐下降 2.42%，P99 上升 16.15%；固定 60k RTT P99 上升约 9.5% | 否决，恢复 read-complete 合并 flush |
| `max-messages-per-read=1` | 128B 饱和只有 236,963 ops/s，且 RTT P99 未改善 | 否决 |
| 1 worker/1 socket | 固定 60k 出现 1ms 到 9ms 级 P99 | 否决 |
| 2 个物理 worker | 128B 饱和吞吐基本不变，P99 比 4 个超线程 worker 低 11.6%；1KiB 固定 RTT P99 降至 0.235ms | 保留为 171 的推荐配置 |

## 8. 最终判断

- Gnalloy 在 128B 饱和场景同时领先 gnet 和 Netty 的吞吐与 P99。
- Gnalloy 在移除 1GbE 瓶颈后的 1KiB 场景同时领先 gnet 和 Netty 的吞吐与 P99。
- 跨机 1KiB 吞吐已接近 1GbE 线速，不能用于区分服务端 CPU 极限。
- 固定 60k 的常态 RTT 已与对手处于同档，但多轮应用 P99 仍受客户端调度和网络尾部噪声影响，当前证据不支持“任何场景都全面超过”的绝对结论。
- 下一轮若要继续区分 1KiB 以上吞吐，必须使用 10GbE 或更高速链路；若要建立低负载 P99 门禁，应把内核 RTT 时间戳采集正式纳入 benchmark，并增加重复次数与置信区间。

## 9. 验证状态

- Windows：`go test ./...`、`go test -race ./...`、`go vet ./...` 通过。
- Debian 171：原生 `go test ./...`、`go test -race ./...`、`go vet ./...` 通过。
- UDP IPv4/IPv6、空接收队列、批量接收器复用测试通过。
- 所有性能 case 串行执行；测试结束后 171/172 governor 均恢复为 `powersave`，无残留 benchmark 进程。
