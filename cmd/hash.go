/*
Copyright © 2024 Seednode <seednode@seedno.de>
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
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidHashAlgorithm = errors.New("invalid hash algorithm provided")
)

func serveHash(algorithm string, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		value := ""

		switch r.Method {
		case http.MethodGet:
			value = strings.TrimPrefix(p.ByName("string"), "/")
		case http.MethodPost:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusInternalServerError)

				_, err := w.Write([]byte("Failed to hash string\n"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

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
			errorChannel <- Error{ErrInvalidHashAlgorithm, realIP(r, true), r.URL.Path}

			w.WriteHeader(http.StatusBadRequest)

			_, err := w.Write([]byte("Invalid hash algorithm requested\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		_, err := io.WriteString(h, value)
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested %s hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				algorithm,
				value)
		}
	}
}

func registerHash(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "hash"

	mux.GET("/hash/", serveUsage(module, usage, errorChannel))

	mux.GET("/hash/md5/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/md5/:string", serveHash("MD5", errorChannel))
	mux.POST("/hash/md5/", serveHash("MD5", errorChannel))

	mux.GET("/hash/sha1/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha1/:string", serveHash("SHA-1", errorChannel))
	mux.POST("/hash/sha1/", serveHash("SHA-1", errorChannel))

	mux.GET("/hash/sha224/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha224/:string", serveHash("SHA-224", errorChannel))
	mux.POST("/hash/sha224/", serveHash("SHA-224", errorChannel))

	mux.GET("/hash/sha256/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha256/:string", serveHash("SHA-256", errorChannel))
	mux.POST("/hash/sha256/", serveHash("SHA-256", errorChannel))

	mux.GET("/hash/sha384/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha384/:string", serveHash("SHA-384", errorChannel))
	mux.POST("/hash/sha384/", serveHash("SHA-384", errorChannel))

	mux.GET("/hash/sha512/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha512/:string", serveHash("SHA-512", errorChannel))
	mux.POST("/hash/sha512/", serveHash("SHA-512", errorChannel))

	mux.GET("/hash/sha512-224/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha512-224/:string", serveHash("SHA-512/224", errorChannel))
	mux.POST("/hash/sha512-224/", serveHash("SHA-512/224", errorChannel))

	mux.GET("/hash/sha512-256/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha512-256/:string", serveHash("SHA-512/256", errorChannel))
	mux.POST("/hash/sha512-256/", serveHash("SHA-512/256", errorChannel))

	usage.Store(module, []string{
		"/hash/md5/foo",
		"/hash/sha1/foo",
		"/hash/sha224/foo",
		"/hash/sha256/foo",
		"/hash/sha384/foo",
		"/hash/sha512/foo",
		"/hash/sha512-224/foo",
		"/hash/sha512-256/foo",
	})
}
