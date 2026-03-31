---
name: ctx-journal-enrich-all
description: "Full journal pipeline: import, catalog, then semantically enrich all unenriched entries. Use when the user says 'process the journal' or to catch up on the backlog."
allowed-tools: Bash(ctx:*), Read, Glob, Grep, Edit, Write, Task, Agent
---

Full journal pipeline: import, catalog, then deep semantic enrichment.

The value of this skill is **semantic analysis** — reading the actual
conversation content and extracting meaning that no heuristic can
infer from filenames alone. If the output could have been produced
by a regex, the skill is not doing its job.

## When to Use

- When the user says "enrich everything" or "process the journal"
- When there is a backlog of unenriched or unimported sessions
- Periodically to catch up on recent sessions
- After the `check-journal` hook reports unenriched entries

## When NOT to Use

- For a single specific session (use `/ctx-journal-enrich` instead)

## Process

### Phase 0: Import If Needed

Check for unimported sessions. If the journal directory has no
`.md` files, or `.jsonl` session files are newer than the newest
journal entry, import first:

```bash
ctx journal import --all --yes
```

Report how many were imported (or "none needed").

### Phase 1: Build the Eligible List

**Scope: ALL entries in the journal directory, not just newly
imported ones.** The backlog is cumulative.

An entry needs enrichment if its YAML frontmatter is missing
**any** of: `type`, `outcome`, `topics`, `technologies`,
`summary`, `key_files`. Check for the presence of enrichment-
only fields — `type` and `outcome` are never set by import.

**Filter out:**
- **Locked entries**: frontmatter `locked: true` or state file
  says locked. Never modify locked files.
- **Multi-part continuations**: files ending in `-p2.md`,
  `-p3.md`, etc. Enrich only the first part.
- **Suggestion sessions**: files under ~20 lines.

Report the count: "N entries eligible for enrichment, M filtered."

### Phase 2: Catalog Pass (Read-Only)

Before enriching any file, **read all eligible entries** and build
a catalog. This gives the agent a landscape view before committing
to labels.

For each entry, read the first 150-300 lines (enough to understand
the session) and note:
- What the session was about (one line)
- Candidate type (see starter list below, but discover new ones)
- Key topics mentioned
- Technologies and libraries used

**Output the catalog as a table** for the user to review:

```markdown
| # | Date | Title | Type | Topics |
|---|------|-------|------|--------|
| 1 | 2026-01-27 | Implement caching layer | feature | caching, redis, auth |
| 2 | 2026-01-28 | Fix OOM in parser | bugfix | memory, parser |
```

**Discover types organically.** The starter list is a floor, not
a ceiling:

| Starter type    | Use when                              |
|-----------------|---------------------------------------|
| feature         | Building new functionality            |
| bugfix          | Fixing broken behavior                |
| refactor        | Restructuring without behavior change |
| exploration     | Research, learning, experimentation   |
| debugging       | Investigating issues (not yet fixed)  |
| documentation   | Writing docs, comments, blog posts    |
| maintenance     | Cleanup, audits, dependency updates   |
| planning        | Design, spec writing, brainstorming   |
| review          | Code review, PR review, audit review  |
| migration       | Moving between systems or patterns    |
| release         | Release prep, versioning, changelogs  |
| configuration   | Settings, hooks, CI/CD, tooling setup |
| onboarding      | Session start, context loading         |

If a session doesn't fit any of these, **create a new type** and
add it to the catalog. After the catalog pass, report any new
types discovered so the user can approve them.

Similarly for topics: popular topics that appear 3+ times become
**canonical tags**. Long-tail topics that appear once are fine too
— they aid search. Report the top 10 topics and any surprising
clusters.

**Wait for user confirmation** before proceeding to Phase 3.

### Phase 3: Semantic Enrichment

For each entry, apply the full enrichment from `/ctx-journal-enrich`:

**Frontmatter fields** (add to existing YAML block):

```yaml
type: feature
outcome: completed
topics:
  - authentication
  - caching
technologies:
  - go
  - redis
libraries:
  - cobra
  - fatih/color
key_files:
  - internal/auth/token.go
  - internal/db/cache.go
summary: >-
  Implemented Redis-backed caching layer for auth tokens.
  Reduced p99 latency from 450ms to 12ms. Added cache
  invalidation on token refresh.
```

**All enrichment fields:**

| Field          | Required | Description                        |
|----------------|----------|------------------------------------|
| `type`         | Yes      | Session type (from catalog)        |
| `outcome`      | Yes      | completed/partial/abandoned/blocked|
| `topics`       | Yes      | 2-7 topic tags, kebab-case         |
| `technologies` | Yes      | Languages, frameworks, tools       |
| `libraries`    | If any   | Specific libraries/packages used   |
| `key_files`    | If any   | Files that were central to work    |
| `summary`      | Yes      | 2-4 sentences, specific not generic|

**Body sections** (append after frontmatter if conversation
contains them):

- `## Decisions` — design choices made, link to DECISIONS.md
  if persisted
- `## Learnings` — gotchas discovered, link to LEARNINGS.md
  if persisted
- `## Tasks` — tasks completed or created during the session

**Auto-populated fields** (set by import, NEVER overwrite):
`date`, `time`, `project`, `session_id`, `model`, `tokens_in`,
`tokens_out`, `branch`

**Quality bar for summaries:** A good summary answers "what was
accomplished and why does it matter?" A bad summary restates the
title. Compare:
- Bad: "Worked on caching implementation."
- Good: "Implemented Redis-backed token cache, reducing auth
  latency from 450ms to 12ms. Added invalidation on refresh to
  prevent stale tokens."

**No angle brackets in YAML values.** Literal `<` and `>` in
frontmatter strings break HTML rendering on the journal site.
Replace angle brackets with parentheses in summaries, titles,
and all other YAML string fields:
- `<pre><code>` → `(pre)(code)`
- `<slug>--<sha8>` → `(slug)--(sha8)`
- `<command-message>` → `(command-message)`

Do NOT use HTML entities (`&lt;`/`&gt;`) — another agent may
unescape or double-escape them. Do NOT use square brackets —
`[text]` can become Markdown link syntax. Parentheses have no
special meaning in Markdown, HTML, or YAML.

### Phase 4: Mark and Report

After each file, update state:

```bash
ctx system mark-journal <filename> enriched
```

Final report:
- Entries enriched (with type breakdown)
- New types discovered (if any)
- Top 10 canonical topics
- Entries skipped (locked, short, multipart)
- Remind: `ctx journal site --build` or `make journal`

## Parallelization

For backlogs over 10 entries, use subagents. Each agent gets a
batch of 5-8 entries plus the catalog from Phase 2. The catalog
ensures consistent type/topic usage across agents.

## Heuristic Fallback

For backlogs over 50 entries where agent enrichment is too slow,
use `references/enrich-heuristic.py` as a **first pass**, then
follow up with agent enrichment on the heuristic output to add
`libraries`, `key_files`, `summary`, and body sections that the
heuristic cannot produce.

## Quality Checklist

- [ ] All entries in the journal were scanned, not just new ones
- [ ] Catalog pass completed before any writes
- [ ] User approved the catalog before enrichment
- [ ] Each entry has ALL enrichment fields (not just type+outcome)
- [ ] Summaries are specific, not restated titles
- [ ] New types/topics were reported to the user
- [ ] key_files and libraries populated where applicable
- [ ] Body sections (Decisions/Learnings/Tasks) added where found
- [ ] State file updated for each enriched entry
