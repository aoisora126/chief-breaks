# FAQ

Frequently asked questions about Chief.

## General

### What is Chief?

Chief is an autonomous PRD agent. You write a Product Requirements Document with user stories, run Chief, and watch as Claude builds your code—story by story.

### Why "Chief"?

Chief manages the project. You write the requirements, Chief orchestrates the work.

### Is Chief free?

Chief itself is open source and free. However, it uses Claude Code, which requires an Anthropic API subscription.

### What models does Chief use?

Chief uses whatever model is configured in Claude Code. By default, this is Claude 3 Sonnet.

## Usage

### Can I run Chief on a remote server?

Yes! Chief works great on remote servers. SSH in, run `chief`, and let it work. Use `screen` or `tmux` if you want to disconnect.

```bash
ssh my-server
tmux new -s chief
chief --dangerously-skip-permissions
# Ctrl+B D to detach
```

### How do I resume after stopping?

Just run `chief` again. It reads state from `prd.json` and continues where it left off.

### Can I edit the PRD while Chief is running?

Yes, but be careful. Chief re-reads `prd.json` between iterations. Edits to the current story might cause confusion.

Best practice: pause Chief (Ctrl+C), edit, then resume.

### Can I have multiple PRDs?

Yes. Create separate directories under `.chief/prds/`:

```
.chief/prds/
├── feature-a/
└── feature-b/
```

Run with: `chief --prd feature-a`

### How do I skip a story?

Mark it as passed manually:

```json
{
  "id": "US-003",
  "passes": true,
  "inProgress": false
}
```

Or remove it from the PRD entirely.

## Technical

### Why stream-json?

Claude Code outputs JSON in a streaming format. Chief uses stream-json to parse this in real-time, allowing it to:
- Display progress as it happens
- React to completion signals immediately
- Handle large outputs efficiently

### Why conventional commits?

Conventional commits (`feat:`, `fix:`, etc.) provide:
- Clear history of what each story added
- Easy to review changes per-story
- Works with changelog generators

### What if Claude makes a mistake?

Git is your safety net. Each story is committed separately, so you can:

```bash
# See what changed
git log --oneline

# Revert a story
git revert HEAD

# Or reset and re-run
git reset --hard HEAD~1
chief
```

### Does Chief work with any language?

Yes. Chief doesn't know or care what language you're using. It passes your PRD to Claude, which handles the implementation.

### How does Chief handle tests?

Chief instructs Claude to run quality checks (tests, lint, typecheck) before committing. The specific commands come from your PRD's settings or Claude's inference from your codebase.

## Troubleshooting

### See [Common Issues](/troubleshooting/common-issues)

For specific problems and solutions.

## Getting Help

### Where can I report bugs?

[GitHub Issues](https://github.com/minicodemonkey/chief/issues)

### Is there a community chat?

Coming soon. For now, use GitHub Discussions.

### Can I contribute?

Yes! See [CONTRIBUTING.md](https://github.com/minicodemonkey/chief/blob/main/CONTRIBUTING.md) in the repository.
