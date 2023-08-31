/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	Version string = "0.5.1"
)

var (
	bind    string
	port    uint16
	verbose bool
	version bool

	rootCmd = &cobra.Command{
		Use:   "query",
		Short: "Serves a variety of web-based utilities.",
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
	rootCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "log tool usage to stdout")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.SilenceErrors = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.SetVersionTemplate("query v{{.Version}}\n")
	rootCmd.Version = Version
}
