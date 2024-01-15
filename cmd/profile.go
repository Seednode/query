/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

func registerProfileHandlers(mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) []string {
	mux.HandlerFunc("GET", "/pprof/", pprof.Index)
	mux.HandlerFunc("GET", "/pprof/cmdline", pprof.Cmdline)
	mux.HandlerFunc("GET", "/pprof/profile", pprof.Profile)
	mux.HandlerFunc("GET", "/pprof/symbol", pprof.Symbol)
	mux.HandlerFunc("GET", "/pprof/trace", pprof.Trace)

	var examples []string
	examples = append(examples, "/pprof/")
	examples = append(examples, "/pprof/cmdline")
	examples = append(examples, "/pprof/profile")
	examples = append(examples, "/pprof/symbol")
	examples = append(examples, "/pprof/trace")

	return examples
}
