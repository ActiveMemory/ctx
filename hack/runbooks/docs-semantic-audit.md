# Documentation Semantic Audit Runbook

**When to use**: Before a release, after adding several new pages, when the
site feels sprawling, or when you suspect narrative gaps. This audit finds
structural problems that linters and link checkers cannot: weak pages that
should be merged, heavy pages that should be split, missing cross-links,
and narrative arcs that don't land.

**Frequency**: Per release cycle, or when the docs surface area grows by
more than 3-4 pages.

**Time**: ~20-40 minutes with an LLM session.

---

## Why This Is a Runbook

These judgments are inherently subjective and context-dependent. A page is
"weak" relative to its neighbors; a narrative arc only matters if the docs
intend to tell a story. Deterministic tools (broken-link checkers, word
counters) can't do this. An LLM reading the full doc set can.

---

## Prompt

Paste or adapt the following into a Claude Code session. The agent needs
read access to `docs/` and `site/` (for nav structure).

```
Read every file under docs/ (including docs/blog/ and docs/recipes/).
For each file, note: title, word count, outbound links, inbound links
(how many other pages link to it), and a one-line summary of its purpose.

Then produce a report with these sections:

## 1. Weak Dangling Pages

Pages that are thin, isolated, or redundant. Signs:
- Under ~300 words with no unique content (just restates what another page says)
- Zero or one inbound links (orphaned in the nav)
- Content that would be stronger merged into an adjacent page
- "Try it in 5 minutes" sections that assume installation already happened
- Pages whose title doesn't work as a nav entry (too long, too vague)

For each: identify the page, explain why it's weak, and recommend
merge target or deletion.

## 2. Overly Heavy Pages

Pages doing too much. Signs:
- Over ~1500 words with multiple distinct topics
- More than 4 H2 sections that could stand alone
- Reader has to scroll past irrelevant content to find what they need
- Mixed audience (beginner setup + advanced config on same page)

For each: identify the page, list the distinct topics, and suggest
split points.

## 3. Missing Cross-Links

Places where a reader would naturally want to jump to related content
but no link exists. Look for:
- Concepts mentioned but not linked (e.g., "scratchpad" without linking
  to the scratchpad page)
- Blog posts that describe features without linking to the reference docs
- Recipes that reference workflows without linking to the relevant
  getting-started section
- Pages that end without a "Next Up" or "See Also" pointer

For each: source page, anchor text, suggested link target.

## 4. Narrative Gaps

The docs should tell a coherent story: problem → install → first session
→ daily workflow → advanced patterns → contributing. Look for:
- Gaps in the progression (e.g., no bridge from "first session" to
  "daily habits")
- Blog posts that introduce concepts the reference docs don't cover
- Recipes that assume knowledge no other page teaches
- Features documented in CLI reference but missing from workflows/recipes

For each: describe the gap and suggest what page or section would fill it.

## 5. Blog Cross-Linking Opportunities

Blog posts are often written in isolation. Look for:
- Posts that cover the same theme but don't reference each other
- Posts that describe the evolution of a feature (natural "part 1 / part 2")
- Posts that would benefit from a "Related posts" footer
- Thematic clusters that could be linked from a recipe or reference page

For each: list the posts, the shared theme, and the suggested links.

## Output Format

For every finding, include:
- File path (docs/whatever.md)
- Severity: high (actively confusing), medium (missed opportunity),
  low (nice to have)
- Concrete recommendation (merge into X, split at H2 Y, add link to Z)

End with a prioritized action list: what to fix first.
```

---

## After the Audit

1. **Triage findings** — not everything needs fixing. Focus on high severity.
2. **Merge weak pages first** — fewer pages is almost always better.
3. **Add cross-links** — cheapest improvement, highest reader impact.
4. **File split decisions in DECISIONS.md** — page splits are architectural.
5. **Regenerate the site** and spot-check nav after structural changes.

## History

- 2026-02-17: Created after merging `docs/re-explaining.md` into `docs/about.md`,
  which surfaced the pattern of weak standalone pages that dilute rather than add.
