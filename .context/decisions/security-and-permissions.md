# security-and-permissions

## [2026-03-14-110748] System path deny-list as safety net, not security boundary

**Status**: Accepted

**Context**: Replacing nolint:gosec directives with centralized I/O wrappers in
internal/io

**Decision**: System path deny-list as safety net, not security boundary

**Rationale**: ctx paths are internally constructed from config constants. The
deny-list catches agent hallucinations (writing to /etc), not adversarial input.
Public security docs would imply a threat model that does not exist.

**Consequence**: internal/io/doc.go documents limitations honestly for
contributors. No user-facing security docs. The deny-list is a modicum of
protection, not a promise.

---

## [2026-02-26-100006] Security and permissions (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-21 to 2026-02-24)

- Keep CONSTITUTION.md minimal: only truly inviolable rules (security,
  correctness, process invariants). Style preferences go in CONVENTIONS.md.
  Overly strict constitution gets ignored.
- Centralize constants with semantic prefixes in `internal/config/config.go`:
  `Dir*` for directories, `File*` for paths, `Filename*` for names,
  `UpdateType*` for entry types. Single source of truth, compile-time typo
  checks.
- Hooks use `ctx` from PATH, not hardcoded absolute paths. Standard Unix
  practice; portable across machines/users. `ctx init` checks PATH availability
  before proceeding.
- Drop absolute-path-to-ctx regex from block-dangerous-commands shell script.
  The block-non-path-ctx Go subcommand already covers this with better patterns;
  duplicating creates two sources of truth.

---

