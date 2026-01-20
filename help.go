/*
Copyright © 2025 Seednode <seednode@seedno.de>
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

func serveUsage(module string, usage *sync.Map, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		securityHeaders(w)

		var output strings.Builder

		var help []string

		output.WriteString("Examples:\n")

		usage.Range(func(key, value any) bool {
			if key == module {
				help = append(help, value.([]string)...)
			}

			return true
		})

		slices.Sort(help)

		for _, line := range help {
			output.WriteString(fmt.Sprintf("- %s\n", line))
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}

		_, err := w.Write([]byte(output.String()))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}
		}
	}
}

func serveHelp(usage *sync.Map, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		securityHeaders(w)

		var output strings.Builder

		output.WriteString(fmt.Sprintf("query v%s\n\n", ReleaseVersion))

		output.WriteString("Examples:\n")

		var help []string

		usage.Range(func(key, value any) bool {
			help = append(help, value.([]string)...)

			return true
		})

		slices.Sort(help)

		for _, line := range help {
			output.WriteString(fmt.Sprintf("- %s\n", line))
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

func registerHelp(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	mux.GET("/", serveHelp(usage, errorChannel))
}
