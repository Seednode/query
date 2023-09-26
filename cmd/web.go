/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

const (
	redirectStatusCode int = http.StatusSeeOther
)

func serveVersion() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := []byte(fmt.Sprintf("query v%s\n", Version))

		w.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(data))))

		w.Write(data)
	}
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

	mux.GET("/", serveVersion())

	mux.GET("/ip/*ip", serveIp())

	mux.GET("/time/*time", serveTime(errorChannel))

	mux.GET("/roll/*roll", serveDiceRoll(errorChannel))

	mux.GET("/qr/*qr", serveQRCode(errorChannel))

	mux.GET("/dns/a/*host", getHostRecord("ip4", errorChannel))

	mux.GET("/dns/aaaa/*host", getHostRecord("ip6", errorChannel))

	mux.GET("/dns/host/*host", getHostRecord("ip", errorChannel))

	mux.GET("/dns/mx/*host", getMXRecord(errorChannel))

	mux.GET("/dns/ns/*host", getNSRecord(errorChannel))

	if profile {
		mux.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
		mux.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
		mux.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
		mux.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)
	}

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
