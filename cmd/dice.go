/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveDiceRoll(errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		wantsVerbose := r.URL.Query().Has("verbose")

		w.Header().Set("Content-Type", "text/plain")

		c, d, _ := strings.Cut(strings.TrimPrefix(p[0].Value, "/"), "d")
		if c == "" {
			c = "1"
		}

		count, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			errorChannel <- err

			return
		}

		die, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			errorChannel <- err

			return
		}

		var i, total int64

		for i = 0; i < count; i++ {
			v, err := rand.Int(rand.Reader, big.NewInt(die))
			if err != nil {
				errorChannel <- err

				return
			}

			v.Add(v, big.NewInt(1))

			if wantsVerbose {
				w.Write([]byte(fmt.Sprintf("Rolled d%d, result %d\n", die, v)))
			}

			total += v.Int64()
		}

		result := strconv.FormatInt(total, 10)
		if err != nil {
			w.Write([]byte("An error occurred while calculating sum of all rolls.\n"))

			return
		}

		if wantsVerbose {
			w.Write([]byte("\nTotal: "))
		}
		w.Write([]byte(result + "\n"))

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
