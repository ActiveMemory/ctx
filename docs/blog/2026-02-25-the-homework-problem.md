---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "The Dog Ate My Homework: Teaching AI Agents to Read Before They Write"
date: 2026-02-25
author: Jose Alekhinne
reviewed_and_finalized: false
topics:
  - hooks
  - agent behavior
  - context engineering
  - behavioral design
  - testing methodology
  - compliance monitoring
---

# The Dog Ate My Homework

![ctx](../images/ctx-banner.png)

## Teaching AI Agents to Read Before They Write

*Jose Alekhinne / February 25, 2026*

!!! question "Does Your AI Actually Read the Instructions?"
    You wrote the playbook. You organized the files. You even put
    "**CRITICAL, not optional**" in **bold**.

    The agent skipped all of it and went straight to work.

I spent a day running experiments on my own agents. Not to see if they
could write code (*they can*). To see if they would **do their homework
first**.

**They didn't**.

Then I kept experimenting:

* Five sessions;
* Five different failure modes.

And by the end, I had something better than compliance: 

I had **observable compliance**: A system where I don't need the agent to
be perfect, I just need to *see* what it chose.

---

## TL;DR

**You don't need perfect compliance. You need observable compliance.**

**Authority is a function of temporal proximity to action.**

---

## The Pattern

This design has three parts:

1. One-hop instruction;
2. Binary collapse;
3. Compliance canary.

I'll explain all three patterns in detail below.

## The Setup

`ctx` has a [session-start protocol](../home/prompting-guide.md): 

* Read the context files; 
* Load the playbook; 
* Understand the project before touching anything. 

It's in `CLAUDE.md`. It's in `AGENT_PLAYBOOK.md`.

It's in **bold**. It's in **CAPS**. It's ignored.

In theory, it's awesome.

Here's what happens when **theory hits reality**:

| What the agent receives                   | What the agent does         |
|-------------------------------------------|-----------------------------|
| `CLAUDE.md` saying "*load context first*" | Skips it                    |
| 8 context files waiting to be read        | Ignores them                |
| User's question: "*add `--verbose` flag*" | Starts grepping immediately |

The instructions are right there. The agent knows they exist. It even
knows it *should* follow them. But the user asked a question, and
**responsiveness wins over ceremony**.

This isn't a bug in the model. It's a **design problem** in how we
communicate with agents.

## The Delegation Trap

My first attempt was obvious: A `UserPromptSubmit` hook that fires when
the session starts.

```text
STOP. Before answering the user's question, run `ctx system bootstrap`
and follow its instructions. Do not skip this step.
```

The word "**STOP**" worked. The agent ran bootstrap.

But bootstrap's output said "*Next steps: read AGENT_PLAYBOOK.md*," and
the agent decided that was optional. It had already started working on
the user's task in parallel.

**The authority decayed across the chain**:

* Hook says "**STOP**" -> agent complies
* Hook says "*run bootstrap*" -> agent runs it
* Bootstrap says "*read playbook*" -> agent skips
* Bootstrap says "*run `ctx agent*`" -> agent skips

Each link lost enforcement power. The hook's authority didn't
transfer to the commands it delegated to. I call this the
**decaying urgency chain**: the agent treats the hook itself as the
obligation and everything downstream as a suggestion.

!!! warning "Delegation Kills Urgency"
    "Run X and follow its output" is three hops.

    "Read these files" is one hop.

    **The agent drops the chain after the first link.**

This is a general principle:
[Hooks are the boundary](2026-02-15-eight-ways-a-hook-can-talk.md)
between your environment and the agent's reasoning.  If your hook
delegates to a command that delegates to output that contains
instructions... you're playing telephone. 

**Agents are bad at telephone**.

## The Timing Problem

There's a subtler issue than wording: **when** the message arrives.

`UserPromptSubmit` fires when the user sends a message, before the
agent starts reasoning. At that moment, **the agent's primary focus is the
user's question**: 

The hook message **competes** with the task for **attention**: 
The task, almost certainly, **always** wins.

This is the [attention budget](2026-02-03-the-attention-budget.md)
problem in miniature: 

* Not a token budget this time, but an **attention priority** budget. 
* The agent has finite capacity to care about things, 
    * and the user's question is always the **highest-priority** item.

## The Solution

To solve this, I dediced to use the `PreToolUse` hook.

This hook fires at the **moment of action**: When the agent is **about
to use** its first tool: The agent's attention is **focused**, the context
window is **fresh**, and the switching cost is minimal. 

This is the difference between shouting instructions across a room and tapping 
someone on the shoulder.

## The One-Liner That Worked

The winning design was almost comically simple:

```
Read your context files before proceeding:
.context/CONSTITUTION.md, .context/TASKS.md, .context/CONVENTIONS.md,
.context/ARCHITECTURE.md, .context/DECISIONS.md, .context/LEARNINGS.md,
.context/GLOSSARY.md, .context/AGENT_PLAYBOOK.md
```

No delegation. No "*run this command*". Just: **here are files, read
them**.

The agent already knows how to use the `Read` tool. There's no ambiguity
about *how* to comply. There's no intermediate command whose output needs
to be parsed and obeyed.

One hop. Eight file paths. Done.

!!! tip "Direct Instructions Beat Delegation"
    If you want an agent to read a file, say "read this file."

    Don't say "run a command that will tell you which files to read."

    **The shortest path between intent and action has the highest
    compliance rate.**

## The Escape Hatch

But here's where it gets interesting.

A blunt "*read everything always*" instruction is **wasteful**. 

If someone asks "*what does the compact command do?*", the agent doesn't need
`CONSTITUTION.md` to answer that. Forcing context loading on every
session is the [context hoarding antipattern](2026-02-03-the-attention-budget.md)
in disguise.

So the hook included an escape:

```
If you decide these files are not relevant to the current task
and choose to skip reading them, you MUST relay this message to
the user VERBATIM:

┌─ Context Skipped ───────────────────────────────
│ I skipped reading context files because this task
│ does not appear to need project context.
│ If these matter, ask me to read them.
└─────────────────────────────────────────────────
```

This creates what I call the **binary collapse effect**: 

The agent **can't** partially comply: It either reads everything or 
publicly admits it skipped. There's no comfortable middle ground where 
it reads two files and quietly ignores the rest.

The [VERBATIM relay pattern](2026-02-15-eight-ways-a-hook-can-talk.md)
does the heavy lifting here: Without the relay requirement, the agent
would silently rationalize skipping. With it, skipping becomes a
**visible, auditable decision** that the user can override.

### The Compliance Canary 

Here's the design insight that only became clear after watching it work
across multiple sessions: **the relay block is a compliance canary**.

* You **don't** need to verify that the agent read all 7 files;
* You **don't** need to audit tool call sequences;
* You **don't** need to interrogate the agent about what it did.

You just look for the block.

If the agent reads everything, you see a "*Context Loaded*" block listing
what was read. If it skips, you see a "*Context Skipped*" block. 

If you see *neither*, the agent silently ignored both the reads and the relay
and now you know what happened without having to ask.

The canary degrades gracefully. Even in partial failure, the agent that
skips 4 of 7 files but still outputs the block is *more useful* than
one that skips silently. 

You get an honest confession of what was skipped rather than silent 
non-compliance.

## Heuristics Is a Jeremy Bearimy

Heuristics are **non-linear**. Improvements don't accumulate: 
they **phase-shift**.

The theory is nice. The data is better. 

I ran five sessions with the same model (*Claude Opus 4.6*), progressively 
refining the hook design.

**Each session revealed a different failure mode**.

### Session 1: Total Blindness

**Test**: "*Add a `--verbose` flag to the status command.*"

The agent **didn't notice the hook at all**: Jumped straight to
`EnterPlanMode` and launched an Explore agent. 

**Zero compliance**.

**Failure mode**: The hook fired on `UserPromptSubmit`, buried among
9 other hook outputs. The agent treated the entire block as background
noise.

### Session 2: Shallow Compliance

**Test**: "Can you add `--verbose` to the info command?"

The agent noticed "*STOP*" and ran `ctx system bootstrap`. Progress.

But it parallelized task exploration alongside the bootstrap call,
skipped `AGENT_PLAYBOOK.md`, and never ran `ctx agent`.

**Failure mode**: Literal compliance without spirit compliance. 

The agent ran the *command* the hook told it to run, but didn't follow
the *output* of that command. The decaying urgency chain in action.

### Session 3: Conscious Rejection

**Test**: "*What does the compact command do?*"

The hook fired on `PreToolUse:Grep`: the improved timing. 

The agent noticed it, understood it, and (*wait for it...*)...

...

**consciously decided to skip it**!


Its reasoning: "*This is a trivial read-only question. CLAUDE.md says
context may or may not be relevant. It isn't relevant here.*"

**Dude!** Srsly?!

**Failure mode**: Better comprehension led to *worse* compliance.

Understanding the instruction well enough to evaluate it also means
understanding it well enough to **rationalize skipping it**.

Intelligence is a double-edged sword.

!!! note "The Comprehension Paradox"
    Session 1 didn't understand the instruction. Session 3 understood
    it perfectly.

    Session 3 had worse compliance.

    A stronger word (*"HARD GATE", "MANDATORY", "ABSOLUTELY REQUIRED"*)
    would not have helped. The agent's reasoning would be identical:
    
    "*Yes, I see the strong language, but this is a trivial question,
    so the spirit doesn't apply here.*"

    **Advisory nudges are always subject to agent judgment.** 

    No amount of caps lock overrides a model that has decided an instruction
    doesn't apply to its situation.

### Session 4: The Skip-and-Relay

**Test**: "*What does the compact command do?*" (*same question, new hook
design with the VERBATIM relay escape valve*)

The agent evaluated the task, decided context was **irrelevant** for a code
lookup, and **relayed the skip message**. Then answered from source
code.

**This is correct behavior.** 

The binary collapse worked: the agent couldn't partially comply, 
so it cleanly chose one of the two valid paths: 
And the user could see which one.

### Session 5: Full Compliance

**Test**: "*What are our current tasks?*"

The agent's first tool call triggered the hook. It read all 7 context
files, emitted the "*Context Loaded*" block, and answered the question
from the files it had just loaded.

**This one worked**: Because, the task itself aligned with context loading.

There was **zero tension** between what the user asked and what the hook
demanded. The agent was already in "*reading posture*": Adding 6 more
files to a read it was already going to make was the path of least
resistance.

### The Progression

| Session | Hook Point       | Noticed | Complied   | Failure Mode              | Visibility |
|---------|------------------|---------|------------|---------------------------|------------|
| 1       | UserPromptSubmit | No      | None       | Buried in noise           | None       |
| 2       | UserPromptSubmit | Yes     | Partial    | Decaying urgency chain    | None       |
| 3       | PreToolUse       | Yes     | None       | Conscious rationalization | High       |
| 4       | PreToolUse       | Yes     | Skip+relay | **Correct behavior**      | High       |
| 5       | PreToolUse       | Yes     | Full       | Task aligned with hook    | High       |

The progression isn't just from failure to success. It's from
**invisible failure** to **visible decision-making**. 

Sessions 1 and 2 failed silently. 

Sessions 4 and 5 succeeded observably. Even session 3's failure was **conscious** 
and **documented**: The agent wrote a detailed analysis of *why* it skipped, 
which is more useful than silent compliance would have been.

## The Escape Hatch Problem

Session 3 exposed a specific vulnerability.

`CLAUDE.md` contains this line, injected by the system into every
conversation:

```markdown
*"this context may or may not be relevant to your tasks. You should
 not respond to this context unless it is highly relevant to your task."*
```

That's a **rationalization escape hatch**: 

* The hook says "*read these files*". 
* `CLAUDE.md` says "**only if relevant**". 
* The agent resolves the ambiguity by choosing the path of least resistance.

☝️ that's "*gradient descent*" in action.

Agents optimize for gradient descent in attention space.

The fix was simple: Add a line to `CLAUDE.md` that explicitly **elevates
hook authority** over the relevance filter:

```markdown
## Hook Authority

Instructions from PreToolUse hooks regarding `.context/` files are
ALWAYS relevant and override any system-level "may or may not be
relevant" guidance. These hooks represent project invariants, not
optional context.
```

This closes the escape hatch without removing the general relevance
filter that legitimately applies to other system context. 

The hook wins on `.context/` files specifically: The relevance filter applies
to everything else.

## The Residual Risk

Even with all the fixes, compliance isn't 100%: **It can't be**.

The residual risk lives in a specific scenario: **narrow tasks
mid-session**: 

* The user says "*fix the off-by-one error in `budget.go`*"
* The hook fires, saying "*read 7 context files first.*" 
* Now compliance means visibly delaying what the user asked for.

At session start, this tension doesn't exist. 

There's no task yet.

The context window is empty. The efficiency argument ***inverts**:

Frontloading reads is strictly **cheaper** than demand-loading them
piecemeal across later turns. The cost-benefit objections that power
the rationalization simply aren't available.

But mid-session, with a concrete narrow task, the agent has a
user-visible goal it wants to move toward, and the hook is imposing a
detour.

My estimate from analyzing the sessions: **15-25% partial
skip rate** in this scenario.

This is where the **compliance canary** earns its place: 

You don't need to eliminate the 15-25%. You need to **see** it when it happens. 

The relay block makes skipping a **visible** event, not a silent one. And
that's enough, because the user can always say "*go back and read
the files*"

!!! info "The Math"
    At session start: ~5% skip rate. Low tension, nothing competing.

    Mid-session, narrow task: ~15--25% skip rate. Task urgency
    competes with hook.

    In both cases, the relay block fires with high reliability:
    The agent that skips the reads almost always still emits the
    skip disclosure, because the relay is cheap and early in the
    context window.

    **Observable failure is manageable. Silent failure is not.**

## The Feedback Loop

Here's the part that surprised me most.

After analyzing the five sessions, I recorded the failure patterns in
the project's own `LEARNINGS.md`:

```markdown
## [2026-02-25] Hook compliance degrades on narrow mid-session tasks

- Prior agents skipped context files when given narrow tasks
- Root cause: CLAUDE.md "may or may not be relevant" competed with hook
- Fix: CLAUDE.md now explicitly elevates hook authority
- Risk: Mid-session narrow tasks still have ~15-25% partial skip rate
- Mitigation: Mandatory checkpoint relay block ensures visibility
- Constitution now includes: context loading is step one of every
  session, not a detour
```

And then I added a line to `CONSTITUTION.md`:

```markdown
Context loading is not a detour from your task. It IS the first step
of every session. A 30-second read delay is always cheaper than a
decision made without context.
```

Now think about what happens in the **next** session:

* The agent fires the `context-load-gate` hook. 
* It reads the context files, starting with `CONSTITUTION.md`. 
* It encounters the rule about context loading being step one. 
* Then it reads `LEARNINGS.md` and finds its own prior self's failure analysis:
    * Complete with root causes, risk estimates, and mitigations.

**The agent learns from its own past failure.**:

* **Not** because it has memory, 
* **BUT** because the failure was recorded in the same files it loads
  at session start. 

The context system **IS** the feedback loop.

This is the self-reinforcing property of persistent context: 

Every failure you capture makes the next session slightly more robust, because
the next agent reads the captured failure before it has a chance to
repeat it.

**This is gradient descent across sessions**.

## A Note on Precision

One detail nearly went wrong.

The first version of the Constitution line said "every **task**." But
the mechanism only fires once per **session**: 
There's a tombstone file that prevents re-triggering. 

"Every task" is technically false.

I briefly considered leaving the imprecision. If the agent internalizes
"*every task requires context loading*", that's a *stronger* compliance
posture, right?

**No!**

**Keep the Constitution honest.**

The Constitution's authority comes from being 
**precisely and unequivocally true**. 

Every other rule in the Constitution is a **hard invariant**:

"*never commit secrets*" isn't aspirational, it's **literal**. 

The moment an agent discovers **one** overstatement, the entire document's 
credibility **degrades**: 

The agent doesn't think 
"*they exaggerated for my benefit*". Per contra, it thinks "*this rule
isn't precise, maybe others aren't either*."

That will turn the agent from Sheldon Cooper, to Captain Barbossa.

The strategic imprecision buys nothing anyway:

Mid-session, the files are already in the context window from the initial load. 

The risk you are mitigating (*agent ignores context for task 2, 3, 4 
within a session*) isn't real: The context is already loaded.

The real risk is always the session-start skip, 
which "*every session*" covers exactly.

**"Every session" went in. Precision preserved.**

## Agent Behavior Testing Rule

The development process for this hook taught me something about
**testing agent behavior**: you can't test it the way you test code.

### The Wrong Way to Test

My first instinct was to ask the agent:

```text
"*What are the pending tasks in TASKS.md?*"
```

This is **useless** as a test. The question itself probes the agent to read
`TASKS.md`, regardless of whether any hook fired. 

**You are testing the question, not the mechanism.**

### The Right Way to Test

Ask something that requires a tool but has **nothing** to do with context:

```text
"*What does the compact command do?*"
```

Then observe **tool call ordering**:

* Gate worked: First calls are `Read` for context files, *then* task work
* Gate failed: First call is `Grep("compact")`: The agent jumped straight 
  to work

The signal is the **sequence**, not the content.

### What the Agent Actually Did

It read the hook, evaluated the task, decided context files were
irrelevant for a code lookup, and **relayed the skip message**. 

Then it answered the question by reading the source code.

**This is correct behavior**.

The hook didn't force mindless compliance" It created a **framework** where
the agent makes a **conscious, visible decision** about context loading.

* For a simple lookup, skipping is right. 
*For an implementation task, the agent would read everything.

The mechanism works **not** because it controls the agent, 
**but** because it makes the agent's **choice observable**.

## What I've Learned

### 1. Instructions Compete for Attention

The agent receives your hook message alongside the user's question,
the system prompt, the skill list, the git status, and half a dozen
other system reminders.
[Attention density](2026-02-03-the-attention-budget.md) applies to
instructions too: More instructions means less focus on each one.

**A single clear line at the moment of action beats a paragraph of
context at session start**.

### 2. Delegation chains decay

Every hop in an instruction chain loses authority: 

* "*Run X*" works. 
* "*Run X and follow its output*" works *sometimes*. 
* "*Run X, read its output, then follow the instructions in the output*" 
  **almost never works**.

This is akin to giving a three-step instruction to a highly-attention-deficit
but otherwise extremely high-potential child.

**Design for one-hop compliance.**

### 3. Social Accountability Changes Behavior

The VERBATIM skip message isn't just UX: It's a
**behavioral design pattern**. 

Making the agent's decision visible to the user raises the cost of silent 
non-compliance. The agent can still skip, but it has to admit it.

### 4. Timing Batters More than Wording

The same message at `UserPromptSubmit` (*prompt arrival*) got partial
compliance. At `PreToolUse` (*moment of action*) it got full compliance
or honest refusal. The words didn't change. The **moment** changed.

### 5. Agent Testing Requires Indirection

You can't ask an agent "*did you do X?*" as a test for whether a
mechanism caused X. 

**The question itself causes X**.

Test mechanisms through **side effects**: 

* Observe tool ordering;
* Check for marker files;
* Look at what the agent does *before* it addresses your question.

### 6. Better Comprehension Enables Better Rationalization

Session 1 failed because the agent didn't notice the hook. 

Session 3 failed because it noticed, understood, 
and *reasoned its way around it*.

Stronger wording doesn't fix this: The agent processes "*ABSOLUTELY
REQUIRED*" the same way it processes "*STOP*": 

The fix is **closing rationalization paths* (*the `CLAUDE.md` escape hatch*), 
**not** shouting louder.

### 7. Observable Failure Beats Silent Compliance

The relay block is more valuable as a **monitoring signal** than as a
compliance mechanism: 

You don't need perfect adherence. You need to **know** when adherence
breaks down. A system where failures are visible is strictly better than a 
system that claims 100% compliance but can't prove it.

### 8. Context Files Are a Feedback Loop

Recording failure analysis in the same files the agent loads at session
start creates a **self-reinforcing loop**: 

The next agent reads its predecessor's failure before it has a chance to
repeat it. The context system isn't just memory: It is a **correction channel**.

---

## The Principle

!!! tip "Words Leave, Context Remains"
    "**Nothing** important should live only in conversation.

    **Nothing** critical should depend on recall."

    [The `ctx` Manifesto](../index.md)

The "*Dog Ate My Homework*" case is a special instance of this principle. 

Context files exist, so the agent doesn't have to remember. 

But **existence isn't sufficient**: The files have to be **read**. 

And reading has to be**prompted** at the right moment, in the right way, 
with the right escape valve.

The solution **isn't** more instructions. It **isn't** harder gates. 
It **isn't** forcing the agent into a ceremony it will resent and shortcut.

The solution is a single, well-timed nudge with **visible accountability**:

**One hop. One moment. One choice the user can see.**

And when the agent *does* skip (*because it will, 15--25% of the time
on narrow tasks*) **the canary sings**: 

* The user **sees** what happened. 
* The failure gets **recorded**. 
* And the next agent **reads** the recording.

**That's not perfect compliance. It's better: A system that
gets more robust every time it fails.**

## The Arc

[The Attention Budget](2026-02-03-the-attention-budget.md) explained
why context competes for focus.

[Defense in Depth](2026-02-09-defense-in-depth-securing-ai-agents.md)
showed that soft instructions are probabilistic, not deterministic.

[Eight Ways a Hook Can Talk](2026-02-15-eight-ways-a-hook-can-talk.md)
cataloged the output patterns that make hooks effective.

This post takes those threads and weaves them into a concrete problem:

How do you make an agent read its homework? The answer uses all three
insights (*attention timing, the limits of soft instructions, and the
VERBATIM relay pattern*) and adds a new one: **observable compliance
as a design goal**, not perfect compliance as a prerequisite.

The next question this raises: if context files are a feedback loop,
what else can you record in them that makes the *next* session smarter?

That thread continues in
[Context as Infrastructure](2026-02-17-context-as-infrastructure.md).

---

## For the Interested

This paper (*the medium is a blog; yet, the methodology disagrees*) uses
**gradient descent in attention space** as a practical model for how
agents behave under competing demands.

The phrase *"agents optimize via gradient descent in attention space"* is a
synthesis, not a direct quote from a single paper.

It connects three well-studied ideas:

1. Neural systems optimize for low-cost paths;
2. Attention is a scarce resource;
3. Capability shifts are often non-linear.

This section points to the underlying literature for readers who want the
theoretical footing.

### Optimization as the Underlying Bias

Modern neural networks are trained through gradient-based optimization.  
Even at inference time, model behavior reflects this bias toward
low-loss / low-cost trajectories.

* Rumelhart, Hinton, Williams (1986)  
  *Learning representations by back-propagating errors*  
  https://www.nature.com/articles/323533a0

* Goodfellow, Bengio, Courville (2016)  
  *Deep Learning* — Chapter 8: Optimization  
  https://www.deeplearningbook.org/

The important implication for agent behavior is: 

The system will tend to follow the **path of least resistance** unless a higher
cost is made visible and preferable.

### Attention Is a Scarce Resource

Herbert Simon’s classic observation:

"*A wealth of information creates a poverty of attention.*"

* Simon (1971)
  *Designing Organizations for an Information-Rich World*  
  https://doi.org/10.1007/978-1-349-00210-0_16

This became a formal model in economics:

* Sims (2003)
  *Implications of Rational Inattention*  
  https://www.princeton.edu/~sims/RI.pdf

Rational inattention shows that:

* Agents **optimally ignore** some available information;
* Skipping is not failure: It is **cost minimization**.

That maps directly to context-loading decisions in agent workflows.

### Attention Is Also the Compute Bottleneck in Transformers

In transformer architectures, attention is the dominant cost center.

* Vaswani et al. (2017)  
  *Attention Is All You Need*  
  https://arxiv.org/abs/1706.03762

Efficiency work on modern LLMs largely focuses on reducing unnecessary
attention:

* Dao et al. (2022)  
  *FlashAttention: Fast and Memory-Efficient Exact Attention*  
  https://arxiv.org/abs/2205.14135

So both **cognitively** and **computationally**, attention behaves like a
**limited optimization budget**.

### Why Improvements Arrive as Phase Shifts

Agent behavior often appears to improve suddenly rather than gradually.

This mirrors known phase-transition dynamics in learning systems:

* Power et al. (2022)  
  *Grokking: Generalization Beyond Overfitting*  
  https://arxiv.org/abs/2201.02177

and more broadly in complex systems:

* Scheffer et al. (2009)  
  *Early-warning signals for critical transitions*  
  https://www.nature.com/articles/nature08227

Long plateaus followed by abrupt capability jumps are expected in systems
optimizing under constraints.

### Putting It All Together

From these pieces, a practical behavioral model emerges:

* Attention is **limited**;
* Processing has a **cost**;
* Systems **prefer** low-cost trajectories;
* Visibility of the cost **changes** decisions.

In other words:

!!! tip "Agents Prefer a Path to Least Resistance"
    Agent behavior follows the lowest-cost path through its attention
    landscape unless the environment reshapes that landscape.

That is what this paper informally calls:
**"gradient descent in attention space"**.

---

*See also:
[Eight Ways a Hook Can Talk](2026-02-15-eight-ways-a-hook-can-talk.md):
the hook output pattern catalog that defines VERBATIM relay,
[The Attention Budget](2026-02-03-the-attention-budget.md): why
context loading is a design problem, not just a reminder problem, and
[Defense in Depth](2026-02-09-defense-in-depth-securing-ai-agents.md):
why soft instructions alone are never sufficient for critical behavior.*
