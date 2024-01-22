/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
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

func getHostname(host net.IP, resolver *net.Resolver) (string, error) {
	hosts, err := resolver.LookupAddr(context.Background(), host.String())
	if err != nil {
		return "", err
	}

	return hosts[0], nil
}

func getIP(host string, resolver *net.Resolver) (net.IP, error) {
	hosts, err := resolver.LookupHost(context.Background(), host)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(hosts[0]), nil
}

func parseHost(host, protocol string, ctx *ipisp.BulkClient, resolver *net.Resolver) (string, error) {
	ips, err := resolver.LookupIP(context.Background(), protocol, host)
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

		hostname, err := getHostname(ip, resolver)
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

func serveHostRecord(protocol string, resolver *net.Resolver, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx, err := getBulkClient()
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			_, err = w.Write([]byte("Lookup failed.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		host := strings.TrimPrefix(p.ByName("host"), "/")

		parsedHost, err := parseHost(host, protocol, ctx, resolver)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			_, err = w.Write([]byte("Lookup failed.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		_, err = w.Write([]byte(parsedHost + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested host records for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}

func parseMX(ctx *ipisp.BulkClient, resolver *net.Resolver, host string) (string, error) {
	records, err := resolver.LookupMX(context.Background(), host)
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
		ip, err := getIP(hosts[h], resolver)
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

func serveMXRecord(resolver *net.Resolver, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx, err := getBulkClient()
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			_, err = w.Write([]byte("Lookup failed.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		host := strings.TrimPrefix(p.ByName("host"), "/")

		parsedHost, err := parseMX(ctx, resolver, host)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			_, err = w.Write([]byte("Lookup failed.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		_, err = w.Write([]byte(parsedHost + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested MX records for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}

func parseNS(ctx *ipisp.BulkClient, resolver *net.Resolver, host string) (string, error) {
	records, err := resolver.LookupNS(context.Background(), host)
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
		ip, err := getIP(hosts[h], resolver)
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

	_, err = retVal.WriteString(fmt.Sprintf("%v:\n", host))
	if err != nil {
		return "", err
	}

	for response := 0; response < len(responses); response++ {
		retVal.WriteString(fmt.Sprintf("\n  %v:\n    IP: %v\n    Provider: %v (%v)\n",
			strings.TrimRight(hosts[response], "."),
			responses[response].IP,
			responses[response].ASN,
			responses[response].ISPName))
	}

	return retVal.String(), nil
}

func serveNSRecord(resolver *net.Resolver, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		ctx, err := getBulkClient()
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			_, err = w.Write([]byte("Lookup failed.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		host := strings.TrimPrefix(p.ByName("host"), "/")

		parsedHost, err := parseNS(ctx, resolver, host)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			_, err = w.Write([]byte("Lookup failed.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		_, err = w.Write([]byte(parsedHost + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested NS records for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				host)
		}
	}
}

func registerDNS(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "dns"

	var resolver *net.Resolver

	if dnsResolver != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(10000),
				}
				return d.DialContext(ctx, network, dnsResolver)
			},
		}
	} else {
		resolver = net.DefaultResolver
	}

	mux.GET("/dns/", serveUsage(module, usage))

	mux.GET("/dns/a/:host", serveHostRecord("ip4", resolver, errorChannel))
	mux.GET("/dns/a/", serveUsage(module, usage))

	mux.GET("/dns/aaaa/:host", serveHostRecord("ip6", resolver, errorChannel))
	mux.GET("/dns/aaaa/", serveUsage(module, usage))

	mux.GET("/dns/host/:host", serveHostRecord("ip", resolver, errorChannel))
	mux.GET("/dns/host/", serveUsage(module, usage))

	mux.GET("/dns/mx/:host", serveMXRecord(resolver, errorChannel))
	mux.GET("/dns/mx/", serveUsage(module, usage))

	mux.GET("/dns/ns/:host", serveNSRecord(resolver, errorChannel))
	mux.GET("/dns/ns/", serveUsage(module, usage))

	usage.Store(module, []string{
		"/dns/a/google.com",
		"/dns/aaaa/google.com",
		"/dns/host/google.com",
		"/dns/mx/google.com",
		"/dns/ns/google.com",
	})
}
