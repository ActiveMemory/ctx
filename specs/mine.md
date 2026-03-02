# ctx mine

## Problem

Projects accumulate knowledge in git history that never gets captured
in context files. Every commit message encodes a decision, a lesson,
or a direction — but this knowledge is locked in `git log` where no
AI session will ever see it.

New projects get the worst of it: hundreds of commits of institutional
knowledge sitting unextracted while DECISIONS.md and LEARNINGS.md
start empty.

Inspired by the "archeology" pattern of mining PR review history to
build review skill files, `ctx mine` applies the same Extract →
Distill → Codify loop to git history — populating context files from
what the team already wrote.

## Approach

Skill-first design. `/ctx-mine` is a user-invocable skill that:

1. Reads a ledger to skip already-processed commits
2. Gathers unmined commits from `git log`
3. Processes them in token-friendly chunks
4. Categorizes findings into five buckets
5. Presents findings for human review
6. Persists approved findings via `ctx add`
7. Updates the ledger

No CLI command for v1 — the skill manages everything using `git`,
`ctx add`, and direct ledger file I/O. A `ctx mine ledger` CLI
subcommand can be added later if programmatic access is needed.

### Why skill-only?

The intelligence is in the distillation, not the data gathering.
`git log` is the data source and it's already a CLI tool. The skill's
value is in the prompting strategy that turns raw commit metadata into
structured context entries. Writing Go code to wrap `git log` adds
complexity without adding capability.

## Behavior

### Five Finding Buckets

| Bucket | Target file | Description |
|--------|-------------|-------------|
| **Decisions** | DECISIONS.md | Architectural choices: "switched from X to Y", "chose library Z" |
| **Learnings** | LEARNINGS.md | Gotchas, bugs, patterns: "nil pointer when X", "must call Y before Z" |
| **Todos** | TASKS.md | Follow-up work: "TODO in commit msg", incomplete migrations, tech debt |
| **Spec candidates** | Reported only | Feature ideas or large efforts worth a design doc |
| **Miscellaneous** | Reported only | Interesting observations that don't fit a bucket |

Decisions, learnings, and todos can be persisted directly via `ctx add`.
Spec candidates and miscellaneous are presented in the summary for the
user to act on manually.

### Happy Path

```
User: /ctx-mine

Agent:
1. Reads .context/state/mine-ledger.json
2. Runs: git log --format='%H|%ai|%s' <since-last-mined>
3. Identifies 47 unmined commits
4. Chunks into groups of ~20 commits
5. For each chunk, reads commit messages + file lists (git show --stat)
6. Analyzes patterns, repeated themes, architectural shifts
7. Produces categorized findings

Presents:
  ## Findings (12 items from 47 commits)

  ### Decisions (3)
  1. Adopted AES-256-GCM for scratchpad encryption
     Evidence: commits abc123, def456 — "add encrypted pad", "switch to GCM"
  2. Moved from shell hooks to Go-native hook runner
     Evidence: commits ghi789, jkl012

  ### Learnings (4)
  1. Hook output must use stderr — stdout gets swallowed by Claude Code
     Evidence: commit mno345 — "fix: route hook output to stderr"
  2. ...

  ### Todos (2)
  1. Split monolithic test file (500+ lines) in cli/pad
     Evidence: commit pqr678 — "TODO: split pad_test.go"

  ### Spec Candidates (1)
  1. Smart context retrieval — multiple commits reference token overflow
     Evidence: commits stu901, vwx234

  ### Miscellaneous (2)
  1. Consistent use of filepath.Join over string concat (convention?)
     Evidence: 6 commits with path-related fixes

  ---
  Which findings should I persist?
  Type "all", specific numbers (e.g., "D1, D3, L2, T1"), or "none".

User: D1, D2, L1, T1

Agent:
  Persisted:
  - ctx add decision "Adopted AES-256-GCM for scratchpad encryption"
  - ctx add decision "Moved from shell hooks to Go-native hook runner"
  - ctx add learning "Hook output must use stderr, not stdout"
  - ctx add task "Split monolithic test file in cli/pad"
  Updated ledger: 47 commits marked as mined.
```

### Two-Pass Processing

To keep token costs manageable:

**Pass 1 — Scan (cheap):** Process commit messages and file lists
only. No diffs. This is enough to identify most decisions and
learnings. A commit message like "fix: handle nil pointer in session
parser" already encodes the learning.

**Pass 2 — Drill (optional, on demand):** For ambiguous findings
or when the user asks "tell me more about finding L3", fetch the
actual diff for those specific commits. This avoids loading hundreds
of diffs upfront.

### Chunking Strategy

Commits are processed in chunks sized to fit comfortably in a
sub-agent's context window:

- ~20 commits per chunk (adjustable via `--chunk-size`)
- Each chunk includes: hash, date, author, message, file list
- Estimated ~100-200 tokens per commit at this detail level
- A chunk of 20 commits ≈ 2,000-4,000 tokens — leaves ample room
  for the analysis prompt and response

For large histories (500+ commits), the skill should summarize
across chunks at the end, merging duplicate themes.

### Edge Cases

| Case | Behavior |
|------|----------|
| No unmined commits | "All commits already mined. Use --force to re-mine or --since for a date range." |
| Empty repo / no commits | "No git history found." |
| Very large history (1000+) | Process in chunks, warn about token cost, suggest `--limit` or `--since` |
| Merge commits | Include — merge commit messages often summarize a whole feature |
| No `.context/` dir | Error: "Run `ctx init` first." |
| Ledger file missing | Create new ledger, treat all commits as unmined |
| Rebase / force-push changed history | Ledger uses commit hashes — orphaned hashes are harmless, new commits get mined |
| No findings in a range | "Processed 20 commits — no context-worthy findings. Ledger updated." |
| User aborts mid-review | Ledger NOT updated — commits stay unmined for next run |

### Validation Rules

- Findings must cite specific commit hashes as evidence
- Decisions must describe a choice between alternatives (not just "we did X")
- Learnings must describe something non-obvious (not "added a test")
- Todos must be actionable (not "maybe someday consider...")
- Duplicate detection: before presenting, check existing DECISIONS.md
  and LEARNINGS.md entries for overlap

## Interface

### Skill

```
/ctx-mine
```

**Trigger phrases:**
- "Mine our git history"
- "What can we learn from our commits?"
- "Extract patterns from git"
- "Bootstrap context from history"
- "Populate context from commits"

**Flags (passed as natural language or explicit):**

| Flag | Default | Purpose |
|------|---------|---------|
| `--since` | Last mined date | Only process commits after this date |
| `--until` | HEAD | Only process commits before this date |
| `--limit N` | 100 | Max commits to process per run |
| `--force` | false | Re-mine already-processed commits |
| `--chunk-size` | 20 | Commits per analysis chunk |
| `--github` | false | Also mine PR descriptions and review comments (v2) |
| `--diff` | false | Include diffs in pass 1 (expensive, more thorough) |

These aren't CLI flags — they're parameters the skill interprets
from the user's natural language or explicit arguments:
- "Mine the last 50 commits" → `--limit 50`
- "Mine everything since January" → `--since 2026-01-01`
- "Re-mine the whole history" → `--force --since=<first-commit>`

## Ledger

### Location

`.context/state/mine-ledger.json`

Follows the existing state directory pattern (alongside `events.jsonl`
and `journal/.state.json`).

### Format

```json
{
  "version": 1,
  "last_mined": "2026-03-01T10:00:00Z",
  "stats": {
    "total_commits_mined": 142,
    "total_findings": 23,
    "total_persisted": 15
  },
  "commits": {
    "abc1234def5678": "2026-03-01",
    "fed8765cba4321": "2026-03-01",
    "...": "..."
  }
}
```

Design choices:
- **Commit hash → date mined**: minimal per-entry data. The date is
  enough for "when did we mine this?" without bloating the file.
- **Stats**: running totals for quick status reporting without
  scanning the full commits map.
- **No per-commit finding details**: findings live in the context
  files they were persisted to. The ledger only tracks "was this
  commit processed?"

### Size Considerations

40-char hash + 12-char date + JSON overhead ≈ 70 bytes per commit.
At 10,000 commits: ~700KB. Acceptable for a state file. For truly
massive repos (100k+ commits), a future version could use date-range
compaction: replace individual hashes with "all commits before X
were mined."

### Ledger Operations (skill-managed)

| Operation | When |
|-----------|------|
| **Read** | Start of every `/ctx-mine` run |
| **Update** | After user confirms findings (not before) |
| **Reset** | User says "re-mine everything" or "reset the ledger" |
| **Stats** | User says "mining status" or "what have we mined?" |

## Skill File Structure

```
internal/assets/claude/skills/ctx-mine/SKILL.md
```

The SKILL.md instructs the agent to:

1. Read the ledger (or create if missing)
2. Run `git log` with appropriate filters
3. Subtract already-mined commits
4. Chunk the remaining commits
5. For each chunk, analyze and extract findings
6. Deduplicate against existing context files
7. Present findings grouped by bucket
8. Wait for user selection
9. Persist approved findings via `ctx add`
10. Update ledger only after successful persistence

### Analysis Prompt (embedded in SKILL.md)

The skill contains a structured analysis prompt for each chunk:

```
You are analyzing git commits to extract project knowledge.
For each commit, consider:

DECISIONS: Did this commit make an architectural choice?
  Look for: library selections, pattern changes, migrations,
  "switch from X to Y", "adopt Z", "replace A with B"

LEARNINGS: Did this commit fix a non-obvious bug or reveal a gotcha?
  Look for: "fix:", bug descriptions, workarounds, edge cases,
  "handle case where...", regression fixes

TODOS: Does this commit mention incomplete work?
  Look for: "TODO", "FIXME", "HACK", "temporary", "Part 1 of N",
  partial migrations, commented-out code notes

SPEC CANDIDATES: Does this commit hint at a larger feature?
  Look for: multiple related commits, "initial support for",
  "phase 1", experimental features, feature flags

MISCELLANEOUS: Anything interesting that doesn't fit above?
  Look for: recurring patterns, naming conventions, structural
  preferences, tooling choices

For each finding, cite the specific commit hash(es) and quote
the relevant commit message text.

Skip: version bumps, merge commits with no useful message,
routine dependency updates, formatting-only changes.
```

## Configuration

| Key | Source | Default | Purpose |
|-----|--------|---------|---------|
| `mine.chunk_size` | `.ctxrc` | 20 | Commits per chunk |
| `mine.default_limit` | `.ctxrc` | 100 | Default `--limit` |
| `mine.exclude_patterns` | `.ctxrc` | `["^Merge branch", "^bump version"]` | Commit message patterns to skip |

Not needed for v1 — hardcode sensible defaults in the skill. Add
`.ctxrc` support when users ask for customization.

## Testing

- **Manual**: Run `/ctx-mine --limit 10` on the ctx repo itself,
  verify findings are sensible and ledger updates correctly
- **Skill review**: Use `/_ctx-skill-creator` to evaluate the
  SKILL.md for clarity, completeness, and trigger reliability
- **Ledger integrity**: After multiple runs, verify no duplicate
  hashes, stats match, dates are monotonic

No Go code to unit test in v1 (skill-only).

## Non-Goals

- **Full diff analysis in v1**: Commit messages + file lists provide
  80% of the signal at 5% of the token cost. Diffs are opt-in.
- **PR comment mining in v1**: Requires `gh` CLI and GitHub access.
  Valuable but adds scope. The `--github` flag is reserved for v2.
- **Automatic persistence**: Every finding goes through human review.
  The user is the quality gate, not a confidence score.
- **Confidence scoring**: Adds state complexity for marginal benefit.
  If a finding is weak, the user skips it. If they want to re-mine,
  `--force` handles it.
- **Cross-repo mining**: Only the current repo's history.
- **CLI command**: No `internal/cli/mine/` package for v1. The skill
  manages everything directly. Add Go code when there's a clear need
  (e.g., programmatic ledger access, CI integration).

## Future Directions (v2+)

- **`--github` flag**: Mine PR descriptions, review comments, and
  issue discussions via `gh api`. PR comments are richer than commit
  messages for learnings and decisions.
- **Domain splitting**: Group commits by file path prefix (frontend/,
  backend/, infra/) and analyze each domain separately for more
  targeted findings.
- **Continuous mining**: A post-commit hook or periodic skill that
  mines incrementally after each commit, surfacing findings while
  they're fresh.
- **Cross-session correlation**: Combine mining findings with journal
  entries to build a richer project narrative.
- **Convention extraction**: Specifically look for repeated code
  patterns across diffs to propose CONVENTIONS.md entries (closer to
  the original article's skill-file generation).
- **`ctx mine ledger` CLI subcommand**: Show stats, reset ranges,
  export mined data for analysis.

## Open Questions

1. **Should merge commits be included or excluded by default?** They
   often summarize features but can also be noise ("Merge branch
   'main' into feature-x"). Current spec: include, with an exclude
   pattern for bare merge messages.

2. **Token cost for large histories**: Mining 500 commits at 20/chunk
   = 25 sub-agent calls. At ~4k tokens per call, that's ~100k tokens.
   Acceptable? Should we warn above a threshold?

3. **Multi-author attribution**: Should findings note who authored the
   relevant commits? Useful for teams, irrelevant for solo projects.
   Leaning toward: include author in evidence, let user decide.
