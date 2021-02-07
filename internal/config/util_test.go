package config

import (
	"strings"
	"testing"
	"time"
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

func TestFormatTime(t *testing.T) {
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Zero",
			args: args{time.Second * 0},
			want: "∞",
		},
		{
			name: "1 Second",
			args: args{time.Second},
			want: "1 second",
		},
		{
			name: "59 Seconds",
			args: args{time.Second * 59},
			want: "59 seconds",
		},
		{
			name: "1 Minute",
			args: args{time.Minute},
			want: "1 minute",
		},
		{
			name: "59 Minutes",
			args: args{time.Minute * 59},
			want: "59 minutes",
		},
		{
			name: "1 Hour",
			args: args{time.Minute * 60},
			want: "1 hour",
		},
		{
			name: "23 Hours",
			args: args{time.Minute * 60 * 23},
			want: "23 hours",
		},
		{
			name: "1 Day",
			args: args{time.Minute * 60 * 24},
			want: "1 day",
		},

		{
			name: "364 days",
			args: args{time.Minute * 60 * 24 * 364},
			want: "364 days",
		},
		{
			name: "1 Year",
			args: args{time.Minute * 60 * 24 * 365},
			want: "1 year",
		},

		{
			name: "2 Years",
			args: args{time.Minute * 60 * 24 * 365 * 2},
			want: "2 years",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTime(tt.args.t); got != tt.want {
				t.Errorf("FormatTime() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func BenchmarkGetTimeStringLong(b *testing.B) {

	testDuration := (time.Second * 1 * 59) + // 59 seconds
		(time.Minute * 59) + // 59 minutes
		(time.Minute * 60 * 23) + // 23 hours
		(time.Minute * 60 * 24 * 364) + // 364 days
		(time.Minute * 60 * 24 * 365 * 9) // 9 years

	want := "9 years, 364 days, 23 hours, 59 minutes, 59 seconds"

	for i := 0; i < b.N; i++ {
		if got := FormatTime(testDuration); got != want {
			b.Errorf("FormatTime() = '%v', want '%v'", got, want)
		}
	}
}
