package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
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
		PrintInfo("Source Caddyfile: /etc/frankenphp/Caddyfile")
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

			domains := getDomainsFromLocalCaddy(configPath)

			for _, domain := range domains {

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
	}
	return sites, scanner.Err()
}

func getDomainsFromLocalCaddy(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return []string{"unknown"}
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "{") || strings.HasPrefix(line, "}") {
			continue
		}

		// Only process lines that appear to be site addresses (no indentation in original file)
		// Note: For simplicity, we assume the site block starts at the beginning of the line
		parts := strings.Fields(line)
		if len(parts) > 0 {
			// Check if the line ends a block or is a directive (very basic heuristic)
			firstToken := strings.TrimSuffix(parts[0], "{")
			if !isDirective(firstToken) {
				// Split by comma for multiple domains on one line: "a.com, b.com {"
				rawDomains := strings.SplitSeq(firstToken, ",")
				for d := range rawDomains {
					cleanD := strings.TrimSpace(d)
					if cleanD != "" {
						domains = append(domains, cleanD)
					}
				}
			}
		}
	}
	if len(domains) == 0 {
		return []string{"unknown"}
	}
	return domains
}

func isDirective(token string) bool {
	directives := []string{"root", "php_fastcgi", "file_server", "reverse_proxy", "log", "tls", "import", "health_port", "health_uri", "php_server"}
	return slices.Contains(directives, token)
}
