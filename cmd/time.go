/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strings"
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

func serveTime(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		startTime := time.Now()

		var format string = ""

		requestedFormat := r.URL.Query().Get("format")

		if requestedFormat != "" {
			for k, v := range timeFormats {
				if strings.EqualFold(requestedFormat, k) {
					format = v

					break
				}
			}
		}

		if format == "" {
			format = timeFormats["RFC822"]
		}

		tz, err := time.LoadLocation(strings.TrimPrefix(p[0].Value, "/"))
		if err != nil {
			errorChannel <- err

			http.Redirect(w, r, "/time/", redirectStatusCode)
		} else {
			startTime = startTime.In(tz)
		}

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(startTime.Format(format) + "\n"))

		if verbose {
			fmt.Printf("%s | %s requested the current time\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerTimeHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/time/*time", serveTime(errorChannel))

	var usage []string
	usage = append(usage, "/time/America/Chicago")
	usage = append(usage, "/time/EST")
	usage = append(usage, "/time/UTC?format=kitchen")

	return usage
}
