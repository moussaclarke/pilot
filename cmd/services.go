package cmd

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"

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
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all development services",
	Long:  "Stop the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense, mailpit and garage",
	Run: func(cmd *cobra.Command, args []string) {
		manageServices("stop")
	},
}

func manageServices(action string) {
	for _, service := range systemdServices {
		PrintInfo(fmt.Sprintf("%s %s...", cases.Title(language.English).String(service), action))
		runServiceCommand("sudo", "systemctl", action, service)
	}

	PrintInfo(fmt.Sprintf("Mailpit %s...", action))
	runServiceCommand("brew", "services", action, "mailpit")
}

func runServiceCommand(name string, arg ...string) {
	finalName := name
	finalArgs := arg

	if name == "brew" {
		finalName = brewPath

		// If running via sudo (e.g. from a system script), switch to the user
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			finalName = "sudo"
			finalArgs = append([]string{"-u", sudoUser, brewPath}, arg...)
		}
	}

	cmd := exec.Command(finalName, finalArgs...)

	// Inherit environment to ensure Homebrew variables are present
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		PrintError(fmt.Sprintf("Failed to execute %s: %v", name, err))
	}
}
