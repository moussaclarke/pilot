package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "delete", "del"},
	Short:   "Remove site configuration",
	Long:    "Remove the configuration for a site.\nRun this command from the project root. It removes the .pilot directory, cleans /etc/hosts and /etc/frankenphp/Caddyfile, and restarts frankenphp.",
	Run: func(cmd *cobra.Command, args []string) {
		reqs := []string{"frankenphp"}
		if !checkPreflight(reqs) {
			return
		}
		if _, err := os.Stat(".pilot"); os.IsNotExist(err) {
			PrintWarning("No .pilot directory found.")
			return
		}

		files, _ := os.ReadDir(".pilot")
		var domain string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".crt") {
				domain = strings.TrimSuffix(file.Name(), ".crt")
				break
			}
		}

		if domain == "" {
			PrintError("Could not determine domain.")
			return
		}

		PrintWarning(fmt.Sprintf("Remove configuration for %s? (y/n): ", domain))
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) != "y" {
			return
		}

		pwd, _ := os.Getwd()

		// Remove from hosts
		exec.Command("sudo", "sed", "-i", fmt.Sprintf("/127.0.0.1 %s/d", domain), "/etc/hosts").Run()
		PrintInfo(fmt.Sprintf("%s removed from /etc/hosts.", domain))

		// Remove from Caddyfile
		importLine := fmt.Sprintf("import %s/.pilot/Caddyfile", pwd)
		escapedLine := strings.ReplaceAll(importLine, "/", "\\/")
		exec.Command("sudo", "sed", "-i", fmt.Sprintf("/%s/d", escapedLine), globalCaddyPath).Run()
		PrintInfo(fmt.Sprintf("%s/.pilot/Caddyfile removed from %s.", pwd, globalCaddyPath))
		os.RemoveAll(".pilot")
		PrintInfo(".pilot directory removed.")
		exec.Command("sudo", "systemctl", "restart", "frankenphp").Run()
		PrintInfo("frankenphp restarted.")
		PrintSuccess(fmt.Sprintf("Done! %s has been removed.", domain))
	},
}
