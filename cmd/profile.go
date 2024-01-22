/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http/pprof"
	"sync"

	"github.com/julienschmidt/httprouter"
)

func registerProfile(module string, mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) []string {
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

	examples := make([]string, 10)
	examples[0] = "/pprof/allocs"
	examples[1] = "/pprof/block"
	examples[2] = "/pprof/cmdline"
	examples[3] = "/pprof/goroutine"
	examples[4] = "/pprof/heap"
	examples[5] = "/pprof/mutex"
	examples[6] = "/pprof/profile"
	examples[7] = "/pprof/symbol"
	examples[8] = "/pprof/threadcreate"
	examples[9] = "/pprof/trace"

	return examples
}
