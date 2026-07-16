# Spec: ceremony nudge — live-session credit and self-suppression

## Problem

`ctx system check-ceremonies` (the UserPromptSubmit hook wired as
`checkceremony`) nudges the operator to open sessions with
`/ctx-remember` and close them with `/ctx-wrap-up`. It decides whether
to nudge by scanning the most recent *imported* journal files
(`ScanJournalsForCeremonies` over `RecentJournalFiles`). Two defects
follow from that journal-only signal:

1. **Self-nudge (task :2634).** The hook fires on every
   UserPromptSubmit, including the prompt that *is* `/ctx-remember`.
   Because `entity.HookInput` parses only `session_id` and
   `tool_input.command`, the hook cannot see the prompt text and so
   cannot tell that the current prompt is the ceremony itself. Result:
   running `/ctx-remember` is answered with "try starting with
   `/ctx-remember`."

2. **Journal-import lag (task :2632).** The live session's
   `/ctx-remember` is not in any imported journal until the journal is
   imported later. Until then the scan reports the ceremony as missing
   and the hook keeps nudging within the very session that already ran
   it.

Both reduce to the same root: the check has no view of the *current*
prompt, only of past, imported journals.

## Design

Parse the UserPromptSubmit `prompt` field and treat "the current prompt
is a ceremony command" as an authoritative live signal.

- Add `Prompt string` (`json:"prompt"`) to `entity.HookInput`.
  `session.ReadInput` already unmarshals the whole envelope, so the
  field populates with no call-site change.

- Add `ceremony.InvokedByPrompt(prompt string) bool` to the ceremony
  core: true when the prompt's first whitespace-delimited token is a
  ceremony slash command in either the bare (`/ctx-remember`) or
  plugin-scoped (`/ctx:ctx-remember`) form, for both `ctx-remember` and
  `ctx-wrap-up`. First-token equality (not substring/prefix) avoids
  matching `/ctx-remembering` or a command merely named in prose.

- In `checkceremony.Run`, after the pause/init preamble and the daily
  throttle check, add one guard: when `InvokedByPrompt(input.Prompt)`
  is true, touch the daily marker and return without nudging. Touching
  the marker credits the live ceremony for the rest of the day (fixes
  the import lag); returning without a nudge suppresses the self-nudge.

The marker is the existing per-day throttle (`ceremony-reminded` under
the state dir), so crediting a live ceremony simply means "the ceremony
question is settled for today" — consistent with the check's existing
once-per-day cadence.

### Trade-off

Touching the day marker on a live `/ctx-remember` also suppresses a
would-be `/ctx-wrap-up` nudge for the rest of that day. This is
deliberate: the check is a coarse daily cadence, and an operator who is
actively running ceremonies does not need to be nagged about the other
one on the same prompt. A finer per-ceremony credit is not worth the
added state.

## Implementation

- `internal/entity/hook.go`: add `Prompt` field + doc.
- `internal/config/ceremony/ceremony.go`: add `SlashPrefix` (`/`) and
  `PluginSlashPrefix` (`/ctx:`) constants.
- `internal/cli/system/core/ceremony/ceremony.go`: add
  `InvokedByPrompt`.
- `internal/cli/system/cmd/checkceremony/run.go`: add the live-ceremony
  guard.

## Tests

- `ceremony_test.go`: table test for `InvokedByPrompt` — bare and
  plugin forms of both commands (true), leading whitespace and trailing
  args (true), non-ceremony prompt, empty prompt, and the boundary case
  `/ctx-remembering` (false).
- `checkceremony/run_test.go`: with an initialized temp project, a
  `/ctx-remember` prompt makes `Run` return nil, emit no nudge, and
  create the `ceremony-reminded` marker (live-credit + self-suppress).

## Acceptance

- Running `/ctx-remember` (bare or `/ctx:` form) never produces a
  ceremony nudge on that prompt.
- After a live `/ctx-remember`, the day marker exists, so subsequent
  prompts the same day are not nudged even before journal import.
- Non-ceremony prompts are unaffected: the journal-driven nudge still
  fires when a ceremony is genuinely missing.

Closes TASKS :2632 (live-credit) and :2634 (self-suppress).
