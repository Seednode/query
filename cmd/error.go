/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"time"
)

type Error struct {
	Message error
	Host    string
	Path    string
}

func serverError(w http.ResponseWriter, r *http.Request, i interface{}) {
	startTime := time.Now()

	if verbose {
		fmt.Printf("%s | Invalid request for %s from %s\n",
			startTime.Format(timeFormats["RFC3339"]),
			r.URL.Path,
			r.RemoteAddr,
		)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/plain")

	w.Write([]byte("500 Internal Server Error\n"))
}

func serverErrorHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	return serverError
}
