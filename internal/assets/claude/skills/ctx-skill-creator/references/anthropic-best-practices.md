# Anthropic Prompting Best Practices for Skill Auditing

Condensed from Anthropic's official prompting best practices
documentation. This reference covers principles relevant to
writing and evaluating Claude Code skills (agent instructions).

## Table of Contents

1. [Clarity and Directness](#clarity-and-directness)
2. [Positive Framing](#positive-framing)
3. [Context and Motivation](#context-and-motivation)
4. [Examples](#examples)
5. [XML Structure](#xml-structure)
6. [Tool Use Guidance](#tool-use-guidance)
7. [Subagent Orchestration](#subagent-orchestration)
8. [Autonomy and Safety](#autonomy-and-safety)
9. [Overtriggering and Verbosity](#overtriggering-and-verbosity)
10. [Overengineering](#overengineering)
11. [Long-Horizon and State Management](#long-horizon-and-state-management)
12. [Hallucination Prevention](#hallucination-prevention)

---

## Clarity and Directness

Claude responds well to clear, explicit instructions. Vague
prompts produce vague results.

**Golden rule:** show the prompt to a colleague with minimal
context. If they'd be confused, Claude will be too.

- Be specific about desired output format and constraints.
- Use numbered lists when step order or completeness matters.
- Provide sequential steps as ordered lists, not prose paragraphs.

## Positive Framing

Tell Claude what to do, not what to avoid. Positive instructions
("write flowing prose paragraphs") outperform negative ones
("don't use markdown").

**Why this matters for skills:** a skill full of "don't do X"
leaves the agent guessing what *to* do. Positive instructions
give a clear target. Negative guards are fine as supplements,
but the primary instruction should describe the desired behavior.

<example>
<poor>Do not use markdown in your response.</poor>
<good>Write your response as flowing prose paragraphs.</good>
</example>

<example>
<poor>NEVER use ellipses.</poor>
<good>Your response will be read aloud by a text-to-speech
engine, so avoid ellipses — the engine cannot pronounce them.</good>
</example>

## Context and Motivation

Explain *why* an instruction matters. Claude generalizes from
reasoning better than it memorizes rigid rules. A rule with
motivation lets the model adapt to edge cases the rule author
didn't anticipate.

- Instead of "ALWAYS sort by date": explain that users typically
  want the most recent items first.
- Instead of "NEVER skip tests": explain that untested code
  creates silent regressions that compound.

**Pattern to watch for:** instructions that rely heavily on
MUST, NEVER, ALWAYS, CRITICAL in caps without explaining the
consequence. These suggest missing motivation.

## Examples

Few-shot examples are one of the most reliable ways to steer
output format, tone, and structure. Skills that describe complex
output without showing it tend to drift over time.

Best practices for examples in skills:
- **Relevant**: mirror realistic use cases, not toy scenarios.
- **Diverse**: cover edge cases; avoid patterns Claude might
  overfit to from a single example.
- **Structured**: wrap examples in `<example>` tags so Claude
  distinguishes them from instructions.
- 3-5 examples is the sweet spot for reliable format adherence.

**Good/bad pairs** set boundaries without being prescriptive:

```
**Bad:** "Should pass now" (claimed without evidence)
**Good:** ran `make audit` -> "All checks pass" (evidence-based)
```

## XML Structure

XML tags help Claude parse complex prompts unambiguously. When
a skill mixes instructions, context, examples, and variable
inputs, wrapping each type in its own tag reduces misinterpretation.

Best practices:
- Use consistent, descriptive tag names across prompts.
- Nest tags when content has natural hierarchy.
- Tags are especially valuable when the skill injects external
  content (file contents, user input, tool output) alongside
  instructions — the tags prevent the agent from confusing
  injected content with skill instructions.

**When XML tags help most:** skills that template in variable
content (code snippets, file paths, user descriptions) alongside
fixed instructions. Without delimiters, the model may treat
injected content as part of the instruction.

## Tool Use Guidance

Claude Opus 4.6 follows explicit tool instructions well. Key
patterns:

- Be explicit about which tool to use and when: "Use the Edit
  tool for modifications" beats "modify the file."
- If a skill references tools, state expected behavior clearly:
  "Read the file first, then edit" — not "look at the file."

**Overtriggering risk:** Claude 4.5/4.6 models are more
responsive to system prompts than earlier models. Skills that
were written to combat undertriggering (with aggressive language
like "CRITICAL: You MUST use this tool") may now overtrigger.
Dial back to natural phrasing: "Use this tool when..."

## Subagent Orchestration

Claude Opus 4.6 has a strong predilection for spawning
subagents and may do so when a simpler, direct approach suffices
(e.g., spawning a subagent for code exploration when a direct
grep call is faster).

**Guidance for skills that invoke subagents:**
- State when subagents are warranted: parallel independent
  tasks, isolated context, independent workstreams.
- State when they are not: simple sequential operations,
  single-file edits, tasks requiring shared state across steps.
- If the skill's workflow can be done directly, say so. Don't
  default to subagent delegation.

## Autonomy and Safety

Without guidance, Claude may take actions that are hard to
reverse or affect shared systems. Skills should match autonomy
to reversibility:

- **Local, reversible actions** (editing files, running tests):
  the skill can encourage autonomous execution.
- **Hard-to-reverse or shared-state actions** (force push,
  deleting branches, posting to external services): the skill
  should instruct the agent to confirm with the user first.

**Pattern:** check if a skill encourages autonomous destructive
actions without a confirmation step. That's a safety gap.

## Overtriggering and Verbosity

Claude 4.5/4.6 models are less verbose and more direct. Skills
written for earlier models may have compensating instructions
that are now counterproductive:

- **Excessive emphasis**: CRITICAL, MUST, NEVER, ALWAYS in caps
  — earlier models needed strong signals; current models may
  overtrigger or treat these as higher priority than intended.
- **Redundant capability reminders**: "You are an expert at X"
  or "You have the ability to Y" — the model already knows its
  capabilities.
- **Verbose output templates**: asking for detailed summaries
  after every action — current models skip unnecessary summaries
  by default, which is usually better.

**Calibration test:** read the skill's instructions and ask:
"Would a senior colleague need this much emphasis to follow
these instructions?" If not, the emphasis is calibrated for
a less capable model.

## Overengineering

Claude Opus 4.5/4.6 tend to overengineer: creating extra files,
adding unnecessary abstractions, or building in flexibility that
wasn't requested.

Skills should:
- Scope actions to what's requested — a bug fix skill shouldn't
  also clean up surrounding code.
- Avoid encouraging "while you're in there" improvements.
- State the minimum viable outcome, not the maximum possible.

## Long-Horizon and State Management

For skills that span multiple steps or potentially multiple
context windows:

- **Checkpoint progress**: encourage saving state at natural
  breakpoints so work isn't lost if context refreshes.
- **Use structured formats for state**: JSON or markdown
  checklists for tracking progress.
- **Use git for persistence**: commits provide both state
  tracking and rollback capability.
- **Incremental progress over big-bang**: "complete and verify
  each step before moving on" beats "implement everything then
  test."

## Hallucination Prevention

Claude's latest models are less prone to hallucination, but
skills can still encourage or prevent it:

- **Investigate before answering**: if a skill references files
  or code, instruct the agent to read them before making claims.
- **Ground assertions in evidence**: "ran the tests and they
  pass" (with actual output) beats "tests should pass."
- **Don't reference phantom files**: every file path mentioned
  in a skill must exist. Broken references are a form of
  hallucination in the skill itself.
