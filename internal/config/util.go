package config

import "fmt"

// FormatBaseAddress returns the public base address of the server given the server config
// it ignores the port if it's default for the given protocol
func (c *Config) FormatBaseAddress() string {

	base := "http"

	if c.ProtoTLS {
		base += "s"
	}

	base += "://" + c.Address

	if c.Port > 0 && ((c.ProtoTLS && c.Port != 443) || (!c.ProtoTLS && c.Port != 80)) {
		base += fmt.Sprintf(":%d", c.Port)
	}

	return base + "/"
}
