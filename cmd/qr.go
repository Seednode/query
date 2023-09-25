/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	qrcode "github.com/skip2/go-qrcode"
)

func serveQRCode() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "image/png")

		value := strings.TrimPrefix(p[0].Value, "/")

		png, err := qrcode.Encode(value, qrcode.Medium, 256)
		if err != nil {
			w.Write([]byte("Failed to encode string.\n"))

			return
		}

		w.Write([]byte(png))

		if verbose {
			fmt.Printf("%s | %s encoded %q as a QR code\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}
