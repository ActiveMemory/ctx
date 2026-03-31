# Stuttery name cleanup

## Problem

288 exported symbols across the codebase contain their package name,
causing stutter when qualified: `score.ScoreEntry`, `format.FormatSize`,
`token.EstimateTokens`, `tpl.TplJournalSiteReadme`.

Go convention: the package name is part of the qualified name, so
symbols should not repeat it.

## Scope

| Category | Count | Example | Fix pattern |
|---|---|---|---|
| config/embed/text DescKey | 57 | `text.DescKeyCheckContextSize...` | Keep — these are YAML key paths, stutter is the key format not the Go name |
| config/* constants | 55 | `event.EventSilent` → `event.Silent` | Drop package prefix |
| write/* output functions | 38 | `remind.ReminderAdded` → `remind.Added` | Drop package prefix |
| cli/* command functions | 86 | `score.ScoreEntry` → `score.Entry` | Drop prefix; think about semantics |
| tpl/* template constants | 26 | `tpl.TplJournalSiteReadme` → `tpl.JournalSiteReadme` | Drop Tpl prefix |
| other core packages | 26 | `token.EstimateTokens` → `token.Estimate` | Drop package name from symbol |

## Exclusions

- **config/embed/text DescKey constants** (57): These map to YAML keys
  like `"check-context-size.billing-box-title"`. The stutter is in the
  constant name matching the key path — changing the Go name without
  changing the YAML key creates a naming disconnect. Exclude from this
  cleanup.

## Decision: naming strategy

Not just "drop the prefix" — think about what each name means:

| Current | Naive fix | Better name | Why |
|---|---|---|---|
| `score.ScoreEntry` | `score.Entry` | `score.Entry` | Entry is what we score |
| `score.ScoreEntries` | `score.Entries` | `score.Entries` | Plural of above |
| `score.RecencyScore` | `score.Recency` | `score.Recency` | The score type is implied |
| `score.RelevanceScore` | `score.Relevance` | `score.Relevance` | Same |
| `collapse.CollapseToolOutputs` | `collapse.ToolOutputs` | `collapse.ToolOutputs` | Action (collapse) is the package |
| `format.FormatSize` | `format.Size` | `format.Size` | Action is the package |
| `generate.GenerateIndex` | `generate.Index` | `generate.Index` | Action is the package |
| `normalize.NormalizeContent` | `normalize.Content` | `normalize.Content` | Action is the package |
| `parse.ParseJournalEntry` | `parse.JournalEntry` | `parse.JournalEntry` | Action is the package |
| `tpl.TplJournalSiteReadme` | `tpl.JournalSiteReadme` | `tpl.JournalSiteReadme` | Tpl is redundant with package |
| `event.EventSilent` | `event.Silent` | `event.Silent` | Event is the package |
| `heartbeat.HeartbeatCounterPrefix` | `heartbeat.CounterPrefix` | `heartbeat.CounterPrefix` | Heartbeat is the package |
| `token.EstimateTokens` | `token.Estimate` | `token.Estimate` | Token is the package |
| `remind.ReminderAdded` | `remind.Added` | `remind.Added` | Reminder is the domain |
| `status.StatusHeader` | `status.Header` | `status.Header` | Status is the package |

## Phases

### Phase 1: cli/* functions (86) — highest impact, most callers
### Phase 2: config/* constants (55) — many callers via config refs
### Phase 3: write/* output functions (38) — moderate callers
### Phase 4: tpl/* template constants (26) — moderate callers
### Phase 5: other core packages (26) — scattered callers

Each phase: rename symbols, update all callers, build, test, commit.

## Non-goals

- Renaming config/embed/text DescKey constants (57) — YAML key coupling
- Renaming private functions (convention applies to exported only)
- Changing package names themselves
