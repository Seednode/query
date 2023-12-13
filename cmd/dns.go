/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"net"

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

func registerDNSHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/dns/a/*host", serveHostRecord("ip4", errorChannel))
	mux.GET("/dns/aaaa/*host", serveHostRecord("ip6", errorChannel))
	mux.GET("/dns/host/*host", serveHostRecord("ip", errorChannel))
	mux.GET("/dns/mx/*host", serveMXRecord(errorChannel))
	mux.GET("/dns/ns/*host", serveNSRecord(errorChannel))

	var usage []string
	usage = append(usage, "/dns/a/google.com")
	usage = append(usage, "/dns/aaaa/google.com")
	usage = append(usage, "/dns/host/google.com")
	usage = append(usage, "/dns/mx/google.com")
	usage = append(usage, "/dns/ns/google.com")

	return usage
}
