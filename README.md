# Pilot

Pilot is a command-line tool written in Go designed to manage local development environments. It serves as a very lightweight alternative to Laravel Valet, orchestrating the specific stack required on my Ubuntu machine.

> [!IMPORTANT]
> **Warning:** This tool is highly opinionated and strictly coupled to my specific local machine configuration. It assumes exact naming conventions for services and specific installation methods.

---

## Service Dependencies

Pilot expects the following services to be present and named exactly as shown:

| Service | Installation Method | Management |
| --- | --- | --- |
| **frankenphp** | Manual binary download from GitHub | `systemctl` (manual service unit) |
| **mysql** | `apt` | `systemctl` |
| **postgresql** | `apt` | `systemctl` |
| **typesense-server** | `deb` / `apt` | `systemctl` |
| **mailpit** | `homebrew` | `brew services` |
| **garage** | `homebrew` | `systemctl` (manual service unit) |

---

## Features

* **Service Management**: Start and stop the entire stack with single commands.
* **Automatic SSL**: Uses `mkcert` to generate local trusted certificates.
* **Host Management**: Updates `/etc/hosts` (requires `sudo`).
* **Caddy Integration**: Generates site-specific Caddyfiles and imports them into the global configuration at `/etc/frankenphp/Caddyfile`.

---

## Installation

### Prerequisites

* **Go** (to build this binary)
* **mkcert**
* **Homebrew**
* **Systemd**
* **mysql** (installed via apt)
* **postgresql** (installed via apt)
* **typesense** (installed via deb/apt)
* **frankenphp** (manually installed via d/l and with an appropriate service unit)
* **mailpit** (installed via homebrew)
* **garage** (installed via homebrew but with a manual service unit)

### Build and install

```bash
go build -o pilot main.go
sudo mv pilot /usr/local/bin/
```

or

```bash
./install.sh
```

---

## Usage

### Site Initialisation

Run this command from your project root. It creates a `.pilot` directory, generates certificates, updates your hosts file, configures Caddy, and restarts the web server.

```bash
pilot init example.test
```

You can edit the Caddyfile manually if you need any kind of custom configuration.

### Site Removal

Removes the `.pilot` configuration, cleans the hosts file entry, removes the Caddyfile import and restarts the web server.

```bash
pilot rm
```

### List Sites

```bash
pilot list
```

Displays the sites currently managed by Pilot, as well as any other sites found in the global Caddyfile.

### Global Service Control

| Command | Description |
| --- | --- |
| `pilot up` | Starts all managed services. |
| `pilot down` | Stops all managed services. |
| `pilot status` | Displays the current running status of the stack. |
| `pilot status --simple` | Displays the current running status of the stack in a more compact format for scripting. |

---

## Technical Details

### Configuration Paths

* **Site Caddyfile**: `./.pilot/Caddyfile`
* **Site Certificates**: `./.pilot/certs/*`
* **Global Caddyfile**: `/etc/frankenphp/Caddyfile`
* **Hosts File**: `/etc/hosts`

## TODO

- [ ] Prettier output
- [ ] Automated Prerequisite check including installation or installation instructions
- [ ] Check service status before up/down
- [x] Site list command
- [ ] Site info command
- [ ] Better error handling
- [ ] Tests
- [ ] Changelog
- [ ] More documentation - including .pilot directory structure, overview of the stack

## Non-goals

- Other OSes or architectures. Or other people's machines in general.
- Multi version support (e.g. there will be no `valet use` equivalent)
- DNSmasq (KISS `/etc/hosts` for now)
