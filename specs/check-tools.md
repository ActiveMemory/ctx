# Spec: tooling dependency version check

## Problem

The project depends on a spread of developer tooling — the Go
toolchain, golangci-lint, Node (for GitNexus and npx-launched MCP
servers), pipx + zensical (site builds), the `claude` CLI with
registered MCP servers (gitnexus, gemini-search) — but nothing
verifies any of it. Gaps surface mid-task instead of up front:
a doc-link fix session discovered zensical *and* pipx were both
missing only after the source edits were done and the site needed
rebuilding.

`make gitnexus-version` covers exactly one tool. There is no
single command that answers "is this machine ready to work on
ctx?"

A second failure class is worse than absence: tooling that is
present but *wrong*. A `make site` run with a zensical 14 versions
behind the one that produced the committed `site/` rewrote all 178
pages with stale generator output — a diff indistinguishable from a
real docs change, on a tree CI never rebuilds. Presence-only checks
cannot see this, and a checker nobody runs before `make site`
cannot prevent it.

## Design

Three pieces:

1. **`hack/tool-versions.txt`** — the manifest. One line per
   tool: name, type, requirement, minimum version. Version
   minimums live here, not in script logic, so bumps are
   one-line diffs.
2. **`hack/check-tools.sh`** — the checker. Reads the manifest,
   probes each tool, prints one aligned table with a verdict per
   row, exits non-zero only if a *required* tool fails.
3. **`make check-tools`** — the entry point developers actually
   remember.

### Tool types

| Type  | Probe                                                       |
|-------|-------------------------------------------------------------|
| `bin` | `command -v` + version extraction (`--version` / `go version`) |
| `npm` | installed version vs `npm view <pkg> version` (drift check) |
| `mcp` | registered in `claude mcp list` (config presence)           |

### Verdicts

- `OK` — present, and ≥ minimum (or == exact pin) when declared
- `MISSING` — not found on PATH / not registered
- `OUTDATED` — below declared minimum or exact pin, or (npm type)
  behind the registry's latest
- `DRIFT` — *ahead* of an exact pin. Only reachable for `=`-pinned
  rows; a floor has no upper bound to exceed
- `SKIP` — probe impossible (no `claude` CLI for mcp rows, no
  network for npm registry) — reported, never fatal

### Version column: floor vs exact pin

| Manifest value | Semantics | Below | Above    |
|----------------|-----------|-------|----------|
| `1.26`         | minimum   | `OUTDATED` | `OK`  |
| `=0.0.51`      | exact pin | `OUTDATED` | `DRIFT` |
| `-`            | presence  | —     | —        |

A floor fits tools whose newer versions are strictly fine (`go`,
`node`). An exact pin fits tools whose *output is committed*, where
any delta rewrites a tracked artifact. `DRIFT` and `OUTDATED` both
route through the normal failure tally, so they warn on an
`optional` row and turn fatal under `--strict` — which is what makes
a build gate out of them.

### Requirement levels

`required` failures exit 1 (go, git, golangci-lint — the
build/lint core). `optional` failures print as warnings and do
not affect the exit code. `--strict` promotes optional failures
to fatal for CI-style use.

### Build-target preconditions

`--only <name>` restricts the run to one manifest row. Combined
with `--strict` it turns any row — including an `optional` one —
into a hard precondition a Make target can depend on:

```make
site: site-guard
site-guard:
	@./hack/check-tools.sh --only zensical --strict || { ...; exit 1; }
```

This is how `make site` refuses to build with a stale generator.
The target supplies the remediation copy; the script supplies the
verdict. An `--only` name absent from the manifest exits 1 rather
than reporting a vacuous success, so a typo cannot disarm a gate.

### Decisions

- **MCP checks are config-presence, not connectivity.** `claude
  mcp list` already health-checks; parsing its status marks is
  brittle across CLI versions. Presence answers the actual
  question ("did I register it?") deterministically.
- **Network probes are best-effort.** `npm view` runs under a
  timeout and degrades to `SKIP`-style notes; a checker that
  hangs offline would never get run.
- **FireCrawl MCP ships as a commented manifest line** — enabling
  it later is a one-character diff, per the original request.
- **The zensical pin is an equality, not a floor.** A floor only
  catches half the failure: it blocks an older generator but waves
  a newer one straight through, silently, into the same churn diff.
  For a tool whose *output is a committed artifact*, any version
  delta is reportable, so the manifest's `=` prefix makes the check
  two-sided. The pin records the generator stamped into
  `site/index.html`; moving it is a deliberate act that travels in
  the same commit as the regenerated `site/`.
- **The diff size is not predictable and must not be asserted.**
  One observed 0.0.33-vs-0.0.47 rebuild touched 178 files, but the
  blast radius depends on the version gap — a template tweak may
  touch a handful of pages, an asset-hash change touches everything
  referencing it. The gate's justification is that the churn is
  *unintended*, not that it is large.
- **Only the manifest holds the version.** `make site-setup`
  derives its install constraint from it via `ZENSICAL_MIN`, so
  the installer and the gate cannot disagree — the drift class
  that produced the original bug.
- **`site-serve` and `journal` stay ungated.** Neither writes
  tracked artifacts (`.context/journal-site/` is gitignored), and
  gating `make journal` would block session *import*, which needs
  no site generator at all.

## Non-Goals

- Auto-installing or upgrading anything (that stays with
  `make site-setup`, `make gitnexus-update`, `make register-mcp`)
- Verifying MCP server *liveness*
- Gating CI on optional tooling. Unchanged: `--only … --strict`
  gates a *local build target* on the tool that target shells out
  to, which is a precondition, not a CI policy. CI runs neither
  `check-tools` nor `make site`.

## Acceptance

- [ ] `make check-tools` prints one row per manifest entry with
      verdict, version, and note
- [ ] Exit 0 on a machine with go/git/golangci-lint present, even
      when optional tools are missing
- [ ] Exit 1 when a required tool is missing or below minimum
- [ ] Manifest comment lines and blank lines are ignored
- [ ] Absence of `claude`, `npm`, or network yields SKIP notes,
      never a hang or crash
- [ ] `--only <name>` prints exactly the matching row(s); with
      `--strict`, an `optional` row's failure exits 1
- [ ] `--only <unknown>` exits 1 with a diagnostic, never 0
- [ ] `make site` exits non-zero without invoking `zensical build`
      when the local generator is *either* below or above the pin,
      and leaves `site/` byte-identical
- [ ] A generator matching the pin exactly yields `OK` and the build
      proceeds
- [ ] `make site-setup` installs the constraint recorded in the
      manifest, with the version written in exactly one place
