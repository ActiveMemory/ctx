# context-model-and-state

## [2026-06-07-162840] Pin an on-disk contract before splitting work across parallel agents

**Context**: The ctx-dream skill (proposal writer) and the ctx dream CLI
(proposal reader) were built by parallel tracks; each independently invented the
proposals layout — one per-file proposals/<id>.json, the other a single
proposals.json array — requiring reconciliation at integration.

**Lesson**: When parallel agents share a serialized artifact (a file path +
schema), the exact on-disk shape must be fixed in BOTH prompts up front. 'Each
track invents it' guarantees a mismatch that only surfaces at integration.

**Application**: Before fanning out producer/consumer work to parallel agents,
write the precise file path + JSON schema into both agents' specs as a fixed
contract, and reconcile/verify it first thing at integration.

---

## [2026-06-07-170016] State, tombstones, logs & filesystem hygiene (consolidated)

**Consolidated from**: 6 entries (2026-02-11 to 2026-03-06)

- Permission drift is distinct from code drift — settings.local.json is
  gitignored so no review catches stale entries; it accumulates session debris
  (run /sanitize-permissions + /ctx-drift). Skill() permissions don't support
  name-prefix globs (list each); wildcard trusted binaries (Bash(ctx:*),
  Bash(make:*)) but keep git granular (never Bash(git:*)).
- Gitignored directories are invisible to git status — stale artifacts persist
  indefinitely (periodically ls them). Add editor artifacts (*.swp,*.swo,*~) to
  .gitignore from day one. Gitignore entries for sensitive paths are security
  controls, not documentation — never remove during cleanup.
- The state directory accumulates write-only session tombstones and grows
  unbounded without auto-prune (234 files found); autoPrune(7) now runs once per
  session at startup via context-load-gate (manual `ctx system prune` still
  available).
- A session-scoped tombstone must include the session ID in its filename, else
  it suppresses hooks across ALL concurrent and future sessions (memory-drift
  fixed; backup-reminded, ceremony-reminded, check-knowledge, journal-reminded,
  version-checked, ctx-wrapped-up still carry this bug). Use the UUID pattern so
  prune can clean them.
- New log sinks must follow the established rotation pattern (size-based, single
  previous generation): eventlog rotated at 1MB but logMessage() in state.go was
  append-only with no size check.
- If a directory is recreated (auto-prune), an SSH shell holding the old inode
  won't see new files (ls returns "no such file" though cat with the full path
  works elsewhere); after `ctx system prune` or any state recreation, SSH
  sessions need cd-. or re-login.

---

## [2026-06-01-174927] Guard managed blocks before regenerating; don't trust the span to be machine-owned

**Context**: ctx learning add silently deleted entry bodies that lived between
INDEX:START/END markers: index.Update replaced the whole marker span with a
regenerated table, and ParseHeaders scanning the full file made the result look
complete, hiding the loss.

**Lesson**: Code that 'replaces the managed block' (index regen, KB managed
blocks, moc.go) assumes the span between its markers is disposable and
machine-owned. That assumption breaks the moment user content drifts inside the
markers, and the regenerated output looks correct so the loss is invisible. The
fix is a precondition guard that refuses to mutate when regeneration would lose
data — not smarter parsing of the trapped content.

**Application**: Before any 'replace between markers' write, validate the span:
refuse on entry/content found where only generated output belongs, and on
malformed/duplicated/out-of-order markers. Fail loud and leave the file
byte-identical rather than regenerate. Run the guard at the read-before-mutate
choke point so nothing is written on refusal.

---

## [2026-05-20-214830] Handover filenames are archaeology; parse by generated-at, not filename

**Context**: User observed three coexisting handover filename shapes:
.context/HANDOVER-2026-04-22.md (pre-skill root file),
.context/handovers/YYYY-MM-DD-HHMMSS-slug.md (skill-era pre-CLI),
.context/handovers/<RFC3339Compact>-slug.md (current CLI). User asked whether
this was a regression or a skill-interpretation problem.

**Lesson**: Neither. The .context/HANDOVER-* root file predates the handovers/
directory contract entirely (the body even said 'delete this file after
reading'). The YYYY-MM-DD-HHMMSS shape was an earlier skill iteration writing
free-form before ctx handover write existed (commit 60543e46, 2026-05-17,
introduced the CLI as sole writer per the anti-pattern note in /ctx-handover
SKILL.md). The current parser at internal/write/handover/parse.go:75-107 keys on
the 'generated-at' YAML frontmatter, not the filename — so legacy shapes still
sort correctly via LatestHandoverCursor. Only files without frontmatter (the
root April file) are invisible.

**Application**: When unifying filename shapes across history, use git mv to
preserve rename detection. Derive the canonical timestamp from the file's own
generated-at frontmatter rather than from the filename — that's the source of
truth the parser uses anyway. If a handover predates frontmatter entirely (rare,
pre-skill era), it's safe to delete because the parser never read it.

---

## [2026-04-14-010134] Raft-lite trade-off is the load-bearing choice in internal/hub

**Context**: Discovered while writing thorough doc.go for internal/hub. The
package embeds HashiCorp Raft for leader election only; data replication is
sequence-based gRPC sync over the append-only JSONL store.

**Lesson**: A leader crash window between accept and replicate can lose the most
recent write. Append-only storage plus idempotent clients make this acceptable;
full Raft log replication would not be needed and would not be simpler.

**Application**: Any future make hub stronger proposal must engage with this
trade-off explicitly. Do not abandon Raft-lite accidentally by introducing
log-replicated state; that would invalidate the simplicity argument.

---

## [2026-04-13-153618] rc.ContextDir() is the single source of truth — fix the resolver, not callers

**Context**: When ctx init failed with a boundary error, my first instinct was
to have init bypass rc.ContextDir() and use filepath.Join(cwd, dir.Context)
directly. Volkan shut that down: rc.ContextDir() encodes invariants (team
shares, symlinks, network mounts, .ctxrc overrides) that individual commands
cannot reason about.

**Lesson**: Resolution chains with multiple fallbacks are contracts. If one
command bypasses the chain, it silently diverges from every other command's
notion of 'the context directory.' When a resolver produces a wrong answer for a
specific case, fix the resolver — don't let callers opt out.

**Application**: Any time you see rc.ContextDir(), rc.RC(), or similar central
resolvers producing a bad result, the fix belongs in the resolver itself (or in
its input data like .ctxrc). Caller-side bypasses create drift.

---

## [2026-04-09-001323] Pad index shifting is a real UX bug in batch operations

**Context**: ctx pad rm 10; rm 11; rm 12 deleted wrong entries because indices
shifted after each deletion

**Lesson**: Any ID-based system where users chain operations needs stable IDs.
Look-then-act is safe for single ops; look-then-batch-act breaks with shifting
indices

**Application**: Both pad and remind now use stable IDs with batch delete and
range support. Apply same pattern to any future numbered-list subsystem

---

## [2026-02-24-032945] CLI tools don't benefit from in-memory caching of context files

**Context**: Discussed whether ctx should read and cache LEARNINGS.md,
DECISIONS.md etc. in memory

**Lesson**: ctx is a short-lived CLI process, not a daemon. Context files are
tiny (few KB), sub-millisecond to read. Cache invalidation complexity exceeds
the read cost. Caching only makes sense if ctx becomes a long-lived process (MCP
server, watch daemon).

**Application**: Don't add caching layers to ctx's file reads. If an MCP server
mode is ever added, revisit then.

---

