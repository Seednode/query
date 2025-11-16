## About

A bunch of useless tools, all available via convenient HTTP endpoints!

Feature requests, code criticism, bug reports, general chit-chat, and unrelated angst accepted at `query@seedno.de`.

Static binary builds available [here](https://cdn.seedno.de/builds/query).

x86_64 and ARM Docker images of latest version: `oci.seedno.de/seednode/query:latest`.

Dockerfile available [here](https://raw.githubusercontent.com/Seednode/query/master/docker/Dockerfile).

An example instance with all features enabled can be found [here](https://q.seedno.de/).

### Configuration
The following configuration methods are accepted, in order of highest to lowest priority:
- Command-line flags
- Environment variables

## Currently available tools

### Dice roll
Roll a specified number of dice.

Optionally display individual roll results, as well as total, by appending `?verbose`.

Examples:
- [/roll/5d20](https://q.seedno.de/roll/5d20)
- [/roll/d6?verbose](https://q.seedno.de/roll/d6?verbose)
- [/roll/4d6,5d8,d4?verbose](https://q.seedno.de/roll/4d6,5d8,d4?verbose)

### DNS
Look up DNS records for a given host.

An alternate DNS resolver can be specified via `--dns-resolver` (e.g. `--dns-resolver "1.1.1.1:53"`). If none is provided, the system default is used.

This uses Team Cymru's [IP to ASN mapping service](https://www.team-cymru.com/ip-asn-mapping), so please be considerate about traffic volume.

Examples:
- [/dns/a/google.com](https://q.seedno.de/dns/a/google.com)
- [/dns/aaaa/google.com](https://q.seedno.de/dns/aaaa/google.com)
- [/dns/host/google.com](https://q.seedno.de/dns/host/google.com)
- [/dns/mx/google.com](https://q.seedno.de/dns/mx/google.com)
- [/dns/ns/google.com](https://q.seedno.de/dns/ns/google.com)

### Hashing
Hash the provided string using the requested algorithm.

Examples:
- [/hash/md5/foo](https://q.seedno.de/hash/md5/foo)
- [/hash/sha1/foo](https://q.seedno.de/hash/sha1/foo)
- [/hash/sha224/foo](https://q.seedno.de/hash/sha224/foo)
- [/hash/sha256/foo](https://q.seedno.de/hash/sha256/foo)
- [/hash/sha384/foo](https://q.seedno.de/hash/sha384/foo)
- [/hash/sha512/foo](https://q.seedno.de/hash/sha512/foo)
- [/hash/sha512-224/foo](https://q.seedno.de/hash/sha512-224/foo)
- [/hash/sha512-256/foo](https://q.seedno.de/hash/sha512-256/foo)

In addition to providing the value to be hashed in the URL, you can submit it as the body of a GET request, so long as no value is provided in the URL.

For example, `curl -X GET https://q.seedno.de/hash/sha512-224/ -d "test"` will return the SHA512/224 hash for `test`.

### HTTP Status Codes
Receive the requested HTTP response status code.

Examples:
- [/http/status/200](https://q.seedno.de/http/status/200)
- [/http/status/404](https://q.seedno.de/http/status/404)
- [/http/status/500](https://q.seedno.de/http/status/500)

### IP address
View your current public IP.

Examples:
- [/ip/](https://q.seedno.de/ip/)

### MAC Lookup
Look up the vendor associated with any MAC address.

The [Wireshark manufacturer database](https://www.wireshark.org/download/automated/data/manuf) is embedded in the generated binary, but a local version can be used instead by providing the `--oui-file` argument.

Examples:
- [/mac/3c-7c-3f-1e-b9-a0](https://q.seedno.de/mac/3c-7c-3f-1e-b9-a0)
- [/mac/e0:00:84:aa:aa:bb](https://q.seedno.de/mac/e0:00:84:aa:aa:bb)
- [/mac/4C445BAABBCC](https://q.seedno.de/mac/4C445BAABBCC)

### QR Codes
Encode a string as a QR code (either a PNG or an ASCII string). Add `url` parm to prepend `https://` to string value.

Examples:
- [/qr/Test](https://q.seedno.de/qr/Test)
- [/qr/Test?string](https://q.seedno.de/qr/Test?string)
- [/qr/google.com?url](https://q.seedno.de/qr/google.com?url)

### Time
Look up the current time in a given timezone and format.

Values can optionally be formatted via the `?format=` query parameter by specifying any layout from the Go [time package](https://pkg.go.dev/time#pkg-constants).

Format values are case-insensitive.

Examples:
- [/time/America/Chicago](https://q.seedno.de/time/America/Chicago)
- [/time/EST](https://q.seedno.de/time/EST)
- [/time/UTC?format=kitchen](https://q.seedno.de/time/UTC?format=kitchen)

### Environment variables
Almost all options configurable via flags can also be configured via environment variables. 

The associated environment variable is the prefix `QUERY_` plus the flag name, with the following changes:
- Leading hyphens removed
- Converted to upper-case
- All internal hyphens converted to underscores

For example:
- `--dns-resolver 1.1.1.1:53` becomes `QUERY_DNS_RESOLVER=1.1.1.1:53`
- `--max-dice-rolls 256` becomes `QUERY_MAX_DICE_ROLLS=256`.

## Usage output
```
Serves a variety of web-based utilities.

Usage:
  query [flags]

Flags:
      --all                   enable all features
  -b, --bind string           address to bind to (default "0.0.0.0")
      --dns                   enable DNS lookup
      --dns-resolver string   custom DNS server IP and port to query (e.g. 8.8.8.8:53)
      --exit-on-error         shut down webserver on error, instead of just printing the error
      --hash                  enable hashing
  -h, --help                  help for query
      --http-status           enable HTTP response status codes
      --ip                    enable IP lookups
      --mac                   enable MAC lookups
      --max-dice-rolls int    maximum number of dice per roll (default 1024)
      --max-dice-sides int    maximum number of sides per die (default 1024)
      --oui-file string       path to Wireshark manufacturer database file
  -p, --port uint16           port to listen on (default 8080)
      --profile               register net/http/pprof handlers
      --qr                    enable QR code generation
      --qr-size int           height/width of PNG-encoded QR codes (in pixels) (default 256)
      --roll                  enable dice rolls
      --subnet                enable subnet calculator
      --time                  enable time lookup
      --tls-cert string       path to TLS certificate
      --tls-key string        path to TLS keyfile
  -v, --verbose               log tool usage to stdout
  -V, --version               display version and exit
```

## Building the Docker image
From inside the cloned repository, build the image using the following command:

`REGISTRY=<registry url> LATEST=yes TAG=alpine ./build-docker.sh`

## Run the container from local image repo

`docker run --rm -d -p 8080:8080/tcp localhost:5000/query:1.23.1-debug`