# Ralph TUI - Feature Specification

## Overview

Ralph is an autonomous agent loop that orchestrates Claude Code to work through PRD user stories. This spec describes a TUI application that wraps the agent loop with monitoring, controls, and a delightful developer experience.

## Goals

1. **Delightful DX** - Make monitoring and controlling the agent loop a pleasure
2. **Easy Distribution** - Single binary, no dependencies, cross-platform
3. **Simple Core** - The actual loop should be ~80 lines, easy to understand and debug
4. **Self-Contained** - Embed the agent prompt, PRD skills, and completion sound

## Non-Goals

- Branch management (removed - let users handle git themselves)
- Headless/CI mode (not needed for v1)
- Settings persistence (CLI flags are sufficient)

## Technology Choice: Go + Bubble Tea

**Why Go?**
- Single binary distribution (no runtime dependencies)
- Cross-compilation via goreleaser (darwin/linux/windows, amd64/arm64)
- Built-in JSON parsing, no external deps needed
- Excellent TUI library ecosystem

**Why Bubble Tea?**
- Modern, composable TUI framework
- Great keyboard handling and focus management
- Built-in support for async operations
- Active community and maintenance

**Alternatives Considered:**
| Option | Pros | Cons |
|--------|------|------|
| Bash + dialog | Simple | Limited, ugly, no Windows |
| Rust + ratatui | Fast, single binary | Steeper learning curve |
| Python + textual | Quick to build | Requires Python runtime |
| Node + ink | React-like | Requires Node runtime |

## Architecture

```
ralph/
├── cmd/ralph/
│   └── main.go                  # CLI entry, flag parsing
├── internal/
│   ├── loop/
│   │   ├── loop.go              # Core loop (~80 lines)
│   │   └── parser.go            # Parse stream-json → events
│   ├── prd/
│   │   ├── types.go             # PRD structs
│   │   ├── loader.go            # Load, watch, list PRDs
│   │   └── generator.go         # `ralph init` (launches Claude)
│   ├── progress/
│   │   └── progress.go          # Append to progress.txt
│   ├── tui/
│   │   ├── app.go               # Main Bubble Tea model
│   │   ├── dashboard.go         # Dashboard view (tasks + details)
│   │   ├── log.go               # Pretty log viewer
│   │   ├── picker.go            # PRD picker modal
│   │   └── styles.go            # Lip Gloss styles
│   └── notify/
│       └── sound.go             # Embed + play completion sound
├── embed/
│   ├── prompt.txt               # Agent prompt
│   ├── prd_skill.txt            # PRD generator prompt
│   ├── convert_skill.txt        # PRD→JSON converter prompt
│   └── complete.wav             # ~30KB completion chime
└── go.mod
```

## Core Loop Design

The loop must be **dead simple** - anyone reading the code should immediately understand it:

```go
// internal/loop/loop.go - The ENTIRE loop logic

func (l *Loop) RunIteration(ctx context.Context) error {
    // This is the only "magic" - just calling claude with args
    cmd := exec.CommandContext(ctx, "claude",
        "--dangerously-skip-permissions",
        "-p", l.prompt,
        "--output-format", "stream-json",
        "--verbose",
    )

    stdout, _ := cmd.StdoutPipe()
    cmd.Start()

    // Parse stream-json and emit events to TUI
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        l.handleLine(scanner.Text())
    }

    return cmd.Wait()
}
```

**Key principle**: No magic. Just `claude` with flags.

## File Structure

When ralph runs in a project:

```
your-project/
├── ralph/
│   ├── prd.md                # Human-readable PRD (from `ralph init`)
│   ├── prd.json              # Machine-readable PRD (from `ralph convert`)
│   ├── prd-backend.json      # Optional additional PRD
│   ├── prd-auth.json         # Optional additional PRD
│   ├── progress.txt          # Human-readable progress log
│   └── .output.log           # Raw Claude output
├── src/
└── ...
```

## PRD Schema

```json
{
  "project": "Project Name",
  "description": "Feature description",
  "userStories": [
    {
      "id": "US-001",
      "title": "Story title",
      "description": "As a..., I need... so that...",
      "acceptanceCriteria": [
        "Criterion 1",
        "Criterion 2",
        "Typecheck passes"
      ],
      "priority": 1,
      "passes": false
    }
  ]
}
```

**Priority ordering:** Lower number = higher priority = do first. Stories should be ordered by dependency (schema → backend → frontend → polish).

**Status tracking via PRD (set by Claude at runtime):**
- `inProgress: true` - Claude sets this when starting a story
- `passes: true` - Claude sets this when story is complete
- `inProgress: false` - Claude sets this when story is complete (along with passes)
- The TUI watches prd.json for changes to update the display

**Note:** `inProgress` is not in the initial prd.json — Claude adds it at runtime.

## CLI Interface

```bash
# Main usage
ralph                      # Auto-detect PRD in ./ralph/, start TUI
ralph ./ralph/prd.json     # Explicit PRD file

# PRD generation (launches Claude as subprocess)
ralph init                 # Interactive: describe feature → generate PRD
ralph init "user auth"     # Non-interactive: generate PRD for "user auth"
ralph convert prd.md       # Convert markdown PRD to prd.json

# Options
ralph --max-iterations 40  # Iteration limit (default: 10)
ralph --no-sound           # Disable completion sound
ralph --verbose            # Show raw Claude output in log

# Note: One iteration = one Claude invocation = typically one story.
# If you have 15 stories, set --max-iterations to at least 15.
# The limit prevents runaway loops and excessive API usage.

# Quick commands (no TUI)
ralph status               # Print current progress, exit
ralph list                 # List all PRDs in ./ralph/
```

## TUI Design

### Design Principles

- **Modern & minimal** — Clean lines, generous spacing, clear hierarchy
- **Information-dense but not cluttered** — Show what matters, hide what doesn't
- **Keyboard-first** — All actions accessible via keyboard, shortcuts always visible
- **Status at a glance** — Current state obvious within 1 second of looking
- **Responsive** — Gracefully handles narrow terminals (min 80 cols) and wide terminals (120+ cols)

### Color Palette (Lip Gloss)

| Element | Color | Hex |
|---------|-------|-----|
| Primary accent | Cyan | `#00D7FF` |
| Success | Green | `#5AF78E` |
| Warning | Yellow | `#F3F99D` |
| Error | Red | `#FF5C57` |
| Muted text | Gray | `#6C7086` |
| Border | Dim gray | `#45475A` |
| Background | Terminal default | — |

### Task Status Indicators

| Symbol | State | Color |
|--------|-------|-------|
| `▶` | In progress | Cyan (animated pulse) |
| `✓` | Completed | Green |
| `○` | Pending | Muted gray |
| `✗` | Failed | Red |
| `⏸` | Paused | Yellow |

---

## Main Dashboard View

The primary view showing task list and details side-by-side.

### Running State

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                          ● RUNNING  Iteration 3/40  00:12:34    │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ─────────────────────────────────────╮ ╭─ Details ─────────────────────────────────────────────╮
│                                               │ │                                                       │
│  ✓  US-101  Set up Tailwind CSS with base     │ │  ▶ US-102 · Configure design tokens                   │
│  ▶  US-102  Configure design tokens           │ │                                                       │
│  ○  US-103  Create color theme system         │ │  ─────────────────────────────────────────────────    │
│  ○  US-104  Build Typography component        │ │                                                       │
│  ○  US-105  Create Button component           │ │  As a developer, I need Tailwind configured with      │
│  ○  US-106  Create Card component             │ │  presentation-appropriate design tokens so that       │
│  ○  US-107  Build responsive grid system      │ │  themes can use consistent, large-scale typography    │
│  ○  US-108  Create navigation header          │ │  and spacing values.                                  │
│  ○  US-109  Implement dark mode toggle        │ │                                                       │
│  ○  US-110  Add page transition animations    │ │  ─────────────────────────────────────────────────    │
│  ○  US-111  Create loading skeleton states    │ │                                                       │
│  ○  US-112  Build toast notification system   │ │  Acceptance Criteria                                  │
│                                               │ │                                                       │
│                                               │ │  ○  Extend fontSize scale (slide-sm to slide-hero)    │
│                                               │ │  ○  Extend spacing scale (slide-1 to slide-32)        │
│                                               │ │  ○  Add fontFamily variants (sans, serif, mono)       │
│                                               │ │  ○  Configure custom breakpoints for slides           │
│                                               │ │  ○  Typecheck passes                                  │
│                                               │ │                                                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │                                                       │
│  1 of 12 complete                         8%  │ │  Priority P1                                          │
│                                               │ │                                                       │
╰───────────────────────────────────────────────╯ ╰───────────────────────────────────────────────────────╯

╭─ Activity ──────────────────────────────────────────────────────────────────────────────────────────────╮
│  Reading tailwind.config.ts to understand current configuration...                                      │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

  p Pause   x Stop   t Log   l Switch PRD   ↑↓ Navigate   ? Help                            prd.json   q Quit
```

### Idle State (Ready to Start)

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                                ○ READY  prd.json  12 stories    │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ─────────────────────────────────────╮ ╭─ Details ─────────────────────────────────────────────╮
│                                               │ │                                                       │
│  ○  US-101  Set up Tailwind CSS with base     │ │  ○ US-101 · Set up Tailwind CSS with base config      │
│  ○  US-102  Configure design tokens           │ │                                                       │
│  ○  US-103  Create color theme system         │ │  ─────────────────────────────────────────────────    │
│  ○  US-104  Build Typography component        │ │                                                       │
│  ○  US-105  Create Button component           │ │  As a developer, I need Tailwind CSS installed and    │
│  ○  US-106  Create Card component             │ │  configured with a base setup so that I can start     │
│  ○  US-107  Build responsive grid system      │ │  building components with utility classes.            │
│  ○  US-108  Create navigation header          │ │                                                       │
│  ○  US-109  Implement dark mode toggle        │ │  ─────────────────────────────────────────────────    │
│  ○  US-110  Add page transition animations    │ │                                                       │
│  ○  US-111  Create loading skeleton states    │ │  Acceptance Criteria                                  │
│  ○  US-112  Build toast notification system   │ │                                                       │
│                                               │ │  ○  Install tailwindcss, postcss, autoprefixer        │
│                                               │ │  ○  Create tailwind.config.ts with TypeScript         │
│                                               │ │  ○  Configure content paths for all components        │
│                                               │ │  ○  Add Tailwind directives to global CSS             │
│                                               │ │  ○  Typecheck passes                                  │
│                                               │ │                                                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │                                                       │
│  0 of 12 complete                         0%  │ │  Priority P1                                          │
│                                               │ │                                                       │
╰───────────────────────────────────────────────╯ ╰───────────────────────────────────────────────────────╯



  s Start   l Switch PRD   ↑↓ Navigate   ? Help                                             prd.json   q Quit
```

### Paused State

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                         ⏸ PAUSED  Iteration 3/40  00:12:34      │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ─────────────────────────────────────╮ ╭─ Details ─────────────────────────────────────────────╮
│                                               │ │                                                       │
│  ✓  US-101  Set up Tailwind CSS with base     │ │  ⏸ US-102 · Configure design tokens                   │
│  ⏸  US-102  Configure design tokens           │ │                                                       │
│  ○  US-103  Create color theme system         │ │  ─────────────────────────────────────────────────    │
│  ...                                          │ │                                                       │
│                                               │ │  Paused after iteration 3. Press s to resume.         │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │                                                       │
│  1 of 12 complete                         8%  │ │                                                       │
╰───────────────────────────────────────────────╯ ╰───────────────────────────────────────────────────────╯

  s Resume   l Switch PRD   ↑↓ Navigate   ? Help                                            prd.json   q Quit
```

### Complete State

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                       ✓ COMPLETE  12 iterations  00:47:23       │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ─────────────────────────────────────╮ ╭─ Summary ─────────────────────────────────────────────╮
│                                               │ │                                                       │
│  ✓  US-101  Set up Tailwind CSS with base     │ │  ✓ All 12 stories complete!                           │
│  ✓  US-102  Configure design tokens           │ │                                                       │
│  ✓  US-103  Create color theme system         │ │  ─────────────────────────────────────────────────    │
│  ✓  US-104  Build Typography component        │ │                                                       │
│  ✓  US-105  Create Button component           │ │  Duration      47m 23s                                │
│  ✓  US-106  Create Card component             │ │  Iterations    12                                     │
│  ✓  US-107  Build responsive grid system      │ │  Stories       12/12                                  │
│  ✓  US-108  Create navigation header          │ │                                                       │
│  ✓  US-109  Implement dark mode toggle        │ │  ─────────────────────────────────────────────────    │
│  ✓  US-110  Add page transition animations    │ │                                                       │
│  ✓  US-111  Create loading skeleton states    │ │  View progress.txt for detailed implementation        │
│  ✓  US-112  Build toast notification system   │ │  notes and learnings.                                 │
│                                               │ │                                                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │                                                       │
│  12 of 12 complete                      100%  │ │                                                       │
│                                               │ │                                                       │
╰───────────────────────────────────────────────╯ ╰───────────────────────────────────────────────────────╯

  l Switch PRD   t View Log   ? Help                                                        prd.json   q Quit
```

---

## Log Viewer

Full-screen view showing Claude's streaming output. Toggle with `t` key.

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                             ● RUNNING  US-102  Iteration 3/40  00:12:34         │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Log ───────────────────────────────────────────────────────────────────────────────────────────────────╮
│                                                                                                         │
│  Reading prd.json to find the next task to work on...                                                   │
│                                                                                                         │
│  The next story is US-102: Configure design tokens. This story has inProgress: false                    │
│  and passes: false, so I'll start working on it now.                                                    │
│                                                                                                         │
│  First, let me update prd.json to mark this story as in progress.                                       │
│                                                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│  │  ✏️  Edit  ralph/prd.json                                                                        │    │
│  └─────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                         │
│  Now let me examine the current Tailwind configuration to understand what's already set up.             │
│                                                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│  │  📖  Read  tailwind.config.ts                                                                    │    │
│  └─────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                         │
│  The config has a basic setup. I need to extend it with presentation-specific scales.                   │
│  I'll add custom fontSize, spacing, and fontFamily values optimized for slide presentations.            │
│                                                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│  │  ✏️  Edit  tailwind.config.ts                                                                    │    │
│  └─────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                         │
│  Let me verify the typecheck still passes with these changes.                                           │
│                                                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│  │  🔨  Bash  npm run typecheck                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                                         │
│  ▌                                                                                                      │
│                                                                                                         │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

  t Dashboard   p Pause   x Stop   ↑↓ jk Scroll   G Bottom   g Top                          prd.json   q Quit
```

**Tool Icons:**

| Tool | Icon |
|------|------|
| Read | 📖 |
| Edit | ✏️ |
| Write | 📝 |
| Bash | 🔨 |
| Glob | 🔍 |
| Grep | 🔎 |
| Task | 🤖 |
| WebFetch | 🌐 |

---

## PRD Picker

Modal overlay for switching between PRDs. Toggle with `l` key.

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                                ○ READY  prd.json  12 stories    │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

        ╭─ Select PRD ────────────────────────────────────────────────────────────────────────╮
        │                                                                                      │
        │   ▶  prd.json                                                        ● Running      │
        │      Tap Documentation Website                                                       │
        │      ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  8/12  67%             │
        │                                                                                      │
        │      prd-api.json                                                    ○ Ready        │
        │      REST API Refactoring                                                            │
        │      ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  0/18   0%             │
        │                                                                                      │
        │      prd-auth.json                                                   ⏸ Paused       │
        │      User Authentication System                                                      │
        │      ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  4/12  33%             │
        │                                                                                      │
        │      prd-mobile.json                                                 ✓ Complete     │
        │      Mobile Responsive Layouts                                                       │
        │      ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  6/6  100%             │
        │                                                                                      │
        ╰──────────────────────────────────────────────────────────────────────────────────────╯

                        ↑↓ Navigate   Enter Select   n New PRD   Esc Back
```

---

## Help Overlay

Modal showing all keyboard shortcuts. Toggle with `?` key.

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                          ● RUNNING  Iteration 3/40  00:12:34    │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

                ╭─ Keyboard Shortcuts ────────────────────────────────────────────────╮
                │                                                                      │
                │   Loop Control                      Navigation                       │
                │   ────────────                      ──────────                       │
                │   s   Start / Resume                ↑ k   Previous story             │
                │   p   Pause after iteration         ↓ j   Next story                 │
                │   x   Stop immediately              g     Go to top                  │
                │                                     G     Go to bottom               │
                │   Views                                                              │
                │   ─────                             Scrolling (Log View)             │
                │   t   Toggle log view               ─────────────────────            │
                │   l   PRD picker                    Ctrl+D   Page down               │
                │   ?   This help                     Ctrl+U   Page up                 │
                │                                                                      │
                │   General                                                            │
                │   ───────                                                            │
                │   r       Refresh PRD                                                │
                │   q       Quit / Back                                                │
                │   Ctrl+C  Force quit                                                 │
                │                                                                      │
                ╰──────────────────────────────────────────────────────────────────────╯

                                           Esc or ? to close
```

---

## Empty State

Shown when no PRD exists in the ralph/ directory.

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                                                   No PRD loaded  │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯




                              ╭──────────────────────────────────────────────╮
                              │                                              │
                              │                  ◇                           │
                              │                                              │
                              │         No PRD found in ./ralph/             │
                              │                                              │
                              │    Get started by creating a new PRD:        │
                              │                                              │
                              │    $ ralph init                              │
                              │      Create a PRD interactively              │
                              │                                              │
                              │    $ ralph init "user authentication"        │
                              │      Generate PRD for a specific feature     │
                              │                                              │
                              │    $ ralph convert ./docs/spec.md            │
                              │      Convert existing spec to prd.json       │
                              │                                              │
                              ╰──────────────────────────────────────────────╯




                                                                                                    q Quit
```

---

## Error State

Shown when an error occurs (e.g., Claude crashes, file not found).

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                            ✗ ERROR  Iteration 3/40  00:12:34    │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ─────────────────────────────────────╮ ╭─ Error ───────────────────────────────────────────────╮
│                                               │ │                                                       │
│  ✓  US-101  Set up Tailwind CSS with base     │ │  ✗ Claude process exited unexpectedly                 │
│  ▶  US-102  Configure design tokens           │ │                                                       │
│  ○  US-103  Create color theme system         │ │  ─────────────────────────────────────────────────    │
│  ○  US-104  Build Typography component        │ │                                                       │
│  ○  US-105  Create Button component           │ │  Exit code: 1                                         │
│  ○  US-106  Create Card component             │ │  Story US-102 was interrupted and will resume         │
│  ○  US-107  Build responsive grid system      │ │  on next iteration.                                   │
│  ○  US-108  Create navigation header          │ │                                                       │
│  ○  US-109  Implement dark mode toggle        │ │  ─────────────────────────────────────────────────    │
│  ○  US-110  Add page transition animations    │ │                                                       │
│  ○  US-111  Create loading skeleton states    │ │  Check .output.log for full error details.            │
│  ○  US-112  Build toast notification system   │ │                                                       │
│                                               │ │                                                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │                                                       │
│  1 of 12 complete                         8%  │ │                                                       │
│                                               │ │                                                       │
╰───────────────────────────────────────────────╯ ╰───────────────────────────────────────────────────────╯

  s Retry   t View Log   l Switch PRD   ? Help                                              prd.json   q Quit
```

---

## Interrupted Story Warning

Shown when ralph starts and detects an `inProgress: true` story from a previous session.

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│  ralph                                                               ⚠ INTERRUPTED  prd.json           │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ─────────────────────────────────────╮ ╭─ Notice ──────────────────────────────────────────────╮
│                                               │ │                                                       │
│  ✓  US-101  Set up Tailwind CSS with base     │ │  ⚠ Previous session was interrupted                   │
│  ▶  US-102  Configure design tokens           │ │                                                       │
│  ○  US-103  Create color theme system         │ │  ─────────────────────────────────────────────────    │
│  ...                                          │ │                                                       │
│                                               │ │  Story US-102 has inProgress: true from a             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │  previous session that didn't complete.               │
│  1 of 12 complete                         8%  │ │                                                       │
│                                               │ │  Press s to resume — the story will be                │
╰───────────────────────────────────────────────╯ │  automatically picked up.                             │
                                                  │                                                       │
                                                  ╰───────────────────────────────────────────────────────╯

  s Resume   l Switch PRD   ↑↓ Navigate   ? Help                                            prd.json   q Quit
```

---

## Narrow Terminal (80 columns)

Graceful degradation for narrower terminals — single column layout.

```
╭──────────────────────────────────────────────────────────────────────────────╮
│  ralph                               ● RUNNING  Iteration 3/40  00:12:34    │
╰──────────────────────────────────────────────────────────────────────────────╯

╭─ Stories ────────────────────────────────────────────────────────────────────╮
│                                                                              │
│  ✓  US-101  Set up Tailwind CSS with base config                             │
│  ▶  US-102  Configure design tokens                                          │
│  ○  US-103  Create color theme system                                        │
│  ○  US-104  Build Typography component                                       │
│  ○  US-105  Create Button component                                          │
│  ○  US-106  Create Card component                                            │
│                                                                              │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  1/12  8%    │
│                                                                              │
╰──────────────────────────────────────────────────────────────────────────────╯

╭─ US-102 ─────────────────────────────────────────────────────────────────────╮
│                                                                              │
│  As a developer, I need Tailwind configured with presentation-appropriate    │
│  design tokens so that themes can use consistent, large-scale typography.    │
│                                                                              │
│  ○  Extend fontSize scale (slide-sm to slide-hero)                           │
│  ○  Extend spacing scale (slide-1 to slide-32)                               │
│  ○  Add fontFamily variants                                                  │
│  ○  Typecheck passes                                                         │
│                                                                              │
╰──────────────────────────────────────────────────────────────────────────────╯

  p Pause  x Stop  t Log  l PRD  ↑↓ Nav  ? Help                          q Quit
```

---

**Multiple loops:** Users can run multiple ralph instances on different PRDs in the same project. Each instance is independent. Trust the user to avoid file conflicts between PRDs.

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `q` | Quit / Back |
| `?` | Show help |
| `Ctrl+C` | Force quit |

### Dashboard

| Key | Action |
|-----|--------|
| `s` | Start/resume agent loop |
| `p` | Pause (after current iteration completes) |
| `x` | Stop immediately (kill Claude process) |
| `r` | Refresh (reload PRD file) |
| `l` | Open loop/PRD picker |
| `t` | Toggle log view |
| `↑/k` | Previous task |
| `↓/j` | Next task |
| `Tab` | Switch panel focus |

### Log View

| Key | Action |
|-----|--------|
| `t` | Back to dashboard |
| `f` | Toggle fullscreen |
| `j/↓` | Scroll down |
| `k/↑` | Scroll up |
| `Ctrl+D` | Page down |
| `Ctrl+U` | Page up |
| `G` | Go to bottom |
| `g` | Go to top |

## Notifications

**Completion sound:** A small (~30KB) pleasant chime embedded in the binary, played when user attention is needed:
- All stories complete successfully (`<ralph-complete/>` received)
- Max iterations reached (loop stops, user needs to decide next steps)

**Cross-platform playback:**
```go
import "github.com/hajimehoshi/oto/v2"  // Cross-platform audio

//go:embed complete.wav
var completeSound []byte

func playComplete() {
    // Use oto for cross-platform WAV playback
}
```

Sound can be disabled with `--no-sound` flag.

## Embedded Prompts

### Agent Prompt (embed/prompt.txt)

```markdown
# Ralph Agent

You are an autonomous agent working through a product requirements document.

## Files

- `ralph/prd.json` — The PRD with user stories
- `ralph/progress.txt` — Progress log (read Codebase Patterns section first)

## Task

1. Read prd.json and select the next story:
   - FIRST: Any story with `inProgress: true` (resume interrupted work)
   - THEN: Story with lowest `priority` number where `passes: false`
2. Set `inProgress: true` on the selected story in prd.json
3. Implement the story completely
4. Run quality checks (typecheck, lint, test as appropriate)
5. For UI changes, verify in browser using Playwright if available
6. Commit changes using conventional commits (see below)
7. Update prd.json: set `passes: true` and `inProgress: false`
8. Append to progress.txt (see format below)

## Conventional Commits

Use this format for all commits:
```
<type>[optional scope]: <description>
```

Types: `feat` (new feature), `fix` (bug fix), `refactor`, `test`, `docs`, `chore`

Examples:
- `feat(auth): add login form validation`
- `fix: prevent race condition in request handler`
- `refactor(api): extract shared validation logic`

Rules:
- Only commit files you modified during this iteration
- Split into multiple commits if logically appropriate
- Never mention Claude or AI in commit messages

## Progress Format

Append to progress.txt (never replace):
```
## YYYY-MM-DD - US-XXX: [Title]
- What was implemented
- Files changed
- **Learnings:** (patterns, gotchas, context for future iterations)
---
```

Add reusable patterns to `## Codebase Patterns` at the top of progress.txt.

## Completion

After each story, check if ALL stories have `passes: true`.
If complete, output: <ralph-complete/>

## Rules

- One story per iteration
- Never commit broken code
- Follow existing code patterns
- Keep changes focused and minimal
```

### PRD Generator Prompt (embed/prd_skill.txt)

Used by `ralph init` - launches Claude to interactively generate a PRD:

```markdown
# PRD Generator

You are helping create a Product Requirements Document.

## Process

1. Ask 3-5 clarifying questions with lettered options (A, B, C, D) about:
   - Problem being solved / goal
   - Core functionality
   - Scope boundaries
   - Success criteria

2. Generate a PRD with:
   - Introduction
   - Goals (measurable)
   - User Stories with acceptance criteria
   - Functional requirements (numbered)
   - Non-Goals (explicit scope boundaries)
   - Design considerations
   - Technical considerations
   - Success metrics
   - Open questions

3. Save to `ralph/prd.md`

## User Story Format

Each story should be:
- Small enough to complete in ONE Claude context window (one iteration)
- Have specific, verifiable acceptance criteria (not vague)
- Include "Typecheck passes" as criterion
- For UI changes, include "Verify in browser using Playwright"

**Right-sized:** database column addition, single UI component, server action update
**Too large (split these):** complete dashboard, full auth system, API refactor

## Output

Save the PRD as markdown to `ralph/prd.md`, then inform the user:
"PRD saved to ralph/prd.md. Run `ralph convert` to generate prd.json"
```

### PRD Converter Prompt (embed/convert_skill.txt)

Used by `ralph convert` - converts markdown PRD to JSON:

```markdown
# PRD Converter

Convert the PRD markdown file to ralph's prd.json format.

## Input

Read the PRD from `ralph/prd.md` (or path provided by user).

## Output Format

```json
{
  "project": "[Project name from PRD]",
  "description": "[Brief description]",
  "userStories": [
    {
      "id": "US-001",
      "title": "[Short title]",
      "description": "[Full story: As a..., I need..., so that...]",
      "acceptanceCriteria": ["Criterion 1", "Criterion 2", "Typecheck passes"],
      "priority": 1,
      "passes": false
    }
  ]
}
```

**Note:** `inProgress` is NOT set here — Claude adds it at runtime.

## Rules

1. **Story sizing**: Each story must complete in ONE iteration (one context window). If describing the change takes more than 2-3 sentences, split it.
2. **Priority order** (lower number = do first): Schema/migrations → Backend/server actions → Frontend/UI → Dashboards/aggregations
3. **Acceptance criteria**: Must be verifiable, not vague. Always include "Typecheck passes". For UI, include "Verify in browser using Playwright".
4. **Dependencies**: No forward dependencies. Story N can only depend on stories 1 to N-1.

## Save

Save to `ralph/prd.json` and confirm to user.
```

## Data Flow

```
┌──────────────┐     ┌───────────────┐     ┌─────────────┐
│   PRD File   │────▶│  Agent Loop   │────▶│  Progress   │
│  (prd.json)  │◀────│   (Claude)    │     │ (progress.txt)
└──────────────┘     └───────────────┘     └─────────────┘
       │                    │
       │  watches for       │  streams
       │  inProgress/passes │  output
       ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                    TUI (Bubble Tea)                     │
│  ┌─────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Tasks   │  │   Details   │  │    Log Viewer       │  │
│  │ Panel   │  │   Panel     │  │    (streaming)      │  │
│  └─────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**Source of truth:** `prd.json` is the only state file. The TUI reads it to display task status and watches for changes.

## State Management

### Loop States

```go
type LoopState int

const (
    StateReady LoopState = iota    // Waiting to start
    StateRunning                    // Claude is executing
    StatePaused                     // Will stop after current iteration
    StateStopping                   // Stop requested, waiting for Claude
    StateComplete                   // All tasks done
    StateError                      // Something went wrong
)
```

### TUI Model

```go
type Model struct {
    // State (derived from prd.json)
    state        LoopState
    prd          *PRD
    selectedTask int

    // Loop
    iteration    int
    maxIter      int
    claudeCmd    *exec.Cmd

    // Views
    activeView   View  // Dashboard, Log, Picker
    logBuffer    *ring.Buffer

    // Components
    taskList     list.Model
    viewport     viewport.Model
    help         help.Model
}
```

**Note:** All persistent state lives in `prd.json`. The TUI model is ephemeral — if ralph restarts, it re-reads prd.json to determine current status (any story with `inProgress: true` was interrupted).

## Error Handling

### Claude Process Errors

- Detect non-zero exit codes
- Parse error messages from stream-json
- Display in TUI with option to retry or skip
- Log full error context to `.output.log`

### Recovery

- If Claude crashes mid-story, `inProgress` stays true in prd.json
- Next iteration automatically resumes the interrupted story (prompt prioritizes `inProgress: true`)
- Failed iterations still count toward max-iterations limit
- TUI shows warning: "Story US-XXX was interrupted — resuming"

### File System Errors

- Handle missing prd.json gracefully (show picker or init prompt)
- Auto-create progress.txt if missing
- Watch for external file changes (hot reload PRD)

## Distribution

### Build Targets

```bash
# Via goreleaser
goreleaser release --snapshot --clean
```

Targets:
- darwin/amd64
- darwin/arm64
- linux/amd64
- linux/arm64
- windows/amd64

### Installation Methods

```bash
# Homebrew (macOS/Linux)
brew install ralph

# Go install
go install github.com/snarktank/ralph@latest

# Download binary
curl -fsSL https://ralph.sh/install.sh | sh
```

## Implementation Phases

### Phase 1: Core

- [ ] Go project setup with Bubble Tea
- [ ] Embedded agent prompt
- [ ] Core loop (~80 lines)
- [ ] Stream-json parser
- [ ] Basic dashboard view (task list + details)
- [ ] Start/pause/stop controls
- [ ] PRD file watching

### Phase 2: Full TUI

- [ ] Pretty log viewer with tool cards
- [ ] PRD picker for multiple loops
- [ ] Progress bar component
- [ ] Keyboard navigation
- [ ] Help overlay

### Phase 3: PRD Generation

- [ ] `ralph init` command (Claude subprocess)
- [ ] `ralph convert` command (Claude subprocess)
- [ ] Embedded skill prompts

### Phase 4: Polish

- [ ] Completion sound (embedded WAV)
- [ ] Error recovery UX
- [ ] `ralph status` quick command
- [ ] `ralph list` quick command

### Phase 5: Distribution

- [ ] goreleaser config
- [ ] Homebrew formula
- [ ] Install script
- [ ] README and docs

## Future Enhancements (Post-MVP)

- Subagent monitoring (track Task tool spawns)
- Cost tracking (parse API usage from stream-json)
- Git integration (show commits made during session)
- Diff preview (show pending changes)
- Web UI (optional browser-based dashboard)
- Team mode (multiple users watching same session)
