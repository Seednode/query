/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func serveHashMd5(errorChannel chan<- error) httprouter.Handle {
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

		h := md5.New()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested MD5 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha1(errorChannel chan<- error) httprouter.Handle {
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

		h := sha1.New()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA1 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha224(errorChannel chan<- error) httprouter.Handle {
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

		h := sha256.New224()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA-224 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha256(errorChannel chan<- error) httprouter.Handle {
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

		h := sha256.New()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA-256 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha384(errorChannel chan<- error) httprouter.Handle {
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

		h := sha512.New384()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA-384 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha512(errorChannel chan<- error) httprouter.Handle {
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

		h := sha512.New()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA-512 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha512_224(errorChannel chan<- error) httprouter.Handle {
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

		h := sha512.New512_224()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA-512/224 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func serveHashSha512_256(errorChannel chan<- error) httprouter.Handle {
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

		h := sha512.New512_256()

		io.WriteString(h, value)

		w.Write([]byte(fmt.Sprintf("%x\n", h.Sum(nil))))

		if verbose {
			fmt.Printf("%s | %s requested SHA-512/256 hash of %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				value)
		}
	}
}

func registerHashHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	mux.GET("/hash/md5/*md5", serveHashMd5(errorChannel))
	mux.GET("/hash/sha1/*sha1", serveHashSha1(errorChannel))
	mux.GET("/hash/sha224/*sha224", serveHashSha224(errorChannel))
	mux.GET("/hash/sha256/*sha256", serveHashSha256(errorChannel))
	mux.GET("/hash/sha384/*sha384", serveHashSha384(errorChannel))
	mux.GET("/hash/sha512/*sha512", serveHashSha512(errorChannel))
	mux.GET("/hash/sha512-224/*sha512_224", serveHashSha512_224(errorChannel))
	mux.GET("/hash/sha512-256/*sha512_256", serveHashSha512_256(errorChannel))

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
