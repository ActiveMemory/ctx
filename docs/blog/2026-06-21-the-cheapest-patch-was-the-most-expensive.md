---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "The Cheapest Patch Was the Most Expensive:
  What Seven AI Coding Runs Taught Me About Cost"
date: 2026-06-21
author: Volkan Özçelik
reviewed_and_finalized: true
topics:
  - spec-driven development
  - model selection
  - context compression
  - agentic coding cost
  - field notes
---

# The Cheapest Patch Was the Most Expensive

![ctx](../images/ctx-banner.png)

## What Seven AI Coding Runs Taught Me About Cost

*Volkan Özçelik / June 21, 2026*

!!! question "What Does a Cheap Patch Actually Cost?"
    The cheapest run fixed the visible bug in four minutes for
    thirty-five cents, and missed the contract entirely.

    The most expensive run wrote a strong patch, and quietly changed
    a product decision nobody asked it to change.

    Neither number on the invoice told you either of those things.

I ran a small AI coding experiment on a real CLI bug.

The bug was boring, which made it useful.

A command accepted a comma-separated list of secret versions. This
worked:

```text
--versions "1,2,3"
```

This looked like it worked:

```text
--versions "1, 2, 3"
```

But one command silently sent only the first version.

The validation path handled whitespace. The conversion path did not.

For example, `"1, 2, 3"` became `[1]`.

A sibling command had similar parsing code, but not the exact same
failure. The right fix was not "*make this one line trim spaces*".
The right fix was to stop duplicating the parsing logic and send both
commands through the same parser.

It was easy to see if you knew the codebase, but deceptively complex
if you were unfamiliar with the project. For reference, the project is
[SPIKE](https://github.com/spiffe/spike).

So I thought I could run a controlled experiment on how spec-driven
development methodologies, context-compression techniques, `ctx`, and
different model choices play together.

If I were to write a paper (*and I am planning to write one*), the
thesis would read something like this:

```text
We evaluate whether context-engineering tools reduce the real cost of
agentic coding under spec-driven development. Rather than measuring token
savings alone, we measure accepted-patch cost: model cost, repair loops,
hidden acceptance failures, and human review burden. On a substantial
cross-layer task in the SPIFFE/SPIKE repository, we compare direct
issue-to-code prompting, frontier-authored SDD artifacts, weak-authored
artifacts with frontier ratification, and multiple context conditions
including structured context manifests, shell-output compression, and
context-runtime filtering. Our results show whether token savings translate
into accepted patches, and identify when context tools help, hurt, or merely
move cost into review.
```

!!! note "This Is a Weekend Hack, Not the Paper"
    To make this paper-grade, I figured I would need to run **~500**
    controlled agentic experiments, each spanning at least half an
    hour. That is not a weekend hack. So I picked a meaningful subset,
    ran them end-to-end, and that is what you are reading. 

    **Treat these numbers as signal, not as proof**.

## The Shared Parser

One more detail about the task at hand: there was already a shared
parser: An agent leveraging the shared parser would already implement
half of the solution and have a head-start. An agent that missed the parser would burn tokens rebuilding what
was already there.

And that detail changed the entire task:

* The job was **not** to design a parser.
* It was to wire an existing helper into two commands, preserve the
  product contract, and add focused command tests.

But the agents that were going to implement this did not know that a
priori.

That is where the whole experiment became interesting.

## The Task at Hand

The accepted behavior was:

```text
"1, 2, 3" -> [1, 2, 3]
```

Other decisions mattered too:

* `""` or whitespace-only selector -> `[0]`
* `0` remains the current-version sentinel
* empty inner tokens are rejected
* non-integers are rejected
* negative integers are rejected
* duplicates are preserved
* no SDK/API/backend/state changes
* no framework rewrite

The invariant was simple:

```text
Accept the whole selector,
or reject the whole selector before any API call.

Do not silently drop a token.
```

The bug was small enough that any strong model could patch something.
Yet it was large enough that a quick patch could be wrong in ways that
looked correct from the outside.

A perfect setup.

## The Runs

I ran multiple implementations of the same task.

Some runs used context compression (*reducing token count by
eliminating unnecessary content that flows to and from the model,
while keeping the compression as lossless as possible*). Some did not.

Here is an uncompressed CLI call:

```bash
volkan@sdd:~/WORKSPACE$ ls -al
total 534724
drwxrwxr-x  6 volkan volkan     4096 Jun 20 22:02 .
drwxr-x--- 21 volkan volkan     4096 Jun 21 14:04 ..
drwxr-xr-x 25 volkan volkan     4096 Jun 20 11:35 ctx
drwxrwxr-x 19 volkan volkan     4096 Jun 20 12:46 ctx-bak
drwxrwxr-x  2 volkan volkan     4096 Jun 20 19:55 harness
drwxrwxr-x 22 volkan volkan     4096 Jun 20 21:35 spike
-rw-rw-r--  1 volkan volkan 78242134 Jun 20 20:51 spike-haiku-for-spec-tooling-on-haiku-for-exec.zip
-rw-rw-r--  1 volkan volkan 78116329 Jun 20 18:00 spike-opus-4-8-xhigh-no-tooling.zip
-rw-rw-r--  1 volkan volkan 78225976 Jun 20 19:54 spike-sonnet-4-6-medium-tooling-on-haiku-for-exec.zip
-rw-rw-r--  1 volkan volkan 78176133 Jun 20 19:39 spike-sonnet-4-6-medium-tooling-on-sonnet-for-exec.zip
-rw-rw-r--  1 volkan volkan 78245337 Jun 20 21:52 spike-sonnet-medium-for-all-tooling-off.zip
-rw-rw-r--  1 volkan volkan 78282827 Jun 20 22:02 spike-sonnet-medium-for-specs-haiku-for-exec-no-tooling.zip
-rw-rw-r--  1 volkan volkan 78217444 Jun 20 20:06 spike-yolo-no-specs-haiku-for-exec.zip
```

And here is the compressed version for comparison:

```txt
755  ctx/
775  ctx-bak/
775  harness/
775  spike/
664  spike-haiku-for-spec-tooling-on-haiku-for-exec.zip  74.6M
664  spike-opus-4-8-xhigh-no-tooling.zip  74.5M
664  spike-sonnet-4-6-medium-tooling-on-haiku-for-exec.zip  74.6M
664  spike-sonnet-4-6-medium-tooling-on-sonnet-for-exec.zip  74.6M
664  spike-sonnet-medium-for-all-tooling-off.zip  74.6M
664  spike-sonnet-medium-for-specs-haiku-for-exec-no-tooling.zip  74.7M
664  spike-yolo-no-specs-haiku-for-exec.zip  74.6M

Summary: 7 files, 4 dirs (7 .zip)
```

The goal of the compression was to preserve meaningful content while
cutting the fluff that would not typically benefit the agent.

In this experiment:

```text
compression ON  = both context compression layers enabled
compression OFF = both context compression layers disabled
```

Except for one "*YOLO this thing end to end*" negative-control case,
every serious run went through a structured debrief/spec/task
workflow, following a formal spec-driven-development methodology.

The decisions the agent made were not necessarily caused by
information loss during compression. They were more about the
**quality** and the **shape** of the context available while the plan
hardened. Which also meant the quality of the agent (**and the
human**) mattered **a lot** during the planning and spec-development
phase.

Here is the short summary of the experiments. For simplicity, and to
keep this a weekend hack, I only used Anthropic models.

| Run                                                         | Planning / discovery                                                             | Implementation                                      | Quality / caveat                                                                                                           |
|-------------------------------------------------------------|----------------------------------------------------------------------------------|-----------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------|
| **Opus end-to-end, compression OFF**<br>$16.25 · 59m27s     | Handled the debrief/spec/task work and implementation.                           | Opus implemented.                                   | Strong patch. It also tightened behavior beyond the final compatibility decision by rejecting whitespace-padded selectors. |
| **Sonnet, compression ON**<br>$6.43 · 1h04m15s              | Completed the structured workflow, but planned larger parser work.               | Sonnet implemented.                                 | Acceptable, but larger than necessary.                                                                                     |
| **Sonnet, compression OFF**<br>$6.11 · 50m13s               | Found that `parseVersionList()` already existed and narrowed the task to wiring. | Sonnet implemented.                                 | Preferred patch shape.                                                                                                     |
| **Haiku, compression ON**<br>$2.15 · 39m53s                 | Completed the structured workflow after steering.                                | Haiku implemented.                                  | Cheap, but needed steering. Weak at repo discovery.                                                                        |
| **Sonnet OFF plan, Haiku implementation**<br>~$4.92 · ~50m15s | Sonnet (compression OFF) found and specified the smaller wiring task.                              | Haiku implemented the ratified task list. | Worth repeating as a follow-up experiment. Not a default rule.                                                             |
| **Sonnet ON plan, Haiku implementation**<br>~$4.9 · composite | Sonnet (compression ON) planned the larger parser task. | Haiku implemented the ratified task list. | Acceptable, but inherited the larger premise.                                                                              |
| **Haiku YOLO**<br>$0.35 attempt · 4m23s                     | No structured workflow.                                                          | Haiku implemented directly.                         | **Rejected**: Fixed the visible symptom, but missed the accepted contract.                                                 |

The two Sonnet end-to-end rows above deserve closer attention.

Same model family. Same general workflow. Different compression
setting.

* With compression OFF, the model found the existing helper and
  narrowed the work.
* With compression ON, the model planned a larger parser task.

That single repo fact mattered more than the model price.

## The Cost Table That Looks Boring But Isn't

At one checkpoint, the numbers looked almost tied.

| Sonnet planning run |  Cost | API time | Wall time | Resulting implementation delta |
|---------------------|------:|---------:|----------:|-------------------------------:|
| **Compression ON**  | $4.39 |   16m48s |    49m42s |                   +1331 / -137 |
| **Compression OFF** | $4.40 |   17m22s |    44m00s |                   +1152 / -165 |

A quick read says compression did not matter.

That read misses the implementation shape.

By this checkpoint, the compression-OFF run had already found
`parseVersionList()` and narrowed the task to wiring the helper. The
compression-ON run was still carrying a larger parser-work premise.

Read it again with the shape in mind. Both runs cost about the same.
The compression-OFF run reused an existing helper; the compression-ON
run was set up to rebuild that logic from scratch. So compression
*did* save context budget. The saving was then spent carrying a less
accurate premise. The dollars came out even; the work did not.

!!! warning "Compression Fails Quietly"
    Compression did not fail loudly. It produced coherent artifacts.
    
    They were useful artifacts. It was solving the problem and meeting
    every product requirement.

    It was also expanding **the wrong-sized job**. That is the dangerous
    part: a failure that ships clean.

## The Patch Quality Review

A frontier-model-assisted static review told a cleaner story than raw
cost.

| Case                                            | Verdict                                      | Why / operating read                                                                                                                     |
|-------------------------------------------------|----------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| Sonnet compression OFF / Sonnet implementation  | **Preferred** quality baseline               | Found existing parser, wired both commands, strongest command tests. Minor caveat: version parsing moved before some source/auth checks. |
| Sonnet compression OFF / Haiku implementation | **Acceptable** bounded-executor trial        | Haiku followed the ratified plan, but tests were thinner. This supports retesting cheap execution after task shape is fixed. |
| Haiku compression ON / Haiku implementation     | Usable after steering, weak as scout         | Core behavior was right, but proof and cleanup were weaker.                                                                              |
| Sonnet compression ON / Haiku implementation  | **Acceptable**, but inherited larger premise | The implementation stayed inside the earlier parser-work shape.                                                                          |
| Sonnet compression ON / Sonnet implementation   | **Acceptable**, but larger than necessary    | Fixed the bug, but missed the existing helper during planning.                                                                           |
| Opus 4.8 x-high                                 | **Strong patch**, contract drift             | Strong and conservative. It rejected whitespace-padded selectors, which diverged from the final product decision.                        |
| Haiku YOLO                                      | **Rejected as Incomplete**                   | Fixed `"1, 2, 3"` but kept duplicated parsing and missed whitespace-only `[0]`.                                                          |

This table should not be read as "*use the cheapest model*". It says
something more significant:

Cheap execution can be tested **after** the task shape is fixed.

Cheap discovery, without a comprehensive spec, was not reliable in
this experiment. Cheap YOLO produced a patch-shaped answer, **not** an
accepted patch.

There is a large difference between:

```text
the model produced a diff
```

and:

```text
the model produced the accepted patch
```

The accepted patch is the metric that matters. That gap between cheap
production and directed judgment is the whole subject of
[Code Is Cheap. Judgment Is Not.][judgment-post], and this experiment
is the same lesson with an invoice attached.

## The Most Expensive Model Is Not Automatically Safer

The Opus run was strong. It was also **too eager**.

It rejected whitespace-padded selectors. That is a defensible CLI
grammar if you are designing from scratch. It is not what the final
compatibility contract said.

A stronger model can preserve more context, reason more carefully, and
still make a product decision you did not ask it to make.

You can argue that the model is **taking initiative** here, thinking
like a senior engineer to make the product more secure and reliable.
But this is a distinct failure mode worth watching:

* The cheap YOLO model missed parts of the contract.
* The expensive model tried to improve the contract.

!!! warning "Both Require Review"
    The lesson is not "*small models bad, large models good*". The
    useful split is three separate questions:

    * which model is **deciding** the task shape;
    * which model is **executing** a ratified task;
    * and which **human checkpoint** catches contract drift.

    Under-reach and over-reach are both contract drift. You cannot
    afford to review for only one of them.

## Scout Versus Executor

The transcripts showed the models **behaving differently**, not just
*costing differently*.

The compression-OFF Sonnet run read the repo and *changed the task*.
It found the helper and stated the situation outright:

```text
undelete.go does not have this bug; delete.go does;
parseVersionList already fixes it.
```

The structured Haiku run, from the same starting point, tended to hand
repo questions back to the human instead of answering them from the
code:

```text
Does undelete.go already have the correct behavior?
```

That is a question a careful reading pass should have closed: 

It is the line between a **scout** that establishes the task shape and a
**bounded executor** that needs the shape handed to it. 

Haiku was a **capable executor** once the task was **pinned**, 
and an *unreliable scout* before it.

The scout did not just find the smaller job; it built less of it. 

Both Sonnet runs went end-to-end through the same workflow, but the
compression-OFF run produced a smaller final diff (+1300 / -260 versus
+1542 / -202), because once the job was "*wire the parser*" there was
simply less to build.

## Token Telemetry

The telemetry had a few surprises. These counts are computed from the
raw session logs across every run (token counts only; the dollar
figures come from the run summaries):

| Scope        | Requests | Input tokens | Output tokens |  Cache read | Cache write |
|--------------|---------:|-------------:|--------------:|------------:|------------:|
| All sessions |      874 |       51,102 |       554,094 |  63,946,649 |   2,268,906 |
| Top-level    |      664 |       42,623 |       510,465 |  58,253,517 |   1,854,589 |
| Subagents    |      210 |        8,479 |        43,629 |   5,693,132 |     414,317 |

Read the cache-read column again: 

Fresh input was about 51K tokens and output about 554K, but the runs read back 
roughly **64 million** cached tokens. 

**Cache reuse**, not fresh reasoning, is where the token activity lived; 
and subagents accounted for only ~5.7M of those ~64M reads, so
they were **not** the sink either.

This is an **activity** table, not a billing table. Cache reads are
cheap per token, which is why a run can move 64 million of them
without the dollar figure exploding: the money lives in the cost
tables above; the attention lives here.

Findings:

1. Cache-read tokens dominated the token profile.
2. **Subagents were not the main token sink**.
3. **Haiku was cheaper in dollars, not necessarily smaller in raw
   token activity**.
4. The preferred Sonnet result was not better because it reasoned
   less. It was better because it found the smaller job: wire the
   existing `parseVersionList()` helper instead of creating or
   extracting a new parser.

The raw activity matters because "*cheap*" can mean several different
things:

* it can mean cheaper dollars;
* it can mean fewer tokens;
* it can mean less wall-clock time;
* it can mean fewer review minutes.

Those are **not** the same thing.

!!! danger "Accounting Fraud With a Patch File"
    In this experiment, the cheap YOLO attempt had the lowest cost. 

    It also **galactically** missed the contract.

    Counting that as a win would be accounting fraud with a patch file
    attached.

## Workflow Enhancement

The code bug was easy to describe.

The workflow enhancement requirement was more subtle.

The structured flow did many things right:

* it found behavior,
* produced artifacts,
* recorded non-goals,
* and guided implementation.

The gap was earlier and more mechanical:

```text
before the spec expands, prove what already exists
```

The helper existed. The workflow needed to force that fact into the
first controlling artifact.

Without that inventory, the pipeline can faithfully expand a plausible
task that is larger than necessary.

That is how you get a detailed spec for the wrong-sized job. It is the
same failure that [The Dog Ate My Homework][homework-post] documented
from the other direction: the expensive mistakes happen when an agent
writes before it has truly read.

## The Fix: Make Problem-Space Inventory a Hard Gate

The operating model I would use after this experiment is
artifact-gated:

```text
/plan:
    produce a repo-grounded debated brief
    include implementation inventory and task-shape correction

human ratification:
    confirm the debated brief before it becomes spec input

/spec:
    turn the ratified brief into a product/engineering contract

human ratification:
    confirm the spec intent before spec-kit expands it

spec-kit:
    generate spec/plan/tasks/analyzer output from the ratified intent

human ratification:
    confirm the generated tasks before coding starts

implementation:
    execute the ratified task list
    do not re-open task shape unless review sends it back

acceptance:
    measure accepted patch cost, not attempt cost
```

Implementation model choice happens **only after** the task list is
ratified.

!!! tip "Which Model for Which Job"
    * Use a **cheaper** model only when the task is well-defined and
      the spec is crystal clear beyond any reasonable doubt.
    * Use the **default** model most of the time; you will still need a
      decent spec, not a two-paragraph prompt.
    * Use a **stronger** model only when the implementation requires
      judgment, security reasoning, broad refactoring, fresh
      discovery, or adversarial scrutiny. Ironically, your spec here
      needs to be *crisper*, not looser: the model will attack it,
      find gaps, and fix them if it decides that is the right call. Be
      very clear about why you need what you need.

## What `/plan` Must Prove

For this class of task, `/plan` should not finish until it records:

* existing helpers
* existing tests
* similar commands
* already-correct behavior
* files that should remain unchanged
* whether the task is wiring, deletion, extraction, or new behavior

For the CLI bug, that inventory would have found:

* `parseVersionList()` already exists
* delete uses duplicated parsing
* undelete has similar code but not the same bug
* the accepted fix is wiring, not parser design

That would have prevented the larger parser-work premise from
surviving into the spec.

The spec was not the problem; the input to the spec was
underspecified.

## The `spec-kit` Gotcha

Spec generation is seductive because it makes the work look settled.

A generated task list feels like progress:

* it has IDs,
* it has dependencies,
* it has phases,
* it has checkboxes.

But if the task shape is **wrong**, the checkboxes become **a very
tidy way to do extra work**.

That does not make `spec-kit` (*or equivalent tools*) bad. It means
`spec-kit` should not be the first place where repo understanding
becomes concrete. Use spec expansion **after** a debated brief is
ratified.

Also, do not assume the implementation command is an interactive
review loop. Treat it as an executor. If you need a checkpoint after
every task or phase, enforce that in the wrapper or in the prompt:

```text
implement T001
stop
summarize diff and tests
wait for approval before T002
```

There is a subtler trap. The generated artifacts are themselves model
output. In one run the structured flow held the "*stop before commit*"
line well, but the generated task list quietly reintroduced "*commit
after each phase*" language that had to be edited back out. The
workflow can constrain a cheap model; the workflow's own artifacts
still need a human read.

!!! warning "Trust Is Not a Workflow"
    Trusting the final result is **not** a workflow. It is a hope with
    a diff.

## Context Compression

Compression is attractive because it reduces what the model has to
carry, saving valuable dollars.

That is also the risk:

Compression can preserve conclusions while dropping the dull facts
that made the conclusion safe.

In this experiment, the dull fact was a parser helper.

* No architecture diagram screams about an existing helper.
* No spec requirement says "*check whether this already exists*"
  unless you make it say that.
* No generated task list rescues you if the earlier artifact already
  chose the wrong implementation shape.

This is the same shape as the [watermelon-rind
anti-pattern][watermelon-post]: a mechanism that answers the question
asked can quietly prevent the discovery of the question you should
have asked. A graph tool did it there by collapsing the search space; in this
run, compression did it by dropping the boring line that would have
changed the plan. And it is the mirror image of
[The Attention Budget][attention-post]: more context is not
automatically better, but less context is not automatically cheaper
either.

!!! tip "Compression Needs a Counterweight"
    Do not treat compression as free savings. Pair it with:

    * **inventory** before compression hardens into a plan;
    * **ratification** before spec expansion;
    * **review** before implementation.

## The YOLO Fun

The cheap YOLO run was useful because it showed the **trap**.

It quickly fixed the visible symptom:

* a shallow test could have passed;
* the diff would look reasonable in a hurry.

However, it did not preserve the full accepted behavior:

* it did not remove duplicated parsing;
* it did not handle whitespace-only `[0]`.

This is the difference between symptom repair and contract repair.

A model can pass the obvious bug report while failing the actual
engineering task.

## What I Would Repeat

I would repeat the Sonnet compression comparison on more tasks.

This task says:

```text
compression OFF found the smaller job
compression ON planned the larger job
```

That is one task, not a universal law.

I would also repeat the "*strong model plans, cheaper model
executes*" pattern, but now with a **strict acceptance review**. The
result is interesting because it **may** reduce cost after the task
shape is fixed.

The numbers hint at why it is worth a look. Once the task was pinned,
the Haiku edit on top of the ratified Sonnet plan was about ninety
lines and cost roughly fifty cents. The composite came to about $4.92
against $6.11 for Sonnet end-to-end: close to 20% cheaper. 

That saving is real **only if** the cheap patch survives acceptance review, 
which is a big *if*, not a default rule.

!!! danger "Never Be Frugal on Planning"
    After this set of experiments I am fairly convinced that **a
    cheaper model should never own planning and spec development**.

    If there is a place you should not be frugal, that is the place.
    Skimp there and you may confidently implement the wrong thing:
    something that passes all the tests and looks right at first
    glance, even under the review of an excellent engineer who is not
    fully familiar with the domain.

The next experiments should separate:

* repo discovery quality
* spec quality
* implementation quality
* review cost
* accepted patch cost

Because each of those has to be judged on its own. Roll them into a
single "cost" number and you lose the distinction that matters most:
attempt cost versus accepted-patch cost.

## What This Experiment Can't Tell You

I want to be honest about the edges of this, so the numbers are not
read as more than they are:

* It is one small, local task in one repository. Larger cross-file or
  cross-repo work could move the budget picture either way.
* Tests were not run on every implementation, so some patches are
  judged by static review, not by a green test suite.
* Static review ran on the final working trees, and a few rows in the
  underlying cost ledger are interpolated rather than separately
  captured.
* Dollar costs come from the run summaries; the token counts come from
  the session logs. They are two lenses, not one ledger.

None of this changes the shape of the finding. It does mean the right
reading is "*strong signal from a weekend*", not "*proven law*".

## The Takeaway

**The expensive part of AI coding is not always the diff.**

Sometimes the expensive part is missing the smaller job.

In this experiment, the best result came from finding an existing
helper before the task shape hardened. Once the task became "*wire the
parser*", implementation was straightforward. When the helper was
missed, the workflow still produced coherent specs and acceptable
patches, but it carried a larger premise.

The workflow change is small:

```text
make implementation inventory mandatory before spec expansion
```

* Find what already exists.
* Ratify that understanding.
* Then generate the spec.
* Then implement.

!!! quote "**If You Remember One Thing from This Post...**"
    **A patch is cheap only after you know which patch you are asking
    for.**

    The cheapest model can miss the contract; the most expensive model
    can rewrite it. Neither is safe without a ratified task shape and a
    human checkpoint that reads the diff against the contract, not
    against the bug report.

## Where This Connects

This experiment is one more data point in a thread that runs through
these field notes.

* [The Dog Ate My Homework][homework-post] argued that the hard part
  is getting an agent to **read before it writes**. This is the same
  failure with money attached: the agent that did not inventory the
  repo first wrote a spec for the wrong-sized job.
* [The Watermelon-Rind Anti-Pattern][watermelon-post] showed that a
  mechanism which answers the question asked can **prevent the
  discovery of the question you should have asked**. Compression did
  exactly that here.
* [Code Is Cheap. Judgment Is Not.][judgment-post] put it in one line:
  production is the easy part, judgment is the hard part. The judgment
  that mattered most was not in the diff. It was in deciding the task
  shape before any model started typing.
* [The Attention Budget][attention-post] explained why more context is
  not automatically better. Compression is the same coin flipped: less
  context is not automatically cheaper.

If you want the operating model in tool form, `ctx` already ships most
of it. [Design Before Coding][design-recipe] walks the
brainstorm / plan / spec / implement chain, and [Scrutinizing a
Plan][scrutinize-recipe] is the `/ctx-plan` step that produces the
repo-grounded debated brief this post keeps asking for: the artifact
that forces "*prove what already exists*" before a spec can expand.

---

*This post is part of the [`ctx` field notes][blog] series,
documenting what we learn building persistent context infrastructure
for AI coding sessions. The experiment ran against the
[SPIFFE/SPIKE](https://github.com/spiffe/spike) repository using
Anthropic models only. The numbers are signal from a weekend's worth
of runs, not a peer-reviewed result.*

[homework-post]: 2026-02-25-the-homework-problem.md
[watermelon-post]: 2026-04-06-the-watermelon-rind-anti-pattern.md
[judgment-post]: 2026-02-17-code-is-cheap-judgment-is-not.md
[attention-post]: 2026-02-03-the-attention-budget.md
[design-recipe]: ../recipes/design-before-coding.md
[scrutinize-recipe]: ../recipes/scrutinizing-a-plan.md
[blog]: index.md
