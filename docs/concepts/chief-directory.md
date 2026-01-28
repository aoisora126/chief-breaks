# The .chief Directory

Chief stores all state in a `.chief/` directory at the root of your project. This design is intentional—everything is self-contained and portable.

## Directory Structure

```
.chief/
└── prds/
    └── my-feature/
        ├── prd.md        # Human-readable PRD
        ├── prd.json      # Machine-readable PRD
        ├── progress.md   # Progress log
        └── claude.log    # Raw Claude output
```

## Multiple PRDs

You can have multiple PRDs in the same project:

```
.chief/
└── prds/
    ├── auth-system/
    │   ├── prd.md
    │   └── prd.json
    ├── payment-integration/
    │   ├── prd.md
    │   └── prd.json
    └── admin-dashboard/
        ├── prd.md
        └── prd.json
```

Run Chief with a specific PRD:

```bash
chief --prd auth-system
```

## File Explanations

### prd.md

Your product requirements in markdown. Write context, background, technical notes—anything that helps Claude understand the project.

### prd.json

Structured PRD data. Chief reads and writes this file:
- Reads to find the next story
- Writes to mark stories complete

### progress.md

Auto-generated log of completed work. Chief appends to this after each story:

```markdown
## 2024-01-15 - US-001
- What was implemented
- Files changed
- Learnings for future iterations
---
```

### claude.log

Raw output from Claude. Useful for debugging if something goes wrong.

## Portability

Move your project, state moves with it:

```bash
# Works perfectly
mv my-project /new/location/
cd /new/location/my-project
chief  # Picks up right where it left off
```

No global config files. No hidden state in your home directory. Everything in `.chief/`.

## Git Considerations

### What to Commit

- `prd.md` - Your requirements (definitely commit)
- `prd.json` - Story state (commit to share progress)
- `progress.md` - Progress log (commit for history)

### What to Ignore

- `claude.log` - Can be large, regenerated each run

Add to `.gitignore`:

```
.chief/prds/*/claude.log
```

## Why Self-Contained?

This design enables:

1. **Reproducibility** - Clone the repo, run chief, same result
2. **Collaboration** - Team members see the same PRD state
3. **Portability** - Move projects without breaking anything
4. **Transparency** - All state is inspectable plain text
