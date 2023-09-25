/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"net"

	"github.com/ammario/ipisp/v2"
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
