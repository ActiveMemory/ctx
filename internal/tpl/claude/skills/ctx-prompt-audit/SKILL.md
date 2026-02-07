---
name: ctx-prompt-audit
description: "Audit prompting patterns. Use periodically to help users improve prompt quality and reduce clarification cycles."
---

Analyze recent session transcripts to identify prompts that led to
unnecessary clarification back-and-forth.

## Before Auditing

1. **Check for session data**: look in `.context/sessions/` and
   `.context/journal/` for transcripts to analyze
2. **Need at least 3 sessions**: fewer than that gives too small a
   sample; tell the user to try again later
3. **Confirm scope**: if the user specifies sessions or a date
   range, use that; otherwise default to the 5 most recent

## When to Use

- Periodically to help users improve their prompting
- When the user asks for feedback on their prompting style
- After noticing many clarification cycles in recent sessions
- After a session with unusually high back-and-forth

## When NOT to Use

- Immediately after a user's first session (not enough data)
- When the user is frustrated; coaching lands poorly when someone
  is already annoyed
- Unsolicited; only run when the user invokes it or explicitly
  asks for feedback

## Usage Examples

```text
/ctx-prompt-audit
/ctx-prompt-audit --sessions 10
/ctx-prompt-audit 2026-01-24
```

## Data Sources

Check both locations for session transcripts:

| Source                  | Format                             |
|-------------------------|------------------------------------|
| `.context/sessions/`    | Manual session dumps (markdown)    |
| `.context/journal/`     | Exported session journals (richer) |

Journal entries contain full turn-by-turn conversation and are
the better source when available. Session dumps are shorter
summaries but still useful for pattern detection.

## Process

1. **Gather transcripts**: read 3-5 recent sessions from journal
   or session dumps
2. **Extract user prompts**: isolate the human turns
3. **Identify vague prompts**: flag those that caused clarifying
   questions (see criteria below)
4. **Cross-reference patterns**: look for repeated habits across
   sessions, not one-off mistakes
5. **Generate coaching report**: use the output format below
6. **Present and discuss**: share the report, ask if the user
   wants to dig into any example

## What Makes a Prompt "Vague"

Look for prompts where the agent asked clarifying questions
instead of acting:

- **Missing file context**: "fix the bug" without specifying
  which file or error
- **Ambiguous scope**: "optimize it" without what to optimize
  or success criteria
- **Undefined targets**: "update the component" when multiple
  components exist
- **Missing error details**: "it's not working" without symptoms
- **Vague action words**: "make it better", "clean this up"

## Important Nuance

Not every short prompt is vague. Consider context:
- "fix the bug" after discussing a specific error: **not vague**
- "fix the bug" as the first message: **vague**
- "same:" after a pattern is established: **not vague** (the
  user set a convention and is being efficient)
- Shorthand that references shared context is good prompting,
  not lazy prompting

## Output Format

```markdown
## Prompt Audit Report

**Sessions analyzed**: 5
**User prompts reviewed**: 47
**Vague prompts found**: 4 (8.5%)

---

### Example 1: Missing File Context

**Your prompt**: "fix the bug"

**What happened**: I had to ask which file and what error.

**Better prompt**: "fix the authentication error in
src/auth/login.ts where JWT validation fails with 401"

---

## Patterns to Watch

Based on your sessions, you tend to:
1. Skip mentioning file paths (3 occurrences)
2. Use "it" without establishing what "it" refers to
   (2 occurrences)

## What You Do Well

- You provide error output when debugging (4 of 5 sessions)
- You reference specific files by path in most prompts

## Tips

- Start prompts with the **file path** when discussing
  specific code
- Include **error messages** when debugging
- Specify **success criteria** for optimization tasks
```

## Guidelines

- **Constructive, not critical**: frame suggestions as
  improvements, not corrections
- **Show actual prompts**: quote from their sessions so
  examples are concrete, not hypothetical
- **Explain the consequence**: what happened because the prompt
  was vague (extra round-trip, wrong file edited, etc.)
- **Provide rewrites**: show a concrete better alternative for
  each example
- **Acknowledge strengths**: include a "What You Do Well"
  section; people learn better when not purely criticized
- **Look for patterns**: one vague prompt is noise; three of the
  same kind is a habit worth addressing
- **End with actionable tips**: 3-5 specific, memorable tips

## Quality Checklist

Before presenting the report, verify:
- [ ] At least 3 sessions were analyzed (not a tiny sample)
- [ ] Every "vague" example includes the actual quoted prompt
- [ ] Every example has a concrete rewrite (not just "be more
      specific")
- [ ] Context was considered (short != vague)
- [ ] Report includes positive observations, not just criticism
- [ ] Tips are specific to this user's patterns, not generic
      advice
