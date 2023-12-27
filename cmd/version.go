/*
Copyright © 2023 Seednode <seednode@seedno.de>
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
			fmt.Printf("%s | %s requested version info for query\n",
				time.Now().Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerVersionHandlers(mux *httprouter.Router, errorChannel chan<- Error) []string {
	mux.GET("/version", serveVersion())
	mux.GET("/version/*version", serveVersion())

	var usage []string
	usage = append(usage, "/version/")

	return usage
}
