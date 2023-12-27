/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	ReleaseVersion string = "0.17.2"
)

var (
	bind         string
	exitOnError  bool
	maxDiceRolls int
	maxDiceSides int
	ouiFile      string
	noDice       bool
	noDNS        bool
	noHash       bool
	noHttpStatus bool
	noIP         bool
	noOUI        bool
	noQR         bool
	noTime       bool
	port         uint16
	profile      bool
	qrSize       int
	verbose      bool
	version      bool

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
			err := ServePage(args)
			if err != nil {
				return err
			}

			return nil
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
	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "address to bind to")
	rootCmd.Flags().BoolVar(&exitOnError, "exit-on-error", false, "shut down webserver on error, instead of just printing the error")
	rootCmd.Flags().IntVar(&maxDiceRolls, "max-dice-rolls", 1024, "maximum number of dice per roll")
	rootCmd.Flags().IntVar(&maxDiceSides, "max-dice-sides", 1024, "maximum number of sides per die")
	rootCmd.Flags().BoolVar(&noDice, "no-dice", false, "disable dice rolling functionality")
	rootCmd.Flags().BoolVar(&noDNS, "no-dns", false, "disable dns lookup functionality")
	rootCmd.Flags().BoolVar(&noHash, "no-hash", false, "disable hashing functionality")
	rootCmd.Flags().BoolVar(&noHttpStatus, "no-http-status", false, "disable http response status code functionality")
	rootCmd.Flags().BoolVar(&noIP, "no-ip", false, "disable IP lookup functionality")
	rootCmd.Flags().BoolVar(&noOUI, "no-oui", false, "disable OUI lookup functionality")
	rootCmd.Flags().BoolVar(&noQR, "no-qr", false, "disable QR code generation functionality")
	rootCmd.Flags().BoolVar(&noTime, "no-time", false, "disable time lookup functionality")
	rootCmd.Flags().StringVar(&ouiFile, "oui-file", "", "path to wireshark manufacturer database file (https://www.wireshark.org/download/automated/data/manuf)")
	rootCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().BoolVar(&profile, "profile", false, "register net/http/pprof handlers")
	rootCmd.Flags().IntVar(&qrSize, "qr-size", 256, "height/width of PNG-encoded QR codes (in pixels)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "log tool usage to stdout")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	rootCmd.Flags().SetInterspersed(true)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.SilenceErrors = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.SetVersionTemplate("query v{{.Version}}\n")
	rootCmd.Version = ReleaseVersion
}
