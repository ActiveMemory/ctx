---
Subject: vLLM examples landing page as a substrate for ctx recipe-surface lessons
Last verified: 2026-05-21
Author: agent-ingested
Confidence: medium
sha: 8c02b754
branch: feat/cwd-anchored-context
---

# vLLM

vLLM is a Python-based, high-throughput LLM serving and inference
engine that publishes its example surface as a thirteen-category
index of GitHub directory pointers rather than an in-docs recipe
library [EV-001][EV-002]. This kb tracks vLLM as a *contrastive
study*, not as a tool ctx will use: vLLM's example taxonomy is
useful precisely because it is the opposite of ctx's prose-rich
recipe surface — studying the contrast surfaces design choices
ctx has made implicitly. The open question driving this topic:
what, if anything, from vLLM's example surface is worth lifting
into ctx's recipe library or doc structure?

---

## What It Is

The vLLM examples landing page at `docs.vllm.ai/en/latest/examples/`
is a single-page index that groups example code into thirteen named
categories: `basic`, `generate`, `pooling`, `speech_to_text`,
`features`, `reasoning`, `tool_calling`, `applications`, `rl`,
`deployment`, `ray_serving`, `disaggregated`, `observability` [EV-001].

The page itself contains no inline code. Each category resolves to a
GitHub directory under `github.com/vllm-project/vllm/tree/main/examples/`,
and the example files live next to vLLM's source tree, not under
its docs [EV-002]. The page is a catalogue of pointers; the
examples are versioned with the code.

The taxonomy operates on two dimensions simultaneously, presented
as a single flat list [EV-003]:

- **Capability dimension** — what the engine can do: `generate`
  (text generation, including multimodal), `pooling` (embedding,
  classification, scoring, reward), `speech_to_text` (transcription,
  translation, real-time audio), `reasoning`, `tool_calling`.
- **Deployment dimension** — how the engine is operated:
  `basic` (minimal offline + serving), `deployment` (production),
  `ray_serving` (scalable serving via Ray), `disaggregated` (KV
  cache connectors, failure recovery), `observability` (metrics,
  logging, tracing, dashboards).
- **Crosscut categories** — `features` (per-capability vLLM features
  like prefix caching, speculative decoding, LoRA), `applications`
  (chatbots, RAG), `rl` (reinforcement learning).

The page contains no numbered workflows, no procedural recipes,
no "if you want X, follow steps 1..N" framing [EV-004]. Its
intro reads, verbatim, *"vLLM's examples are organized into the
following categories"* — taxonomy presented as taxonomy, no
narration about which entry-point a new user should pick [EV-005].

---

## Why This KB Cares

ctx's recipe surface at `docs/recipes/` is the structural opposite
of vLLM's: a small set of named procedural workflows (`activating-context.md`,
`build-a-knowledge-base.md`, etc.), each containing prose, command
sequences, and rationale. The contrast surfaces three choices
ctx has made implicitly that are worth examining:

1. **Prose vs index.** ctx assumes the reader needs guidance through
   a workflow; vLLM assumes the reader needs only a pointer to
   versioned code. Different audiences (ctx: humans + AI agents
   reading recipes as instructions; vLLM: ML engineers reading
   examples as templates to adapt).
2. **Docs-tree vs code-tree.** ctx's recipes live in `docs/recipes/`,
   physically separated from `internal/`. vLLM's examples live
   in `examples/` next to source, with the docs site just indexing
   them. Both approaches have a trade-off: docs-tree lets prose
   evolve independently of code; code-tree keeps examples
   compile-fresh and version-pinned.
3. **Taxonomy depth.** vLLM ships thirteen categories despite the
   page lacking any guided entry point; the breadth IS the value
   when each entry is itself a directory of working code. ctx
   currently has a flatter recipe surface; if it grows beyond ten
   recipes, the question becomes whether to keep the flat list or
   adopt a vLLM-style category split.

The kb tracks this topic because future ctx design discussions
about recipe surface, doc/code separation, and example versioning
will benefit from having a worked contrastive reference. It also
seeds the broader pattern of *studying adjacent AI-infrastructure
tools' doc strategies* as a recurring kb activity.

---

## Sources and Further Reading

- `VLLM-EXAMPLES`: vLLM's examples landing page on the public docs
  site; canonical index of category pointers. Source recorded in
  `../../source-map.md`; rows extracted on 2026-05-21 in
  `../../evidence-index.md`.

This pass ingested only the landing page itself, not the per-category
GitHub directories. Per-category deep dives (e.g. what does
`deployment/` actually demonstrate beyond a directory name?) are
deferred. A future pass with discovery enabled could follow the
thirteen GitHub links and mint per-category EV rows; that work
would land as sub-pages under this topic, not as a single mega-page.

---

## Related Concepts in This KB

`none surfaced` — this is the first topic in this kb. The
topic-adjacency pre-flight ran against an empty
`source-coverage.md` and found nothing to surface. Future passes
that ingest comparable doc structures (Claude Code's recipe
surface, zensical's example tree, GitNexus's docs layout) should
be cross-linked here.
