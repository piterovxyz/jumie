<p align="center">
  <img src="logo.svg" alt="jumie logo" width="200"/>
</p>

# jumie

A smart Unix terminal assistant powered by AI.
It runs as a lightweight background daemon and a quick CLI client, translating natural language queries into safe, context-aware command pipelines.
This tool is designed to make command line navigation accessible and secure, giving you a full, human-friendly explanation of actions before executing them.

## Key Features

* **Client-Server Architecture:** Uses Unix Domain Sockets (UDS) for low-latency communication between the client CLI and the background daemon.
* **Environment Scanning (Indexer):** The daemon automatically detects your operating system, package manager, active shell interpreter, root privileges, and all available executables in `$PATH` to generate compatible, safe commands.
* **AI-Powered (Gemini):** Integrated with the official Google GenAI Go SDK using `gemini-3.1-flash-lite`.
* **Zero-Trust Verification:** Never executes destructive commands blindly. Jumie waits for explicit `y/n` confirmation `(o_-)`.
* **BYOK:** Local configuration at `~/.config/jumie/config.json`.
* **Easy Init & Control:** Makefile templates to register the user-space daemon on macOS (`launchd`) and Linux (`systemd`) without requiring sudo privileges.

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

### Authentication (BYOK)

Before using the tool, authenticate with your Google Gemini API Key:

```bash
jum login YOUR_GEMINI_API_KEY
```
This command will validate the key with a quick test API request and save it securely in `~/.config/jumie/config.json`. If you run any command without a configured key, `jum` will automatically prompt you to enter it.

## Example Usage

Once the daemon is installed and you are logged in, simply ask `jum` to do anything:

```bash
jum show hardware stats
```

```text
✦ jumie plan:
➜  let's gather key hardware specs of your Mac to present it nicely
$  system_profiler SPHardwareDataType SPSoftwareDataType | grep -E 'Model Name|Processor Name|Memory|OS Version'

(o_o) execute? (y/n): y

running: system_profiler SPHardwareDataType SPSoftwareDataType | grep -E 'Model Name|Processor Name|Memory|OS Version'
      Model Name: MacBook Pro
      Processor Name: Apple M4 Pro
      Memory: 16 GB
      System Version: macOS 15.5 (24F70)

successfully executed all commands!
```
