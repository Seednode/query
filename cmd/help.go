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

func serveHelp(usage []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		var output strings.Builder

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

func registerHelpHandlers(mux *httprouter.Router, usage []string, errorChannel chan<- error) {
	mux.GET("/", serveHelp(usage))
}
