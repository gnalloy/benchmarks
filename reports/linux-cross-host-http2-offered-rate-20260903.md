# Linux 跨机 HTTP/2 固定速率对比（2026-09-03）

## 目的

在相同客户端、相同 offered rate 下分离客户端调度积压和实际请求往返耗时，复验 Gnalloy 与 Netty 的 HTTPS/2 尾延迟。该测试走真实 Gnalloy TCP transport、TLS handler 和 HTTP/2 codec，不包含 benchmark 自写帧解析或固定响应旁路。

## 环境与方法

- 服务端：`172.16.8.171`，Debian 13，Linux `6.12.48+deb13-amd64`，Intel Xeon E3-1535M v5。
- 客户端：`172.16.8.172`，Debian 13，Linux `6.12.48+deb13-amd64`。
- 服务端和客户端分别固定到 CPU `0,1,2,3`，`GOMAXPROCS=4`，测试期间使用 performance governor。
- 协议：HTTPS/2、TLS 1.2、128B 响应；64 个连接，每连接 20,000 个测量请求和 1,000 个预热请求。
- 每个 offered rate 执行 5 轮；Gnalloy 和 Netty 服务端轮换先后顺序，任意时刻只运行一个测试案例。
- 所有服务端均由同一个 Gnalloy HTTP/2 客户端驱动。结果行中的 `framework=gnalloy` 是客户端标签，实际服务端由前置 `case=` 行确定。
- Gnalloy 服务端和客户端二进制 SHA-256：`30a495329f8e73604d60616fe047c543646ea651d7a211149692fe1d8d508f0d`。
- Netty 服务端 JAR SHA-256：`127ee936c86ac5d4cf0f84559845b0121e35f9537d8f577acbce4feff112a25f`。
- 两机空载 ICMP RTT：平均 `0.209ms`，范围 `0.109ms` 至 `0.298ms`。

## 五轮中位数

| Offered rate | 服务端 | 实际吞吐 (ops/s) | Total P99 | Schedule delay P99 | RTT P50 | RTT P95 | RTT P99 | RTT P99.9 |
|---:|---|---:|---:|---:|---:|---:|---:|---:|
| 60,000 | Gnalloy | 59,998.89 | 2.619ms | 1.515ms | 0.401ms | 0.753ms | 1.359ms | 4.492ms |
| 60,000 | Netty | 59,998.80 | 280.663ms | 278.979ms | 0.458ms | 0.994ms | 2.239ms | 6.212ms |
| 75,000 | Gnalloy | 74,997.47 | 3.244ms | 2.470ms | 0.451ms | 0.847ms | 1.400ms | 3.009ms |
| 75,000 | Netty | 74,996.45 | 362.829ms | 362.225ms | 0.588ms | 1.130ms | 1.722ms | 6.329ms |
| 85,000 | Gnalloy | 84,996.58 | 7.606ms | 6.326ms | 0.481ms | 0.922ms | 1.539ms | 2.785ms |
| 85,000 | Netty | 84,838.58 | 623.182ms | 622.323ms | 0.633ms | 1.145ms | 1.611ms | 5.874ms |

## 结论

- 60k 和 75k 是目标速率约束场景，吞吐相同，不用于声称吞吐领先。Gnalloy 的 RTT P99 分别比 Netty 低约 `39.3%` 和 `18.7%`。
- 85k 时 Gnalloy 仍维持约 `84.997k ops/s`，Netty 中位数降至约 `84.839k ops/s`。Gnalloy RTT P99 低约 `4.4%`，优势缩小但没有反转。
- Netty 在三档速率均出现明显 schedule backlog，85k 的 P99 达 `622ms`；Gnalloy 在 85k 的对应中位数为 `6.326ms`。Total latency 保留了这部分排队，RTT 则隔离实际发送后的网络、TLS、HTTP/2 codec 和服务端处理耗时。
- 本次结果说明此前“饱和模式下 Gnalloy P99 高于 Netty”主要混入了不同饱和吞吐所造成的排队效应。在等 offered rate 下，本场景的 Gnalloy RTT P99 五轮中位数均低于 Netty。
- 该结论只覆盖 HTTPS/2、TLS 1.2、128B 和本次两台主机。TLS 1.3、1KiB、HTTP/3 与 QUIC 必须分别测试，不能由本表外推。

## 复现命令

```powershell
.\scripts\run-cross-host-http2-common-client.ps1 -Scenarios https2-tls12 -Payloads 128 -Frameworks gnalloy,netty -Repetitions 5 -TargetRate 60000
.\scripts\run-cross-host-http2-common-client.ps1 -Scenarios https2-tls12 -Payloads 128 -Frameworks gnalloy,netty -Repetitions 5 -TargetRate 75000
.\scripts\run-cross-host-http2-common-client.ps1 -Scenarios https2-tls12 -Payloads 128 -Frameworks gnalloy,netty -Repetitions 5 -TargetRate 85000
```
