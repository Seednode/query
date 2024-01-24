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

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		die, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		switch {
		case count > int64(maxDiceRolls):
			w.WriteHeader(http.StatusBadRequest)

			_, err = w.Write([]byte(fmt.Sprintf("Dice roll count must be no greater than %d", maxDiceRolls)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		case count < 1:
			w.WriteHeader(http.StatusBadRequest)

			_, err = w.Write([]byte("Cannot roll zero dice"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		case die > int64(maxDiceSides):
			w.WriteHeader(http.StatusBadRequest)

			_, err = w.Write([]byte(fmt.Sprintf("Dice side count must be no greater than %d", maxDiceSides)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		case die < 1:
			w.WriteHeader(http.StatusBadRequest)

			_, err = w.Write([]byte("Dice cannot have zero sides"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		w.Header().Set("Cache-Control", "no-store")

		var i, total int64

		padCountTo := len(fmt.Sprintf("%d", count))
		padValueTo := len(fmt.Sprintf("%d", die))

		length := 0

		for i = 0; i < count; i++ {
			v, err := rand.Int(rand.Reader, big.NewInt(die))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			v.Add(v, big.NewInt(1))

			if wantsVerbose {
				written, err := w.Write([]byte(fmt.Sprintf("%*d | d%d -> %*d\n", padCountTo, i+1, die, padValueTo, v)))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}

					w.WriteHeader(http.StatusInternalServerError)

					return
				}

				if written > length {
					length = written
				}
			}

			total += v.Int64()
		}

		result := strconv.FormatInt(total, 10)
		if err != nil {
			errorChannel <- Error{Message: err, Path: "serveDiceRoll()"}

			_, err = w.Write([]byte("An error occurred while calculating sum of all rolls\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		if wantsVerbose {
			_, err = w.Write([]byte(fmt.Sprintf("%s\nTotal: ", strings.Repeat("-", length-1))))
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
	const module = "roll"

	mux.GET("/roll/:roll", serveDiceRoll(errorChannel))
	mux.GET("/roll/", serveUsage(module, usage, errorChannel))

	usage.Store(module, []string{
		"/roll/5d20",
		"/roll/d6?verbose",
	})
}
