# Spec: MCP entry tools accept provenance

Issue: ActiveMemory/ctx#190.

`ctx_add` and `ctx_watch_update` cannot add a task, decision, or
learning for any schema-driven MCP client in a project without a
relaxing `.ctxrc`: entry validation demands fields the tool schemas
never advertise.

## Problem

### A. Provenance is required but not advertised

`entry.Validate` (`internal/entry/validate.go:44-58`) requires
`session_id`, `branch`, and `commit` for tasks, decisions, and
learnings unless `.ctxrc` `provenance_required` relaxes them; the
`rc.Provenance*Required` accessors default to true. The MCP layer
reads the three values only from tool arguments (`extract.Opts`), but
`EntryAttrProps` (`internal/mcp/server/def/tool/prop.go`) never
declares them. MCP clients send only what a schema advertises, so
every such call fails:

    learning missing required fields: --session-id, --branch, --commit

`ctx_watch_update` shares `EntryAttrProps` and the same validation,
and fails identically.

### B. `ctx_watch_update` does not advertise `section` or `priority`

Validation also requires `section` for tasks (`validate.go:62-65`).
`ctx_add` advertises `section` and `priority`; the `ctx_watch_update`
schema (`internal/mcp/server/def/tool/tool.go`) declares neither,
although its handler reads both. `ctx_watch_update type=task` fails
the same way for schema-driven clients: same defect class as A.

### C. What an MCP client can supply is weak, absent, or unsafe

- `session_id`: only the client can know it. Claude Code's hook relay
  puts the real ID in the agent's context on every prompt
  (`Session: … | Branch: … @ …`); clients without a relay have none.
  The server has no session identity to fall back on: `MCPSession`
  holds governance counters only, and `ctx_sessionevent start` merely
  resets them (`internal/mcp/handler/tool.go:531-533`).
- `branch` / `commit`: the server can resolve them from git itself,
  and its values are better than relayed ones: read at write time (a
  relayed commit predates any commit the agent made since the prompt)
  and never retyped by a model.
- `SanitizedOpts` (`internal/mcp/server/extract/extract.go`) promises
  sanitization for every text field but leaves `priority`, `section`,
  `branch`, and `commit` raw. `format.Task` writes priority, branch,
  and commit into the task line, and `insert.TaskAfterSection` writes
  an unmatched section as a new `### ` heading, so a line break in any
  of them forges task lines or headings in `TASKS.md`. Advertising
  `branch` and `commit` makes schema-driven clients send them.

## Approach

MCP layer only. The CLI, `entry.Validate`, and `.ctxrc` are unchanged.

1. **Advertise** `section`, `priority`, `session_id`, `branch`, and
   `commit` in `EntryAttrProps`, so both entry tools declare every
   field validation can require from one definition.
2. **Resolve the git half server-side.** In the route layer
   (`internal/mcp/server/route/tool/tool.go`), `add` and `watchUpdate`
   build their options through one helper that wraps
   `extract.SanitizedOpts` and fills `Branch` and `Commit` when they
   are empty or `unknown`: `git.CurrentBranch()` and `git.ShortHead()`,
   falling back to `unknown`. These are the values the hook relay
   emits (`provenance.Emit`), so MCP-written provenance matches the
   CLI-plus-relay path, including `unknown` when git cannot answer.
   The fallback is `cmp.Or`; importing `provenance.DefaultVal` would
   pull cobra into the MCP route for a one-line default. An explicit
   value wins, for an agent working in another checkout than the
   server's cwd. No cache: caching would freeze `commit` at first use,
   the staleness this step removes.
3. **Flatten single-line fields.** A new `sanitize.Line(s, maxLen)`
   strips every control character, including the tab, CR, and LF that
   `StripControl` keeps for multi-line content, and caps the length on
   a rune boundary. `SanitizedOpts` applies it to `priority`,
   `section`, `branch`, and `commit`, capped by a new
   `cfg.MaxOptsLineLen` (256) in the MCP-SAN.1 limits block. No
   legitimate value contains a control character, so none changes.
   Over-long values are truncated, not rejected, as `sanitize.SessionID`
   does. `sanitize.Reflect` does not fit: it keeps line breaks. The
   cap is the one `specs/sanitize-hardening-followup.md` deferred for
   these fields ("tracked separately if needed"); advertising `branch`
   and `commit` makes it needed.
4. **`session_id` stays caller-supplied and required.** Its schema
   description directs the agent to the relayed `Session:` value, or an
   explicit `"unknown"` when it has none: no silent default and no
   fabricated ID. This is not a relaxation: omitting the field still
   fails, and only `.ctxrc` relaxes the gate
   (`docs/recipes/knowledge-capture.md`). `"unknown"` is what the relay
   itself emits without a session ID (`provenance.ShortSessionID`), and
   the CLI already accepts it.
5. **Describe** `branch` and `commit` as optional over MCP: omit them
   to let the server resolve them from git.
6. **Prompts.** The `ctx-decision-add` and `ctx-learning-add` footers
   say "call ctx_add with … all fields above", and `session_id` is not
   among those fields; the footers name it.
7. **Document** the missing argument rows for both tools in
   `docs/cli/mcp.md`, and note in `docs/home/configuration.md` that
   `provenance_required` governs the MCP entry tools too.

| Caller sends                     | Task line stores                                   |
|----------------------------------|----------------------------------------------------|
| all three values                 | the values as sent, flattened                      |
| `session_id` only                | that `session_id` plus the server's branch/commit  |
| `branch: "unknown"`              | the server's branch                                |
| `branch: "x\n- [x] forged"`      | `#branch:x- [x] forged`, on one line               |
| no `session_id`                  | nothing: `missing required fields: --session-id`   |
| git absent, or no commits yet    | `unknown` for whichever git cannot resolve         |
| detached HEAD                    | `#branch:unknown`, as the relay emits              |

## Tests

Each test below fails before its fix, except the explicit-values
case, which guards the "explicit value wins" rule:

- `internal/mcp/server/def/tool`: both entry tools advertise
  `section`, `priority`, `session_id`, `branch`, `commit`, `context`,
  `rationale`, `consequence`, `lesson`, and `application`.
- `internal/sanitize`: `Line` drops line breaks, tabs, and U+2028,
  keeps spaces and ref characters (`feat/x.y_z-1`), and caps length.
- `internal/mcp/server/extract`: `SanitizedOpts` flattens and caps
  `priority`, `section`, `branch`, and `commit`.
- `internal/mcp/server`, end to end through `tools/call`:
  - in a temporary git repository, for both tools,
    `session_id: "unknown"` with `branch: "unknown"` and no `commit`
    stores the repository's branch and short HEAD, read from git
    rather than assumed (the default branch name varies);
  - explicit values are stored verbatim, even in a repository whose
    real branch differs;
  - the issue's probe (`type=learning`, no `branch`/`commit`) with
    `session_id: "unknown"` succeeds;
  - without `session_id` the call still fails, naming only
    `--session-id`.

## Acceptance

- The issue's JSON-RPC probe with `"session_id": "unknown"` added
  returns `Added learning to LEARNINGS.md` on a fresh `ctx init`.
- `tools/list` shows `section`, `priority`, `session_id`, `branch`,
  and `commit` on both entry tools.
- `golangci-lint run` is clean, and there are no new test failures
  against a clean `upstream/main`.

## Non-Goals

- **Persisting provenance for decisions and learnings.** `entry.Write`
  passes session/branch/commit only to `format.Task`
  (`internal/entry/write.go:93-97`); `format.Decision` and
  `format.Learning` take none, so `DECISIONS.md` and `LEARNINGS.md`
  carry no provenance on any path. Persisting it, or narrowing the
  gate to tasks, is a maintainer decision raised on #190.
- A server-minted session ID, or a server-side `unknown` default for
  `session_id`. The latter is offered on #190 as a one-line
  alternative.
- MCP-shaped field names in validation errors (`--session-id`): the
  validator is shared with the CLI.
- Changing CLI provenance behavior or flag sanitization, `.ctxrc`
  `provenance_required`, or the `ctx mcp serve` exemption from the
  git gate.
