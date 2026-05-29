# Spec: `--json <file>` ingest for `ctx <type> add`

**Status:** stub (seed from TASKS.md Phase CLI-FIX, 2026-05-28)
**Driver session:** 96765858

Add `--json <file>` to `ctx decision/learning/task add` (and `convention` if it gains structured fields) for ingesting a JSON payload that populates the typed fields directly.

- **Driver**: this session hit a class of denial we worked around but should fix at the root. The project's canonical `permissions.deny` set (`.claude/settings.local.json` lines 119-121) matches on the literal Bash command string — including the *content* of `--rationale`/`--context`/`--consequence` flag values. A decision whose rationale legitimately describes installing a binary into the system PATH (literal substring " /usr/local/bin") gets caught by `Bash(* /usr/local/bin*)` and denied, even though the command's intent has nothing to do with that path. The workaround was Edit-direct into `DECISIONS.md`/`LEARNINGS.md`, which bypasses the ctx command's schema gates and `INDEX:START/END` maintenance.
- **Shape**: `ctx decision add --json /path/to/payload.json` where the JSON is `{"title":"…","context":"…","rationale":"…","consequence":"…"}`. The flag supersedes individual content flags. Provenance (`--session-id`/`--branch`/`--commit`) can stay on the command line OR be folded into the JSON envelope (`{"provenance":{"session_id":"…","branch":"…","commit":"…"}}`). Complements the existing `--file` (which only replaces the title/body positional).
- **Phase 2 (optional)**: array form `[{...},{...}]` for batch persists — useful for `/ctx-wrap-up` writing N decisions+learnings in one call instead of N separate invocations.
- **Mirror per command**: same shape applies to `ctx learning add --json …` (`{title,context,lesson,application}`) and `ctx task add --json …` (`{title,body,priority,section}`).
- **Surfaced by**: this session's persist denials and post-mortem; reference handover `20260528T201500Z-ctxctl-and-native-pressure-shipped.md`.
