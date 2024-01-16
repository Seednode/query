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
			s[len(s)-1] = strings.ToUpper(fmt.Sprintf("%x", i))
			ouis = append(ouis, strings.Join(s, ""))
		}
	} else {
		ouis = append(ouis, oui)
	}

	vendor := strings.Replace(s.String(), ",", " ", -1)

	re.ReplaceAllString(vendor, " ")

	return ouis, vendor
}

func scan() (func(), *bufio.Scanner, error) {
	var readFile fs.File
	var err error

	if ouiFile == "" {
		readFile, err = ouis.Open("oui.txt")
	} else {
		readFile, err = os.Open(ouiFile)
	}

	if err != nil {
		return func() {}, nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	buffer := make([]byte, 0, 64*1024)
	fileScanner.Buffer(buffer, 1024*1024)

	return func() { _ = readFile.Close() }, fileScanner, nil
}

func ouiMap() (map[string]string, error) {
	whiteSpace := regexp.MustCompile(`\s+`)

	retVal := make(map[string]string)

	closeFile, scanner, err := scan()
	if err != nil {
		return retVal, err
	}
	defer closeFile()

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		oui, vendor := format(line, whiteSpace)

		for i := 0; i < len(oui); i++ {
			retVal[oui[i]] = vendor
		}
	}

	return retVal, err
}

func serveOui(ouis map[string]string, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		mac := strings.TrimPrefix(p.ByName("mac"), "/")

		val := ""

		for i := 12; i >= 6; i -= 2 {
			s := chunks(firstN(strip(strings.ToUpper(mac)), i), 2)
			m := strings.Join(s, ":")

			v, ok := ouis[m]

			if ok {
				val = v

				break
			}
		}

		if val == "" {
			val = fmt.Sprintf("No OUI found for MAC %q\n", mac)
		}

		w.Write([]byte(val + "\n"))

		if verbose {
			fmt.Printf("%s | %s requested vendor info for MAC %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				mac)
		}
	}
}

func registerOUIHandlers(module string, mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) ([]string, error) {
	ouiMap, err := ouiMap()
	if err != nil {
		return []string{}, err
	}

	mux.GET("/mac/:mac", serveOui(ouiMap, errorChannel))
	mux.GET("/mac/", serveUsage(module, usage))

	examples := make([]string, 3)
	examples[0] = "/mac/3c-7c-3f-1e-b9-a0"
	examples[1] = "/mac/e0:00:84:aa:aa:bb"
	examples[2] = "/mac/4C445BAABBCC"

	return examples, nil
}
