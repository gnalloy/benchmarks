package tlsbench

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"

	gnalloytls "gnalloy.org/handler-tls"
)

const (
	Version11 = "1.1"
	Version12 = "1.2"
	Version13 = "1.3"
)

var ErrInvalidConfig = errors.New("tlsbench: invalid config")

// ClientOptions 描述跨框架基准使用的 TLS 客户端参数。
type ClientOptions struct {
	Enabled                   bool
	ServerName                string
	Version                   string
	ALPN                      string
	CipherSuites              string
	AllowInsecureCipherSuites bool
	InsecureSkipVerify        bool
}

// Build 构造精确锁定协议版本的客户端配置。
func (o ClientOptions) Build(addr string) (*tls.Config, error) {
	if !o.Enabled {
		return nil, nil
	}
	version, _, err := parseVersion(o.Version)
	if err != nil {
		return nil, err
	}
	if version == tls.VersionTLS13 && strings.TrimSpace(o.CipherSuites) != "" {
		return nil, fmt.Errorf("%w: TLS 1.3 cipher suites are runtime-managed", ErrInvalidConfig)
	}
	cipherSuites, err := gnalloytls.ParseCipherSuites(o.CipherSuites, gnalloytls.CipherSuiteOptions{
		IncludeInsecure: o.AllowInsecureCipherSuites,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	if err := gnalloytls.ValidateConfigurableCipherSuites(cipherSuites); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	serverName, err := resolveServerName(o.ServerName, addr)
	if err != nil {
		return nil, err
	}
	nextProtos := splitNames(o.ALPN)
	if len(nextProtos) == 0 {
		return nil, fmt.Errorf("%w: empty ALPN", ErrInvalidConfig)
	}
	return &tls.Config{
		ServerName:         serverName,
		MinVersion:         version,
		MaxVersion:         version,
		CipherSuites:       cipherSuites,
		NextProtos:         nextProtos,
		InsecureSkipVerify: o.InsecureSkipVerify,
	}, nil
}

// NormalizeVersion 返回标准化的 TLS 版本名称。
func NormalizeVersion(value string) (string, error) {
	_, name, err := parseVersion(value)
	return name, err
}

// CipherSuiteNames 返回稳定的 IANA cipher suite 名称。
func CipherSuiteNames(ids []uint16) string {
	return strings.Join(gnalloytls.CipherSuiteNames(ids), ",")
}

func parseVersion(value string) (uint16, string, error) {
	version := strings.TrimSpace(value)
	if version == "" {
		version = Version13
	}
	switch version {
	case Version11:
		return tls.VersionTLS11, version, nil
	case Version12:
		return tls.VersionTLS12, version, nil
	case Version13:
		return tls.VersionTLS13, version, nil
	default:
		return 0, "", fmt.Errorf("%w: unsupported TLS version %q", ErrInvalidConfig, value)
	}
}

func resolveServerName(configured, addr string) (string, error) {
	if serverName := strings.TrimSpace(configured); serverName != "" {
		return serverName, nil
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("%w: parse address: %v", ErrInvalidConfig, err)
	}
	return host, nil
}

func splitNames(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';'
	})
	names := parts[:0]
	for _, part := range parts {
		if name := strings.TrimSpace(part); name != "" {
			names = append(names, name)
		}
	}
	return names
}
