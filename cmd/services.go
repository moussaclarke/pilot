package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var systemdServices = []string{"frankenphp", "mysql", "postgresql", "typesense-server"}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start all development services",
	Long:  "Start the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense and mailpit",
	Run: func(cmd *cobra.Command, args []string) {
		manageServices("start")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all development services",
	Long:  "Stop the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense and mailpit",
	Run: func(cmd *cobra.Command, args []string) {
		manageServices("stop")
	},
}

func manageServices(action string) {
	for _, service := range systemdServices {
		fmt.Printf("%s %s...\n", action, service)
		runServiceCommand("sudo", "systemctl", action, service)
	}

	fmt.Printf("%s mailpit...\n", action)
	runServiceCommand("brew", "services", action, "mailpit")
}

func runServiceCommand(name string, arg ...string) {
	err := exec.Command(name, arg...).Run()
	if err != nil {
		fmt.Printf("Failed to execute %s %v: %v\n", name, arg, err)
	}
}
