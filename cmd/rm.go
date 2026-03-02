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
	Long:    "Remove the configuration for a site.\nRun this command from the project root. It removes the .pilot directory or certs/Caddyfile, cleans /etc/hosts and /etc/frankenphp/Caddyfile, and restarts frankenphp.",
	Run: func(cmd *cobra.Command, args []string) {
		reqs := []string{"frankenphp"}
		if !checkPreflight(reqs) {
			return
		}

		pwd, _ := os.Getwd()

		// Detect setup type: proxy mode (certs folder) or PHP mode (.pilot folder)
		isProxyMode := false
		if _, err := os.Stat(".pilot"); !os.IsNotExist(err) {
			isProxyMode = false
		} else if _, err := os.Stat("certs"); !os.IsNotExist(err) {
			isProxyMode = true
		} else {
			PrintWarning("No .pilot directory or certs folder found.")
			return
		}

		// Determine domain from certificate files
		var domain string
		var certDir string
		if isProxyMode {
			certDir = "certs"
		} else {
			certDir = ".pilot"
		}

		files, _ := os.ReadDir(certDir)
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

		// Remove from hosts
		exec.Command("sudo", "sed", "-i", fmt.Sprintf("/127.0.0.1 %s/d", domain), "/etc/hosts").Run()
		PrintInfo(fmt.Sprintf("%s removed from /etc/hosts.", domain))

		// Remove from Caddyfile
		var importLine string
		if isProxyMode {
			importLine = fmt.Sprintf("import %s/Caddyfile", pwd)
		} else {
			importLine = fmt.Sprintf("import %s/.pilot/Caddyfile", pwd)
		}
		escapedLine := strings.ReplaceAll(importLine, "/", "\\/")
		exec.Command("sudo", "sed", "-i", fmt.Sprintf("/%s/d", escapedLine), globalCaddyPath).Run()
		PrintInfo(fmt.Sprintf("Import line removed from %s.", globalCaddyPath))

		// Remove directory/files
		if isProxyMode {
			os.RemoveAll("certs")
			os.Remove("Caddyfile")
			PrintInfo("certs folder and Caddyfile removed.")
		} else {
			os.RemoveAll(".pilot")
			PrintInfo(".pilot directory removed.")
		}

		exec.Command("sudo", "systemctl", "restart", "frankenphp").Run()
		PrintInfo("frankenphp restarted.")
		PrintSuccess(fmt.Sprintf("Done! %s has been removed.", domain))
	},
}
