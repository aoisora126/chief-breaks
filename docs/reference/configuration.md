---
description: Chief configuration reference. CLI flags, environment variables, and settings for customizing Chief's behavior.
---

# Configuration

Chief is designed to work with zero configuration. All state lives in `.chief/` and settings are passed via CLI flags.

## CLI Flags vs Environment Variables

Most settings can be provided either way:

| CLI Flag | Environment Variable | Description |
|----------|---------------------|-------------|
| `--prd <name>` | `CHIEF_PRD` | Which PRD to use |
| `--max-iterations <n>` | `CHIEF_MAX_ITERATIONS` | Loop iteration limit |

CLI flags take precedence over environment variables.

## Configuration File (Optional)

For convenience, you can create a `.chief/config.json`:

```json
{
  "defaultPrd": "main-feature",
  "maxIterations": 150,
  "sound": true
}
```

## Claude Code Configuration

Chief invokes Claude Code under the hood. Claude Code has its own configuration:

```bash
# Authentication
claude login

# Model selection (if you have access)
claude config set model claude-3-opus-20240229
```

See [Claude Code documentation](https://github.com/anthropics/claude-code) for details.

## Permission Handling

By default, Claude Code asks for permission before:
- Executing bash commands
- Writing files
- Making network requests

For autonomous operation, use:

```bash
chief --dangerously-skip-permissions
```

::: warning
Only use `--dangerously-skip-permissions` when you trust the PRD and are prepared for Claude to make changes to your codebase.
:::

## Project-Specific Settings

Some settings are best kept in the PRD itself:

```json
{
  "project": "My Feature",
  "settings": {
    "testCommand": "npm test",
    "buildCommand": "npm run build",
    "lintCommand": "npm run lint"
  },
  "userStories": [...]
}
```

Chief passes these to Claude, which uses them when running quality checks.

## No Global Config

Intentionally, Chief has no global configuration file. This ensures:

1. **Portability** - Project works the same on any machine
2. **Reproducibility** - No hidden state affecting behavior
3. **Simplicity** - One place to look for all settings
