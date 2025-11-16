/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ReleaseVersion string = "1.23.1"
)

var (
	all          bool
	bind         string
	exitOnError  bool
	maxDiceRolls int
	maxDiceSides int
	ouiFile      string
	dns          bool
	dnsResolver  string
	hashing      bool
	httpStatus   bool
	ip           bool
	mac          bool
	qr           bool
	qrSize       int
	roll         bool
	subnet       bool
	timezones    bool
	tlsCert      string
	tlsKey       string
	port         uint16
	profile      bool
	whoami       bool
	verbose      bool
	version      bool

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
		"whoami",
	}
)

func main() {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Serves a variety of web-based utilities.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initializeConfig(cmd)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case tlsCert == "" && tlsKey != "" || tlsCert != "" && tlsKey == "":
				return errors.New("TLS certificate and keyfile must both be specified to enable HTTPS")
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
			return servePage()
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "enable all features")
	cmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "address to bind to")
	cmd.Flags().BoolVar(&dns, "dns", false, "enable DNS lookup")
	cmd.Flags().StringVar(&dnsResolver, "dns-resolver", "", "custom DNS server IP and port to query (e.g. 8.8.8.8:53)")
	cmd.Flags().BoolVar(&exitOnError, "exit-on-error", false, "shut down webserver on error, instead of just printing the error")
	cmd.Flags().BoolVar(&hashing, "hash", false, "enable hashing")
	cmd.Flags().BoolVar(&httpStatus, "http-status", false, "enable HTTP response status codes")
	cmd.Flags().BoolVar(&ip, "ip", false, "enable IP lookups")
	cmd.Flags().BoolVar(&mac, "mac", false, "enable MAC lookups")
	cmd.Flags().IntVar(&maxDiceRolls, "max-dice-rolls", 1024, "maximum number of dice per roll")
	cmd.Flags().IntVar(&maxDiceSides, "max-dice-sides", 1024, "maximum number of sides per die")
	cmd.Flags().StringVar(&ouiFile, "oui-file", "", "path to Wireshark manufacturer database file")
	cmd.Flags().Uint16VarP(&port, "port", "p", 8080, "port to listen on")
	cmd.Flags().BoolVar(&profile, "profile", false, "register net/http/pprof handlers")
	cmd.Flags().BoolVar(&qr, "qr", false, "enable QR code generation")
	cmd.Flags().IntVar(&qrSize, "qr-size", 256, "height/width of PNG-encoded QR codes (in pixels)")
	cmd.Flags().BoolVar(&roll, "roll", false, "enable dice rolls")
	cmd.Flags().BoolVar(&subnet, "subnet", false, "enable subnet calculator")
	cmd.Flags().BoolVar(&timezones, "time", false, "enable time lookup")
	cmd.Flags().StringVar(&tlsCert, "tls-cert", "", "path to TLS certificate")
	cmd.Flags().StringVar(&tlsKey, "tls-key", "", "path to TLS keyfile")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "log tool usage to stdout")
	cmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")
	cmd.Flags().BoolVar(&whoami, "whoami", false, "enable whoami endpoint")

	cmd.Flags().SetInterspersed(true)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	cmd.MarkFlagsOneRequired(requiredArgs...)

	cmd.SilenceErrors = true
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	cmd.SetVersionTemplate("query v{{.Version}}\n")
	cmd.Version = ReleaseVersion

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func initializeConfig(cmd *cobra.Command) {
	v := viper.New()

	v.SetEnvPrefix("query")

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.AutomaticEnv()

	bindFlags(cmd, v)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
