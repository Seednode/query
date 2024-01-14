/*
Copyright © 2024 Seednode <seednode@seedno.de>
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
	fields := strings.SplitAfter(r.RemoteAddr, ":")

	host := strings.TrimSuffix(strings.Join(fields[:len(fields)-1], ""), ":")
	port := fields[len(fields)-1]

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
			fmt.Printf("%s | %s requested their IP\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true))
		}
	}
}

func registerIPHandlers(mux *httprouter.Router, errorChannel chan<- Error) []string {
	mux.GET("/ip", serveIp())
	mux.GET("/ip/*ip", serveIp())

	var usage []string
	usage = append(usage, "/ip/")

	return usage
}
