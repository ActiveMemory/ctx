# Spec: MCP Server (`ctx mcp`)

Expose `ctx` context over the Model Context Protocol: making `ctx`'s
full behavioral layer available to any MCP-compatible AI tool, not
just Claude Code.

**Catalyst**: [PR #27](https://github.com/ActiveMemory/ctx/pull/27)
by @CoderMungan: proof-of-concept MCP server with resources and tools.

This spec defines the target architecture that the PR should align to.

---

## Problem

`ctx`'s most powerful features are Claude Code-exclusive:

| Feature           | Claude Code | Cursor | Windsurf | Copilot | Aider |
|-------------------|:-----------:|:------:|:--------:|:-------:|:-----:|
| Context files     | yes         | yes    | yes      | yes     | yes   |
| `ctx agent` output| yes         | yes    | yes      | yes     | yes   |
| Hooks (behavioral)| yes         | no     | no       | no      | no    |
| Skills (workflows)| yes         | no     | no       | no      | no    |
| CLAUDE.md inject  | yes         | no     | no       | no      | no    |

The hook and skill architecture is what makes `ctx` **disciplined**, not
just *available*. Without it, agents can read context but don't get
nudged to persist learnings, load context before acting, or checkpoint
when context runs low.

Today, non-Claude-Code tools get context via `ctx hook generate`
(*which produces tool-specific rule files*) or by pasting `ctx agent`
output. Both are manual, static, and lack behavioral scaffolding.

**MCP changes this.**:

It's a standard protocol (*JSON-RPC 2.0 over stdin/stdout*) supported by 
Claude Desktop, Cursor, Windsurf, VS Code Copilot, and a growing ecosystem. 

A single `ctx mcp serve` replaces per-tool integrations with a universal, 
live, interactive channel.

But **MCP is more than a read API**: It has three capability categories
that map directly to ctx's architecture:

| MCP Capability | ctx Equivalent              | What it enables                    |
|----------------|-----------------------------|------------------------------------|
| **Resources**  | Context files, `ctx agent`  | Live read access to `.context/`    |
| **Tools**      | CLI commands                | Validated writes + queries         |
| **Prompts**    | Skills, AGENT_PLAYBOOK      | Behavioral scaffolding & workflows |

The Prompts capability is the key insight: **MCP Prompts are the
protocol-native equivalent of ctx skills.** They let the server
encode *how* to work, not just *what* to read; making ctx's
behavioral discipline available to every MCP client.

## How Does MCP Fit `ctx`'s Evolution and Manifesto?

MCP is a natural and important fit, but with nuance.

`ctx`'s core thesis is: "Without durable context, intelligence resets." 

The delivery mechanism has been tool-specific hooks (*Claude Code hooks, 
Cursor rules, Aider conventions*). MCP offers a protocol-level delivery 
mechanism: one implementation that works across all MCP-compatible tools.

**Where it aligns**:

* "*Context, not prompts*": MCP resources are exactly this. Instead of users  
  pasting context, the tool can pull it via protocol
* Tool-agnostic vision: `ctx` already supports Claude, Cursor, Aider, Copilot 
  via `ctx hook generate`. MCP would be the cleanest version of this
* Markdown-on-filesystem: the MCP server is a read layer over `.context/` files, 
  no new state
* Local-first, no telemetry: `stdin/stdout`, no network

**Where it needs care**:

* `ctx`'s hook architecture (`UserPromptSubmit`, `PreToolUse`, `PostToolUse`) 
  does far more than serve files: **it enforces invariants**, nudges
  behavior, tracks context size, gates commits. MCP resources/tools alone 
  cannot replace hooks. The hooks are ctx's nervous system; MCP
  would be a read/write API
* The manifesto says "**Verification, not vibes**": the current MCP tools 
  (`ctx_add`, `ctx_complete`) let AI tools mutate context without the
  guardrails that hooks provide. This is a philosophical gap
* `ctx`'s strategic position is L3 (system memory): MCP fits as an access 
  protocol, but ctx's value is in the structure and discipline around context,
  not just serving files over a protocol

  **Bottom line**: MCP is a great delivery channel for `ctx` context. It is 
  **not** a replacement for the hook/skill architecture. It should be
  positioned as "*universal read access and guarded write operations*".

## Design Principles

**Six invariants** govern the MCP server design:

1. **Markdown-on-filesystem**: all state remains in `.context/` files;
   the MCP server is a *view* over them, never a separate store
2. **Delegate, don't reimplement**: every MCP tool handler calls an
   existing internal package; zero duplicated logic
3. **Read-first, write-guarded**: resources are freely available;
   tools that mutate validate structurally and advise behaviorally
4. **Prompts encode discipline**: behavioral scaffolding (*when to
   persist, when to reflect, when to checkpoint*) lives in MCP Prompts
5. **Session state is advisory, never blocking**: the server tracks
   what has been read/called and adds advisory text to responses; it
   never refuses a valid operation
6. **Token budgeting preserved**: the agent resource respects the
   configured token budget, same as `ctx agent --budget`

These extend (**not replace**) the six `ctx` design invariants documented
in `ARCHITECTURE.md`.

## Design Invariants and Constraints

The MCP server operates within **hard boundaries** derived from the
**`ctx` Manifesto**, security architecture, and design philosophy. 

**These are not trade-offs to revisit**: Yhey are **load-bearing walls**.

### Local-only. No remote serving.

The MCP server communicates over `stdin/stdout` with a process on
the same machine. It will **never** expose an `HTTP` endpoint, listen
on a network socket, or serve context to remote clients.

This follows from three invariants:

- **No telemetry: no data leaves the local machine**: Context
  files contain the project's reasoning: *why* decisions were
  made, *what* was tried and failed, *which* conventions emerged.
  **This is more sensitive than source code**:" it's the intent
  behind the code. Sending it over the wire is a fundamentally
  different trust decision.
- **Local-first: no network required**: `ctx` is infrastructure.
  Infrastructure that depends on someone else's uptime is a
  dependency, **not** infrastructure. The MCP server must work on an
  airplane, behind a firewall, and during an outage.
- **Security model is trust-by-locality**: `ctx` has no auth, no
  TLS, no access control. The machine boundary **IS** the trust
  boundary. `AES-256-GCM` protects data at rest; locality protects
  data in use. Adding remote access would require an entirely
  new security layer that doesn't exist and isn't planned.

### Teams share context through git, not servers.

`.context/` is version-controlled alongside the code it describes.
Teams share context the same way they share code: push, pull,
review, merge. This is **deliberate**:

* Context changes are **reviewable**: a new decision or
  convention goes through the same PR process as a code change.
* Context has **history**: `git log .context/` shows when and
  why context evolved.
* Context has **branching**: experimental decisions live on
  feature branches and merge with the code they inform.
* Context has **ownership**: the repo's access controls govern
  who can modify context.

This pattern extends beyond context files. Journal import/export
handles session history from external locations through the
filesystem; **not** through a service. Cross-machine sharing is a
file operation (copy, rsync, git merge), **not** a protocol operation.

The entire data model is designed so that `cp`, `git pull`, and
`diff` are the collaboration primitives.

A remote MCP server would bypass all of this: no review, no
history, no branching, no ownership. The **`ctx` Manifesto** says "*structure
enables scale*" and "*ambiguity is a system failure*": A remote
context endpoint that anyone  or write is **the opposite
of structure**.

### The machine boundary is the trust boundary.

Every security layer in `ctx` assumes local execution:

* Layer 2 (symlink rejection) assumes local filesystem
* Layer 3 (boundary validation) assumes local path resolution
* Layer 4 (permission deny list) assumes local process control
* Layer 5 (plugin hooks) assumes local hook execution

The MCP server inherits these assumptions. It does not add
network-facing security because it does not face a network.

### `stdin/stdout` is the only transport.

MCP supports two transports: stdio and SSE (HTTP-based). ctx
uses stdio exclusively. This is not a "*start with stdio, add
HTTP later*" plan: It's a **permanent design choice** that follows
from the local-only constraint. If the MCP spec adds new
local-only transports in the future (*e.g., Unix domain sockets*),
those could be evaluated; network transports **will not**.

## Architecture

```
MCP Client (Claude Desktop, Cursor, VS Code, ...)
    |
    | stdin/stdout (JSON-RPC 2.0)
    |  (each request carries context_dir or uses roots/fallback)
    |
    v
ctx mcp serve                        [path-agnostic]
    |
    +-- Resolve ────> context_dir resolution chain
    +-- Resources ──> context.Load()    [read .context/ files]
    +-- Tools ──────> internal packages [validated operations]
    +-- Prompts ────> embedded templates [behavioral scaffolding]
    +-- Session ────> in-memory state   [advisory tracking, per context_dir]
```

**The server is path-agnostic**: It does not hardcode a single
project's `.context/` path at construction time. 

Instead, each request resolves the context directory through a fallback chain:

| Priority | Source                              | When it applies                    |
|----------|-------------------------------------|------------------------------------|
| 1        | Explicit `context_dir` on the call  | Client knows exactly where         |
| 2        | MCP roots (client declares roots)   | Client supports roots capability   |
| 3        | `--context-dir` flag at server start| Single-project deployment          |
| 4        | CWD-based discovery (walk up)       | Fallback                           |

This means a single `ctx mcp serve` process can serve three (*or more*)
projects in the same session: The client just passes different
`context_dir` values. No need to spawn one server per project.

The server is **stateless across invocations**: Each `ctx mcp serve`
process starts fresh. Within a single process lifetime, it maintains
lightweight **session state** (*a set of flags, keyed by resolved
context directory*) used solely for advisory messages in tool
responses.

### Package Layout

```
internal/mcp/
├── doc.go          # Package documentation
├── protocol.go     # JSON-RPC 2.0 + MCP type definitions
├── server.go       # Server lifecycle, dispatch, I/O
├── resources.go    # Resource list + read handlers
├── tools.go        # Tool list + call handlers
├── prompts.go      # Prompt list + get handlers          [v0.2]
├── session.go      # Session state tracking + advisory   [v0.2]
└── server_test.go  # Tests

internal/cli/mcp/
└── mcp.go          # Cobra command: ctx mcp serve [+ config]
```

---

## Phase Plan

### v0.1 — Foundation (Present)

**Goal**: Ship a correct, **minimally-delightful** MCP server that any tool
can use to read `ctx` context and perform validated writes.

Aligns with PR #27 scope, plus corrections for spec compliance.

#### Resources

Expose context files as read-only MCP resources. Resource URIs
encode the context directory path, making the server path-agnostic:

```
ctx://{context_dir}/{file}
```

When the server starts with `--context-dir` or discovers a single
project via CWD, it registers static resources with resolved paths.
When serving multiple projects (via roots or per-request
`context_dir`), clients use the full URI form.

**Static resources** (single-project mode):

| URI                          | Source                | Description                            |
|------------------------------|-----------------------|----------------------------------------|
| `ctx://context/constitution` | CONSTITUTION.md       | Hard rules that must never be violated |
| `ctx://context/tasks`        | TASKS.md              | Current work items and their status    |
| `ctx://context/conventions`  | CONVENTIONS.md        | Code patterns and standards            |
| `ctx://context/architecture` | ARCHITECTURE.md       | System architecture documentation      |
| `ctx://context/decisions`    | DECISIONS.md          | Architectural decisions with rationale |
| `ctx://context/learnings`    | LEARNINGS.md          | Gotchas, tips, and lessons learned     |
| `ctx://context/glossary`     | GLOSSARY.md           | Project-specific terminology           |
| `ctx://context/playbook`     | AGENT_PLAYBOOK.md     | How agents should use this system      |
| `ctx://context/agent`        | Assembled packet      | All files in priority order, budgeted  |

**Resource template** (*multi-project mode, registered via
`resources/templates/list`*):

```
ctx://{context_dir}/{name}
```

The client fills `context_dir` with the absolute path to a
`.context/` directory. `ValidateBoundary()` ensures the resolved
path stays within the declared directory.

The `agent` resource assembles files using `config.FileReadOrder` and
respects the token budget from `rc.TokenBudget()`. This is the same
output as `ctx agent --budget N`, not a naive concatenation.

**Difference from PR #27**: PR #27 concatenates files without budgeting
and omits AGENT_PLAYBOOK. **Both must be corrected**.

Resource implementation uses `context.Load()`: the same path as
every other `ctx` consumer.

#### Tools

Four tools, all delegating to existing internal packages. Every
tool accepts an optional `context_dir` argument: if omitted, the
server resolves it via the fallback chain (roots > flag > CWD).

| Tool           | Package               | Mutates? | Description                          |
|----------------|-----------------------|----------|--------------------------------------|
| `ctx_status`   | `context.Load()`      | No       | File count, tokens, per-file summary |
| `ctx_drift`    | `drift.Detect()`      | No       | Stale paths, missing files, warnings |
| `ctx_add`      | `add.WriteEntry()`    | Yes      | Add task/decision/learning/convention|
| `ctx_complete` | `cli/complete` logic  | Yes      | Mark a task done by number or text   |

**Critical**: `ctx_complete` **MUST** delegate to the same code path as
`ctx complete`, not reimplement task parsing: 

PR #27 reimplements this in `tools.go` (*~60 lines of task matching logic*): 
This **must** use the existing `internal/task` package and the complete 
command's entry point instead.

**`ctx_add` context-dir threading**: `add.WriteEntry()` resolves
the context directory via `rc.ContextDir()` (line 76 of `run.go`),
while the MCP server receives `contextDir` as a constructor argument
and passes it explicitly to `context.Load()`, `drift.Detect()`, and
`toolComplete`. This means `toolAdd` bypasses `s.contextDir` entirely
and goes through `rc` instead.

This (*purely coindidentially, and by sheer luck*) works "*today*"" 
because `PersistentPreRunE` calls `rc.OverrideContextDir()` before `RunE`, 
so `rc.ContextDir()` returns the same value as `s.contextDir`. 

**But it's fragile**: The two paths are coupled through global state rather 
than explicit parameter passing. If the `add` package is ever used in a 
context where `rc` isn't initialized (*e.g., tests, library use*), it will 
silently use the default `.context` path.

**Recommendation**: Add a `ContextDir` field to `EntryParams`. When
set, `WriteEntry` uses it; when empty, it falls back to
`rc.ContextDir()`. This makes the dependency explicit without
breaking existing callers. The MCP server sets
`params.ContextDir = s.contextDir`; the CLI doesn't set it and
gets the current behavior.

Tool annotations (MCP spec 2025-03-26):

```json
{
  "ctx_status": {"readOnlyHint": true},
  "ctx_drift":  {"readOnlyHint": true},
  "ctx_add":    {"readOnlyHint": false, "destructiveHint": false, "idempotentHint": false},
  "ctx_complete": {"readOnlyHint": false, "destructiveHint": false, "idempotentHint": true}
}
```

Annotations tell MCP clients whether to auto-approve or prompt the
user. Read-only tools can be auto-approved; **write tools should
surface for confirmation in clients that support this**.

#### CLI

```
ctx mcp serve
```

Starts the MCP server on `stdin/stdout`. Intended to be launched by
MCP clients, not run directly by users.

Inherits context directory from the standard resolution chain:
`--context-dir` flag > `CTX_DIR` env > `.ctxrc` > default `.context/`.

Uses `config.AnnotationSkipInit` to bypass the normal `ctx init`
guard (same as `ctx system` commands).

#### Protocol

* JSON-RPC 2.0 over stdin/stdout (newline-delimited)
* MCP protocol version: `2025-03-26` (latest stable)
* Server capabilities: `resources`, `tools`
* Handles: `initialize`, `ping`, `notifications/initialized`,
  `notifications/cancelled`, `resources/list`, `resources/read`,
  `tools/list`, `tools/call`
* Scanner buffer: 1 MB (matches PR #27)

**Version negotiation**: the MCP `initialize` handshake includes the
client's `protocolVersion`. The server should:

1. If the client sends `2025-03-26`: full feature set including
   tool annotations (`readOnlyHint`, `destructiveHint`, etc.)
2. If the client sends an older version (e.g., `2024-11-05`):
   graceful degradation: omit `annotations` from tool definitions,
   omit any fields the older spec doesn't recognize. Resources and
   tools still work; the client just doesn't get annotation hints.
3. The server always *responds* with its own `protocolVersion`
   (`2025-03-26`). Per the MCP spec, the effective version is the
   lower of the two: the server respects this when building
   responses.

This means: **target the newest spec, degrade for older clients**,
**never** refuse to serve. The degradation is cosmetic (*missing hints*),
not functional (*missing capabilities*).

#### What Changes from PR #27

| PR #27 Current                           | Spec Requirement                                |
|------------------------------------------|-------------------------------------------------|
| Agent resource: naive concatenation      | Token-budgeted assembly via `context.Load()`    |
| Missing AGENT_PLAYBOOK resource          | Add `ctx://context/playbook`                    |
| `toolComplete` reimplements task parsing | Delegate to existing `internal/task` + complete |
| `//nolint:gosec` inline suppression      | Handle the error properly (see action items)    |
| No `TestToolAdd` test                    | Required: test all tools                        |
| Phantom `eventlog` in ARCHITECTURE.md    | Remove (not part of this PR)                    |
| No tool annotations                      | Add `readOnlyHint` annotations                  |
| `site/` not updated with docs            | Regenerate site after docs changes              |

#### Immediate Action Items for PR Merge

These must be resolved before the PR can merge. Each item includes
the concern, resolution, and justification.

**1. Remove `//nolint:gosec`: handle the error instead**

*Concern*: The PR suppresses gosec G705 with an inline `//nolint`:

```go
_, _ = fmt.Fprintf(s.out, "%s\n", out) //nolint:gosec // G705: stdout, not HTTP response
```

*Resolution*: Don't suppress the lint; handle the error. The
`_, _ =` already discards the error; instead, check it:

```go
if _, writeErr := fmt.Fprintf(s.out, "%s\n", out); writeErr != nil {
    return writeErr
}
```

A write failure to stdout in the server's main loop is an
unrecoverable I/O error (*broken pipe, closed fd*). Returning the
error lets `Serve()` exit cleanly instead of silently writing
into the void.

*Justification*:

* Project convention (CONVENTIONS.md): "zero `//nolint:errcheck`
  policy: handle errors, don't suppress them"
* Leaking this into `.golangci.yml` as a global suppression would
  weaken the lint for the entire codebase to accommodate one line
* LEARNINGS.md (2026-03-04): "golangci-lint v2 ignores inline
  nolint directives for some linters": the suppression may not
  even work reliably with the project's linter version
* The `writeError` helper has the same pattern and should also
  handle its write error (*or accept that it's best-effort and
  use `_ =` with a comment explaining why*)

**2. `toolComplete` must delegate, not reimplement**

*Concern*: `tools.go` contains ~60 lines of task matching logic
(line scanning, regex matching, number vs text search, checkbox
replacement) that duplicates what `internal/task` and
`internal/cli/complete` already provide.

*Resolution*: Refactor `toolComplete` to call the existing
completion code path. If the existing package API doesn't support
being called with an explicit `contextDir` (i.e., it uses
`rc.ContextDir()` internally), extend the package API — don't
copy the logic.

*Justification*:

* `CONVENTIONS.md`: "Non-test code: apply the rule of three —
  extract when a block appears 3+ times." This is appearance
  number two, but the logic is complex enough that even two
  copies will drift (task regex changes, checkbox format changes,
  `#done` timestamp additions)
* Design principle #2 of this spec: "**Delegate, don't reimplement**"
* The existing code handles edge cases (*subtasks, `#in-progress`
  labels, timestamp tagging*) that the PR's reimplementation misses

**3. Agent resource must respect token budgeting**

*Concern*: The `ctx://context/agent` resource concatenates all
files naively. The `ctx agent` CLI command applies token budgeting
to fit within context window limits: this is a key differentiator.

*Resolution*: Use `context.Load()` with `rc.TokenBudget()` and
apply the same truncation/omission logic as `ctx agent`. Lower-
priority files should be truncated or listed as "Also Noted" when
the budget is exceeded.

*Justification*:

* Design principle #6: "Token budgeting preserved"
* Without budgeting, the agent resource can return unbounded
  content, which defeats the purpose of ctx's priority ordering

**4. Add AGENT_PLAYBOOK resource**

*Concern*: The playbook is in `config.FileReadOrder` and defines
how agents should use the context system. It's missing from the
resource table.

*Resolution*: Add `ctx://context/playbook` mapping to
`AGENT_PLAYBOOK.md` in `resourceTable`.

**5. Add `TestToolAdd` test**

*Concern*: `ctx_add` is the most complex tool (7 argument fields,
4 entry types, conditional validation) and has **zero test coverage**.

*Resolution*: Add tests for each entry type (*task, decision,
learning, convention*), validation failures (*missing required
fields*), and verify the file is actually written.

**6. Remove phantom `eventlog` from ARCHITECTURE.md**

*Concern*: The PR adds `internal/eventlog` to the ARCHITECTURE.md
core packages table, but no `eventlog` package exists in this PR.

*Resolution*: Remove the `eventlog` line from the ARCHITECTURE.md
diff. Only document packages that the PR introduces.

**7. Regenerate `site/` after docs changes**

*Concern*: The PR adds `docs/cli/mcp.md` and updates
`docs/cli/index.md` but doesn't include the regenerated `site/`
directory.

*Resolution*: Run the site generator and include `site/` changes
in the commit.

*Justification*:

- `CONVENTIONS.md`: "*Always stage site/ when committing docs/
  changes: the generated HTML is tracked in git with no CI
  build step*"
- `LEARNINGS.md` (2026-02-27): "*site/ directory must be committed
  with docs changes*"

---

### v0.2 — Behavioral (Near Future)

**Goal**: Bring `ctx`'s behavioral discipline to MCP clients via
Prompts and session-aware advisory responses.

#### MCP Prompts

Prompts are server-defined prompt templates that clients discover
via `prompts/list` and invoke via `prompts/get`. The server returns
a list of messages that the client injects into the model's context.

This is how ctx encodes *workflows*: The protocol-native equivalent
of skills.

**Source of truth**: skill files. Skills that should be exposed as
MCP prompts declare `mcp_prompt: true` in their frontmatter. At
build time, prompt definitions are extracted from skill files and
embedded in `internal/assets/`. This ensures a single source of
truth: Edit the skill, the MCP prompt updates automatically.

**Drift detection**: `ctx drift` gains a check that compares the
embedded MCP prompt definitions against the live skill content.
If a skill with `mcp_prompt: true` has changed since the last
build (content hash mismatch), drift reports a warning:
"MCP prompt for ctx-reflect is stale — rebuild to sync."
The `/ctx-skill-audit` skill also checks for this.

| Prompt               | Arguments           | Returns                                                                      |
|----------------------|---------------------|------------------------------------------------------------------------------|
| `ctx-session-start`  | `budget` (optional) | Context packet + playbook summary + current tasks + "what to work on next"   |
| `ctx-add-decision`   | `title`             | Structured template: context, rationale, consequences: with instructions     |
| `ctx-add-learning`   | `title`             | Structured template: context, lesson, application: with instructions         |
| `ctx-reflect`        | none                | "What did you learn? What decisions were made? What's left?": guides capture |
| `ctx-checkpoint`     | none                | "Summarize progress, persist to TASKS.md, note blockers"                     |

**How prompts encode discipline**: example `ctx-session-start`:

```json
{
  "name": "ctx-session-start",
  "description": "Load project context and orient for a new work session. Use this at the beginning of every session.",
  "arguments": [
    {
      "name": "budget",
      "description": "Token budget for context assembly (default: from .ctxrc or 8000)",
      "required": false
    }
  ]
}
```

When the client calls `prompts/get` for `ctx-session-start`, the
server returns messages containing:

1. The assembled context packet (same as `ctx://context/agent`)
2. Playbook summary: "Read context before acting. Persist decisions
   and learnings as you go. Checkpoint before context runs low."
3. Current task status: N pending, M in-progress
4. Suggestion: "Consider starting with: [highest-priority pending task]"

This is the MCP equivalent of what Claude Code gets from hooks
(`context-load-gate`) + skills (`/ctx-remember`) + CLAUDE.md
instructions. One prompt, universal reach.

**How `ctx-reflect` encodes end-of-session ceremony:**

The returned messages ask the model to:

1. List decisions made during this session
2. List gotchas or surprises encountered
3. List tasks completed and tasks remaining
4. For each decision/learning, call `ctx_add` to persist it

This mirrors the `/ctx-reflect` skill but works in any MCP client.

#### Session State and Advisory Responses

The server tracks lightweight session state: not persisted, just
in-memory flags for the current process lifetime:

```go
type sessionState struct {
    resourcesRead  map[string]bool  // Which resources have been read?
    toolCalls      int              // Total tool calls this session
    tasksCompleted int              // Tasks completed this session
    addsPerformed  map[string]int   // Adds by type (decision, learning, etc.)
}
```

Advisory rules: Soft signals appended to tool response text:

| Condition                                        | Advisory text                                                          |
|-------------------------------------- -----------|------------------------------------------------------------------------|
| First tool call, no resources read yet           | "Tip: Load context first with resources/read for ctx://context/agent"  |
| 3+ tasks completed, 0 learnings added            | "You've completed N tasks. Consider capturing learnings."              |
| 10+ tool calls, no resources read                | "No context loaded this session. Decisions may miss existing context." |
| `ctx_add(decision)` without context/rationale    | Structural validation rejects — this is a hard error, not advisory     |

Advisory text is **appended to the normal response**, not returned as
errors. The AI reads it and decides whether to act. This matches how
ctx hooks work today: they output nudges to stderr, and the agent
chooses to follow them.

**Key principle**: advisory never blocks. A valid `ctx_add` call
succeeds even if the agent hasn't read context. The advisory just
says "you might want to."

**Opt-out**: the existing `ctx pause` / `ctx resume` mechanism
controls advisory. When paused, the MCP server suppresses advisory
text in tool responses: same as hooks going silent during a pause.
The server checks pause state (via `.context/state/`) before
appending advisory messages. Agents that find advisories noisy
call `ctx pause`; `ctx resume` re-enables them. No new config key
needed — this reuses an established pattern already documented in
the AGENT_PLAYBOOK.

#### Resource Templates

Add support for dynamic resource URIs:

```json
{
  "uriTemplate": "ctx://context/{name}",
  "name": "Context file by name",
  "description": "Read any context file by its filename"
}
```

This allows clients to request arbitrary files from `.context/`
without the server pre-enumerating all of them. Bounded by
`validation.ValidateBoundary()` to prevent path traversal.

#### Resource Subscriptions

Register for change notifications on specific resources:

```
Client → resources/subscribe { uri: "ctx://context/tasks" }
Server → (watches .context/TASKS.md for changes)
Server → notifications/resources/updated { uri: "ctx://context/tasks" }
```

Implementation: lightweight polling (stat mtime every N seconds)
rather than fsnotify, to avoid platform-specific dependencies.
Polling interval configurable via `.ctxrc`.

#### Additional Tools

| Tool           | Package                 | Mutates? | Description                     |
|----------------|-------------------------|----------|---------------------------------|
| `ctx_recall`   | `recall/parser`         | No       | Recent session summaries        |
| `ctx_compact`  | `cli/compact`           | Yes      | Archive completed tasks         |
| `ctx_next`     | task analysis           | No       | Suggest next task to work on    |

#### CLI Addition

```
ctx mcp config [tool]
```

Generate MCP client configuration for a specific tool:

```bash
ctx mcp config claude-desktop   # Outputs JSON for claude_desktop_config.json
ctx mcp config cursor           # Outputs JSON for .cursor/mcp.json
ctx mcp config vscode           # Outputs JSON for .vscode/mcp.json
```

This replaces manual configuration — users copy-paste the output.

---

### v0.3 — Convergence (Far Future)

**Goal**: MCP becomes **the primary integration channel**, with hooks
as a Claude Code-specific optimization layer on top.

#### Full Prompt Parity with Skills

Every user-invocable ctx skill becomes an MCP prompt:

| Skill              | MCP Prompt             |
|--------------------|------------------------|
| `/ctx-remember`    | `ctx-session-start`    |
| `/ctx-status`      | `ctx-health-check`     |
| `/ctx-next`        | `ctx-pick-task`        |
| `/ctx-commit`      | `ctx-commit`           |
| `/ctx-reflect`     | `ctx-reflect`          |
| `/ctx-implement`   | `ctx-implement`        |
| `/ctx-brainstorm`  | `ctx-brainstorm`       |
| `/ctx-wrap-up`     | `ctx-wrap-up`          |

Prompts are generated from skill definitions: single source of
truth, dual deployment (*skill file + MCP prompt*).

#### Sampling: Server-Initiated Reasoning

MCP's Sampling capability allows the server to ask the *client's
model* to perform work. Use cases:

* **Auto-enrichment**: when a new decision is added via `ctx_add`,
  the server asks the model to check for conflicts with existing
  decisions
* **Drift explanation**: when `ctx_drift` finds violations, the
  server asks the model to suggest fixes
* **Session summary**: at session end (or after N tool calls), the
  server asks the model to summarize what happened

This is powerful but intrusive — it means the server initiates
computation. Design with explicit opt-in.

#### `ctx hook generate mcp`

Instead of generating tool-specific rule files, `ctx hook generate`
gains an `mcp` target that outputs the MCP client config for the
detected tool. This positions MCP as the recommended integration path
for non-Claude-Code tools:

```bash
ctx hook generate cursor     # Today: generates .cursorrules
ctx hook generate cursor     # Future: generates .cursor/mcp.json (MCP config)
```

The rule-file approach remains as a fallback for tools that don't
support MCP.

#### Multi-Project Serving

Handled by the path-agnostic design from v0.1: a single
`ctx mcp serve` process can serve multiple projects because
every tool and resource accepts an explicit `context_dir`.

No special multi-project code needed. The client passes different
paths; the server resolves, validates boundaries, and serves.
Session state (advisory tracking) is keyed by resolved
`context_dir` so advisories for project A don't leak into
project B.

---

## Behavioral Discipline: The Full Picture

This section explains *why* the three-layer design (*resources +
tools + prompts*) is necessary and how the layers compose.

### The Problem with Tools Alone

An MCP server with only resources and tools is merely a CRUD API over
`.context/`. Any AI agent  and write; but nothing guides
*when* or *why*. This is the equivalent of giving someone access to
a codebase without a contributing guide.

`ctx`'s value is **not** the files. It's the **discipline around the files**:
read before acting, persist after deciding, reflect before wrapping up.
Without that discipline, context accumulates entropy instead of value.

### How Hooks Work Today (Claude Code)

```
UserPromptSubmit → "Have you loaded context?" (context-load-gate)
                 → "Context is 80% full" (check-context-size)
                 → "You have pending reminders" (reminders)

PostToolUse      → "Commit detected — did you capture learnings?" (post-commit)
                 → "3 tasks done — consider reflecting" (task-completion-nudge)
```

Hooks fire *automatically* at lifecycle points. The agent doesn't
choose to invoke them: they are injected.

### How MCP Prompts Approximate This

MCP prompts don't fire automatically: The client must invoke them.
But prompts can *instruct the model* to follow a behavioral protocol:

**`ctx-session-start` prompt content includes:**

> You are working with ctx, a persistent context system.
>
> Before making any changes, read the context files provided above.
> As you work, follow this cycle: **Work -> Reflect -> Persist.**
>
> After completing a task, ask yourself: "If this session ended now,
> would the next session know what happened?" If not, persist
> something before continuing.
>
> When you complete 3+ tasks, call ctx_add with type "learning" to
> capture what you discovered. When you make a design choice between
> alternatives, call ctx_add with type "decision" to record the
> rationale.
>
> Available tools: ctx_status, ctx_add, ctx_complete, ctx_drift.

This is ~80% of what hooks achieve, delivered through the standard
protocol. The remaining 20% (*automatic firing, lifecycle triggers*)
requires native client support (*which is what Claude Code hooks
provide*). 

**For non-Claude-Code tools, 80% is a massive improvement over 0%**.

### Advisory Responses: The Safety Net

Even if the agent ignores the prompt instructions, session-aware
advisory responses provide a second layer:

```
Agent calls ctx_complete("task 3") without having read any resources:
→ Tool succeeds, response includes:
  "Completed: Implement caching layer

   Note: No context has been loaded this session.
   Consider reading ctx://context/agent to review
   existing decisions before continuing."
```

The agent sees this in the tool response and can choose to act on
it. This is strictly advisory: the operation already succeeded.

### The Discipline Stack

```
Layer 3: Hooks (Claude Code only)        — automatic, lifecycle-triggered
Layer 2: Prompts (any MCP client)        — invoked, workflow-encoded
Layer 1: Advisory (any MCP client)       — passive, response-embedded
Layer 0: Structural validation (always)  — hard enforcement, rejects invalid input
```

Each layer catches what the layer below misses. Layer 0 is always
active. Layers 1-2 require an MCP client. Layer 3 requires Claude
Code. The more layers active, the more disciplined the session — but
even Layer 0 alone prevents garbage writes.

---

## Edge Cases

| Case                                      | Expected behavior                                         |
|-------------------------------------------|-----------------------------------------------------------|
| `.context/` doesn't exist                 | Resources return empty list; tools return clear error     |
| Context file is empty                     | Resource returns empty text; status shows "EMPTY"         |
| Token budget exceeded                     | Agent resource truncates lower-priority files             |
| `ctx_add` with missing required fields    | Structural validation error (Layer 0)                     |
| `ctx_complete` with ambiguous match       | Error: "multiple tasks match, be more specific"           |
| Client never calls `initialize`           | Server rejects all requests per MCP spec                  |
| Malformed JSON on stdin                   | Parse error response, continue reading                    |
| Very large context file (>1MB)            | Scanner buffer handles up to 1MB; beyond = error          |
| Concurrent MCP + direct file edits        | No locking; next resource read reflects latest disk state |
| `ctx_add` called for project without init | Clear error: ".context/ not found, run ctx init"          |

## Validation Rules

**Structural (Layer 0, always enforced):**

* `ctx_add(decision)` requires: content, context, rationale, consequences
* `ctx_add(learning)` requires: content, context, lesson, application
* `ctx_add(task)` requires: content
* `ctx_add(convention)` requires: content
* `ctx_complete` query must match exactly one pending task
* Resource URIs must match known patterns (no path traversal)
* All file operations bounded by `validation.ValidateBoundary()`

**Advisory (Layer 1, v0.2+, never blocks):**

* No resources read before first write tool → advisory note
* 3+ completions without a learning → advisory note
* 10+ tool calls without resource read → advisory note

## Security

* **Boundary enforcement**: all resource reads and tool writes
  validated against `.context/` boundary via `ValidateBoundary()`
* **No shell execution**: tool handlers call Go functions, never
  `exec.Command`
* **No secrets in responses**: `.ctx.key`, `scratchpad.enc`, webhook
  URLs are never exposed as resources
* **No network access**: server communicates only via stdin/stdout
* **Input sanitization**: tool arguments validated before use;
  `SanitizeFilename()` for any file-derived inputs
* **Resource filtering**: only `.md` files in `.context/` are
  exposed; dotfiles, state files, and subdirectories excluded from
  default resource list

## Configuration

### `.ctxrc` Keys

| Key                    | Type   | Default | Description                           |
|------------------------|--------|---------|---------------------------------------|
| `mcp.poll_interval`    | int    | 5       | Subscription poll interval (seconds)  |

Advisory opt-out uses `ctx pause` / `ctx resume`, not a `.ctxrc`
key. See Session State & Advisory section.

### Environment Variables

No new environment variables. MCP server inherits `CTX_DIR` and
`CTX_TOKEN_BUDGET` from the standard resolution chain.

## Testing

### v0.1

* **Unit**: each handler tested in isolation with injected in/out
  * Initialize handshake
  * Ping/pong
  * Resource list (count, URIs, MIME types)
  * Resource read (each file, agent packet with budgeting)
  * Resource read unknown URI → error
  * Tool list (count, names, schemas)
  * Tool call: status (verify output format)
  * Tool call: add (each type, validation errors)
  * Tool call: complete (by number, by text, ambiguous, missing)
  * Tool call: drift (verify report structure)
  * Tool call: unknown → error
  * Notification → no response
  * Parse error → error response
  * Empty lines → skipped
* **Integration**: end-to-end via `ctx mcp serve` subprocess
  * Send initialize + resources/list, verify resource count
  * Send tools/call ctx_add, verify file written
  * Verify `site/` regeneration not triggered (MCP is headless)

### v0.2

* Advisory: verify advisory text appears after N tool calls
  without resource reads
* Prompts: verify prompt content includes context packet
* Subscriptions: modify file, verify notification sent

## Helpers to Reuse

| Existing Package                | What it provides                          |
|---------------------------------|-------------------------------------------|
| `context.Load()`                | Load all `.context/` files with tokens    |
| `context.EstimateTokens()`      | Token counting                            |
| `config.FileReadOrder`          | Priority ordering                         |
| `config.FileType`               | Map entry type → filename                 |
| `config.RegExTask`              | Task checkbox regex                       |
| `drift.Detect()`                | Context quality checks                    |
| `add.ValidateEntry()`           | Structural validation for entries         |
| `add.WriteEntry()`              | Append entries to context files           |
| `task.Pending()`                | Check if task is pending                  |
| `task.Content()`                | Extract task text                         |
| `validation.ValidateBoundary()` | Path traversal prevention                 |
| `validation.SanitizeFilename()` | Input sanitization                        |
| `rc.ContextDir()`               | Resolved context directory                |
| `rc.TokenBudget()`              | Configured token budget                   |

**Zero reimplementation.** If a tool handler needs logic that exists
in an internal package, it calls that package. If the package API
doesn't support the MCP use case (e.g., missing `contextDir`
parameter), the package API is extended: the logic is not copied.

## Non-Goals

* **Replacing hooks**: MCP complements hooks, doesn't replace them.
  Claude Code hooks remain the most reliable behavioral enforcement.
  MCP provides 80% of the value to 100% of the tools.
* **MCP client implementation**: ctx is an MCP *server*. It does not
  consume MCP services from other servers.
* **Network transport**: MCP supports SSE (HTTP-based) transport.
  ctx uses stdio exclusively: This is a permanent design choice,
  not a "start here, add HTTP later" plan. See Design Invariants.
* **Authentication**: not needed. The machine boundary is the trust
  boundary. See Design Invariants.
* **GUI**: no web dashboard, no visual resource browser. MCP clients
  provide their own UI.
* **Real-time sync**: the server reads files on each request. No
  in-memory caching, no file watchers (except for subscriptions in
  v0.2). This matches ctx's CLI-tool-not-daemon philosophy.
