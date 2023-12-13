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

func parseMX(ctx *ipisp.BulkClient, host string) (string, error) {
	records, err := net.LookupMX(host)
	if len(records) == 0 || err != nil {
		return "", err
	}

	var hosts []string
	var priorities []uint16

	for h := 0; h < len(records); h++ {
		record := records[h]
		hosts = append(hosts, record.Host)
		priorities = append(priorities, record.Pref)
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
		retVal.WriteString(fmt.Sprintf("\n  (%v) %v:\n    IP: %v\n    Provider: %v (%v)\n",
			priorities[response],
			strings.TrimRight(hosts[response], "."),
			r.IP,
			r.ASN,
			r.ISPName))
	}

	return retVal.String(), nil
}

func getMXRecord(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx, err := getBulkClient()
		if err != nil {
			errorChannel <- err

			return
		}

		host := strings.TrimPrefix(p[0].Value, "/")

		parsedHost, err := parseMX(ctx, host)
		if err != nil {
			errorChannel <- err

			return
		}

		w.Write([]byte(parsedHost + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up MX records for %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}
