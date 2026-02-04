---
name: ctx-save
description: "Save session snapshot. Use when significant progress made, complex task completed, or before starting a risky operation."
allowed-tools: Bash(ctx:*)
---

Save the current context state to `.context/sessions/`.

## Usage

```
/ctx-save
/ctx-save auth-refactor
/ctx-save "database migration"
```

Saves session to `.context/sessions/YYYY-MM-DD-<topic>.md`. Topic is optional.

## Execution

```bash
ctx session save $ARGUMENTS
```

Report the saved session file path to the user.
