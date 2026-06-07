<p align="center">
  <img src="logo.svg" alt="jumie logo" width="200"/>
</p>

# jumie

A smart Unix terminal assistant powered by AI.
It runs as a lightweight background daemon and a quick CLI client, translating natural language queries into safe, context-aware command pipelines.
This tool is designed to make command line navigation accessible and secure, giving you a full, human-friendly explanation of actions before executing them.

## Key Features

* **100% Local & Private:** Uses an isolated, embedded Ollama instance running `gemma4:e2b` (with MLX acceleration on macOS). It automatically downloads and manages its own binaries and models without touching your system-wide installations or ports.
* **Client-Server Architecture:** Uses Unix Domain Sockets (UDS) for low-latency communication between the client CLI and the background daemon. The daemon manages the lifecycle of the local AI sandbox.
* **Environment Scanning (Indexer):** The daemon automatically detects your operating system, package manager, active shell interpreter, root privileges, and all available executables in `$PATH` to generate compatible, safe commands.
* **Multilingual:** Full native support for any language. The AI thinks, explains, and spins in the exact same language as your query.
* **Optimized VRAM Usage:** Generates commands instantly with aggressively optimized context limits, preventing memory bloat on local setups.
* **Zero-Trust Verification:** Never executes destructive commands blindly. Jumie waits for explicit `y/n` confirmation `(o_o)`.
* **Zero Config Setup:** Run `jum` once, and it will automatically setup the isolated AI sandbox for you.

## Build and Run

Requires Go version `1.26.3` or newer (due to the usage of sequence iterators).

### Compiling

Build optimized binaries inside the `dist/` directory:

```bash
make build
```

### Installation

To install the CLI client and register the background daemon in your system startup (user-level):

```bash
make install
```
* On **macOS**, this will register a LaunchAgent plist.
* On **Linux**, this will enable and start a user-level Systemd unit.

To completely uninstall and wipe the daemon configuration:

```bash
make uninstall
```

## Example Usage

Once the daemon is installed, simply ask `jum` to do anything. If the isolated AI sandbox isn't setup yet, it will prompt you to download it on the first run.

```bash
jum show hardware stats
```

```text
✦ jumie reasoning: The user wants to gather key hardware specs. The OS is macOS. I will use system_profiler to extract hardware and software details, and filter the output using grep.

✦ jumie plan:
➜  gathering key hardware specs of your Mac using system_profiler
$  system_profiler SPHardwareDataType SPSoftwareDataType | grep -E 'Model Name|Processor Name|Memory|OS Version'

(o_o) execute? (y/n): y

running: system_profiler SPHardwareDataType SPSoftwareDataType | grep -E 'Model Name|Processor Name|Memory|OS Version'
      Model Name: MacBook Pro
      Processor Name: Apple M4 Pro
      Memory: 16 GB
      System Version: macOS 15.5 (24F70)
```
