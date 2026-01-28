# How Chief Works

Chief is an autonomous PRD agent that turns product requirements into working code.

## The Concept

You write a PRD (Product Requirements Document) with user stories. Chief reads the PRD, invokes Claude, and watches as code gets written—one story at a time.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   You       │────▶│   Chief     │────▶│   Claude    │
│   (PRD)     │     │   (Loop)    │     │   (Code)    │
└─────────────┘     └─────────────┘     └─────────────┘
```

## One Iteration, One Story

Chief works through your PRD methodically:

1. **Read** - Finds the highest-priority story that hasn't passed yet
2. **Execute** - Invokes Claude with context about the story
3. **Complete** - Claude works until the story is done
4. **Repeat** - Moves to the next story

Each story is isolated. Claude commits its changes before moving on. If something breaks, you know exactly which story caused it.

## Why This Works

Traditional AI coding assistants require constant interaction. You prompt, it responds, you prompt again.

Chief is different. You define what you want upfront, then step back. Claude works autonomously, following the PRD as its guide.

This enables:

- **Background execution** - SSH in, run `chief`, disconnect
- **Predictable output** - Conventional commits, tracked progress
- **Resumable work** - Stop anytime, continue where you left off

## Further Reading

- [The Ralph Loop](/concepts/ralph-loop) - Deep dive into the execution loop
- [PRD Format](/concepts/prd-format) - How to write effective PRDs
- [The .chief Directory](/concepts/chief-directory) - Where state is stored
