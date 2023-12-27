/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidHashAlgorithm = errors.New("invalid hash algorithm provided")
)

func serveHash(algorithm string, errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		value := strings.TrimPrefix(p[0].Value, "/")
		if value == "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				errorChannel <- err

				w.Write([]byte("Failed to hash string.\n"))

				return
			}

			value = string(body)
		}

		var h hash.Hash

		switch algorithm {
		case "MD5":
			h = md5.New()
		case "SHA-1":
			h = sha1.New()
		case "SHA-224":
			h = sha256.New224()
		case "SHA-256":
			h = sha256.New()
		case "SHA-384":
			h = sha512.New384()
		case "SHA-512":
			h = sha512.New()
		case "SHA-512/224":
			h = sha512.New512_224()
		case "SHA-512/256":
			h = sha512.New512_256()
		default:
			errorChannel <- ErrInvalidHashAlgorithm

			return
		}

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested %s hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				algorithm,
				value)
		}
	}
}

func registerHashHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/hash/md5/*string", serveHash("MD5", errorChannel))
	mux.GET("/hash/sha1/*string", serveHash("SHA-1", errorChannel))
	mux.GET("/hash/sha224/*string", serveHash("SHA-224", errorChannel))
	mux.GET("/hash/sha256/*string", serveHash("SHA-256", errorChannel))
	mux.GET("/hash/sha384/*string", serveHash("SHA-384", errorChannel))
	mux.GET("/hash/sha512/*string", serveHash("SHA-512", errorChannel))
	mux.GET("/hash/sha512-224/*string", serveHash("SHA-512/224", errorChannel))
	mux.GET("/hash/sha512-256/*string", serveHash("SHA-512/256", errorChannel))

	var usage []string
	usage = append(usage, "/hash/md5/<string>")
	usage = append(usage, "/hash/sha1/<string>")
	usage = append(usage, "/hash/sha224/<string>")
	usage = append(usage, "/hash/sha256/<string>")
	usage = append(usage, "/hash/sha384/<string>")
	usage = append(usage, "/hash/sha512/<string>")
	usage = append(usage, "/hash/sha512-224/<string>")
	usage = append(usage, "/hash/sha512-256/<string>")

	return usage
}
