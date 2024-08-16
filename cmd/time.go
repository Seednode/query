/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var timeFormats = map[string]string{
	"ANSIC":       `Mon Jan _2 15:04:05 2006`,
	"DateOnly":    `2006-01-02`,
	"DateTime":    `2006-01-02 15:04:05`,
	"Kitchen":     `3:04PM`,
	"Layout":      `01/02 03:04:05PM '06 -0700`,
	"RFC1123":     `Mon, 02 Jan 2006 15:04:05 MST`,
	"RFC1123Z":    `Mon, 02 Jan 2006 15:04:05 -0700`,
	"RFC3339":     `2006-01-02T15:04:05Z07:00`,
	"RFC3339Nano": `2006-01-02T15:04:05.999999999Z07:00`,
	"RFC822":      `02 Jan 06 15:04 MST`,
	"RFC822Z":     `02 Jan 06 15:04 -0700`,
	"RFC850":      `Monday, 02-Jan-06 15:04:05 MST`,
	"RubyDate":    `Mon Jan 02 15:04:05 -0700 2006`,
	"Stamp":       `Jan _2 15:04:05`,
	"StampMicro":  `Jan _2 15:04:05.000000`,
	"StampMilli":  `Jan _2 15:04:05.000`,
	"StampNano":   `Jan _2 15:04:05.000000000`,
	"TimeOnly":    `15:04:05`,
	"UnixDate":    `Mon Jan _2 15:04:05 MST 2006`,
}

var timeAbbrevations = map[string][]*time.Location{
	"ACDT": {
		time.FixedZone("Australian Central Daylight Saving Time", 10.5*60*60),
	},
	"ACST": {
		time.FixedZone("Australian Central Standard Time", 9.5*60*60),
	},
	"ACT": {
		time.FixedZone("Acre Time", -5*60*60),
		time.FixedZone("ASEAN Common Time", 8*60*60),
	},
	"ACWST": {
		time.FixedZone("Australian Central Western Standard Time", 8.75*60*60),
	},
	"ADT": {
		time.FixedZone("Atlantic Daylight Time", -3*60*60),
	},
	"AEDT": {
		time.FixedZone("Australian Eastern Daylight Saving Time", 11*60*60),
	},
	"AEST": {
		time.FixedZone("Australian Eastern Standard Time", 10*60*60),
	},
	"AFT": {
		time.FixedZone("Afghanistan Time", 4.5*60*60),
	},
	"AKDT": {
		time.FixedZone("Alaska Daylight Time", -8*60*60),
	},
	"AKST": {
		time.FixedZone("Alaska Standard Time", -9*60*60),
	},
	"ALMT": {
		time.FixedZone("Alma-Ata Time", 6*60*60),
	},
	"AMST": {
		time.FixedZone("Amazon Summer Time", -3*60*60),
	},
	"AMT": {
		time.FixedZone("Amazon Time", -4*60*60),
		time.FixedZone("Armenia Time", 4*60*60),
	},
	"ANAT": {
		time.FixedZone("Anadyr Time", 12*60*60),
	},
	"AQTT": {
		time.FixedZone("Aqtobe Time", 5*60*60),
	},
	"ART": {
		time.FixedZone("Argentina Time", -3*60*60),
	},
	"AST": {
		time.FixedZone("Arabia Standard Time", 3*60*60),
		time.FixedZone("Atlantic Standard Time", -4*60*60),
	},
	"AWST": {
		time.FixedZone("Australian Western Standard Time", 8*60*60),
	},
	"AZOST": {
		time.FixedZone("Azores Summer Time", 0*60*60),
	},
	"AZOT": {
		time.FixedZone("Azores Standard Time", -1*60*60),
	},
	"AZT": {
		time.FixedZone("Azerbaijan Time", 4*60*60),
	},
	"BNT": {
		time.FixedZone("Brunei Time", 8*60*60),
	},
	"BIOT": {
		time.FixedZone("British Indian Ocean Time", 6*60*60),
	},
	"BIT": {
		time.FixedZone("Baker Island Time", -12*60*60),
	},
	"BOT": {
		time.FixedZone("Bolivia Time", -4*60*60),
	},
	"BRST": {
		time.FixedZone("Brasília Summer Time", -2*60*60),
	},
	"BRT": {
		time.FixedZone("Brasília Time", -3*60*60),
	},
	"BST": {
		time.FixedZone("Bangladesh Standard Time", 6*60*60),
		time.FixedZone("Bougainville Standard Time", 11*60*60),
		time.FixedZone("British Summer Time", 1*60*60),
	},
	"BTT": {
		time.FixedZone("Bhutan Time", 6*60*60),
	},
	"CAT": {
		time.FixedZone("Central Africa Time", 2*60*60),
	},
	"CCT": {
		time.FixedZone("Cocos Islands Time", 6.5*60*60),
	},
	"CDT": {
		time.FixedZone("Central Daylight Time", -5*60*60),
		time.FixedZone("Cuba Daylight Time", -4*60*60),
	},
	"CEST": {
		time.FixedZone("Central European Summer Time", 2*60*60),
	},
	"CET": {
		time.FixedZone("Central European Time", 1*60*60),
	},
	"CHADT": {
		time.FixedZone("Chatham Daylight Time", 13.75*60*60),
	},
	"CHAST": {
		time.FixedZone("Chatham Standard Time", 12.75*60*60),
	},
	"CHOT": {
		time.FixedZone("Choibalsan Standard Time", 8*60*60),
	},
	"CHOST": {
		time.FixedZone("Choibalsan Summer Time", 9*60*60),
	},
	"CHST": {
		time.FixedZone("Chamorro Standard Time", 10*60*60),
	},
	"CHUT": {
		time.FixedZone("Chuuk Time", 10*60*60),
	},
	"CIST": {
		time.FixedZone("Clipperton Island Standard Time", -8*60*60),
	},
	"CKT": {
		time.FixedZone("Cook Island Time", -10*60*60),
	},
	"CLST": {
		time.FixedZone("Chile Summer Time", -3*60*60),
	},
	"CLT": {
		time.FixedZone("Chile Standard Time", -4*60*60),
	},
	"COST": {
		time.FixedZone("Colombia Summer Time", -4*60*60),
	},
	"COT": {
		time.FixedZone("Colombia Time", -5*60*60),
	},
	"CST": {
		time.FixedZone("Central Standard Time", -6*60*60),
		time.FixedZone("China Standard Time", 8*60*60),
		time.FixedZone("Cuba Standard Time", -5*60*60),
	},
	"CVT": {
		time.FixedZone("Cape Verde Time", -1*60*60),
	},
	"CWST": {
		time.FixedZone("Central Western Standard Time", 8.75*60*60),
	},
	"CXT": {
		time.FixedZone("Christmas Island Time", 7*60*60),
	},
	"DAVT": {
		time.FixedZone("Davis Time", 7*60*60),
	},
	"DDUT": {
		time.FixedZone("Dumont d'Urville Time", 10*60*60),
	},
	"DFT": {
		time.FixedZone("AIX-specific equivalent of Central European Time", 1*60*60),
	},
	"EASST": {
		time.FixedZone("Easter Island Summer Time", -5*60*60),
	},
	"EAST": {
		time.FixedZone("Easter Island Standard Time", -6*60*60),
	},
	"EAT": {
		time.FixedZone("East Africa Time", 3*60*60),
	},
	"ECT": {
		time.FixedZone("Eastern Caribbean Time", -4*60*60),
		time.FixedZone("Ecuador Time", -5*60*60),
	},
	"EDT": {
		time.FixedZone("Eastern Daylight Time", -4*60*60),
	},
	"EEST": {
		time.FixedZone("Eastern European Summer Time", 3*60*60),
	},
	"EET": {
		time.FixedZone("Eastern European Time", 2*60*60),
	},
	"EGST": {
		time.FixedZone("Eastern Greenland Summer Time", 0*60*60),
	},
	"EGT": {
		time.FixedZone("Eastern Greenland Time", -1*60*60),
	},
	"EST": {
		time.FixedZone("Eastern Standard Time", -5*60*60),
	},
	"FET": {
		time.FixedZone("Further-eastern European Time", 3*60*60),
	},
	"FJT": {
		time.FixedZone("Fiji Time", 12*60*60),
	},
	"FKST": {
		time.FixedZone("Falkland Islands Summer Time", -3*60*60),
	},
	"FKT": {
		time.FixedZone("Falkland Islands Time", -4*60*60),
	},
	"FNT": {
		time.FixedZone("Fernando de Noronha Time", -2*60*60),
	},
	"GALT": {
		time.FixedZone("Galápagos Time", -6*60*60),
	},
	"GAMT": {
		time.FixedZone("Gambier Islands Time", -9*60*60),
	},
	"GET": {
		time.FixedZone("Georgia Standard Time", 4*60*60),
	},
	"GFT": {
		time.FixedZone("French Guiana Time", -3*60*60),
	},
	"GILT": {
		time.FixedZone("Gilbert Island Time", 12*60*60),
	},
	"GIT": {
		time.FixedZone("Gambier Island Time", -9*60*60),
	},
	"GMT": {
		time.FixedZone("Greenwich Mean Time", 0*60*60),
	},
	"GST": {
		time.FixedZone("South Georgia and the South Sandwich Islands Time", -2*60*60),
		time.FixedZone("Gulf Standard Time", 4*60*60),
	},
	"GYT": {
		time.FixedZone("Guyana Time", -4*60*60),
	},
	"HDT": {
		time.FixedZone("Hawaii–Aleutian Daylight Time", -9*60*60),
	},
	"HAEC": {
		time.FixedZone("Heure Avancée d'Europe Centrale", 2*60*60),
	},
	"HST": {
		time.FixedZone("Hawaii–Aleutian Standard Time", -10*60*60),
	},
	"HKT": {
		time.FixedZone("Hong Kong Time", 8*60*60),
	},
	"HMT": {
		time.FixedZone("Heard and McDonald Islands Time", 5*60*60),
	},
	"HOVST": {
		time.FixedZone("Hovd Summer Time", 8*60*60),
	},
	"HOVT": {
		time.FixedZone("Hovd Time", 7*60*60),
	},
	"ICT": {
		time.FixedZone("Indochina Time", 7*60*60),
	},
	"IDLW": {
		time.FixedZone("International Date Line West", -12*60*60),
	},
	"IDT": {
		time.FixedZone("Israel Daylight Time", 3*60*60),
	},
	"IOT": {
		time.FixedZone("Indian Ocean Time", 6*60*60),
	},
	"IRDT": {
		time.FixedZone("Iran Daylight Time", 4.5*60*60),
	},
	"IRKT": {
		time.FixedZone("Irkutsk Time", 8*60*60),
	},
	"IRST": {
		time.FixedZone("Iran Standard Time", 3.5*60*60),
	},
	"IST": {
		time.FixedZone("Indian Standard Time", 5.5*60*60),
		time.FixedZone("Irish Standard Time", 1*60*60),
		time.FixedZone("Israel Standard Time", 2*60*60),
	},
	"JST": {
		time.FixedZone("Japan Standard Time", 9*60*60),
	},
	"KALT": {
		time.FixedZone("Kaliningrad Time", 2*60*60),
	},
	"KGT": {
		time.FixedZone("Kyrgyzstan Time", 6*60*60),
	},
	"KOST": {
		time.FixedZone("Kosrae Time", 11*60*60),
	},
	"KRAT": {
		time.FixedZone("Krasnoyarsk Time", 7*60*60),
	},
	"KST": {
		time.FixedZone("Korea Standard Time", 9*60*60),
	},
	"LHST": {
		time.FixedZone("Lord Howe Standard Time", 10.5*60*60),
		time.FixedZone("Lord Howe Summer Time", 11*60*60),
	},
	"LINT": {
		time.FixedZone("Line Islands Time", 14*60*60),
	},
	"MAGT": {
		time.FixedZone("Magadan Time", 12*60*60),
	},
	"MART": {
		time.FixedZone("Marquesas Islands Time", -9.5*60*60),
	},
	"MAWT": {
		time.FixedZone("Mawson Station Time", 5*60*60),
	},
	"MDT": {
		time.FixedZone("Mountain Daylight Time", -6*60*60),
	},
	"MET": {
		time.FixedZone("Middle European Time", 1*60*60),
	},
	"MEST": {
		time.FixedZone("Middle European Summer Time", 2*60*60),
	},
	"MHT": {
		time.FixedZone("Marshall Islands Time", 12*60*60),
	},
	"MIST": {
		time.FixedZone("Macquarie Island Station Time", 11*60*60),
	},
	"MIT": {
		time.FixedZone("Marquesas Islands Time", -9.5*60*60),
	},
	"MMT": {
		time.FixedZone("Myanmar Standard Time", 6.5*60*60),
	},
	"MSK": {
		time.FixedZone("Moscow Time", 3*60*60),
	},
	"MST": {
		time.FixedZone("Malaysia Standard Time", 8*60*60),
		time.FixedZone("Mountain Standard Time", -7*60*60),
	},
	"MUT": {
		time.FixedZone("Mauritius Time", 4*60*60),
	},
	"MVT": {
		time.FixedZone("Maldives Time", 5*60*60),
	},
	"MYT": {
		time.FixedZone("Malaysia Time", 8*60*60),
	},
	"NCT": {
		time.FixedZone("New Caledonia Time", 11*60*60),
	},
	"NDT": {
		time.FixedZone("Newfoundland Daylight Time", -2.5*60*60),
	},
	"NFT": {
		time.FixedZone("Norfolk Island Time", 11*60*60),
	},
	"NOVT": {
		time.FixedZone("Novosibirsk Time", 7*60*60),
	},
	"NPT": {
		time.FixedZone("Nepal Time", 5.75*60*60),
	},
	"NST": {
		time.FixedZone("Newfoundland Standard Time", -3.5*60*60),
	},
	"NT": {
		time.FixedZone("Newfoundland Time", -3.5*60*60),
	},
	"NUT": {
		time.FixedZone("Niue Time", -11*60*60),
	},
	"NZDT": {
		time.FixedZone("New Zealand Daylight Time", 13*60*60),
	},
	"NZST": {
		time.FixedZone("New Zealand Standard Time", 12*60*60),
	},
	"OMST": {
		time.FixedZone("Omsk Time", 6*60*60),
	},
	"ORAT": {
		time.FixedZone("Oral Time", 5*60*60),
	},
	"PDT": {
		time.FixedZone("Pacific Daylight Time", -7*60*60),
	},
	"PET": {
		time.FixedZone("Peru Time", -5*60*60),
	},
	"PETT": {
		time.FixedZone("Kamchatka Time", 12*60*60),
	},
	"PGT": {
		time.FixedZone("Papua New Guinea Time", 10*60*60),
	},
	"PHOT": {
		time.FixedZone("Phoenix Island Time", 13*60*60),
	},
	"PHT": {
		time.FixedZone("Philippine Time", 8*60*60),
	},
	"PHST": {
		time.FixedZone("Philippine Standard Time", 8*60*60),
	},
	"PKT": {
		time.FixedZone("Pakistan Standard Time", 5*60*60),
	},
	"PMDT": {
		time.FixedZone("Saint Pierre and Miquelon Daylight Time", -2*60*60),
	},
	"PMST": {
		time.FixedZone("Saint Pierre and Miquelon Standard Time", -3*60*60),
	},
	"PONT": {
		time.FixedZone("Pohnpei Standard Time", 11*60*60),
	},
	"PST": {
		time.FixedZone("Pacific Standard Time", -8*60*60),
	},
	"PWT": {
		time.FixedZone("Palau Time", 9*60*60),
	},
	"PYST": {
		time.FixedZone("Paraguay Summer Time", -3*60*60),
	},
	"PYT": {
		time.FixedZone("Paraguay Time", -4*60*60),
	},
	"RET": {
		time.FixedZone("Réunion Time", 4*60*60),
	},
	"ROTT": {
		time.FixedZone("Rothera Research Station Time", -3*60*60),
	},
	"SAKT": {
		time.FixedZone("Sakhalin Island Time", 11*60*60),
	},
	"SAMT": {
		time.FixedZone("Samara Time", 4*60*60),
	},
	"SAST": {
		time.FixedZone("South African Standard Time", 2*60*60),
	},
	"SBT": {
		time.FixedZone("Solomon Islands Time", 11*60*60),
	},
	"SCT": {
		time.FixedZone("Seychelles Time", 4*60*60),
	},
	"SDT": {
		time.FixedZone("Samoa Daylight Time", -10*60*60),
	},
	"SGT": {
		time.FixedZone("Singapore Time", 8*60*60),
	},
	"SLST": {
		time.FixedZone("Sri Lanka Standard Time", 5.5*60*60),
	},
	"SRET": {
		time.FixedZone("Srednekolymsk Time", 11*60*60),
	},
	"SRT": {
		time.FixedZone("Suriname Time", -3*60*60),
	},
	"SST": {
		time.FixedZone("Samoa Standard Time", -11*60*60),
	},
	"SYOT": {
		time.FixedZone("Showa Station Time", 3*60*60),
	},
	"TAHT": {
		time.FixedZone("Tahiti Time", -10*60*60),
	},
	"THA": {
		time.FixedZone("Thailand Standard Time", 7*60*60),
	},
	"TFT": {
		time.FixedZone("French Southern and Antarctic Time", 5*60*60),
	},
	"TJT": {
		time.FixedZone("Tajikistan Time", 5*60*60),
	},
	"TKT": {
		time.FixedZone("Tokelau Time", 13*60*60),
	},
	"TLT": {
		time.FixedZone("Timor Leste Time", 9*60*60),
	},
	"TMT": {
		time.FixedZone("Turkmenistan Time", 5*60*60),
	},
	"TRT": {
		time.FixedZone("Turkey Time", 3*60*60),
	},
	"TOT": {
		time.FixedZone("Tonga Time", 13*60*60),
	},
	"TST": {
		time.FixedZone("Taiwan Standard Time", 8*60*60),
	},
	"TVT": {
		time.FixedZone("Tuvalu Time", 12*60*60),
	},
	"ULAST": {
		time.FixedZone("Ulaanbaatar Summer Time", 9*60*60),
	},
	"ULAT": {
		time.FixedZone("Ulaanbaatar Standard Time", 8*60*60),
	},
	"UTC": {
		time.FixedZone("Coordinated Universal Time", 0*60*60),
	},
	"UYST": {
		time.FixedZone("Uruguay Summer Time", -2*60*60),
	},
	"UYT": {
		time.FixedZone("Uruguay Standard Time", -3*60*60),
	},
	"UZT": {
		time.FixedZone("Uzbekistan Time", 5*60*60),
	},
	"VET": {
		time.FixedZone("Venezuelan Standard Time", -4*60*60),
	},
	"VLAT": {
		time.FixedZone("Vladivostok Time", 10*60*60),
	},
	"VOLT": {
		time.FixedZone("Volgograd Time", 3*60*60),
	},
	"VOST": {
		time.FixedZone("Vostok Station Time", 6*60*60),
	},
	"VUT": {
		time.FixedZone("Vanuatu Time", 11*60*60),
	},
	"WAKT": {
		time.FixedZone("Wake Island Time", 12*60*60),
	},
	"WAST": {
		time.FixedZone("West Africa Summer Time", 2*60*60),
	},
	"WAT": {
		time.FixedZone("West Africa Time", 1*60*60),
	},
	"WEST": {
		time.FixedZone("Western European Summer Time", 1*60*60),
	},
	"WET": {
		time.FixedZone("Western European Time", 0*60*60),
	},
	"WIB": {
		time.FixedZone("Western Indonesian Time", 7*60*60),
	},
	"WIT": {
		time.FixedZone("Eastern Indonesian Time", 9*60*60),
	},
	"WITA": {
		time.FixedZone("Central Indonesia Time", 8*60*60),
	},
	"WGST": {
		time.FixedZone("West Greenland Summer Time", -2*60*60),
	},
	"WGT": {
		time.FixedZone("West Greenland Time", -3*60*60),
	},
	"WST": {
		time.FixedZone("Western Standard Time", 8*60*60),
	},
	"YAKT": {
		time.FixedZone("Yakutsk Time", 9*60*60),
	},
	"YEKT": {
		time.FixedZone("Yekaterinburg Time", 5*60*60),
	},
}

func serveTime(errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		var format string = ""

		requestedFormat := r.URL.Query().Get("format")
		if requestedFormat == "" {
			requestedFormat = "RFC822"
		}

		for k, v := range timeFormats {
			if strings.EqualFold(requestedFormat, k) {
				format = v

				break
			}
		}

		if format == "" {
			format = timeFormats["RFC822"]
		}

		adjustedStartTime := startTime

		location := strings.TrimPrefix(p.ByName("time"), "/") + p.ByName("rest")

		zones := []*time.Location{}

		abbrev, exists := timeAbbrevations[location]
		if exists {
			zones = append(zones, abbrev...)
		} else {
			tz, err := time.LoadLocation(location)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				w.WriteHeader(http.StatusBadRequest)

				_, err = w.Write([]byte("Invalid timezone requested\n"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			}

			zones = append(zones, tz)
		}

		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		if verbose {
			fmt.Printf("%s | %s => %s\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				r.RequestURI)
		}

		for i := 0; i < len(zones); i++ {
			adjustedStartTime = adjustedStartTime.In(zones[i])

			_, err := w.Write([]byte(adjustedStartTime.Format(format) + "\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				return
			}
		}
	}
}

func registerTime(mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) {
	const module = "time"

	mux.GET("/time/:time", serveTime(errorChannel))
	mux.GET("/time/:time/*rest", serveTime(errorChannel))
	mux.GET("/time/", serveUsage(module, usage, errorChannel))

	usage.Store(module, []string{
		"/time/America/Chicago",
		"/time/EST",
		"/time/UTC?format=kitchen",
	})
}
