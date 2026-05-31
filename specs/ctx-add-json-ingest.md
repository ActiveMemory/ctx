# Spec: `--json-file <path>` ingest for `ctx <type> add`

**Status:** accepted (impl 2026-05-30)
**Driver session:** 96765858

Add `--json-file <path>` to `ctx decision/learning/task/convention add`
for ingesting a JSON payload that populates the typed fields directly,
keeping flag-value *content* off the literal Bash command string.

## Driver

This class of denial should be fixed at the root. The project's canonical
`permissions.deny` set (`.claude/settings.local.json`) matches on the
literal Bash command string — including the *content* of
`--rationale`/`--context`/`--consequence` flag values. A decision whose
rationale legitimately describes installing a binary into the system PATH
(literal substring `" /usr/local/bin"`) is caught by
`Bash(* /usr/local/bin*)` and denied, even though the command's intent has
nothing to do with that path. The workaround was Edit-direct into
`DECISIONS.md`/`LEARNINGS.md`, which bypasses the ctx command's schema
gates and `INDEX:START/END` maintenance. Moving the values into a JSON file
keeps them out of the command string entirely.

## Flag name

`--json-file <path>`, **not** `--json`. Across the rest of the CLI
(`ctx status --json`, `ctx drift --json`, `ctx doctor --json`, …) `--json`
is a *bool* flag meaning "format output as JSON". Overloading it as a
*string input-path* flag on the add commands would break that convention,
so the input-payload flag is named `--json-file` (long-only; `-j` is
already `ShortJSON`). Parallels the existing `--file`/`-f` source flag.

## Payload shape

A single JSON object. All keys optional; only the ones relevant to the
noun are consumed:

```json
{
  "title": "…",
  "body": "…",
  "context": "…",
  "rationale": "…",
  "consequence": "…",
  "lesson": "…",
  "application": "…",
  "priority": "…",
  "section": "…",
  "provenance": { "session_id": "…", "branch": "…", "commit": "…" }
}
```

Per-noun field consumption (extra keys are tolerated but ignored by the
noun's formatter):

| Noun       | content   | typed fields                       |
|------------|-----------|------------------------------------|
| decision   | `title`   | `context`, `rationale`, `consequence` |
| learning   | `title`   | `context`, `lesson`, `application` |
| task       | `title` (+ `body`) | `priority`, `section`     |
| convention | `title`   | —                                  |

- **Content** is `title`; for `task`, a non-empty `body` is appended to
  `title` (space-joined) since `TASKS.md` entries are single-line.
- **Provenance** (`session_id`/`branch`/`commit`) may stay on the command
  line OR be folded into the `provenance` envelope.
- Decoding is **strict** (`DisallowUnknownFields`): a misspelled key is a
  hard error so typos surface instead of silently dropping a field.

## Precedence

`--json-file` **supersedes** the individual content/typed flags: any
non-empty payload field overrides the corresponding CLI flag. Empty or
absent payload fields leave the CLI flag value intact, so a caller may mix
CLI flags and a partial JSON envelope. For content, `--json-file` outranks
`--file`, positional args, and stdin.

## Implementation

The overlay runs in two visible touchpoints; both load the (tiny) file via
the shared `jsonpayload.Load`:

1. **Typed fields → cobra flags, in `PreRunE`.** `jsonpayload.OverlayFlags`
   reads `--json-file`, loads the payload, and `flags.Set`s each non-empty
   mapped field that exists on the command (context, rationale,
   consequence, lesson, application, priority, section, session-id, branch,
   commit). This must happen in `PreRunE` because the decision/learning
   `PreRunE` placeholder gate (`validate.RejectPlaceholder`) validates the
   *effective* flag values — so JSON-supplied values are placeholder-checked
   too, closing the bypass hole. `decision`/`learning`/`task` call
   `OverlayFlags` first in their `PreRunE`; `convention` has no typed fields
   and needs no overlay.
2. **Content → positional, in `extract.Content`.** When `flags.JSONFile`
   is set and the payload yields non-empty content, that content wins;
   otherwise extraction falls through to `--file`/args/stdin.

`run.Run` is unchanged: the bound `--*` variables already reflect the
`flags.Set` overlay, so the assembled `entity.AddConfig` carries the
effective values.

## Out of scope

- **Phase 2 (batch array form `[{…},{…}]`)** for N entries in one call
  (useful for `/ctx-wrap-up`). Deferred; the object form lands first.

## Surfaced by

This session's persist denials and post-mortem; reference handover
`20260528T201500Z-ctxctl-and-native-pressure-shipped.md`.
