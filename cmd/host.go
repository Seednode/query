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

func parseHost(ctx *ipisp.BulkClient, host, protocol string) (string, error) {
	ips, err := net.DefaultResolver.LookupIP(context.Background(), protocol, host)
	if len(ips) == 0 || err != nil {
		return "", err
	}

	responses, err := ctx.LookupIPs(ips...)
	if err != nil {
		return "", err
	}

	var retVal strings.Builder

	retVal.WriteString(fmt.Sprintf("%s:\n\n", host))

	for response := 0; response < len(responses); response++ {
		r := responses[response]
		ip := r.IP

		hostname, err := getHostname(ip)
		if err != nil {
			return "", err
		}

		hostname = strings.TrimRight(hostname, ".")
		asn := r.ASN
		provider := r.ISPName
		subnet := r.Range

		retVal.WriteString(fmt.Sprintf("  %v:\n    Provider: %v (%v)\n    Hostname: %v\n    Range: %v\n\n",
			ip, asn, provider, hostname, subnet))
	}

	return retVal.String(), nil
}

func serveHostRecord(protocol string, errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx, err := getBulkClient()
		if err != nil {
			errorChannel <- err

			return
		}

		host := strings.TrimPrefix(p[0].Value, "/")

		parsedHost, err := parseHost(ctx, host, protocol)
		if err != nil {
			errorChannel <- err

			return
		}

		w.Write([]byte(parsedHost + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up host records for %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}
