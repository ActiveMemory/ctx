---
title: "The Anatomy of a Skill That Works"
date: 2026-02-07
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - skill design
  - context engineering
  - quality gates
  - E/A/R framework
  - practical patterns
---

# The Anatomy of a Skill That Works

!!! note "Update (2026-02-11)"
    As of `v0.4.0`, `ctx` consolidated sessions into the journal mechanism.
    References to `ctx-save`, `ctx session`, and `.context/sessions/`
    in this post reflect the architecture at the time of writing.

![ctx](../images/ctx-banner.png)

## What 20 Skill Rewrites Taught Me About Guiding AI

*Jose Alekhinne / 2026-02-07*

!!! question "Why do some skills produce great results while others get ignored or produce garbage?"
    I had 20 skills. Most were well-intentioned stubs: a description,
    a command to run, and a wish for the best.

    Then I **rewrote all of them in a single session**. This is what I learned.

In [Skills That Fight the Platform][fight-post], I described what skills
should *not* do. In [You Can't Import Expertise][import-post], I showed
why templates fail. This post completes the trilogy: the concrete
patterns that make a skill actually work.

[fight-post]: 2026-02-04-skills-that-fight-the-platform.md
[import-post]: 2026-02-05-you-cant-import-expertise.md

## The Starting Point

Here is what a typical skill looked like before the rewrite:

```markdown
---
name: ctx-save
description: "Save session snapshot."
---

Save the current context state to `.context/sessions/`.

## Execution

ctx session save $ARGUMENTS

Report the saved session file path to the user.
```

Seven lines of body. A vague description. No guidance on *when* to use
it, *when not to*, what the command actually accepts, or how to tell
if it worked.

As a result, the agent would either never trigger the skill (*the
description was too vague*), or trigger it and produce shallow output
(*no examples to calibrate quality*).

**A skill without boundaries is just a suggestion.**

More precisely: the most effective boundary I found was a **quality gate
that runs *before* execution**, not during it.

## The Pattern That Emerged

After rewriting 20 skills, a **repeatable anatomy** emerged (*independent
of the skill's purpose*). Not every skill needs every section, but the
effective ones share the same bones:

| Section           | What It Does                                    |
|-------------------|-------------------------------------------------|
| Before X-ing      | Pre-flight checks; prevents premature execution |
| When to Use       | Positive triggers; narrows activation           |
| When NOT to Use   | Negative triggers; prevents misuse              |
| Usage Examples    | Invocation patterns the agent can pattern-match |
| Process/Execution | What to do; commands, steps, flags              |
| Good/Bad Examples | Desired vs undesired output; sets boundaries    |
| Quality Checklist | Verify before claiming completion               |

I realized **the first three sections matter more than the rest**;
because **a skill with great execution steps but no activation guidance
is like a manual for a tool nobody knows they have**.

!!! warning "Anti-Pattern: The Perfect Execution Trap"
    A skill with detailed execution steps but no activation guidance
    will *fail* **more often** than a vague skill because it executes
    confidently at the wrong time.

## Lesson 1: Quality Gates Prevent Premature Execution

The single most impactful addition was a "*Before X-ing*" section at the
top of each skill. Not process steps; pre-flight checks.

```markdown
## Before Recording

1. **Check if it belongs here**: is this learning specific
   to this project, or general knowledge?
2. **Check for duplicates**: search LEARNINGS.md for similar
   entries
3. **Gather the details**: identify context, lesson, and
   application before recording
```

* Without this gate, the agent would execute immediately on trigger.
* With it, the agent pauses to verify preconditions.

The difference is dramatic: instead of shallow, reflexive execution, you get
considered output.

!!! tip "Readback"
    For the astute readers, the aviation parallel is intentional:

    Pilots do not skip the pre-flight checklist because they have flown before.

    The checklist exists precisely because the stakes are high enough that
    "*I know what I'm doing*" is not sufficient.

## Lesson 2: "*When NOT to Use*" Is **Not** Optional

Every skill had a "*When to Use*" section. Almost none had "*When NOT
to Use*". This is a problem.

AI agents are biased toward action. Given a skill that says "*use
when journal entries need enrichment*", the agent will find reasons
to enrich.

Without explicit negative triggers, **over-activation is not a bug; it
is the default behavior**.

Some examples of negative triggers that made a real difference:

| Skill        | Negative Trigger                                           |
|--------------|------------------------------------------------------------|
| ctx-reflect  | "*When the user is in flow; do not interrupt*"             |
| ctx-save     | "*After trivial changes; a typo does not need a snapshot*" |
| prompt-audit | "*Unsolicited; only when the user invokes it*"             |
| qa           | "*Mid-development when code is intentionally incomplete*"  |

These are not just nice-to-have. They are **load-bearing**. 

Withoutthem, the agent will trigger the skill at the *wrong* time, produce
*unwanted* output, and **erode the user's trust** in the skill system.

## Lesson 3: Examples Set Boundaries Better Than Rules

The most common failure mode of thin skills was not wrong behavior
but *vague* behavior. The agent would do roughly the right thing,
but at a quality level that required human cleanup.

Rules like "*be constructive, not critical*" are too abstract. What
does "*constructive*" look like in a prompt audit report? The agent
has to guess.

Good/bad example pairs avoid guessing:

```markdown
### Good Example

> This session implemented the cooldown mechanism for
> `ctx agent`. We discovered that `$PPID` in hook context
> resolves to the Claude Code PID.
>
> I'd suggest persisting:
> - **Learning**: `$PPID` resolves to Claude Code PID
>   `ctx add learning --context "..." --lesson "..."`
> - **Task**: mark "Add cooldown" as done

### Bad Examples

* "*We did some stuff. Want me to save it?*"
* Listing 10 trivial learnings that are general knowledge
* Persisting without asking the user first
```

The good example shows the exact format, level of detail, and
command syntax. The bad examples show where the boundary is.

Together, they define a quality corridor without prescribing
every word.

**Rules describe. Examples demonstrate.**

## Lesson 4: Skills Are Read by Agents, Not Humans

This seems obvious, but it has **non-obvious** consequences. During
the rewrite, one skill included guidance that said "*use a blog or
notes app*" for general knowledge that does not belong in the
project's learnings file.

The agent does not have a notes app. It does not browse the web
to find one. This instruction, clearly written for a human
audience, was *dead weight* in a skill consumed by an AI.

!!! tip "Skills are for the Agents"
    **Every sentence in a skill should be actionable by the agent**.

    If the guidance requires human judgment or human tools, it belongs in
    documentation, **not** in a skill.

    The corollary: **command references must be exact**. 

    A skill that says "*save it somewhere*" is useless. 

    A skill that says `ctx add learning --context "..." --lesson "..." --application "..."`
    is **actionable**.

    The agent can pattern-match and fill in the blanks.

**Litmus test**: If a sentence starts with "*you could...*" or assumes
external tools, it does not belong in a skill.

## Lesson 5: The Description Field Is the Trigger

This was covered in [Skills That Fight the Platform][fight-post],
but the rewrite reinforced it with data. Several skills had good
bodies but vague descriptions:

```yaml
# Before: vague, activates too broadly or not at all
description: "Show context summary."

# After: specific, activates at the right time
description: "Show context summary. Use at session start or
  when unclear about current project state."
```

The description is not a title. It is the **activation condition**.

The platform's skill matching reads this field to decide whether
to surface the skill. A vague description means the skill either
*never* triggers or triggers when it **should not**.

## Lesson 6: Flag Tables Beat Prose

Most skills wrap CLI tools. The thin versions described flags in
prose, if at all. The rewritten versions use tables:

```markdown
| Flag        | Short | Default | Purpose                  |
|-------------|-------|---------|--------------------------|
| `--limit`   | `-n`  | 20      | Maximum sessions to show |
| `--project` | `-p`  | ""      | Filter by project name   |
| `--full`    |       | false   | Show complete content    |
```

Tables are **scannable**, **complete**, and **unambiguous**. 

The agent can read them faster than parsing prose, and they serve as both
reference and validation: If the agent invokes a flag not in the table, 
something is wrong.

## Lesson 7: Template Drift Is a Real Maintenance Burden

// TODO: this has changed; we deploy from the marketplace; update it.
// at least add an admonition saying thing are different now.

`ctx` deploys skills through templates (via `ctx init`). Every
skill exists in two places: the live version (`.claude/skills/`)
and the template (`internal/assets/claude/skills/`).

They must match.

During the rewrite, every skill update required editing both files
and running `diff` to verify. This sounds trivial, but across 16
template-backed skills, it was the most error-prone part of the
process.

Template drift is dangerous because it creates **false confidence**:
the agent appears to follow rules that no longer exist.

The lesson: **if your skills have a deployment mechanism, build
the drift check into your workflow.** We added a row to the
`update-docs` skill's mapping table specifically for this:

```markdown
| `internal/assets/claude/skills/` | `.claude/skills/` (live) |
```

Intentional differences (*like project-specific scripts in the
live version but not the template*) should be documented, not
discovered later as bugs.

## The Rewrite Scorecard

| Metric                   | Before    | After     |
|--------------------------|-----------|-----------|
| Average skill body       | ~15 lines | ~80 lines |
| Skills with quality gate | 0         | 20        |
| Skills with "When NOT"   | 0         | 20        |
| Skills with examples     | 3         | 20        |
| Skills with flag tables  | 2         | 12        |
| Skills with checklist    | 0         | 20        |

More lines, but almost entirely Expert content (per the
[E/A/R framework][import-post]). No personality roleplay, no
redundant guidance, no capability lists. Just project-specific
knowledge the platform does not have.

## The Meta-Lesson

The previous two posts argued that skills should provide
knowledge, not personality; that they should complement the
platform, not fight it; that they should grow from project
history, not imported templates.

This post adds the missing piece: **structure**.

**A skill without a structure is a wish**.

A skill with quality gates, negative triggers, examples, and
checklists is a **tool**: the difference is not the content; it
is **whether the agent can reliably execute it without human
intervention**.

!!! tip "Skills are Interfaces"
    **Good skills are not instructions. They are contracts.**:

    * They **specify** preconditions, postconditions, and boundaries.
    * They **show** what success looks like and what failure looks like.
    * They **trust** the agent's intelligence but do not trust its assumptions.

---

!!! quote "**If You Remember One Thing From This Post...**"
    **Skills that work have bones, not just flesh.**

    Quality gates, negative triggers, examples, and checklists
    are the skeleton. The domain knowledge is the muscle.

    Without the skeleton, the muscle has nothing to attach to.

---

*This post was written during the same session that rewrote all
22 skills. The skill-creator skill was updated to encode these
patterns. The meta continues.*
