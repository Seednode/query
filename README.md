## About

A bunch of useless tools, all available via convenient HTTP endpoints!

Feature requests, code criticism, bug reports, general chit-chat, and unrelated angst accepted at `query@seedno.de`.

Static binary builds available [here](https://cdn.seedno.de/builds/query).

x86_64 and ARM Docker images of latest version: `oci.seedno.de/seednode/query:latest`.

Dockerfile available [here](https://git.seedno.de/seednode/query/docker/Dockerfile).

## Currently available tools

### Dice roll
Roll a specified number of dice.

Optionally display individual roll results, as well as total, by appending `?verbose`

Examples:
- `/roll/5d20`
- `/roll/d6?verbose`

### DNS
Look up DNS records for a given host

Examples:
- `/dns/a/<host>`
- `/dns/aaaa/<host>`
- `/dns/mx/<host>`
- `/dns/ns/<host>`

### IP address
View your current public IP

Examples:
- `/ip/`

### Time
Look up the current time in a given timezone

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
  -b, --bind string   address to bind to (default "0.0.0.0")
  -h, --help          help for query
  -p, --port uint16   port to listen on (default 8080)
  -v, --verbose       log tool usage to stdout
  -V, --version       display version and exit
```

## Building the Docker container
From inside the `docker/` subdirectory, build the image using the following command:

`REGISTRY=<registry url> LATEST=yes TAG=alpine ./build.sh`