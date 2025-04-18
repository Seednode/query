/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"net/http/pprof"
	"sync"

	"github.com/julienschmidt/httprouter"
)

func registerProfile(mux *httprouter.Router, usage *sync.Map) {
	const module = "profile"

	mux.Handler("GET", "/pprof/allocs", pprof.Handler("allocs"))
	mux.Handler("GET", "/pprof/block", pprof.Handler("block"))
	mux.Handler("GET", "/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handler("GET", "/pprof/heap", pprof.Handler("heap"))
	mux.Handler("GET", "/pprof/mutex", pprof.Handler("mutex"))
	mux.Handler("GET", "/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.HandlerFunc("GET", "/pprof/cmdline", pprof.Cmdline)
	mux.HandlerFunc("GET", "/pprof/profile", pprof.Profile)
	mux.HandlerFunc("GET", "/pprof/symbol", pprof.Symbol)
	mux.HandlerFunc("GET", "/pprof/trace", pprof.Trace)

	usage.Store(module, []string{
		"/pprof/allocs",
		"/pprof/block",
		"/pprof/cmdline",
		"/pprof/goroutine",
		"/pprof/heap",
		"/pprof/mutex",
		"/pprof/profile",
		"/pprof/symbol",
		"/pprof/threadcreate",
		"/pprof/trace",
	})
}
