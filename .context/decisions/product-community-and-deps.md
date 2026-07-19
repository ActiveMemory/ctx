# product-community-and-deps

## [2026-07-16-215955] Progressive disclosure: canonical knowledge files become bounded roots with themed rollout

**Status**: Accepted

**Context**: Canonical knowledge files (LEARNINGS/DECISIONS/CONVENTIONS) grow without bound as entries accumulate, and the entries stay valid, so nothing can be dropped: time-sharding plus a load-excluded 'cold' bucket was already rejected (a supersession pass found only ~1.5% cold across 162 entries; recency-gating is dangerous because old approximately equals live). At sufficient scale an agent that legitimately wants system understanding reads every decision, then every learning, and exhausts its context window. Consumption discipline (headings-first via ctx index) is necessary but NOT sufficient: an agent can always choose to read the whole file, and will when it wants completeness. Consolidation does not help either: a 2026-07-16 pass moved LEARNINGS only 98 to 88 because the corpus is dense with distinct signal, not redundancy.

**Decision**: Progressive disclosure: canonical knowledge files become bounded roots with themed rollout

**Rationale**: Structural boundedness makes 'read it all and die' impossible rather than merely discouraged, applying the project's own 'mechanical gates over prose' principle to context loading. Storing the generated synthesis is justified precisely BECAUSE it is expensive to recompute (an LLM pass), unlike the trivially-recomputable INDEX block whose storage was pure waste and pure drift: the cost/benefit of storing derived state inverts when recomputation is costly. Reading the bounded root alone yields compressed history PLUS verbatim recent delta, a complete current picture with no staleness gap, because staging IS the un-digested remainder by construction. The staging zone therefore serves as the watermark, so no state file is needed: a .context/state/digest.json would be a second source of truth about a fact the file already states (the INDEX-block mistake in a new coat), and its one advantage (detecting a misplaced entry) is recovered by a cheap structural invariant instead. Verified against the code: the add path anchors on the first line-start '## [' and falls back to AfterHeader, so entries always land above '## Themes' even when staging is empty, meaning ctx decision/learning add needs ZERO change; and because the root is bounded, the existing ctx agent packet needs no rewire either.

**Consequence**: Entry bodies now move between files, which is the clobber risk class that index.Validate exists for, so the pass must append to the theme file, verify byte-presence, and only then remove from staging; it validates preconditions, fails loud with no auto-repair (matching the learning-add clobber fix precedent), and prefers duplication over loss on crash. New structural invariants become mechanically checkable: no line-start '## [' below '## Themes'; root gists and theme files are 1:1; every entry lives in exactly one place. CONVENTIONS needs an extra trailing '## Recent' staging heading because AppendAtEnd plus '###' prose sections would nest ambiguously inside '## Themes', and conventions edits-in-place now happen in theme files behind a link. Theme proliferation remains a slow unbounded growth vector on the root; the structure is self-similar, so an overgrown theme file can itself become a root (nesting deferred, not precluded). Scope is LEARNINGS, DECISIONS, CONVENTIONS; CONSTITUTION (small by design) and TASKS (auto-archived) are excluded. The pass is agent-driven and human-gated (agent suggests themes, human can override), triggered by suggestion only from the growth nudge, /ctx-remember, and /ctx-wrap-up, never performed inline at wrap-up because the human is leaving at that point. Spec: specs/progressive-disclosure.md.

---

## [2026-07-15-000000] Live ceremony credit reuses the daily throttle marker, suppressing the day's other ceremony nudge

**Status**: Accepted

**Context**: The `check-ceremony` hook nudges the operator to open
sessions with `/ctx-remember` and close them with `/ctx-wrap-up`. It was
journal-driven only, so it self-nudged on the very prompt that ran a
ceremony and kept nudging until the session was journal-imported. The
fix (specs/ceremony-nudge-live-session.md) parses the live prompt and,
when it is a ceremony command, credits it immediately. The question was
*how* to record that live credit.

**Alternatives Considered**:
- Reuse the existing per-day throttle marker (`ceremony-reminded`):
  touching it on a live ceremony credits the session and settles the
  ceremony question for the day. Pros: zero new state; one guard;
  matches the check's existing once-per-day cadence. Cons: crediting a
  live `/ctx-remember` also suppresses a would-be `/ctx-wrap-up` nudge
  for the rest of that day (and vice-versa).
- Per-ceremony live markers (separate remember/wrap-up credit): Pros:
  the other ceremony can still nudge the same day. Cons: new state files,
  a second throttle axis, and more moving parts for a coarse daily nudge.

**Decision**: Reuse the single daily throttle marker. On a live ceremony
prompt, `check-ceremony` touches `ceremony-reminded` and returns without
nudging.

**Rationale**: The check is a deliberately coarse daily cadence, not a
per-ceremony ledger. An operator actively running one ceremony does not
need to be nagged about the other on that same prompt, and the marker's
"settled for today" semantics already express exactly that. The extra
state a per-ceremony scheme buys is not worth it for a once-a-day hint.

**Consequence**: Running either ceremony live suppresses both ceremony
nudges for the rest of that day. This is intended; the trade-off is
documented in the spec's Trade-off section. If a per-ceremony cadence is
ever wanted, this is the decision to revisit.

**Related**: See spec specs/ceremony-nudge-live-session.md.

## [2026-07-04-152957] Statusline informs, never gates

**Status**: Accepted

**Context**: Porting a cost-aware status line whose reference design escalates to a spend alarm with a model-switch suggestion when an expensive model family is detected

**Decision**: Statusline informs, never gates

**Rationale**: A family-substring rule carries no task context: the expensive model is often the cheaper choice per outcome, and a statusline cannot see outcomes, so it must not prescribe. Alarms also only reach operators who are present and already see the dollar figure; unattended overspend belongs to loop/cron notifications

**Consequence**: ctx system statusline renders model, ctx%, and plain cost only; show_cost exists for screen-sharing; any future budget enforcement must be a separate, deliberate feature, not statusline creep

---

## [2026-07-03-182236] Keep sonnet-4-6 at the 200k default despite the API catalog listing 1M

**Status**: Accepted

**Context**: Claude 5 window-detection fix (specs/model-context-window-fable.md); Anthropic's current catalog shows Sonnet 4.6 with a 1M window at the API level

**Decision**: Keep sonnet-4-6 at the 200k default despite the API catalog listing 1M

**Rationale**: ctx models Claude Code's per-session gating, not raw API capability: Sonnet 4.6's 1M is an explicit opt-in already handled by ModelSuffix1M and ClaudeSettingsHas1M, and the existing test deliberately pins the 200k default

**Consequence**: Revisit only with session-level evidence (a sonnet-4-6 JSONL showing 1M without the [1m] suffix); over-reporting the window would silence the context hook while a session genuinely fills

---

## [2026-04-01-233247] IRC to Discord as primary community channel

**Status**: Accepted

**Context**: Discord server exists at https://ctx.ist/discord; IRC/libera.chat
references were stale

**Decision**: IRC to Discord as primary community channel

**Rationale**: Discord is faster for async community support; IRC was historical

**Consequence**: Updated zensical.toml, README, community docs, journal
template. Added community footer to ctx help and ctx init output via YAML assets
pipeline

---

## [2026-03-06-200306] Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant

**Status**: Accepted

**Context**: fatih/color was used in 32 files for green checkmarks, yellow
warnings, cyan headings, dim text

**Decision**: Drop fatih/color dependency — Unicode symbols are sufficient for
terminal output, color was redundant

**Rationale**: Every colored output already had a semantic symbol (✓, ⚠,
○) that conveyed the same meaning; color added visual noise in non-terminal
contexts (logs, pipes)

**Consequence**: Removed --no-color flag (only existed for color.NoColor); one
fewer external dependency; FlagNoColor retained in config for CLI compatibility

---




---

