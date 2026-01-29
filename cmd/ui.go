package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colours
	subtle    = lipgloss.AdaptiveColor{Light: "#D9D9D9", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#FFA500", Dark: "#FFA500"}
	danger    = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}

	// Text Styles
	styleTitle   = lipgloss.NewStyle().Foreground(highlight).Bold(true).MarginBottom(1)
	styleHeading = lipgloss.NewStyle().Foreground(highlight).Bold(true)
	styleSuccess = lipgloss.NewStyle().Foreground(special)
	styleWarning = lipgloss.NewStyle().Foreground(warning)
	styleDanger  = lipgloss.NewStyle().Foreground(danger)
	styleDim     = lipgloss.NewStyle().Foreground(subtle)

	// Layout Styles
	styleTableHead = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true).BorderForeground(subtle).Bold(true).Padding(0, 1)
	styleTableCell = lipgloss.NewStyle().Padding(0, 1)
)

func PrintSuccess(msg string) {
	prefix := styleSuccess.Render("✔")
	fmt.Printf("%s %s\n", prefix, msg)
}

func PrintInfo(msg string) {
	prefix := styleHeading.Render("•")
	fmt.Printf("%s %s\n", prefix, msg)
}

func PrintWarning(msg string) {
	prefix := styleWarning.Bold(true).Render("!")
	fmt.Printf("%s %s\n", prefix, styleWarning.Render(msg))
}

func PrintError(msg string) {
	prefix := styleDanger.Bold(true).Render("!")
	fmt.Printf("%s %s\n", prefix, styleWarning.Render(msg))
}

func PrintDim(msg string) {
	fmt.Println(styleDim.Render(msg))
}
