package matrix

import (
	"strings"
	"testing"
	"time"

	"gnalloy.org/benchmarks/parity"
)

func TestLinuxFullSpecCoversRequestedProtocolsAndFrameworks(t *testing.T) {
	spec := LinuxFullSpec()
	if err := spec.Validate(); err != nil {
		t.Fatal(err)
	}
	requireScenario(t, spec, "gnalloy epoll tcp-echo 16KiB", false)
	requireScenario(t, spec, "netty epoll udp-echo 1KiB", false)
	requireScenario(t, spec, "fasthttp https1 tls12 1KiB", false)
	requireScenario(t, spec, "gnalloy epoll https2 tls13 1KiB", false)
	requireScenario(t, spec, "netty epoll http3 tls13 1KiB", false)
	requireScenario(t, spec, "gnalloy rfc9000 quic-stream tls13 1KiB", false)
	requireScenario(t, spec, "netpoll udp echo unsupported", true)
	netty := requireScenario(t, spec, "netty epoll tcp-echo 1KiB", false)
	if command := strings.Join(netty.Command, " "); !strings.Contains(command, "--payload 1024") || strings.Contains(command, " -payload ") {
		t.Fatalf("netty command=%q, want double-dash benchmark flags", command)
	}

	frameworks := map[string]bool{}
	protocols := map[string]bool{}
	tlsVersions := map[string]bool{}
	for _, scenario := range spec.Scenarios {
		frameworks[scenario.Framework] = true
		protocols[scenario.Protocol] = true
		for _, arg := range scenario.Command {
			switch arg {
			case "1.1", "1.2", "1.3":
				tlsVersions[arg] = true
			}
		}
	}
	for _, framework := range []string{"gnalloy", "netty", "gnet", "fasthttp", "netpoll"} {
		if !frameworks[framework] {
			t.Fatalf("framework %q missing", framework)
		}
	}
	for _, protocol := range []string{"tcp-echo", "udp-echo", "http1", "https1", "http2", "https2", "http3", "quic-stream"} {
		if !protocols[protocol] {
			t.Fatalf("protocol %q missing", protocol)
		}
	}
	for _, version := range []string{"1.1", "1.2", "1.3"} {
		if !tlsVersions[version] {
			t.Fatalf("tls version %q missing", version)
		}
	}
}

func TestLinuxFullSpecRejectsIllegalTLSCombinationsBySkipping(t *testing.T) {
	spec := LinuxFullSpec()
	for _, scenario := range spec.Scenarios {
		command := strings.Join(scenario.Command, " ")
		if scenario.Protocol == "https2" && strings.Contains(command, " 1.1") {
			t.Fatalf("HTTP/2 over TLS 1.1 must not be executable: %+v", scenario)
		}
		if (scenario.Protocol == "http3" || scenario.Protocol == "quic-stream") && strings.Contains(command, " 1.1") {
			t.Fatalf("QUIC protocol must not use TLS 1.1: %+v", scenario)
		}
		if (scenario.Protocol == "http3" || scenario.Protocol == "quic-stream") && strings.Contains(command, " 1.2") {
			t.Fatalf("QUIC protocol must not use TLS 1.2: %+v", scenario)
		}
	}
	requireScenario(t, spec, "gnalloy epoll https2 tls11 unsupported", true)
	requireScenario(t, spec, "fasthttp http2 unsupported", true)
	requireScenario(t, spec, "gnet quic stream unsupported", true)
}

func TestWindowsFullSpecUsesWindowsBackends(t *testing.T) {
	spec := WindowsFullSpec()
	if err := spec.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, scenario := range spec.Scenarios {
		if scenario.Skip {
			continue
		}
		if scenario.Framework == "gnalloy" && scenario.Protocol != "http3" && scenario.Protocol != "quic-stream" && scenario.Backend != "iocp" {
			t.Fatalf("gnalloy windows backend=%q scenario=%+v", scenario.Backend, scenario)
		}
		if scenario.Framework == "netty" && scenario.Backend != "nio" {
			t.Fatalf("netty windows backend=%q scenario=%+v", scenario.Backend, scenario)
		}
	}
	requireScenario(t, spec, "netpoll tcp echo unsupported", true)
}

func TestFullSpecAppliesProfileTimeout(t *testing.T) {
	spec := FullSpec(Profile{
		Name:               "timeout test",
		Timeout:            17 * time.Second,
		CommandTimeout:     9 * time.Second,
		NetpollSupported:   true,
		IncludeUnsupported: true,
	})
	scenario := requireScenario(t, spec, "gnalloy epoll tcp-echo 64B", false)
	if got := time.Duration(scenario.Timeout); got != 17*time.Second {
		t.Fatalf("timeout=%s, want 17s", got)
	}
	if got := spec.Variables["COMMAND_TIMEOUT"]; got != "9s" {
		t.Fatalf("command timeout=%q, want 9s", got)
	}
	command := strings.Join(scenario.Command, " ")
	if !strings.Contains(command, "-timeout ${COMMAND_TIMEOUT}") {
		t.Fatalf("command=%q, want command timeout flag", command)
	}
}

func requireScenario(t *testing.T, spec parity.Spec, name string, skipped bool) parity.Scenario {
	t.Helper()
	for _, scenario := range spec.Scenarios {
		if scenario.Name != name {
			continue
		}
		if scenario.Skip != skipped {
			t.Fatalf("scenario %q skip=%t, want %t", name, scenario.Skip, skipped)
		}
		return scenario
	}
	t.Fatalf("scenario %q missing", name)
	return parity.Scenario{}
}
