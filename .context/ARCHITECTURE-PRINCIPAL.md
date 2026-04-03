# Architecture - Principal Analysis

_Generated 2026-04-03. Strategic analysis only; see ARCHITECTURE.md
for the authoritative architecture reference._

## Current State Summary

ctx is a 527-package Go CLI (34 commands) with an MCP server, built
around a `.context/` markdown directory that provides structured,
token-budgeted context for AI coding assistants. The codebase
recently underwent a major restructuring: flat packages exploded
into layered taxonomies (config -> 60+ sub-packages, cli -> cmd/root
+ core/, output -> write/*, errors -> err/*), adding an MCP server,
4 session parsers, and 15+ new commands. The architecture is
clean-layered with strict dependency direction. Two direct Go
dependencies (cobra, yaml.v3). All user-facing text externalized
to YAML for i18n readiness.

The project is in a "structural maturity" phase — 40+ audit tests
enforce conventions (no magic strings, no raw file I/O, no flag
binding outside flagbind, etc.), but the feature surface is still
expanding (memory bridge, MCP, multi-agent support).

## Vision Alignment

<!-- TODO: create a vision document to make this explicit, make it detailed -->

[inferred] Based on the codebase trajectory (MCP server, 4 parsers,
setup command supporting 5 tools, memory bridge, peer MCP model
decision), the vision is: **ctx as the universal context layer for
any AI coding assistant**.

The current architecture **partially** supports this:

**What aligns well**:
- MCP server is 100% agent-agnostic — pure JSON-RPC 2.0 protocol
- Session parser is extensible (4 formats: Claude, Copilot, Copilot
  CLI, Markdown)
- Context files are plain markdown — any tool can read them
- Token budgeting is agent-independent
- `ctx setup` generates configs for 5 tools

<!-- todo: these needs to be specced out in depth -->

**What constrains the vision**:
- **Hook system is Claude Code-only**: hooks.json, settings.local.json,
  CLAUDE.md are deeply embedded in init, assets, and 34 system
  subcommands. No equivalent lifecycle for Cursor/Copilot/Aider.
- **Skills are Claude Code-specific**: 30 live skills in `.claude/skills/`
  with SKILL.md format; MCP prompts cover only 5 basic operations.
- **Session import is Claude-heavy**: The Claude Code JSONL parser is
  the most mature; Copilot parsers are newer and thinner.
- **Memory bridge assumes Claude Code**: DiscoverPath() uses Claude's
  project slug format; MergePublished() targets MEMORY.md.
- **Bootstrap injects Claude Code config**: init deploys .claude/
  directory, hooks, settings, plugin version.
- **AGENT_PLAYBOOK.md references Claude Code**: behavioral instructions
  are Claude-specific.

## Future Direction

<!-- TODO: document this make it explicit; create specs and tasks -->

[inferred] The architecture should evolve toward a **multi-agent
runtime abstraction** where Claude Code is one integration among
many, not the assumed default.

Key structural changes needed:

1. **Agent Integration Interface**: Abstract the hook lifecycle
   (UserPromptSubmit / PreToolUse / PostToolUse) behind an
   interface. Each agent gets an adapter implementing its native
   lifecycle hooks. Claude Code's hooks.json becomes one adapter.

2. **Skill -> Prompt Parity**: MCP prompts currently cover 5 basic
   operations vs 30 Claude skills. Either expand MCP prompts or
   create a tool-agnostic skill format that compiles to each tool's
   native format.

3. **Multi-Agent Memory Discovery**: memory.DiscoverPath() should
   accept a registry of slug conventions per tool, not hardcode
   Claude's format.

4. **Agent-Neutral Playbook**: AGENT_PLAYBOOK.md should be a
   generic agent playbook; Claude-specific instructions should live
   in the Claude integration layer.

## Known Bottlenecks

2. **write/ package proliferation**: 46 packages that are mostly
   3-10 line functions. The overhead of a package per command for
   output is high relative to the content. Many write packages have
   only 1-2 functions. <!-- will solve itself after go templatification; then better to do a new arch analysis -->

3. **Single-threaded MCP server**: All tool calls processed
   sequentially. Read-only tools (status, drift, next, remind)
   could safely run concurrently.

5. **Three-file ceremony for new text**: Adding a user-facing
   message requires a DescKey constant, a YAML entry, and a
   write/err function. This is the cost of externalization — paid
   back by agent legibility and mechanical linkage verification,
   but real friction for rapid iteration. <!-- given an agent is rapidly iterating rather than a human; it's a non-issue; and there is nothing a good ide cannot solve -->

## Implementation Alternatives

### 1. Config Organization: Explosion vs Domains

**Current**: 60+ atomic sub-packages (one per concept).
**Alternative**: 5-8 domain packages (config/mcp, config/cli,
config/entry, config/infra, config/text).
**Tradeoffs**: Current gives surgical imports but noisy navigation.
Alternative reduces package count but creates larger dependency
units. Current approach is correct for compilation; alternative is
correct for developer experience.
**Recommendation**: Keep current structure but add a config/README
documenting the organization principle and import patterns.

### 2. Output Layer: Package-per-Command vs Functional

**Current**: 46 write/ packages (one per command).
**Alternative**: 5-8 write packages by output type (write/table,
write/status, write/entry, write/json, write/progress).
**Tradeoffs**: Current aligns with CLI structure but many packages
have only 1-2 functions. Alternative reduces package count but
breaks the 1:1 CLI correspondence.
**Recommendation**: Keep current structure — it scales cleanly and
matches the CLI taxonomy. The overhead is tooling noise, not
architectural cost.

### 3. MCP Threading: Sequential vs Concurrent

**Current**: Single-threaded main loop.
**Alternative**: Goroutine pool with read-only tools executed
concurrently.
**Tradeoffs**: Current is simple and correct. Concurrent execution
requires classifying tools as read-only vs mutating and managing
shared state (context files, session state).
**Recommendation**: Add request timeout first (easy win). Consider
concurrent execution only if latency becomes a user complaint.

## Gaps

1. **No agent abstraction layer**: The codebase has "Claude Code
   integration" baked into init, assets, hooks, and memory — but
   no equivalent for other agents. Adding Cursor support means
   touching 10+ packages. There should be an `internal/agent/`
   package defining the integration interface.

2. **No lifecycle events for non-Claude agents**: The 34 system
   hook subcommands are Claude Code's lifecycle. Cursor (via
   .cursorrules), Copilot (via copilot-instructions.md), and Aider
   (via .aider.conf) have no equivalent event-driven integration.

3. **No tool capability negotiation**: ctx exposes the same 11 MCP
   tools to all clients. Different agents have different
   capabilities — some can write files directly, others need
   explicit confirmation. No capability-based tool filtering.

4. **No prompt template versioning**: MCP prompts are hardcoded.
   No way to version, A/B test, or customize prompt content per
   agent or per user.

5. **No context diff/delta protocol**: Every agent interaction
   loads the full context. No mechanism to send only what changed
   since last load, which matters for long-running MCP sessions.

6. **No multi-project support**: ctx assumes one .context/ per
   project root. Monorepo users with multiple logical projects
   in one repo have no subdivision mechanism.

## Areas of Improvement

**High impact, low effort** (do first):
- Add request timeout to MCP server handler calls (prevents hangs)
- Create `internal/agent/` interface package with Agent type
  (Name, HookFormat, ConfigPath, DiscoverMemory methods)
- Add `ctx setup --list` to show supported agents and their
  integration status
- Document the config/ organization principle in a README

**High impact, high effort** (plan for):
- Abstract hook lifecycle behind agent adapter interface; refactor
  34 system subcommands to use adapter dispatch
- Expand MCP prompts to cover skill-equivalent operations (from 5
  to 15-20 prompts matching the most-used skills)
- Implement context delta protocol for MCP sessions (track last
  resource read per client, send diffs)
- Create agent-specific playbook generator (generic core + agent
  overlay)

**Low impact** (defer or skip):
- Consolidate small write/ packages (cosmetic, not architectural)
- Add concurrent tool execution in MCP (premature optimization)
- Merge config/ sub-packages into domains (trading one cost for
  another)

## Risks

1. **Claude Code coupling calcification**: Every feature built
   against Claude Code's specific lifecycle (hooks.json, JSONL
   sessions, MEMORY.md, settings.local.json) raises the cost of
   supporting other agents. The 34 hook subcommands are accumulating
   Claude-specific assumptions that will be expensive to generalize.

2. **MCP server as afterthought**: The MCP server exposes 11 tools
   and 5 prompts while the CLI has 34 commands and 30 skills. This
   gap means MCP clients get a degraded experience. If MCP becomes
   the primary integration path (as the market is trending), the
   CLI-first design creates a capability gap.

3. **Audit test brittleness**: 40+ AST-based audit tests enforce
   conventions, but they also encode assumptions about package
   structure. Major refactoring requires updating quarantine lists,
   grandfathered maps, and allowlists across multiple audit files.

4. **YAML text externalization overhead**: Developers must maintain
   879+ YAML entries alongside code. Adding a new error message
   requires changes in 3 places (YAML file, DescKey constant, error
   constructor). This friction slows development and creates sync
   drift opportunities.

5. **Memory bridge is one-way strategic bet**: The memory bridge
   assumes Claude Code's auto-memory is the primary external memory
   source. If Cursor, Copilot, or other tools implement their own
   memory systems with different formats, ctx needs N memory
   adapters.

## Intervention Points

Top 5 highest-leverage places to implement agent-agnostic features:

1. **`internal/bootstrap/group.go`** - Command registration is
   centralized here. Adding agent-aware command variants (e.g.,
   `ctx setup cursor` with full lifecycle) is a single-file change
   that immediately surfaces in help output. Impact: all users.

2. **`internal/cli/setup/`** - The setup command already supports 5
   agents. This is the natural place to add agent-specific lifecycle
   deployment (hooks, config files, playbooks). Extending this
   creates parity without touching other packages.

3. **`internal/mcp/server/def/prompt/`** - Prompt definitions are
   the MCP equivalent of skills. Expanding from 5 to 20 prompts
   here immediately enriches the MCP experience for all agents.
   Each prompt is self-contained (no cross-prompt dependencies).

4. **`internal/journal/parser/`** - Parser registration is
   interface-based. Adding new session formats (Cursor, Aider,
   Windsurf) here enables journal import for those agents without
   touching any other package.

5. **`internal/memory/`** - Memory discovery is the only package
   that reads external agent state. Abstracting DiscoverPath()
   into an agent-keyed registry makes memory bridge work for any
   agent with persistent memory.

<-- TODO: create specs for these -->

## Upstream Proposals

### 1. Agent Integration Protocol Spec

**What**: Define a formal spec for how AI coding assistants
integrate with project context systems. Covers: lifecycle hooks,
context loading, session export format, memory sync format.

**Why**: Every AI tool is inventing its own integration format
(hooks.json, .cursorrules, copilot-instructions.md, .aider.conf).
A shared spec would let ctx (and similar tools) support all agents
through one interface.

**Where**: Crosses the boundary between ctx's integration layer
and upstream agent implementations. Affects bootstrap, setup,
system, memory packages.

**Risk**: Adoption requires buy-in from multiple tool vendors. May
be premature — the market is still evolving. Start with a de-facto
spec based on ctx's existing multi-agent support.

### 2. MCP Context Resource Subscription Standard

**What**: Propose a standard pattern for MCP servers that provide
project context resources with token-budgeted assembly, change
notifications, and delta delivery.

**Why**: ctx's MCP server solves a general problem (expose project
knowledge to AI agents) that every context-management tool will
face. Standardizing the resource URI scheme, subscription model,
and assembly protocol would enable interoperability.

**Where**: Crosses the MCP protocol boundary. Affects mcp/server/,
mcp/proto/, mcp/handler/ and the MCP specification itself.

**Risk**: MCP spec is evolving rapidly. A proposal now may be
obsoleted by upstream changes. Mitigate by implementing first,
proposing from experience.

### 3. Session Export Format Standard

**What**: Define a common session export format that any AI tool can
produce and any analysis tool can consume. JSON-based, covering:
messages, tool calls, tool results, timing, token counts, metadata.

**Why**: ctx currently reverse-engineers Claude Code's JSONL format
and parses Copilot's format separately. Every AI tool has its own
session format. A standard would eliminate per-tool parsers and
enable ecosystem tooling (analytics, audit, replay).

**Where**: Crosses journal/parser boundary and upstream session
storage in Claude Code, Copilot, Cursor, etc.

**Risk**: Tool vendors may resist standardization to maintain lock-in.
Start by publishing ctx's internal session representation (entity.
Session) as a proposed schema and building import/export around it.

## Productization Gaps

**Multi-tenant / multi-project**:
- ctx assumes one .context/ per repo root. Monorepo users with 5
  services in one repo can't maintain separate context per service.
- No team/shared context model — each developer maintains their
  own .context/ independently. No merge, no conflict resolution.

**Observability**:
- Event log (log/event) is local JSONL with size-based rotation.
  No aggregation, no dashboard, no alerting.
- MCP session state is in-memory only — no telemetry on tool usage
  patterns, governance warning effectiveness, or agent behavior.
- No way to measure whether context is actually helping agents
  make better decisions (no outcome tracking).

**Operational hardening**:
- File locking is absent — concurrent CLI and MCP writes to the
  same context file can corrupt data.
- No backup/restore for .context/ beyond manual git (the backup
  command exists but targets SMB shares).
- No graceful degradation — missing .context/ files cause hard
  errors rather than operating in reduced mode.

**What a large customer would hit first**:
- "How do I deploy ctx to 50 engineers?" — no provisioning,
  no central config, no onboarding automation beyond ctx init.
- "How do I enforce context quality?" — drift checks are local;
  no CI integration beyond manual Makefile targets.
- "How do I share context between team members?" — no shared
  context, no merge strategy, no conflict resolution.

## Failure-First Analysis

**Hidden assumptions**:
- Assumes single-user, single-agent access to .context/ at any
  time. No file locking, no conflict detection. Two terminals
  running ctx add simultaneously can lose writes.
- Assumes the curated tier (.context/*.md) stays small. load.Do()
  reads only top-level .md files, not journal/ — but the curated
  files themselves have no size guard beyond token budgeting.
- Assumes stable upstream formats. Claude Code JSONL format,
  Copilot session format, MEMORY.md slug convention — all are
  undocumented and can change without notice.
- Assumes linear session lifecycle. The hook system expects
  UserPromptSubmit -> tool calls -> PostToolUse in sequence.
  Parallel tool execution (which Claude Code supports) can
  cause interleaved hook events.

**What breaks silently**:
- Missing YAML text keys produce empty strings, not errors.
  A user sees missing output but no crash.
- Failed state writes (journal state, event log) are logged to
  stderr and swallowed. Data loss happens quietly.
- Governance warnings in MCP are appended to response text —
  if the agent doesn't surface the full response, warnings vanish.

**What breaks loudly**:
- Missing .context/ directory -> immediate error on any command
- Corrupted .ctxrc YAML -> parse warning + defaults used
- Scanner buffer overflow in MCP -> parse error response

**Cascade risk**:
- If config/embed/text YAML is malformed, every command that calls
  desc.Text() gets empty strings. The entire CLI becomes a wall of
  blank output. This is the single highest-impact failure point.

## Silent Choices

The codebase is making bets that are not documented as decisions
but shape the architecture significantly:

1. **File-per-entity package structure will scale indefinitely.**
   The explosion to 527 packages (60+ config, 46 write, 35 err)
   bet that Go tooling handles this well. It does today. At 1000+
   packages, build times, IDE performance, and import management
   become real costs. The abstraction boundary is the package, and
   every new entity creates a new package.

2. **YAML text externalization pays for itself before i18n.**
   879+ description keys with 3-file-change overhead per new message.
   The justification is NOT future i18n — it's agent legibility
   (named constants are traversable graphs, not guessable literals),
   drift prevention (TestDescKeyYAMLLinkage catches orphans), and
   the archaeology itself (finding 879 scattered strings was the
   hard part; YAML is the artifact of that work). i18n is a free
   downstream consequence, not the bet.

3. **MCP is a secondary interface, CLI is primary.**
   The CLI has 34 commands and 30 skills; MCP has 11 tools and 5
   prompts. This bets that human-in-the-loop CLI usage remains
   dominant. If the market moves to autonomous agents that prefer
   MCP, ctx's CLI-first design becomes a liability.

4. **Claude Code is the reference agent.**
   34 hook subcommands, Claude-specific init, MEMORY.md bridge,
   JSONL parser maturity — all bet that Claude Code remains the
   primary user. If Cursor or Copilot overtake Claude Code, the
   migration cost to equal-citizen support is high.

5. **Single-project, single-developer scope is sufficient.**
   No multi-project, no team sharing, no central server. This
   bets that AI coding assistants remain individual-developer
   tools. If team-based AI workflows emerge (shared context,
   coordinated agents), ctx has no foundation for it.

## Onboarding Friction

**What a new engineer hits in week one:**

1. **Package discovery**: 527 packages with no search or guide
   beyond ARCHITECTURE.md. Finding where a feature lives requires
   grep or asking. The config/ explosion (60 packages) is
   especially disorienting — "where do I put a new constant?"
   has a non-obvious answer.

2. **Three-file output ceremony**: Adding a new user-facing message
   requires (a) YAML entry in assets, (b) DescKey constant in
   config/embed/text, (c) error or write function using the key.
   This is not documented in onboarding materials.

3. **Audit test quarantine**: Making changes often triggers audit
   test failures (magic strings, raw file I/O, flag binding). The
   fix is to use the "right" abstraction, but which abstraction
   is right isn't obvious without reading the audit tests.

4. **Hook system opacity**: 34 hidden system subcommands that fire
   on agent lifecycle events. Understanding what happens when
   "Claude starts a session" requires reading hooks.json +
   12 check_ commands + input.go stdin protocol.

5. **What isn't written down**: The relationship between
   .claude/skills/ (agent instructions) and MCP prompts
   (protocol-level instructions) is implicit. Why some operations
   are skills and others are prompts, and when to use which, is
   tribal knowledge.

## Domain Clustering Comparison (enriched 2026-04-03 via GitNexus)

GitNexus auto-detected 94 functional clusters from the call graph.
Comparing against the 5 manual DETAILED_DESIGN domain splits:

| Manual Domain | GitNexus Clusters | Match | Notes |
|---------------|-------------------|-------|-------|
| Foundation | Format (97%), Sysinfo (82%), Rc (57%), Io (53%) | Partial | Format is extremely cohesive; Rc and Io are looser than expected |
| Domain | Drift (263 sym, 84%), Session (110, 88%), Memory (91, 61%), Parser (89, 70%), Journal (74, 47%), Trace (59, 68%) | Partial | Drift is the largest cluster — pulls in cli/drift + write/drift. Journal is low cohesion (47%) |
| MCP | Server (49, 90%) | Yes | Tight, well-bounded cluster |
| CLI | Root (85, 83%), Initialize (72, 50%), Bootstrap (41, 67%) | Partial | Initialize is low cohesion — it touches too many domains during deployment |
| Output | Lock (42, 76%), Notify (46, 65%), Watch (39, 57%) | No | GitNexus doesn't group write/* as a cluster; it groups by behavioral coupling instead |

### Hidden Coupling (GitNexus groups, manual splits)

- **Drift cluster (263 symbols)** includes symbols from drift/,
  cli/drift/, write/drift/, and context/load — all tightly coupled
  through the Detect() -> load.Do() -> report chain. Manual split
  puts these in 3 different domains.
- **Pad cluster (200 symbols)** is self-contained — pad, crypto,
  and store form a tight unit. Manual split disperses these across
  Foundation (crypto) and Domain (pad).

### Artificial Grouping (GitNexus splits, manual groups)

- **Journal (47% cohesion)** — lowest cohesion cluster. The site
  pipeline (zensical, CSS, nav links) and the import pipeline
  (parser, state, execute) have minimal internal coupling. These
  are effectively two separate subsystems sharing a command name.
- **Initialize (50% cohesion)** — touches every domain during
  deployment (templates, hooks, skills, settings, vscode, makefile,
  gitignore). Low cohesion is structural, not a bug — it's a
  deployment orchestrator, not a domain.

### Key Insight

The write/* and err/* layers don't form clusters because they have
no internal coupling — each write package is called by exactly one
CLI command. GitNexus correctly identifies them as leaves, not
communities. The manual "Output" domain is an organizational
convenience, not an architectural boundary.

## Questions That Would Sharpen This Analysis

Answering any of these would move speculative sections to grounded:

1. **Vision** - Is "universal context layer for any AI agent" the
   actual 12-24 month goal, or is Claude Code primacy intentional?
2. **Future direction** - Is MCP intended to reach parity with CLI,
   or stay as a lightweight complement?
3. **Known bottlenecks** - Is the 3-file text ceremony causing
   measurable dev velocity drag?
4. **Assumptions marked** - These sections are labeled [inferred]:
   Vision Alignment, Future Direction
