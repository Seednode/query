## About

A bunch of useless tools, all available via convenient HTTP endpoints!

Feature requests, code criticism, bug reports, general chit-chat, and unrelated angst accepted at `query@seedno.de`.

Static binary builds available [here](https://cdn.seedno.de/builds/query).

x86_64 and ARM Docker images of latest version: `oci.seedno.de/seednode/query:latest`.

Dockerfile available [here](https://git.seedno.de/seednode/query/raw/branch/master/docker/Dockerfile).

## Currently available tools

### Dice roll
Roll a specified number of dice.

Optionally display individual roll results, as well as total, by appending `?verbose`.

Examples:
- `/roll/5d20`
- `/roll/d6?verbose`

### DNS
Look up DNS records for a given host.

Examples:
- `/dns/a/<host>`
- `/dns/aaaa/<host>`
- `/dns/mx/<host>`
- `/dns/ns/<host>`

### HTTP Status Codes
Receive the requested HTTP response status code.

Examples:
- `/http/status/200`
- `/http/status/404`
- `/http/status/500`

### IP address
View your current public IP.

Examples:
- `/ip/`

### OUI Lookup
Look up the vendor associated with any MAC address.

The [Wireshark manufacturer database](https://www.wireshark.org/download/automated/data/manuf) is embedded in the generated binary, but a local version can be used instead by providing the `--oui-file` argument.

Examples:
- `/oui/00:00:08`
- `/oui/00-50-C2`
- `/oui/70b3d5`

### QR Codes
Encode a string as a QR code (either a PNG or an ASCII string).

Examples:
- `/qr/Test`
- `/qr/Test?string`

### Time
Look up the current time in a given timezone and format.

Values can optionally be formatted via the `?format=` query parameter by specifying any layout from the Go [time package](https://pkg.go.dev/time#pkg-constants).

Format values are case-insensitive.

Examples:
- `/time/America/Chicago`
- `/time/EST`
- `/time/UTC?format=kitchen`

## Usage output
```
Serves a variety of web-based utilities.

Usage:
  query [flags]

Flags:
  -b, --bind string          address to bind to (default "0.0.0.0")
      --exit-on-error        shut down webserver on error, instead of just printing the error
  -h, --help                 help for query
      --max-dice-rolls int   maximum number of dice per roll (default 1024)
      --max-dice-sides int   maximum number of sides per die (default 1024)
      --no-dice              disable dice rolling functionality
      --no-dns               disable dns lookup functionality
      --no-http-status       disable http response status code functionality
      --no-ip                disable IP lookup functionality
      --no-oui               disable OUI lookup functionality
      --no-qr                disable QR code generation functionality
      --no-time              disable time lookup functionality
      --oui-file string      path to wireshark manufacturer database file (https://www.wireshark.org/download/automated/data/manuf)
  -p, --port uint16          port to listen on (default 8080)
      --profile              register net/http/pprof handlers
      --qr-size int          height/width of PNG-encoded QR codes (in pixels) (default 256)
  -v, --verbose              log tool usage to stdout
  -V, --version              display version and exit
```

## Building the Docker container
From inside the `docker/` subdirectory, build the image using the following command:

`REGISTRY=<registry url> LATEST=yes TAG=alpine ./build.sh`
