# Installation

Chief is distributed as a single binary with no runtime dependencies.

## Homebrew (Recommended)

The easiest way to install Chief on macOS or Linux:

```bash
brew install minicodemonkey/chief/chief
```

## Install Script

Download and install with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/minicodemonkey/chief/main/install.sh | bash
```

### Script Options

```bash
# Install a specific version
curl -fsSL ... | bash -s -- --version v0.1.0

# Install to a custom directory
curl -fsSL ... | bash -s -- --dir /opt/chief
```

## Manual Download

Download the binary for your platform from the [releases page](https://github.com/minicodemonkey/chief/releases).

| Platform | Architecture | File |
|----------|-------------|------|
| macOS | Apple Silicon | `chief-darwin-arm64` |
| macOS | Intel | `chief-darwin-amd64` |
| Linux | x64 | `chief-linux-amd64` |
| Linux | ARM64 | `chief-linux-arm64` |

After downloading:

```bash
chmod +x chief-*
mv chief-* /usr/local/bin/chief
```

## Building from Source

Requires Go 1.21 or later:

```bash
git clone https://github.com/minicodemonkey/chief.git
cd chief
go build -o chief
```

## Verify Installation

```bash
chief --version
```

## Prerequisites

Chief requires Claude Code CLI to be installed and authenticated:

```bash
# Install Claude Code
npm install -g @anthropic-ai/claude-code

# Authenticate (opens browser)
claude login
```
