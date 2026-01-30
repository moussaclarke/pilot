# Pilot

🧑‍✈️ Pilot is a lightweight, opinionated development manager built for Ubuntu. It orchestrates a modern stack using **FrankenPHP** to provide a seamless local environment similar to Laravel Valet, but tailored specifically for Linux systemd and Homebrew services.

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

### Prerequisite Checks

```bash
pilot diagnose
```

Checks for any missing system dependencies and suggests how to resolve them. Preflight checks are also run for most commands.

---

## Technical Details

### Core Architecture

Pilot transforms your local machine into a development server by managing three primary layers:

1. **Domain Resolution**: Maps custom `.test` domains to `127.0.0.1` via `/etc/hosts`.
2. **Request Handling**: Uses **FrankenPHP** as a web server and PHP runtime.
3. **Automatic SSL**: Uses `mkcert` to generate locally trusted certificates for every site, ensuring a full HTTPS development experience.

### Comparison to Laravel Valet

While inspired by Valet, Pilot differs in several key areas:

* **Web Server**: Uses **FrankenPHP** instead of Nginx + PHP-FPM.
* **Service Manager**: Native **Systemd** integration for Linux instead of macOS `launchd`.
* **DNS**: Explicitly manages `/etc/hosts` for transparency rather than running a background `dnsmasq` proxy.
* **Configuration**: Stores site-specific settings and certs within a project-local `.pilot` folder instead of a global hidden directory.

---

### Project Conventions

#### The `.pilot` Directory

Upon running `pilot init`, a `.pilot` folder is created in your project root.

> [!IMPORTANT]
> You should add `/.pilot/` to your global or project-specific `.gitignore`.

This directory contains:

* **SSL Certificates**: Site-specific `.crt` and `.key` files.
* **Caddyfile**: The local routing rules for the project.

#### Non-PHP Projects

Although optimised for PHP, Pilot can manage any project. By editing the generated `.pilot/Caddyfile`, you can use it as a **reverse proxy** for Node.js, Go, or Python applications, or as a **static file server** for frontend builds.

---

### Configuration Paths

* **Site Caddyfile**: `./.pilot/Caddyfile`
* **Site Certificates**: `./.pilot/certs/*`
* **Global Caddyfile**: `/etc/frankenphp/Caddyfile`
* **Hosts File**: `/etc/hosts`

## TODO

- [x] Prettier output
- [x] Automated Prerequisite check including installation or installation instructions
- [x] Site init command
- [x] Site rm command
- [ ] Extend prerequisite checks to also check for global conf files /etc/hosts and /etc/frankenphp/Caddyfile
- [ ] Check service status before up/down
- [x] Site list command
- [ ] Site info command
- [ ] Better error handling
- [ ] Tests
- [ ] Changelog
- [x] More documentation - including .pilot directory structure, overview of the stack

## Non-goals

- Other OSes or architectures. Or other people's machines in general.
- Multi version support (e.g. there will be no `valet use` equivalent)
- DNSmasq (KISS `/etc/hosts` for now)
