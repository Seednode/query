/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

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

func getTimeAbbrevations() *sync.Map {
	retVal := sync.Map{}

	retVal.Store("ACDT", []*time.Location{
		time.FixedZone("Australian Central Daylight Saving Time", 10.5*60*60),
	})
	retVal.Store("ACST", []*time.Location{
		time.FixedZone("Australian Central Standard Time", 9.5*60*60),
	})
	retVal.Store("ACT", []*time.Location{
		time.FixedZone("Acre Time", -5*60*60),
		time.FixedZone("ASEAN Common Time", 8*60*60),
	})
	retVal.Store("ACWST", []*time.Location{
		time.FixedZone("Australian Central Western Standard Time", 8.75*60*60),
	})
	retVal.Store("ADT", []*time.Location{
		time.FixedZone("Atlantic Daylight Time", -3*60*60),
	})
	retVal.Store("AEDT", []*time.Location{
		time.FixedZone("Australian Eastern Daylight Saving Time", 11*60*60),
	})
	retVal.Store("AEST", []*time.Location{
		time.FixedZone("Australian Eastern Standard Time", 10*60*60),
	})
	retVal.Store("AFT", []*time.Location{
		time.FixedZone("Afghanistan Time", 4.5*60*60),
	})
	retVal.Store("AKDT", []*time.Location{
		time.FixedZone("Alaska Daylight Time", -8*60*60),
	})
	retVal.Store("AKST", []*time.Location{
		time.FixedZone("Alaska Standard Time", -9*60*60),
	})
	retVal.Store("ALMT", []*time.Location{
		time.FixedZone("Alma-Ata Time", 6*60*60),
	})
	retVal.Store("AMST", []*time.Location{
		time.FixedZone("Amazon Summer Time", -3*60*60),
	})
	retVal.Store("AMT", []*time.Location{
		time.FixedZone("Amazon Time", -4*60*60),
		time.FixedZone("Armenia Time", 4*60*60),
	})
	retVal.Store("ANAT", []*time.Location{
		time.FixedZone("Anadyr Time", 12*60*60),
	})
	retVal.Store("AQTT", []*time.Location{
		time.FixedZone("Aqtobe Time", 5*60*60),
	})
	retVal.Store("ART", []*time.Location{
		time.FixedZone("Argentina Time", -3*60*60),
	})
	retVal.Store("AST", []*time.Location{
		time.FixedZone("Arabia Standard Time", 3*60*60),
		time.FixedZone("Atlantic Standard Time", -4*60*60),
	})
	retVal.Store("AWST", []*time.Location{
		time.FixedZone("Australian Western Standard Time", 8*60*60),
	})
	retVal.Store("AZOST", []*time.Location{
		time.FixedZone("Azores Summer Time", 0),
	})
	retVal.Store("AZOT", []*time.Location{
		time.FixedZone("Azores Standard Time", -1*60*60),
	})
	retVal.Store("AZT", []*time.Location{
		time.FixedZone("Azerbaijan Time", 4*60*60),
	})
	retVal.Store("BNT", []*time.Location{
		time.FixedZone("Brunei Time", 8*60*60),
	})
	retVal.Store("BIOT", []*time.Location{
		time.FixedZone("British Indian Ocean Time", 6*60*60),
	})
	retVal.Store("BIT", []*time.Location{
		time.FixedZone("Baker Island Time", -12*60*60),
	})
	retVal.Store("BOT", []*time.Location{
		time.FixedZone("Bolivia Time", -4*60*60),
	})
	retVal.Store("BRST", []*time.Location{
		time.FixedZone("Brasília Summer Time", -2*60*60),
	})
	retVal.Store("BRT", []*time.Location{
		time.FixedZone("Brasília Time", -3*60*60),
	})
	retVal.Store("BST", []*time.Location{
		time.FixedZone("Bangladesh Standard Time", 6*60*60),
		time.FixedZone("Bougainville Standard Time", 11*60*60),
		time.FixedZone("British Summer Time", 1*60*60),
	})
	retVal.Store("BTT", []*time.Location{
		time.FixedZone("Bhutan Time", 6*60*60),
	})
	retVal.Store("CAT", []*time.Location{
		time.FixedZone("Central Africa Time", 2*60*60),
	})
	retVal.Store("CCT", []*time.Location{
		time.FixedZone("Cocos Islands Time", 6.5*60*60),
	})
	retVal.Store("CDT", []*time.Location{
		time.FixedZone("Central Daylight Time", -5*60*60),
		time.FixedZone("Cuba Daylight Time", -4*60*60),
	})
	retVal.Store("CEST", []*time.Location{
		time.FixedZone("Central European Summer Time", 2*60*60),
	})
	retVal.Store("CET", []*time.Location{
		time.FixedZone("Central European Time", 1*60*60),
	})
	retVal.Store("CHADT", []*time.Location{
		time.FixedZone("Chatham Daylight Time", 13.75*60*60),
	})
	retVal.Store("CHAST", []*time.Location{
		time.FixedZone("Chatham Standard Time", 12.75*60*60),
	})
	retVal.Store("CHOT", []*time.Location{
		time.FixedZone("Choibalsan Standard Time", 8*60*60),
	})
	retVal.Store("CHOST", []*time.Location{
		time.FixedZone("Choibalsan Summer Time", 9*60*60),
	})
	retVal.Store("CHST", []*time.Location{
		time.FixedZone("Chamorro Standard Time", 10*60*60),
	})
	retVal.Store("CHUT", []*time.Location{
		time.FixedZone("Chuuk Time", 10*60*60),
	})
	retVal.Store("CIST", []*time.Location{
		time.FixedZone("Clipperton Island Standard Time", -8*60*60),
	})
	retVal.Store("CKT", []*time.Location{
		time.FixedZone("Cook Island Time", -10*60*60),
	})
	retVal.Store("CLST", []*time.Location{
		time.FixedZone("Chile Summer Time", -3*60*60),
	})
	retVal.Store("CLT", []*time.Location{
		time.FixedZone("Chile Standard Time", -4*60*60),
	})
	retVal.Store("COST", []*time.Location{
		time.FixedZone("Colombia Summer Time", -4*60*60),
	})
	retVal.Store("COT", []*time.Location{
		time.FixedZone("Colombia Time", -5*60*60),
	})
	retVal.Store("CST", []*time.Location{
		time.FixedZone("Central Standard Time", -6*60*60),
		time.FixedZone("China Standard Time", 8*60*60),
		time.FixedZone("Cuba Standard Time", -5*60*60),
	})
	retVal.Store("CVT", []*time.Location{
		time.FixedZone("Cape Verde Time", -1*60*60),
	})
	retVal.Store("CWST", []*time.Location{
		time.FixedZone("Central Western Standard Time", 8.75*60*60),
	})
	retVal.Store("CXT", []*time.Location{
		time.FixedZone("Christmas Island Time", 7*60*60),
	})
	retVal.Store("DAVT", []*time.Location{
		time.FixedZone("Davis Time", 7*60*60),
	})
	retVal.Store("DDUT", []*time.Location{
		time.FixedZone("Dumont d'Urville Time", 10*60*60),
	})
	retVal.Store("DFT", []*time.Location{
		time.FixedZone("AIX-specific equivalent of Central European Time", 1*60*60),
	})
	retVal.Store("EASST", []*time.Location{
		time.FixedZone("Easter Island Summer Time", -5*60*60),
	})
	retVal.Store("EAST", []*time.Location{
		time.FixedZone("Easter Island Standard Time", -6*60*60),
	})
	retVal.Store("EAT", []*time.Location{
		time.FixedZone("East Africa Time", 3*60*60),
	})
	retVal.Store("ECT", []*time.Location{
		time.FixedZone("Eastern Caribbean Time", -4*60*60),
		time.FixedZone("Ecuador Time", -5*60*60),
	})
	retVal.Store("EDT", []*time.Location{
		time.FixedZone("Eastern Daylight Time", -4*60*60),
	})
	retVal.Store("EEST", []*time.Location{
		time.FixedZone("Eastern European Summer Time", 3*60*60),
	})
	retVal.Store("EET", []*time.Location{
		time.FixedZone("Eastern European Time", 2*60*60),
	})
	retVal.Store("EGST", []*time.Location{
		time.FixedZone("Eastern Greenland Summer Time", 0),
	})
	retVal.Store("EGT", []*time.Location{
		time.FixedZone("Eastern Greenland Time", -1*60*60),
	})
	retVal.Store("EST", []*time.Location{
		time.FixedZone("Eastern Standard Time", -5*60*60),
	})
	retVal.Store("FET", []*time.Location{
		time.FixedZone("Further-eastern European Time", 3*60*60),
	})
	retVal.Store("FJT", []*time.Location{
		time.FixedZone("Fiji Time", 12*60*60),
	})
	retVal.Store("FKST", []*time.Location{
		time.FixedZone("Falkland Islands Summer Time", -3*60*60),
	})
	retVal.Store("FKT", []*time.Location{
		time.FixedZone("Falkland Islands Time", -4*60*60),
	})
	retVal.Store("FNT", []*time.Location{
		time.FixedZone("Fernando de Noronha Time", -2*60*60),
	})
	retVal.Store("GALT", []*time.Location{
		time.FixedZone("Galápagos Time", -6*60*60),
	})
	retVal.Store("GAMT", []*time.Location{
		time.FixedZone("Gambier Islands Time", -9*60*60),
	})
	retVal.Store("GET", []*time.Location{
		time.FixedZone("Georgia Standard Time", 4*60*60),
	})
	retVal.Store("GFT", []*time.Location{
		time.FixedZone("French Guiana Time", -3*60*60),
	})
	retVal.Store("GILT", []*time.Location{
		time.FixedZone("Gilbert Island Time", 12*60*60),
	})
	retVal.Store("GIT", []*time.Location{
		time.FixedZone("Gambier Island Time", -9*60*60),
	})
	retVal.Store("GMT", []*time.Location{
		time.FixedZone("Greenwich Mean Time", 0),
	})
	retVal.Store("GST", []*time.Location{
		time.FixedZone("South Georgia and the South Sandwich Islands Time", -2*60*60),
		time.FixedZone("Gulf Standard Time", 4*60*60),
	})
	retVal.Store("GYT", []*time.Location{
		time.FixedZone("Guyana Time", -4*60*60),
	})
	retVal.Store("HDT", []*time.Location{
		time.FixedZone("Hawaii–Aleutian Daylight Time", -9*60*60),
	})
	retVal.Store("HAEC", []*time.Location{
		time.FixedZone("Heure Avancée d'Europe Centrale", 2*60*60),
	})
	retVal.Store("HST", []*time.Location{
		time.FixedZone("Hawaii–Aleutian Standard Time", -10*60*60),
	})
	retVal.Store("HKT", []*time.Location{
		time.FixedZone("Hong Kong Time", 8*60*60),
	})
	retVal.Store("HMT", []*time.Location{
		time.FixedZone("Heard and McDonald Islands Time", 5*60*60),
	})
	retVal.Store("HOVST", []*time.Location{
		time.FixedZone("Hovd Summer Time", 8*60*60),
	})
	retVal.Store("HOVT", []*time.Location{
		time.FixedZone("Hovd Time", 7*60*60),
	})
	retVal.Store("ICT", []*time.Location{
		time.FixedZone("Indochina Time", 7*60*60),
	})
	retVal.Store("IDLW", []*time.Location{
		time.FixedZone("International Date Line West", -12*60*60),
	})
	retVal.Store("IDT", []*time.Location{
		time.FixedZone("Israel Daylight Time", 3*60*60),
	})
	retVal.Store("IOT", []*time.Location{
		time.FixedZone("Indian Ocean Time", 6*60*60),
	})
	retVal.Store("IRDT", []*time.Location{
		time.FixedZone("Iran Daylight Time", 4.5*60*60),
	})
	retVal.Store("IRKT", []*time.Location{
		time.FixedZone("Irkutsk Time", 8*60*60),
	})
	retVal.Store("IRST", []*time.Location{
		time.FixedZone("Iran Standard Time", 3.5*60*60),
	})
	retVal.Store("IST", []*time.Location{
		time.FixedZone("Indian Standard Time", 5.5*60*60),
		time.FixedZone("Irish Standard Time", 1*60*60),
		time.FixedZone("Israel Standard Time", 2*60*60),
	})
	retVal.Store("JST", []*time.Location{
		time.FixedZone("Japan Standard Time", 9*60*60),
	})
	retVal.Store("KALT", []*time.Location{
		time.FixedZone("Kaliningrad Time", 2*60*60),
	})
	retVal.Store("KGT", []*time.Location{
		time.FixedZone("Kyrgyzstan Time", 6*60*60),
	})
	retVal.Store("KOST", []*time.Location{
		time.FixedZone("Kosrae Time", 11*60*60),
	})
	retVal.Store("KRAT", []*time.Location{
		time.FixedZone("Krasnoyarsk Time", 7*60*60),
	})
	retVal.Store("KST", []*time.Location{
		time.FixedZone("Korea Standard Time", 9*60*60),
	})
	retVal.Store("LHST", []*time.Location{
		time.FixedZone("Lord Howe Standard Time", 10.5*60*60),
		time.FixedZone("Lord Howe Summer Time", 11*60*60),
	})
	retVal.Store("LINT", []*time.Location{
		time.FixedZone("Line Islands Time", 14*60*60),
	})
	retVal.Store("MAGT", []*time.Location{
		time.FixedZone("Magadan Time", 12*60*60),
	})
	retVal.Store("MART", []*time.Location{
		time.FixedZone("Marquesas Islands Time", -9.5*60*60),
	})
	retVal.Store("MAWT", []*time.Location{
		time.FixedZone("Mawson Station Time", 5*60*60),
	})
	retVal.Store("MDT", []*time.Location{
		time.FixedZone("Mountain Daylight Time", -6*60*60),
	})
	retVal.Store("MET", []*time.Location{
		time.FixedZone("Middle European Time", 1*60*60),
	})
	retVal.Store("MEST", []*time.Location{
		time.FixedZone("Middle European Summer Time", 2*60*60),
	})
	retVal.Store("MHT", []*time.Location{
		time.FixedZone("Marshall Islands Time", 12*60*60),
	})
	retVal.Store("MIST", []*time.Location{
		time.FixedZone("Macquarie Island Station Time", 11*60*60),
	})
	retVal.Store("MIT", []*time.Location{
		time.FixedZone("Marquesas Islands Time", -9.5*60*60),
	})
	retVal.Store("MMT", []*time.Location{
		time.FixedZone("Myanmar Standard Time", 6.5*60*60),
	})
	retVal.Store("MSK", []*time.Location{
		time.FixedZone("Moscow Time", 3*60*60),
	})
	retVal.Store("MST", []*time.Location{
		time.FixedZone("Malaysia Standard Time", 8*60*60),
		time.FixedZone("Mountain Standard Time", -7*60*60),
	})
	retVal.Store("MUT", []*time.Location{
		time.FixedZone("Mauritius Time", 4*60*60),
	})
	retVal.Store("MVT", []*time.Location{
		time.FixedZone("Maldives Time", 5*60*60),
	})
	retVal.Store("MYT", []*time.Location{
		time.FixedZone("Malaysia Time", 8*60*60),
	})
	retVal.Store("NCT", []*time.Location{
		time.FixedZone("New Caledonia Time", 11*60*60),
	})
	retVal.Store("NDT", []*time.Location{
		time.FixedZone("Newfoundland Daylight Time", -2.5*60*60),
	})
	retVal.Store("NFT", []*time.Location{
		time.FixedZone("Norfolk Island Time", 11*60*60),
	})
	retVal.Store("NOVT", []*time.Location{
		time.FixedZone("Novosibirsk Time", 7*60*60),
	})
	retVal.Store("NPT", []*time.Location{
		time.FixedZone("Nepal Time", 5.75*60*60),
	})
	retVal.Store("NST", []*time.Location{
		time.FixedZone("Newfoundland Standard Time", -3.5*60*60),
	})
	retVal.Store("NT", []*time.Location{
		time.FixedZone("Newfoundland Time", -3.5*60*60),
	})
	retVal.Store("NUT", []*time.Location{
		time.FixedZone("Niue Time", -11*60*60),
	})
	retVal.Store("NZDT", []*time.Location{
		time.FixedZone("New Zealand Daylight Time", 13*60*60),
	})
	retVal.Store("NZST", []*time.Location{
		time.FixedZone("New Zealand Standard Time", 12*60*60),
	})
	retVal.Store("OMST", []*time.Location{
		time.FixedZone("Omsk Time", 6*60*60),
	})
	retVal.Store("ORAT", []*time.Location{
		time.FixedZone("Oral Time", 5*60*60),
	})
	retVal.Store("PDT", []*time.Location{
		time.FixedZone("Pacific Daylight Time", -7*60*60),
	})
	retVal.Store("PET", []*time.Location{
		time.FixedZone("Peru Time", -5*60*60),
	})
	retVal.Store("PETT", []*time.Location{
		time.FixedZone("Kamchatka Time", 12*60*60),
	})
	retVal.Store("PGT", []*time.Location{
		time.FixedZone("Papua New Guinea Time", 10*60*60),
	})
	retVal.Store("PHOT", []*time.Location{
		time.FixedZone("Phoenix Island Time", 13*60*60),
	})
	retVal.Store("PHT", []*time.Location{
		time.FixedZone("Philippine Time", 8*60*60),
	})
	retVal.Store("PHST", []*time.Location{
		time.FixedZone("Philippine Standard Time", 8*60*60),
	})
	retVal.Store("PKT", []*time.Location{
		time.FixedZone("Pakistan Standard Time", 5*60*60),
	})
	retVal.Store("PMDT", []*time.Location{
		time.FixedZone("Saint Pierre and Miquelon Daylight Time", -2*60*60),
	})
	retVal.Store("PMST", []*time.Location{
		time.FixedZone("Saint Pierre and Miquelon Standard Time", -3*60*60),
	})
	retVal.Store("PONT", []*time.Location{
		time.FixedZone("Pohnpei Standard Time", 11*60*60),
	})
	retVal.Store("PST", []*time.Location{
		time.FixedZone("Pacific Standard Time", -8*60*60),
	})
	retVal.Store("PWT", []*time.Location{
		time.FixedZone("Palau Time", 9*60*60),
	})
	retVal.Store("PYST", []*time.Location{
		time.FixedZone("Paraguay Summer Time", -3*60*60),
	})
	retVal.Store("PYT", []*time.Location{
		time.FixedZone("Paraguay Time", -4*60*60),
	})
	retVal.Store("RET", []*time.Location{
		time.FixedZone("Réunion Time", 4*60*60),
	})
	retVal.Store("ROTT", []*time.Location{
		time.FixedZone("Rothera Research Station Time", -3*60*60),
	})
	retVal.Store("SAKT", []*time.Location{
		time.FixedZone("Sakhalin Island Time", 11*60*60),
	})
	retVal.Store("SAMT", []*time.Location{
		time.FixedZone("Samara Time", 4*60*60),
	})
	retVal.Store("SAST", []*time.Location{
		time.FixedZone("South African Standard Time", 2*60*60),
	})
	retVal.Store("SBT", []*time.Location{
		time.FixedZone("Solomon Islands Time", 11*60*60),
	})
	retVal.Store("SCT", []*time.Location{
		time.FixedZone("Seychelles Time", 4*60*60),
	})
	retVal.Store("SDT", []*time.Location{
		time.FixedZone("Samoa Daylight Time", -10*60*60),
	})
	retVal.Store("SGT", []*time.Location{
		time.FixedZone("Singapore Time", 8*60*60),
	})
	retVal.Store("SLST", []*time.Location{
		time.FixedZone("Sri Lanka Standard Time", 5.5*60*60),
	})
	retVal.Store("SRET", []*time.Location{
		time.FixedZone("Srednekolymsk Time", 11*60*60),
	})
	retVal.Store("SRT", []*time.Location{
		time.FixedZone("Suriname Time", -3*60*60),
	})
	retVal.Store("SST", []*time.Location{
		time.FixedZone("Samoa Standard Time", -11*60*60),
	})
	retVal.Store("SYOT", []*time.Location{
		time.FixedZone("Showa Station Time", 3*60*60),
	})
	retVal.Store("TAHT", []*time.Location{
		time.FixedZone("Tahiti Time", -10*60*60),
	})
	retVal.Store("THA", []*time.Location{
		time.FixedZone("Thailand Standard Time", 7*60*60),
	})
	retVal.Store("TFT", []*time.Location{
		time.FixedZone("French Southern and Antarctic Time", 5*60*60),
	})
	retVal.Store("TJT", []*time.Location{
		time.FixedZone("Tajikistan Time", 5*60*60),
	})
	retVal.Store("TKT", []*time.Location{
		time.FixedZone("Tokelau Time", 13*60*60),
	})
	retVal.Store("TLT", []*time.Location{
		time.FixedZone("Timor Leste Time", 9*60*60),
	})
	retVal.Store("TMT", []*time.Location{
		time.FixedZone("Turkmenistan Time", 5*60*60),
	})
	retVal.Store("TRT", []*time.Location{
		time.FixedZone("Turkey Time", 3*60*60),
	})
	retVal.Store("TOT", []*time.Location{
		time.FixedZone("Tonga Time", 13*60*60),
	})
	retVal.Store("TST", []*time.Location{
		time.FixedZone("Taiwan Standard Time", 8*60*60),
	})
	retVal.Store("TVT", []*time.Location{
		time.FixedZone("Tuvalu Time", 12*60*60),
	})
	retVal.Store("ULAST", []*time.Location{
		time.FixedZone("Ulaanbaatar Summer Time", 9*60*60),
	})
	retVal.Store("ULAT", []*time.Location{
		time.FixedZone("Ulaanbaatar Standard Time", 8*60*60),
	})
	retVal.Store("UTC", []*time.Location{
		time.FixedZone("Coordinated Universal Time", 0),
	})
	retVal.Store("UYST", []*time.Location{
		time.FixedZone("Uruguay Summer Time", -2*60*60),
	})
	retVal.Store("UYT", []*time.Location{
		time.FixedZone("Uruguay Standard Time", -3*60*60),
	})
	retVal.Store("UZT", []*time.Location{
		time.FixedZone("Uzbekistan Time", 5*60*60),
	})
	retVal.Store("VET", []*time.Location{
		time.FixedZone("Venezuelan Standard Time", -4*60*60),
	})
	retVal.Store("VLAT", []*time.Location{
		time.FixedZone("Vladivostok Time", 10*60*60),
	})
	retVal.Store("VOLT", []*time.Location{
		time.FixedZone("Volgograd Time", 3*60*60),
	})
	retVal.Store("VOST", []*time.Location{
		time.FixedZone("Vostok Station Time", 6*60*60),
	})
	retVal.Store("VUT", []*time.Location{
		time.FixedZone("Vanuatu Time", 11*60*60),
	})
	retVal.Store("WAKT", []*time.Location{
		time.FixedZone("Wake Island Time", 12*60*60),
	})
	retVal.Store("WAST", []*time.Location{
		time.FixedZone("West Africa Summer Time", 2*60*60),
	})
	retVal.Store("WAT", []*time.Location{
		time.FixedZone("West Africa Time", 1*60*60),
	})
	retVal.Store("WEST", []*time.Location{
		time.FixedZone("Western European Summer Time", 1*60*60),
	})
	retVal.Store("WET", []*time.Location{
		time.FixedZone("Western European Time", 0),
	})
	retVal.Store("WIB", []*time.Location{
		time.FixedZone("Western Indonesian Time", 7*60*60),
	})
	retVal.Store("WIT", []*time.Location{
		time.FixedZone("Eastern Indonesian Time", 9*60*60),
	})
	retVal.Store("WITA", []*time.Location{
		time.FixedZone("Central Indonesia Time", 8*60*60),
	})
	retVal.Store("WGST", []*time.Location{
		time.FixedZone("West Greenland Summer Time", -2*60*60),
	})
	retVal.Store("WGT", []*time.Location{
		time.FixedZone("West Greenland Time", -3*60*60),
	})
	retVal.Store("WST", []*time.Location{
		time.FixedZone("Western Standard Time", 8*60*60),
	})
	retVal.Store("YAKT", []*time.Location{
		time.FixedZone("Yakutsk Time", 9*60*60),
	})
	retVal.Store("YEKT", []*time.Location{
		time.FixedZone("Yekaterinburg Time", 5*60*60),
	})

	return &retVal
}

func serveTime(timeAbbrevations *sync.Map, errorChannel chan<- Error) httprouter.Handle {
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

		abbrev, exists := timeAbbrevations.Load(location)
		if exists {
			val := abbrev.([]*time.Location)

			zones = append(zones, val...)
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

		securityHeaders(w)

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

	timeAbbreviations := getTimeAbbrevations()

	mux.GET("/time/:time", serveTime(timeAbbreviations, errorChannel))
	mux.GET("/time/:time/*rest", serveTime(timeAbbreviations, errorChannel))
	mux.GET("/time/", serveUsage(module, usage, errorChannel))

	usage.Store(module, []string{
		"/time/America/Chicago",
		"/time/EST",
		"/time/UTC?format=kitchen",
	})
}
