package benchtls

import (
	"crypto/tls"
	"testing"
)

func TestParseCipherSuitesSupportsLegacyAndModernSuites(t *testing.T) {
	ids, names, err := ParseCipherSuites("TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA || ids[1] != tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 {
		t.Fatalf("ids=%v", ids)
	}
	if len(names) != 2 || names[0] != "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA" || names[1] != "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" {
		t.Fatalf("names=%v", names)
	}
}

func TestParseCipherSuitesRejectsUnknownSuite(t *testing.T) {
	if _, _, err := ParseCipherSuites("TLS_UNKNOWN"); err == nil {
		t.Fatal("expected unknown cipher suite rejection")
	}
}
