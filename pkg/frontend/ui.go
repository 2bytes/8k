package frontend

import (
	"time"
)

// Data contains templating data inserted into the index page request
type Data struct {
	Title        string
	AccentColour string
	MaxBytes     int
	MaxItems     int
	TTL          time.Duration
	BaseAddress  string
	RandomPath   string
}
