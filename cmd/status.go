package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check service status",
	Long:  "Check the active status of the currently managed development services.\nSpecifically frankenphp, mysql, postgresql, typesense, mailpit and garage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Systemd Services:")
		for _, s := range systemdServices {
			out, _ := exec.Command("systemctl", "is-active", s).Output()
			fmt.Printf("  - %s: %s", s, out)
		}

		fmt.Println("Brew Services:")
		out, _ := exec.Command("brew", "services", "list").Output()
		if strings.Contains(string(out), "mailpit") && strings.Contains(string(out), "started") {
			fmt.Println("  - mailpit: active")
		} else {
			fmt.Println("  - mailpit: inactive")
		}
	},
}
