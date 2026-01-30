package cmd

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"
)

var systemdServices = []string{"frankenphp", "mysql", "postgresql", "typesense-server", "garage"}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start all development services",
	Long:  "Start the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense, mailpit and garage",
	Run: func(cmd *cobra.Command, args []string) {
		manageServices("start")
		PrintSuccess("All development services started.")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all development services",
	Long:  "Stop the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense, mailpit and garage",
	Run: func(cmd *cobra.Command, args []string) {
		manageServices("stop")
		PrintSuccess("All development services stopped.")
	},
}

func manageServices(action string) {
	for _, service := range systemdServices {
		PrintInfo(fmt.Sprintf("%s %s...", cases.Title(language.English).String(service), action))
		_, err := runServiceCommand("sudo", "systemctl", action, service)
		if err != nil {
			PrintError(fmt.Sprintf("Error %s %s: %v", cases.Title(language.English).String(service), action, err))
		}
	}

	PrintInfo(fmt.Sprintf("Mailpit %s...", action))
	_, err := runServiceCommand("brew", "services", action, "mailpit")
	if err != nil {
		PrintError(fmt.Sprintf("Error Mailpit %s: %v", action, err))
	}
}
