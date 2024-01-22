/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidMaxDiceCount = errors.New("max dice roll count must be a positive integer")
	ErrInvalidMaxDiceSides = errors.New("max dice side count must be a positive integer")
)

func serveDiceRoll(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		wantsVerbose := r.URL.Query().Has("verbose")

		w.Header().Set("Content-Type", "text/plain")

		c, d, _ := strings.Cut(strings.TrimPrefix(p.ByName("roll"), "/"), "d")
		if c == "" {
			c = "1"
		}

		count, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		die, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		switch {
		case count > int64(maxDiceRolls):
			_, err = w.Write([]byte(fmt.Sprintf("Dice roll count must be no greater than %d", maxDiceRolls)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}

			return
		case die > int64(maxDiceSides):
			_, err = w.Write([]byte(fmt.Sprintf("Dice side count must be no greater than %d", maxDiceSides)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}

			return
		}

		var i, total int64

		for i = 0; i < count; i++ {
			v, err := rand.Int(rand.Reader, big.NewInt(die))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}

			v.Add(v, big.NewInt(1))

			if wantsVerbose {
				_, err = w.Write([]byte(fmt.Sprintf("Rolled d%d, result %d\n", die, v)))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}

					return
				}
			}

			total += v.Int64()
		}

		result := strconv.FormatInt(total, 10)
		if err != nil {
			errorChannel <- Error{Message: err, Path: "serveDiceRoll()"}

			_, err = w.Write([]byte("An error occurred while calculating sum of all rolls.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}

			return
		}

		if wantsVerbose {
			_, err = w.Write([]byte("\nTotal: "))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}
		}

		_, err = w.Write([]byte(result + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s rolled %dd%d, resulting in %d\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				count,
				die,
				total)
		}
	}
}

func registerRoll(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	module := "roll"

	mux.GET("/roll/:roll", serveDiceRoll(errorChannel))
	mux.GET("/roll/", serveUsage(module, usage))

	usage.Store(module, []string{
		"/roll/5d20",
		"/roll/d6?verbose",
	})
}
