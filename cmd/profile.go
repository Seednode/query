/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

func registerProfileHandlers(mux *httprouter.Router, errorChannel chan<- Error) []string {
	mux.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
	mux.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
	mux.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	mux.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)

	var usage []string
	usage = append(usage, "/debug/pprof/")
	usage = append(usage, "/debug/pprof/cmdline")
	usage = append(usage, "/debug/pprof/profile")
	usage = append(usage, "/debug/pprof/symbol")
	usage = append(usage, "/debug/pprof/trace")

	return usage
}
