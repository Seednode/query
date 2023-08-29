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

func serveTime(TimeFormats map[string]string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		startTime := time.Now()

		var format string = ""

		requestedFormat := r.URL.Query().Get("format")

		if requestedFormat != "" {
			for k, v := range TimeFormats {
				if strings.EqualFold(requestedFormat, k) {
					format = v

					break
				}
			}
		}

		if format == "" {
			format = TimeFormats["RFC822"]
		}

		tz, err := time.LoadLocation(strings.TrimPrefix(p[0].Value, "/"))

		if err != nil {
			http.Redirect(w, r, "/time/", RedirectStatusCode)
		} else {
			startTime = startTime.In(tz)
		}

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(startTime.Format(format) + "\n"))

		if verbose {
			fmt.Printf("%s | %s checked the time\n",
				startTime.Format(LogDate),
				realIP(r, true))
		}
	}
}
