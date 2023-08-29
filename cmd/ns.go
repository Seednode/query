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

func GetNS(host string) []string {
	records, err := net.LookupNS(host)
	if err != nil {
		fmt.Println("\nNo NS records exist for host " + host + ".")
	}

	var hosts []string

	for h := 0; h < len(records); h++ {
		record := records[h]
		hosts = append(hosts, record.Host)
	}

	return hosts
}

func ParseNS(ctx *ipisp.BulkClient, host string) string {
	hosts := GetNS(host)

	if len(hosts) == 0 {
		return ""
	}

	var ips []net.IP
	for h := 0; h < len(hosts); h++ {
		ips = append(ips, GetIP(hosts[h]))
	}

	var retVal strings.Builder

	responses, err := ctx.LookupIPs(ips...)
	if err != nil {
		retVal.WriteString("Lookup failed.\n")

		return retVal.String()
	}

	retVal.WriteString(fmt.Sprintf("%v:\n", host))

	for response := 0; response < len(responses); response++ {
		r := responses[response]
		host := strings.TrimRight(hosts[response], ".")
		ip := r.IP
		asn := r.ASN
		provider := r.ISPName
		retVal.WriteString(fmt.Sprintf("\n  %v:\n    IP: %v\n    Provider: %v (%v)\n",
			host, ip, asn, provider))
	}

	return retVal.String()
}

func getNSRecord() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx := GetBulkClient()

		host := strings.TrimPrefix(p[0].Value, "/")

		w.Write([]byte(ParseNS(ctx, host) + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up NS records for %s\n",
				startTime.Format(LogDate),
				realIP(r, true),
				host)
		}
	}
}
