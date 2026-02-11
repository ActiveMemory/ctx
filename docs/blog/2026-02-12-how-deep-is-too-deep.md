---
title: "How Deep Is Too Deep?"
date: 2026-02-12
author: Jose Alekhinne
topics:
  - AI foundations
  - abstraction boundaries
  - agentic systems
  - context engineering
  - failure modes
---

# How Deep Is Too Deep?

![ctx](../images/ctx-banner.png)

## When "Master ML" Is the Wrong Next Step

*Jose Alekhinne / 2026-02-12*

!!! question "Have you ever felt like you should understand more of the stack beneath you?"
    You can talk about transformers at a whiteboard. You can explain attention
    to a colleague. You can use agentic AI to ship real software.

    But somewhere in the back of your mind, there is a voice:
    *"Maybe I should go deeper. Maybe I need to master machine learning."*

I had that voice for months. Then I spent a week debugging an agent failure
that had nothing to do with ML theory and everything to do with **knowing
which abstraction was leaking**.

This post is about when depth compounds and when it doesn't.

## The Hierarchy Nobody Questions

There is an implicit stack most people carry around when thinking about AI:

| Layer            | What Lives Here                                  |
|------------------|--------------------------------------------------|
| Agentic AI       | Autonomous loops, tool use, multi-step reasoning |
| Generative AI    | Text, image, code generation                     |
| Deep Learning    | Transformer architectures, training at scale     |
| Neural Networks  | Backpropagation, gradient descent                |
| Machine Learning | Statistical learning, optimization               |
| Classical AI     | Search, planning, symbolic reasoning             |

At some point down that stack, you hit a comfortable plateau: the layer
where you can hold a conversation but not debug a failure.

The instinctive response is to go deeper.

**But that instinct hides a more important question**: Does depth still
compound when the abstractions above you are moving hyper-exponentially?

## The Uncomfortable Observation

If you squint hard enough, a large chunk of modern ML intuition collapses
into older fields:

| ML Concept       | Older Field                        |
|------------------|------------------------------------|
| Gradient descent | Numerical optimization             |
| Backpropagation  | Reverse-mode autodiff              |
| Loss landscapes  | Non-convex optimization            |
| Generalization   | Statistics                         |
| Scaling laws     | Asymptotics and information theory |

Nothing here is *uniquely* "AI."

Most of this math predates the term **deep learning**. 

In some cases, **by decades**.

So what changed?

## Same Tools, Different Regime

The mistake is assuming this is a **new theory problem**. It isn't.

It is a **new operating regime**.

Classical numerical methods were developed under assumptions like:

* manageable dimensionality
* reasonably well-conditioned objectives
* losses that actually represent the goal

Modern ML violates all three. *On purpose.*

Today's models operate with millions to trillions of parameters, wildly
underdetermined systems, and objective functions we know are wrong but
optimize anyway. At this scale, familiar concepts warp:

* "*Local minima*" mostly dissolve into saddle points
* Noise stops being noise and starts becoming structure
* Overfitting can coexist with generalization
* Bigger models outperform "better" ones

The math didn't change. **The phase did.**

This is **less numerical analysis and more statistical physics*. 

Same equations: Turbulent systems.

## Why Scaling Laws Feel Alien

In classical statistics, asymptotics describe what happens *eventually*.

In modern ML, scaling laws describe **where you can operate today**.

They don't say: *"Given enough time, things converge."*

They say: *"Cross this threshold and behavior qualitatively changes."*

This is why **dumb architectures** plus **scale** beat clever ones. 

Why small theoretical gains disappear under data. 

Why "**just make it bigger**", ironically, keeps working longer than it should.

That is not a triumph of ML theory. 

It is a property of high-dimensional systems under loose objectives.

## Where Depth Actually Pays Off

This reframes the original question.

You don't need depth because this is "AI." You need depth where
**failure modes propagate upward**.

I learned this building `ctx`. The agent failures I have spent the most
time debugging were never about the model's architecture. 

They were about:

* **Misplaced trust**: The model was confident. The output was wrong.
  Knowing *when* confidence and correctness diverge is not something you
  learn from a textbook. You learn it from watching patterns across
  hundreds of sessions.
* **Distribution shift**: The model performed well on common patterns
  and fell apart on edge cases specific to this project. Recognizing
  that shift before it compounds requires understanding *why*
  generalization has limits, not just *that* it does.
* **Error accumulation**: In a single prompt, model quirks are tolerable.
  In [autonomous loops][loops-post] running overnight, they compound.
  A small bias in how the model interprets instructions becomes a large
  drift by iteration 20.
* **Scale hiding errors**: The model's raw capability masked problems
  that only surfaced under specific conditions. More parameters didn't
  fix the issue. It just made the failure mode rarer and harder to
  reproduce.

[loops-post]: 2026-02-09-defense-in-depth-securing-ai-agents.md

This is the kind of **depth that compounds**. 

Not deriving backprop.

Understanding **when the math lies to you**.

## The Connection to Context Engineering

This is the same pattern I keep finding at different altitudes.

In [The Attention Budget][attention-post], I wrote about how dumping
everything into the context window degrades the model's focus. The fix
was not a better model. It was **better curation**: load less, load the
right things, preserve signal per token.

In [Skills That Fight the Platform][fight-post], I wrote about how
custom instructions can conflict with the model's built-in behavior.
The fix was not deeper ML knowledge. It was **understanding** that the
model already has judgment and you should **extend** it, not override it.

In [You Can't Import Expertise][import-post], I wrote about how generic
templates fail because they don't encode project-specific knowledge. A
consolidation skill with eight Rust-based analysis dimensions was 70%
noise for a Go project. The fix was not a better template. It was
**growing expertise from this project's own history**.

[attention-post]: 2026-02-03-the-attention-budget.md
[fight-post]: 2026-02-04-skills-that-fight-the-platform.md
[import-post]: 2026-02-05-you-cant-import-expertise.md

In every case, the answer was not "go deeper into ML." The answer was
**know which abstraction is leaking and fix it at the right layer**.

## Agentic Systems Are Not an ML Problem

Agentic AI is a **systems problem under uncertainty**:

* Feedback loops between the agent and its environment
* Error accumulation across iterations
* Brittle representations that break outside training distribution
* Misplaced trust in outputs that look correct

In short-lived interactions, model quirks are tolerable. In long-running
autonomous loops, they compound. That is where shallow understanding
becomes expensive.

But the understanding you need is not about optimizer internals. It is
about:

| What Matters                                              | What Doesn't (for Most Practitioners)             |
|-----------------------------------------------------------|---------------------------------------------------|
| Why gradient descent fails in specific regimes            | How to derive it from scratch                     |
| When memorization masquerades as reasoning                | The formal definition of VC dimension             |
| Recognizing distribution shift before it compounds        | Hand-tuning learning rate schedules               |
| Predicting when scale hides errors instead of fixing them | Chasing theoretical purity divorced from practice |

The depth that matters is **diagnostic**, not **theoretical**.

## The Real Answer

Not "all the way down."

**Go deep enough to**:

* Diagnose failures instead of cargo-culting fixes
* Reason about uncertainty instead of trusting confidence
* Design guardrails that align with model behavior, not hope

**Stop before**:

* Hand-deriving gradients for the sake of it
* Obsessing over optimizer internals you will never touch
* Chasing theoretical purity divorced from the scale you actually
  operate at

This is not about mastering ML.

It is about **knowing which abstractions you can safely trust and which
ones leak**.

## The `ctx` Lesson

Every design decision in `ctx` is downstream of this principle.

The [attention budget][attention-post] exists because the model's
internal attention mechanism has real limits. You don't need to
understand the math of softmax to build around it. But you do need to
understand that **more context is not always better** and that attention
density degrades with scale.

The [skill system][anatomy-post] exists because the model's built-in
behavior is already good. You don't need to understand RLHF to build
effective skills. But you do need to understand that **the model already
has judgment** and your skills should teach it things it doesn't know,
not override how it thinks.

[anatomy-post]: 2026-02-07-the-anatomy-of-a-skill-that-works.md

[Defense in depth][loops-post] exists because soft instructions are
probabilistic. You don't need to understand the transformer architecture
to know that a Markdown file is not a security boundary. But you do need
to understand that **the model follows instructions from context, and
context can be poisoned**.

In each case, the useful depth was one or two layers below the
abstraction I was working at. Not at the bottom of the stack.

**The boundary between "useful understanding" and "academic exercise"
is where your failure modes live.**

## Closing Thought

Most modern AI systems don't fail because the math is wrong.

They fail because we **apply correct math in the wrong regime**, then
build autonomous systems on top of it.

Understanding that boundary, not crossing it blindly, is where depth
still compounds.

And that is a far more useful form of expertise than memorizing another
loss function.

---

!!! quote "If you remember one thing from this post..."
    **Go deep enough to diagnose your failures. Stop before you are
    solving problems that don't propagate to your layer.**

    The abstractions below you are not sacred. But neither are they
    irrelevant.

    The useful depth is wherever your failure modes live:
    usually one or two layers down, not at the bottom.

---

*This post started as a note about whether I should take an ML course.
The answer turned out to be "no, but understand why not." The meta
continues.*
