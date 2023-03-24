package util

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/2bytes/8k/internal/config"
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

	tc := &config.Config{ProtoTLS: true, Address: address, Port: 443}

	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPSDefaultPort) != 0 {
		t.Error(baseHTTPSDefaultPort, base)
	}
}

func TestFormatBaseAddressHTTPDefaultPort(t *testing.T) {

	tc := &config.Config{ProtoTLS: false, Address: address, Port: 80}
	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPDefaultPort) != 0 {
		t.Error(baseHTTPDefaultPort, base)
	}
}

func TestFormatBaseAddressHTTPSDefaultHTTPPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: true, Address: address, Port: 80}
	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPSPort80) != 0 {
		t.Error(baseHTTPSPort80, base)
	}
}

func TestFormatBaseAddressHTTPDefaultHTTPSPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: false, Address: address, Port: 443}
	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPPort443) != 0 {
		t.Error(baseHTTPSPort80, base)
	}
}

func TestFormatBaseAddressHTTPSOtherPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: true, Address: address, Port: 3456}
	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPSNonDefaultPort) != 0 {
		t.Error(baseHTTPSNonDefaultPort, base)
	}
}

func TestFormatBaseAddressHTTPOtherPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: false, Address: address, Port: 3456}
	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPNonDefaultPort) != 0 {
		t.Error(baseHTTPNonDefaultPort, base)
	}
}

func TestFormatBaseAddressBadPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: true, Address: address, Port: 0}
	base := FormatBaseAddress(tc, nil)

	if strings.Compare(base, baseHTTPSDefaultPort) != 0 {
		t.Error(baseHTTPSDefaultPort, base)
	}
}

func TestFormatBaseAddressRequestHost(t *testing.T) {
	tc := &config.Config{ProtoTLS: false, Address: "", Port: 0}
	base := FormatBaseAddress(tc, &http.Request{URL: &url.URL{Scheme: "http"}, Host: address, Header: map[string][]string{"X-Forwarded-Proto": {"https"}}})

	if strings.Compare(base, baseHTTPSDefaultPort) != 0 {
		t.Error(baseHTTPSDefaultPort, base)
	}
}

func TestFormatBaseAddressRequestHostHTTP(t *testing.T) {
	tc := &config.Config{ProtoTLS: false, Address: "", Port: 0}
	base := FormatBaseAddress(tc, &http.Request{URL: &url.URL{Scheme: "invalid"}, Host: address, Header: map[string][]string{"X-Forwarded-Proto": {"http"}}})

	if strings.Compare(base, baseHTTPDefaultPort) != 0 {
		t.Error(baseHTTPDefaultPort, base)
	}
}

func TestFormatBaseAddressRequestHostNonDefaultPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: false, Address: "", Port: 0}
	base := FormatBaseAddress(tc, &http.Request{URL: &url.URL{Scheme: "invalid"}, Host: address + ":3456", Header: map[string][]string{"X-Forwarded-Proto": {"https"}}})

	if strings.Compare(base, baseHTTPSNonDefaultPort) != 0 {
		t.Error(baseHTTPSNonDefaultPort, base)
	}
}

func TestFormatBaseAddressRequestHostHTTPNonDefaultPort(t *testing.T) {
	tc := &config.Config{ProtoTLS: false, Address: "", Port: 0}
	base := FormatBaseAddress(tc, &http.Request{URL: &url.URL{Scheme: "invalid"}, Host: address + ":3456", Header: map[string][]string{"X-Forwarded-Proto": {"http"}}})

	if strings.Compare(base, baseHTTPNonDefaultPort) != 0 {
		t.Error(baseHTTPNonDefaultPort, base)
	}
}

func BenchmarkFormatBaseAddress(b *testing.B) {
	tc := &config.Config{ProtoTLS: false, Address: "", Port: 0}
	base := FormatBaseAddress(tc, &http.Request{URL: &url.URL{Scheme: "https"}, Host: address + ":3456"})

	for i := 0; i < b.N; i++ {
		if strings.Compare(base, baseHTTPSNonDefaultPort) != 0 {
			b.Errorf("FormatBaseAddress() = '%v', want '%v'", base, baseHTTPSNonDefaultPort)
		}
	}
}
