/*
Copyright © 2024 Seednode <seednode@seedno.de>
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

func serveHTTPStatusCode(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		var text string = ""

		trimmed := strings.TrimSuffix(strings.TrimPrefix(p.ByName("status"), "/"), "/")

		value, err := strconv.Atoi(trimmed)
		if err == nil {
			text = http.StatusText(value)
		}

		if text == "" {
			w.WriteHeader(http.StatusBadRequest)

			w.Write([]byte("Invalid status code requested.\n"))
		} else {
			w.WriteHeader(value)

			w.Write([]byte(text + "\n"))
		}

		if verbose {
			fmt.Printf("%s | %s requested status code %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				trimmed)
		}
	}
}

func registerHTTPStatus(module string, mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) []string {
	mux.GET("/http/", serveUsage(module, usage))

	mux.GET("/http/status/:status", serveHTTPStatusCode(errorChannel))
	mux.GET("/http/status/", serveUsage(module, usage))

	examples := make([]string, 3)
	examples[0] = "/http/status/200"
	examples[1] = "/http/status/404"
	examples[2] = "/http/status/500"

	return examples
}
