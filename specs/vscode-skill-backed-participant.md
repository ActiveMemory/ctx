# VS Code `@ctx` Participant: Skill-Backed, Verified Against the Real CLI

Issue: https://github.com/ActiveMemory/ctx/issues/127
Supersedes: https://github.com/ActiveMemory/ctx/pull/128 (closed for rework)

The `@ctx` chat participant on `main` exposes only CLI wrappers, none of
the skill/workflow layer (#127). PR #128 added the workflow commands
but was closed: a dozen of its commands dispatched to `ctx` subcommands
that do not exist, failures rendered as results, and nothing tested the
dispatched argv against the real binary. This spec lands the participant
on current `main`, with every command traceable to a real `ctx` command
or a shipped skill, and a test that fails when that stops being true.

## Problem

1. **`main`'s participant is dead on arrival.** Every invocation passes
   `--no-color`, a flag the CLI no longer has, so every command exits 1
   with a usage dump. Beyond that: `recall list`, `add <type>`,
   `notify ...`, `system resources|message`, `pause`, `resume`,
   `reindex`, `prompt`, `dep`, and `loop <tool>` no longer match the
   CLI.
2. **Failures read as results.** `runCtx` resolved on any output, so an
   error (or a timed-out partial run) rendered like success.
3. **Windows command injection.** `execFile(..., { shell: true })`
   joins arguments unquoted: a multi-word prompt splits into extra
   arguments and `&` in a prompt runs a second command.
4. **The workflow layer is missing** (#127), and #128's version of it
   re-implemented skills in TypeScript against file formats ctx no longer
   uses (`IMPLEMENTATION_PLAN.md`, `- ` bullets in DECISIONS.md), so
   those commands were dead too, just silently.
5. **Nothing guards the surface.** Unit tests mock `execFile`; any
   argv passes.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| What "skill-backed" means | `/<name>` sends the canonical `ctx-<name>` SKILL.md, `ctx agent` output, skill-specific read-only ctx output, earlier turns of the same skill, and `#file` attachments to the chat model (`request.model`) | Delegates the workflow to the skill ctx ships instead of a TypeScript imitation that drifts from it. |
| Skill source | Bundled at build time from `internal/assets/claude/skills/` (esbuild `--loader:.md=text`) | Pins skill text to the commit the VSIX is built from; a renamed skill breaks the build. No CLI command prints skill bodies, and none is added for this. |
| Which skills | brainstorm, spec, implement, next, remember, reflect, wrap-up, blog, consolidate | They work from the inputs the participant can supply. The model here cannot run commands or browse the repo, so architecture (the old `/map`), link-check, worktree, and blog-changelog stay with agent integrations. `/audit` and `/verify` have no shipped skill. |
| Command names | Skill name without `ctx-` (`/wrap-up`, not `/wrapup`); CLI nouns singular (`/task`, `/change`, `/permission`, `/decision`, `/learning`) | One obvious mapping each way; matches the command names already on `main`. |
| Model limits | The preamble forbids claiming to have run or changed anything; the model gives the exact `ctx`/`@ctx` command instead | Honest about the participant's reach. |
| History sent to the model | Only earlier skill exchanges | `/pad`, `/notify`, `/add` turns can hold secrets; CLI output never reaches the model. |
| Process execution | No shell; stdin closed; resolve with the exit code; reject on spawn/buffer failure, cancel, timeout | Fixes injection and splitting; a prompting command (`ctx why`, content-less `add`) fails fast; failures cannot pose as results. |
| Rendering | One `runAndRender`: non-zero exit shows "exited with code N" and the CLI output, plus an `/init` hint when `.context/` is missing | Every handler routes through it, so none can render a failure as success. |
| Init gate | Removed | The CLI owns which commands need `.context/` (`AnnotationSkipInit`); a copy in the extension drifts (it blocked `/guide`, `/why`). |
| `/add` defaults | Provenance filled in (`vscode.env.sessionId`, git branch and commit); no `--section` default | The CLI requires provenance; it deliberately refuses a catch-all section (`internal/cli/add/core/build/section.go`). |
| `/notify setup` | Tells the user to run `ctx hook notify setup` in a terminal | The command prompts for a secret webhook URL; it must not go through chat. |
| Natural language | Routes only to read-only commands; CLI targets run without arguments | A keyword match must never complete a task or add a reminder. |
| Background hooks from #128 | Kept: session start/end, reminder bar (read-only `remind list`). Dropped: save watcher, commit popup, dependency popup, heartbeat file, violation recording | `system check-task-completion` reads hook JSON from stdin and exits silently without it, so the save watcher never did anything; the commit popup fired on every HEAD move (checkout, pull); the rest had no ctx consumer. Violation capture is out of scope (see #128 review). |
| Version | `0.10.0`; new CHANGELOG section, `0.9.0` untouched | Review blocker 4. |

## Command Surface

CLI-backed (27): `/init`, `/status`, `/agent`, `/drift`, `/recall`,
`/setup`, `/add`, `/decision`, `/learning`, `/load`, `/compact`,
`/sync`, `/task`, `/remind`, `/pad`, `/notify`, `/system`, `/memory`,
`/journal`, `/doctor`, `/config`, `/why`, `/change`, `/guide`,
`/permission`, `/pause`, `/resume`. The argv each one runs is recorded
in `editors/vscode/src/ctx-cli-surface.json`.

Skill-backed (9): `/brainstorm`, `/spec`, `/implement`, `/next`,
`/remember`, `/reflect`, `/wrap-up`, `/blog`, `/consolidate`.

Removed from `main`'s surface: `/prompt`, `/dep`, `/reindex` (CLI
commands gone), `/loop` (writes a terminal-agent script), `/site`
(hidden maintainer command).

## Parity Guard

Two halves, one artifact:

1. **vitest** (`editors/vscode/src/commandParity.test.ts`):
   - `package.json` commands == `CLI_COMMANDS` ∪ `SKILLS`;
   - each skill command bundles the skill it names (`ctx-<command>`,
     text byte-equal to the SKILL.md);
   - every follow-up and natural-language route targets a real command;
   - every command has a scenario; each scenario runs through the real
     chat handler with `execFile` mocked, and every `ctx` argv produced
     (plus background invocations) is written to
     `src/ctx-cli-surface.json` via `toMatchFileSnapshot`. An unreviewed
     change to what the extension runs fails the test.
2. **Go** (`internal/bootstrap/vscode_surface_test.go`, in the main
   `go test ./...` job), for every argv in that file:
   - `Find` in a fresh `Initialize(RootCmd())` tree: unknown command,
     deprecated command, or group-only command fails;
   - `ParseFlags`, `ValidateArgs`, `ValidateRequiredFlags`,
     `ValidateFlagGroups`; a group or no-argument command with leftover
     positionals fails (cobra accepts those silently, which is how
     `ctx system stats` "worked");
   - `task|decision|learning|convention add` invocations run in a
     scratch project, because their required fields are checked in
     `RunE`;
   - every listed skill must exist in `internal/assets/claude/skills`.

A CLI rename therefore fails the Go job; an extension change that
alters an argv fails vitest until the snapshot is refreshed and the Go
job re-validates it.

**Refreshing the snapshot**: from `editors/vscode`, run
`npx vitest run -u`, review `git diff src/ctx-cli-surface.json`, then
`go test ./internal/bootstrap -run TestVSCode`.

## Review Points from #128

| Point | Resolution |
|-------|------------|
| B1: commands call non-existent subcommands | Reconciled against the current tree; guarded by the Go test. |
| B2: `runCtx` renders failures as success | Exit code surfaced; `runAndRender` shows failures; cancel/timeout/spawn reject. |
| B3: no test for the new surface | `commandParity.test.ts` + snapshot + Go test; handler tests cover every branch. |
| B4: no version bump | `0.10.0`, new CHANGELOG section. |
| Reminder bell permanent | Reads `ctx remind list`; bar refresh waits for bootstrap (it used to run before it and never updated at activation). |
| `/pause` `/resume` repurposed | `ctx hook pause` / `ctx hook resume`. |
| Pre-init gate blocks `/guide`, `/why` | Gate removed; the CLI decides. |
| Unreachable palette registrations | None registered. |
| `/verify`, `/wrapup` oversell | `/verify` dropped (no shipped skill, `/doctor` + `/drift` cover it); `/wrap-up` runs the real skill and says it writes nothing. |
| Unkillable git handlers | Git runs only for provenance, with timeout and cancellation. |
| `saveWatcher` cross-root misfire | Watcher removed (it never produced a nudge). |
| Guardrails (terminal capture, violations.json) | Not included; separate proposal if at all. |
| Windows `shell: true` injection | Fixed: no shell. |
| Multi-root `folder[0]` | Active editor's folder, else the first. |

## Verification

- `editors/vscode`: `npm ci`, `npm run build`,
  `npx tsc --noEmit -p tsconfig.ci.json`, `npm run lint`, `npm test`,
  `npx vsce package --no-dependencies`.
- `go test ./internal/bootstrap -run TestVSCode`; full `go test ./...`
  shows no new failures; `golangci-lint run` clean.
- The recorded argv were also executed against a freshly built `ctx` in
  a scratch project, and the bundled `dist/extension.js` was driven end
  to end against it with a stub `vscode` module.

## Out of Scope

- Letting the model call tools (`vscode.lm.invokeTool`) so skills can
  explore the workspace; that would bring the excluded skills in.
- A Command Palette surface.
- Terminal-command and sensitive-file capture (#128 guardrails).
