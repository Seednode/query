/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	EmptyColor = color.RGBA{0, 0, 0, 0x0}
)

var defaultColors = map[string]color.Color{
	"aliceblue":            color.RGBA{240, 248, 255, 0xff},
	"antiquewhite":         color.RGBA{250, 235, 215, 0xff},
	"aquamarine":           color.RGBA{127, 255, 212, 0xff},
	"aqua":                 color.RGBA{0, 255, 255, 0xff},
	"azure":                color.RGBA{240, 255, 255, 0xff},
	"beige":                color.RGBA{245, 245, 220, 0xff},
	"bisque":               color.RGBA{255, 228, 196, 0xff},
	"black":                color.RGBA{0, 0, 0, 0xff},
	"blanchedalmond":       color.RGBA{255, 235, 205, 0xff},
	"blue":                 color.RGBA{0, 0, 255, 0xff},
	"blueviolet":           color.RGBA{138, 43, 226, 0xff},
	"brown":                color.RGBA{165, 42, 42, 0xff},
	"burlywood":            color.RGBA{222, 184, 135, 0xff},
	"cadetblue":            color.RGBA{95, 158, 160, 0xff},
	"chartreuse":           color.RGBA{127, 255, 0, 0xff},
	"chocolate":            color.RGBA{210, 105, 30, 0xff},
	"coral":                color.RGBA{255, 127, 80, 0xff},
	"cornflowerblue":       color.RGBA{100, 149, 237, 0xff},
	"cornsilk":             color.RGBA{255, 248, 220, 0xff},
	"crimson":              color.RGBA{220, 20, 60, 0xff},
	"cyan":                 color.RGBA{0, 255, 255, 0xff},
	"darkblue":             color.RGBA{0, 0, 139, 0xff},
	"darkcyan":             color.RGBA{0, 139, 139, 0xff},
	"darkgoldenrod":        color.RGBA{184, 134, 11, 0xff},
	"darkgray":             color.RGBA{169, 169, 169, 0xff},
	"darkgreen":            color.RGBA{0, 100, 0, 0xff},
	"darkkhaki":            color.RGBA{189, 183, 107, 0xff},
	"darkmagenta":          color.RGBA{139, 0, 139, 0xff},
	"darkolivegreen":       color.RGBA{85, 107, 47, 0xff},
	"darkorange":           color.RGBA{255, 140, 0, 0xff},
	"darkorchid":           color.RGBA{153, 50, 204, 0xff},
	"darkred":              color.RGBA{139, 0, 0, 0xff},
	"darksalmon":           color.RGBA{233, 150, 122, 0xff},
	"darkseagreen":         color.RGBA{143, 188, 139, 0xff},
	"darkslateblue":        color.RGBA{72, 61, 139, 0xff},
	"darkslategray":        color.RGBA{47, 79, 79, 0xff},
	"darkturquoise":        color.RGBA{0, 206, 209, 0xff},
	"darkviolet":           color.RGBA{148, 0, 211, 0xff},
	"deeppink":             color.RGBA{255, 20, 147, 0xff},
	"deepskyblue":          color.RGBA{0, 191, 255, 0xff},
	"dimgray":              color.RGBA{105, 105, 105, 0xff},
	"dodgerblue":           color.RGBA{30, 144, 255, 0xff},
	"firebrick":            color.RGBA{178, 34, 34, 0xff},
	"floralwhite":          color.RGBA{255, 250, 240, 0xff},
	"forestgreen":          color.RGBA{34, 139, 34, 0xff},
	"fuchsia":              color.RGBA{255, 0, 255, 0xff},
	"gainsboro":            color.RGBA{220, 220, 220, 0xff},
	"ghostwhite":           color.RGBA{248, 248, 255, 0xff},
	"goldenrod":            color.RGBA{218, 165, 32, 0xff},
	"gold":                 color.RGBA{255, 215, 0, 0xff},
	"gray":                 color.RGBA{128, 128, 128, 0xff},
	"green":                color.RGBA{0, 128, 0, 0xff},
	"greenyellow":          color.RGBA{173, 255, 47, 0xff},
	"honeydew":             color.RGBA{240, 255, 240, 0xff},
	"hotpink":              color.RGBA{255, 105, 180, 0xff},
	"indianred":            color.RGBA{205, 92, 92, 0xff},
	"indigo":               color.RGBA{75, 0, 130, 0xff},
	"ivory":                color.RGBA{255, 255, 240, 0xff},
	"khaki":                color.RGBA{240, 230, 140, 0xff},
	"lavenderblush":        color.RGBA{255, 240, 245, 0xff},
	"lavender":             color.RGBA{230, 230, 250, 0xff},
	"lawngreen":            color.RGBA{124, 252, 0, 0xff},
	"lemonchiffon":         color.RGBA{255, 250, 205, 0xff},
	"lightblue":            color.RGBA{173, 216, 230, 0xff},
	"lightcoral":           color.RGBA{240, 128, 128, 0xff},
	"lightcyan":            color.RGBA{224, 255, 255, 0xff},
	"lightgoldenrodyellow": color.RGBA{250, 250, 210, 0xff},
	"lightgray":            color.RGBA{211, 211, 211, 0xff},
	"lightgreen":           color.RGBA{144, 238, 144, 0xff},
	"lightpink":            color.RGBA{255, 182, 193, 0xff},
	"lightsalmon":          color.RGBA{255, 160, 122, 0xff},
	"lightseagreen":        color.RGBA{32, 178, 170, 0xff},
	"lightskyblue":         color.RGBA{135, 206, 250, 0xff},
	"lightslategray":       color.RGBA{119, 136, 153, 0xff},
	"lightsteelblue":       color.RGBA{176, 196, 222, 0xff},
	"lightyellow":          color.RGBA{255, 255, 224, 0xff},
	"limegreen":            color.RGBA{50, 205, 50, 0xff},
	"lime":                 color.RGBA{0, 255, 0, 0xff},
	"linen":                color.RGBA{250, 240, 230, 0xff},
	"magenta":              color.RGBA{255, 0, 255, 0xff},
	"maroon":               color.RGBA{128, 0, 0, 0xff},
	"mediumaquamarine":     color.RGBA{102, 205, 170, 0xff},
	"mediumblue":           color.RGBA{0, 0, 205, 0xff},
	"mediumorchid":         color.RGBA{186, 85, 211, 0xff},
	"mediumpurple":         color.RGBA{147, 112, 219, 0xff},
	"mediumseagreen":       color.RGBA{60, 179, 113, 0xff},
	"mediumslateblue":      color.RGBA{123, 104, 238, 0xff},
	"mediumspringgreen":    color.RGBA{0, 250, 154, 0xff},
	"mediumturquoise":      color.RGBA{72, 209, 204, 0xff},
	"mediumvioletred":      color.RGBA{199, 21, 133, 0xff},
	"midnightblue":         color.RGBA{25, 25, 112, 0xff},
	"mintcream":            color.RGBA{245, 255, 250, 0xff},
	"mistyrose":            color.RGBA{255, 228, 225, 0xff},
	"moccasin":             color.RGBA{255, 228, 181, 0xff},
	"navajowhite":          color.RGBA{255, 222, 173, 0xff},
	"navy":                 color.RGBA{0, 0, 128, 0xff},
	"oldlace":              color.RGBA{253, 245, 230, 0xff},
	"olivedrab":            color.RGBA{107, 142, 35, 0xff},
	"olive":                color.RGBA{128, 128, 0, 0xff},
	"orangered":            color.RGBA{255, 69, 0, 0xff},
	"orange":               color.RGBA{255, 165, 0, 0xff},
	"orchid":               color.RGBA{218, 112, 214, 0xff},
	"palegoldenrod":        color.RGBA{238, 232, 170, 0xff},
	"palegreen":            color.RGBA{152, 251, 152, 0xff},
	"paleturquoise":        color.RGBA{175, 238, 238, 0xff},
	"palevioletred":        color.RGBA{219, 112, 147, 0xff},
	"papayawhip":           color.RGBA{255, 239, 213, 0xff},
	"peachpuff":            color.RGBA{255, 218, 185, 0xff},
	"peru":                 color.RGBA{205, 133, 63, 0xff},
	"pink":                 color.RGBA{255, 192, 203, 0xff},
	"plum":                 color.RGBA{221, 160, 221, 0xff},
	"powderblue":           color.RGBA{176, 224, 230, 0xff},
	"purple":               color.RGBA{128, 0, 128, 0xff},
	"rebeccapurple":        color.RGBA{102, 51, 153, 0xff},
	"red":                  color.RGBA{255, 0, 0, 0xff},
	"rosybrown":            color.RGBA{188, 143, 143, 0xff},
	"royalblue":            color.RGBA{65, 105, 225, 0xff},
	"saddlebrown":          color.RGBA{139, 69, 19, 0xff},
	"salmon":               color.RGBA{250, 128, 114, 0xff},
	"sandybrown":           color.RGBA{244, 164, 96, 0xff},
	"seagreen":             color.RGBA{46, 139, 87, 0xff},
	"seashell":             color.RGBA{255, 245, 238, 0xff},
	"sienna":               color.RGBA{160, 82, 45, 0xff},
	"silver":               color.RGBA{192, 192, 192, 0xff},
	"skyblue":              color.RGBA{135, 206, 235, 0xff},
	"slateblue":            color.RGBA{106, 90, 205, 0xff},
	"slategray":            color.RGBA{112, 128, 144, 0xff},
	"snow":                 color.RGBA{255, 250, 250, 0xff},
	"springgreen":          color.RGBA{0, 255, 127, 0xff},
	"steelblue":            color.RGBA{70, 130, 180, 0xff},
	"tan":                  color.RGBA{210, 180, 140, 0xff},
	"teal":                 color.RGBA{0, 128, 128, 0xff},
	"thistle":              color.RGBA{216, 191, 216, 0xff},
	"tomato":               color.RGBA{255, 99, 71, 0xff},
	"turquoise":            color.RGBA{64, 224, 208, 0xff},
	"violet":               color.RGBA{238, 130, 238, 0xff},
	"wheat":                color.RGBA{245, 222, 179, 0xff},
	"white":                color.RGBA{255, 255, 255, 0xff},
	"whitesmoke":           color.RGBA{245, 245, 245, 0xff},
	"yellowgreen":          color.RGBA{154, 205, 50, 0xff},
	"yellow":               color.RGBA{255, 255, 0, 0xff},
}

func isValidHex(s string) bool {
	dst := make([]byte, hex.DecodedLen(len(s)))

	if _, err := hex.Decode(dst, []byte(s)); err == nil {
		return true
	}

	return false
}

func getColor(requestedColor string) (color.Color, error) {
	r := chunks(requestedColor, 2)

	for _, val := range r {
		if !isValidHex(val) {
			return EmptyColor, nil
		}
	}

	red, err := strconv.Atoi(fmt.Sprintf("%x", "0x"+r[0]))
	if err != nil {
		return EmptyColor, err
	}

	green, err := strconv.Atoi(fmt.Sprintf("%x", "0x"+r[1]))
	if err != nil {
		return EmptyColor, err
	}

	blue, err := strconv.Atoi(fmt.Sprintf("%x", "0x"+r[2]))
	if err != nil {
		return EmptyColor, err
	}

	return color.RGBA{uint8(red), uint8(blue), uint8(green), 0xff}, nil
}

func drawImage(format string, errorChannel chan<- Error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startTime := time.Now()

		w.Header().Set("Content-Type", "text/plain")

		requested := p.ByName("color")

		c, found := defaultColors[requested]

		var colorToUse color.Color
		var err error

		switch {
		case found:
			colorToUse = c
		case len(requested) != 6:
			_, err := w.Write([]byte("Color codes must be six hex characters.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		default:
			colorToUse, err = getColor(requested)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			requested = "#" + requested
		}

		if colorToUse == EmptyColor {
			_, err := w.Write([]byte("Failed to parse color.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		dimensions := strings.SplitN(p.ByName("dimensions"), "x", 2)

		width, err := strconv.Atoi(dimensions[0])
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			_, err := w.Write([]byte("Failed to parse width.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		height, err := strconv.Atoi(dimensions[1])
		if err != nil {
			errorChannel <- Error{err, realIP(r, true), r.URL.Path}

			_, err := w.Write([]byte("Failed to parse height.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		switch {
		case width > maxImageWidth:
			_, err := w.Write([]byte(fmt.Sprintf("Image width must be no greater than %d", maxImageWidth)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		case height > maxImageHeight:
			_, err := w.Write([]byte(fmt.Sprintf("Image height must be no greater than %d", maxImageHeight)))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		img := image.NewRGBA(image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		})

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				img.Set(x, y, colorToUse)
			}
		}

		switch format {
		case "GIF":
			w.Header().Set("Content-Type", "image/gif")

			err := gif.Encode(w, img, nil)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				_, err := w.Write([]byte("Failed to encode GIF.\n"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			}
		case "JPEG":
			w.Header().Set("Content-Type", "image/jpeg")

			err := jpeg.Encode(w, img, nil)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				_, err := w.Write([]byte("Failed to encode JPEG.\n"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			}
		case "PNG":
			w.Header().Set("Content-Type", "image/png")

			err := png.Encode(w, img)
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}

				_, err := w.Write([]byte("Failed to encode PNG.\n"))
				if err != nil {
					errorChannel <- Error{err, realIP(r, true), r.URL.Path}
				}

				return
			}
		default:
			_, err := w.Write([]byte("Invalid image format requested.\n"))
			if err != nil {
				errorChannel <- Error{err, realIP(r, true), r.URL.Path}
			}

			return
		}

		if verbose {
			fmt.Printf("%s | %s requested a %dx%d %s of color %q\n",
				startTime.Format(timeFormats["RFC3339"]),
				realIP(r, true),
				width,
				height,
				format,
				requested)
		}
	}
}

func registerDraw(module string, mux *httprouter.Router, usage *sync.Map, errorChannel chan<- Error) []string {
	mux.GET("/draw/", serveUsage(module, usage))

	mux.GET("/draw/gif/:color/:dimensions", drawImage("GIF", errorChannel))
	mux.GET("/draw/gif/", serveUsage(module, usage))

	mux.GET("/draw/jpg/:color/:dimensions", drawImage("JPEG", errorChannel))
	mux.GET("/draw/jpg", serveUsage(module, usage))

	mux.GET("/draw/png/:color/:dimensions", drawImage("PNG", errorChannel))
	mux.GET("/draw/png/", serveUsage(module, usage))

	examples := make([]string, 3)
	examples[0] = "/draw/gif/beige/640x480"
	examples[1] = "/draw/jpg/white/320x240"
	examples[2] = "/draw/png/fafafa/1024x768"

	return examples
}
