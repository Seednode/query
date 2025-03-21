/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveHTTPStatusCode(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		securityHeaders(w)

		var text string = ""

		trimmed := strings.TrimSuffix(strings.TrimPrefix(p.ByName("status"), "/"), "/")

		value, err := strconv.Atoi(trimmed)
		if err == nil {
			text = http.StatusText(value)
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}

		if text == "" {
			w.WriteHeader(http.StatusBadRequest)

			_, err = w.Write([]byte("Invalid status code requested\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}
		} else {
			w.WriteHeader(value)

			_, err = w.Write([]byte(text + "\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}
		}
	}
}

func registerHTTPStatus(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "http"

	mux.GET("/http/", serveUsage(module, usage, errorChannel))

	mux.GET("/http/status/:status", serveHTTPStatusCode(errorChannel))
	mux.GET("/http/status/", serveUsage(module, usage, errorChannel))

	usage.Store(module, []string{
		"/http/status/200",
		"/http/status/404",
		"/http/status/500",
	})
}
