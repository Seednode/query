/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

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

const (
	redirectStatusCode int = http.StatusSeeOther
)

type Error struct {
	Message error
	Host    string
	Path    string
}

func serverError(w http.ResponseWriter, r *http.Request, i interface{}) {
	startTime := time.Now()

	if verbose {
		fmt.Printf("%s | Invalid request for %s from %s\n",
			startTime.Format(timeFormats["RFC3339"]),
			r.URL.Path,
			r.RemoteAddr,
		)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/plain")

	w.Write([]byte("500 Internal Server Error\n"))
}

func serverErrorHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	return serverError
}

func ServePage(args []string) error {
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

			fmt.Printf("%s | Error: `%v` (%s => %s)\n",
				time.Now().Format(timeFormats["RFC3339"]),
				err.Message,
				err.Host,
				err.Path)

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

	if draw || all {
		registerDraw(mux, &usage, errorChannel)
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
		registerProfile(mux, &usage, errorChannel)
	}

	if qr || all {
		registerQR(mux, &usage, errorChannel)
	}

	if roll || all {
		registerRoll(mux, &usage, errorChannel)
	}

	if timezones || all {
		registerTime(mux, &usage, errorChannel)
	}

	registerVersion(mux, &usage, errorChannel)

	registerHelp(mux, &usage, errorChannel)

	if verbose {
		fmt.Printf("%s | Listening on http://%s/\n",
			time.Now().Format(timeFormats["RFC3339"]),
			srv.Addr)
	}

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
