# ctx init: Local Plugin Enablement

## Problem

`ctx init` enables the ctx plugin in `~/.claude/settings.json`
(global) but does not write `enabledPlugins` to the project-level
`.claude/settings.local.json`. Plugin hooks only fire when the
plugin is enabled in the settings that Claude Code merges for the
project. Users with global enablement are unaffected, but new
users who only install at project level get silent hook failures.

## Approach

Add `plugin.EnableLocally()` that writes `enabledPlugins` to
`.claude/settings.local.json`, called from `ctx init` alongside
the existing `EnableGlobally()`. Both are gated by the existing
`--no-plugin-enable` flag.

## Behavior

- Reads existing `settings.local.json` (preserves permissions,
  hooks, and other keys)
- Adds `enabledPlugins: {"ctx@activememory-ctx": true}`
- Skips with message if already enabled
- Prints confirmation with file path when written
- Non-fatal: warns on error, continues init

## Non-Goals

- Interactive "choose where to install" prompt
- Removing or changing the global enable behavior
