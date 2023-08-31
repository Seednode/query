/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bytes"
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
	logDate            string = `2006-01-02T15:04:05.000-07:00`
	redirectStatusCode int    = http.StatusSeeOther
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

	TimeFormats := timeFormats()

	mux := httprouter.New()

	mux.PanicHandler = serverErrorHandler()

	mux.GET("/", serveVersion())

	mux.GET("/ip/*ip", serveIp())

	mux.GET("/time/*time", serveTime(TimeFormats))

	mux.GET("/roll/*roll", rollDice())

	mux.GET("/dns/a/*host", getHostRecord("ip4"))

	mux.GET("/dns/aaaa/*host", getHostRecord("ip6"))

	mux.GET("/dns/host/*host", getHostRecord("ip"))

	mux.GET("/dns/mx/*host", getMXRecord())

	mux.GET("/dns/ns/*host", getNSRecord())

	srv := &http.Server{
		Addr:         net.JoinHostPort(bind, strconv.Itoa(int(port))),
		Handler:      mux,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Minute,
	}

	fmt.Printf("Server listening on %s...\n", srv.Addr)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
