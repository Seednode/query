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

func serveTime() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		t := time.Now()

		tz, err := time.LoadLocation(strings.TrimPrefix(p[0].Value, "/"))

		if err != nil {
			http.Redirect(w, r, "/time/", RedirectStatusCode)
		} else {
			t = t.In(tz)
		}

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(t.String()))

		fmt.Printf("%s checked the time!\n", realIP(r))
	}
}
