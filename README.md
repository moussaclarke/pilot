# Pilot

🧑‍✈️ Pilot is a lightweight, opinionated development manager built for Ubuntu. It orchestrates a modern stack using **FrankenPHP** to provide a locally installed php environment similar to Laravel Valet.

> [!IMPORTANT]
> **Warning:** This tool is highly opinionated and strictly coupled to my specific local machine configuration. It assumes exact naming conventions for services and specific installation methods. If you run systemd and homebrew, then it might work on your machine but YMMV.

---

## Features

* **Service Management**: Start and stop the entire stack with single commands.
* **Automatic SSL**: Uses `mkcert` to generate local trusted certificates.
* **Host Management**: Updates `/etc/hosts` (requires `sudo`).
* **Caddy Integration**: Generates site-specific Caddyfiles and imports them into the global configuration at `/etc/frankenphp/Caddyfile`.

---

## Installation

### Service Dependencies

Pilot expects the following services to be present and named exactly as shown:

| Service | Management |
| --- | --- |
| **frankenphp** | `systemd` |
| **mysql** | `systemd` |
| **postgresql** | `systemd` |
| **typesense-server** | `systemd` |
| **mailpit** | `homebrew` |
| **garage** | `systemd` |

The installation method isn't important except for mailpit, which **must** be installed via homebrew since pilot relies on `brew services` to manage it.

### Prerequisites

* **Go** (to build this binary from source)
* **mkcert**
* **Homebrew**
* **Systemd**
* **mysql**
* **postgresql**
* **typesense**
* **frankenphp**
* **mailpit** (installed via homebrew)
* **garage**

On my machine I have these installed as follows as at the time of writing (YMMV):
- homebrew: go, mailpit, garage (latter with manual root service unit)
- apt: mysql, postgresql, mkcert
- deb/apt: typesense, frankenphp (latter with manual root service unit)

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

You can also optionally create a reverse proxy setup without a `.pilot` directory. This is useful if you want to set up a domain for a service running on some ad hoc port rather than from within a project repo.

```bash
pilot init example.test --proxy --port 1234
```

### Site Removal

Removes the `.pilot` configuration, cleans the hosts file entry, removes the Caddyfile import and restarts the web server. Will also detect a reverse proxy setup without a `.pilot` directory.

```bash
pilot rm
```

### List Sites

Displays the sites currently managed by Pilot, as well as any other sites found in the global Caddyfile.

```bash
pilot list
```

### Global Service Control

| Command | Description |
| --- | --- |
| `pilot up` | Starts all managed services. |
| `pilot down` | Stops all managed services. |
| `pilot status` | Displays the current running status of each service in the stack. |
| `pilot status --simple` | Displays the current running status of the stack in a more compact format for scripting. |

### Prerequisite Checks

Checks for any missing system dependencies and suggests how to resolve them.

```bash
pilot diagnose
```

Preflight checks are also run for most commands.

---

## Technical Details

### Core Architecture

Pilot transforms your local machine into a development server by managing three primary layers:

1. **Domain Resolution**: Maps custom domains to `127.0.0.1` via `/etc/hosts`.
2. **Request Handling**: Uses **FrankenPHP** as a web server and PHP runtime.
3. **Automatic SSL**: Uses `mkcert` to generate locally trusted certificates for every site, so https works.

### Comparison to Laravel Valet

While inspired by Valet, Pilot differs in several key areas:

* **Web Server**: Uses **FrankenPHP** instead of Nginx + PHP-FPM.
* **Service Manager**: Native **Systemd** integration for Linux instead of homebrew-managed macOS `launctl` services.
* **DNS**: Explicitly manages `/etc/hosts` rather than running a background `dnsmasq` proxy.
* **Configuration**: Stores site-specific settings and certs within a project-local `.pilot` folder instead of a global hidden directory.
* **Binary**: It's written in go and installed as a single self-contained binary rather than a php application.

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

Although it defaults to PHP, Pilot can manage any project. By editing the generated `.pilot/Caddyfile`, you can use it as a **reverse proxy** for Node.js, Go or Python applications, or as a **static file server** for frontend builds.

See also `pilot init --proxy` for reverse proxy setups where you're not wihin a project repo.

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
- [x] Diagnose command
- [ ] Site info command
- [x] List handle lab/reverse proxy/non php setups explicitly (i.e. no .pilot directory, everything in a visible top dir)
- [x] Init handle lab/reverse proxy/non php setups explicitly (i.e. no .pilot directory, everything in a visible top dir)
- [x] Rm handle lab/reverse proxy/non php setups explicitly (i.e. no .pilot directory, everything in a visible top dir)
- [ ] Better error handling
- [ ] Tests
- [x] Changelog
- [x] More documentation - including .pilot directory structure, overview of the stack

## Current Non-goals

- Installing dependencies (i.e. there will be no `valet install` equivalent)
- Other OSes or architectures. Or other people's machines in general. Might look at MacOS eventually though.
- Multi version support (i.e. there will be no `valet use` equivalent)
- DNSmasq (KISS `/etc/hosts` for now)
