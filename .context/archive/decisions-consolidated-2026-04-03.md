# Archived Decisions (consolidated 2026-04-03)

Originals replaced by consolidated entries in DECISIONS.md.

## Group: Output/write location

## [2026-03-22-084509] No runtime pluralization — use singular/plural text key pairs

**Status**: Accepted

**Context**: Hardcoded English plural rules (+ s, y → ies) were scattered across format.Pluralize, padPluralize, and inline code — all i18n dead-ends

**Decision**: No runtime pluralization — use singular/plural text key pairs

**Rationale**: Different languages have vastly different plural rules. Complete sentence templates with embedded counts (time.minute-count '1 minute', time.minutes-count '%d minutes') let each locale define its own plural forms.

**Consequence**: format.Pluralize and format.PluralWord are deleted. All plural output uses paired text keys with the count embedded in the template.

---

## [2026-03-21-084020] Output functions belong in write/, logic and types in core/

**Status**: Accepted

**Context**: PrintFeedReport was initially placed in cli/site/core/ but it calls cmd.Println — that's output formatting, not business logic

**Decision**: Output functions belong in write/, logic and types in core/

**Rationale**: The project taxonomy separates concerns: core/ owns domain logic, types, and helpers; write/ owns CLI output formatting that takes *cobra.Command for Println. Mixing them blurs the boundary and makes testing harder.

**Consequence**: All functions that call cmd.Print/Println/Printf belong in the write/ package tree. core/ never imports cobra for output purposes.

---


## Group: YAML text externalization

## [2026-04-03-133236] YAML text externalization justification is agent legibility, not i18n

**Status**: Accepted

**Context**: Principal analysis initially framed 879-key YAML externalization as a bet on i18n. Blog post review (v0.8.0) revealed the real justification: agent legibility (named DescKey constants as traversable graphs), drift prevention (TestDescKeyYAMLLinkage catches orphans mechanically), and completing the archaeology of finding all 879 scattered strings.

**Decision**: YAML text externalization justification is agent legibility, not i18n

**Rationale**: The v0.8.0 blog makes it explicit: finding strings is the hard part, translation is mechanical. The externalization already pays for itself through agent legibility and mechanical verification. i18n is a free downstream consequence, not the justification.

**Consequence**: Future architecture analysis should frame externalization as already-justified investment. The 3-file ceremony (DescKey + YAML + write/err function) is the cost of agent-legible, drift-proof output — not speculative i18n prep.

---

## [2026-03-15-101336] TextDescKey exhaustive test verifies all 879 constants resolve to non-empty YAML values

**Status**: Accepted

**Context**: PR #42 merged with ~80 new MCP text keys but no test coverage for key-to-YAML mapping

**Decision**: TextDescKey exhaustive test verifies all 879 constants resolve to non-empty YAML values

**Rationale**: A single table-driven test parsing embed.go source catches typos and missing YAML entries at test time — no manual key list to maintain

**Consequence**: New TextDescKey constants are automatically covered; orphaned keys fail CI

---

## [2026-03-15-040638] Split text.yaml into 6 domain files loaded via loadYAMLDir

**Status**: Accepted

**Context**: text.yaml grew to 1812 lines covering write, errors, mcp, doctor, hooks, and ui domains

**Decision**: Split text.yaml into 6 domain files loaded via loadYAMLDir

**Rationale**: Matches existing split pattern (commands.yaml, flags.yaml, examples.yaml); loadYAMLDir merges all files in commands/text/ transparently so TextDesc() API stays unchanged

**Consequence**: New domain files must go into commands/text/; loadYAMLDir reads all .yaml in the directory at init time

---

## [2026-03-12-133007] Split commands.yaml into 4 domain files

**Status**: Accepted

**Context**: Single 2373-line YAML mixed commands, flags, text, and examples with inconsistent quoting

**Decision**: Split commands.yaml into 4 domain files

**Rationale**: Context is for humans — localization files should be human-readable block scalars. Separate files eliminate the underscore prefix namespace hack

**Consequence**: 4 files (commands.yaml, flags.yaml, text.yaml, examples.yaml) with dedicated loaders in embed.go

---

## [2026-03-06-200257] Externalize all command descriptions to embedded YAML for i18n readiness — commands.yaml holds Short/Long for 105 commands plus flag descriptions, loaded via assets.CommandDesc() and assets.FlagDesc()

**Status**: Accepted

**Context**: Command descriptions were inline strings scattered across 105 cobra.Command definitions

**Decision**: Externalize all command descriptions to embedded YAML for i18n readiness — commands.yaml holds Short/Long for 105 commands plus flag descriptions, loaded via assets.CommandDesc() and assets.FlagDesc()

**Rationale**: Centralizing user-facing text in a single translatable file prepares for i18n without runtime cost (embedded at compile time)

**Consequence**: System's 30 hidden hook subcommands excluded (not user-facing); flag descriptions use _flags.scope.name convention

---


## Group: Package taxonomy and code placement

## [2026-03-06-200247] cmd/root + core taxonomy for all CLI packages — single-command packages use cmd/root/{cmd.go,run.go}, multi-subcommand packages use cmd/<sub>/{cmd.go,run.go}, shared helpers in core/

**Status**: Accepted

**Context**: 35 CLI packages had inconsistent flat structures mixing Cmd(), run logic, helpers, and types in the same directory

**Decision**: cmd/root + core taxonomy for all CLI packages — single-command packages use cmd/root/{cmd.go,run.go}, multi-subcommand packages use cmd/<sub>/{cmd.go,run.go}, shared helpers in core/

**Rationale**: Taxonomical symmetry: every package has the same predictable shape, making navigation instant and agent-friendly

**Consequence**: cmd/ contains only cmd.go + run.go; helpers go to core/; 474 files changed in initial restructuring

---

## [2026-03-06-200227] Shared entry types and API live in internal/entry, not in CLI packages — domain types that multiple packages consume (mcp, watch, memory) belong in a domain package, not a CLI subpackage

**Status**: Accepted

**Context**: External consumers were importing cli/add for EntryParams/ValidateEntry/WriteEntry, creating a leaky abstraction

**Decision**: Shared entry types and API live in internal/entry, not in CLI packages — domain types that multiple packages consume (mcp, watch, memory) belong in a domain package, not a CLI subpackage

**Rationale**: Domain types in CLI packages force consumers to depend on CLI internals; internal/entry provides a clean boundary

**Consequence**: entry aliases Params from add/core to avoid import cycle (entry imports add/core for insert logic); future work may move insert logic to entry to eliminate the cycle

---

## [2026-03-13-151954] Templates and user-facing text live in assets, structural constants stay in config

**Status**: Accepted

**Context**: Ongoing refactoring session moving Tpl* constants out of config/

**Decision**: Templates and user-facing text live in assets, structural constants stay in config

**Rationale**: config/ is for structural constants (paths, limits, regexes); assets/ is for templates, labels, and text that would need i18n. Clean separation of concerns

**Consequence**: All tpl_entry.go, tpl_journal.go, tpl_loop.go, tpl_recall.go moved to assets/

---


## Group: Eager init over lazy loading

## [2026-03-18-193631] Eager Init() for static embedded data instead of per-accessor sync.Once

**Status**: Accepted

**Context**: 4 sync.Once guards + 4 exported maps + 4 Load functions + a wrapper package for YAML that never mutates.

**Decision**: Eager Init() for static embedded data instead of per-accessor sync.Once

**Rationale**: Data is static and required at startup. sync.Once per accessor is cargo cult. One Init() in main.go is sufficient. Tests call Init() in TestMain.

**Consequence**: Maps unexported, accessors are plain lookups, permissions and stopwords also loaded eagerly. Zero sync.Once remains in the lookup pipeline.

---

## [2026-03-16-104143] Explicit Init over package-level init() for resource lookup

**Status**: Accepted

**Context**: server/resource package used init() to silently build the URI lookup map

**Decision**: Explicit Init over package-level init() for resource lookup

**Rationale**: Implicit init hides startup dependencies, makes ordering unclear, and is harder to test. Explicit Init called from NewServer makes the dependency visible.

**Consequence**: res.Init() called explicitly from NewServer before ToList(); no package-level side effects

---


## Group: Pure logic separation of concerns

## [2026-03-15-230640] Pure-logic CompactContext with no I/O — callers own file writes and reporting

**Status**: Accepted

**Context**: MCP server and CLI compact command both implemented task compaction independently, with the MCP handler using a local WriteContextFile wrapper

**Decision**: Pure-logic CompactContext with no I/O — callers own file writes and reporting

**Rationale**: Separating pure logic from I/O lets both MCP (JSON-RPC responses) and CLI (cobra cmd.Println) callers control output and file writes. Eliminates duplication and the unnecessary mcp/server/fs package

**Consequence**: tidy.CompactContext returns a CompactResult struct; callers iterate FileUpdates and write them. Archive logic stays in callers since MCP and CLI have different archive policies

---

## [2026-03-16-122033] Server methods only handle dispatch and I/O, not struct construction

**Status**: Accepted

**Context**: MCP server had ok/error/writeError as methods plus prompt builders that didn't use Server state — they just constructed response structs

**Decision**: Server methods only handle dispatch and I/O, not struct construction

**Rationale**: Methods that don't access receiver state hide their true dependencies and inflate the Server interface. Free functions make the dependency graph explicit and are independently testable.

**Consequence**: New response helpers go in server/out, prompt builders in server/prompt. Server methods are limited to dispatch (handlePromptsGet) and I/O (writeJSON, emitNotification). Same principle applies to future tool/resource builders.

---

## [2026-03-23-003346] Pure-data param structs in entity — replace function pointers with text keys

**Status**: Accepted

**Context**: MergeParams had UpdateFn callback, DeployParams had ListErr/ReadErr function pointers — both smuggled side effects into data structs

**Decision**: Pure-data param structs in entity — replace function pointers with text keys

**Rationale**: Text keys are pure data, keep entity dependency-free, and the consuming function can do the dispatch itself

**Consequence**: All cross-cutting param structs in entity must be function-pointer-free; I/O functions passed as direct parameters

---


