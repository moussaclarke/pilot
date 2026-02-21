package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const brewPath = "/home/linuxbrew/.linuxbrew/bin/brew"
const globalCaddyPath = "/etc/frankenphp/Caddyfile"
const hostsPath = "/etc/hosts"

type CheckType int

const (
	Binary CheckType = iota
	SystemdUnit
	BrewFormula
)

type Requirement struct {
	Name       string
	Type       CheckType
	Identifier string
	Remedy     string
}

var requirements = []Requirement{
	{Name: "Systemd", Type: Binary, Identifier: "systemctl", Remedy: "Your machine needs to be running systemd"},
	{Name: "Homebrew", Type: Binary, Identifier: "brew", Remedy: fmt.Sprintf("Ensure Homebrew is installed and accessible at %s", brewPath)},
	{Name: "FrankenPHP", Type: SystemdUnit, Identifier: "frankenphp", Remedy: "Install binary and ensure /etc/systemd/system/frankenphp.service exists"},
	{Name: "MySQL", Type: SystemdUnit, Identifier: "mysql", Remedy: "Install via 'sudo apt install mysql' and ensure /etc/systemd/system/mysql.service exists"},
	{Name: "PostgreSQL", Type: SystemdUnit, Identifier: "postgresql", Remedy: "Install via 'sudo apt install postgresql' and ensure /etc/systemd/system/postgresql.service exists"},
	{Name: "Typesense", Type: SystemdUnit, Identifier: "typesense-server", Remedy: "Visit https://typesense.org and install via apt/deb and ensure /etc/systemd/system/typesense-server.service exists"},
	{Name: "Garage", Type: SystemdUnit, Identifier: "garage", Remedy: "Ensure manual systemd unit is configured for Garage binary at /etc/systemd/system/garage.service"},
	{Name: "Mailpit", Type: BrewFormula, Identifier: "mailpit", Remedy: "Ensure Homebrew is available and install via 'brew install mailpit'"},
	{Name: "mkcert", Type: Binary, Identifier: "mkcert", Remedy: "Install via 'sudo apt install mkcert' and run 'mkcert -install'"},
}

func checkRequirement(req Requirement) bool {
	switch req.Type {
	case Binary:
		// for brew check the hard coded path
		if req.Identifier == "brew" {
			_, err := os.Stat(brewPath)
			return err == nil
		}
		_, err := exec.LookPath(req.Identifier)
		return err == nil
	case SystemdUnit:
		out, _ := runServiceCommand("systemctl", "list-unit-files", req.Identifier+".service")
		return strings.Contains(string(out), "enabled")
	case BrewFormula:
		out, _ := runServiceCommand("brew", "list")
		return strings.Contains(string(out), req.Identifier)
	default:
		return false
	}
}

func checkPreflight(reqs []string) bool {
	res := true
	for _, id := range reqs {
		req := Requirement{}
		for _, r := range requirements {
			if r.Identifier == id {
				req = r
				break
			}
		}
		// if req empty, it's not a valid requirement
		if req.Identifier == "" {
			PrintError(fmt.Sprintf("Unknown requirement: %s", id))
			return false
		}
		if !checkRequirement(req) {
			res = false
			PrintError(fmt.Sprintf("Missing requirement: %s", req.Name))
			PrintDim(req.Remedy)
		}
	}

	return res
}

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
