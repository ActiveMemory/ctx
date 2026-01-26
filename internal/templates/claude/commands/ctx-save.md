---
description: "Save current context to a session file"
argument-hint: "[topic]"
---

Save the current context state to `.context/sessions/`.

**Usage:**
```
/ctx-save
/ctx-save auth-refactor
/ctx-save "database migration"
```

Saves session to `.context/sessions/YYYY-MM-DD-<topic>.md`. Topic is optional.

```!
ctx session save $ARGUMENTS
```

Report the saved session file path to the user.
