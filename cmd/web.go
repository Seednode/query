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
	LogDate            string = `2006-01-02T15:04:05.000-07:00`
	RedirectStatusCode int    = http.StatusSeeOther
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

	TimeFormats := map[string]string{
		"ANSIC":       `Mon Jan _2 15:04:05 2006`,
		"DateOnly":    `2006-01-02`,
		"DateTime":    `2006-01-02 15:04:05`,
		"Kitchen":     `3:04PM`,
		"Layout":      `01/02 03:04:05PM '06 -0700`,
		"RFC1123":     `Mon, 02 Jan 2006 15:04:05 MST`,
		"RFC1123Z":    `Mon, 02 Jan 2006 15:04:05 -0700`,
		"RFC3339":     `2006-01-02T15:04:05Z07:00`,
		"RFC3339Nano": `2006-01-02T15:04:05.999999999Z07:00`,
		"RFC822":      `02 Jan 06 15:04 MST`,
		"RFC822Z":     `02 Jan 06 15:04 -0700`,
		"RFC850":      `Monday, 02-Jan-06 15:04:05 MST`,
		"RubyDate":    `Mon Jan 02 15:04:05 -0700 2006`,
		"Stamp":       `Jan _2 15:04:05`,
		"StampMicro":  `Jan _2 15:04:05.000000`,
		"StampMilli":  `Jan _2 15:04:05.000`,
		"StampNano":   `Jan _2 15:04:05.000000000`,
		"TimeOnly":    `15:04:05`,
		"UnixDate":    `Mon Jan _2 15:04:05 MST 2006`,
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	mux := httprouter.New()

	mux.PanicHandler = serverErrorHandler()

	mux.GET("/", serveVersion())

	mux.GET("/ip/*ip", serveIp())

	mux.GET("/time/*time", serveTime(TimeFormats))

	mux.GET("/roll/*roll", rollDice())

	mux.GET("/dns/a/*host", getARecord())

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
