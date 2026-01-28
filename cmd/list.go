package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type SiteInfo struct {
	Domain      string
	Path        string
	PilotExists bool
	Certs       bool
}

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed sites",
	Run: func(cmd *cobra.Command, args []string) {
		sites, err := getManagedSites()
		if err != nil {
			fmt.Printf("Error reading configuration: %v\n", err)
			return
		}

		if len(sites) == 0 {
			fmt.Println("No sites currently managed by Pilot.")
			return
		}

		fmt.Printf("%-25s %-10s %-10s %-30s\n", "DOMAIN", "PILOT", "CERTS", "PATH")
		fmt.Println(strings.Repeat("-", 85))
		for _, site := range sites {
			pilotStatus := "Yes"
			certStatus := "OK"

			if !site.PilotExists {
				pilotStatus = "No"
				certStatus = "?"
			} else if !site.Certs {
				certStatus = "MISSING"
			}

			fmt.Printf("%-25s %-10s %-10s %-30s\n", site.Domain, pilotStatus, certStatus, site.Path)
		}
	},
}

func getManagedSites() ([]SiteInfo, error) {
	file, err := os.Open("/etc/frankenphp/Caddyfile")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sites []SiteInfo
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "import ") {
			// Extract path: import /home/user/project/.pilot/Caddyfile
			configPath := strings.TrimSpace(strings.Replace(line, "import ", "", 1))
			projectRoot := filepath.Dir(filepath.Dir(configPath))
			pilotDir := filepath.Join(projectRoot, ".pilot")

			// Check if .pilot directory exists
			_, pilotErr := os.Stat(pilotDir)
			pilotExists := pilotErr == nil

			// Extract domain from local config
			domain := getDomainFromLocalCaddy(configPath)

			// Check for certs specifically if pilot exists
			certStatus := false
			if pilotExists {
				certPath := filepath.Join(pilotDir, domain+".crt")
				_, certErr := os.Stat(certPath)
				certStatus = certErr == nil
			}

			sites = append(sites, SiteInfo{
				Domain:      domain,
				Path:        projectRoot,
				PilotExists: pilotExists,
				Certs:       certStatus,
			})
		}
	}
	return sites, scanner.Err()
}

func getDomainFromLocalCaddy(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) > 0 {
			// Remove trailing brace if present e.g. "domain.test {"
			return strings.TrimSuffix(parts[0], "{")
		}
	}
	return "unknown"
}

