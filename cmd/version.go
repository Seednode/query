/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func serveVersion() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := []byte(fmt.Sprintf("query v%s\n", ReleaseVersion))

		w.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(data))))

		w.Write(data)
	}
}

func registerVersionHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/version", serveVersion())
	mux.GET("/version/*version", serveVersion())

	var usage []string
	usage = append(usage, "/version/")

	return usage
}
