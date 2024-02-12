/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
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
	xRealIP := r.Header.Get("X-Real-Ip")

	switch {
	case cfIP != "" && includePort:
		return cfIP + ":" + port
	case cfIP != "":
		return cfIP
	case xRealIP != "" && includePort:
		return xRealIP + ":" + port
	case xRealIP != "":
		return xRealIP
	case includePort:
		return host + ":" + port
	default:
		return host
	}
}

func serveIP(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}

		_, err := w.Write([]byte(realIP(r, false) + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}
	}
}

func registerIP(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "ip"

	mux.GET("/ip/", serveIP(errorChannel))
	mux.GET("/ip/:ip", serveIP(errorChannel))

	usage.Store(module, []string{
		"/ip/",
	})
}
