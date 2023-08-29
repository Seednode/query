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

func rollDice() httprouter.Handle {
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
			w.Write([]byte("Invalid format: valid option regex '^[0-9]+?d[0-9]+$'\n"))

			return
		}

		die, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			w.Write([]byte("Invalid format: valid option regex '^[0-9]+?d[0-9]+$'\n"))

			return
		}

		var i, total int64

		for i = 0; i < count; i++ {
			v, err := rand.Int(rand.Reader, big.NewInt(die))
			if err != nil {
				w.Write([]byte("Invalid format: valid option regex '^[0-9]+?d[0-9]+$'\n"))

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
			fmt.Println(err)
		}

		w.Write([]byte(result + "\n"))

		if verbose {
			fmt.Printf("%s | %s rolled %dd%d, resulting in %d!\n",
				startTime.Format(LogDate),
				realIP(r, true),
				count,
				die,
				total)
		}
	}
}
