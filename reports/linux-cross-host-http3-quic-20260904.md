# Linux 跨机 HTTP/3 与 QUIC stream 对比（2026-09-04）

## 目的

在 171 服务端、172 客户端上串行对比 Gnalloy 与 Netty 的 HTTP/3 和 QUIC stream。Gnalloy 案例使用正式 `transport-quic`、`transport-http3`、`codec-http3` 与 `transport-quic/application` codec；benchmark 不自写 HTTP/QUIC 帧解析、不复制 codec framing，也没有固定响应旁路。

## 环境与方法

- 服务端：`172.16.8.171`，Debian 13，Linux `6.12.48+deb13-amd64`，Intel Xeon E3-1535M v5 @ 2.90GHz。
- 客户端：`172.16.8.172`，Debian 13，Linux `6.12.48+deb13-amd64`。
- 服务端和客户端分别固定 CPU `0,1,2,3`，`GOMAXPROCS=4`，测试期间使用 performance governor。
- 64 个连接，每连接 5,000 个测量请求和 500 个预热请求；每 64 个请求采样一次延迟。
- 每个案例执行 5 轮，取中位数；Gnalloy 与 Netty 轮换先后顺序，任意时刻只运行一个案例。
- 所有服务端由同一个 Gnalloy 客户端驱动。输出行的 `framework=gnalloy` 是客户端标签，实际服务端由前置 `case=` 确定。
- HTTP/3 和 QUIC 均使用 TLS 1.3；HTTP/3 ALPN 为 `h3`，QUIC stream ALPN 为 `gnalloy-quic`。
- Hertz、gnet、netpoll、fasthttp 没有本矩阵所需的原生 HTTP/3 或通用 QUIC stream codec，因此标记为 N/A，未用自写协议实现补齐。

## HTTP/3

HTTP/3 服务端 Gnalloy 二进制 SHA-256 为 `9e8452cf643a8c2381c63ac4d4097decdbdf8c8c85b2e0b06256daf03bd7af26`，Netty JAR 为 `6c4f6fd727aeb42b177d6d1ac71701a0b7b122c3ac509b6e89401b87fcda2f13`，公共客户端为 `cb50eccda1fd477d1d8587956b69d9b6d423a5b01837e684cd0d952ac5783d54`。饱和测试前两机 ICMP RTT 平均 `0.273ms`。

### 饱和模式

| Payload | 服务端 | 吞吐 (ops/s) | RTT P50 | RTT P95 | RTT P99 | RTT P99.9 |
|---:|---|---:|---:|---:|---:|---:|
| 128B | Gnalloy | 45,329.03 | 1.185ms | 2.887ms | 4.380ms | 10.792ms |
| 128B | Netty | 26,964.62 | 2.117ms | 3.426ms | 4.513ms | 29.445ms |
| 1KiB | Gnalloy | 45,229.91 | 1.197ms | 2.933ms | 4.277ms | 8.401ms |
| 1KiB | Netty | 24,506.41 | 2.381ms | 3.665ms | 5.076ms | 31.685ms |

### 20k offered-rate

| Payload | 服务端 | 实际吞吐 (ops/s) | Total P99 | Schedule P99 | RTT P50 | RTT P95 | RTT P99 | RTT P99.9 |
|---:|---|---:|---:|---:|---:|---:|---:|---:|
| 128B | Gnalloy | 19,999.69 | 2.536ms | 0.745ms | 0.265ms | 0.933ms | 1.805ms | 2.932ms |
| 128B | Netty | 19,935.16 | 101.553ms | 98.286ms | 0.833ms | 3.540ms | 4.501ms | 23.160ms |
| 1KiB | Gnalloy | 19,999.65 | 6.155ms | 3.085ms | 0.290ms | 1.243ms | 2.819ms | 7.136ms |
| 1KiB | Netty | 19,986.04 | 143.648ms | 141.131ms | 1.043ms | 3.756ms | 4.966ms | 29.184ms |

HTTP/3 在本矩阵内闭环：Gnalloy 的饱和吞吐、饱和 RTT P99，以及固定 offered-rate 下的 total、schedule 和 RTT P99 均低于 Netty。

## QUIC stream

QUIC stream 使用 `transport-quic` 提供的连接与 stream、`application.LengthPrefixedCodec` 提供的编解码，以及有界 stream 执行器。Gnalloy 不再为每条短 stream 创建 goroutine，响应缓冲区可复用，codec 将完整帧单次提交。Netty 使用官方 `LengthFieldBasedFrameDecoder` 与 `LengthFieldPrepender`，没有 benchmark 自写长度字段编码。

最终 Gnalloy 服务端二进制 SHA-256 为 `b5d652b76e3237faadf96530461fbdb58e725a11175dc590eed91f8e96a40c57`，使用 Netty 官方长度 codec 的 JAR 为 `62747f076b58ed92bca01b633f3e76803df558e4ac5a3f67a4ab8dafb1b0e69a`，公共客户端为 `cb50eccda1fd477d1d8587956b69d9b6d423a5b01837e684cd0d952ac5783d54`。

### 饱和模式

| Payload | 服务端 | 吞吐 (ops/s) | RTT P50 | RTT P95 | RTT P99 | RTT P99.9 |
|---:|---|---:|---:|---:|---:|---:|
| 128B | Gnalloy | 70,633.86 | 0.767ms | 1.909ms | 3.049ms | 5.390ms |
| 128B | Netty | 38,445.64 | 1.434ms | 2.481ms | 2.876ms | 3.277ms |
| 1KiB | Gnalloy | 65,469.42 | 0.827ms | 2.003ms | 3.138ms | 6.075ms |
| 1KiB | Netty | 40,858.33 | 1.389ms | 2.534ms | 3.077ms | 5.614ms |

Gnalloy 饱和吞吐分别领先约 `83.7%` 和 `60.2%`，RTT P50/P95 也更低。饱和 RTT P99 仍分别高约 `6.0%` 和 `2.0%`；由于两端运行在不同实际吞吐，这一列不能代替等 offered-rate 延迟比较，也不能声称饱和点所有尾延迟均领先。

### 30k offered-rate

| Payload | 服务端 | 实际吞吐 (ops/s) | Total P99 | Schedule P99 | RTT P50 | RTT P95 | RTT P99 | RTT P99.9 |
|---:|---|---:|---:|---:|---:|---:|---:|---:|
| 128B | Gnalloy | 29,999.42 | 1.751ms | 0.108ms | 0.228ms | 0.653ms | 1.719ms | 2.836ms |
| 128B | Netty | 29,993.34 | 57.643ms | 55.683ms | 0.720ms | 2.390ms | 2.836ms | 3.501ms |
| 1KiB | Gnalloy | 29,999.16 | 1.662ms | 0.437ms | 0.270ms | 0.918ms | 1.576ms | 3.187ms |
| 1KiB | Netty | 29,890.14 | 271.104ms | 269.083ms | 1.186ms | 2.603ms | 2.912ms | 3.661ms |

在相同 30k/s offered-rate 下，Gnalloy 的实际吞吐、total P99、schedule P99、RTT P50/P95/P99/P99.9 均领先。128B 与 1KiB RTT P99 分别低约 `39.4%` 和 `45.9%`。

## 结论边界

- 当前已测 HTTP/3 场景中，Gnalloy 的吞吐和 RTT P99 均领先 Netty。
- 当前 QUIC stream 等 offered-rate 场景中，Gnalloy 的吞吐保持率和所有记录的延迟分位均领先 Netty。
- QUIC stream 饱和点的 P99/P99.9 尚未全面领先；高出部分与 Gnalloy 更高的实际饱和吞吐同时存在，需要按容量曲线而不是单个饱和点继续分析。
- 测试只覆盖本次两台主机、4 个固定 CPU、64 连接、128B/1KiB 和 TLS 1.3，不能外推为所有机器和所有负载。
- 一次矩阵运行期间 172 上出现外部 `xz` 占用约 14 核并造成数百毫秒 GC pause；该运行已中止，结果未纳入任何表格。

## 复现命令

```powershell
.\scripts\run-cross-host-quic-common-client.ps1 -Protocol http3 -Frameworks gnalloy,netty -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 0 -Repetitions 5
.\scripts\run-cross-host-quic-common-client.ps1 -Protocol http3 -Frameworks gnalloy,netty -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 20000 -Repetitions 5
.\scripts\run-cross-host-quic-common-client.ps1 -Protocol quic-stream -Frameworks gnalloy,netty -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 0 -Repetitions 5
.\scripts\run-cross-host-quic-common-client.ps1 -Protocol quic-stream -Frameworks gnalloy,netty -Payloads 128,1024 -Connections 64 -Messages 5000 -WarmupMessages 500 -LatencySampleRate 64 -TargetRate 30000 -Repetitions 5
```
