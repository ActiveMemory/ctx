---
name: ctx-handover
description: Per-session handover artifact writer. Wraps `ctx handover write` with `--summary` and `--next` (both required, both validated non-placeholder by the CLI). Always invoked as the final step of `/ctx-wrap-up`; not the user-facing trigger. When `.context/kb/` exists, also folds postdated closeouts into the handover and archives them.
allowed-tools: Bash(ctx:*), Bash(git status), Bash(git log:*), Read
---

# Write a Handover

Capture the session's narrative thread so the next session (a
fresh agent, a different operator, a cold restart the next
morning) can resume without re-deriving context probabilistically
from canonical files plus journal.

This skill is the **sole authoritative recall artifact** writer
(per `KB-RULES.md` §Four inviolable rules: *"the handover is
the sole authoritative recall artifact"*). `SESSION_LOG.md`
entries, closeouts, and journal entries are mid-flight surfaces;
the handover is what `/ctx-remember` reads on session start.

Authoritative background reading:
`.context/ingest/KB-RULES.md` §Four inviolable rules;
`specs/kb-editorial-pipeline.md` §Interface.

## When to Use

`/ctx-wrap-up` owns the user-facing trigger for session-end
("let's wrap up", "save state", "leave a handover", "before I
go", "stepping away") and delegates to this skill as its final
step. Do not advertise this skill as a direct user trigger.

- **Mandatory tail of `/ctx-wrap-up`.** Every `/ctx-wrap-up`
  run ends with this skill.
- Mid-session checkpoint when the user wants to pause without
  consuming closeouts (use `--no-fold`). This is the one case
  where direct invocation is appropriate.

## When NOT to Use

- Nothing meaningful happened (only read files, quick lookup);
  but check with the user. A no-op session still benefits from
  a "nothing changed; next-step is X" handover when the next
  session has zero context.
- The user already ran `/ctx-handover` recently in this session
  and nothing has changed since.
- The user invokes a capture skill (`/ctx-task-add`,
  `/ctx-decision-add`, etc.); those write to canonical files,
  not to a handover artifact.

## Authority Boundary (vs Other Skills)

- **`/ctx-handover`**: writes
  `.context/handovers/<TS>-<slug>.md`; folds postdated
  closeouts from `.context/ingest/closeouts/` into the
  handover's `## Folded closeouts` section; archives folded
  closeouts to `.context/archive/closeouts/`. Single writer of
  this artifact.
- **`/ctx-wrap-up`**: owns the user-facing session-end
  trigger. Drives the broader capture ceremony (learnings,
  decisions, conventions, tasks) and always delegates to
  `/ctx-handover` as its final step.
- **`/ctx-remember`**: reads the latest handover plus any
  closeouts whose `generated-at` postdates the handover; the
  read-side counterpart to this skill's write surface.
- **Capture skills** (`/ctx-task-add`, `/ctx-decision-add`,
  `/ctx-learning-add`, `/ctx-convention-add`): write to the
  five canonical files. This skill never modifies those files;
  the handover narrative *references* them, it does not author
  them.

## Usage Examples

```text
/ctx-handover "kb editorial pipeline phase KB skills drafted"
/ctx-handover "rev2 spec landed; tomorrow start the writer package"
/ctx-handover "research session on cursor hooks"
/ctx-handover --no-fold "mid-session checkpoint before lunch"
```

## Input Contract

The skill wraps `ctx handover write`, which enforces required
flags via `MarkFlagRequired` and rejects placeholder bodies via
the Phase SK validation pattern. Empty `TBD`, `see chat`,
whitespace-only values are rejected by the CLI, not just by the
skill text.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--summary` | string | (required) | Past tense; what happened this session. |
| `--next` | string | (required) | Future tense; what the next agent should do FIRST. Specific, not vague. |
| `--highlights` | string | "" | Notable artifacts produced this session. |
| `--open-questions` | string | "" | Things that remain undecided. |
| `--no-fold` | bool | false | Skip closeout consumption (mid-session checkpoint). |
| `--commit` | string | (resolved) | Override resolved git HEAD for Provenance line (CI replay). |

Positional argument: handover title (becomes filename slug).

## Pre-Write Gates

Two distinct refusals, each leaves zero residue:

- `.context/` missing → suggest `ctx init` and stop.
- `.context/handovers/` missing → suggest `ctx init --upgrade`
  and stop.

`.context/kb/` is **not** required for handover; the artifact
exists for code-dev sessions as well. KB-state folding is
conditional on the directory's existence (see §Process).

## Process

1. **Verify pre-write gates.** Refuse cleanly if any gate
   fails. Zero residue on refusal.

2. **Gather signal silently** (mirror `/ctx-wrap-up` Phase 1
   when invoked standalone):

   ```bash
   git status --short
   git diff --stat
   git log --oneline @{upstream}..HEAD 2>/dev/null || git log --oneline -5
   ```

   Scan the conversation history for:
   - The session's arc: what shifted from start to now.
   - Concrete artifacts produced (files, commits, decisions,
     spec entries).
   - Open questions surfaced but not resolved.
   - The specific first action the next session should take.

3. **Draft `--summary` and `--next`.** Both are required, both
   are validated non-placeholder by the CLI:

   - **`--summary`**: past tense. One paragraph. Names what
     was done, not what was attempted. Concrete: *"drafted six
     Phase KB skill files; reconciled rev2 spec changes;
     deferred CLI wiring to next session"*, not *"made
     progress on KB stuff"*.
   - **`--next`**: future tense. One paragraph. Names the
     specific first action the next agent should take.
     Concrete: *"start `internal/cli/handover/cmd/write/cmd.go`
     using Phase SK validation pattern"*, not *"continue
     work" or "look at the kb"*.

   Surface the drafts to the user for confirmation before
   running the CLI. The user is the final authority on what
   the handover says.

4. **Run `ctx handover write`** with the confirmed values:

   ```bash
   ctx handover write "<title>" \
     --summary "<one-paragraph past tense>" \
     --next "<one-paragraph future tense>" \
     [--highlights "<bullet list>"] \
     [--open-questions "<bullet list>"] \
     [--no-fold] \
     [--commit <sha>]
   ```

   The CLI:
   - Validates flags (placeholder rejection per Phase SK).
   - Resolves git HEAD via `gitmeta.ResolveHead` (honors
     `CTX_TASK_COMMIT` and `GITHUB_SHA` for CI replay).
   - Reads `LatestHandoverCursor` to find the postdated
     closeout window.
   - Lists `UnconsumedCloseouts` (closeouts whose
     `generated-at` postdates the cursor).
   - For each unconsumed closeout, folds its body into the
     handover's `## Folded closeouts` section. Malformed
     closeouts (missing `generated-at`, malformed frontmatter)
     are skipped with a warning.
   - Calls `ArchiveCloseouts` to move folded closeouts to
     `.context/archive/closeouts/`. Archived closeouts are
     immutable.
   - Writes `.context/handovers/<TS>-<slug>.md`.

   When `--no-fold` is set, the fold + archive steps are
   skipped; closeouts stay in place. Use for mid-session
   checkpoints where the user wants the handover artifact but
   intends to keep ingesting before the next session boundary.

5. **Report the result.** Surface:
   - The handover filename written.
   - Count of closeouts folded (or *"none postdated the prior
     handover"*).
   - Count of malformed closeouts skipped (with filenames so
     the user can fix or delete; site-review's job to flag,
     but the warning here is opportunistic).
   - Any CLI validation failures (with the placeholder text
     that triggered rejection).

## Closeout Fold Mechanics

The fold mechanism is the integration point between the
editorial pipeline (`/ctx-kb-*` closeouts) and session continuity
(handover artifacts). Mechanically:

- `LatestHandoverCursor` reads `.context/handovers/` and returns
  the `generated-at` of the newest handover (or zero time if
  none exists).
- `UnconsumedCloseouts` walks `.context/ingest/closeouts/` and
  returns every closeout whose `generated-at` postdates the
  cursor.
- Each folded closeout's body is embedded under
  `## Folded closeouts` in the new handover, in `generated-at`
  order. The frontmatter is preserved verbatim so the audit
  trail survives the fold.
- After the fold, `ArchiveCloseouts` moves the source files to
  `.context/archive/closeouts/`. Archived closeouts are
  immutable; subsequent passes never re-fold them.

A handover with no postdated closeouts to fold writes a
`## Folded closeouts` section with the body *"none"*; never
omit the section, so the audit trail is explicit.

## Edge Cases

| Case | Expected behavior |
|------|-------------------|
| `.context/` missing | Refuse; suggest `ctx init`. No residue. |
| `.context/handovers/` missing | Refuse; suggest `ctx init --upgrade`. No residue. |
| Empty `--summary` or `--next` | The CLI rejects with the placeholder-rejection message; surface verbatim. |
| Placeholder values (`TBD`, `see chat`, whitespace-only) for `--summary` or `--next` | The CLI rejects; surface verbatim and ask the user to redraft. |
| No postdated closeouts to fold | Write the handover with `## Folded closeouts` body *"none"*. Never omit the section. |
| Postdated closeout has malformed frontmatter | The CLI skips the file with a warning naming it. Report the warning to the user so they can fix or delete. |
| `--no-fold` set | Skip the fold + archive steps; the handover stands alone; closeouts stay in `.context/ingest/closeouts/` for the next invocation. |
| Mid-session re-invocation | Each invocation writes a new handover file. The newest one is what `/ctx-remember` reads next session. Multiple per session are fine. |
| Session aborted before wrap-up | Closeouts stay in place; next session's `/ctx-remember` reads canonical files + the last handover + any postdated unfolded closeouts. Editorial work survives. |
| User runs `/ctx-wrap-up` without `.context/kb/` present | `/ctx-wrap-up` still drives `/ctx-handover` as its final step; kb-presence affects what gets folded, not whether the handover is written. |
| `gitmeta.ResolveHead` returns an error (no git, detached HEAD with no fallback) | The CLI surfaces the typed `MissingGitError`; relay verbatim. Phase RG owns the recovery path; this skill does not invent one. |
| `CTX_TASK_COMMIT` or `GITHUB_SHA` set | Honoured for the Provenance line per `gitmeta.ResolveHead`'s precedence rules; no special handling here. |

## Anti-Patterns

- Writing a handover with `--summary "TBD"` or `--next "see
  chat"`. The CLI rejects these; do not work around the
  rejection by inventing prose that technically passes the
  placeholder check but is still vague.
- Skipping the fold to *"keep closeouts available for a future
  pass"*. The fold is the integration point; closeouts that
  outlive their relevant handover are recall noise. Use
  `--no-fold` explicitly when the user wants the checkpoint
  behavior; do not infer it.
- Hand-writing a handover file. The CLI is the sole writer.
  Hand-edits drift from the schema the read side expects.
- Modifying an archived closeout. Archived closeouts are
  immutable per `KB-RULES.md` §Closeout shape.
- Inventing `--highlights` or `--open-questions` content the
  session did not actually produce. Light compression for
  clarity is allowed; new facts are not.

## Output Contract

For pre-write refusals, return only the specified refusal text
and stop.

For successful handover writes, end with this structured
summary:

- **Handover**: filename on its own line.
- **Folded closeouts**: count + filenames (or *"none
  postdated the prior handover"*).
- **Malformed skipped**: count + filenames (or `none`).
- **Provenance**: `sha=<short> branch=<name>` as resolved by
  the CLI.
- **Next-session focus**: the `--next` value, verbatim, so the
  operator sees what the next agent will read first.

## Quality Checklist

Before reporting completion, verify:

- [ ] Pre-write gates passed (or the matching refusal was
      returned with zero residue).
- [ ] `--summary` is past tense and concrete (no placeholder).
- [ ] `--next` is future tense and specific (no placeholder).
- [ ] User confirmed the drafts before the CLI ran.
- [ ] Closeouts were folded (or `--no-fold` was explicitly
      requested).
- [ ] Folded closeouts were archived to
      `.context/archive/closeouts/`.
- [ ] Handover filename + provenance were reported back to the
      user.
