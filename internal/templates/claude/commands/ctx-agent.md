---
description: "Get AI-ready context packet"
argument-hint: "[--budget N]"
---

Load the full context packet for AI consumption.

**Usage:**
```
/ctx-agent
/ctx-agent --budget 4000
/ctx-agent --budget 8000
```

Returns context optimized for AI assistants. Use `--budget` to limit token count.

```!
ctx agent $ARGUMENTS
```

This provides the complete context state optimized for AI assistants.
