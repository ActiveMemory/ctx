# internal/assets README

## Problem

`internal/assets/` is the embed root for every non-Go artifact
ctx ships — TypeScript plugin source, Bash and PowerShell
hook scripts, Markdown skill bodies, JSON manifests, YAML
registries. The Go binary literally bakes these bytes in at
compile time (`//go:embed`) and writes them to the user's
filesystem at `ctx setup` time, where the consumer tool
(Claude Code, OpenCode/Bun, Copilot CLI, …) executes them in
its own runtime.

A reader browsing the tree has no signal that this is the
contract. The directory reads — at first glance — as a dumping
ground: `internal/` (a Go-private import boundary) containing
TypeScript, shell scripts, and Markdown. The lack of a
contract document made an experienced reviewer mistake this
for an architectural smell and propose lifting `integrations/`
out of `internal/`, before realising every file in there is
genuinely embedded and the `//go:embed` no-`../` rule makes
relocation costly.

The artifact-vs-tooling confusion has a second cost: there is
no obvious home for *dev tooling* about the embedded payload
(e.g. `package.json` / `tsconfig.json` for the OpenCode plugin
TypeScript). Anyone trying to add a type-check or linter gate
faces the same first question we faced — "where does this go,
and won't it get embedded?"

## Approach

Add `internal/assets/README.md` as the long-form contract
document for the embed tree. Audience: contributors, code
reviewers, security auditors. Complement to `doc.go` (which
remains the Go-doc summary).

The README documents:

1. **Why non-Go code lives under `internal/`** — Go's
   `internal/` convention is about *import privacy*, not
   *source language*. ctx ships as one statically-linked
   binary; `//go:embed` is how foreign-language payloads
   become build-time inputs to that binary. Concrete OpenCode
   trace (TS source → embedded bytes → ctx setup → user
   disk → Bun runtime) makes the lifecycle vivid.
2. **The embed contract** — three-clause definition of what
   belongs in the tree; the hard `//go:embed` no-`../`
   constraint that anchors the directory shape.
3. **Directory map** — every embedded subdirectory with
   language, consumer, deploy target. Single source of truth
   for "where does X live?" questions.
4. **Quality gates today** — honest accounting of what
   `embed_test.go` actually verifies (presence, format
   parse, schema integrity, spot-content) and the known gaps
   (TypeScript type-check, shellcheck, PSScriptAnalyzer,
   skill frontmatter validity).
5. **Adding a new embedded asset** — six-step checklist.
6. **What does not belong here** — negative space; the
   "dev tooling for the payload" exclusion that gives the
   line-30 type-check task its home (sibling tooling
   directory, not nested inside the payload).

## Behavior

### Happy Path

A contributor opens `internal/assets/` in their editor, sees
`README.md` at the top, reads it, and now knows:
- This tree is the embed root for the Go binary.
- Foreign-language files are intentional, not stray.
- Where to add a new asset and what tests to add alongside.
- What test coverage exists today and what does not.
- Where dev tooling for embedded assets should live (sibling,
  not nested).

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Reader only looks at `doc.go` | `doc.go` already exists; README does not replace it, only complements it. |
| Reader wants the Go-doc summary | `doc.go` is still authoritative for `go doc`; README is for humans browsing the filesystem. |
| Directory layout changes | The README's directory map must be updated as part of the same change. Treated as part of "Adding a new embedded asset" checklist. |
| Quality gates added later | The "Quality gates today" section must be updated when a new gate lands (e.g. when the TS type-check task ships, move TypeScript from "not gated" to "gated"). |

### Validation Rules

This is a documentation change. The validation is editorial:
claims in the README must match the actual `embed.go`
directives and the actual `embed_test.go` coverage.

### Error Handling

Not applicable (documentation file).

## Interface

Not applicable. README is a static Markdown file at
`internal/assets/README.md`. No CLI, no skill, no flag.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/assets/README.md` | New file. |
| `.context/TASKS.md` | Redirect Phase 0 line-30 task: original "add TS type-check" becomes "completed via README + spawned four gap tasks"; add four new Phase 0 items for the known gaps (TS type-check, shellcheck, PSScriptAnalyzer, skill frontmatter validation). |
| `.context/DECISIONS.md` | Add entry: "Embedded foreign-language assets under `internal/assets/` is intentional, not a smell; the lift-out instinct is wrong because of `//go:embed` no-`../`." |

### Helpers to Reuse

None. Single file authored from scratch using the project's
existing Markdown conventions.

## Configuration

None.

## Testing

- Editorial review: every claim about embed coverage cross-
  checked against `embed.go`.
- Editorial review: every claim about test coverage cross-
  checked against `embed_test.go`.
- `make build` and `make test` must pass (the README itself
  is not loaded by Go; this is a sanity gate that no nearby
  changes broke anything).

## Non-Goals

- This spec does **not** restructure `internal/assets/`. The
  tree stays where it is; only the contract document changes.
- This spec does **not** implement any of the four gap items
  (TS type-check, shellcheck, PSScriptAnalyzer, frontmatter
  validation). Those are spawned as separate Phase 0 tasks.
- This spec does **not** add dev tooling files for the
  embedded TypeScript plugin. That work belongs to the TS
  type-check task once it is picked up.
- This spec does **not** modify `doc.go`. `doc.go` continues
  to serve the Go-doc audience; the README serves the
  filesystem-browsing audience.

## Open Questions

None.

## Implementation Outcome (2026-05-11)

After the README landed, a follow-up clarification rejected the
"redirect by spawning four follow-up tasks" framing as a thinly-
disguised deferral: the original Phase 0 line-30 task ("add
`tsc --noEmit` for the embedded OpenCode plugin") still had to
be implemented, not just rescoped. The implementation was done
in the same session:

- `tools/typecheck/opencode/` — new top-level dev-tooling
  directory (sibling to `internal/assets/`, honoring the embed
  contract). Contains `package.json` (`@opencode-ai/plugin`,
  `@types/bun`, `typescript`), `tsconfig.json` (with `baseUrl`
  + `paths` to resolve modules from the tools dir despite the
  source file living elsewhere), `package-lock.json` (npm,
  matching `editors/vscode/`), and a local `README.md`.
- `.github/workflows/ci.yml` — new job
  `typecheck-opencode-plugin` running `npm ci` then
  `npx tsc --noEmit` against the embedded TS source.
- `.gitignore` — `tools/typecheck/opencode/node_modules/`.
- TASKS.md — original line-30 task marked `[x]` complete;
  duplicate "spawned" item marked `[-]` as duplicate of the
  now-completed original.

Three of the four originally-spawned gap tasks remain genuinely
new and stay open: shellcheck for embedded Bash, PSScriptAnalyzer
for embedded PowerShell, and skill-frontmatter validity tests.
