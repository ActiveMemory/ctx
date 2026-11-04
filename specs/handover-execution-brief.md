# Handover Execution Brief Schema

## Problem

`ctx handover write` (`internal/cli/handover/cmd/write/cmd.go`) currently
requires only `--summary` (framed as a past-tense session narrative) and
`--next` (the next session's first action), alongside optional
`--highlights` and `--open-questions`.

Empirical evaluation in `arXiv:2611.05306` (*You don't need long context*,
Appendix A and Appendix G.1) shows that treating session compaction as a
narrative summary of past history is suboptimal for software engineering
workflows. Narrative summaries give equal weight to early exploration and
repeated trace segments, causing the next agent instance to:

1. Repeat dead-end approaches that already failed in the previous stretch.
2. Confuse unverified hypotheses with established facts observed in files
   or command outputs.
3. Re-run expensive verification suites or assume unverified changes were
   already tested.
4. Lose the overarching completion criteria across multiple compactions.

To serve as an accuracy-preserving Semantic State Space Model (SSSM)
transition artifact across both session boundaries and intra-task EoCC/SAC
resets, the handover must enforce a strict forward-looking execution brief
schema.

## Approach

Upgrade `ctx handover write` and `internal/write/handover` to require the
five core execution-brief fields derived from Appendix A and Appendix G.1
of `arXiv:2611.05306`:

1. **`--summary` (Required):** Objective and concise past-tense summary of
   what was accomplished in this stretch, retaining the verbatim or
   near-verbatim goal being solved.
2. **`--facts` (Required):** Established facts (exact file paths, relevant
   symbols, commands, observed values, and errors), explicitly separating
   verified observations from unverified assumptions.
3. **`--verified` (Required):** Applied changes (each created or modified
   file and what changed) together with exact verification checks run,
   their outcomes, obstacles, and workarounds.
4. **`--dead-ends` (Required):** Failed approaches and rejected paths that
   the next stretch must not repeat.
5. **`--next` (Required):** Remaining plan in execution order, the
   immediate next command or edit to execute first, and the concrete
   condition for completion (plus a narrow retrieval command if any needed
   output was truncated).

All five required flags are validated at the CLI boundary with
`validate.RejectPlaceholder`, which rejects empty strings, whitespace-only
strings, and placeholder tokens (`TBD`, `n/a`, `none`, `pending`, `see chat`,
and `.ctxrc` `placeholders` extensions). When a session legitimately has no
failed approaches or no file edits (for example, a read-only investigation),
the caller must supply a concrete descriptive sentence (such as
`"no failed approaches in this stretch"` or
`"read-only investigation; no files modified"`) rather than a bare
placeholder token.

Optional flags (`--highlights`, `--open-questions`, `--commit`, `--no-fold`)
remain unchanged for backward compatibility with existing editorial and KB
workflows.

## Behavior

### Happy Path

1. The agent (via `/ctx-handover`, `/ctx-wrap-up`, or an EoCC/SAC reset)
   gathers working-tree status, recent commits, test results, and session
   observations.
2. The agent drafts concrete values for `--summary`, `--facts`,
   `--verified`, `--dead-ends`, and `--next`, plus optional `--highlights`
   and `--open-questions`.
3. The CLI validates all five required flags through
   `validate.RejectPlaceholder`.
4. `handover.Write` composes the YAML frontmatter (`sha`, `branch`,
   `generated-at`, `title`) and renders the Markdown sections in canonical
   order:
   - `## Summary`
   - `## Established facts`
   - `## Applied changes and verification`
   - `## Dead ends`
   - `## Next session`
   - `## Highlights` (when non-empty)
   - `## Open questions` (when non-empty)
   - `## Folded closeouts` (when closeouts are folded)
5. The handover file is written to `.context/handovers/<TS>-<slug>.md`.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Any of `--summary`, `--facts`, `--verified`, `--dead-ends`, or `--next` is omitted | Cobra rejects the invocation naming the missing required flag(s). Zero files written. |
| Caller passes `"none"`, `"N/A"`, `"TBD"`, or whitespace to `--dead-ends` or `--verified` | `validate.RejectPlaceholder` rejects the flag with `errCli.FlagPlaceholder` or `errCli.FlagEmpty`. Zero residue on disk. |
| Read-only research session with no file edits or failed approaches | Caller passes concrete prose such as `--verified "read-only session; no code modified"` and `--dead-ends "no failed approaches encountered"`. Validation succeeds. |
| Existing handover files on disk written under the legacy 2-section schema | Readers (`/ctx-remember`, `ctx agent`, `handover.Latest`) parse legacy files without error; missing new sections are treated as absent rather than malformed. |
| `--no-fold` passed during mid-session EoCC compaction | Closeouts in `.context/ingest/closeouts/` are left untouched on disk while all five execution-brief sections are written to the new handover file. |

### Validation Rules

- Enforced in `internal/cli/handover/cmd/write/cmd.go` via
  `MarkFlagRequired` for `cFlag.Summary`, `cFlag.Facts`, `cFlag.Verified`,
  `cFlag.DeadEnds`, and `cFlag.Next`.
- Enforced in `internal/cli/handover/cmd/write/run.go` via a loop over the
  five `(flag, value)` pairs calling `validate.RejectPlaceholder(flag, val)`.
- Enforced at the library boundary in `internal/write/handover/write.go` via
  non-empty checks returning `errHandover.ErrSummaryRequired`,
  `errHandover.ErrFactsRequired`, `errHandover.ErrVerifiedRequired`,
  `errHandover.ErrDeadEndsRequired`, and `errHandover.ErrNextRequired`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| Missing required flag | `Error: required flag(s) "dead-ends", "facts", "verified" not set` | Supply concrete values for every required execution-brief flag. |
| Placeholder value supplied | `Error: flag --dead-ends cannot be placeholder value "none"` | Replace the placeholder with a specific description (e.g., `"no failed approaches in this session"`). |
| `.context/handovers/` missing | Existing `errHandover` directory refusal message | Run `ctx init --upgrade` to scaffold `.context/handovers/`. |

## Interface

### CLI

```bash
ctx handover write <title> \
  --summary "<objective and completed work>" \
  --facts "<verified observations vs assumptions>" \
  --verified "<applied file changes and verification outcomes>" \
  --dead-ends "<failed approaches not to repeat>" \
  --next "<immediate next command/edit, remaining subgoals, completion condition>" \
  [--highlights "<notable artifacts>"] \
  [--open-questions "<undecided items>"] \
  [--no-fold] \
  [--commit <sha>]
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--summary` | `string` | (required) | Objective and concise past-tense summary of the stretch. |
| `--facts` | `string` | (required) | Established facts (paths, symbols, observed values) separated from assumptions. |
| `--verified` | `string` | (required) | Applied file changes and exact verification checks run with outcomes. |
| `--dead-ends` | `string` | (required) | Failed approaches and obstacles that must not be repeated. |
| `--next` | `string` | (required) | Immediate next command or edit, ordered subgoals, and completion condition. |
| `--highlights` | `string` | `""` | Optional notable artifacts produced this session. |
| `--open-questions` | `string` | `""` | Optional unresolved questions. |
| `--no-fold` | `bool` | `false` | Skip closeout consumption (used for mid-session EoCC checkpoints). |
| `--commit` | `string` | (resolved) | Override resolved git HEAD SHA. |

### Skill

```text
/ctx-handover
/ctx-wrap-up
```

Both skills and their mirror copies under `internal/assets/codex/`,
`internal/assets/integrations/copilot-cli/`, and
`internal/assets/integrations/opencode/` must prompt the agent to draft and
confirm all five required fields following the Appendix G.1 checklist.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/config/flag/flag.go` | Add flag name constants `Facts = "facts"`, `Verified = "verified"`, and `DeadEnds = "dead-ends"`. |
| `internal/config/embed/flag/handover.go` | Add description keys `DescKeyHandoverFacts`, `DescKeyHandoverVerified`, and `DescKeyHandoverDeadEnds`. |
| `internal/assets/commands/flags.yaml` | Add user-facing descriptions for `handover.facts`, `handover.verified`, and `handover.dead-ends`. |
| `internal/assets/commands/commands.yaml` | Update `handover` and `handover.write` examples to include `--facts`, `--verified`, and `--dead-ends`. |
| `internal/config/handover/handover.go` | Add section header constants `SectionFacts = "## Established facts"`, `SectionVerified = "## Applied changes and verification"`, and `SectionDeadEnds = "## Dead ends"`. |
| `internal/config/embed/text/err_handover.go` | Add `DescKeyErrHandoverFactsRequired`, `DescKeyErrHandoverVerifiedRequired`, and `DescKeyErrHandoverDeadEndsRequired`. |
| `internal/assets/commands/text/errors.yaml` | Add error messages for empty `facts`, `verified`, and `dead-ends`. |
| `internal/err/handover/handover.go` | Add sentinels `ErrFactsRequired`, `ErrVerifiedRequired`, and `ErrDeadEndsRequired`. |
| `internal/write/handover/types.go` | Add `Facts`, `Verified`, and `DeadEnds` string fields to `handover.Entry`. |
| `internal/write/handover/write.go` | Validate non-empty `Facts`, `Verified`, and `DeadEnds` on `Entry`. |
| `internal/write/handover/markdown.go` | Render `SectionFacts`, `SectionVerified`, and `SectionDeadEnds` between `SectionSummary` and `SectionNext`. |
| `internal/cli/handover/cmd/write/cmd.go` | Bind `--facts`, `--verified`, and `--dead-ends` and mark all five brief flags required. |
| `internal/cli/handover/cmd/write/run.go` | Accept the three new parameters and validate all five required flags with `validate.RejectPlaceholder`. |
| `internal/cli/handover/cmd/write/doc.go` | Document the five required execution-brief flags. |
| `internal/write/handover/handover_test.go` | Update existing tests and add section-ordering and required-field rejection tests for the new fields. |
| `internal/assets/claude/skills/ctx-handover/SKILL.md` | Update input contract, process steps, and examples to require `--facts`, `--verified`, and `--dead-ends` (and sync to `internal/assets/codex/skills/`, `internal/assets/integrations/copilot-cli/skills/`, and `internal/assets/integrations/opencode/skills/`). |
| `internal/assets/claude/skills/ctx-wrap-up/SKILL.md` | Update Phase 5 handover delegation instructions to draft all five brief fields (and sync mirrored skill trees). |
| `docs/cli/handover.md` | Document the new required flags and the forward-looking execution brief structure. |

### Key Functions

```go
func Run(
	cobraCmd *cobra.Command,
	title, summary, facts, verified, deadEnds, next string,
	highlights, openQuestions, commit string,
	noFold bool,
) error
```

### Helpers to Reuse

- `validate.RejectPlaceholder` (`internal/validate/rejectplaceholder.go`)
- `flagbind.StringFlag` (`internal/flagbind`)
- `handover.Write` (`internal/write/handover/write.go`)

## Configuration

No new `.ctxrc` keys. Existing `.ctxrc` `placeholders:` extensions
automatically apply to `--facts`, `--verified`, and `--dead-ends` via
`validate.RejectPlaceholder`.

## Testing

- **Unit (`internal/write/handover/handover_test.go`):** Verify
  `renderBody` emits `## Summary`, `## Established facts`,
  `## Applied changes and verification`, `## Dead ends`, and
  `## Next session` in exact order, followed by optional sections.
- **CLI validation (`internal/cli/handover/cmd/write`):** Verify omitting
  any of the five flags errors, passing `"none"` or `"TBD"` to
  `--dead-ends`, `--facts`, or `--verified` errors with the flag name, and
  passing concrete prose succeeds.
- **Skill sync compliance:** Verify `make check-codex-skills` and
  `make check-copilot-skills` pass after updating skill files.

## Non-Goals

- Rejecting or migrating historical handover files in `.context/handovers/`
  that were written with only `## Summary` and `## Next session`.
- Replacing `.context/TASKS.md` or `.context/DECISIONS.md` with handover
  files; canonical files remain the long-term project memory while
  handovers carry the immediate execution state between stretches.
