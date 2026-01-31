package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var simple bool

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&simple, "simple", "s", false, "Return only 'up', 'down', or 'partial' instead of a formatted list of services")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check service status",
	Long:  "Check the active status of the currently managed development services.",
	Run: func(cmd *cobra.Command, args []string) {
		reqs := []string{"brew", "systemctl"}
		if !checkPreflight(reqs) {
			return
		}
		activeCount := 0
		totalServices := len(systemdServices) + 1 // +1 for mailpit

		// Check Systemd Services
		var systemdStatus strings.Builder
		for _, s := range systemdServices {
			statusCol := styleDim
			statusCheck := checkSystemdServiceStatus(s)
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
		mailpitActive := checkBrewServiceStatus("mailpit")

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
