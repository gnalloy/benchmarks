package matrix

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gnalloy.org/benchmarks/parity"
)

const (
	PlatformLinux   = "linux"
	PlatformWindows = "windows"
)

const (
	protocolTCP        = "tcp-echo"
	protocolUDP        = "udp-echo"
	protocolHTTP1      = "http1"
	protocolHTTPS1     = "https1"
	protocolHTTP2      = "http2"
	protocolHTTPS2     = "https2"
	protocolHTTP3      = "http3"
	protocolQUICStream = "quic-stream"
)

const (
	frameworkGnalloy  = "gnalloy"
	frameworkNetty    = "netty"
	frameworkGnet     = "gnet"
	frameworkFastHTTP = "fasthttp"
	frameworkNetpoll  = "netpoll"
)

var payloadsSmall = []payload{
	{Bytes: "128", Label: "128B"},
	{Bytes: "1024", Label: "1KiB"},
}

var payloadsTCP = []payload{
	{Bytes: "64", Label: "64B"},
	{Bytes: "1024", Label: "1KiB"},
	{Bytes: "16384", Label: "16KiB"},
}

var tlsVersions = []tlsVersion{
	{Value: "1.1", Tag: "tls11", CipherVariable: "TLS11_CIPHER"},
	{Value: "1.2", Tag: "tls12", CipherVariable: "TLS12_CIPHER"},
	{Value: "1.3", Tag: "tls13"},
}

type payload struct {
	Bytes string
	Label string
}

type tlsVersion struct {
	Value          string
	Tag            string
	CipherVariable string
}

// Profile 描述一个平台上的完整对标矩阵。
type Profile struct {
	Name               string
	Notes              string
	Platform           string
	GnalloyBackend     string
	GnalloyUDPBackend  string
	NettyBackend       string
	NetpollSupported   bool
	RepoRoot           string
	Connections        string
	Messages           string
	WarmupMessages     string
	LatencySampleRate  string
	Timeout            time.Duration
	CommandTimeout     time.Duration
	IncludeUnsupported bool
}

// LinuxFullSpec 返回 Linux 原生完整协议对标矩阵。
func LinuxFullSpec() parity.Spec {
	return FullSpec(LinuxProfile())
}

// WindowsFullSpec 返回 Windows 原生完整协议对标矩阵。
func WindowsFullSpec() parity.Spec {
	return FullSpec(WindowsProfile())
}

// LinuxProfile 返回 Debian/Linux 推荐的完整矩阵配置。
func LinuxProfile() Profile {
	return Profile{
		Name:               "gnalloy linux full parity matrix",
		Notes:              "Linux 完整对标矩阵：TCP、UDP、HTTP/1.1、HTTP/2、HTTP/3、HTTPS TLS 1.1/1.2/1.3 与纯 QUIC stream。HTTP/2 over TLS 跳过 TLS 1.1，HTTP/3/QUIC 固定 TLS 1.3；gnet、fasthttp、netpoll 仅执行各自真实支持的等价协议。",
		Platform:           PlatformLinux,
		GnalloyBackend:     "epoll",
		GnalloyUDPBackend:  "epoll",
		NettyBackend:       "epoll",
		NetpollSupported:   true,
		Connections:        "64",
		Messages:           "5000",
		WarmupMessages:     "500",
		LatencySampleRate:  "64",
		Timeout:            16 * time.Minute,
		CommandTimeout:     15 * time.Minute,
		IncludeUnsupported: true,
	}
}

// WindowsProfile 返回 Windows 原生完整矩阵配置。
func WindowsProfile() Profile {
	return Profile{
		Name:               "gnalloy windows full parity matrix",
		Notes:              "Windows 完整对标矩阵：TCP、UDP、HTTP/1.1、HTTP/2、HTTP/3、HTTPS TLS 1.1/1.2/1.3 与纯 QUIC stream。Netty 使用 NIO，Gnalloy 使用 IOCP/std；netpoll 没有 Windows 等价实现时标记为 unsupported。",
		Platform:           PlatformWindows,
		GnalloyBackend:     "iocp",
		GnalloyUDPBackend:  "iocp",
		NettyBackend:       "nio",
		NetpollSupported:   false,
		Connections:        "64",
		Messages:           "5000",
		WarmupMessages:     "500",
		LatencySampleRate:  "64",
		Timeout:            16 * time.Minute,
		CommandTimeout:     15 * time.Minute,
		IncludeUnsupported: true,
	}
}

// FullSpec 根据平台 profile 生成完整对标规格。
func FullSpec(profile Profile) parity.Spec {
	p := profile.withDefaults()
	builder := specBuilder{profile: p}
	builder.addTCP()
	builder.addUDP()
	builder.addHTTP1()
	builder.addHTTPS1()
	builder.addHTTP2()
	builder.addHTTPS2()
	builder.addHTTP3()
	builder.addQUICStream()
	return parity.Spec{
		Name:      p.Name,
		Notes:     p.Notes,
		Variables: p.variables(),
		Scenarios: builder.scenarios,
	}
}

func (p Profile) withDefaults() Profile {
	if strings.TrimSpace(p.Platform) == "" {
		p.Platform = PlatformLinux
	}
	if strings.TrimSpace(p.Name) == "" {
		p.Name = "gnalloy full parity matrix"
	}
	if strings.TrimSpace(p.GnalloyBackend) == "" {
		p.GnalloyBackend = "epoll"
	}
	if strings.TrimSpace(p.GnalloyUDPBackend) == "" {
		p.GnalloyUDPBackend = p.GnalloyBackend
	}
	if strings.TrimSpace(p.NettyBackend) == "" {
		p.NettyBackend = "epoll"
	}
	if strings.TrimSpace(p.RepoRoot) == "" {
		p.RepoRoot = "."
	}
	if strings.TrimSpace(p.Connections) == "" {
		p.Connections = "64"
	}
	if strings.TrimSpace(p.Messages) == "" {
		p.Messages = "5000"
	}
	if strings.TrimSpace(p.WarmupMessages) == "" {
		p.WarmupMessages = "500"
	}
	if strings.TrimSpace(p.LatencySampleRate) == "" {
		p.LatencySampleRate = "64"
	}
	if p.Timeout <= 0 {
		p.Timeout = 16 * time.Minute
	}
	if p.CommandTimeout <= 0 {
		p.CommandTimeout = p.Timeout - time.Minute
		if p.CommandTimeout <= 0 {
			p.CommandTimeout = p.Timeout
		}
	}
	return p
}

func (p Profile) variables() map[string]string {
	return map[string]string{
		"REPO_ROOT":                 p.RepoRoot,
		"CONNECTIONS":               p.Connections,
		"MESSAGES":                  p.Messages,
		"WARMUP_MESSAGES":           p.WarmupMessages,
		"LATENCY_SAMPLE_RATE":       p.LatencySampleRate,
		"COMMAND_TIMEOUT":           formatCommandTimeout(p.CommandTimeout),
		"GNALLOY_HTTP1_READ_BUFFER": "384",
		"HTTP1_ALPN":                "http/1.1",
		"HTTP2_ALPN":                "h2",
		"HTTP3_ALPN":                "h3",
		"QUIC_ALPN":                 "gnalloy-quic",
		"TLS11_CIPHER":              "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		"TLS12_CIPHER":              "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"GNALLOY_BENCH":             "./benchmarks/external/bin/gnalloy-bench",
		"NETTY_BENCH_JAR":           "./benchmarks/external/bin/netty-bench.jar",
		"GNET_BENCH":                "./benchmarks/external/bin/gnet-bench",
		"FASTHTTP_BENCH":            "./benchmarks/external/bin/fasthttp-bench",
		"NETPOLL_BENCH":             "./benchmarks/external/bin/netpoll-bench",
	}
}

type specBuilder struct {
	profile   Profile
	scenarios []parity.Scenario
}

func (b *specBuilder) addTCP() {
	for _, size := range payloadsTCP {
		b.add(gnalloyCommand(b.profile.GnalloyBackend, protocolTCP, size, nil))
		b.add(nettyCommand(b.profile.NettyBackend, protocolTCP, size, nil))
		b.add(gnetCommand(protocolTCP, size))
		if b.profile.NetpollSupported {
			b.add(netpollCommand(protocolTCP, size))
		}
	}
	if !b.profile.NetpollSupported {
		b.unsupported(frameworkNetpoll, protocolTCP, "tcp echo", "CloudWeGo netpoll harness is not supported on this platform")
	}
	b.unsupported(frameworkFastHTTP, protocolTCP, "tcp echo", "fasthttp is an HTTP server framework, not a raw TCP echo framework")
}

func (b *specBuilder) addUDP() {
	for _, size := range payloadsSmall {
		b.add(gnalloyCommand(b.profile.GnalloyUDPBackend, protocolUDP, size, nil))
		b.add(nettyCommand(b.profile.NettyBackend, protocolUDP, size, nil))
		b.add(gnetCommand(protocolUDP, size))
	}
	b.unsupported(frameworkNetpoll, protocolUDP, "udp echo", "CloudWeGo netpoll does not expose an equivalent UDP server API in this harness")
	b.unsupported(frameworkFastHTTP, protocolUDP, "udp echo", "fasthttp is HTTP-only and cannot execute UDP datagram echo")
}

func (b *specBuilder) addHTTP1() {
	for _, size := range payloadsSmall {
		b.add(gnalloyCommand(b.profile.GnalloyBackend, protocolHTTP1, size, nil))
		b.add(nettyCommand(b.profile.NettyBackend, protocolHTTP1, size, nil))
		b.add(gnetCommand(protocolHTTP1, size))
		b.add(fasthttpCommand(protocolHTTP1, size, nil))
		if b.profile.NetpollSupported {
			b.add(netpollCommand(protocolHTTP1, size))
		}
	}
	if !b.profile.NetpollSupported {
		b.unsupported(frameworkNetpoll, protocolHTTP1, "http1", "CloudWeGo netpoll harness is not supported on this platform")
	}
}

func (b *specBuilder) addHTTPS1() {
	for _, version := range tlsVersions {
		for _, size := range payloadsSmall {
			options := &tlsOptions{Version: version, ALPNVariable: "HTTP1_ALPN"}
			b.add(gnalloyCommand(b.profile.GnalloyBackend, protocolHTTPS1, size, options))
			b.add(nettyCommand(b.profile.NettyBackend, protocolHTTPS1, size, options))
			b.add(fasthttpCommand(protocolHTTPS1, size, options))
		}
	}
	b.unsupported(frameworkGnet, protocolHTTPS1, "https1", "gnet benchmark harness intentionally measures the native poller and does not wrap TLS")
	b.unsupported(frameworkNetpoll, protocolHTTPS1, "https1", "netpoll benchmark harness intentionally measures the native poller and does not wrap TLS")
}

func (b *specBuilder) addHTTP2() {
	for _, size := range payloadsSmall {
		b.add(gnalloyCommand(b.profile.GnalloyBackend, protocolHTTP2, size, nil))
		b.add(nettyCommand(b.profile.NettyBackend, protocolHTTP2, size, nil))
	}
	b.unsupported(frameworkGnet, protocolHTTP2, "http2", "gnet does not provide an HTTP/2 codec in this benchmark harness")
	b.unsupported(frameworkFastHTTP, protocolHTTP2, "http2", "fasthttp does not provide an HTTP/2 server implementation")
	b.unsupported(frameworkNetpoll, protocolHTTP2, "http2", "netpoll benchmark harness does not include an HTTP/2 codec")
}

func (b *specBuilder) addHTTPS2() {
	b.unsupported(frameworkGnalloy, protocolHTTPS2, b.profile.GnalloyBackend+" https2 tls11", "HTTP/2 over TLS requires TLS 1.2 or newer")
	b.unsupported(frameworkNetty, protocolHTTPS2, b.profile.NettyBackend+" https2 tls11", "HTTP/2 over TLS requires TLS 1.2 or newer")
	for _, version := range tlsVersions[1:] {
		for _, size := range payloadsSmall {
			options := &tlsOptions{Version: version, ALPNVariable: "HTTP2_ALPN"}
			b.add(gnalloyCommand(b.profile.GnalloyBackend, protocolHTTPS2, size, options))
			b.add(nettyCommand(b.profile.NettyBackend, protocolHTTPS2, size, options))
		}
	}
	b.unsupported(frameworkGnet, protocolHTTPS2, "https2", "gnet benchmark harness does not provide HTTP/2 over TLS")
	b.unsupported(frameworkFastHTTP, protocolHTTPS2, "https2", "fasthttp does not provide an HTTP/2 server implementation")
	b.unsupported(frameworkNetpoll, protocolHTTPS2, "https2", "netpoll benchmark harness does not provide HTTP/2 over TLS")
}

func (b *specBuilder) addHTTP3() {
	for _, size := range payloadsSmall {
		options := &tlsOptions{Version: tlsVersions[2], ALPNVariable: "HTTP3_ALPN"}
		b.add(gnalloyCommand("rfc9000", protocolHTTP3, size, options))
		b.add(nettyCommand(b.profile.NettyBackend, protocolHTTP3, size, options))
	}
	b.unsupported(frameworkGnet, protocolHTTP3, "http3", "gnet does not provide an HTTP/3 over QUIC codec in this benchmark harness")
	b.unsupported(frameworkFastHTTP, protocolHTTP3, "http3", "fasthttp is HTTP/1 oriented and does not provide HTTP/3")
	b.unsupported(frameworkNetpoll, protocolHTTP3, "http3", "netpoll benchmark harness does not provide HTTP/3 over QUIC")
}

func (b *specBuilder) addQUICStream() {
	for _, size := range payloadsSmall {
		options := &tlsOptions{Version: tlsVersions[2], ALPNVariable: "QUIC_ALPN"}
		b.add(gnalloyCommand("rfc9000", protocolQUICStream, size, options))
		b.add(nettyCommand(b.profile.NettyBackend, protocolQUICStream, size, options))
	}
	b.unsupported(frameworkGnet, protocolQUICStream, "quic stream", "gnet does not provide a QUIC transport")
	b.unsupported(frameworkFastHTTP, protocolQUICStream, "quic stream", "fasthttp does not provide a QUIC transport")
	b.unsupported(frameworkNetpoll, protocolQUICStream, "quic stream", "netpoll does not provide a QUIC transport")
}

func (b *specBuilder) add(scenario parity.Scenario) {
	if !scenario.Skip && b.profile.Timeout > 0 {
		scenario.Timeout = parity.Duration(b.profile.Timeout)
	}
	b.scenarios = append(b.scenarios, scenario)
}

func (b *specBuilder) unsupported(framework string, protocol string, label string, reason string) {
	if !b.profile.IncludeUnsupported {
		return
	}
	b.add(parity.Scenario{
		Name:       fmt.Sprintf("%s %s unsupported", framework, label),
		Framework:  framework,
		Protocol:   protocol,
		Backend:    "unsupported",
		Skip:       true,
		SkipReason: reason,
		Tags:       []string{b.profile.Platform, framework, protocol, "unsupported"},
	})
}

type tlsOptions struct {
	Version      tlsVersion
	ALPNVariable string
}

func gnalloyCommand(backend string, protocol string, size payload, tls *tlsOptions) parity.Scenario {
	command := []string{"${GNALLOY_BENCH}", "-protocol", protocol}
	if backend != "rfc9000" {
		command = append(command, "-backend", backend)
	}
	if protocol == protocolHTTP1 || protocol == protocolHTTPS1 {
		command = append(command, "-read-buffer-size", "${GNALLOY_HTTP1_READ_BUFFER}")
	}
	command = appendBaseFlags(command, size, "-")
	command = appendTLSFlags(command, tls, "-")
	return executableScenario(frameworkGnalloy, protocol, backend, size, command, tls, "parity-harness")
}

func nettyCommand(backend string, protocol string, size payload, tls *tlsOptions) parity.Scenario {
	command := []string{"java", "-jar", "${NETTY_BENCH_JAR}", "--protocol", protocol, "--backend", backend}
	command = appendBaseFlags(command, size, "--")
	command = appendTLSFlags(command, tls, "--")
	return executableScenario(frameworkNetty, protocol, backend, size, command, tls, "external")
}

func gnetCommand(protocol string, size payload) parity.Scenario {
	command := []string{"${GNET_BENCH}", "-protocol", protocol}
	command = appendBaseFlags(command, size, "-")
	return executableScenario(frameworkGnet, protocol, "poller", size, command, nil, "external")
}

func fasthttpCommand(protocol string, size payload, tls *tlsOptions) parity.Scenario {
	command := []string{"${FASTHTTP_BENCH}", "-protocol", protocol}
	command = appendBaseFlags(command, size, "-")
	command = appendTLSFlags(command, tls, "-")
	return executableScenario(frameworkFastHTTP, protocol, "net", size, command, tls, "external")
}

func netpollCommand(protocol string, size payload) parity.Scenario {
	command := []string{"${NETPOLL_BENCH}", "-protocol", protocol}
	command = appendBaseFlags(command, size, "-")
	return executableScenario(frameworkNetpoll, protocol, "poller", size, command, nil, "external")
}

func appendBaseFlags(command []string, size payload, prefix string) []string {
	return append(command,
		prefix+"payload", size.Bytes,
		prefix+"connections", "${CONNECTIONS}",
		prefix+"messages", "${MESSAGES}",
		prefix+"latency-sample-rate", "${LATENCY_SAMPLE_RATE}",
		prefix+"warmup-messages", "${WARMUP_MESSAGES}",
		prefix+"timeout", "${COMMAND_TIMEOUT}",
	)
}

func appendTLSFlags(command []string, tls *tlsOptions, prefix string) []string {
	if tls == nil || tls.Version.Value == "" {
		return command
	}
	name := prefix + "tls-version"
	alpn := prefix + "alpn"
	cipher := prefix + "cipher-suites"
	command = append(command, alpn, "${"+tls.ALPNVariable+"}", name, tls.Version.Value)
	if tls.Version.CipherVariable != "" {
		command = append(command, cipher, "${"+tls.Version.CipherVariable+"}")
	}
	return command
}

func executableScenario(framework string, protocol string, backend string, size payload, command []string, tls *tlsOptions, firstTag string) parity.Scenario {
	name := fmt.Sprintf("%s %s %s %s", framework, backend, displayProtocol(protocol, tls), size.Label)
	if framework == frameworkFastHTTP {
		name = fmt.Sprintf("%s %s %s", framework, displayProtocol(protocol, tls), size.Label)
	}
	tags := []string{firstTag, framework, protocol, payloadTag(size)}
	if tls != nil && tls.Version.Tag != "" {
		tags = append(tags, tls.Version.Tag)
	}
	if protocol == protocolHTTP3 || protocol == protocolQUICStream {
		tags = append(tags, "quic")
	}
	return parity.Scenario{
		Name:      name,
		Framework: framework,
		Protocol:  protocol,
		Backend:   backend,
		Payload:   size.Label,
		Warmup:    1,
		Repeat:    3,
		WorkDir:   "${REPO_ROOT}",
		Command:   command,
		Timeout:   parity.Duration(5 * time.Minute),
		Tags:      tags,
	}
}

func displayProtocol(protocol string, tls *tlsOptions) string {
	switch {
	case tls != nil && tls.Version.Tag != "":
		return protocol + " " + tls.Version.Tag
	case protocol == protocolQUICStream:
		return "quic stream"
	default:
		return protocol
	}
}

func payloadTag(size payload) string {
	return "payload-" + strings.TrimSuffix(strings.ToLower(size.Label), "b")
}

func formatCommandTimeout(timeout time.Duration) string {
	switch {
	case timeout <= 0:
		return "15m"
	case timeout%time.Minute == 0:
		return strconv.FormatInt(int64(timeout/time.Minute), 10) + "m"
	case timeout%time.Second == 0:
		return strconv.FormatInt(int64(timeout/time.Second), 10) + "s"
	default:
		millis := timeout / time.Millisecond
		if millis <= 0 {
			millis = 1
		}
		return strconv.FormatInt(int64(millis), 10) + "ms"
	}
}
