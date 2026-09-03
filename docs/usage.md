# Usage

[简体中文](usage.zh-CN.md) | [Docs Index](README.md)

## Requirements

- Go 1.25 or newer, matching the module `go` directive.
- A Gnalloy application, recipe, example, or benchmark harness that owns lifecycle and deployment configuration.
- Standalone module verification should set `GOWORK=off` so the module is tested through its published dependency graph.

## Install
```bash
go get gnalloy.org/benchmarks@dev
```

## Import
```go
import "gnalloy.org/benchmarks"
```

## Integration Pattern
- Configure scenario, protocol, payload size, concurrency, warmup count, measured repetitions, timeout, and report output before comparing results.
- Report throughput and latency separately. Do not compare numbers collected on different hosts as a ranking.

## Full Parity Matrix

Build external harnesses first, then run the built-in matrix from the workspace root that contains the `benchmarks` directory:

```bash
(cd benchmarks && GOWORK=off GOTOOLCHAIN=local go build -o external/bin/gnalloy-parity ./cmd/gnalloy-parity)
./benchmarks/external/bin/gnalloy-parity -matrix linux-full -strict-external -format markdown -out benchmarks/reports/linux-full-parity.md -timeout 6h
```

The Linux matrix covers TCP, UDP, QUIC stream, HTTP/1.1, HTTP/2, HTTP/3, HTTPS TLS 1.1/1.2/1.3, Gnalloy, Netty, gnet, fasthttp, and netpoll. Illegal or unsupported combinations are recorded as skipped scenarios with explicit reasons. Executable scenarios pass a 15-minute harness timeout and use a 16-minute outer scenario guard.

## Cross-Host HTTP/1 Benchmark

Run the common HTTP/1 client on Debian `172.16.8.172` while each server runs serially on Debian `172.16.8.171`:

```powershell
.\scripts\run-cross-host-http1-common-client.ps1
```

The runner uses Gnalloy's real TCP transport, channel pipeline, HTTP/1 codec, TLS handler when enabled, and application handler. Its default `auto` flush policy selects event-loop batching for plaintext HTTP/1 and immediate flush for HTTPS/1. Plaintext benefits from write batching, while the TLS pipeline already coalesces each small response before encryption and avoids an additional event-loop-tail delay with immediate flush. Pass `-GnalloyFlushStrategy` to test one fixed public channel strategy across all cases. gnet and netpoll are reported as `N/A` because neither framework provides an HTTP codec; the benchmark does not add its own parser or fixed-response bypass.

## Focused Linux TCP Benchmark

After building `external/bin/gnalloy-bench` and `external/bin/netpoll-bench`, run the focused TCP comparison on one Linux host:

```bash
./scripts/run-linux-tcp-optimized.sh
```

The script runs one process at a time, alternates Gnalloy/netpoll order across five repetitions, and records host metadata, parameters, binary hashes, throughput, and latency. Environment variables control payloads, load, cooldown, executable paths, output path, and process priority. It does not change the CPU governor or stop host services. See `reports/linux-tcp-read-write-optimization.md` for the verified 16KiB result and syscall analysis.

## Focused Linux UDP Benchmark

Build `external/bin/gnalloy-bench`, `external/bin/gnet-bench`, `external/bin/netty-bench.jar`, and `external/bin/udp-load`, then run all UDP servers through the same client on one Linux host:

```bash
./scripts/run-linux-udp-common-client.sh
TARGET_RATE=60000 ./scripts/run-linux-udp-common-client.sh
```

The default run measures saturation throughput. `TARGET_RATE` sets a shared aggregate request rate for a comparable latency run. Fixed-rate `latency` starts at each request's scheduled send time, `scheduleDelay` measures how late the client starts the send, and `roundTripLatency` measures from the actual send attempt until the echo arrives. This decomposition keeps scheduling delay and backlog visible while separating them from network and server processing. The script pins server and client processes to disjoint physical cores, verifies package/core IDs through Linux sysfs, rotates framework order, runs cases serially, and records the CPU topology, governor, parameters, and binary hashes. Set `REQUIRE_PHYSICAL_CPU_ISOLATION=0` only when the host cannot provide physically isolated CPU sets, and do not compare that result with an isolated run. Netpoll and fasthttp remain excluded because their benchmark harnesses do not provide comparable UDP servers.

To validate UDP over the LAN with fixed host roles, run the common client on Debian `172.16.8.172` and each server serially on Debian `172.16.8.171`. The default paths are `/opt/test/gnalloy/benchmarks-cross-host` for the client checkout and `/opt/test/gnalloy/benchmarks-e481c9a` for the server checkout. Build `external/bin/udp-load` on the client and all server binaries on the server, then run the orchestrator from the control workstation:

```powershell
.\scripts\run-cross-host-udp-common-client.ps1
.\scripts\run-cross-host-udp-common-client.ps1 -TargetRate 60000
```

For diagnosis, select only Gnalloy and capture one probe type per run. Profiles are written below the server checkout and downloaded to `reports/raw/cross-host-profiles`; instrumented results must not be compared with normal performance runs:

```powershell
.\scripts\run-cross-host-udp-common-client.ps1 -Frameworks gnalloy -Payloads 128 -Repetitions 1 -CaptureCPUProfile
.\scripts\run-cross-host-udp-common-client.ps1 -Frameworks gnalloy -Payloads 128 -Repetitions 1 -TargetRate 60000 -CaptureRuntimeTrace
```

The client defaults to CPUs `0,1,2,4`, avoiding its NIC IRQ on CPU 11 and using four distinct physical cores. The runner saves each host's CPU governors, selects `performance` for the run, and restores every original value during cleanup. Cross-host results include both NICs, the switch, and both operating-system network stacks. Keep them separate from same-host CPU-isolated results.

## Cross-host HTTP/2 Benchmark

Run the shared Gnalloy HTTP/2 client on Debian `172.16.8.172` against each server on Debian `172.16.8.171`. Cases and frameworks run serially. A zero target rate measures saturation; a positive rate measures all servers at the same aggregate offered load:

```powershell
.\scripts\run-cross-host-http2-common-client.ps1
.\scripts\run-cross-host-http2-common-client.ps1 -TargetRate 60000
```

For fixed-rate runs, `latency` covers scheduled send time through response completion, `scheduleDelay` covers client-side dispatch delay, and `roundTripLatency` covers actual send through response completion. Use `roundTripLatency` to compare the network, protocol stack, codec, handler, and server processing path; keep `scheduleDelay` visible to identify client backlog. HTTP/2 over TLS supports TLS 1.2 and TLS 1.3. gnet, netpoll, and fasthttp are reported as `N/A` because their benchmark servers do not provide native HTTP/2 codecs.

## API Selection

Use the API inventory to choose the exact constructor or option type for your protocol path:

```bash
go doc gnalloy.org/benchmarks
```

Common current entry points:
- `var ErrInvalidRunner = errors.New("gnalloy/benchmarks/benchdiff: invalid runner") ...`
- `type Config struct{ ... }`
- `type Config struct{ ... }`
- `const ProtocolHTTP1 = "http1"`
- `type Config struct{ ... }`
- `const DefaultServerName = "gnalloy.local"`
- `var ErrInvalidSpec = errors.New("gnalloy/benchmarks/parity: invalid spec") ...`
- `type ExternalHarnessOptions struct{ ... }`

## Cross-Module Assembly

When multiple Gnalloy repositories are developed together, create a local `go.work` file in your chosen workspace. Keep application-local `replace` directives out of published library modules unless the change is intentionally temporary and never committed.

## Error Handling

Network input, peer behavior, platform capability, and timeout failures must be handled as normal errors. Do not recover protocol correctness by panicking. Return or propagate the module error and close the affected Channel when ownership requires it.
