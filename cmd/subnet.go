/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	tpl4 = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<meta name="Description" content="Serves a variety of web-based utilities." />
    <meta charset="utf-8" />
	<title>Query v{{.Version}}</title>
    <link rel="stylesheet" href="/css/subnet.css" />
    <meta property="og:site_name" content="https://github.com/Seednode/trivia"/>
    <meta property="og:title" content="Query v{{.Version}}"/>
    <meta property="og:description" content="Serves a variety of web-based utilities."/>
    <meta property="og:url" content="https://github.com/Seednode/query"/>
    <meta property="og:type" content="website"/>
  </head>

  <body>
	<p id="table">
      <table>
        <tr>
		  <th></th>
		  <th>Binary</th>
		  <th>Decimal</th>
		</tr>
		<tr>
		  <th>Address</th>
		  <td>{{.Address_Binary}}</td>
		  <td>{{.Address_Decimal}}</td>
		</tr>
		<tr>
		  <th>Mask</th>
		  <td>{{.Mask_Binary}}</td>
		  <td>{{.Mask_Decimal}}</td>
		</tr>
		<tr>
		  <th>First</th>
		  <td>{{.First_Binary}}</td>
		  <td>{{.First_Decimal}}</td>
		</tr>
		<tr>
		  <th>Last</th>
		  <td>{{.Last_Binary}}</td>
		  <td>{{.Last_Decimal}}</td>
		</tr>
		<tr>
		  <th>Total</th>
		  <td colspan="2">{{.Total}}</td>
		</tr>
	  </table>
	</p>
  </body>
</html>
`

	tpl6 = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<meta name="Description" content="Serves a variety of web-based utilities." />
    <meta charset="utf-8" />
	<title>Query v{{.Version}}</title>
    <link rel="stylesheet" href="/css/subnet.css" />
    <meta property="og:site_name" content="https://github.com/Seednode/trivia"/>
    <meta property="og:title" content="Query v{{.Version}}"/>
    <meta property="og:description" content="Serves a variety of web-based utilities."/>
    <meta property="og:url" content="https://github.com/Seednode/query"/>
    <meta property="og:type" content="website"/>
  </head>

  <body>
	<p id="table">
      <table>
        <tr>
		  <th></th>
		  <th>Binary</th>
		  <th>Hex (Full)</th>
		  <th>Hex (Shortened)</th>
		</tr>
		<tr>
		  <th>Address</th>
		  <td>{{.Address_Binary}}</td>
		  <td>{{.Address_Hex}}</td>
		  <td>{{.Address_Short}}</td>
		</tr>
		<tr>
		  <th>Mask</th>
		  <td>{{.Mask_Binary}}</td>
		  <td>{{.Mask_Hex}}</td>
		  <td>{{.Mask_Short}}</td>
		</tr>
		<tr>
		  <th>First</th>
		  <td>{{.First_Binary}}</td>
		  <td>{{.First_Hex}}</td>
		  <td>{{.First_Short}}</td>
		</tr>
		<tr>
		  <th>Last</th>
		  <td>{{.Last_Binary}}</td>
		  <td>{{.Last_Hex}}</td>
		  <td>{{.Last_Short}}</td>
		</tr>
		<tr>
		  <th>Total</th>
		  <td colspan="3">{{.Total}}</td>
		</tr>
	  </table>
	</p>
  </body>
</html>
`
)

type Template4 struct {
	Version         string
	Address_Binary  string
	Address_Decimal string
	Mask_Binary     string
	Mask_Decimal    string
	First_Binary    string
	First_Decimal   string
	Last_Binary     string
	Last_Decimal    string
	Total           string
}

type Template6 struct {
	Version        string
	Address_Binary string
	Address_Hex    string
	Address_Short  string
	Mask_Binary    string
	Mask_Hex       string
	Mask_Short     string
	First_Binary   string
	First_Hex      string
	First_Short    string
	Last_Binary    string
	Last_Hex       string
	Last_Short     string
	Total          string
}

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

func toDottedDecimal(b []byte) string {
	var s strings.Builder

	for i := 0; i < len(b); i++ {
		s.WriteString(fmt.Sprintf("%d", b[i]))

		if i != len(b)-1 {
			s.WriteString(".")
		}
	}

	return s.String()
}

func calculateV4Subnet(cidr string) (Template4, error) {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return Template4{}, errors.New("not valid CIDR notation")
	}

	as4 := ip.To4()

	if as4 == nil {
		return Template4{}, errors.New("not a valid IPv6 address")
	}

	first, err := and(as4, net.Mask)
	if err != nil {
		return Template4{}, err
	}

	last, err := or(as4, invert(net.Mask))
	if err != nil {
		return Template4{}, err
	}

	return Template4{
		Version:         ReleaseVersion,
		Address_Binary:  toBinary(as4),
		Address_Decimal: toDottedDecimal(as4),
		Mask_Binary:     toBinary(net.Mask),
		Mask_Decimal:    toDottedDecimal(net.Mask),
		First_Binary:    toBinary(first),
		First_Decimal:   toDottedDecimal(first),
		Last_Binary:     toBinary(last),
		Last_Decimal:    toDottedDecimal(last),
		Total:           subtract(first, last),
	}, nil
}

func serveV4Subnet(template *template.Template, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/html;charset=UTF-8")

		data, err := calculateV4Subnet(strings.TrimPrefix(p.ByName("v4"), "/"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		err = template.Execute(w, data)
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

func calculateV6Subnet(cidr string) (Template6, error) {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return Template6{}, errors.New("not valid CIDR notation")
	}

	as4 := ip.To4()

	if as4 != nil {
		return Template6{}, errors.New("not a valid IPv6 address")
	}

	first, err := and(ip, net.Mask)
	if err != nil {
		return Template6{}, err
	}

	last, err := or(ip, invert(net.Mask))
	if err != nil {
		return Template6{}, err
	}

	return Template6{
		Version:        ReleaseVersion,
		Address_Binary: toBinary(ip),
		Address_Hex:    toColonedHex(ip),
		Address_Short:  ip.String(),
		Mask_Binary:    toBinary(net.Mask),
		Mask_Hex:       toColonedHex(net.Mask),
		Mask_Short:     net.Mask.String(),
		First_Binary:   toBinary(first),
		First_Hex:      toColonedHex(first),
		First_Short:    first.String(),
		Last_Binary:    toBinary(last),
		Last_Hex:       toColonedHex(last),
		Last_Short:     last.String(),
		Total:          subtract(first, last),
	}, nil
}

func serveV6Subnet(template *template.Template, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/html;charset=UTF-8")

		data, err := calculateV6Subnet(strings.TrimPrefix(p.ByName("v6"), "/"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		err = template.Execute(w, data)
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

	template4, err := template.New("subnet").Parse(tpl4)
	if err != nil {
		errorChannel <- Error{err, "", ""}

		return
	}

	template6, err := template.New("subnet").Parse(tpl6)
	if err != nil {
		errorChannel <- Error{err, "", ""}

		return
	}

	mux.GET("/subnet/", serveUsage(module, usage, errorChannel))
	mux.GET("/subnet/v4/*v4", serveV4Subnet(template4, errorChannel))
	mux.GET("/subnet/v6/*v6", serveV6Subnet(template6, errorChannel))

	usage.Store(module, []string{
		"/subnet/v4/192.168.0.1/24",
		"/subnet/v4/10.10.100.0/22",
		"/subnet/v6/fdd8:0c61:bf60:590f::/64",
		"/subnet/v6/2606:4700:a560::/48",
	})
}
