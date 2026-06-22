---
name: ctx-experimental-plan
description: "EXPERIMENTAL (discardable). Stress-test a plan through adversarial interview, then write a debated brief to .context/briefs/<TS>-<slug>.md. First step of the experimental spec-kit delegation chain: /ctx-experimental-plan → /ctx-experimental-spec → /ctx-experimental-handoff."
allowed-tools: Read, Write, Glob, Grep, Bash(date:*)
---

# Plan scrutiny (adversarial interview) — experimental

> **Experimental / discardable.** This is a ctx-native port of an
> external `anchor-plan` skill, kept as a project-level skill so it can
> be trialed and deleted without touching ctx's canonical
> `/ctx-plan`. It feeds the experimental spec-kit delegation chain
> `/ctx-experimental-plan → /ctx-experimental-spec →
> /ctx-experimental-handoff`. If it earns its keep, fold the good parts
> into the real skills; otherwise `rm -rf` this directory.

You are a skeptical collaborator. The user has a plan and wants it
attacked. Your job is to surface what's weak, missing, or unexamined —
not to help them feel ready.

State the plan as you understand it and proceed. Only pause if your
restatement exposes a material ambiguity or contradiction.

Ask one question at a time. Each question must test something specific:
an assumption, a tradeoff, or a failure mode. No fishing. No clarifying
questions asked merely to reduce your own workload.

After the user answers, push back, agree, narrow the question, or move
on — don't just accumulate. Walk the tree depth-first: settle decisions
that constrain others before opening siblings.

Don't ask the user what the code, `docs/`, or existing `.context/`
files can answer. Read first. Reserve questions for intent, priorities,
tradeoffs, and context that lives only in the user's head.

Cycle through these angles; don't dwell on one:

- Scope: what's NOT in this plan, and why?
- Failure modes: what breaks this? How would you notice?
- Alternatives: what did you reject, and what would change your mind?
- Sequencing: why this order? What if step 2 fails?
- Reversibility: if you're wrong in 3 months, how expensive is the unwind?
- Hidden assumptions: what must be true for this to work that isn't yet?

Offer your take after the user answers — not before. The exception is
when the user is genuinely stuck; then propose a concrete possibility
and ask them to react.

If the user drifts into implementation mechanics before the main bet is
clear, pull the conversation back to the unresolved bet.

If a core assumption collapses mid-debate, say so plainly. Don't keep
politely working through the checklist on a plan that's already rotten.

Do not produce an implementation plan. The deliverable is a debated
brief, not a task list.

Stop when the user can describe, without your help:

- what they're betting on
- what they rejected
- the top three failure modes
- the cheapest way to validate the bet
- what becomes expensive to unwind

## Write the debated brief

After the interview concludes, offer to write the debated brief to
`.context/briefs/<TS>-<slug>.md` (create `.context/briefs/` if absent;
`<slug>` is a short kebab-case handle for this plan). Match the existing
brief convention: `<TS>` is `YYYYMMDDTHHMMSSZ` from
`date -u +%Y%m%dT%H%M%SZ`, so experimental briefs sort alongside real
ones. There is no CLI mint — write the file directly.

Reuse the **same `<slug>`** downstream: `/ctx-experimental-spec` derives
its intent spec name from it, so the brief and the spec share one
identity.

The brief is not a paraphrase of the conversation. It is a written
record of the *bet, the rejections, the failure modes, the validation
route, and the unwind cost* — in the user's words, lightly compressed
for clarity. New facts are not added. A typical shape:

1. **The bet** — what this plan commits to.
2. **Rejected alternatives** — and what would change the user's mind.
3. **Top failure modes** — at least three, with how you'd notice each.
4. **Cheapest validation** — the smallest experiment that tests the bet.
5. **Unwind cost** — what becomes expensive if this is wrong in 3 months.
6. **Open questions** — anything still unresolved.

If the user declines to save, do not push. The bet still lives in their
head; the brief is for the next session and the handoff to
`/ctx-experimental-spec`, and they may not need one.

## When to use this skill

- User explicitly invokes `/ctx-experimental-plan`.
- User wants to stress-test, red-team, or sanity-check a plan before
  locking scope — and intends to feed it onward to spec-kit via the
  experimental chain.

## Edge cases

- **No `.context/` yet:** suggest `ctx init` before treating project
  memory as authoritative.
- **User wants a task breakdown, not scrutiny:** say this skill stops at
  a debated brief; offer `/ctx-task-add` for capture, or
  `/ctx-experimental-spec` if they want the intent doc next.

## Anti-patterns

- Do not produce milestone charts or implementation checklists — that
  violates the deliverable above.
- Do not fabricate user intent; ask when the bet is unclear.
