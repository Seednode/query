/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"net"
	"strings"

	"github.com/ammario/ipisp/v2"
	"github.com/julienschmidt/httprouter"
)

func getBulkClient() (*ipisp.BulkClient, error) {
	c, err := ipisp.DialBulkClient(context.Background())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getHostname(host net.IP) (string, error) {
	hosts, err := net.LookupAddr(host.String())
	if err != nil {
		return "", err
	}

	return hosts[0], nil
}

func getIP(host string) (net.IP, error) {
	hosts, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(hosts[0]), nil
}

func registerDNSHandlers(mux *httprouter.Router, helpText *strings.Builder, errorChannel chan<- error) {
	mux.GET("/dns/a/*host", serveHostRecord("ip4", errorChannel))
	helpText.WriteString("/dns/a/google.com\n")

	mux.GET("/dns/aaaa/*host", serveHostRecord("ip6", errorChannel))
	helpText.WriteString("/dns/aaaa/google.com\n")

	mux.GET("/dns/host/*host", serveHostRecord("ip", errorChannel))
	helpText.WriteString("/dns/host/google.com\n")

	mux.GET("/dns/mx/*host", serveMXRecord(errorChannel))
	helpText.WriteString("/dns/mx/google.com\n")

	mux.GET("/dns/ns/*host", serveNSRecord(errorChannel))
	helpText.WriteString("/dns/ns/google.com\n")
}
