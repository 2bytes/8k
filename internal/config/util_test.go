package config

import (
	"strings"
	"testing"
)

const (
	address = "127.0.0.1"

	baseHTTPSDefaultPort    = "https://" + address + "/"
	baseHTTPDefaultPort     = "http://" + address + "/"
	baseHTTPSPort80         = "https://" + address + ":80/"
	baseHTTPPort443         = "http://" + address + ":443/"
	baseHTTPSNonDefaultPort = "https://" + address + ":3456/"
	baseHTTPNonDefaultPort  = "http://" + address + ":3456/"
)

func TestFormatBaseAddressHTTPSDefaultPort(t *testing.T) {

	tc := &Config{ProtoTLS: true, Address: address, Port: 443}

	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPSDefaultPort) != 0 {
		t.Error(baseHTTPSDefaultPort, base)
	}
}

func TestFormatBaseAddressHTTPDefaultPort(t *testing.T) {

	tc := &Config{ProtoTLS: false, Address: address, Port: 80}
	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPDefaultPort) != 0 {
		t.Error(baseHTTPDefaultPort, base)
	}
}

func TestFormatBaseAddressHTTPSDefaultHTTPPort(t *testing.T) {
	tc := &Config{ProtoTLS: true, Address: address, Port: 80}
	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPSPort80) != 0 {
		t.Error(baseHTTPSPort80, base)
	}
}

func TestFormatBaseAddressHTTPDefaultHTTPSPort(t *testing.T) {
	tc := &Config{ProtoTLS: false, Address: address, Port: 443}
	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPPort443) != 0 {
		t.Error(baseHTTPSPort80, base)
	}
}

func TestFormatBaseAddressHTTPSOtherPort(t *testing.T) {
	tc := &Config{ProtoTLS: true, Address: address, Port: 3456}
	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPSNonDefaultPort) != 0 {
		t.Error(baseHTTPSNonDefaultPort, base)
	}
}

func TestFormatBaseAddressHTTPOtherPort(t *testing.T) {
	tc := &Config{ProtoTLS: false, Address: address, Port: 3456}
	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPNonDefaultPort) != 0 {
		t.Error(baseHTTPNonDefaultPort, base)
	}
}

func TestFormatBaseAddressBadPort(t *testing.T) {
	tc := &Config{ProtoTLS: true, Address: address, Port: 0}
	base := tc.FormatBaseAddress()

	if strings.Compare(base, baseHTTPSDefaultPort) != 0 {
		t.Error(baseHTTPSDefaultPort, base)
	}
}
