/*
Copyright © 2026 Seednode <seednode@seedno.de>
*/

package main

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveWhoAmI(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		securityHeaders(w)

		var output strings.Builder

		request := make([]string, 0)

		for header := range r.Header {
			for _, value := range r.Header[header] {
				request = append(request, fmt.Sprintf("%s: %s\n", header, value))
			}
		}

		slices.SortStableFunc(request, func(a, b string) int {
			return strings.Compare(a, b)
		})

		for _, v := range request {
			output.WriteString(v)
		}

		_, err := w.Write([]byte(output.String()))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}
	}
}

func registerWhoAmI(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "whoami"

	mux.GET("/whoami", serveWhoAmI(errorChannel))

	usage.Store(module, []string{
		"/whoami",
	})
}
