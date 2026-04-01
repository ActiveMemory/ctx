# Context Files Instructions

Files in `.context/` are the project's persistent memory. Follow these
rules when reading or modifying them.

## Reading Context

- Read `.context/CONSTITUTION.md` first: it contains hard rules
- Read `.context/TASKS.md` to understand current work items
- Use `ctx agent` for a token-budgeted summary instead of reading all files

## Modifying Context

- NEVER delete content from context files without explicit user approval
- Use append-only patterns: add new entries, mark old ones complete
- Task format: `- [ ] description #added:YYYY-MM-DD-HHMMSS`
- Decision format: date, decision, rationale as a section entry
- Mark completed tasks with `[x]`, never delete them

## File Permissions

- `CONSTITUTION.md`: Read-only unless user explicitly approves changes
- `TASKS.md`: Append new tasks, mark existing ones complete
- `DECISIONS.md`: Append only
- `LEARNINGS.md`: Append only
- `CONVENTIONS.md`: Append only, propose changes to user first
- `sessions/`: Create new files freely, never modify existing ones
