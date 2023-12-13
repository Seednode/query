/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func serveVersion() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := []byte(fmt.Sprintf("query v%s\n", ReleaseVersion))

		w.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(data))))

		w.Write(data)
	}
}

func registerVersionHandlers(mux *httprouter.Router, helpText *strings.Builder, errorChannel chan<- error) {
	mux.GET("/version", serveVersion())
	mux.GET("/version/*version", serveVersion())
	helpText.WriteString("/version/\n")
}
