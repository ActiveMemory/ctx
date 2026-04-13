---
name: ctx-config
description: "Manage runtime configuration profiles. Use when asked to switch to dev mode, check active profile, toggle verbose logging, or switch to base."
allowed-tools: Bash(ctx config*)
---

Manage `.ctxrc` configuration profiles (dev vs base).

## When to Use

- User says "switch to dev mode" / "switch to dev"
- User says "switch to base" / "switch to prod"
- User says "what profile am I on?"
- User says "toggle verbose logging"
- User says "show config status"

## Commands

```bash
# Switch to a specific profile
ctx config switch dev
ctx config switch base

# Toggle between profiles (no argument)
ctx config switch

# Show which profile is active
ctx config status
```

## Profiles

| Profile | Description                                 |
|---------|---------------------------------------------|
| `dev`   | Verbose logging, webhook notifications on   |
| `base`  | All defaults, notifications off             |

Source files (`.ctxrc.base`, `.ctxrc.dev`) are committed to git.
The working copy (`.ctxrc`) is gitignored.

## Process

1. Determine which operation the user wants (switch or status)
2. Run the appropriate `ctx config` command
3. Report the result to the user
