# Plan: Restructure Recipes Navigation into Sub-groups

## Context

The recipes section has grown to 20 entries. The flat sidebar is getting
too tall vertically, risking double-scrollbar UX issues. The index page
listing is also hard to scan at this length.

## Approach

Reorganize into 6 collapsible sub-groups in the sidebar. **No file moves** —
all recipe .md files stay in `docs/recipes/`. Only two files change:

### 1. Update `zensical.toml` nav (lines 41-61)

Replace the flat Recipes list with nested sub-groups:

```toml
{ "Recipes" = [
  { "Getting Started" = [
    "recipes/guide-your-agent.md",
    "recipes/multi-tool-setup.md",
    "recipes/external-context.md",
  ]},
  { "Sessions" = [
    "recipes/session-lifecycle.md",
    "recipes/session-ceremonies.md",
    "recipes/session-archaeology.md",
  ]},
  { "Knowledge & Tasks" = [
    "recipes/knowledge-capture.md",
    "recipes/task-management.md",
    "recipes/scratchpad-with-claude.md",
    "recipes/scratchpad-sync.md",
  ]},
  { "Hooks & Notifications" = [
    "recipes/hook-output-patterns.md",
    "recipes/system-hooks-audit.md",
    "recipes/webhook-notifications.md",
  ]},
  { "Maintenance" = [
    "recipes/context-health.md",
    "recipes/claude-code-permissions.md",
    "recipes/permission-snapshots.md",
    "recipes/publishing.md",
  ]},
  { "Agents & Automation" = [
    "recipes/autonomous-loops.md",
    "recipes/when-to-use-agent-teams.md",
    "recipes/parallel-worktrees.md",
  ]},
]},
```

Remove `recipes/index.md` from nav — "Recipes" becomes a pure collapsible
container. The index page still exists for direct linking but doesn't
appear in the sidebar.

### 2. Update `docs/recipes/index.md`

Align the section headings and recipe order to match the new nav grouping.
Replace the current 5 sections (Getting Started, Daily Workflow,
Maintenance, History and Discovery, Advanced) with the 6 new groups.

### Files Modified

- `zensical.toml` — nav restructure
- `docs/recipes/index.md` — section reorder to match nav

### Verification

1. `make site` — build succeeds
2. `make site-serve` — check sidebar renders collapsed sub-groups
3. Verify no broken links in recipes cross-references
