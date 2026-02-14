# ctx journal obsidian — Obsidian Vault Export

## Problem

The journal site generator produces a zensical-compatible static site. Users
who prefer Obsidian for knowledge management cannot import the journal into
their vault with proper navigation, graph connectivity, or tag integration.
The enriched journal entries are already 80% Obsidian-compatible (YAML
frontmatter, Markdown content), but the linking model and navigation structure
are incompatible.

## Solution

A new subcommand `ctx journal obsidian` that exports enriched journal entries
as an Obsidian vault. Reuses the existing scan/parse/index infrastructure from
`internal/cli/journal/` and adds an Obsidian-specific output layer.

## Design Principles

- **Reuse, don't duplicate**: Use the same `scanJournalEntries`,
  `buildTopicIndex`, `buildKeyFileIndex`, `buildTypeIndex` functions.
  The new command is an alternative output backend, not a parallel pipeline.
- **Wikilinks native**: All internal links use `[[target|display]]` format.
  No markdown relative links.
- **Graph-optimized**: Every entry links to its topic/type MOCs. MOCs link
  to entries. This creates hub-and-spoke graph topology.
- **Minimal vault config**: Just `.obsidian/` with `app.json` enforcing
  wikilink mode. Obsidian auto-generates everything else on first open.
- **Non-destructive**: Writes to a separate output directory (default
  `.context/journal-obsidian/`). Never modifies source journal files.

## CLI Surface

```
ctx journal obsidian [--output DIR]
```

### Flags

| Flag       | Default                          | Description                   |
|------------|----------------------------------|-------------------------------|
| `--output` | `.context/journal-obsidian/`     | Output directory for vault    |

No `--build` or `--serve` flags — Obsidian opens the folder directly.

## Output Structure

```
journal-obsidian/
├── .obsidian/
│   └── app.json                    # Minimal config (wikilinks on)
├── Home.md                         # Root MOC
├── entries/
│   ├── 2026-01-23--76fe2ab9.md     # Transformed journal entries
│   └── ...
├── topics/
│   ├── _Topics.md                  # Topics MOC
│   ├── context-persistence.md      # Per-topic page (popular only)
│   └── ...
├── files/
│   ├── _Key Files.md               # Files MOC
│   ├── internal-cli-journal.md     # Per-file page (popular only)
│   └── ...
└── types/
    ├── _Session Types.md           # Types MOC
    ├── feature.md                  # Per-type page
    └── ...
```

### Why `_` prefix on MOC filenames

Obsidian sorts alphabetically. `_Topics.md` sorts before entry files,
making MOCs easy to find in the file explorer. This is a common Obsidian
convention.

## Transformations

### 1. Frontmatter Mapping

Source enriched frontmatter:

```yaml
---
title: "Session title"
date: 2026-01-23
type: feature
outcome: completed
topics:
  - context-persistence
  - journal
technologies:
  - go
  - cobra
key_files:
  - internal/cli/journal/run.go
---
```

Obsidian output frontmatter:

```yaml
---
title: "Session title"
date: 2026-01-23
type: feature
outcome: completed
tags:
  - context-persistence
  - journal
technologies:
  - go
  - cobra
key_files:
  - internal/cli/journal/run.go
aliases:
  - "Session title"
---
```

Changes:
- `topics` renamed to `tags` (Obsidian-recognized key)
- `aliases` added with the title (makes entries findable by title in search)
- All other fields preserved as custom properties

### 2. Link Conversion

All markdown links to journal entries become wikilinks:

| Source (site format)                        | Obsidian output                              |
|---------------------------------------------|----------------------------------------------|
| `[title](2026-01-23-slug.md)`              | `[[2026-01-23-slug\|title]]`                 |
| `[← Previous](file-p1.md)`                 | `[[file-p1\|← Previous]]`                    |
| `[topic](../topics/caching.md)`            | `[[caching\|topic]]`                         |

External links (`https://...`) are left unchanged.

### 3. Related Sessions Footer

Each entry gets a footer section linking to related entries (same topics):

```markdown
---

## Related Sessions

**Topics**: [[_Topics|Topics MOC]] · [[context-persistence]] · [[journal]]

**Type**: [[feature]]

**See also**:
- [[2026-01-24-other-session|Other session title]]
- [[2026-01-25-another-session|Another session title]]
```

This creates bidirectional graph edges between entries that share topics.

### 4. MOC Pages

#### Home.md (root MOC)

```markdown
# Session Journal

Navigation hub for all journal entries.

## Browse by
- [[_Topics|Topics]] — sessions grouped by topic
- [[_Key Files|Key Files]] — sessions grouped by file touched
- [[_Session Types|Session Types]] — sessions grouped by type

## Recent Sessions

- [[2026-02-14-entry|Entry title]] — `feature` · `completed`
- [[2026-02-13-entry|Entry title]] — `bugfix` · `completed`
- ...
```

#### Topic MOC (_Topics.md)

Same structure as the site's `topics/index.md` but with wikilinks:

```markdown
# Topics

## Popular Topics

- [[context-persistence]] (12 sessions)
- [[journal]] (8 sessions)

## Long-tail Topics

- **caching** — [[2026-01-23-entry|Session title]]
```

Popular topics get dedicated pages; long-tail topics link inline
(same threshold logic as the site generator).

#### Individual Topic Page

```markdown
# context-persistence

12 sessions on this topic.

## 2026-02

- [[2026-02-14-entry|Entry title]] — `feature` · `completed`
- [[2026-02-13-entry|Entry title]] — `bugfix` · `partial`

## 2026-01

- [[2026-01-28-entry|Entry title]] — `refactor` · `completed`
```

Key file and type MOCs follow the same pattern.

## Content Normalization

Apply the same normalization pipeline as the site generator:
- `stripSystemReminders`
- `cleanToolOutputJSON`
- `consolidateToolRuns`
- `mergeConsecutiveTurns`
- `softWrapContent`

But do NOT modify source files (unlike `runJournalSite` which writes back).
Only write normalized content to the output directory.

## .obsidian/app.json

Minimal config:

```json
{
  "useMarkdownLinks": false,
  "showFrontmatter": true,
  "strictLineBreaks": false
}
```

This ensures:
- New links created in Obsidian use wikilink format (matching our output)
- Frontmatter is visible in the properties panel
- Markdown rendering matches standard behavior

## Non-Goals

- **No Obsidian plugin integration**: We produce vanilla Markdown that works
  with core Obsidian. No Dataview queries, no Templater, no community plugins.
- **No incremental sync**: Each run produces a full vault. If the user wants
  incremental updates, that's a future feature.
- **No theme customization**: `.obsidian/` contains only `app.json`.
  Users customize the vault in Obsidian after import.
- **No modification of source entries**: Unlike `journal site` which
  soft-wraps source files in-place, `journal obsidian` is read-only on
  the source journal.

## Error Cases

| Condition                    | Behavior                                  |
|------------------------------|-------------------------------------------|
| No `.context/journal/`       | Error: "No journal directory found"       |
| No entries found             | Error: "No journal entries found"         |
| Output dir exists            | Overwrite (same as `journal site`)        |
| Entry without frontmatter    | Copy as-is, no related footer, warn       |
| Entry without topics         | No topic links in footer, no graph edges  |

## Implementation Notes

### Code Organization

New files in `internal/cli/journal/`:

| File            | Purpose                                            |
|-----------------|----------------------------------------------------|
| `obsidian.go`   | `journalObsidianCmd()` — Cobra command definition  |
| `vault.go`      | `runJournalObsidian()` — orchestration             |
| `wikilink.go`   | Link conversion and wikilink formatting helpers     |
| `moc.go`        | MOC page generation (Home, Topics, Files, Types)   |

### Reused Functions

From existing code (no changes needed):

- `scanJournalEntries()` — scan and parse journal files
- `buildTopicIndex()` — aggregate by topic
- `buildKeyFileIndex()` — aggregate by file
- `buildTypeIndex()` — aggregate by type
- `groupByMonth()` — month grouping for page sections
- `continuesMultipart()` — detect part 2+ files
- `softWrapContent()`, `mergeConsecutiveTurns()`, `consolidateToolRuns()`,
  `cleanToolOutputJSON()`, `stripSystemReminders()` — normalization pipeline

### New Constants

Add to `internal/config/`:

```go
// Obsidian vault constants
const (
    ObsidianDirName       = "journal-obsidian"
    ObsidianDirEntries    = "entries"
    ObsidianConfigDir     = ".obsidian"
    ObsidianAppConfig     = `{"useMarkdownLinks":false,"showFrontmatter":true,"strictLineBreaks":false}`
    ObsidianHomeMOC       = "Home.md"
    ObsidianMOCPrefix     = "_"
    ObsidianTopicsMOC     = "_Topics.md"
    ObsidianFilesMOC      = "_Key Files.md"
    ObsidianTypesMOC      = "_Session Types.md"
)
```

### Wikilink Conversion Strategy

1. Parse content line-by-line
2. Match markdown links with regex: `\[([^\]]+)\]\(([^)]+)\)`
3. For each match:
   - If target starts with `http://` or `https://` — skip (external)
   - Strip `.md` extension from target
   - Strip path prefix (`../topics/` etc.)
   - Emit `[[target|display]]`
4. Handle the source link injection differently — instead of the
   `[View source](file://...)` link, just add the source path as a
   frontmatter field: `source_file: .context/journal/filename.md`

## Testing

- Unit tests for wikilink conversion (markdown link → wikilink)
- Unit tests for frontmatter transformation (topics → tags, aliases)
- Unit tests for MOC generation (verify wikilink format in output)
- Integration test: run full pipeline on test fixtures, verify vault structure
