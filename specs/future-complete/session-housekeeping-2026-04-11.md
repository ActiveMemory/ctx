---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Session housekeeping 2026-04-11 — parent examples, journal decision, docs restructure
status: accepted
date: 2026-04-11
owner: jose
scope: three small independent cleanups rolled into one commit
---

# Spec — Session housekeeping bundle

## What this covers

Three small, unrelated cleanups that were completed during the
2026-04-11 working session and need a spec to anchor the commit
trailer. Each is self-contained and did not warrant its own standalone
spec at the scale-to-work level.

### 1. Parent command `Example:` wiring

`ctx pad`, `ctx notify`, and `ctx remind` parent commands were the
only three (of 135 total) Cobra commands in the tree without an
`Example:` field. Their `examples.yaml` entries already existed from
an earlier pass; the Go wiring was missing.

Change: add `Example: desc.Example(cmd.DescKey<Parent>)` to the
three parent `Cmd()` functions in:

- `internal/cli/pad/pad.go`
- `internal/cli/notify/notify.go`
- `internal/cli/remind/remind.go`

Result: 100% `Example:` coverage across all 135 commands and
subcommands.

### 2. Journal-stays-local decision

New decision record in `.context/DECISIONS.md` under id
`[2026-04-11-200000]` titled "Journal stays local; LEARNINGS.md is
the shareable layer". Resolves the question of whether `ctx hub`
should ever sync raw journal entries, or whether a future "export
enriched entries as shareable learning items" pipeline should exist.

Decision: **no**. The journal is Tier-0 personal; LEARNINGS /
DECISIONS / CONVENTIONS are Tier-1 shareable; `/ctx-journal-enrich`
is the promotion boundary. Hub sync code paths must exclude
`.context/journal/` at the code level, not just via gitignore.
Follow-up tasks (test-level enforcement, hub doc updates) are
latent and will be scheduled as needed.

### 3. `docs/home/context-files.md` restructure

Removed `templates/`, `steering/`, `hooks/`, `skills/` rows from the
"File Overview" table because they conflated:

- `.context/` substrate files (user-facing, read by priority order)
- `.context/`-internal subdirectories that are implementation
  details (`templates/`, `steering/`)
- Things not under `.context/` at all (`hooks/` — plugin-owned,
  `skills/` — plugin-owned or `.claude/skills/`)

Replaced with:

- The 8-row table of actual core context files.
- A note labelling `templates/` and `steering/` as
  implementation-detail-but-user-editable with links.
- An "Outside `.context/`" subsection pointing at the skill and hook
  home pages.
- A new `## steering/` dedicated section parallel to the existing
  `## templates/` section.

Generated site counterparts (`site/home/context-files/index.html`,
`site/search.json`) regenerate and are committed alongside for tree
hygiene.

## Validation

Each file passed the hard gate (build, lint, test) during the
session immediately after its change landed. Bundling for the commit
does not change the per-file validation status.

## Out of scope

- Any further restructure of the `File Overview` table layout.
- Promoting parent commands out of `ctx system` (tracked by
  `specs/cli-namespace-cleanup.md`).
- Enforcing the journal-stays-local boundary in code (latent
  follow-up from the decision record).
