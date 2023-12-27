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

		retVal.WriteString(fmt.Sprintf("  %v:\n    Provider: %v (%v)\n    Hostname: %v\n    Range: %v\n\n",
			ip,
			r.ASN,
			r.ISPName,
			strings.TrimRight(hostname, "."),
			r.Range))
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
			fmt.Printf("%s | %s requested host records for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}

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

func serveMXRecord(errorChannel chan<- error) httprouter.Handle {
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
			fmt.Printf("%s | %s requested MX records for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}

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
		retVal.WriteString(fmt.Sprintf("\n  %v:\n    IP: %v\n    Provider: %v (%v)\n",
			strings.TrimRight(hosts[response], "."),
			responses[response].IP,
			responses[response].ASN,
			responses[response].ISPName))
	}

	return retVal.String(), nil
}

func serveNSRecord(errorChannel chan<- error) httprouter.Handle {
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
			fmt.Printf("%s | %s requested NS records for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}

func registerDNSHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/dns/a/*host", serveHostRecord("ip4", errorChannel))
	mux.GET("/dns/aaaa/*host", serveHostRecord("ip6", errorChannel))
	mux.GET("/dns/host/*host", serveHostRecord("ip", errorChannel))
	mux.GET("/dns/mx/*host", serveMXRecord(errorChannel))
	mux.GET("/dns/ns/*host", serveNSRecord(errorChannel))

	var usage []string
	usage = append(usage, "/dns/a/google.com")
	usage = append(usage, "/dns/aaaa/google.com")
	usage = append(usage, "/dns/host/google.com")
	usage = append(usage, "/dns/mx/google.com")
	usage = append(usage, "/dns/ns/google.com")

	return usage
}
