/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveVersion() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := []byte(fmt.Sprintf("query v%s\n", ReleaseVersion))

		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		w.Write(data)

		if verbose {
			fmt.Printf("%s | %s requested version info\n",
				time.Now().Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerVersion(module string, mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) []string {
	mux.GET("/version/", serveVersion())
	mux.GET("/version/:version", serveVersion())

	examples := make([]string, 1)
	examples[0] = "/version/"

	return examples
}
