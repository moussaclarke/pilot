package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var simple bool

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&simple, "simple", "s", false, "Return only 'up', 'down', or 'partial' instead of list")
}

func CheckSystemdServiceStatus(service string) bool {
	out, _ := exec.Command("systemctl", "is-active", service).Output()
	return strings.TrimSpace(string(out)) == "active"
}

func CheckBrewServiceStatus(service string) bool {
	var brewCmd *exec.Cmd
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		brewCmd = exec.Command("sudo", "-u", sudoUser, brewPath, "services", "info", service, "--json")
	} else {
		brewCmd = exec.Command(brewPath, "services", "info", service, "--json")
	}
	out, _ := brewCmd.Output()
	// unmarshal json output to check if service is active - take the first object in the array, and check for running : true
	var result []map[string]any
	json.Unmarshal(out, &result)

	return result[0]["running"].(bool)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check service status",
	Long:  "Check the active status of the currently managed development services.",
	Run: func(cmd *cobra.Command, args []string) {
		activeCount := 0
		totalServices := len(systemdServices) + 1 // +1 for mailpit

		// Check Systemd Services
		var systemdStatus strings.Builder
		for _, s := range systemdServices {
			statusCol := styleDim
			statusCheck := CheckSystemdServiceStatus(s)
			statusStr := "inactive"
			if statusCheck {
				statusStr = "active"
				activeCount++
				statusCol = styleSuccess
			}
			if !simple {
				systemdStatus.WriteString(fmt.Sprintf("  - %s: %s\n", styleTableCell.Width(20).Render(s), statusCol.Render(statusStr)))
			}
		}

		// Check Brew Services
		mailpitActive := CheckBrewServiceStatus("mailpit")

		if mailpitActive {
			activeCount++
		}

		// Handle Simple Output
		if simple {
			switch {
			case activeCount == totalServices:
				fmt.Println(styleSuccess.Render("up"))
			case activeCount > 0:
				fmt.Println(styleWarning.Render("partial"))
			default:
				fmt.Println(styleDim.Render("down"))
			}
			return
		}
		// Handle Detailed Output
		fmt.Println(styleHeading.Render("Systemd Services"))
		fmt.Print(systemdStatus.String())

		fmt.Println(styleHeading.Render("Brew Services"))
		if mailpitActive {
			fmt.Printf("  - %s: %s\n", styleTableCell.Width(20).Render("mailpit"), styleSuccess.Render("active"))
		} else {
			fmt.Printf("  - %s: %s\n", styleTableCell.Width(20).Render("mailpit"), styleDim.Render("inactive"))
		}
	},
}
