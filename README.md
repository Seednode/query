## About

A bunch of useless tools, all available via convenient HTTP endpoints!

Feature requests, code criticism, bug reports, general chit-chat, and unrelated angst accepted at `query@seedno.de`.

Static binary builds available [here](https://cdn.seedno.de/builds/query).

x86_64 and ARM Docker images of latest version: `oci.seedno.de/seednode/query:latest`.

Dockerfile available [here](https://git.seedno.de/seednode/query/raw/branch/master/docker/Dockerfile).

An example instance with all features enabled can be found [here](https://q.seedno.de/).

## Currently available tools

### Dice roll
Roll a specified number of dice.

Optionally display individual roll results, as well as total, by appending `?verbose`.

Examples:
- `/roll/5d20`
- `/roll/d6?verbose`

### DNS
Look up DNS records for a given host.

An alternate DNS resolver can be specified via `--dns-resolver` (e.g. `--dns-resolver "1.1.1.1:53"`). If none is provided, the system default is used.

This uses Team Cymru's [IP to ASN mapping service](https://www.team-cymru.com/ip-asn-mapping), so please be considerate about traffic volume.

Examples:
- `/dns/a/google.com`
- `/dns/aaaa/google.com`
- `/dns/host/google.com`
- `/dns/mx/google.com`
- `/dns/ns/google.com`

### Draw
Outputs a rectangle with the specified width and height. The color can be specified as either a hex value (without the leading `#`) or any [HTML color name](https://www.w3schools.com/tags/ref_colornames.asp).

Examples:
- `/draw/gif/beige/640x480`
- `/draw/jpg/white/320x240`
- `/draw/png/fafafa/1024x768`

### Hashing
Hash the provided string using the requested algorithm.

Examples:
- `/hash/md5/foo`
- `/hash/sha1/foo`
- `/hash/sha224/foo`
- `/hash/sha256/foo`
- `/hash/sha384/foo`
- `/hash/sha512/foo`
- `/hash/sha512-224/foo`
- `/hash/sha512-256/foo`

In addition to providing the value to be hashed in the URL, you can submit it as the body of a GET request, so long as no value is provided in the URL.

For example, `curl -X GET https://q.seedno.de/hash/sha512-224/ -d "test"` will return the SHA512/224 hash for `test`.

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

### MAC Lookup
Look up the vendor associated with any MAC address.

The [Wireshark manufacturer database](https://www.wireshark.org/download/automated/data/manuf) is embedded in the generated binary, but a local version can be used instead by providing the `--oui-file` argument.

Examples:
- `/mac/3c-7c-3f-1e-b9-a0`
- `/mac/e0:00:84:aa:aa:bb`
- `/mac/4C445BAABBCC`

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
      --all                    enable all functionality
  -b, --bind string            address to bind to (default "0.0.0.0")
      --dns                    enable DNS lookup functionality
      --dns-resolver string    DNS server IP and port to query (e.g. 1.1.1.1:53, uses system default if empty)
      --draw                   enable drawing functionality
      --exit-on-error          shut down webserver on error, instead of just printing the error
      --hash                   enable hashing functionality
  -h, --help                   help for query
      --http-status            enable HTTP response status code functionality
      --ip                     enable IP lookup functionality
      --mac                    enable MAC lookup functionality
      --max-dice-rolls int     maximum number of dice per roll (default 1024)
      --max-dice-sides int     maximum number of sides per die (default 1024)
      --max-image-height int   maximum height of generated images (default 1024)
      --max-image-width int    maximum width of generated images (default 1024)
      --oui-file string        path to Wireshark manufacturer database file (https://www.wireshark.org/download/automated/data/manuf)
  -p, --port uint16            port to listen on (default 8080)
      --profile                register net/http/pprof handlers
      --qr                     enable QR code generation functionality
      --qr-size int            height/width of PNG-encoded QR codes (in pixels) (default 256)
      --roll                   enable dice rolling functionality
      --time                   enable time lookup functionality
  -v, --verbose                log tool usage to stdout
  -V, --version                display version and exit
```

## Building the Docker container
From inside the `docker/` subdirectory, build the image using the following command:

`REGISTRY=<registry url> LATEST=yes TAG=alpine ./build.sh`
