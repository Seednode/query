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

func realIP(r *http.Request, includePort bool) string {
	host, port, _ := strings.Cut(r.RemoteAddr, ":")

	if host == "" {
		return r.RemoteAddr
	}

	cfIP := r.Header.Get("Cf-Connecting-Ip")
	xRealIp := r.Header.Get("X-Real-Ip")

	switch {
	case cfIP != "" && includePort:
		return cfIP + ":" + port
	case cfIP != "":
		return cfIP
	case xRealIp != "" && includePort:
		return xRealIp + ":" + port
	case xRealIp != "":
		return xRealIp
	case includePort:
		return host + ":" + port
	default:
		return host
	}
}

func serveIp() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(realIP(r, false) + "\n"))

		if verbose {
			fmt.Printf("%s | %s checked their IP\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}
