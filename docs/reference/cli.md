# CLI Reference

Chief provides a minimal but powerful CLI.

## Commands

### chief (default)

Run the Ralph Loop on the active PRD.

```bash
chief
```

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `--prd <name>` | Select which PRD to run | Auto-detect |
| `--max-iterations <n>` | Maximum loop iterations | 100 |
| `--dangerously-skip-permissions` | Skip permission prompts | false |
| `--no-sound` | Disable completion sound | false |

**Examples:**

```bash
# Run with a specific PRD
chief --prd auth-system

# Increase iteration limit
chief --max-iterations 200

# Skip permission prompts (for CI/automation)
chief --dangerously-skip-permissions
```

### chief init

Create a new PRD in the current project.

```bash
chief init
```

Creates `.chief/prds/<name>/` with template `prd.md` and `prd.json` files.

**Interactive prompts:**
- PRD name (e.g., "auth-system")
- Project description
- First user story

### chief edit

Open the PRD for editing.

```bash
chief edit
```

Opens `prd.md` and `prd.json` in your `$EDITOR`.

**Options:**

| Flag | Description |
|------|-------------|
| `--prd <name>` | Edit a specific PRD |

### chief status

Show current PRD status.

```bash
chief status
```

Displays:
- Current PRD name
- Total stories
- Completed / In Progress / Pending counts
- Next story to be worked on

### chief list

List all PRDs in the project.

```bash
chief list
```

## Keyboard Shortcuts (TUI)

When Chief is running, the TUI supports these controls:

| Key | Action |
|-----|--------|
| `Tab` | Switch between panels |
| `↑/↓` | Scroll log output |
| `q` | Quit (gracefully stops Claude) |
| `Ctrl+C` | Force quit |
| `?` | Show help |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success - all stories complete |
| 1 | Error - something went wrong |
| 2 | Interrupted - user quit |
| 3 | Max iterations reached |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CHIEF_PRD` | Default PRD to use |
| `CHIEF_MAX_ITERATIONS` | Default iteration limit |
| `EDITOR` | Editor for `chief edit` |
