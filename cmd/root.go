package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "pilot",
	Short:   "Pilot is a lightweight local development environment manager",
	Long:    "Pilot is a command-line tool that helps you manage local development environments for web projects.\n\nIt is a very lightweight alternative to Laravel valet, orchestrating the specific stack required on my Ubuntu machine",
	Version: "0.4.3",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
