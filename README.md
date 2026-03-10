# Chief

<p align="center">
  <img src="assets/hero.png" alt="Chief" width="500">
</p>

Build big projects with Claude. Chief breaks your work into tasks and runs Claude Code in a loop until they're done.
## Install

```bash
brew install minicodemonkey/chief/chief
```

Or via install script:

```bash
curl -fsSL https://raw.githubusercontent.com/MiniCodeMonkey/chief/refs/heads/main/install.sh | sh
```

## Usage

```bash
# Create a new project
chief new

# Launch the TUI and press 's' to start
chief
```

Chief runs Claude in a [Ralph Wiggum loop](https://ghuntley.com/ralph/): each iteration starts with a fresh context window, but progress is persisted between runs. This lets Claude work through large projects without hitting context limits.

## How It Works

1. **Describe your project** as a series of tasks
2. **Chief runs Claude** in a loop, one task at a time
3. **One commit per task** — clean git history, easy to review

See the [documentation](https://minicodemonkey.github.io/chief/concepts/how-it-works) for details.

## Requirements

- **[Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code)**, **[Codex CLI](https://developers.openai.com/codex/cli/reference)**, or **[OpenCode CLI](https://opencode.ai)** installed and authenticated

Use Claude by default, or configure Codex or OpenCode in `.chief/config.yaml`:

```yaml
agent:
  provider: opencode
  cliPath: /usr/local/bin/opencode   # optional
```

Or run with `chief --agent opencode` or set `CHIEF_AGENT=opencode`.

