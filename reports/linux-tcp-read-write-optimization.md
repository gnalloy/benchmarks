# Gnalloy Linux TCP 读写模型优化报告

## 结论

在同一台 Debian 主机上，以 AB/BA 顺序严格串行执行五轮 16KiB TCP echo 测试后，Gnalloy 的吞吐中位数为 `166,981.55 ops/s`，netpoll 为 `156,392.37 ops/s`，Gnalloy 高 `6.77%`。Gnalloy 的 P99 中位数为 `1.082 ms`，netpoll 为 `1.276 ms`，Gnalloy 低 `15.22%`。

该结论只适用于本报告记录的主机、版本和参数，不代表所有硬件、payload、连接数与业务 handler 的普遍排名。

## 环境与版本

| 项目 | 值 |
| --- | --- |
| 时间 | 2026-09-02 01:47 CST |
| 主机 | `xigexb-dev2` / `172.16.8.171` |
| 系统 | Debian, Linux `6.12.48+deb13-amd64`, amd64 |
| CPU | Intel Xeon E3-1535M v5, 4 核 8 线程 |
| Go | `go1.27.0 linux/amd64` |
| benchmark | `9a043cb112b4b9b18bbf171096c6ed8c524aa406` |
| gnalloy core | `v0.0.0-20260901173305-c5aafaacc185` |
| transport-tcp | `v0.0.0-20260901170722-14a382ab4ab3` |
| Gnalloy 二进制 SHA-256 | `3d088f3fb18beb091e970a3d100a1f5642e10e67061c3cd02fb39a359907017e` |
| netpoll 二进制 SHA-256 | `38d9f4822c2b2013dc3e5b329d817926ecf5d0fe4d35cef3a471da7aa97b82cc` |

两个二进制均从干净 checkout 使用 `GOWORK=off GOTOOLCHAIN=local` 构建。正式采样期间 CPU governor 临时设为 `performance`，测试结束后已恢复为 `powersave`。主机存在常驻服务，因此报告保留五轮中位数，并使用相同优先级和冷却时间公平处理两个框架。

## 参数

| 参数 | Gnalloy | netpoll |
| --- | --- | --- |
| payload | 16,384 bytes | 16,384 bytes |
| connections | 64 | 64 |
| messages/connection | 20,000 | 20,000 |
| warmup/connection | 1,000 | 1,000 |
| latency sampling | 1/128 | 1/128 |
| GOMAXPROCS | 8 | 8 |
| poller | 1 epoll loop | netpoll default |
| handler | 8 fixed workers, per-connection FIFO | netpoll callback |
| read buffer | 32,768 bytes | framework default |
| flush | event-loop batch | framework default |
| cooldown | 10 seconds | 10 seconds |

每轮只有一个 benchmark 进程运行。奇数轮顺序为 Gnalloy -> netpoll，偶数轮顺序为 netpoll -> Gnalloy。

## 五轮结果

| Run | First | Gnalloy throughput ops/s | Gnalloy P99 ms | netpoll throughput ops/s | netpoll P99 ms |
| ---: | --- | ---: | ---: | ---: | ---: |
| 1 | Gnalloy | 161,023.11 | 1.148 | 158,910.98 | 1.276 |
| 2 | netpoll | 168,970.33 | 1.062 | 162,167.75 | 1.048 |
| 3 | Gnalloy | 159,566.39 | 1.200 | 156,392.37 | 1.207 |
| 4 | netpoll | 166,981.55 | 1.082 | 130,402.23 | 1.934 |
| 5 | Gnalloy | 174,395.59 | 0.943 | 154,870.01 | 1.302 |
| **Median** | - | **166,981.55** | **1.082** | **156,392.37** | **1.276** |

| 中位数资源 | Gnalloy | netpoll | Gnalloy 相对变化 |
| --- | ---: | ---: | ---: |
| RSS bytes | 24,539,136 | 27,365,376 | -10.33% |
| heap allocated bytes | 6,184,840 | 10,497,760 | -41.08% |
| GC count | 1 | 17 | -94.12% |

## 根因与修复

旧路径存在三个独立开销：handler worker 的写请求需要再次跳转到 owner loop；Linux `unix.Write` 引入额外的 Go runtime syscall 包装；自适应读缓冲从过小容量启动，并在边界容量间振荡，使一个请求拆成多次 read、handler 和 write。

修复分为三个独立提交：

| 仓库 | 提交 | 作用 |
| --- | --- | --- |
| gnalloy | `423e2b3` | readiness 后端允许安全的 handler worker 并发写；存在 outbound/flush-complete/writability handler 时自动回退 owner loop |
| transport-tcp | `14a382a` | Linux readiness 路径使用 raw nonblocking `read/write/writev`，保留 EINTR/EAGAIN 与 iovec 边界处理 |
| gnalloy | `c5aafaa` | 读缓冲从配置容量启动，采用滞后缩容，避免 1KiB/2KiB 振荡与请求碎片化 |

benchmark harness 使用固定 worker executor 和每连接 FIFO。该模型保留连接内消息顺序与消息所有权，同时让不同连接并行处理；它不改变上层 Channel、Pipeline、codec 或 handler API。

## 16KiB 读缓冲验证

16KiB payload 与 16KiB read buffer 完全相等时，第一次 read 填满缓冲，read loop 仍需再调用一次 read 才能通过 EAGAIN 确认数据耗尽。32KiB read buffer 使同一请求成为 short read，可以直接结束本轮读取。

三轮交替筛选结果：

| Read buffer | Throughput median ops/s | P99 median ms |
| ---: | ---: | ---: |
| 16KiB | 136,719.87 | 1.120 |
| 32KiB | 164,874.22 | 1.104 |

32KiB read buffer 将吞吐提高 `20.59%`，P99 降低 `1.38%`。

`strace -f -c` 使用 64 连接、每连接 5,000 条消息进行单变量验证。strace 会显著降低绝对吞吐，因此这里只比较 syscall 数量：

| Read buffer | read calls | read EAGAIN | write calls |
| ---: | ---: | ---: | ---: |
| 16KiB | 1,415,376 | 704,118 | 711,247 |
| 32KiB | 1,057,088 | 352,140 | 704,937 |

成功 read 和 write 数量基本不变，而总 read 调用减少 `25.31%`，EAGAIN read 减少 `49.99%`，与吞吐提升方向一致。

## 复验

提交前后执行了以下检查：

```text
Windows:
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test -race ./... -count=1

Debian:
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test -race ./... -count=1
```

可复现测试入口为 `scripts/run-linux-tcp-optimized.sh`。脚本记录主机、内核、governor、参数和二进制哈希，按 payload 执行五轮 AB/BA 串行测试；它不会修改 governor 或停止主机服务。
