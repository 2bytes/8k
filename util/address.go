package util

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/2bytes/8k/internal/config"
)

const xForwardedProto = "X-Forwarded-Proto"

// FormatBaseAddress returns the public base address of the server given the server config
// it ignores the port if it's default for the given protocol
func FormatBaseAddress(c *config.Config, r *http.Request) string {

	scheme := "http"
	host := c.Address

	// We want to use all Host-sourced configs, or none, so check the Host header is non-empty too.
	if r != nil && r.Host != "" && r.Header.Get(xForwardedProto) != "" {
		scheme = r.Header.Get(xForwardedProto)
	} else if c.ProtoTLS {
		scheme = "https"
	} else if r != nil && r.URL != nil && r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}

	if r != nil && r.Host != "" { // Host headers override all local config
		host = r.Host
	} else if c.Port > 0 && ((c.ProtoTLS && c.Port != 443) || (!c.ProtoTLS && c.Port != 80)) { // else config settings apply for non-standard ports
		host = fmt.Sprintf("%s:%d", c.Address, c.Port)
	}

	url, err := url.Parse(fmt.Sprintf("%s://%s/", scheme, host))
	if err != nil {
		log.Println("failed to parse URL:", err)
		return ""
	}

	return url.String()
}
