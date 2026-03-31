# CLI Taxonomy Redesign

**Status**: Approved
**Date**: 2026-03-26
**Author**: Jose + Claude

## Problem

ctx has 32 visible top-level commands in a flat, alphabetically-sorted list.
Users see a wall of text and must parse each description to find what they need.
Related commands (`recall`, `journal`, `memory`) sit far apart alphabetically.
The skill `/ctx-remember` and the command `ctx recall` sound like synonyms but
serve different purposes. There is no visual hierarchy distinguishing daily
drivers from plumbing.

### Symptoms

1. **Cognitive overload**: 32 commands in `ctx --help` is too many to scan
2. **Semantic collisions**: `recall` (browse history) vs. `/ctx-remember`
   (context readback); `journal` (site gen) vs. `recall import` (session ingest)
3. **No workflow guidance**: help output is alphabetical, not task-oriented
4. **Audience mixing**: user commands (`status`, `add`) sit next to automation
   plumbing (`watch`, `loop`, `mcp`)
5. **Orphaned commands**: `dep`, `reindex`, `why` are rarely referenced in
   skills, recipes, or playbook

### What users actually run

Usage analysis across 44 skills, 34 recipes, playbook, and Makefiles:

| Rank | Command  | Refs | Primary audience |
|------|----------|------|------------------|
| 1    | pad      | 81   | User             |
| 2    | system   | 80   | Automation/hooks |
| 3    | recall   | 69   | Mixed            |
| 4    | add      | 67   | Mixed            |
| 5    | init     | 36   | User (one-time)  |
| 6    | status   | 30+  | Mixed            |
| 7    | drift    | 25+  | Mixed            |
| 8    | agent    | 20+  | Automation       |
| 9    | journal  | 20+  | Mixed            |
| 10   | memory   | 15+  | Automation       |

Commands with zero external references: several skill-only commands that are
invoked through skills, not directly.

## Prior Art

### Patterns from established CLIs

| Tool      | Cmds | Strategy                  | Nesting | Key insight                              |
|-----------|------|---------------------------|---------|------------------------------------------|
| kubectl   | ~50  | 8 workflow groups          | None    | Groups in help output only (Cobra AddGroup) |
| gh        | 44   | 4 groups (Core + buckets)  | noun-verb | "Most important first" with catch-all   |
| docker    | 60+  | Management + Legacy        | Nested  | Migration path: both forms work forever  |
| git       | 169  | Curated 30 / full 169      | None    | Two views: default shows workflow subset |
| terraform | ~20  | 5 Main + All Other         | Minimal | Core workflow in execution order          |

### Key takeaways

- **Groups are presentation, not structure**: Cobra's `AddGroup()` changes
  help output without changing command routing or breaking existing invocations
- **Nesting is a one-way door**: `docker container ls` coexists with `docker ps`
  but the alias layer adds permanent maintenance burden
- **Curation beats reorganization**: git shows 30 of 169 commands by default;
  terraform highlights 5 of 20. The "less is more" approach works better than
  restructuring
- **Task-oriented labels**: git uses "start a working area", "examine the
  history and state" — verbs and goals, not category nouns

## Design

### Approach: Cobra command groups (presentation-only)

Use Cobra's `AddGroup()` API to organize commands into labeled sections in
`--help` output. This:

- Changes **zero** command paths or invocations
- Requires no migration, aliases, or deprecation
- Affects only the help template rendering
- Is the same mechanism used by kubectl, gh, and other Cobra-based CLIs

No nesting. No renames. No breaking changes.

### Approved layout

The Long description of the root command includes a prominent preamble box
with the critical ceremony path. Developers never RTFM — the help command is
the battlefield to counter that.

```
ctx - persistent context for AI coding assistants
──────────────────────────────────────────────────────
!  START: /ctx-remember                              !
!    END: /ctx-wrap-up                               !
!                                                    !
!  MAINTENANCE (every 2–3 sessions):                 !
!         ctx recall import --all                    !
!         /ctx-journal-enrich-all                    !
!                                                    !
!  Skip this, and you will suffer.                   !
!                                                    !
!  ctx is not a skill: ctx is the persistence layer. !
──────────────────────────────────────────────────────

Usage:
  ctx [command]

Getting Started:
  init          Initialize a new .context/ directory with template files
  status        Show context summary with token estimate
  guide         Quick-reference cheat sheet for ctx

Context (source of truth):
  add           Add a new item to a context file
  load          Output assembled context Markdown
  agent         Print AI-ready context packet
  sync          Reconcile context with codebase
  drift         Detect stale or invalid context
  compact       Archive completed tasks and clean up context

Artifacts (.context/ files):
  decision      Manage DECISIONS.md
  learning      Manage LEARNINGS.md
  task          Manage task archival and snapshots

Sessions:
  recall        Browse and search AI session history
  journal       Analyze and synthesize imported sessions
  memory        Bridge Claude Code auto memory into .context/
  remind        Session-scoped reminders
  pad           Encrypted scratchpad for sensitive notes

Runtime:
  config        Manage runtime configuration
  permission    Manage permission snapshots
  pause         Pause context hooks for this session
  resume        Resume context hooks for this session

Integration:
  hook          Generate AI tool integration configs
  mcp           Model Context Protocol server
  watch         Watch for context-update commands in AI output
  notify        Send a webhook notification
  loop          Generate a Ralph loop script

Diagnostics:
  doctor        Structural health check
  change        Show what changed since last session
  dep           Show package dependency graph
  why           Read the philosophy behind ctx

Site / Output:
  serve         Serve a static site locally via zensical
  site          Site management commands

Utilities:
  reindex       Regenerate indices for DECISIONS.md and LEARNINGS.md
  completion    Generate shell autocompletion script
  help          Help about any command

Flags:
      --allow-outside-cwd   Allow context directory outside project root
      --context-dir string  Override context directory path
  -h, --help                help for ctx
  -v, --version             version for ctx

Use "ctx [command] --help" for more information about a command.
```

### Groups (10 total)

| Group ID           | Title                       | Commands                                    |
|--------------------|-----------------------------|---------------------------------------------|
| getting-started    | Getting Started:            | init, status, guide                         |
| context            | Context (source of truth):  | add, load, agent, sync, drift, compact      |
| artifacts          | Artifacts (.context/ files):| decision, learning, task                    |
| sessions           | Sessions:                   | recall, journal, memory, remind, pad        |
| runtime            | Runtime:                    | config, permission, pause, resume           |
| integration        | Integration:                | hook, mcp, watch, notify, loop              |
| diagnostics        | Diagnostics:                | doctor, change, dep, why                    |
| site-output        | Site / Output:              | serve, site                                 |
| utilities          | Utilities:                  | reindex, completion, help                   |

`system` is hidden and gets no GroupID — it won't appear in help output.

### Design rationale

**Preamble box**: The most important information in the entire CLI sits at
the top of `--help`. Three things every user must know: start ceremony, end
ceremony, maintenance cadence. The tone is deliberately direct.

**Group ordering** follows the session lifecycle from the playbook:
Load → Orient → Pick → Work → Commit → Reflect. "Getting Started" maps to
setup, "Context" to the daily work loop, "Sessions" to reflection, etc.

**Group size**: 3–6 commands per group. The largest group (Context) has 6,
which matches the core daily workflow.

**"Artifacts"** separates the three .context/ file managers (decision,
learning, task) from the broader context commands. These are CRUD operations
on specific files.

**"Sessions"** groups the five commands related to session history and
ephemeral session data: `recall` (browse/import), `journal` (publish),
`memory` (auto-memory bridge), `remind` (session-scoped), `pad`
(session-scoped scratchpad).

**"Runtime"** groups config and session hook control. `hook` moved to
Integration because it generates configs for external tools, not for ctx
itself.

**"Site / Output"** separates static site concerns from the rest — these
are about publishing, not about context management.

**"Utilities"** replaces Cobra's default "Additional Commands" label for
the catch-all bucket.

**`system` is not listed** — it's already hidden (all its subcommands are
hidden plumbing). It won't appear in grouped help output either.

### The remember vs. recall naming collision

| Current name          | What it does                                    |
|-----------------------|-------------------------------------------------|
| `ctx recall`          | Browse/import AI session history (CLI command)   |
| `/ctx-remember`       | Context readback ceremony (skill)                |

These sound like synonyms but serve entirely different purposes. Options:

**Option A — Rename `recall` to `session`**:
- `ctx session list`, `ctx session import`, `ctx session show`
- Pro: "session" is unambiguous — it's about session transcripts
- Con: touches the command path, requiring migration across docs/skills/hooks

**Option B — Rename the skill `/ctx-remember` to `/ctx-readback`**:
- Pro: skills are internal, no public API surface
- Con: "readback" is jargon; "remember" is the natural user trigger phrase

**Option C — Do nothing; let groups disambiguate**:
- With "Session History" grouping, `recall` sits next to `journal` and
  `memory`, making its purpose clear from context
- The skill `/ctx-remember` is invoked by phrase ("do you remember?"), not
  by browsing help output
- Pro: zero changes, zero migration
- Con: the collision persists in documentation

**Recommendation: Option C now, revisit if users report confusion.** The
grouping provides enough context. The naming collision is between a CLI
command and a skill — two different invocation surfaces that rarely overlap.
If we later want Option A, the migration is the same shape as the
`export` → `import` rename we just completed.

## Implementation

### Phase 1: Add command groups (presentation-only)

**Scope**: `internal/bootstrap/bootstrap.go` + constants + YAML

1. Define group constants in `internal/config/embed/cmd/group.go`:
   ```go
   const (
       GroupGettingStarted = "getting-started"
       GroupContext         = "context"
       GroupArtifacts      = "artifacts"
       GroupSessions       = "sessions"
       GroupRuntime        = "runtime"
       GroupIntegration    = "integration"
       GroupDiagnostics    = "diagnostics"
       GroupSiteOutput     = "site-output"
       GroupUtilities      = "utilities"
   )
   ```

2. Define group titles in YAML (for i18n readiness):
   ```yaml
   group.getting-started:
     short: "Getting Started:"
   group.context:
     short: "Context (source of truth):"
   group.artifacts:
     short: "Artifacts (.context/ files):"
   group.sessions:
     short: "Sessions:"
   group.runtime:
     short: "Runtime:"
   group.integration:
     short: "Integration:"
   group.diagnostics:
     short: "Diagnostics:"
   group.site-output:
     short: "Site / Output:"
   group.utilities:
     short: "Utilities:"
   ```

3. Update the root command Long description (in commands.yaml) to include
   the preamble box with ceremony instructions.

4. Register groups in `bootstrap.go` and set GroupID on each command.
   Keep the mapping centralized in bootstrap.go — the taxonomy should be
   visible in one place, not scattered across 32 packages.

**Files changed**: ~5 (bootstrap.go, group constants, YAML, possibly a
custom help template for spacing)

**Zero breaking changes**: all existing `ctx <command>` invocations work
identically.

### Phase 2: Custom help template (optional polish)

Cobra's default grouped help output is functional but could be improved:

- Add blank lines between groups for readability
- Optionally show a one-line "tip" at the bottom:
  `Run "ctx guide" for workflows and recipes.`
- Consider matching terraform's style of bolding the group header

This is cosmetic and can ship separately.

### Phase 3: Update documentation

Per the rename/refactor documentation checklist (CONVENTIONS.md):

1. **Docstrings**: None needed — groups don't change command signatures
2. **User-facing docs** (`docs/`): Update `docs/cli/` if it shows help output
3. **Recipes**: No changes unless recipes reference help output formatting
4. **Skills**: No changes — skills invoke commands by name, not by group
5. **Blog**: Consider a short post about the new help organization

### Phase 4: Revisit session/recall naming (future, if needed)

If Option A (rename `recall` → `session`) is pursued later:
- Same migration shape as `export` → `import`
- Touch count: ~40-60 files (similar to this rename)
- Would also rename `internal/cli/recall/` → `internal/cli/session/`

## Non-goals

- **Command nesting**: No `ctx context add` or `ctx history recall`. Flat
  commands with groups is the right complexity level for 32 commands.
- **Command removal**: No commands are removed. Even rarely-used ones like
  `dep` and `why` stay — they're just grouped into "Diagnostics."
- **Aliases**: No Docker-style aliasing. The command surface is small enough
  that aliases add complexity without proportional value.
- **Tab completion changes**: Groups don't affect completion. All commands
  remain top-level completions.

## Risks

- **Cobra version**: `AddGroup()` requires Cobra v1.6.0+. Current project
  uses Cobra — verify version in `go.mod`.
- **CI/test impact**: Tests that assert on help output (e.g.,
  `TestGuideDefaultOutput`) may need updating if the guide command references
  the grouped help structure.
- **Third-party docs**: External tutorials or blog posts referencing the old
  help output won't break, just look different.

## Verification

- [ ] `ctx --help` shows grouped output with correct section ordering
- [ ] `ctx <command>` still works for every command (no routing changes)
- [ ] `ctx completion bash/zsh/fish` still works
- [ ] `ctx help <command>` still works
- [ ] Hidden commands (`system` subcommands) don't appear in any group
- [ ] `make lint && make test` pass
- [ ] Bootstrap test (`TestInitializeSubcommandCount`) still passes
