package frontend

import "fmt"

// UI defined the variables for templating into the frontend
type UI struct {
	Title        string
	AccentColour string
	ProtoTLS     bool
	Address      string
	Port         int
}

// ProtoString returns the protocol string for the configured UI mode
func (ui *UI) ProtoString() string {
	if ui.ProtoTLS {
		return "https"
	}
	return "http"
}

// Data contains templating data inserted into the index page request
type Data struct {
	Title        string
	AccentColour string
	MaxBytes     int
	Proto        string
	Address      string
	Port         string
	RandomPath   string
}

// BaseAddress creates the public address string for use in the UI
// This may differ from the bind address if behind a reverse proxy or NAT (e.g. Docker)
func (ui *UI) BaseAddress() string {

	if ui.ProtoTLS && ui.Port != 443 {
		return fmt.Sprintf("%s://%s:%d/", ui.ProtoString(), ui.Address, ui.Port)
	} else if !ui.ProtoTLS && ui.Port != 80 {
		return fmt.Sprintf("%s://%s:%d/", ui.ProtoString(), ui.Address, ui.Port)
	} else {
		return fmt.Sprintf("%s://%s/", ui.ProtoString(), ui.Address)
	}
}
