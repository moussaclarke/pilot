package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var simple bool

func init() {
	rootCmd.AddCommand(statusCmd)
	// Register the local flag
	statusCmd.Flags().BoolVarP(&simple, "simple", "s", false, "Return only 'up' or 'down' instead of list")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check service status",
	Long:  "Check the active status of the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense, mailpit and garage",
	Run: func(cmd *cobra.Command, args []string) {
		allUp := true

		// Check Systemd Services
		var systemdStatus strings.Builder
		for _, s := range systemdServices {
			out, _ := exec.Command("systemctl", "is-active", s).Output()
			statusStr := string(out)
			if !strings.Contains(statusStr, "active") {
				allUp = false
			}
			if !simple {
				systemdStatus.WriteString(fmt.Sprintf("  - %s: %s", s, statusStr))
			}
		}

		// Check Brew Services
		var brewCmd *exec.Cmd
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			brewCmd = exec.Command("sudo", "-u", sudoUser, brewPath, "services", "list")
		} else {
			brewCmd = exec.Command(brewPath, "services", "list")
		}

		brewOut, err := brewCmd.Output()
		mailpitActive := err == nil && strings.Contains(string(brewOut), "mailpit") && strings.Contains(string(brewOut), "started")

		if !mailpitActive {
			allUp = false
		}

		// Conditional Output
		if simple {
			if allUp {
				fmt.Println("up")
			} else {
				fmt.Println("down")
			}
			return
		}

		// Standard Detailed Output
		fmt.Println("Systemd Services:")
		fmt.Print(systemdStatus.String())
		fmt.Println("Brew Services:")
		if err != nil {
			fmt.Println("  - mailpit: error checking status")
		} else if mailpitActive {
			fmt.Println("  - mailpit: active")
		} else {
			fmt.Println("  - mailpit: inactive")
		}
	},
}

