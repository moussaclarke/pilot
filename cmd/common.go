package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

const brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"

func runServiceCommand(name string, arg ...string) ([]byte, error) {
	finalName := name
	finalArgs := arg

	if name == "brew" {
		finalName = brewPath

		// If running via sudo (e.g. from a system script), switch to the user
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			finalName = "sudo"
			finalArgs = append([]string{"-u", sudoUser, brewPath}, arg...)
		}
	}

	cmd := exec.Command(finalName, finalArgs...)

	// Inherit environment to ensure Homebrew variables are present
	cmd.Env = os.Environ()

	out, err := cmd.Output()

	return out, err
}

func checkSystemdServiceStatus(service string) bool {
	out, _ := runServiceCommand("systemctl", "is-active", service)
	return strings.TrimSpace(string(out)) == "active"
}

func checkBrewServiceStatus(service string) bool {
	out, _ := runServiceCommand(brewPath, "services", "info", service, "--json")
	// unmarshal json output to check if service is active - take the first object in the array, and check for running : true
	var result []map[string]any
	json.Unmarshal(out, &result)

	return result[0]["running"].(bool)
}
