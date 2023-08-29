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
	hosts, err := net.LookupHost(host)
	if err != nil {
		fmt.Println("\nNo A records exist for host " + host + ".")
	}

	return hosts
}

func ParseA(ctx *ipisp.BulkClient, host string) string {
	hosts := GetA(host)

	if len(hosts) == 0 {
		return ""
	}

	var ips []net.IP
	for h := 0; h < len(hosts); h++ {
		ipAddr := GetIP(hosts[h])
		ips = append(ips, ipAddr)
	}

	var retVal strings.Builder

	responses, err := ctx.LookupIPs(ips...)
	if err != nil {
		retVal.WriteString("No information found for provided addresses.")
		for i := 0; i < len(ips); i++ {
			retVal.WriteString(fmt.Sprintf("IP: %s\n", ips[i]))
		}
		return retVal.String()
	}

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
