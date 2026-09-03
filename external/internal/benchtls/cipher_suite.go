package benchtls

import (
	"crypto/tls"
	"fmt"
	"strings"
)

// ParseCipherSuites 将 Go 标准密码套件名称解析为配置标识。
func ParseCipherSuites(value string) ([]uint16, []string, error) {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == ':'
	})
	if len(parts) == 0 {
		return nil, nil, nil
	}
	suites := append(tls.CipherSuites(), tls.InsecureCipherSuites()...)
	ids := make([]uint16, 0, len(parts))
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.ToUpper(strings.TrimSpace(part))
		if name == "" {
			continue
		}
		found := false
		for _, suite := range suites {
			if suite.Name == name {
				ids = append(ids, suite.ID)
				names = append(names, suite.Name)
				found = true
				break
			}
		}
		if !found {
			return nil, nil, fmt.Errorf("benchtls: unsupported cipher suite %q", part)
		}
	}
	return ids, names, nil
}
