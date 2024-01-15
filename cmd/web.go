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
	"slices"
	"strconv"
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

	usage := make(map[string][]string)

	if !noDns {
		usage["dns"] = registerDNSHandlers("dns", mux, usage, errorChannel)
	}

	if !noDraw {
		usage["draw"] = registerDrawHandlers("draw", mux, usage, errorChannel)
	}

	if !noHash {
		usage["hash"] = registerHashHandlers("hash", mux, usage, errorChannel)
	}

	if !noHttpStatus {
		usage["http"] = registerHttpStatusHandlers("http", mux, usage, errorChannel)
	}

	if !noIp {
		usage["ip"] = registerIPHandlers("ip", mux, usage, errorChannel)
	}

	if !noMac {
		usage["mac"] = registerOUIHandlers("mac", mux, usage, errorChannel)
	}

	if profile {
		usage["profile"] = registerProfileHandlers("profile", mux, usage, errorChannel)
	}

	if !noQr {
		usage["qr"] = registerQRHandlers("qr", mux, usage, errorChannel)
	}

	if !noRoll {
		usage["roll"] = registerRollHandlers("roll", mux, usage, errorChannel)
	}

	if !noTime {
		usage["time"] = registerTimeHandlers("time", mux, usage, errorChannel)
	}

	usage["version"] = registerVersionHandlers("version", mux, usage, errorChannel)

	help := getUsage(usage)

	slices.Sort(help)

	registerHelpHandlers(mux, help, errorChannel)

	srv := &http.Server{
		Addr:         net.JoinHostPort(bind, strconv.Itoa(int(port))),
		Handler:      mux,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if verbose {
		fmt.Printf("%s | query v%s\n",
			time.Now().Format(timeFormats["RFC3339"]),
			ReleaseVersion)

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
