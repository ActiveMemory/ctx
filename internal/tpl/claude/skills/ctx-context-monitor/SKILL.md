---
name: ctx-context-monitor
description: "Respond to context checkpoint signals. Triggered automatically by the check-context-size hook — not user-invocable."
---

When you see a "Context Checkpoint" message from the UserPromptSubmit hook:

1. **Assess** your current context usage relative to the model's context window
2. **If usage appears high (>80%)**:
   - Inform the user concisely: "Context is getting full. Consider wrapping up or starting a new session."
   - Offer to persist any unsaved learnings, decisions, or session notes
3. **If usage is moderate**, continue silently — do not mention the checkpoint

Do NOT mention the checkpoint mechanism unless the user asks about it.
