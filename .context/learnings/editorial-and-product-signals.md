# editorial-and-product-signals

## [2026-06-07-170006] Editorial KB pipeline: design epistemology (consolidated)

**Consolidated from**: 5 entries (all 2026-05-10)

- An ongoing user paying concrete workaround tax (disabled skills, hand-typed
  closeouts, colliding root constitution files) is the strongest validation
  evidence — beats user research, N=2 discussion, "seems useful." Use the
  workaround details as the inverse-spec; ship the shape they hand-rolled and
  use their project as the regression corpus.
- When lifting from a battle-tested external design, lift the renames and
  disambiguation moves alongside the features: intentional renames encode
  resolved conflicts (KB-RULES.md not CONSTITUTION.md; domain-decisions.md not
  DECISIONS.md). Treating them as cosmetic re-litigates the underlying fight.
- KB epistemology: a knowledge base has no "decide" moment — only
  evidence-capture events with confidence bands (>0.9 = decided by contract).
  Even NL assertions ("anchor on this") are evidence-capture, not
  decision-capture. So a parallel /ctx-kb-decide skill is the wrong shape; the
  pipeline-only-writer model is ontologically correct. General check: "I chose
  between alternatives" vs "I learned about the world."
- Recursive composability eliminates feature classes: a KB of KBs is a KB
  (source-map kind: kb + the standard ingest pipeline covers federation; no v1
  schema lockout). Ask whether the standard pipeline pointed at its own output
  covers a "thing-of-things" before designing a new mechanism.
- The LLM is the migration tool: every category of being-wrong about a schema
  (ID renumbering, taxonomy reshuffle, band remapping, path renames) is cheap
  because LLM cleanup absorbs the migration. Commit to the readable, opinionated
  v1 schema instead of hedging with abstract types; surface dirty state via
  doctor advisories so the agent has a work surface.

---

## [2026-05-28-215214] ctx kb: single topic-enumeration site; life-stage count is consumer-side

**Context**: kb reindex blanked the CTX:KB:TOPICS block for grouped kbs
(things-wtf-dr regrouped 49 topics into folders); the task speculated a sibling
life-stage topic-count glob was also affected.

**Lesson**: reindex.ListTopics (internal/cli/kb/core/reindex/topic.go) is the
ONLY topic enumeration/count in ctx, and CTX:KB:TOPICS is the only managed
block. The life-stage concept in ctx is the ingest/closeout frontmatter field,
unrelated to topics. Any per-life-stage topic count lives in the consumer kb,
which ctx neither generates nor owns.

**Application**: Localize nested-topic fixes to ListTopics; treat
per-group/per-life-stage topic counts as consumer territory (same recurse +
exclude-group-landing pattern, fixed in their repo).

---

## [2026-05-17-200000] Creator confusion is the strongest doc-quality signal — louder than any user signal

**Context**: In this session the project author asked *"why
external sources only? I can ground on a repo, a MCP query, a
markdown I dropped into ./inbox — are they also considered
'external'. Or is there a nomenclature confusion here?"* — and
explicitly noted *"it is confusing to the very creator of this
pipeline. -- and that's not a good sign."* Investigation
confirmed the input contract accepted in-tree paths and MCP
resources all along, but the SKILL.md ledes, the CLI docs table
row, and the recipe all framed ground as "external" — which
the creator's own mental model couldn't reconcile with the
contract.

**Lesson**: A normal-user reading-confusion signal is "I don't
understand this." A creator reading-confusion signal is "this
contradicts what I built." The second is louder by an order
of magnitude — the creator has the full internal model and a
strong prior on what the system should say. If they trip over
the words, the words are wrong, full stop. Don't defend the
existing framing; don't explain what was meant. Rewrite to
match what the contract actually does. The creator was a
control instrument; if even that instrument deflected, the
docs are mis-anchoring everyone.

**Application**: When the project's own creator asks a
"do we even need X?" or "wait, isn't X actually doing Y?"
question, treat it as a doc-bug report, not an architecture
question. Investigate the literal contract (input/output
shapes, code-level reality) before debating semantics. If the
contract is correct but the docs misframe it, the action is
"rewrite the framing across every doc surface that touches
it" — skills, recipes, CLI tables, anything user-facing.
Concrete instance handled this session: dropped "external"
from ctx-kb-ground's prose, description, pass-mode value,
recipe Step 4, and CLI table row across all three skill trees.

---

## [2026-03-05-023941] Blog post editorial feedback is higher-leverage than drafting

**Context**: Draft of Agent Memory Is Infrastructure was publication-quality on
first pass; user editorial feedback (structural emphasis, rhetorical sharpening,
amnesia/archaeology bridge) elevated it significantly more than initial
generation

**Lesson**: For narrative content, the first draft captures the argument; the
editorial pass captures the voice. Both are necessary but the editorial pass has
disproportionate impact on quality.

**Application**: For future blog posts, invest more in the editorial cycle
(structural feedback then targeted refinements) rather than trying to nail voice
on first generation.

---

## [2026-02-26-100010] ctx add and decision recording (consolidated)

**Consolidated from**: 4 entries (2026-01-27 to 2026-02-14)

- `ctx add learning` requires `--context`, `--lesson`, `--application` flags.
  `ctx add decision` requires `--context`, `--rationale`, `--consequence`. A
  bare string only sets the title and the command will fail without required
  flags.
- Structured entries with Context/Lesson/Application are more useful than
  one-liners. Agents are guided via AGENT_PLAYBOOK.md.
- Always complete decision record sections — placeholder text like "[Add
  context here]" is a code smell. Decisions without rationale lose their value
  over time.
- Slash commands using `!` bash syntax require matching permissions in
  settings.local.json. When adding new /ctx-* commands, ensure ctx init
  pre-seeds the required `Bash(ctx <subcommand>:*)` permissions.

---

## [2026-01-28-051426] IDE is already the UI

**Context**: Considering whether to build custom UI for .context/ files

**Lesson**: Discovery, search, and editing of .context/ markdown files works
better in VS Code/IDE than any custom UI we'd build. Full-text search,
git integration, extensions - all free.

**Application**: Don't reinvent the editor. Let users use their preferred IDE.

---


*Module-specific, niche, and historical learnings:
[learnings-reference.md](learnings-reference.md)*
