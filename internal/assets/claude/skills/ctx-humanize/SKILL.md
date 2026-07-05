---
name: ctx-humanize
description: "Review, rewrite, or edit human-facing prose to remove formulaic LLM writing patterns while preserving meaning, intent, and certainty. Use when the user says 'humanize this', 'this sounds like AI', 'remove the AI tells', 'de-AI this', 'make it sound human', asks whether text reads as AI-generated, or wants a blog post, doc, README, or announcement polished before publishing."
allowed-tools: Read, Write, Edit, Grep, Glob, AskUserQuestion
---

Make prose sound like a person wrote it, not like a model
assembled the most statistically likely version of an answer.

This applies to human-facing prose: blog posts, essays,
user-facing docs, READMEs, PR descriptions, release
announcements, leadership summaries. It does not apply to code,
logs, structured data, test output, or machine-readable content.
In ctx projects, also leave `.context/` knowledge files,
`specs/`, and `SKILL.md` files alone unless the user explicitly
asks: they are agent-facing operational text, and "sounding
human" is not their job.

The goal is not to make everything casual. It is to make the
writing appropriate for its audience, specific in its claims,
and free of the tells that make readers discount it.

## Invariants

These hold in every mode. They exist because a humanizing pass
that changes what the text *says* is worse than no pass at all.

1. Preserve the author's meaning.
2. Preserve the author's level of certainty. Do not make weak
   claims sound stronger just because the prose gets cleaner.
3. Preserve document structure unless asked to restructure.
4. Do not add facts, examples, citations, names, numbers,
   anecdotes, opinions, or emotions that are not in the source.
   If a rewrite needs voice, draw it from stances the author
   actually expressed. Never invent a first-person reaction.
5. Do not flatten a real voice into corporate-neutral mush.
6. Prefer the smallest edit that makes the prose better.

When the text is already good, say so and stop. Do not rewrite
to prove the skill ran.

## Modes

Default to **Review** unless the user asks for a rewrite or
asks you to edit a file.

**Review**: the user wants a verdict or feedback. Report:
the main AI tells found (with locations), the highest-value
edits, and any places where rewriting would risk meaning,
voice, or precision. Do not modify files.

**Rewrite**: the user asks to humanize, polish, or rewrite.
Return the rewritten text plus a brief summary of what changed.

**Apply**: the user asks you to edit a file. Prefer Edit over
Write; use Write only for new files or explicit full
replacement. Afterwards report files changed, kinds of changes,
and anything intentionally left as-is.

## Protected Content

Improve the prose *around* these; the items themselves survive
unchanged, because their value is exactness:

- fenced code blocks, inline code, shell commands
- file paths, URLs, package/API/function names, flags, env vars
- issue/ticket IDs, CVE IDs, commit SHAs, PR numbers
- quoted text, unless the user explicitly asks to edit a quote
- legal, compliance, or security wording where exact language
  matters
- RFC-2119 keywords (MUST, SHOULD, MAY, SHALL, REQUIRED,
  OPTIONAL)

## Process

1. Read the full input. For file targets, Read the whole file
   rather than sampling it.
2. Identify *clustered* AI patterns using the index below. Read
   `references/pattern-catalog.md` (from this skill's directory)
   before any substantial rewrite: it has before/after examples
   for every pattern and the guardrails for adding voice.
3. Pick the mode and act on it. Steps 4 through 6 apply when
   you produce a rewrite (Rewrite and Apply modes); in Review
   mode, report the findings those steps would have caught.
4. Draft, then ask: "what still makes this sound generated?"
   and revise once more. First drafts of de-AI'd text usually
   still carry essay symmetry.
5. Run the typography check (below).
6. Verify no facts were invented, no claims strengthened or
   weakened, no coverage lost: if the original made five
   substantive points, the rewrite still makes five.
   Zero-content filler (generic optimism, throat-clearing,
   chatbot residue) does not count as a substantive point;
   deleting it is not a coverage loss.

## Typography

Many readers now treat certain typography as machine residue,
so normalize it in human-facing prose unless the user opts out:

- smart quotes and apostrophes → ASCII quotes
- en dashes → hyphens
- em dashes → replaced *semantically*, never by blind
  search-and-replace: a period when the dash starts a new
  thought, a comma for a tight aside, a colon before an
  explanation, parentheses for a true aside, or restructure
  the sentence. Also catch ` -- ` double-hyphens used as dashes.

The final text must contain no em or en dashes. Verify
mechanically, not by eye: Grep the result for `—`, `–`, ` -- `,
and smart-quote characters before returning. For inline
rewrites, write the draft to a scratch file and grep that. In
repos that ship `hack/detect-ai-typography.sh` (the ctx repo
does), the script runs the same check, but note it takes a
*directory* of markdown files, not a single file path.

## Detection Guidance

Flag clusters, not isolated hits. Clean human writing can
trip several patterns; one "additionally" or one bold phrase
is not a defect. A paragraph that stacks inflated significance
on a forced triplet on a generic conclusion is.

Not reliable indicators on their own: perfect grammar,
consistent style, dry prose, formal vocabulary, curly quotes,
one short emphatic sentence, unsourced claims, correct or
complex formatting.

Preserve signs of an actual human: specific hard-to-fabricate
detail, mixed feelings, unresolved tension, era-bound
references, first-person editorial choices, varied sentence
length, useful asides, parentheticals, self-corrections. These
are often the point of the piece. Do not sand them away.

## Pattern Index

Full catalog with before/after examples and watch lists:
`references/pattern-catalog.md`. The short version:

| # | Pattern | Tell |
|---|---------|------|
| 1 | Inflated significance | "pivotal moment", "testament to", "broader trends" |
| 2 | Forced notability | listing media coverage instead of saying what happened |
| 3 | Superficial -ing analysis | "highlighting...", "underscoring...", "reflecting..." |
| 4 | Brochure language | "vibrant", "nestled", "seamless", "stunning" |
| 5 | Vague attribution | "experts argue", "observers have noted" |
| 6 | Challenges-and-outlook filler | "despite these challenges", "the future looks bright" |
| 7 | Forced triplets | three parallel fragments where one point would do |
| 8 | Overused AI vocabulary | "delve", "crucial", "landscape", "tapestry" clustering |
| 9 | Copula avoidance | "serves as", "boasts", "features" instead of "is"/"has" |
| 10 | Negative parallelisms | "not just X, it is Y"; tailing negations |
| 11 | Elegant variation | synonym-cycling when repetition is clearer |
| 12 | False ranges | "from X to Y" with no real scale |
| 13 | Reflexive passive voice | passives where the actor matters |
| 14 | Boldface spray | mechanically bolded phrases |
| 15 | Inline-header bullets | "**Label:** sentence" lists that should be prose |
| 16 | Mechanical headings | heading + one-line warm-up restating the heading |
| 17 | Emoji decoration | emojis as bullet/heading ornaments |
| 18 | Chatbot residue | "Great question!", "I hope this helps", "let's dive in" |
| 19 | Cutoff disclaimers / gap-filling | "as of my last update"; inventing plausible filler around missing facts |
| 20 | Sycophancy | flattering the reader instead of answering |
| 21 | Filler phrases | "in order to", "due to the fact that", "it is important to note" |
| 22 | Excessive hedging | "could potentially possibly" |
| 23 | Generic positive conclusions | "exciting times lie ahead" |
| 24 | Authority tropes | "at its core", "the real question is" |
| 25 | Diff-anchored writing | docs narrating a change instead of the current state |
| 26 | Staccato drama | chains of clipped fragments engineered for quotability |
| 27 | Aphorism formulas | "X is the Y of Z", "the currency of" |
| 28 | Theatrical openers | standalone "Honestly?", "Here's the thing" |

## Voice Calibration

If the user provides a writing sample, read it first. Match the
deeper rhythm (how the writer argues, how much context they
give, how directly they land claims, how much mess they allow),
not surface quirks. When editing a project's blog or docs and
no sample is given, read one or two existing published pieces
from the same directory; an established publication voice beats
the default. Absent both, write natural, direct, varied,
specific, and mildly opinionated where the content allows.

## Example

Before:

> Additionally, this groundbreaking tool serves as a testament
> to the transformative potential of context persistence,
> ensuring seamless workflows and fostering collaboration
> across teams — from solo developers to large enterprises.

After:

> The tool keeps context across sessions. Solo developers use
> it to resume work without re-explaining; teams use it to
> share decisions.

What changed: dropped the vocabulary cluster (pattern 8), the
copula avoidance (9), the -ing chain (3), the false range (12),
and the em dash; replaced abstractions with the two concrete
claims the original was gesturing at. Nothing new was invented.

## Quality Checklist

Before returning:

- [ ] Mode honored (Review touched no files)
- [ ] Meaning, certainty, and coverage preserved
- [ ] No invented facts, opinions, or first-person emotion
- [ ] Protected content byte-identical
- [ ] Grep confirms no em/en dashes or smart quotes remain
- [ ] The text no longer reads as generated. If it was already
      fine, you said so instead of editing
