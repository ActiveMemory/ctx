# Spec: MCP Handler Flatten & Entity Promotion

## Goal

Eliminate the `mcp/handler.Handler` god-object and the last
grandfathered cross-package MCP type violations by promoting
the ambient runtime to `entity.MCPDeps`, moving session data
to `entity.MCPSession`, and converting domain methods to free
functions. Close out `TASKS.md:45` ("Move 6 grandfathered
cross-package MCP types to entity/").

## Motivation

- The 6-item grandfather list in `internal/audit/cross_package_types_test.go`
  had been parked since the `sameModule` heuristic was hardened,
  blocking progress on strict cross-package-type enforcement.
- `handler.Handler` was carrying three fields (`ContextDir`,
  `TokenBudget`, `*session.State`) purely so the server could
  thread them through route dispatch. It acted as a parameter
  object pretending to be a god object: 17 methods, of which
  13 only read `ContextDir` and 2 actually mutated session state.
- Moving thick behavioral types to `entity/` naively would
  violate entity's "no I/O methods" charter, so a pure/
  behavioral split was needed.

## Scope

### Phase 1 — sub-task moves

- `def/prompt.EntrySpec`, `def/prompt.EntryField` →
  `entity.PromptEntrySpec`, `entity.PromptEntryField`
- `internal/mcp/session/*` collapsed into `internal/mcp/handler/*`
  (single consumer, no other packages referenced session)
- `internal/mcp/server/poll/` → `internal/mcp/server/dispatch/poll/`
  (ancestor→descendant relationship exempts the crossing)

### Phase 2 — god-object flatten

- `entity.MCPSession` — pure data + pure mutation methods
  (`RecordToolCall`, `RecordAdd`, etc.). No I/O methods.
- `entity.PendingUpdate` — data type
- `entity.MCPDeps { ContextDir, TokenBudget, *MCPSession }` —
  parameter object, replaces the Handler struct
- All 17 former `(h *Handler) Foo(...)` methods → free
  functions `handler.Foo(d *entity.MCPDeps, ...)`
- `handler.CheckGovernance` stays in `handler/` as a free
  function (does violations-file I/O via `readAndClearViolations`)
- `server.Server.handler *handler.Handler` →
  `server.Server.deps *entity.MCPDeps`
- `governance_test.go` rewritten to construct `*entity.MCPDeps`
  instead of `*State`
- Route layer (`route/tool/*`, `route/prompt/*`) updated to
  pass `*entity.MCPDeps` instead of `*handler.Handler`

### Phase 3 — absorbed into Phase 2

Splitting `CheckGovernance` I/O from the data struct fell out
naturally when `State` moved to entity and the method became
a free function in handler.

### Follow-up sweep — purge all remaining grandfathered types

The user cleared `grandfatheredTypes` to the empty map and
demanded the underlying violations be fixed, not re-granted.
Fixed 35 violations across 18 packages by moving type
declarations into per-package `types.go` files (methods
stay in their source files, which then carry no type
declarations and are skipped by the audit).

Hardened the `grandfatheredTypes` guardrail comment to
explicitly forbid drive-by additions by any agent and
require a dedicated PR with per-entry justification.

### Follow-up — interface segregation

Interface types moved out of `types.go` into their own
`<name>.go` files because they are behavioral contracts,
not data blueprints. The audit (`TestTypeFileConvention`)
was updated to treat interface types as pure — they cannot
carry method receivers, so the "must have exported receiver"
rule does not apply.

- `type Session interface` → `internal/journal/parser/session.go`
- `type GraphBuilder interface` → `internal/cli/dep/core/builder/graph_builder.go`

## Non-goals

- No behavior change. All MCP tool/resource/prompt dispatch
  flows continue to work exactly as before.
- No spec/protocol changes.
- No changes to the grandfather-list convention itself, only
  to what is listed and the enforcement comment.

## Verification

- `go build ./...` clean
- `make test` all pass
- `make lint` 0 issues
- `TestCrossPackageTypes` grandfathered count: 6 → 0
- `TestTypeFileConvention` grandfathered list: 34 → 0
- `TestDocCommentStructure` grandfathered count: 0 (unchanged)
- `TestNoDeadExports` passes
- `handler.Handler` struct no longer exists
- `mcp/session/` directory no longer exists
- `mcp/server/poll/` directory no longer exists (now under
  `mcp/server/dispatch/poll/`)
- `def/prompt/entry.go` no longer exists (types moved to entity)
