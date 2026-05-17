<!--
Topic-page template. `ctx kb topic new "<name>"` copies this
file into `.context/kb/topics/<slug>/index.md` and
substitutes the slug + display name. Sub-pages
(`security.md`, `multi-surface.md`, ...) split off lazily
from this index when the cold-reader rubric's "boundaries
clear?" check fails.

The `topic-page` pass of `ctx kb ingest` is expected to fill
`TBD` placeholders with cited prose. Do not hand-edit this
template unless you are changing the topic-page shape for
the whole project.
-->

---
Subject: TBD; short subject line, max 80 chars
Last verified: TBD-YYYY-MM-DD
Author: TBD; agent or human name
Confidence: TBD; high | medium | low | speculative
sha: TBD; short SHA at last verification (in-repo context)
branch: TBD; branch at last verification
---

# <Topic name>

TBD: one-paragraph lede. State what this topic is, in plain
language, in 2-4 sentences. The cold-reader rubric's
"concept clear?" item is judged against this paragraph. Cite
at least one `EV-###` row from `evidence-index.md` here.

---

## What It Is

TBD: definition prose. Describe the thing this page is
about: shape, mechanism, surface area, key behaviors. Cite
evidence rows inline as `[EV-###]`. Confidence band on each
substantive claim, demoted to the floor of cited bands per
`../../../ingest/KB-RULES.md` "Confidence bands".

If this section grows past the cold-reader rubric's
"boundaries clear?" threshold, propose a sub-page split
(e.g. `mechanics.md`) on the next ingest pass; do not
auto-split.

---

## Why This KB Cares

TBD: relevance statement. Why does this kb track this
topic? What downstream questions does it answer? What
decisions or claims elsewhere in the kb lean on it? The
cold-reader rubric's "why this kb cares clear?" item is
judged against this section.

---

## Sources and Further Reading

TBD: bullet list of the canonical sources this page draws
from, by source short name as defined in
`../../source-map.md`. Each bullet pairs a source short name
with a one-line description.

- `<SOURCE-SHORTNAME>`: TBD one-line description; locator
  hint.

For deep evidence, the cold-reader rubric's "canonical
evidence reachable?" item requires that walking from this
page to `../../evidence-index.md` to the original source
takes one or two clicks. If reachability degrades, this
section needs revision.

---

## Related Concepts in This KB

TBD: bullet list of adjacent topic pages, surfaced by the
topic-adjacency pre-flight on the most recent ingest pass.
Each bullet pairs a topic slug with a one-line description
of the relationship.

- `<adjacent-slug>`: TBD one-line description of the
  relationship.

When the adjacency pre-flight finds no incomplete adjacent
topics, this section records `none surfaced` explicitly.
Silence is not the same as a clean pre-flight; see
`../../../ingest/KB-RULES.md` "Topic-adjacency pre-flight".
