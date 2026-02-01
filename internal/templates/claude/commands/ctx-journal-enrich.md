---
description: "Enrich a journal entry with frontmatter and tags"
---

Enrich a session journal entry with structured metadata.

## Input

The user specifies a journal entry by partial match:
- `twinkly-stirring-kettle` (slug)
- `twinkly` (partial slug)
- `2026-01-24` (date)
- `76fe2ab9` (short ID)

Find matching files:
```!
ls .context/journal/*.md | grep -i "<pattern>"
```

If multiple matches, show them and ask which one.
If no argument given, show recent entries and ask.

## Enrichment Tasks

Read the journal entry and extract:

### 1. Frontmatter (YAML at top of file)
```yaml
---
title: "Session title"
date: 2026-01-27
type: feature|bugfix|refactor|exploration|debugging|documentation
outcome: completed|partial|abandoned|blocked
topics:
  - authentication
  - caching
technologies:
  - go
  - postgresql
libraries:
  - cobra
  - fatih/color
error_types:
  - nil-pointer
  - timeout
key_files:
  - internal/auth/token.go
  - internal/db/cache.go
---
```

### 2. Summary
If `## Summary` says "[Add your summary...]", replace with 2-3 sentences.

### 3. Extracted Items
Scan the conversation and extract:

**Decisions made** - Link to DECISIONS.md if persisted:
```markdown
## Decisions
- Used Redis for caching ([D12](../DECISIONS.md#d12))
- Chose JWT over sessions (not yet persisted)
```

**Learnings discovered** - Link to LEARNINGS.md if persisted:
```markdown
## Learnings
- Token refresh requires cache invalidation ([L8](../LEARNINGS.md#l8))
- Go's defer runs LIFO (new insight)
```

**Tasks completed/created**:
```markdown
## Tasks
- [x] Implement caching layer
- [ ] Add cache metrics (created this session)
```

## Process

1. Find and read the journal file
2. Analyze the conversation
3. Propose enrichment (type, topics, outcome)
4. Ask user for confirmation/adjustments
5. Show diff and write if approved

## Example

User: `/ctx-journal-enrich twinkly`

You:
> Found: `2026-01-24-twinkly-stirring-kettle-76fe2ab9.md`
>
> After reading, I'd categorize this as:
> - **Type**: refactor
> - **Outcome**: completed
> - **Topics**: cli, code-cleanup, conventions
> - **Technologies**: go
> - **Libraries**: cobra, spf13/pflag
> - **Key files**: internal/cli/task/tasks.go, internal/cli/session/run.go
>
> I also found:
> - 1 decision: "Use cmd.Print* instead of fmt.Print*"
> - 2 learnings: AST parsing gotcha, flag naming convention
> - 3 tasks completed
>
> Want me to add these? Any adjustments to the categories?
