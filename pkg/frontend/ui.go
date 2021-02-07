package frontend

// Data contains templating data inserted into the index page request
type Data struct {
	Title        string
	AccentColour string
	MaxBytes     int
	MaxItems     int
	TTL          string
	BaseAddress  string
	RandomPath   string
	Version      string
}
