/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

func registerProfile(module string, mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) []string {
	mux.HandlerFunc("GET", "/pprof/", pprof.Index)

	mux.HandlerFunc("GET", "/pprof/cmdline", pprof.Cmdline)

	mux.HandlerFunc("GET", "/pprof/profile", pprof.Profile)

	mux.HandlerFunc("GET", "/pprof/symbol", pprof.Symbol)

	mux.HandlerFunc("GET", "/pprof/trace", pprof.Trace)

	examples := make([]string, 5)
	examples[0] = "/pprof/"
	examples[1] = "/pprof/cmdline"
	examples[2] = "/pprof/profile"
	examples[3] = "/pprof/symbol"
	examples[4] = "/pprof/trace"

	return examples
}
