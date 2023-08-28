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
		var total int

		roll := strings.TrimPrefix(p[0].Value, "/")

		s := strings.Split(roll, "d")

		count, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println(err)
		}

		die, err := strconv.Atoi(s[1])
		if err != nil {
			fmt.Println(err)
		}

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

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte(result))

		fmt.Printf("%s rolled the dice!\n", realIP(r))
	}
}
