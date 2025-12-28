# GlowNet 🌟

GlowNet is a high-performance, viral-potential network monitoring tool built from scratch in Go. It serves as a modern, beautiful replacement for `netstat` and `nethogs`.

## Features

- **Real-time Monitoring**: Polls `/proc/net/tcp` and `/proc/net/tcp6` for active connections.
- **Process Mapping**: Automatically links socket inodes to process names and PIDs.
- **Global Bandwidth**: Displays real-time upload and download speeds in a stylish header.
- **Enrichment**:
    - **Geo-IP Lookup**: Identifies the country of remote hosts using a local MMDB.
    - **Reverse DNS**: Asynchronously resolves remote IP addresses to hostnames.
- **Modern TUI**: Built with `charmbracelet/bubbletea` and `charmbracelet/lipgloss`.
- **Zero C-Dependencies**: Pure Go implementation for maximum portability and performance.
- **Vim Bindings**: Native support for `j`/`k` navigation, `f` for filtering, and `s` for sorting.

## Installation

```bash
# Clone the repository
git clone https://github.com/knbr13/glow-net
cd glow-net

# Build the binary
go build -o glownet ./cmd/glownet

# Run (requires read access to /proc)
sudo ./glownet
```

## Why this exists?

Existing tools like `netstat` are often bloated or lack modern UI features. `nethogs` is great but depends on `libpcap` and C-wrappers, which can be a pain to distribute and may have performance overhead. GlowNet aims to provide a "single binary" experience with a beautiful interface and high-performance `/proc` parsing.

## Usage

- `j`/`k`: Navigate the connection table.
- `s`: Cycle sorting (by PID, Process Name, or Remote Host).
- `f`: Enter filter mode. Type to filter by process name or host. Press `Enter` or `Esc` to exit filter mode.
- `q`: Quit.

---

Built with ❤️ by Junie.
