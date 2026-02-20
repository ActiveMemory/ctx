# Smart Retrieval: Budget-Aware Context Packet Assembly

Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 1)

## Problem

`ctx agent --budget N` is cosmetic. The budget value is displayed in the
packet header but has zero effect on content selection. The actual behavior:

| Section       | Source            | Selection Logic       | Limit     |
|---------------|-------------------|-----------------------|-----------|
| Constitution  | CONSTITUTION.md   | All checkbox items    | Unlimited |
| Tasks         | TASKS.md          | All unchecked items   | Unlimited |
| Conventions   | CONVENTIONS.md    | First N bullet items  | 5         |
| Decisions     | DECISIONS.md      | Last N titles only    | 3         |
| Learnings     | LEARNINGS.md      | —                     | Excluded  |
| Architecture  | ARCHITECTURE.md   | —                     | Excluded  |
| Glossary      | GLOSSARY.md       | —                     | Excluded  |

Problems:
1. **Budget is ignored.** A project with 50 pending tasks and 30 decisions
   produces the same packet whether budget is 2000 or 20000.
2. **LEARNINGS.md is entirely excluded.** Learnings are often the most
   operationally useful context (gotchas, workarounds, non-obvious behavior).
3. **Decisions are title-only.** The rationale/consequences body — the most
   valuable part — is discarded.
4. **No relevance filtering.** A learning about "hook edge cases" has the same
   weight whether the current task is about hooks or about blog posts.
5. **No graceful degradation.** When budget is tight, there's no summarization
   — entries are either fully included or silently dropped.

## Solution

Make `ctx agent` actually use the budget by scoring entries and allocating
tokens across sections with priority-based budgeting.

### Section Priority Tiers

Sections are allocated budget in tiers. Higher tiers are always fully included
before lower tiers get any budget.

| Tier | Section       | Budget Rule                        | Rationale                    |
|------|---------------|------------------------------------|------------------------------|
| 1    | Constitution  | Always full (all rules)            | Inviolable, small, fixed     |
| 1    | Read Order    | Always full (file list)            | Negligible cost              |
| 1    | Instruction   | Always full (static string)        | Negligible cost              |
| 2    | Tasks         | All unchecked, up to 40% of budget | Current work drives session  |
| 3    | Conventions   | All items, up to 20% of budget     | Always relevant to coding    |
| 4    | Decisions     | Scored, up to 20% of budget        | Full body, not just titles   |
| 5    | Learnings     | Scored, up to 20% of budget        | NEW: finally included        |

**Remaining budget** after Tier 1-3 fills is split between Decisions and
Learnings proportionally. If the combined Tier 4+5 content fits in the
remaining budget, include everything. Otherwise, score and rank.

### Entry Scoring

For Decisions and Learnings, each entry gets a relevance score:

```
score = recency_score + task_relevance_score + superseded_penalty
```

**Recency score** (0.0–1.0):
- Entries from the last 7 days: 1.0
- 8–30 days: 0.7
- 31–90 days: 0.4
- 90+ days: 0.2

**Task relevance score** (0.0–1.0):
- Extract keywords from active task text (split on whitespace, lowercase,
  drop common words, drop short words <3 chars)
- For each entry, count keyword matches in title + body
- Normalize: `min(matches / 3, 1.0)` (3+ matches = full score)

**Superseded penalty**:
- Entries marked `~~Superseded` get score = 0.0 (excluded unless budget
  allows everything)

Final score is `recency + relevance` (range 0.0–2.0). Entries are ranked
by score descending.

### Graceful Degradation

When scored entries exceed their budget allocation:

1. **Include full entries** in score order until budget is ~80% consumed.
2. **Summarize remaining entries** as one-line items: just the title from
   the `## [timestamp] Title` header. Group under a "Also noted:" subheader.
3. **Drop superseded entries** entirely (they're archived context).

This ensures the packet always communicates *what* exists even when it can't
include *why*.

### Token Estimation

Use the existing `context.EstimateTokensString()` (len/4 heuristic) for
budget accounting. It's conservative (overestimates), which is correct for
budgeting — better to include slightly less than to overflow.

## Changes

### Modified Files

| File | Change |
|------|--------|
| `internal/cli/agent/extract.go` | Add `extractLearnings`, `extractDecisionBlocks`, keyword extraction |
| `internal/cli/agent/out.go` | Replace hardcoded assembly with budget-aware assembly |
| `internal/cli/agent/types.go` | Add `Learnings` field to `Packet`, add `ScoredEntry` type |
| `internal/cli/agent/score.go` | **NEW**: Entry scoring logic (recency, task relevance) |
| `internal/cli/agent/budget.go` | **NEW**: Budget allocation and section assembly |
| `internal/cli/agent/score_test.go` | **NEW**: Tests for scoring |
| `internal/cli/agent/budget_test.go` | **NEW**: Tests for budget allocation |

### New Types

```go
// ScoredEntry is an entry block with a computed relevance score.
type ScoredEntry struct {
    index.EntryBlock
    Score    float64
    Tokens   int      // pre-computed token estimate of full body
    Included bool     // true = full body, false = title-only summary
}

// SectionBudget tracks token allocation for a packet section.
type SectionBudget struct {
    Name      string
    MaxTokens int      // allocated budget
    Used      int      // tokens consumed
    Items     []string // rendered items for output
}
```

### New Extraction Functions

```go
// extractDecisionBlocks parses DECISIONS.md into scored entries.
// Uses index.ParseEntryBlocks for parsing, then scores each block.
func extractDecisionBlocks(ctx *context.Context, keywords []string) []ScoredEntry

// extractLearningBlocks parses LEARNINGS.md into scored entries.
func extractLearningBlocks(ctx *context.Context, keywords []string) []ScoredEntry

// extractTaskKeywords extracts keywords from active task text for
// relevance matching. Splits on whitespace, lowercases, drops stop
// words and words shorter than 3 characters.
func extractTaskKeywords(tasks []string) []string
```

### Budget Assembly Flow

```
runAgent(cmd, budget, format, cooldown, session)
  |
  context.Load("")
  |
  +--> Tier 1: constitution, readOrder, instruction (always full)
  |    -> subtract token cost from remaining budget
  |
  +--> Tier 2: tasks (all unchecked, cap at 40% of original budget)
  |    -> if tasks exceed cap, include newest first, summarize rest
  |    -> subtract from remaining
  |
  +--> Tier 3: conventions (all items, cap at 20% of original budget)
  |    -> if conventions exceed cap, include first N that fit
  |    -> subtract from remaining
  |
  +--> extract keywords from included tasks
  |
  +--> Tier 4+5: score decisions + learnings
  |    -> rank by score
  |    -> fill remaining budget: full entries first, then summaries
  |
  +--> assemble Packet / Markdown output
```

### Output Changes

**Markdown format**: New "Recent Learnings" section. Decision entries now
include body content (not just titles). Entries that didn't fit get a
"Also noted:" section with title-only summaries.

**JSON format**: New `learnings` field (array of strings). Decision entries
contain full body. New `summaries` field for entries that were title-only.

### Backward Compatibility

- Default `--budget 8000` produces a richer packet than before (includes
  learnings, decision bodies) but stays within budget.
- The existing extraction functions (`extractConventions`, etc.) remain
  available — the new budget system calls them internally.
- JSON `Packet` struct gains new fields (`learnings`, `summaries`) — additive,
  not breaking for consumers that ignore unknown fields.

## Non-Goals

- **Semantic/embedding-based similarity**: Too heavy for a CLI tool. Keyword
  matching is sufficient for structured Markdown entries.
- **Cross-file references**: Scoring entries based on which files they mention
  is interesting but adds complexity. Revisit later.
- **User-tunable weights**: Hardcode sensible defaults. If users need
  different weights, `.contextrc` extension is a future enhancement.
- **Architecture/Glossary inclusion**: These files are reference material,
  not session-specific. Keep them in the "read these files" list without
  extracting content. Revisit if users request it.
- **Consolidation (Issue #19 Phase 3)**: Separate skill, separate spec.
  Smart retrieval makes consolidation less urgent by surfacing the right
  entries regardless of file size.

## Testing

### Unit Tests

- `score_test.go`:
  - Recency scoring with various age brackets
  - Task keyword extraction (stop words, short words, deduplication)
  - Relevance scoring with 0, 1, 3+ keyword matches
  - Superseded entries get score 0.0
  - Score ordering (newer + relevant > older + irrelevant)

- `budget_test.go`:
  - Tier 1 always fits (even with budget=100)
  - Tasks capped at 40% with overflow to summaries
  - Conventions capped at 20%
  - Remaining budget split between decisions and learnings
  - Graceful degradation: full entries → summaries → titles
  - Empty files produce empty sections (no errors)
  - Budget=0 produces constitution + instruction only

### Integration

- Existing `ctx agent` tests should pass unchanged (output is richer but
  still valid Markdown/JSON).
- New test: load a `.context/` fixture with large DECISIONS.md and
  LEARNINGS.md, verify budget is respected (output token count ≤ budget).

## Implementation Order

1. Add `score.go` with scoring functions + tests
2. Add `budget.go` with allocation logic + tests
3. Update `extract.go` with new extraction functions
4. Update `types.go` with new Packet fields
5. Update `out.go` to use budget-aware assembly
6. Update existing tests
7. Manual verification with real `.context/` directory

## Open Questions (from Issue #19)

> Should `ctx agent` budget allocation weights be configurable?

Not in v1. Hardcode sensible defaults. If the issue comes up, add
`.contextrc` keys later.

> Is 30 the right soft cap for learnings?

Soft caps are Phase 2 (drift nudges), not this spec. But for reference,
the retrieval system makes the cap less critical — even with 50 learnings,
only the relevant ones consume budget.
