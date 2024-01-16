/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	ErrInvalidQRSize = errors.New("qr code size must be between 256 and 2048 pixels")
)

func serveQRCode(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		value := strings.TrimPrefix(p.ByName("string"), "/")
		if value == "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			value = string(body)
		}

		qrCode, err := qrcode.New(value, qrcode.Medium)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.Write([]byte("Failed to encode string.\n"))

			return
		}

		if r.URL.Query().Has("string") {
			w.Write([]byte("\n" + qrCode.ToString(false) + "\n"))
		} else {
			png, err := qrCode.PNG(qrSize)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			w.Header().Set("Content-Type", "image/png")

			w.Write(png)
		}

		if verbose {
			fmt.Printf("%s | %s requested QR code of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func registerQR(module string, mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) []string {
	mux.GET("/qr/:string", serveQRCode(errorChannel))
	mux.GET("/qr/", serveUsage(module, usage))

	examples := make([]string, 2)
	examples[0] = "/qr/Test"
	examples[1] = "/qr/Test?string"

	return examples
}
