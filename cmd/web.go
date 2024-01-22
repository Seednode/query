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

	bindHost, err := net.LookupHost(bind)
	if err != nil {
		return err
	}

	bindAddr := net.ParseIP(bindHost[0])
	if bindAddr == nil {
		return errors.New("invalid bind address provided")
	}

	mux := httprouter.New()

	mux.PanicHandler = serverErrorHandler()

	errorChannel := make(chan Error)

	usage := sync.Map{}

	if !noDNS {
		usage.Store("dns", registerDNS("dns", mux, &usage, errorChannel))
	}

	if !noDraw {
		usage.Store("draw", registerDraw("draw", mux, &usage, errorChannel))
	}

	if !noHash {
		usage.Store("hash", registerHash("hash", mux, &usage, errorChannel))
	}

	if !noHTTPStatus {
		usage.Store("http", registerHTTPStatus("http", mux, &usage, errorChannel))
	}

	if !noIP {
		usage.Store("ip", registerIP("ip", mux, &usage, errorChannel))
	}

	if !noMAC {
		usage.Store("mac", registerMAC("mac", mux, &usage, errorChannel))
	}

	if profile {
		usage.Store("profile", registerProfile("profile", mux, &usage, errorChannel))
	}

	if !noQR {
		usage.Store("qr", registerQR("qr", mux, &usage, errorChannel))
	}

	if !noRoll {
		usage.Store("roll", registerRoll("roll", mux, &usage, errorChannel))
	}

	if !noTime {
		usage.Store("time", registerTime("time", mux, &usage, errorChannel))
	}

	usage.Store("version", registerVersion("version", mux, &usage, errorChannel))

	registerHelp(mux, &usage, errorChannel)

	srv := &http.Server{
		Addr:         net.JoinHostPort(bind, strconv.Itoa(int(port))),
		Handler:      mux,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Minute,
	}

	if verbose {
		fmt.Printf("%s | Listening on http://%s/\n",
			time.Now().Format(timeFormats["RFC3339"]),
			srv.Addr)
	}

	go func() {
		for err := range errorChannel {
			fmt.Printf("%s | Error: `%v` (%s => %s)\n",
				time.Now().Format(timeFormats["RFC3339"]),
				err.Message,
				err.Host,
				err.Path)

			if exitOnError {
				fmt.Printf("%s | Error: Shutting down...\n", time.Now().Format(timeFormats["RFC3339"]))

				srv.Shutdown(context.Background())
			}
		}
	}()

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
