/*
Copyright © 2026 Seednode <seednode@seedno.de>
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Error struct {
	Message error
	Host    string
	Path    string
}

func securityHeaders(w http.ResponseWriter) {
	w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
	w.Header().Set("Permissions-Policy", "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(), payment=()")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-Xss-Protection", "1; mode=block")
}

func serverError(w http.ResponseWriter, r *http.Request, i any) {
	if verbose {
		fmt.Printf("%s | %s => %s (Invalid request)\n",
			time.Now().Format(timeFormats["RFC3339"]),
			realIP(r, true),
			r.RequestURI)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/plain")

	securityHeaders(w)

	w.Write([]byte("500 Internal Server Error\n"))
}

func serverErrorHandler() func(http.ResponseWriter, *http.Request, any) {
	return serverError
}

func servePage() error {
	timeZone := os.Getenv("TZ")
	if timeZone != "" {
		var err error

		time.Local, err = time.LoadLocation(timeZone)
		if err != nil {
			return err
		}
	}

	if verbose {
		fmt.Printf("%s | query v%s\n",
			time.Now().Format(timeFormats["RFC3339"]),
			ReleaseVersion)
	}

	bindAddr := net.ParseIP(bind)
	if bindAddr == nil {
		return errors.New("invalid bind address provided")
	}

	mux := httprouter.New()

	mux.PanicHandler = serverErrorHandler()

	srv := &http.Server{
		Addr:         net.JoinHostPort(bind, strconv.Itoa(int(port))),
		Handler:      mux,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Minute,
	}

	errorChannel := make(chan Error)

	go func() {
		for err := range errorChannel {
			if err.Host == "" {
				err.Host = "local"
			}

			fmt.Printf("%s | %s => %s (Error: `%v`)\n",
				time.Now().Format(timeFormats["RFC3339"]),
				err.Host,
				err.Path,
				err.Message)

			if exitOnError {
				fmt.Printf("%s | Error: Shutting down...\n", time.Now().Format(timeFormats["RFC3339"]))

				srv.Shutdown(context.Background())

				break
			}
		}
	}()

	usage := sync.Map{}

	if dns || all {
		registerDNS(mux, &usage, errorChannel)
	}

	if hashing || all {
		registerHash(mux, &usage, errorChannel)
	}

	if httpStatus || all {
		registerHTTPStatus(mux, &usage, errorChannel)
	}

	if ip || all {
		registerIP(mux, &usage, errorChannel)
	}

	if mac || all {
		registerMAC(mux, &usage, errorChannel)
	}

	if profile {
		registerProfile(mux, &usage)
	}

	if qr || all {
		registerQR(mux, &usage, errorChannel)
	}

	if roll || all {
		registerRoll(mux, &usage, errorChannel)
	}

	if subnet || all {
		registerSubnetting(mux, &usage, errorChannel)
	}

	if timezones || all {
		registerTime(mux, &usage, errorChannel)
	}

	if whoami || all {
		registerWhoAmI(mux, &usage, errorChannel)
	}

	registerVersion(mux, &usage, errorChannel)

	registerHelp(mux, &usage, errorChannel)

	registerCss(mux, errorChannel)

	var err error

	if verbose {
		if tlsKey != "" && tlsCert != "" {
			fmt.Printf("%s | Listening on https://%s/\n",
				time.Now().Format(timeFormats["RFC3339"]),
				srv.Addr)

			err = srv.ListenAndServeTLS(tlsCert, tlsKey)
		} else {
			fmt.Printf("%s | Listening on http://%s/\n",
				time.Now().Format(timeFormats["RFC3339"]),
				srv.Addr)

			err = srv.ListenAndServe()
		}
	}

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
