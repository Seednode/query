/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	ReleaseVersion string = "1.6.8"
)

var (
	all            bool
	bind           string
	exitOnError    bool
	maxDiceRolls   int
	maxDiceSides   int
	maxImageHeight int
	maxImageWidth  int
	ouiFile        string
	dns            bool
	dnsResolver    string
	hashing        bool
	httpStatus     bool
	ip             bool
	mac            bool
	qr             bool
	qrSize         int
	roll           bool
	timezones      bool
	port           uint16
	profile        bool
	verbose        bool
	version        bool

	requiredArgs = []string{
		"all",
		"dns",
		"hash",
		"http-status",
		"ip",
		"mac",
		"qr",
		"roll",
		"time",
	}

	rootCmd = &cobra.Command{
		Use:   "query",
		Short: "Serves a variety of web-based utilities.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case qrSize < 256 || qrSize > 2048:
				return ErrInvalidQRSize
			case maxDiceRolls < 1:
				return ErrInvalidMaxDiceCount
			case maxDiceSides < 1:
				return ErrInvalidMaxDiceSides
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := servePage(args)

			return err
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&all, "all", false, "enable all functionality")
	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "address to bind to")
	rootCmd.Flags().BoolVar(&dns, "dns", false, "enable DNS lookup functionality")
	rootCmd.Flags().StringVar(&dnsResolver, "dns-resolver", "", "custom DNS server IP and port to query (e.g. 8.8.8.8:53)")
	rootCmd.Flags().BoolVar(&exitOnError, "exit-on-error", false, "shut down webserver on error, instead of just printing the error")
	rootCmd.Flags().BoolVar(&hashing, "hash", false, "enable hashing functionality")
	rootCmd.Flags().BoolVar(&httpStatus, "http-status", false, "enable HTTP response status code functionality")
	rootCmd.Flags().BoolVar(&ip, "ip", false, "enable IP lookup functionality")
	rootCmd.Flags().BoolVar(&mac, "mac", false, "enable MAC lookup functionality")
	rootCmd.Flags().IntVar(&maxDiceRolls, "max-dice-rolls", 1024, "maximum number of dice per roll")
	rootCmd.Flags().IntVar(&maxDiceSides, "max-dice-sides", 1024, "maximum number of sides per die")
	rootCmd.Flags().IntVar(&maxImageHeight, "max-image-height", 1024, "maximum height of generated images")
	rootCmd.Flags().IntVar(&maxImageWidth, "max-image-width", 1024, "maximum width of generated images")
	rootCmd.Flags().StringVar(&ouiFile, "oui-file", "", "path to Wireshark manufacturer database file")
	rootCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().BoolVar(&profile, "profile", false, "register net/http/pprof handlers")
	rootCmd.Flags().BoolVar(&qr, "qr", false, "enable QR code generation functionality")
	rootCmd.Flags().IntVar(&qrSize, "qr-size", 256, "height/width of PNG-encoded QR codes (in pixels)")
	rootCmd.Flags().BoolVar(&roll, "roll", false, "enable dice rolling functionality")
	rootCmd.Flags().BoolVar(&timezones, "time", false, "enable time lookup functionality")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "log tool usage to stdout")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	rootCmd.Flags().SetInterspersed(true)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.MarkFlagsOneRequired(requiredArgs...)

	rootCmd.SilenceErrors = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.SetVersionTemplate("query v{{.Version}}\n")
	rootCmd.Version = ReleaseVersion
}
