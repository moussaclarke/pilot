package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "pilot",
	Short:   "Pilot is a lightweight local development environment manager",
	Long:    "🧑‍✈️ Pilot is a command-line tool that helps you manage local development environments for web projects.\nYou can think of it as a very lightweight alternative to Laravel Valet, orchestrating the specific stack required on my Ubuntu machine.\nCurrently this stack consists of frankenphp, mysql, postgresql, typesense, garage, and mailpit.\nCertificates are handled via mkcert and we KISS for domain resolution with /etc/hosts.",
	Version: "0.4.6",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
