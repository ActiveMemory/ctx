---
description: "Generate a Ralph loop script"
argument-hint: "[--tool claude|aider] [--prompt FILE] [--max-iterations N]"
---

Generate a ready-to-use Ralph loop shell script.

**Usage:**
```
/ctx-loop
/ctx-loop --tool aider
/ctx-loop --prompt PROMPT.md --max-iterations 10
```

Generates a shell script for iterative AI development. Defaults to Claude Code.

```!
ctx loop $ARGUMENTS
```

Report the generated script path and how to run it.
