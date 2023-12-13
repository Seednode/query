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

func serveHelp(helpText *strings.Builder) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(helpText.String() + "\n"))

		if verbose {
			fmt.Printf("%s | %s requested usage info\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerHelpHandlers(mux *httprouter.Router, helpText *strings.Builder, errorChannel chan<- error) {
	mux.GET("/", serveHelp(helpText))
	mux.GET("/help", serveHelp(helpText))
	mux.GET("/help/*help", serveHelp(helpText))
	helpText.WriteString("/help/\n")
}
