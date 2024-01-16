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

func parseOUIs() (map[string]string, error) {
	startTime := time.Now()

	whiteSpace := regexp.MustCompile(`\s+`)

	retVal := make(map[string]string)

	var readFile fs.File
	var err error

	if ouiFile == "" {
		readFile, err = ouis.Open("oui.txt")
	} else {
		readFile, err = os.Open(ouiFile)
	}
	defer readFile.Close()
	if err != nil {
		return retVal, err
	}

	scanner := bufio.NewScanner(readFile)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		oui, vendor := format(line, whiteSpace)

		for i := 0; i < len(oui); i++ {
			retVal[oui[i]] = vendor
		}
	}

	if verbose {
		fmt.Printf("%s | Loaded OUI database in %dms\n",
			startTime.Format(timeFormats["RFC3339"]),
			time.Since(startTime).Milliseconds())
	}

	return retVal, err
}

func serveMAC(ouis map[string]string, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		mac := strings.TrimPrefix(p.ByName("mac"), "/")

		val := ""

		for i := 12; i >= 6; i -= 2 {
			v, ok := ouis[strings.Join(chunks(firstN(strip(strings.ToUpper(mac)), i), 2), ":")]

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

func registerMAC(module string, mux *httprouter.Router, usage map[string][]string, errorChannel chan<- Error) ([]string, error) {
	ouis, err := parseOUIs()
	if err != nil {
		return []string{}, err
	}

	mux.GET("/mac/:mac", serveMAC(ouis, errorChannel))
	mux.GET("/mac/", serveUsage(module, usage))

	examples := make([]string, 3)
	examples[0] = "/mac/3c-7c-3f-1e-b9-a0"
	examples[1] = "/mac/e0:00:84:aa:aa:bb"
	examples[2] = "/mac/4C445BAABBCC"

	return examples, nil
}
