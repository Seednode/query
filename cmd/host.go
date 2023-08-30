/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ammario/ipisp/v2"
	"github.com/julienschmidt/httprouter"
)

func ParseHost(ctx *ipisp.BulkClient, host, protocol string) string {
	ips, _ := net.DefaultResolver.LookupIP(context.Background(), protocol, host)

	if len(ips) == 0 {
		return "No records found for specified host.\n"
	}

	responses, err := ctx.LookupIPs(ips...)
	if err != nil {
		return "Lookup failed.\n"
	}

	var retVal strings.Builder

	retVal.WriteString(fmt.Sprintf("%s:\n\n", host))

	for response := 0; response < len(responses); response++ {
		r := responses[response]
		ip := r.IP
		hostname := strings.TrimRight(GetHostname(ip), ".")
		asn := r.ASN
		provider := r.ISPName
		subnet := r.Range

		retVal.WriteString(fmt.Sprintf("  %v:\n    Provider: %v (%v)\n    Hostname: %v\n    Range: %v\n\n",
			ip, asn, provider, hostname, subnet))
	}

	return retVal.String()
}

func getHostRecord(protocol string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx := GetBulkClient()

		host := strings.TrimPrefix(p[0].Value, "/")

		w.Write([]byte(ParseHost(ctx, host, protocol) + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up host records for %s\n",
				startTime.Format(LogDate),
				realIP(r, true),
				host)
		}
	}
}
