/*
Copyright © 2023 Seednode <seednode@seedno.de>
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

func serveQRCode(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		value := strings.TrimPrefix(p[0].Value, "/")
		if value == "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				errorChannel <- err

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			value = string(body)
		}

		qrCode, err := qrcode.New(value, qrcode.Medium)
		if err != nil {
			errorChannel <- err

			w.Write([]byte("Failed to encode string.\n"))

			return
		}

		if r.URL.Query().Has("string") {
			w.Write([]byte("\n" + qrCode.ToString(false) + "\n"))
		} else {
			png, err := qrCode.PNG(qrSize)
			if err != nil {
				errorChannel <- err

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			w.Header().Set("Content-Type", "image/png")

			w.Write(png)
		}

		if verbose {
			fmt.Printf("%s | %s requested %q as a QR code\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func registerQRHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/qr/*qr", serveQRCode(errorChannel))

	var usage []string
	usage = append(usage, "/qr/Test")
	usage = append(usage, "/qr/Test?string")

	return usage
}
