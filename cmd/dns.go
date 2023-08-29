/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/ammario/ipisp/v2"
)

func GetBulkClient() *ipisp.BulkClient {
	c, err := ipisp.DialBulkClient(context.Background())
	if err != nil {
		panic(err)
	}
	return c
}

func GetHostname(host net.IP) string {
	hosts, err := net.LookupAddr(host.String())
	if err != nil {
		return "n/a"
	}
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	hostname := hosts[0]
	return hostname
}

func GetIP(host string) net.IP {
	hosts, err := net.LookupHost(host)
	if err != nil {
		panic(err)
	}
	ip := net.ParseIP(hosts[0])
	return ip
}
