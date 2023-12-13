/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http/pprof"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func registerProfileHandlers(mux *httprouter.Router, helpText *strings.Builder, errorChannel chan<- error) {
	mux.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
	mux.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
	mux.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	mux.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)
	helpText.WriteString("/debug/pprof/\n")
	helpText.WriteString("/debug/pprof/cmdline\n")
	helpText.WriteString("/debug/pprof/profile\n")
	helpText.WriteString("/debug/pprof/symbol\n")
	helpText.WriteString("/debug/pprof/trace\n")
}
