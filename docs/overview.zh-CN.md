# 概览

[English](overview.md) | [文档索引](README.zh-CN.md)

## 目标

Gnalloy 模块化网络栈的基准、性能对标与压测工具。

该模块包含 benchmark、parity 与 pressure-test 工具，用于测量和回归防护，不作为生产运行时依赖。

## 仓库身份

- 模块路径：`gnalloy.org/benchmarks`
- GitHub 仓库：`github.com/gnalloy/benchmarks`
- 默认分支：`dev`
- 许可证：Apache-2.0

## 包结构
- `gnalloy.org/benchmarks/benchdiff`（`benchdiff`）
- `gnalloy.org/benchmarks/cmd/gnalloy-benchdiff`（`main`）
- `gnalloy.org/benchmarks/external/internal/benchh2`（`benchh2`）
- `gnalloy.org/benchmarks/external/internal/benchh3`（`benchh3`）
- `gnalloy.org/benchmarks/external/internal/benchtls`（`benchtls`）
- `gnalloy.org/benchmarks/microbench`（`microbench`）
- `gnalloy.org/benchmarks/parity`（`parity`）

## 直接 Gnalloy 依赖

- `gnalloy.org/codec-http1`
- `gnalloy.org/codec-http3`
- `gnalloy.org/gnalloy`
- `gnalloy.org/transport-http3`
- `gnalloy.org/transport-quic`

## 当前仓库集合中的直接下游

- `gnalloy.org/examples`

## 架构位置

Gnalloy 保持核心小而轻依赖。本仓库围绕单一职责形成可替换模块，通过显式 Go package 连接，而不是依靠运行时发现。

## 兼容性

公共导入路径是 `gnalloy.org/benchmarks`。首个稳定 tag 发布前，请按依赖策略使用 `@dev` 或明确的 pseudo-version。
