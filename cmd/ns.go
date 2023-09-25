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

func parseNS(ctx *ipisp.BulkClient, host string) (string, error) {
	records, err := net.LookupNS(host)
	if len(records) == 0 || err != nil {
		return "", err
	}

	var hosts []string

	for h := 0; h < len(records); h++ {
		record := records[h]
		hosts = append(hosts, record.Host)
	}

	var ips []net.IP
	for h := 0; h < len(hosts); h++ {
		ip, err := getIP(hosts[h])
		if err != nil {
			return "", err
		}

		ips = append(ips, ip)
	}

	responses, err := ctx.LookupIPs(ips...)
	if len(responses) == 0 || err != nil {
		return "", err
	}

	var retVal strings.Builder

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

	return retVal.String(), nil
}

func getNSRecord(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx, err := getBulkClient()
		if err != nil {
			errorChannel <- err

			return
		}

		host := strings.TrimPrefix(p[0].Value, "/")

		parsedHost, err := parseNS(ctx, host)
		if err != nil {
			errorChannel <- err

			return
		}

		w.Write([]byte(parsedHost + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up NS records for %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}
