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
	"sync"
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

		value := ""

		switch r.Method {
		case http.MethodGet:
			value = strings.TrimPrefix(p.ByName("string"), "/")
		case http.MethodPost:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusInternalServerError)

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			value = string(body)
		}

		qrCode, err := qrcode.New(value, qrcode.Medium)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			w.Write([]byte("Failed to encode string.\n"))

			return
		}

		if r.URL.Query().Has("string") {
			_, err = w.Write([]byte("\n" + qrCode.ToString(false) + "\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusInternalServerError)

				w.Write([]byte("Failed to encode string.\n"))

				return
			}
		} else {
			png, err := qrCode.PNG(qrSize)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusInternalServerError)

				w.Write([]byte("Failed to encode string.\n"))

				return
			}

			w.Header().Set("Content-Type", "image/png")

			_, err = w.Write(png)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusInternalServerError)

				w.Write([]byte("Failed to encode string.\n"))

				return
			}
		}

		if verbose {
			fmt.Printf("%s | %s requested QR code of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func registerQR(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "qr"

	mux.GET("/qr/", serveUsage(module, usage))
	mux.GET("/qr/:string", serveQRCode(errorChannel))
	mux.POST("/qr/", serveQRCode(errorChannel))

	usage.Store(module, []string{
		"/qr/Test",
		"/qr/Test?string",
	})
}
