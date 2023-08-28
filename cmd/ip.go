/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func serveIp() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(realIP(r)))
	}
}
