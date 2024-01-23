/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var timeFormats = map[string]string{
	"ANSIC":       `Mon Jan _2 15:04:05 2006`,
	"DateOnly":    `2006-01-02`,
	"DateTime":    `2006-01-02 15:04:05`,
	"Kitchen":     `3:04PM`,
	"Layout":      `01/02 03:04:05PM '06 -0700`,
	"RFC1123":     `Mon, 02 Jan 2006 15:04:05 MST`,
	"RFC1123Z":    `Mon, 02 Jan 2006 15:04:05 -0700`,
	"RFC3339":     `2006-01-02T15:04:05Z07:00`,
	"RFC3339Nano": `2006-01-02T15:04:05.999999999Z07:00`,
	"RFC822":      `02 Jan 06 15:04 MST`,
	"RFC822Z":     `02 Jan 06 15:04 -0700`,
	"RFC850":      `Monday, 02-Jan-06 15:04:05 MST`,
	"RubyDate":    `Mon Jan 02 15:04:05 -0700 2006`,
	"Stamp":       `Jan _2 15:04:05`,
	"StampMicro":  `Jan _2 15:04:05.000000`,
	"StampMilli":  `Jan _2 15:04:05.000`,
	"StampNano":   `Jan _2 15:04:05.000000000`,
	"TimeOnly":    `15:04:05`,
	"UnixDate":    `Mon Jan _2 15:04:05 MST 2006`,
}

func serveTime(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		var format string = ""

		requestedFormat := r.URL.Query().Get("format")
		if requestedFormat == "" {
			requestedFormat = "RFC822"
		}

		for k, v := range timeFormats {
			if strings.EqualFold(requestedFormat, k) {
				format = v

				break
			}
		}

		if format == "" {
			format = timeFormats["RFC822"]
		}

		adjustedStartTime := startTime

		location := strings.TrimPrefix(p.ByName("time"), "/") + p.ByName("rest")

		tz, err := time.LoadLocation(location)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			http.Redirect(w, r, "/time/", redirectStatusCode)
		} else {
			adjustedStartTime = adjustedStartTime.In(tz)
		}

		w.Header().Set("Content-Type", "text/plain")

		_, err = w.Write([]byte(adjustedStartTime.Format(format) + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested the current time for %q in %s format\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				location,
				requestedFormat)
		}
	}
}

func registerTime(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "time"

	mux.GET("/time/:time", serveTime(errorChannel))
	mux.GET("/time/:time/*rest", serveTime(errorChannel))
	mux.GET("/time/", serveUsage(module, usage, errorChannel))

	usage.Store(module, []string{
		"/time/America/Chicago",
		"/time/EST",
		"/time/UTC?format=kitchen",
	})
}
