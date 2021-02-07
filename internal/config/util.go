package config

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

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

	var seconds, secondsCarried float64 = 0, 0
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
