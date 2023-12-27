/*
Copyright © 2023 Seednode <seednode@seedno.de>
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

func addColons(mac string) string {
	s := chunks(mac, 2)
	output := strings.Join(s, ":")

	return output
}

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

func format(mac, line string, re *regexp.Regexp) string {
	words := strings.Split(line, "\t")

	var retVal strings.Builder

	if len(words) < 2 {
		return ""
	}

	if len(words) < 3 {
		retVal.WriteString(prettify(words[1], re))
	} else {
		for i := 2; i < len(words); i++ {
			retVal.WriteString(prettify(words[i], re))
		}
	}

	return retVal.String()
}

func getOui(mac string, re *regexp.Regexp) (string, error) {
	normalizedMac := normalize(mac)
	trimmedMac := trim(normalizedMac)

	closeFile, scanner, err := scan()
	if err != nil {
		return "", err
	}
	defer closeFile()

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if trim(normalize(line)) == trimmedMac {
			return format(addColons(mac), line, re), nil
		}

	}

	return "", nil
}

func normalize(oui string) string {
	oui = strings.ToUpper(oui)
	oui = strip(oui)

	return oui
}

func prettify(s string, re *regexp.Regexp) string {
	s = strings.Replace(s, ",", " ", -1)

	return re.ReplaceAllString(s, " ")
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

func trim(oui string) string {
	return firstN(oui, 6)
}

func serveOui(re *regexp.Regexp, errorChannel chan<- error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		mac := strings.TrimPrefix(p[0].Value, "/")

		oui, err := getOui(mac, re)
		if err != nil {
			errorChannel <- err

			return
		}

		if oui == "" {
			oui = fmt.Sprintf("No OUI found for MAC %q\n", mac)
		}

		w.Write([]byte(oui + "\n"))

		if verbose {
			fmt.Printf("%s | %s looked up OUI for %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				mac)
		}
	}
}

func registerOUIHandlers(mux *httprouter.Router, errorChannel chan<- error) []string {
	whiteSpaceRegex := regexp.MustCompile(`\s+`)

	mux.GET("/oui/*oui", serveOui(whiteSpaceRegex, errorChannel))

	var usage []string
	usage = append(usage, "/oui/00:00:08")
	usage = append(usage, "/oui/00-50-C2")
	usage = append(usage, "/oui/70b3d5")

	return usage
}
