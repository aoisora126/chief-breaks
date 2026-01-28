# PRD Format

Chief uses a structured PRD format with two files: a human-readable markdown file and a machine-readable JSON file.

## File Structure

```
.chief/prds/my-feature/
├── prd.md      # Human-readable description
├── prd.json    # Structured data for Chief
└── progress.md # Auto-generated progress log
```

## prd.md

The markdown file is for humans. Write whatever helps Claude understand the project:

```markdown
# My Feature

## Overview
We're building a user authentication system...

## Technical Context
The backend uses Express.js with PostgreSQL...

## Design Notes
Follow the existing patterns in the codebase...
```

This file is included in Claude's context but not parsed programmatically.

## prd.json

The JSON file is the source of truth for Chief:

```json
{
  "project": "My Feature",
  "description": "A user authentication system",
  "userStories": [
    {
      "id": "US-001",
      "title": "User Registration",
      "description": "As a user, I want to register...",
      "acceptanceCriteria": [
        "Registration form with email/password",
        "Email validation",
        "Password strength requirements"
      ],
      "priority": 1,
      "passes": false,
      "inProgress": false
    }
  ]
}
```

## Story Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (e.g., US-001) |
| `title` | string | Short, descriptive title |
| `description` | string | User story format recommended |
| `acceptanceCriteria` | string[] | What must be true when done |
| `priority` | number | Lower = higher priority |
| `passes` | boolean | Has this story been completed? |
| `inProgress` | boolean | Is Claude working on this now? |

## Story Selection Logic

Chief selects the next story using:

1. Filter to stories where `passes: false`
2. Skip stories where `inProgress: true` (being worked on)
3. Sort by `priority` ascending
4. Take the first one

## Best Practices

### Write Clear Acceptance Criteria

```json
// ✓ Good - specific and testable
"acceptanceCriteria": [
  "Login form with email and password fields",
  "Error message shown for invalid credentials",
  "Redirect to dashboard on success"
]

// ✗ Bad - vague and subjective
"acceptanceCriteria": [
  "Nice login page",
  "Good error handling"
]
```

### Keep Stories Small

One feature per story. If a story has more than 5-7 acceptance criteria, consider splitting it.

### Use Consistent IDs

Stick to a pattern like `US-001`, `US-002`. The IDs appear in commit messages.

## See Also

- [PRD Schema Reference](/reference/prd-schema) - Complete schema documentation
