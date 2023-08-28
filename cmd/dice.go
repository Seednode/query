/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func rollDice(verbose bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "text/plain")

		c, d, _ := strings.Cut(strings.TrimPrefix(p[0].Value, "/"), "d")
		if c == "" {
			c = "1"
		}

		count, err := strconv.Atoi(c)
		if err != nil {
			fmt.Println(err)
		}

		die, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println(err)
		}

		var total int

		for i := 0; i < count; i++ {
			v := rand.Intn(die) + 1

			if verbose {
				fmt.Printf("Rolled d%d, result %d\n", die, v)
			}

			total += v
		}

		result := strconv.Itoa(total)
		if err != nil {
			fmt.Println(err)
		}

		w.Write([]byte(result + "\n"))

		fmt.Printf("%s rolled the dice!\n", realIP(r))
	}
}
