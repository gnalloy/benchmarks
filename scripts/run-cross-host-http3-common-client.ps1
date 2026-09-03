# HTTP/3 兼容入口；统一实现位于 QUIC 协议族 runner。
& (Join-Path $PSScriptRoot "run-cross-host-quic-common-client.ps1") -Protocol http3 @args
