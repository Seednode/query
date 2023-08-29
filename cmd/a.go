/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ammario/ipisp/v2"
	"github.com/julienschmidt/httprouter"
)

func GetA(host string) []string {
	hosts, _ := net.LookupHost(host)

	return hosts
}

func ParseA(ctx *ipisp.BulkClient, host string) string {
	hosts := GetA(host)

	if len(hosts) == 0 {
		return "No A records found for specified host.\n"
	}

	var ips []net.IP
	for h := 0; h < len(hosts); h++ {
		ipAddr := GetIP(hosts[h])
		ips = append(ips, ipAddr)
	}

	responses, err := ctx.LookupIPs(ips...)
	if err != nil {
		return "Lookup failed.\n"
	}

	var retVal strings.Builder

	for response := 0; response < len(responses); response++ {
		r := responses[response]
		host := host
		ip := r.IP
		hostname := strings.TrimRight(GetHostname(ip), ".")
		asn := r.ASN
		provider := r.ISPName
		subnet := r.Range

		retVal.WriteString(fmt.Sprintf("%v:\n\n  %v:\n    Provider: %v (%v)\n    Hostname: %v\n    Range: %v",
			host, ip, asn, provider, hostname, subnet))
	}

	return retVal.String()
}

func getARecord() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx := GetBulkClient()

		host := strings.TrimPrefix(p[0].Value, "/")

		w.Write([]byte(ParseA(ctx, host) + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up A records for %s\n",
				startTime.Format(LogDate),
				realIP(r, true),
				host)
		}
	}
}
