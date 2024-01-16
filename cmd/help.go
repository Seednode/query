/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func getUsage(usage map[string][]string) []string {
	var help []string

	for _, i := range usage {
		help = append(help, i...)
	}

	return help
}

func serveUsage(module string, usage map[string][]string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		var output strings.Builder

		output.WriteString("Examples:\n")

		for _, line := range usage[module] {
			output.WriteString(fmt.Sprintf("- %s\n", line))
		}

		w.Write([]byte(output.String()))

		if verbose {
			fmt.Printf("%s | %s requested usage info for %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				module)
		}
	}
}

func serveHelp(usage []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		var output strings.Builder

		output.WriteString(fmt.Sprintf("query v%s\n\n", ReleaseVersion))

		output.WriteString("Examples:\n")

		for _, line := range usage {
			output.WriteString(fmt.Sprintf("- %s\n", line))
		}

		w.Write([]byte(output.String()))

		if verbose {
			fmt.Printf("%s | %s requested usage info\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerHelp(mux *httprouter.Router, usage []string, errorChannel chan<- Error) {
	mux.GET("/", serveHelp(usage))
}
