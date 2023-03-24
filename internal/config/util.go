package config

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const xForwardedProto = "X-Forwarded-Proto"

// FormatBaseAddress returns the public base address of the server given the server config
// it ignores the port if it's default for the given protocol
func (c *Config) FormatBaseAddress(r *http.Request) string {

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

func (c *Config) FormatBaseAddressStringBuilder(r *http.Request) string {

	var b strings.Builder

	// We want to use all Host-sourced configs, or none, so check the Host header is non-empty too.
	if r != nil && r.Host != "" && r.Header.Get(xForwardedProto) != "" {
		b.WriteString(r.Header.Get(xForwardedProto))
	} else if c.ProtoTLS {
		b.WriteString("https")
	} else if r != nil && r.URL != nil && r.URL.Scheme != "" {
		b.WriteString(r.URL.Scheme)
	} else {
		b.WriteString("http")
	}

	if r != nil && r.Host != "" { // Host headers override all local config
		b.WriteString("://")
		b.WriteString(r.Host)
	} else if c.Port > 0 && ((c.ProtoTLS && c.Port != 443) || (!c.ProtoTLS && c.Port != 80)) { // else config settings apply for non-standard ports
		b.WriteString("://")
		b.WriteString(c.Address)
		b.WriteRune(':')
		b.WriteString(strconv.Itoa(c.Port))
	} else {
		b.WriteString("://")
		b.WriteString(c.Address)
	}

	b.WriteRune('/')

	return b.String()
}

type timeString struct {
	b strings.Builder
}

func (ts *timeString) append(word string, number uint64) {

	if number <= 0 {
		return
	}

	if ts.b.Len() != 0 {
		ts.b.WriteString(", ")
	}

	//fmt.Sprintf("%d %s", number, word)
	ts.b.WriteString(strconv.FormatUint(number, 10))
	ts.b.WriteString(" ")
	ts.b.WriteString(word)
	if number > 1 {
		ts.b.WriteString("s")
	}
}

// FormatTime provides a human readable time string given a time.Duration
func FormatTime(t time.Duration) string {
	if t <= 0 {
		return "∞"
	}

	var seconds, secondsCarried float64
	if t.Seconds() < 60 {
		secondsCarried = t.Seconds()
	} else {
		secondsCarried = math.Mod(t.Seconds(), 60)
		seconds = t.Seconds() - secondsCarried
	}

	mins := seconds / 60
	minsCarried := math.Mod(seconds/60, 60)

	hours := (mins - minsCarried) / 60
	hoursCarried := math.Mod(hours, 24)

	days := (hours - hoursCarried) / 24
	daysCarried := math.Mod(days, 365)

	years := (days - daysCarried) / 365

	ts := timeString{}
	ts.append("year", uint64(years))
	ts.append("day", uint64(daysCarried))
	ts.append("hour", uint64(hoursCarried))
	ts.append("minute", uint64(minsCarried))
	ts.append("second", uint64(secondsCarried))

	return ts.b.String()
}
