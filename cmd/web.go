/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
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

	rand.New(rand.NewSource(time.Now().UnixNano()))

	mux := httprouter.New()

	mux.PanicHandler = serverErrorHandler()

	errorChannel := make(chan error)

	var usage []string

	if profile {
		usage = append(usage, registerProfileHandlers(mux, errorChannel)...)
	}

	if !noDice {
		usage = append(usage, registerRollHandlers(mux, errorChannel)...)
	}

	if !noDNS {
		usage = append(usage, registerDNSHandlers(mux, errorChannel)...)
	}

	if !noHttpStatus {
		usage = append(usage, registerHttpStatusHandlers(mux, errorChannel)...)
	}

	if !noIP {
		usage = append(usage, registerIPHandlers(mux, errorChannel)...)
	}

	if !noOUI {
		usage = append(usage, registerOUIHandlers(mux, errorChannel)...)
	}

	if !noQR {
		usage = append(usage, registerQRHandlers(mux, errorChannel)...)
	}

	if !noTime {
		usage = append(usage, registerTimeHandlers(mux, errorChannel)...)
	}

	usage = append(usage, registerVersionHandlers(mux, errorChannel)...)

	registerHelpHandlers(mux, usage, errorChannel)

	srv := &http.Server{
		Addr:         net.JoinHostPort(bind, strconv.Itoa(int(port))),
		Handler:      mux,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Minute,
	}

	fmt.Printf("Server listening on %s...\n", srv.Addr)

	go func() {
		for err := range errorChannel {
			fmt.Printf("%s | Error: %v\n", time.Now().Format(timeFormats["RFC3339"]), err)

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
