package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

type SiteInfo struct {
	Domain      string
	Path        string
	PilotExists bool
	CaddyExists bool
	Certs       bool
}

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed sites",
	Long:  "List all sites currently found as imports in /etc/frankenphp/Caddyfile. If a pilot directory is present, also whether certs and Caddyfile exist.",
	Run: func(cmd *cobra.Command, args []string) {
		sites, err := getManagedSites()
		if err != nil {
			PrintError(fmt.Sprintf("Error reading configuration: %v", err))
			return
		}

		if len(sites) == 0 {
			PrintInfo("No sites currently managed by Pilot.")
			return
		}

		re := lipgloss.NewRenderer(os.Stdout)
		headerStyle := re.NewStyle().Foreground(highlight).Bold(true).Align(lipgloss.Center)
		borderStyle := re.NewStyle().Foreground(subtle)

		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(borderStyle).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == table.HeaderRow {
					return headerStyle
				}
				return lipgloss.NewStyle().Padding(0, 1)
			}).
			Headers("DOMAIN", "PILOT", "CADDY", "CERTS", "PATH")

		for _, site := range sites {
			pilotStatus := "Yes"
			certStatus := "OK"
			caddyStatus := "OK"

			if !site.PilotExists {
				pilotStatus = "No"
				certStatus = "?"
				caddyStatus = "?"
			} else {
				if !site.CaddyExists {
					caddyStatus = "MISSING"
				}
				if !site.Certs {
					certStatus = "MISSING"
				}
			}

			t.Row(
				site.Domain,
				pilotStatus,
				caddyStatus,
				certStatus,
				site.Path,
			)
		}
		fmt.Println(t.Render())
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
			configPath := strings.TrimSpace(strings.Replace(line, "import ", "", 1))

			var projectRoot string
			var pilotExists bool

			// Detect if this is a Pilot-managed site or a standard Caddyfile import
			if strings.Contains(configPath, "/.pilot/Caddyfile") {
				projectRoot = filepath.Dir(filepath.Dir(configPath))
				pilotExists = true
			} else {
				projectRoot = filepath.Dir(configPath)
				pilotExists = false
			}

			domain := getDomainFromLocalCaddy(configPath)

			certStatus := false
			caddyStatus := false

			// Check filesystem for actual presence of files
			if _, err := os.Stat(configPath); err == nil {
				caddyStatus = true
			}

			if pilotExists {
				pilotDir := filepath.Join(projectRoot, ".pilot")
				certPath := filepath.Join(pilotDir, domain+".crt")
				if _, err := os.Stat(certPath); err == nil {
					certStatus = true
				}
			}

			sites = append(sites, SiteInfo{
				Domain:      domain,
				Path:        projectRoot,
				PilotExists: pilotExists,
				CaddyExists: caddyStatus,
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
