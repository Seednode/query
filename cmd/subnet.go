/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func toBinary(b []byte) string {
	var s strings.Builder

	for i := 0; i < len(b); i++ {
		s.WriteString(fmt.Sprintf("%08b", b[i]))

		if i < (len(b) - 1) {
			s.WriteString(" ")
		}
	}

	return s.String()
}

func subtract(a, b []byte) string {
	var c, d, e big.Int

	c.SetBytes(a)
	d.SetBytes(b)

	comp := c.Cmp(&d)
	switch comp {
	case -1:
		e.Sub(&d, &c)
	case 0:
		e = *big.NewInt(0)
	case 1:
		e.Sub(&c, &d)
	}

	return e.Add(&e, big.NewInt(1)).String()
}

func and(a, b []byte) (net.IP, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length %d does not equal length %d", len(a), len(b))
	}

	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] & b[i]
	}

	return result, nil
}

func or(a, b []byte) (net.IP, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length %d does not equal length %d", len(a), len(b))
	}

	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] | b[i]
	}

	return result, nil
}

func invert(b []byte) net.IP {
	inverted := make([]byte, len(b))

	for i := 0; i < len(b); i++ {
		inverted[i] = b[i] ^ ((2 << 7) - 1)
	}

	return inverted
}

func multiFormat(b []byte) string {
	return fmt.Sprintf("%s | %s", toBinary(b), toColonedHex(b))
}

func toColonedHex(b []byte) string {
	if len(b) != 16 {
		return ""
	}

	return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5], b[6], b[7],
		b[8], b[9], b[10], b[11],
		b[12], b[13], b[14], b[15])
}

func toDottedDecimalUint16(b []byte) string {
	var s strings.Builder

	for i := 0; i < len(b); i += 1 {
		if i%2 != 0 {
			continue
		}

		data := binary.BigEndian.Uint16(b[i : i+2])

		s.WriteString(fmt.Sprintf("%d", data))

		if i != len(b)-2 {
			s.WriteString(".")
		}
	}

	return s.String()
}

func toDottedDecimalUint8(b []byte) string {
	var s strings.Builder

	for i := 0; i < len(b); i++ {
		s.WriteString(fmt.Sprintf("%d", b[i]))

		if i != len(b)-1 {
			s.WriteString(".")
		}
	}

	return s.String()
}

func calculateV4Subnet(cidr string, r *http.Request, errorChannel chan<- Error) string {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return "Not valid CIDR notation.\n"
	}

	as4 := ip.To4()

	if as4 == nil {
		return "Not a valid IPv4 address.\n"
	}

	first, err := and(as4, net.Mask)
	if err != nil {
		errorChannel <- Error{err, realIP(r, true), r.URL.Path}

		return ""
	}

	last, err := or(as4, invert(net.Mask))
	if err != nil {
		errorChannel <- Error{err, realIP(r, true), r.URL.Path}

		return ""
	}

	return fmt.Sprintf("Address: %s | %s\nMask:    %s | %s\nFirst:   %s | %s\nLast:    %s | %s\nTotal:   %s\n",
		toBinary(as4), as4,
		toBinary(net.Mask), toDottedDecimalUint8(net.Mask),
		toBinary(first), first,
		toBinary(last), last,
		subtract(first, last))
}

func serveV4Subnet(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		resp := calculateV4Subnet(strings.TrimPrefix(p.ByName("v4"), "/"), r, errorChannel)

		_, err := w.Write([]byte(resp + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, false),
				r.RequestURI)
		}
	}
}

func calculateV6Subnet(cidr string, r *http.Request, errorChannel chan<- Error) string {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return "Not valid CIDR notation.\n"
	}

	as4 := ip.To4()

	if as4 != nil {
		return "Not a valid IPv6 address.\n"
	}

	first, err := and(ip, net.Mask)
	if err != nil {
		errorChannel <- Error{err, realIP(r, true), r.URL.Path}

		return ""
	}

	last, err := or(ip, invert(net.Mask))
	if err != nil {
		errorChannel <- Error{err, realIP(r, true), r.URL.Path}

		return ""
	}

	return fmt.Sprintf("Address: %s | %s | %s\nMask:    %s | %s | %s\nFirst:   %s | %s | %s\nLast:    %s | %s | %s\nTotal:   %s\n",
		multiFormat(ip), toDottedDecimalUint16(ip), ip.String(),
		multiFormat(net.Mask), toDottedDecimalUint16(net.Mask), net.Mask.String(),
		multiFormat(first), toDottedDecimalUint16(first), first.String(),
		multiFormat(last), toDottedDecimalUint16(last), last.String(),
		subtract(first, last))
}

func serveV6Subnet(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		resp := calculateV6Subnet(strings.TrimPrefix(p.ByName("v6"), "/"), r, errorChannel)

		_, err := w.Write([]byte(resp + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, false),
				r.RequestURI)
		}
	}
}

func registerSubnetting(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "subnet"

	mux.GET("/subnet/", serveUsage(module, usage, errorChannel))
	mux.GET("/subnet/v4/*v4", serveV4Subnet(errorChannel))
	mux.GET("/subnet/v6/*v6", serveV6Subnet(errorChannel))

	usage.Store(module, []string{
		"/subnet/v4/192.168.0.1/24",
		"/subnet/v4/10.10.100.0/22",
		"/subnet/v6/fdd8:0c61:bf60:590f::/64",
		"/subnet/v6/2606:4700:a560::/48",
	})
}
