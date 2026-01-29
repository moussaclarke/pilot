package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [domain]",
	Short: "Initialise a new site configuration",
	Long:  "Initialise a new site configuration.\nRun this command from your project root. It creates a .pilot directory, updates /etc/hosts, creates certs, creates a Caddyfile and imports it into the global Caddyfile. Finally it restarts the frankenphp server.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		pwd, _ := os.Getwd()

		// Validate public folder
		if _, err := os.Stat(filepath.Join(pwd, "public")); os.IsNotExist(err) {
			PrintWarning("The current directory does not contain a public folder. Continue? (y/n): ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return
			}
			PrintInfo("Ok, continuing. We'll set the entry point as 'public' but you'll need to change that manually if that's not what you want.")
		}

		// Check /etc/hosts
		hostsContent, _ := os.ReadFile("/etc/hosts")
		if strings.Contains(string(hostsContent), domain) {
			PrintError(fmt.Sprintf("Domain %s already exists in /etc/hosts", domain))
			PrintDim("Please remove it manually and try again.")
			return
		}

		// Check .pilot folder
		if _, err := os.Stat(".pilot"); !os.IsNotExist(err) {
			PrintError(".pilot folder already exists.")
			PrintDim("Please remove it manually and try again.")
			return
		}

		// Update /etc/hosts
		hostEntry := fmt.Sprintf("\n127.0.0.1 %s", domain)
		f, _ := os.OpenFile("/tmp/hosts_append", os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(hostEntry)
		f.Close()
		exec.Command("bash", "-c", "cat /tmp/hosts_append | sudo tee -a /etc/hosts").Run()
		PrintInfo(fmt.Sprintf("Added %s to /etc/hosts", domain))
		// Create certs
		os.Mkdir(".pilot", 0755)
		exec.Command("mkcert", "-cert-file", ".pilot/"+domain+".crt", "-key-file", ".pilot/"+domain+".key", domain).Run()
		PrintInfo("Created certs")

		// Create Caddyfile
		caddyContent := fmt.Sprintf("%s {\n  root * %s/public #change this to wherever your entry point is\n  php_server\n  tls %s/.pilot/%s.crt %s/.pilot/%s.key\n}",
			domain, pwd, pwd, domain, pwd, domain)
		os.WriteFile(".pilot/Caddyfile", []byte(caddyContent), 0644)
		PrintInfo("Created Caddyfile")

		// Update global Caddyfile
		importLine := fmt.Sprintf("\nimport %s/.pilot/Caddyfile", pwd)
		f2, _ := os.OpenFile("/tmp/caddy_append", os.O_CREATE|os.O_WRONLY, 0644)
		f2.WriteString(importLine)
		f2.Close()
		exec.Command("bash", "-c", "cat /tmp/caddy_append | sudo tee -a /etc/frankenphp/Caddyfile").Run()
		PrintInfo(fmt.Sprintf("Imported %s/.pilot/Caddyfile into /etc/frankenphp/Caddyfile", pwd))

		exec.Command("sudo", "systemctl", "restart", "frankenphp").Run()
		PrintInfo("Restarted frankenphp")
		PrintSuccess(fmt.Sprintf("Done! You can access your site at https://%s", domain))
		PrintDim("Your certs and Caddyfile are stored in the .pilot folder")
		PrintDim("If you need to change the entry point, you'll need to change the 'root' directive in the Caddyfile")
		PrintDim("Remember to restart frankenphp after making changes to the Caddyfile")
	},
}
