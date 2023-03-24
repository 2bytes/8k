package util

import (
	"fmt"
	"net/http"
	"time"
)

type statusCoder struct {
	http.ResponseWriter
	statusCode int
}

func (sc *statusCoder) WriteHeader(code int) {
	sc.ResponseWriter.WriteHeader(code)
	sc.statusCode = code
}

func (sc *statusCoder) Write(b []byte) (int, error) {
	out, err := sc.ResponseWriter.Write(b)
	return out, err
}

// StatusCoder logs status code of request after it completes
func StatusCoder(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		scUtil := &statusCoder{w, http.StatusOK}
		h(scUtil, r)

		// <start> <status> <elapsed> <method> <path> <proto> <user-agent>
		fmt.Printf("%s %3d %7s %4s %s %s %s\n", start.Format("02/Jan/2006:15:04:05 -0700"), scUtil.statusCode, time.Since(start).Truncate(time.Microsecond), r.Method, r.URL.Path, r.Proto, r.UserAgent())
	})
}
