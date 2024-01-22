/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveVersion(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := []byte(fmt.Sprintf("query v%s\n", ReleaseVersion))

		w.Header().Set("Content-Length", strconv.Itoa(len(data)))

		_, err := w.Write(data)
		if err != nil {
			errorChannel <- Error{Message: err, Path: "serveVersion()"}
		}

		if verbose {
			fmt.Printf("%s | %s requested version info\n",
				time.Now().Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerVersion(module string, mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) []string {
	mux.GET("/version/", serveVersion(errorChannel))

	examples := make([]string, 1)
	examples[0] = "/version/"

	return examples
}
