package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Check system prerequisites",
	Long:  "Check for any missing system dependencies that Pilot requires and suggest how to resolve them.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(styleTitle.Render("Pilot System Diagnosis"))

		for _, req := range requirements {
			ok := checkRequirement(req)

			if ok {
				fmt.Printf("%s %-15s %s\n", styleSuccess.Render("✔"), req.Name, styleDim.Render("Found"))
			} else {
				fmt.Printf("%s %-15s %s\n", styleWarning.Render("✘"), req.Name, styleWarning.Render(req.Remedy))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
}
