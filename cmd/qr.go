/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	qrcode "github.com/skip2/go-qrcode"
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
			w.Write([]byte(qrCode.ToString(false)))
		} else {
			size := r.URL.Query().Get("size")

			sizeAsInt, err := strconv.Atoi(size)
			if err != nil || size == "" {
				sizeAsInt = 256
			}

			png, err := qrCode.PNG(sizeAsInt)
			if err != nil {
				errorChannel <- err

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			w.Header().Set("Content-Type", "image/png")

			w.Write(png)
		}

		if verbose {
			fmt.Printf("%s | %s encoded %q as a QR code\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}
