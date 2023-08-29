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

func GetMX(host string) ([]string, []uint16) {
	var hosts []string
	var priorities []uint16

	records, _ := net.LookupMX(host)

	for h := 0; h < len(records); h++ {
		record := records[h]
		hosts = append(hosts, record.Host)
		priorities = append(priorities, record.Pref)
	}

	return hosts, priorities
}

func ParseMX(ctx *ipisp.BulkClient, host string) string {
	hosts, priorities := GetMX(host)

	if len(hosts) == 0 {
		return "No MX records found for specified host.\n"
	}

	var ips []net.IP
	for h := 0; h < len(hosts); h++ {
		ips = append(ips, GetIP(hosts[h]))
	}

	responses, err := ctx.LookupIPs(ips...)
	if err != nil {
		return "Lookup failed.\n"
	}

	var retVal strings.Builder

	retVal.WriteString(fmt.Sprintf("%v:\n", host))

	for response := 0; response < len(responses); response++ {
		r := responses[response]
		host := strings.TrimRight(hosts[response], ".")
		priority := priorities[response]
		ip := r.IP
		asn := r.ASN
		provider := r.ISPName
		retVal.WriteString(fmt.Sprintf("\n  (%v) %v:\n    IP: %v\n    Provider: %v (%v)\n",
			priority, host, ip, asn, provider))
	}

	return retVal.String()
}

func getMXRecord() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx := GetBulkClient()

		host := strings.TrimPrefix(p[0].Value, "/")

		w.Write([]byte(ParseMX(ctx, host) + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up MX records for %s\n",
				startTime.Format(LogDate),
				realIP(r, true),
				host)
		}
	}
}
