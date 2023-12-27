/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveHttpStatusCode(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		var text string

		trimmed := strings.TrimSuffix(strings.TrimPrefix(p[0].Value, "/"), "/")

		value, err := strconv.Atoi(trimmed)
		if err != nil {
			text = ""
		} else {
			text = http.StatusText(value)
		}

		if text == "" {
			w.WriteHeader(http.StatusBadRequest)

			w.Write([]byte("Invalid status code requested.\n"))

			if verbose {
				fmt.Printf("%s | %s requested invalid status code %q\n",
					startTime.Format(timeFormats["RFC3339"]),
					realIP(r, true),
					trimmed)
			}
		} else {
			w.WriteHeader(value)

			w.Write([]byte(text + "\n"))

			if verbose {
				fmt.Printf("%s | %s requested status code %q\n",
					startTime.Format(timeFormats["RFC3339"]),
					realIP(r, true),
					trimmed)
			}
		}
	}
}

func registerHttpStatusHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/http/status/*status", serveHttpStatusCode(errorChannel))

	var usage []string
	usage = append(usage, "/http/status/200")
	usage = append(usage, "/http/status/404")
	usage = append(usage, "/http/status/500")

	return usage
}
