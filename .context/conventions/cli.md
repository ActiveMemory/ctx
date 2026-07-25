# cli

## CLI Structure

- **CLI package taxonomy**: Every package under `internal/cli/` follows:
  parent.go (Cmd wiring), doc.go, `cmd/root/` or `cmd/<sub>/`
  (implementation), `core/` (shared helpers)
- **cmd/ directories**: Only cmd.go, run.go, and tests — helpers and
  output go to `core/`
- **core/ structs**: Consolidated into a single `types.go` file
- **User-facing text via assets**: All text routed through
  `internal/assets` with YAML-backed TextDescKeys — no inline strings
  in `core/` or `cmd/` packages
- **config/ doc.go**: Every package under `internal/config/` must have
  a doc.go with the project header and a one-line package comment
- **DescKey prefix**: Not CmdDescKey — `cmd.DescKeyFoo` not
  `cmd.CmdDescKeyFoo` (Go package hygiene, avoids stutter)
- **Cobra Use: fields**: Must reference `cmd.Use*` constants, never raw
  strings or `cmd.DescKey*`
- **Run functions exported PascalCase**: `Run`, `RunImport`,
  `RunArchive` etc. No private `runXXX` variants
- **write/ packages write to stdio only**: Functions take
  `*cobra.Command`, not `io.Writer`. Exception: `write/rc` writes to
  `os.Stderr` because rc loads before cobra
- **Package directory names singular**: Unless Go convention requires
  plural
- **Import grouping**: stdlib — blank line — external deps (cobra,
  yaml) — blank line — ctx imports. Three groups, always in this order
- **camelCase import aliases**: `cFlag` not `cflag`, `cfgFmt` not
  `cfgfmt`
- **Icons and symbols as token constants**: Not unicode escapes
- **Cross-cutting domain types in internal/entity**: Types used by one
  package stay in that package; types used across packages go to entity

- Warn format strings centralized in config/warn/ — use warn.Close,
  warn.Write, warn.Remove, warn.Mkdir, warn.Rename, warn.Walk, warn.Getwd,
  warn.Readdir, warn.Marshal instead of inline format strings in log.Warn calls

- Nav frontmatter title: fields must not contain ctx — frontmatter does not
  support backticks, so the brand stays out of nav titles entirely (Hub, not The
  ctx Hub). Body headings can use `ctx` since markdown supports backticks.

- CLI flags and slash-commands inside headings or admonition titles must be
  backticked: `--keep-frontmatter=false`, `/ctx-reflect`. The title-case engine
  in hack/title-case-headings.py protects these patterns automatically, but
  authors should still backtick at write time for clarity.

- File extensions inside headings must be backticked when title-case
  capitalization would otherwise apply: write `CONSTITUTION.md`, not
  CONSTITUTION.Md. The title-case engine refuses to capitalize lowercase tokens
  following a literal . dot, but explicit backticks remain the clearest signal.
- New editor integrations include an MCP-merge test covering: create / empty
  file / preserve existing keys / skip when registered / reject malformed JSON

- Substrate vs. artifact placement: cognitive substrate (consumed and mutated
  via ctx-mediated paths — `ctx agent`, `ctx decision add`, `/ctx-kb-ingest`,
  `/ctx-handover`, ceremonies) lives under `.context/`; project artifacts (read
  and edited directly by humans — `specs/`, `CLAUDE.md`, `GETTING_STARTED.md`,
  `docs/`) live at the project root; tool config and tool homes (`.ctxrc`,
  `.claude/`) live at root by dotfile/tool convention. The kb is substrate, not
  artifact: direct file edits remain possible per Invariant 1, but the
  skill-mediated path is the discipline. Rationale recorded in DECISIONS.md.

## User-Facing Surface Completeness

When a change adds or alters a user-facing surface — a new
`ctx` subcommand, a new flag, an observable behavior change,
a new exit shape, a new output line — the work is **not
complete** until every one of the following has been updated
in the same commit (or the same stacked PR, with the user's
explicit OK):

- `internal/assets/commands/commands.yaml` and
  `examples.yaml` for the subcommand description and example
- `internal/assets/claude/skills/ctx-<area>/SKILL.md` so the
  agent knows the surface exists and when to trigger it
- `internal/assets/integrations/copilot-cli/skills/<...>` if
  a parallel skill exists for the integration
- `docs/recipes/<related-recipe>.md` for any recipe that
  already demonstrates the broader feature; consider a new
  recipe if the surface is its own workflow shape
- `docs/cli/<command>.md` if a per-command CLI doc page
  exists for this surface

Splitting these into a "Phase 2 / follow-up commit / future
sweep" is **deferral** in the Constitution's sense, no matter
how the phase is labeled. Docs are part of the deliverable,
not a separable improvement. The "I can create a follow-up
task" prohibition applies verbatim.

Acceptable exceptions (state them in the commit body):

- The surface is internal-only (no human user encounters it).
- A recipe / skill genuinely does not exist for this feature
  area and writing one is itself a larger separable piece of
  work (then file the spec for that piece in the same commit,
  do not just defer).

The Self-check before declaring a feature commit complete is:
*"If a user runs `ctx help` or asks `/ctx-<area>` to do this
new thing today, will the help text / skill / recipe match
what the code does?"* If no, the commit is not complete.

