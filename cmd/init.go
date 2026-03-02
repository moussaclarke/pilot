package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&proxyMode, "proxy", "p", false, "Create a reverse proxy setup instead of PHP setup")
	initCmd.Flags().IntVarP(&proxyPort, "port", "P", 0, "Port for the reverse proxy (required for proxy mode)")
}

var (
	proxyMode bool
	proxyPort int
)

var initCmd = &cobra.Command{
	Use:   "init [domain]",
	Short: "Initialise a new site configuration",
	Long:  "Initialise a new site configuration.\nRun this command from your project root. It creates a .pilot directory, updates hosts file, creates certs, creates a Caddyfile and imports it into the global Caddyfile. Finally it restarts the frankenphp server.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reqs := []string{"mkcert", "systemctl", "frankenphp"}
		if !checkPreflight(reqs) {
			return
		}
		domain := args[0]
		pwd, _ := os.Getwd()

		// Determine if we're in proxy mode
		isProxyMode := proxyMode

		// If not explicitly set, check if public folder exists
		if !isProxyMode {
			if _, err := os.Stat(filepath.Join(pwd, "public")); os.IsNotExist(err) {
				// Ask if they want reverse proxy setup
				fmt.Print("The current directory does not contain a public folder. Would you like to create a reverse proxy setup instead? (y/n/a): ")
				var answer string
				fmt.Scanln(&answer)
				if strings.ToLower(answer) == "y" {
					isProxyMode = true
				} else if strings.ToLower(answer) == "a" {
					PrintInfo("Ok, aborting.")
					return
				} else {
					PrintInfo("Ok, continuing with PHP setup. We'll set the entry point as 'public' but you'll need to change that manually if that's not what you want.")
				}
			}
		}

		// Get port if in proxy mode
		port := proxyPort
		if isProxyMode && port == 0 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter the port for the reverse proxy (e.g., 3000, 8080): ")
			portStr, _ := reader.ReadString('\n')
			portStr = strings.TrimSpace(portStr)
			port, _ = strconv.Atoi(portStr)
			if port == 0 {
				PrintError("Invalid port number. Please try again.")
				return
			}
		}

		// Check /etc/hosts
		hostsContent, _ := os.ReadFile(hostsPath)
		if strings.Contains(string(hostsContent), domain) {
			PrintError(fmt.Sprintf("Domain %s already exists in %s", domain, hostsPath))
			PrintDim("Please remove it manually and try again.")
			return
		}

		// Check .pilot folder (only for non-proxy mode)
		if !isProxyMode {
			if _, err := os.Stat(".pilot"); !os.IsNotExist(err) {
				PrintError(".pilot folder already exists.")
				PrintDim("Please remove it manually and try again.")
				return
			}
		}

		// Check certs folder (only for proxy mode)
		if isProxyMode {
			if _, err := os.Stat("certs"); !os.IsNotExist(err) {
				PrintError("certs folder already exists.")
				PrintDim("Please remove it manually and try again.")
				return
			}
		}

		// Update /etc/hosts
		hostEntry := fmt.Sprintf("\n127.0.0.1 %s", domain)
		f, _ := os.OpenFile("/tmp/hosts_append", os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(hostEntry)
		f.Close()
		exec.Command("bash", "-c", fmt.Sprintf("cat /tmp/hosts_append | sudo tee -a %s", hostsPath)).Run()
		PrintInfo(fmt.Sprintf("Added %s to %s", domain, hostsPath))

		// Create directory structure
		if isProxyMode {
			os.Mkdir("certs", 0755)
		} else {
			os.Mkdir(".pilot", 0755)
		}

		// Create certs
		if isProxyMode {
			exec.Command("mkcert", "-cert-file", "certs/"+domain+".crt", "-key-file", "certs/"+domain+".key", domain).Run()
		} else {
			exec.Command("mkcert", "-cert-file", ".pilot/"+domain+".crt", "-key-file", ".pilot/"+domain+".key", domain).Run()
		}
		PrintInfo("Created certs")

		// Create Caddyfile
		var caddyContent string
		if isProxyMode {
			caddyContent = fmt.Sprintf("%s {\n reverse_proxy 127.0.0.1:%d\n  tls %s/certs/%s.crt %s/certs/%s.key\n}",
				domain, port, pwd, domain, pwd, domain)
		} else {
			caddyContent = fmt.Sprintf("%s {\n  root * %s/public #change this to wherever your entry point is\n  php_server\n  tls %s/.pilot/%s.crt %s/.pilot/%s.key\n}",
				domain, pwd, pwd, domain, pwd, domain)
		}

		// Write Caddyfile to appropriate location
		if isProxyMode {
			os.WriteFile("Caddyfile", []byte(caddyContent), 0644)
		} else {
			os.WriteFile(".pilot/Caddyfile", []byte(caddyContent), 0644)
		}
		PrintInfo("Created Caddyfile")

		// Update global Caddyfile
		var importLine string
		if isProxyMode {
			importLine = fmt.Sprintf("\nimport %s/Caddyfile", pwd)
		} else {
			importLine = fmt.Sprintf("\nimport %s/.pilot/Caddyfile", pwd)
		}
		f2, _ := os.OpenFile("/tmp/caddy_append", os.O_CREATE|os.O_WRONLY, 0644)
		f2.WriteString(importLine)
		f2.Close()
		exec.Command("bash", "-c", fmt.Sprintf("cat /tmp/caddy_append | sudo tee -a %s", globalCaddyPath)).Run()
		PrintInfo(fmt.Sprintf("Imported Caddyfile into %s", globalCaddyPath))

		exec.Command("sudo", "systemctl", "restart", "frankenphp").Run()
		PrintInfo("Restarted frankenphp")
		PrintSuccess(fmt.Sprintf("Done! You can access your site at https://%s", domain))
		if isProxyMode {
			PrintDim("Your certs and Caddyfile are stored in the project root")
			PrintDim("If you need to change the reverse proxy target, edit the 'reverse_proxy' directive in the Caddyfile")
		} else {
			PrintDim("Your certs and Caddyfile are stored in the .pilot folder")
			PrintDim("If you need to change the entry point, you'll need to change the 'root' directive in the Caddyfile")
		}
		PrintDim("Remember to restart frankenphp after making changes to the Caddyfile")
	},
}
