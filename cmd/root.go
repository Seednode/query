/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ReleaseVersion string = "1.20.1"
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
	port         uint16
	profile      bool
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
	}
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "query",
		Short: "Serves a variety of web-based utilities.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initializeConfig(cmd)
		},
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
			return servePage()
		},
	}

	rootCmd.Flags().BoolVar(&all, "all", false, "enable all features")
	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0", "address to bind to")
	rootCmd.Flags().BoolVar(&dns, "dns", false, "enable DNS lookup")
	rootCmd.Flags().StringVar(&dnsResolver, "dns-resolver", "", "custom DNS server IP and port to query (e.g. 8.8.8.8:53)")
	rootCmd.Flags().BoolVar(&exitOnError, "exit-on-error", false, "shut down webserver on error, instead of just printing the error")
	rootCmd.Flags().BoolVar(&hashing, "hash", false, "enable hashing")
	rootCmd.Flags().BoolVar(&httpStatus, "http-status", false, "enable HTTP response status codes")
	rootCmd.Flags().BoolVar(&ip, "ip", false, "enable IP lookups")
	rootCmd.Flags().BoolVar(&mac, "mac", false, "enable MAC lookups")
	rootCmd.Flags().IntVar(&maxDiceRolls, "max-dice-rolls", 1024, "maximum number of dice per roll")
	rootCmd.Flags().IntVar(&maxDiceSides, "max-dice-sides", 1024, "maximum number of sides per die")
	rootCmd.Flags().StringVar(&ouiFile, "oui-file", "", "path to Wireshark manufacturer database file")
	rootCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().BoolVar(&profile, "profile", false, "register net/http/pprof handlers")
	rootCmd.Flags().BoolVar(&qr, "qr", false, "enable QR code generation")
	rootCmd.Flags().IntVar(&qrSize, "qr-size", 256, "height/width of PNG-encoded QR codes (in pixels)")
	rootCmd.Flags().BoolVar(&roll, "roll", false, "enable dice rolls")
	rootCmd.Flags().BoolVar(&subnet, "subnet", false, "enable subnet calculator")
	rootCmd.Flags().BoolVar(&timezones, "time", false, "enable time lookup")
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

	return rootCmd
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
