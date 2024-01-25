/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

//go:embed oui.txt
var ouis embed.FS

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}

	if chunkSize >= len(s) {
		return []string{s}
	}

	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)

	currentLen := 0
	currentStart := 0

	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}

		currentLen++
	}

	chunks = append(chunks, s[currentStart:])

	return chunks
}

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]

		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}

	return result.String()
}

func format(line string, re *regexp.Regexp) ([]string, string) {
	ouis := []string{}

	words := strings.Split(line, "\t")

	if len(words) < 2 {
		return ouis, ""
	}

	var s strings.Builder

	if len(words) < 3 {
		s.WriteString(words[1])
	} else {
		for i := 2; i < len(words); i++ {
			s.WriteString(words[i])
		}
	}

	oui, _, isRange := strings.Cut(strings.TrimSpace(words[0]), "/")

	if isRange {
		for i := 0; i < 16; i++ {
			s := strings.Split(oui, "")
			s[len(s)-1] = fmt.Sprintf("%X", i)
			ouis = append(ouis, strings.Join(s, ""))
		}
	} else {
		ouis = append(ouis, oui)
	}

	vendor := strings.Replace(s.String(), ",", " ", -1)

	re.ReplaceAllString(vendor, " ")

	return ouis, vendor
}

func parseOUIs(errorChannel chan<- Error) *sync.Map {
	startTime := time.Now()

	whiteSpace := regexp.MustCompile(`\s+`)

	retVal := sync.Map{}

	var readFile fs.File
	var err error

	if ouiFile == "" {
		readFile, err = ouis.Open("oui.txt")
		if err != nil {
			errorChannel <- Error{Message: err, Path: "parseOUIs()"}

			return &retVal
		}
	} else {
		readFile, err = os.Open(ouiFile)
		if err != nil {
			errorChannel <- Error{Message: err, Path: "parseOUIs()"}

			return &retVal
		}
	}
	defer func() {
		err = readFile.Close()
		if err != nil {
			errorChannel <- Error{Message: err, Path: "parseOUIs()"}
		}
	}()

	scanner := bufio.NewScanner(readFile)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		oui, vendor := format(line, whiteSpace)

		for i := 0; i < len(oui); i++ {
			retVal.Store(oui[i], vendor)
		}
	}

	if verbose {
		fmt.Printf("%s | Loaded OUI database in %dms\n",
			startTime.Format(timeFormats["RFC3339"]),
			time.Since(startTime).Milliseconds())
	}

	return &retVal
}

func serveMAC(ouis *sync.Map, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		mac := strings.TrimPrefix(p.ByName("mac"), "/")

		val := ""

		for i := 12; i >= 6; i -= 2 {
			v, ok := ouis.Load(strings.Join(chunks(firstN(strip(strings.ToUpper(mac)), i), 2), ":"))

			if ok {
				val = v.(string)

				break
			}
		}

		if val == "" {
			val = fmt.Sprintf("No OUI found for MAC %q", mac)
		}

		_, err := w.Write([]byte(val + "\n"))
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested vendor info for MAC %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				mac)
		}
	}
}

func registerMAC(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "mac"

	ouis := parseOUIs(errorChannel)

	mux.GET("/mac/:mac", serveMAC(ouis, errorChannel))
	mux.GET("/mac/", serveUsage(module, usage, errorChannel))

	usage.Store(module, []string{
		"/mac/3c-7c-3f-1e-b9-a0",
		"/mac/e0:00:84:aa:aa:bb",
		"/mac/4C445BAABBCC",
	})
}
