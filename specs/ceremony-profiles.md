---
title: Ceremony profiles (project-level alias map for session ceremony nudges)
date: 2026-04-26
status: ready
---

# Ceremony Profiles

## Problem

The ceremony nudge hook (`internal/cli/system/core/ceremony/ceremony.go`,
config in `internal/config/ceremony/ceremony.go`) hardcodes two skill
names:

- `ctx-remember` — open-bookend ceremony
- `ctx-wrap-up` — close-bookend ceremony

These names are also baked into the message templates at
`internal/assets/hooks/messages/check-ceremony/{remember,wrapup,both}.txt`.

Non-code-dev projects rebrand those ceremonies. The motivating case is
the DR knowledgebase project at `things-wtf-disaster-recovery`, an
**editorial** project where:

- ctx code-dev skills explicitly do not apply (per its
  `.context/CONSTITUTION.md` and root `10-CONSTITUTION.md`).
- The project has shipped DR-specific replacements: `/dp-remember`,
  `/dp-wrap-up`, `/dp-commit`.

In that project the ceremony nudge fires correctly (3 sessions without
ceremony usage) but recommends the wrong skills. The user gets nudged
toward `/ctx-remember` / `/ctx-wrap-up`, which the project's
`CLAUDE.md` says **do not apply**. The nudge becomes cosmetic noise
the user has to mentally translate.

The leak is not specific to `/ctx-remember`: it's that the ceremony
layer assumes the code-dev ceremony set. Any non-code-dev profile
(editorial, research, ops, writing) hits the same wall.

## Solution

A project-level **ceremony profile**: an alias map declared in
`.ctxrc` that the ceremony scanner and message renderer consult
instead of the hardcoded constants.

Default profile preserves current behavior (`ctx-remember` /
`ctx-wrap-up`). Projects override by declaring a `ceremony:` block.

## Design

### Configuration surface

In `.ctxrc`:

```yaml
ceremony:
  remember: dp-remember
  wrapup: dp-wrap-up
  # commit: dp-commit   # optional, see "Commit ceremony" below
```

When the block is absent, the ceremony layer falls back to the
existing constants in `internal/config/ceremony/ceremony.go`. This
preserves backward compatibility — every existing project keeps
working unchanged.

A profile **may not** redefine other ctx command names; only the
ceremony bookend skills are aliasable. This keeps the abstraction
narrow and prevents projects from forking the entire ctx CLI surface.

### Scanner change

`ScanJournalsForCeremonies` in
`internal/cli/system/core/ceremony/ceremony.go` currently calls
`strings.Contains(content, ceremony.RememberCmd)` with a hardcoded
constant. Change it to receive the resolved command names from the
caller, e.g.:

```go
func ScanJournalsForCeremonies(files []string, names Names) (remember, wrapUp bool)
```

where `Names` is a struct populated by the rc loader (with defaults
applied when the project omits the block).

### Message templates

The three templates in
`internal/assets/hooks/messages/check-ceremony/`
(`remember.txt`, `wrapup.txt`, `both.txt`) currently embed literal
`/ctx-remember` and `/ctx-wrap-up` strings. They must be templated
with the resolved names.

Two viable approaches:

**Approach A — Go templates.** Parse the file as
`text/template`, render with `{{.Remember}}` / `{{.WrapUp}}`. Minor
parser cost on every nudge, but the existing `message.Load` path
already handles loading; only the rendering step is new.

**Approach B — String substitution.** Use a sentinel like `{REMEMBER}`
/ `{WRAPUP}` in the templates and `strings.Replace` at emit time.
Simpler, no template package overhead, no risk of accidental template
syntax in author-written content.

**Recommendation: B.** The templates are short and authored by us;
template syntax is overkill. Strict sentinels are easier to lint.

### Box title and fallback descs

`internal/assets/embed/text` (referenced via
`text.DescKeyCeremonyBoxBoth` etc.) holds the box titles and fallback
text. These also need the same sentinel-substitution treatment if
they include the literal command names. Audit them and apply the
same rule.

### Commit ceremony (separate, smaller question)

`/ctx-commit` is not currently part of the ceremony nudge — only
`remember` and `wrapup` are scanned. Adding `commit` to the alias
map is cheap (it's just another field), but only worth doing if we
also start nudging on missing commit ceremony usage. That's a
separate decision; keep the field reserved but undocumented for
now, or omit it entirely until needed. **Recommendation: omit until
we have a real nudge that needs it.** YAGNI.

### Profile inheritance / multi-project hub

Out of scope. Hub federation
(`specs/hub-federation.md`,
`specs/context-hub.md`) is its own track. If a hub later wants to
push profiles to member repos, it can do so by writing to `.ctxrc`
like any other config field. No special plumbing required here.

## CLI surface

No new commands. The change is purely in:

1. The `.ctxrc` schema (new `ceremony:` block).
2. The rc loader (parse the block, apply defaults).
3. The ceremony scanner (accept resolved names).
4. The message templates (sentinel substitution).

`ctx status` should show the active ceremony profile so users can see
what nudge they'll get without grepping `.ctxrc`. One line:

```
Ceremony: ctx-remember / ctx-wrap-up   (default)
```

or

```
Ceremony: dp-remember / dp-wrap-up   (project)
```

## Implementation plan

1. Add `Ceremony` struct to `internal/rc/types.go` with `Remember`
   and `WrapUp` fields. Defaults applied in `internal/rc/rc.go`
   loader from `internal/config/ceremony/ceremony.go` constants.
2. Change `ScanJournalsForCeremonies` to take resolved names rather
   than reading the constants directly.
3. Change `Emit` (or its caller) to thread the resolved names into
   the message render path.
4. Convert `check-ceremony/{remember,wrapup,both}.txt` to use
   `{REMEMBER}` / `{WRAPUP}` sentinels. Audit
   `internal/config/embed/text` for the same sentinel needs in box
   titles and fallback descs.
5. Add a sentinel-substitution helper (or extend
   `internal/cli/system/core/message.Load`) so substitution happens
   in one place, not per-call.
6. Add a one-line ceremony-profile readout to `ctx status`.
7. Update tests:
   - Default profile renders `/ctx-remember` / `/ctx-wrap-up` as
     before.
   - A project with `ceremony.remember: dp-remember` renders
     `/dp-remember` and the scanner only considers `dp-remember`
     usage as fulfilling the open-bookend.
8. Document in `docs/recipes/` (one short page on declaring a
   project ceremony profile, with the editorial-project example).

## Non-goals

- Aliasing arbitrary ctx commands. Only the bookend ceremonies.
- Profile inheritance across repos / hub-pushed profiles.
- Adding a `commit` ceremony nudge. The field is reserved-or-omitted;
  do not implement a commit-ceremony scanner as part of this work.
- Per-session profile switching. Profile is per-project, set in
  `.ctxrc`.

## Open questions

- **Migration**: should we deprecate `RememberCmd` / `WrapUpCmd` in
  `internal/config/ceremony/ceremony.go` once the rc-driven names
  land, or keep them as fallback constants forever? Recommendation:
  keep as the canonical default values the loader reads when the
  project omits the block. They stop being "the names" and start
  being "the default names."
- **Validation**: should the rc loader reject ceremony aliases that
  collide with built-in ctx command names (e.g.,
  `ceremony.remember: status`)? Probably yes — at least warn — but
  it's a small edge.
- **Empty string semantics**: `ceremony.remember: ""` could mean
  "this project has no open-bookend ceremony, suppress that half of
  the nudge." Clean way to silence one bookend without disabling the
  other. Worth supporting.

## Related

- DR editorial project (consumer):
  `~/Desktop/WORKSPACE/things-wtf-disaster-recovery`
- Existing ceremony plumbing:
  `internal/config/ceremony/ceremony.go`,
  `internal/cli/system/core/ceremony/ceremony.go`,
  `internal/assets/hooks/messages/check-ceremony/`
