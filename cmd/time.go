/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveTime() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := strings.TrimPrefix(p[0].Value, "/")

		t := time.Now()

		tz, err := time.LoadLocation(path)

		if err != nil {
			t = t.In(tz)
		}

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(t.String()))
	}
}
