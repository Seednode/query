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

type HashAlgorithm string

const (
	MD5        HashAlgorithm = "MD5"
	SHA1       HashAlgorithm = "SHA-1"
	SHA224     HashAlgorithm = "SHA-224"
	SHA256     HashAlgorithm = "SHA-256"
	SHA384     HashAlgorithm = "SHA-384"
	SHA512     HashAlgorithm = "SHA-512"
	SHA512_224 HashAlgorithm = "SHA-512/224"
	SHA512_256 HashAlgorithm = "SHA-512/256"
)

func serveHash(algorithm HashAlgorithm, errorChannel chan<- Error) httprouter.Handle {
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
		case MD5:
			h = md5.New()
		case SHA1:
			h = sha1.New()
		case SHA224:
			h = sha256.New224()
		case SHA256:
			h = sha256.New()
		case SHA384:
			h = sha512.New384()
		case SHA512:
			h = sha512.New()
		case SHA512_224:
			h = sha512.New512_224()
		case SHA512_256:
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

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}

		_, err = w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}
	}
}

func registerHash(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "hash"

	mux.GET("/hash/", serveUsage(module, usage, errorChannel))

	mux.GET("/hash/md5/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/md5/:string", serveHash(MD5, errorChannel))
	mux.POST("/hash/md5/", serveHash(MD5, errorChannel))

	mux.GET("/hash/sha1/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha1/:string", serveHash(SHA1, errorChannel))
	mux.POST("/hash/sha1/", serveHash(SHA1, errorChannel))

	mux.GET("/hash/sha224/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha224/:string", serveHash(SHA224, errorChannel))
	mux.POST("/hash/sha224/", serveHash(SHA224, errorChannel))

	mux.GET("/hash/sha256/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha256/:string", serveHash(SHA256, errorChannel))
	mux.POST("/hash/sha256/", serveHash(SHA256, errorChannel))

	mux.GET("/hash/sha384/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha384/:string", serveHash(SHA384, errorChannel))
	mux.POST("/hash/sha384/", serveHash(SHA384, errorChannel))

	mux.GET("/hash/sha512/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha512/:string", serveHash(SHA512, errorChannel))
	mux.POST("/hash/sha512/", serveHash(SHA512, errorChannel))

	mux.GET("/hash/sha512-224/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha512-224/:string", serveHash(SHA512_224, errorChannel))
	mux.POST("/hash/sha512-224/", serveHash(SHA512_224, errorChannel))

	mux.GET("/hash/sha512-256/", serveUsage(module, usage, errorChannel))
	mux.GET("/hash/sha512-256/:string", serveHash(SHA512_256, errorChannel))
	mux.POST("/hash/sha512-256/", serveHash(SHA512_256, errorChannel))

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
