# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-02-04 | Skills that restate or contradict Claude Code's built-in system prompt create tension, not clarity. The system prompt already covers: avoid over-engineering, don't add unnecessary features, prefer simplicity. Skills should complement the system prompt, not compete with it. Before writing a skill, check if the guidance already exists in the platform. |
| 2026-02-04 | Skill files that suppress AI judgment are jailbreak patterns, not productivity tools. Red flags: <EXTREMELY-IMPORTANT> urgency tags, 'you cannot rationalize' overrides, tables that label hesitation as wrong, absurdly low thresholds (1%). The fix for 'AI forgets skills' is better skill descriptions, not overriding reasoning. Discard these entirely — nothing is salvageable. |
<!-- INDEX:END -->

## [2026-02-04-193920] Skills that restate or contradict Claude Code's built-in system prompt create tension, not clarity. The system prompt already covers: avoid over-engineering, don't add unnecessary features, prefer simplicity. Skills should complement the system prompt, not compete with it. Before writing a skill, check if the guidance already exists in the platform.

**Context**: Reviewing entropy.txt skill that duplicated system prompt guidance about code minimalism

**Lesson**: Skills that conflict with system prompts cause unpredictable behavior — the AI has to reconcile contradictory instructions

**Application**: When evaluating or writing skills, first check Claude Code's system prompt defaults. Only create skills for guidance the platform does NOT already provide.

---

## [2026-02-04-192812] Skill files that suppress AI judgment are jailbreak patterns, not productivity tools. Red flags: <EXTREMELY-IMPORTANT> urgency tags, 'you cannot rationalize' overrides, tables that label hesitation as wrong, absurdly low thresholds (1%). The fix for 'AI forgets skills' is better skill descriptions, not overriding reasoning. Discard these entirely — nothing is salvageable.

**Context**: Reviewing power.txt skill that forced skill invocation on every message

**Lesson**: Jailbreak-structured prompts should be identified and discarded, not refined

**Application**: When evaluating skills, check for judgment-suppression patterns before assessing content
