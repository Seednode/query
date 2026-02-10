/*
Copyright © 2026 Seednode <seednode@seedno.de>
*/

package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	number                 = regexp.MustCompile(`\d+`)
	ErrInvalidMaxDiceCount = errors.New("max dice roll count must be a positive integer")
	ErrInvalidMaxDiceSides = errors.New("max dice side count must be a positive integer")
)

func rollDice(count, die int64) ([]int64, []int64, error) {
	var i int64

	rolls := make([]int64, count)
	results := make([]int64, count)

	for i = range count {
		v, err := rand.Int(rand.Reader, big.NewInt(die))
		if err != nil {
			return rolls, results, err
		}

		v.Add(v, big.NewInt(1))

		rolls[i] = die

		results[i] = v.Int64()
	}

	return rolls, results, nil
}

func serveDiceRoll(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		wantsVerbose := r.URL.Query().Has("verbose")

		langHeaders := strings.Split(r.Header.Get("Accept-Language"), ",")

		lang := language.Tag{}

		for _, value := range langHeaders {
			i, _ := language.Parse(value)

			if (i != language.Tag{}) {
				lang = i
				break
			}
		}

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		securityHeaders(w)

		pr := message.NewPrinter(lang)

		var total int64 = 0
		var length int = 0

		rolledDice := make([]int64, 0)
		rolledResults := make([]int64, 0)

		longestDie := 0

		trimmed := strings.TrimPrefix(p.ByName("roll"), "/")

		rolls := strings.Split(trimmed, ",")

		for roll := 0; roll < len(rolls); roll += 1 {
			c, d, _ := strings.Cut(rolls[roll], "d")
			if c == "" {
				c = "1"
			}

			c = strings.Join(number.FindAllString(c, -1), "")
			d = strings.Join(number.FindAllString(d, -1), "")

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

			thisDie := len(fmt.Sprintf("%d", die))
			if thisDie > longestDie {
				longestDie = thisDie
			}

			switch {
			case count > int64(maxDiceRolls):
				if verbose {
					fmt.Printf("%s | %s => %s (too many dice)\n",
						startTime.Format(timeFormats["RFC3339"]),
						realIP(r, true),
						r.URL.Path)
				}

				w.WriteHeader(http.StatusBadRequest)

				_, err = w.Write(fmt.Appendf(nil, "Dice roll count must be no greater than %d", maxDiceRolls))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			case count < 1:
				if verbose {
					fmt.Printf("%s | %s => %s (too few dice)\n",
						startTime.Format(timeFormats["RFC3339"]),
						realIP(r, true),
						r.URL.Path)
				}

				w.WriteHeader(http.StatusBadRequest)

				_, err = w.Write([]byte("Cannot roll zero dice"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			case die > int64(maxDiceSides):
				if verbose {
					fmt.Printf("%s | %s => %s (too many sides)\n",
						startTime.Format(timeFormats["RFC3339"]),
						realIP(r, true),
						r.URL.Path)
				}

				w.WriteHeader(http.StatusBadRequest)

				_, err = w.Write(fmt.Appendf(nil, "Dice side count must be no greater than %d", maxDiceSides))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			case die < 1:
				if verbose {
					fmt.Printf("%s | %s => %s (too few sides)\n",
						startTime.Format(timeFormats["RFC3339"]),
						realIP(r, true),
						r.URL.Path)
				}

				w.WriteHeader(http.StatusBadRequest)

				_, err = w.Write([]byte("Dice cannot have zero sides"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			}

			w.Header().Set("Cache-Control", "no-store")

			theseDice, theseResults, err := rollDice(count, die)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}

			rolledDice = append(rolledDice, theseDice...)

			rolledResults = append(rolledResults, theseResults...)
		}

		padCountTo := len(fmt.Sprintf("%d", len(rolledResults)))
		padDiceTo := longestDie + 1
		padValueTo := longestDie

		for i := 0; i < len(rolledDice); i++ {
			total += rolledResults[i]

			if wantsVerbose {
				written, err := w.Write(fmt.Appendf(nil, "%*d | %*s -> %*d\n", padCountTo, i+1, padDiceTo, fmt.Sprintf("d%d", rolledDice[i]), padValueTo, rolledResults[i]))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}

					return
				}

				if written > length {
					length = written
				}
			}
		}

		if wantsVerbose {
			_, err := w.Write(fmt.Appendf(nil, "%s\nTotal: ", strings.Repeat("-", length-1)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}
		}

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}

		result, _ := strconv.Atoi(strconv.FormatInt(total, 10))

		_, err := w.Write([]byte(pr.Sprintf("%*d\n", length-8, result)))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
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
		"/roll/4d6,5d8,d4?verbose",
	})
}
