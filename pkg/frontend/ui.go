package frontend

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

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

//go:embed tmpl/index.html
var IndexFS embed.FS
var IndexFile = "tmpl/index.html"

func init() {
	_, err := template.ParseFS(IndexFS, IndexFile)
	if err != nil {
		fmt.Printf("Failed to load embedded html: %s\n", err)
		os.Exit(1)
	}
}
