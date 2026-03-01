---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "The Thesis"
icon: lucide/scroll-text
---

# Context as State

![ctx](../images/ctx-banner.png)

## A Persistence Layer for Human-AI Cognition

*Jose Alekhinne* - <jose@ctx.ist>

*February 2026*

### Abstract

As AI tools evolve from code-completion utilities into reasoning collaborators, 
the knowledge that governs their behavior becomes as important as the code they 
produce; yet, that knowledge is routinely discarded at the end of every session.

AI-assisted development systems assemble context at prompt time using heuristic 
retrieval from mutable sources: recent files, semantic search results, 
session history. These approaches optimize relevance at the moment of generation 
but do not persist the cognitive state that produced decisions. Reasoning is 
not reproducible, intent is lost across sessions, and teams cannot audit the 
knowledge that constrains automated behavior.

This paper argues that context should be treated as deterministic,
version-controlled state rather than as a transient query result. We ground this 
argument in three sources of evidence: a landscape analysis of 17 systems 
spanning AI coding assistants, agent frameworks, and knowledge stores; 
a taxonomy of five primitive categories that reveals irrecoverable architectural 
trade-offs; and an experience report from [`ctx`](https://github.com/ActiveMemory/ctx), 
a persistence layer for AI-assisted development, which developed itself using its 
own persistence model across 389 sessions over 33 days. We define a three-tier 
model for cognitive state: authoritative knowledge, delivery views, and 
ephemeral state. Then we present six design invariants empirically validated 
by 56 independent rejection decisions observed across the analyzed landscape. 
We show that context determinism applies to assembly, not to model output, 
and that the curation cost this model requires is offset by compounding returns 
in reproducibility, auditability, and team cognition.

---

## 1. Introduction

The introduction of large language models into software development has shifted 
the primary interface from code execution to interactive reasoning. In this 
environment, the correctness of an output depends not only on source code but 
on the context supplied to the model: the conventions, decisions, architectural 
constraints, and domain knowledge that bound the space of acceptable responses.

Current systems treat context as a query result assembled at the moment of 
interaction. A developer begins a session; the tool retrieves what it estimates 
to be relevant from chat history, recent files, and vector stores; the model 
generates output conditioned on this transient assembly; the session ends, and 
the context evaporates. The next session begins the cycle again.

This model has improved substantially over the past year. `CLAUDE.md` files, 
Cursor rules, Copilot's memory system, and tools such as Mem0, Letta, and Kindex 
each address aspects of the persistence problem. Yet across 17 systems we 
analyzed spanning AI coding assistants, agent frameworks, autonomous coding 
agents, and purpose-built knowledge stores, no system provides all five of the 
following properties simultaneously: deterministic context assembly, 
human-readable file-based persistence, token-budgeted delivery, 
zero runtime dependencies, and local-first operation.

This paper does not propose a universal replacement for retrieval-centric 
workflows. It defines a persistence layer (embodied in `ctx` (https://ctx.ist)) 
whose advantages emerge under specific operational conditions: when 
reproducibility is a requirement, when knowledge must outlive sessions and 
individuals, when teams require shared cognitive authority, or when offline 
operation is necessary. 

The trade-offs (manual curation cost, reduced automatic recall, coarser 
granularity) are intentional and mirror the trade-offs accepted by systems that 
favor reproducibility over convenience, such as reproducible builds and 
immutable infrastructure [^1] [^6].

The contribution is threefold: a three-tier model for cognitive state that 
resolves the ambiguity between authoritative knowledge and ephemeral session 
artifacts; six design invariants empirically grounded in a cross-system 
landscape analysis; and an experience report demonstrating that the model 
produces compounding returns when applied to its own development.

---

## 2. The Limits of Prompt-Time Context

Prompt-time assembly pipelines typically consist of corpus selection, retrieval,
ranking, and truncation. These pipelines are probabilistic and time-dependent,
producing three failure modes that compound over the lifetime of a project.

### 2.1 Non-Reproducibility

If context is derived from mutable sources using heuristic ranking, identical
requests at different times receive different inputs. A developer who asks 
"What is our authentication strategy?" on Tuesday may receive a different 
context window than the same question on Thursday: Not because the strategy 
changed, but because the retrieval heuristic surfaced different fragments.

Reproducibility (the ability to reconstruct the exact inputs that produced a 
given output) is a foundational property of reliable systems. Its loss in 
AI-assisted development mirrors the historical evolution from ad-hoc builds to 
deterministic build systems [^1] [^2]. The build community learned that when 
outputs depend on implicit state (environment variables, system clocks, 
network-fetched dependencies), debugging becomes archaeology. The same principle 
applies when AI outputs depend on non-deterministic context retrieval.

### 2.2 Opaque Knowledge

Embedding-based memory increases recall but reduces inspectability. When a 
vector store determines that a code snippet is "similar" to the current query, 
the ranking function is opaque: the developer cannot inspect why that snippet 
was chosen, whether a more relevant artifact was excluded, or whether the 
ranking will remain stable. This prevents deterministic debugging, policy 
auditing, and causal attribution (properties that information retrieval theory 
identifies as fundamental trade-offs of probabilistic ranking) [^3].

In practice, this opacity manifests as a compliance ceiling. In our experience 
developing a context management system (detailed in Section 7), soft instructions 
(directives that ask an AI agent to read specific files or follow specific 
procedures) achieve approximately 75-85% compliance. The remaining 15-25% 
represents cases where the agent exercises judgment about whether the 
instruction applies, effectively applying a second ranking function on top of 
the explicit directive. When 100% compliance is required, instruction is 
insufficient; the content must be injected directly, removing the agent's option
to skip it.

### 2.3 Loss of Intent

Session transcripts record interaction but not cognition. A transcript captures
what was said but not which assumptions were accepted, which alternatives were 
rejected, or which constraints governed the decision. The distinction matters: 
a decision to use PostgreSQL recorded as a one-line note ("Use PostgreSQL") 
teaches a model what was decided; a structured record with context, rationale, 
and consequences teaches it why (and why is what prevents the model from 
unknowingly reversing the decision in a future session) [^4].

Session transcripts provide history. Cognitive state requires something more: 
the persistent, structured representation of the knowledge required for correct
decision-making.

---

## 3. Cognitive State: A Three-Tier Model

### 3.1 Definitions

We define **cognitive state** as the authoritative, persistent representation of
the knowledge required for correct decision-making within a project. It is
human-authored or human-ratified, versioned, inspectable, and reproducible. It 
is distinct from logs, transcripts, retrieval results, and model-generated 
summaries.

Previous formulations of this idea have treated cognitive state as a monolithic 
concept. In practice, a three-tier model better captures the operational reality:

**Tier 1: Authoritative State**: The canonical knowledge that the system treats 
as ground truth. In a concrete implementation, this corresponds to a set of 
human-curated files with defined schemas: a constitution (inviolable rules), 
conventions (code patterns), an architecture document (system structure), 
decision records (choices with rationale), learnings (captured experience), a 
task list (current work), a glossary (domain terminology), and an agent 
playbook (operating instructions). Each file has a single purpose, a defined 
lifecycle, and a distinct update frequency. Authoritative state is 
version-controlled alongside code and reviewed through the same mechanisms 
(diffs, pull requests, blame annotations).

**Tier 2: Delivery Views**: Derived representations of authoritative state, 
assembled for consumption by a model. A delivery view is produced by a 
deterministic assembly function that takes the authoritative state, a token 
budget, and an inclusion policy as inputs and produces a context window as 
output. The same authoritative state, budget, and policy must always produce the 
same delivery view. Delivery views are ephemeral (they exist only for the 
duration of a session), but their construction is reproducible.

**Tier 3: Ephemeral State**: Session transcripts, scratchpad notes, draft 
journal entries, and other artifacts that exist during or immediately after a 
session but are not authoritative. Ephemeral state is the raw material from 
which authoritative state may be extracted through human review, but it is 
never consumed directly by the assembly function.

This three-tier model resolves confusion present in earlier formulations:
the claim that AI output is a deterministic function of the repository state. 
The corrected claim is that **context selection** is deterministic (the delivery 
view is a function of authoritative state), but **model output** remains 
stochastic, conditioned on the deterministic context. Formally:

```
delivery_view = assemble(authoritative_state, budget, policy)
output = model(delivery_view)   # stochastic
```

The persistence layer's contribution is making `assemble` reproducible, not
making `model` deterministic.

### 3.2 Separation of Concerns

The decision to separate authoritative state into distinct files with distinct
purposes is not cosmetic. Different types of knowledge have different lifecycles:

| Knowledge Type | Update Frequency | Read Frequency   | Load Priority | Example                                                  |
|----------------|------------------|------------------|---------------|----------------------------------------------------------|
| Constitution   | Rarely           | Every session    | Always        | "Never commit secrets to git"                            |
| Tasks          | Every session    | Session start    | Always        | "Implement token budget CLI flag"                        |
| Conventions    | Weekly           | Before coding    | High          | "All errors use structured logging with severity levels" |
| Decisions      | When decided     | When questioning | Medium        | "Use PostgreSQL over MySQL (see ADR-003)"                |
| Learnings      | When learned     | When stuck       | Medium        | "Hook scripts >50ms degrade interactive UX"              |
| Architecture   | When changed     | When designing   | On demand     | "Three-layer pipeline: ingest → enrich → assemble"       |
| Journal        | Every session    | Rarely           | Never auto    | "Session 247: Removed dead-end session copy layer"       |

A monolithic context file would force the assembly function to load everything 
or nothing. Separation enables **progressive disclosure**: the minimum context 
that matters for the current moment, with the option to load more when needed. 
A normal session loads the constitution, tasks, and conventions; a deep 
investigation loads decision history and journal entries from specific dates.

The budget mechanism is the constraint that makes separation valuable. Without a 
budget, the default behavior is to load everything, which destroys the attention 
density that makes loaded context useful. With a budget, the assembly function 
must prioritize ruthlessly: constitution first (always full), then tasks and 
conventions (budget-capped), then decisions and learnings (scored by recency). 
Entries that do not fit receive title-only summaries rather than being silently 
dropped (an application of the "tell me what you don't know" pattern identified 
independently by four systems in our landscape analysis).

---

## 4. Design Invariants

The following six invariants define the constraints that a cognitive state 
persistence layer must satisfy. They are not axioms chosen a priori; they are 
empirically grounded properties whose violation was independently identified as 
producing complexity costs across the 17 systems we analyzed.

### Invariant 1: Markdown-on-Filesystem Persistence

Context files must be human-readable, git-diffable, and editable with any text
editor. No database. No binary storage.

*Validation*: 11 independent rejection decisions across the analyzed landscape 
protected this property. Systems that adopted embedded records, binary 
serialization, or knowledge graphs as their core primitive consistently traded 
away the ability for a developer to run `cat DECISIONS.md` and understand the 
system's knowledge. The inspection cost of opaque storage compounds over the 
lifetime of a project: every debugging session, every audit, every onboarding 
conversation requires specialized tooling to access knowledge that could have 
been a text file.

### Invariant 2: Zero Runtime Dependencies

The tool must work with no installed runtimes, no running services, and no API
keys for core functionality.

*Validation*: 13 independent rejection decisions protected this property  
(the most frequently defended invariant). Systems that required databases 
(PostgreSQL, SQLite, Redis), embedding models, server daemons, container 
runtimes, or cloud APIs for core operation introduced failure modes proportional 
to their dependency count. A persistence layer that depends on infrastructure is 
not a persistence layer; it is a service. Services have uptime requirements, 
version compatibility matrices, and operational costs that simple file 
operations do not.

### Invariant 3: Deterministic Context Assembly

The same files plus the same budget must produce the same output. No 
embedding-based retrieval, no LLM-driven selection, no wall-clock-dependent 
scoring in the assembly path.

*Validation*: 6 independent rejection decisions protected this property. 
Non-deterministic assembly (whether from embedding variance, LLM-based selection, 
or time-dependent scoring) destroys the ability to reproduce a context window 
and therefore to diagnose why a model produced a given output. Determinism in 
the assembly path is what makes the persistence layer auditable.

### Invariant 4: Human Authority Over Persistent State

The agent may propose changes to context files but must not unilaterally modify 
them. All persistent changes go through human-reviewable git commits.

*Validation*: 6 independent rejection decisions protected this property. Systems 
that allowed agents to self-modify their memory (writing freeform notes, 
auto-pruning old entries, generating summaries as ground truth) consistently 
produced lower-quality persistent context than systems that enforced human 
review. Structure is a feature, not a limitation: across the landscape, the 
pattern "structured beats freeform" was independently discovered by four 
systems that evolved from freeform LLM summaries to typed schemas with required 
fields.

### Invariant 5: Local-First, Air-Gap Capable

Core functionality must work offline with no network access. Cloud services may
be used for optional features but never for core context management.

*Validation*: 7 independent rejection decisions protected this property. 
Infrastructure-dependent memory systems cannot operate in classified 
environments, isolated networks, or disaster-recovery scenarios. A 
filesystem-native model continues to function under all conditions where the 
repository is accessible.

### Invariant 6: No Default Telemetry

Any analytics, if ever added, must be strictly opt-in.

*Validation*: 4 independent rejection decisions protected this property. Default 
telemetry erodes the trust model that a persistence layer depends on. If 
developers must trust the system with their architectural decisions, operational 
learnings, and project constraints, the system cannot simultaneously be reporting 
usage data to external services.

These six invariants collectively define a design space. Each feature proposal 
can be evaluated against them: a feature that violates any invariant is rejected 
regardless of how many other systems implement it. The discipline of constraint 
(refusing to add capabilities that compromise foundational properties) is 
itself an architectural contribution. Across the 17 analyzed systems, 56
patterns were explicitly rejected for violating these invariants. The rejection 
count per invariant (11, 13, 6, 6, 7, 4) provides a rough measure of each 
property's vulnerability to architectural erosion. A representative sample of 
these rejections is provided in Appendix A.[^1]

[^1]: The 56 figure counts patterns explicitly rejected for invariant violation 
or on other architectural grounds. The full analysis catalog contains additional 
entries classified as "watch" (adopted provisionally or deferred for future 
evaluation) or "adapt" (partially adopted with modifications) rather than 
"reject." Only entries where the pattern was conclusively rejected are included 
in the count.

---

## 5. Landscape Analysis

The 17 systems were selected to cover the architectural design space rather than 
to achieve completeness. Each included system satisfies three criteria: it 
represents a distinct architectural primitive for AI-assisted development, it is 
actively maintained or widely referenced, and it provides sufficient public 
documentation or source code for architectural inspection. The goal was to 
ensure that every major category of primitive (document, embedded record, state 
snapshot, event/message, construction/derivation) was represented by multiple 
systems, enabling cross-system pattern detection.

The resulting set spans six categories: AI coding assistants (Continue, 
Sourcegraph/Cody, Aider, Claude Code), AI agent frameworks (CrewAI, AutoGen, 
LangGraph, LlamaIndex, Letta/MemGPT), autonomous coding agents (OpenHands, 
Sweep), session provenance tools (Entire), data versioning systems (Dolt, 
Pachyderm), pipeline/build systems (Dagger), and purpose-built knowledge 
stores (QubicDB, Kindex). Each system was analyzed from its source code and 
documentation, producing 34 individual analysis artifacts (an architectural 
profile and a set of insights per system) that yielded 87 adopt/adapt 
recommendations, 56 explicit rejection decisions, and 52 watch items.

### 5.1 Primitive Taxonomy

Every system in the AI-assisted development landscape operates on a core 
primitive: an atomic unit around which the entire architecture revolves. Our 
analysis of 17 systems reveals five categories of primitives, each making 
irrecoverable trade-offs:

**Group A: Document/File Primitives**: Human-readable documents as the primary 
unit. Documents are authored by humans, version-controlled in git, and consumed 
by AI tools. The invariant of this group is that the primitive is always 
human-readable and version-controllable with standard tools. Three systems 
participate in this pattern: the system described in this paper as a pure 
expression, and Continue (via its rules directory) and Claude Code 
(via `CLAUDE.md` files) as partial participants: both use document-based 
context as an input but organize around different core primitives.

**Group B: Embedded Record Primitives**: Vector-embedded records stored with 
numerical embeddings for similarity search, metadata for filtering, and scoring 
mechanisms for ranking. Five systems use this approach 
(LlamaIndex, CrewAI, Letta/MemGPT, QubicDB, Kindex). The invariant is that the 
primitive requires an embedding model or vector database for core operations: 
a dependency that precludes offline and air-gapped use.

**Group C: State Snapshot Primitives**: Point-in-time captures of the complete 
system state. The invariant is that any past state can be reconstructed at any 
historical point. Three systems use this approach (LangGraph, Entire, Dolt).

**Group D: Event/Message Primitives**: Sequential events or messages forming an 
append-only log with causal relationships. Four systems use this approach 
(OpenHands, AutoGen, Claude Code, Sweep). The invariant is temporal ordering 
and append-only semantics.

**Group E: Construction/Derivation Primitives**: Derived or constructed values 
that encode how they were produced. The invariant is that the primitive is a 
function of its inputs; re-executing the same inputs produces the same 
primitive. Three systems use this approach (Dagger, Pachyderm, Aider).

### 5.2 Comparison Matrix

The five primitive categories differ along seven dimensions:

| Property              | Document | Embedded Record | State Snapshot | Event/Message | Construction |
|-----------------------|----------|-----------------|----------------|---------------|--------------|
| Human-readable        | Yes      | No              | Varies         | Partially     | No           |
| Version-controllable  | Yes      | No              | Varies         | Yes           | Yes          |
| Queryable by meaning  | No       | Yes             | No             | No            | No           |
| Rewindable            | Via git  | No              | Yes            | Yes (replay)  | Yes          |
| Deterministic         | Yes      | No              | Yes            | Yes           | Yes          |
| Zero-dependency       | Yes      | No              | Varies         | Varies        | Varies       |
| Offline-capable       | Yes      | No              | Varies         | Varies        | Yes          |

The document primitive is the only one that simultaneously satisfies 
human-readability, version-controllability, determinism, zero dependencies, and 
offline capability. This is not because documents are superior in general
(embedded records provide semantic queryability that documents lack) but because 
the combination of all five properties is what the persistence layer requires. 
The choice between primitive categories is not a matter of capability but of 
which properties are considered invariant.

### 5.3 Convergent Patterns

Across the 17 analyzed systems, six design patterns were independently 
discovered. These convergent patterns carry extra validation weight because 
they emerged from different problem spaces:

**Pattern 1: "Tell me what you don't know"**: When context is incomplete, 
explicitly communicate to the model what information is missing and what 
confidence level the provided context represents. Four systems independently 
converged on this pattern: inserting skip markers, tracking evidence gaps, 
annotating provenance, or naming output quality tiers.

**Pattern 2: "Freshness matters"**: Information relevance decreases over time. 
Three systems independently chose exponential decay with different half-lives 
(30 days, 90 days, and LRU ordering). Static priority ordering with no time 
dimension leaves relevant recent knowledge at the same priority as stale 
entries. This pattern is in productive tension with the persistence model's 
emphasis on determinism: the claim is not that time-dependence is irrelevant, 
but that it belongs in the curation step (a human deciding to consolidate or 
archive stale entries) rather than in the assembly function (an algorithm 
silently down-ranking entries based on age).

**Pattern 3: "Content-address everything"**: Compute a hash of content at 
creation time for deduplication, cache invalidation, integrity verification, 
and change detection. Five systems independently implement content hashing, 
each discovering it solves different problems [^5].

**Pattern 4: "Structured beats freeform"**: When capturing knowledge or session 
state, a structured schema with required fields produces more useful data than 
freeform text. Four systems evolved from freeform summaries to typed schemas:
one moving from LLM-generated prose to a structured condenser with explicit 
fields for completed tasks, pending tasks, and files modified.

**Pattern 5: "Protocol convergence"**: The Model Context Protocol (MCP) is 
emerging as a standard tool integration layer. Nine of 17 systems support it, 
spanning every category in the analysis. MCP's significance for the persistence 
model is that it provides a transport mechanism for context delivery without 
dictating how context is stored or assembled. This makes the approach compatible 
with both retrieval-centric and persistence-centric architectures.

**Pattern 6: "Human-in-the-loop for memory"**: Critical memory decisions should 
involve human judgment. Fully automated memory management produces lower-quality 
persistent context than human-reviewed systems. Four systems independently 
converged on variants of this pattern: ceremony-based consolidation, 
interrupt/resume for human input, confirmation mode for high-risk actions, and 
separated "think fast" vs. "think slow" processing paths.

Pattern 6 directly validates the ceremony model described in this paper. The 
persistence layer requires human curation not because automation is impossible, 
but because the quality of persistent knowledge degrades when the curation step 
is removed. The improvement opportunity is to make curation easier, not to 
automate it away.

---

## 6. Worked Example: Architectural Decision Under Two Models

We now instantiate the three-tier model in a concrete system  (`ctx`) and 
illustrate the difference between prompt-time retrieval and cognitive state 
persistence using a real scenario from its development.

### 6.1 The Problem

During development, the system accumulated three overlapping storage layers for 
session data: raw transcripts (owned by the AI tool), session copies 
(JSONL copies plus context snapshots), and enriched journal entries 
(Markdown summaries). The middle layer (session copies) was a dead-end write 
sink. An auto-save hook copied transcripts to a directory that nothing read 
from, because the journal pipeline already read directly from the raw 
transcripts. Approximately 15 source files, a shell hook, 20 configuration 
constants, and 30 documentation references supported infrastructure with 
no consumers.

### 6.2 Prompt-Time Retrieval Model

In a retrieval-based system, the decision to remove the middle layer depends on 
whether the retrieval function surfaces the relevant context:

The developer asks: "Should we simplify the session storage?" The retrieval 
system must find and rank the original discussion thread where the three layers 
were designed, the usage statistics showing zero reads from the middle layer, 
the journal pipeline documentation showing it reads from raw transcripts 
directly, and the dependency analysis showing 15 files, a hook, and 30 doc 
references. If any of these fragments are not retrieved (because they are in old 
chat history, because the embedding similarity score is low, or because the 
token budget was consumed by more recent but less relevant context), the model 
may recommend preserving the middle layer, or may not realize it exists.

Six months later, a new team member asks the same question. The retrieval results 
will differ: the original discussion has aged out of recency scoring, the usage 
statistics are no longer in recent history, and the model may re-derive the 
answer or arrive at a different conclusion.

### 6.3 Cognitive State Model

In the persistence model, the decision is recorded as a structured artifact at 
write time:

```markdown
## [2026-02-11] Remove .context/sessions/ storage layer

**Status**: Accepted

**Context**: The session/recall/journal system had three overlapping
storage layers. The recall pipeline reads directly from raw transcripts,
making .context/sessions/ a dead-end write sink that nothing reads from.

**Decision**: Remove .context/sessions/ entirely. Two stores remain:
raw transcripts (global, tool-owned) and enriched journal
(project-local).

**Rationale**: Dead-end write sinks waste code surface, maintenance
effort, and user attention. The recall pipeline already proved that
reading directly from raw transcripts is sufficient. Context snapshots
are redundant with git history.

**Consequences**: Deleted internal/cli/session/ (15 files), removed
auto-save hook, removed --auto-save from watch, removed pre-compact
auto-save, removed /ctx-save skill, updated ~45 documentation files.
Four earlier decisions superseded.
```

This artifact is:

* **Deterministically included** in every subsequent session's delivery view
  (budget permitting, with title-only fallback if budget is exceeded)
* **Human-readable** and reviewable as a diff in the commit that introduced it
* **Permanent**: it persists in version control regardless of retrieval heuristics
* **Causally linked**: it explicitly supersedes four earlier decisions, 
  creating an auditable chain

When the new team member asks "Why don't we store session copies?" six months 
later, the answer is the same artifact, at the same revision, with the same 
rationale. The reasoning is reconstructible because it was persisted at write 
time, not discovered at query time.

### 6.4 The Diff When Policy Changes

If a future requirement re-introduces session storage (for example, to 
support multi-agent session correlation), the change appears as a diff to the 
decision record:

```diff
- **Status**: Accepted
+ **Status**: Superseded by [2026-08-15] Reintroduce session storage
+ for multi-agent correlation
```

The new decision record references the old one, creating a chain of reasoning
visible in `git log`. In the retrieval model, the old decision would simply be 
ranked lower over time and eventually forgotten.

---

## 7. Experience Report: A System That Designed Itself

The persistence model described in this paper was developed and tested by using 
it on its own development. Over 33 days and 389 sessions, the system's context 
files accumulated a detailed record of decisions made, reversed, and 
consolidated: providing quantitative and qualitative evidence for the model's 
properties.

### 7.1 Scale and Structure

The development produced the following authoritative state artifacts:

* 8 consolidated decision records covering 24 original decisions spanning 
  context injection architecture, hook design, task management, security, 
  agent autonomy, and webhook systems
* 18 consolidated learning records covering 75 original observations spanning 
  agent compliance, hook behavior, testing patterns, documentation drift, and 
  tool integration
* A constitution with 13 inviolable rules across 4 categories 
  (security, quality, process, context preservation)
* 389 enriched journal entries providing a complete session-level audit trail

The consolidation ratio (24 decisions compressed to 8 records, 75 learnings 
compressed to 18) illustrates the curation cost and its return: authoritative 
state becomes denser and more useful over time as related entries are merged, 
contradictions are resolved, and superseded decisions are marked.

### 7.2 Architectural Reversals

Three architectural reversals during development provide evidence that the 
persistence model captures and communicates reasoning effectively:

**Reversal 1: The two-tier persistence model**: The original design included a 
middle storage tier for session copies. After 21 days of development, the middle 
tier was identified as a dead-end write sink (described in Section 6). The 
decision record captured the full context, and the removal was executed cleanly: 
15 source files, a shell hook, and 45 documentation references. The pattern of 
a "dead-end write sink" was subsequently observed in 7 of 17 systems in our 
landscape analysis that store raw transcripts alongside structured context.

**Reversal 2: The prompt-coach hook**: An early design included a hook that 
analyzed user prompts and offered improvement suggestions. After deployment, 
the hook produced zero useful tips, its output channel was invisible to users, 
and it accumulated orphan temporary files. The hook was removed, and the 
decision record captured the failure mode for future reference.

**Reversal 3: The soft-instruction compliance model**: The original context 
injection strategy relied on soft instructions: directives asking the AI agent 
to read specific files. After measuring compliance across multiple sessions, 
we found a consistent 75-85% compliance ceiling. The revised strategy 
injects content directly, bypassing the agent's judgment about whether to 
comply. The learning record captures the ceiling measurement and the rationale 
for the architectural change.

Each reversal was captured as a structured decision record with context, 
rationale, and consequences. In a retrieval-based system, these reversals would 
exist only in chat history, discoverable only if the retrieval function happens 
to surface them. In the persistence model, they are permanent, indexable 
artifacts that inform future decisions.

### 7.3 Compliance Ceiling

The 75-85% compliance ceiling for soft instructions is the most operationally 
significant finding from the experience report. It means that any context 
management strategy relying on agent compliance with instructions ("read this 
file," "follow this convention," "check this list") has a hard ceiling on 
reliability.

The root cause is structural: the instruction "don't apply judgment" is itself 
evaluated by judgment. When an agent receives a directive to read a file, it 
first assesses whether the directive is relevant to the current task (and that 
assessment is the judgment the directive was trying to prevent).

The architectural response maps directly to the formal model defined in 
Section 3.1. Content requiring 100% compliance is included in 
`authoritative_state` and injected by the deterministic `assemble` function, 
bypassing the agent entirely. Content where 80% compliance is acceptable is 
delivered as instructions within the delivery view. The three-tier architecture 
makes this distinction explicit: authoritative state is injected; delivery 
views are assembled deterministically; ephemeral state is available but 
not pushed.

### 7.4 Compounding Returns

Over 33 days, we observed a qualitative shift in the development experience. 
Early sessions (days 1-7) spent significant time re-establishing context:
explaining conventions, re-stating constraints, re-deriving past decisions. 
Later sessions (days 25-33) began with the agent loading curated context and 
immediately operating within established constraints, because the constraints 
were in files rather than in chat history.

This compounding effect (where each session's context curation improves all 
subsequent sessions) is the primary return on the curation investment. The cost 
is borne once (writing a decision record, capturing a learning, updating the 
task list); the benefit is collected on every subsequent session load.

The effect is analogous to compound interest in financial systems: the 
knowledge base grows not linearly with effort but with increasing marginal 
returns as new knowledge interacts with existing context. A learning captured 
on day 5 prevents a mistake on day 12, which avoids a debugging session that 
would have consumed a day 12 session, freeing that session for productive work 
that generates new learnings. The growth is not literally exponential (it is 
bounded by project scope and subject to diminishing returns as the knowledge 
base matures), but within the observed 33-day window, the returns were 
consistently accelerating.

### 7.5 Scope and Generalizability

This experience report is self-referential by design: the system was developed 
using its own persistence model. This circularity strengthens the internal 
validity of the findings (the model was stress-tested under authentic 
conditions) but limits external generalizability. The two-week crossover point 
was observed on a single project of moderate complexity with a small team 
already familiar with the model's assumptions. Whether the same crossover holds 
for larger teams, for codebases with different characteristics, or for teams 
adopting the model without having designed it remains an open empirical 
question. The quantitative claims in this section should be read as 
existence proofs (demonstrating that the model *can* produce compounding 
returns) rather than as predictions about specific adoption scenarios.

---

## 8. Situating the Persistence Layer

The persistence layer occupies a specific position in the stack of AI-assisted 
development:

```
Application Logic
AI Interaction / Agents
Context Retrieval Systems
Cognitive State Persistence Layer
Version Control / Storage
```

Current systems innovate primarily in the retrieval layer (improving how 
context is discovered, ranked, and delivered at query time). The persistence 
layer sits beneath retrieval and above version control. Its role is to maintain 
the authoritative state that retrieval systems may query but do not own. The 
relationship is complementary: retrieval answers "What in the corpus might be 
relevant?"; cognitive state answers "What must be true for this system to 
operate correctly?" A mature system uses both: retrieval for discovery, 
persistence for authority.

---

## 9. Applicability and Trade-Offs

### 9.1 When to Use This Model

A cognitive state persistence layer is most appropriate when:

**Reproducibility is a requirement**: If a system must be able to answer "Why 
did this output occur, and can it be produced again?" then deterministic, 
version-controlled context becomes necessary. This is relevant in regulated 
environments, safety-critical systems, long-lived infrastructure, and 
security-sensitive deployments.

**Knowledge must outlive sessions and individuals**: Projects with multi-year 
lifetimes accumulate architectural decisions, domain interpretations, and 
operational policy. If this knowledge is stored only in chat history, issue 
trackers, and institutional memory, it decays. The persistence model converts 
implicit knowledge into branchable, reviewable artifacts.

**Teams require shared cognitive authority**: In collaborative environments, 
correctness depends on a stable answer to "What does the system believe to be 
true?" When this answer is derived from retrieval heuristics, authority shifts 
to ranking algorithms. When it is versioned and human-readable, authority 
remains with the team.

**Offline or air-gapped operation is required**: Infrastructure-dependent 
memory systems cannot operate in classified environments, isolated networks, 
or disaster-recovery scenarios.

### 9.2 When Not to Use This Model

**Zero-configuration personal workflows**: For short-lived or exploratory tasks, 
the cost of explicit knowledge curation outweighs its benefits. Heuristic 
retrieval is sufficient when correctness is non-critical, outputs are 
disposable, and historical reconstruction is unnecessary.

**Maximum automatic recall from large corpora**: Vector retrieval systems 
provide superior performance when the primary task is searching vast, weakly 
structured information spaces. The persistence model assumes that what matters 
can be decided and that this decision is valuable to record.

**Fully autonomous agent architectures**: Agent runtimes that generate and 
discard state continuously, optimizing for local goal completion, do not benefit 
from a model that centers human ratification of knowledge.

### 9.3 Incremental Adoption

The transition does not require full system replacement. An incremental path:

**Step 1: Record decisions as versioned artifacts**: Instead of allowing 
conclusions to remain in discussion threads, persist them in reviewable form 
with context, rationale, and consequences [^4]. This alone converts ephemeral 
reasoning into the cognitive state.

**Step 2: Make inclusion deterministic**: Define explicit assembly rules. 
Retrieval may still exist, but it is no longer authoritative.

**Step 3: Move policy into cognitive state**: When system behavior depends on 
stable constraints, encode those constraints as versioned knowledge. Behavior 
becomes reproducible.

**Step 4: Optimize assembly, not retrieval**: Once the authoritative layer 
exists, performance improvements come from budgeting, caching, and structural 
refinement rather than from improving ranking heuristics.

### 9.4 The Curation Cost

The primary objection to this model is the cost of explicit knowledge curation. 
This cost is real. Writing a structured decision record takes longer than 
letting a chatbot auto-summarize a conversation. Maintaining a glossary requires 
discipline. Consolidating 75 learnings into 18 records requires judgment.

The response is not that the cost is negligible but that it is amortized. 
A decision record written once is loaded hundreds of times. A learning captured 
today prevents repeated mistakes across all future sessions. The curation cost 
is paid once; the benefit compounds.

The experience report provides rough order-of-magnitude numbers. Across 389 
sessions over 33 days, curation activities (writing decision records, 
capturing learnings, updating the task list, consolidating entries) averaged 
approximately 3-5 minutes per session. In early sessions (days 1-7), before 
curated context existed, re-establishing context consumed approximately 10-15 
minutes per session: re-explaining conventions, re-stating architectural 
constraints, re-deriving decisions that had been made but not persisted. By the 
final week (days 25-33), the re-explanation overhead had dropped to near zero: 
the agent loaded curated context and began productive work immediately.

At ~12 sessions per day, the curation cost was roughly 35-60 minutes daily. 
The re-explanation cost in the first week was roughly 120-180 minutes daily. 
By the third week, that cost had fallen to under 15 minutes daily while the 
curation cost remained stable. The crossover (where cumulative curation cost 
was exceeded by cumulative time saved) occurred around day 10. These figures are 
approximate and derived from a single project with a small team already familiar 
with the model; the crossover point will vary with project complexity, 
team size, and curation discipline.

---

## 10. Future Work

Several directions are compatible with the model described here:

**Section-level deterministic budgeting**: Current assembly operates at file 
granularity. Section-level budgeting would allow finer-grained control (including 
specific decision records while excluding others within the same file) without 
sacrificing determinism.

**Causal links between decisions**: The experience report shows that decisions 
frequently reference earlier decisions (superseding, extending, or qualifying 
them). Formal causal links would enable traversal of the decision graph and 
automatic detection of orphaned or contradictory constraints.

**Content-addressed context caches**: Five systems in our landscape analysis 
independently discovered that content hashing provides cache invalidation, 
integrity verification, and change detection. Applying content addressing to the 
assembly output would enable efficient cache reuse when the authoritative state 
has not changed.

**Conditional context inclusion**: Five systems independently suggest that 
context entries could carry activation conditions (file patterns, task 
keywords, or explicit triggers) that control whether they are included in a 
given assembly. This would reduce the per-session budget cost of large 
knowledge bases without sacrificing determinism.

**Provenance metadata**: Linking context entries to the sessions, decisions, or 
learnings that motivated them would strengthen the audit trail. Optional 
provenance fields on Markdown entries (session identifier, cause reference, 
motivation) would be lightweight and compatible with the existing file-based 
model.

---

## 11. Conclusion

AI-assisted development has treated context as a "query result" assembled at the 
moment of interaction, discarded at the session end. This paper identifies a 
complementary layer: the persistence of authoritative cognitive state as 
deterministic, version-controlled artifacts.

The contribution is grounded in three sources of evidence. A landscape analysis 
of 17 systems reveals five categories of primitives and shows that no existing 
system provides the combination of human-readability, determinism, zero 
dependencies, and offline capability that the persistence layer requires. Six 
design invariants, validated by 56 independent rejection decisions, define the 
constraints of the design space. An experience report over 389 sessions and 
33 days demonstrates compounding returns: later sessions start faster, 
decisions are not re-derived, and architectural reversals are captured with 
full context.

The core claim is this: persistent cognitive state enables causal reasoning 
across time. A system built on this model can explain not only *what* is true, 
but *why it became true* and *when it changed*.

When context is the state:

* Reasoning is reproducible: the same authoritative state, budget, and policy 
  produce the same delivery view.
* Knowledge is auditable: decisions are traceable to explicit artifacts with 
  context, rationale, and consequences.
* Understanding compounds: each session's curation improves all subsequent 
  sessions.

The choice between retrieval-centric workflows and a persistence layer is not a 
matter of capability but of time horizon. Retrieval optimizes for relevance at 
the moment of interaction. Persistence optimizes for the durability of 
understanding across the lifetime of a project.

---

🐸🖤 *"Gooood... let the deterministic context flow through the repository..."*<br>
- Kermit the Sidious, probably

---

## Appendix A: Representative Rejection Decisions

The 56 rejection decisions referenced in Section 4 were cataloged across all 17 
system analyses, grouped by the invariant they would violate. This appendix 
provides a representative sample (two per invariant) to illustrate the 
methodology.

**Invariant 1: Markdown-on-Filesystem (11 rejections)**: CrewAI's vector 
embedding storage was rejected because embeddings are not human-readable, not 
git-diff-friendly, and require external services. Kindex's knowledge graph as 
core primitive was rejected because it requires specialized commands to inspect 
content that could be a text file (`kin show <id>` vs. `cat DECISIONS.md`).

**Invariant 2: Zero Runtime Dependencies (13 rejections)**: Letta/MemGPT's 
PostgreSQL-backed architecture was rejected because it conflicts with 
local-first, no-database, single-binary operation. Pachyderm's Kubernetes-based 
distributed architecture was rejected as the antithesis of a single-binary 
design for a tool that manages text files.

**Invariant 3: Deterministic Assembly (6 rejections)**: LlamaIndex's 
embedding-based retrieval as the primary selection mechanism was rejected because 
it destroys determinism, requires an embedding model, and removes human 
judgment from the selection process. QubicDB's wall-clock-dependent scoring was 
rejected because it directly conflicts with the "same inputs produce same 
output" property.

**Invariant 4: Human Authority (6 rejections)**: Letta/MemGPT's agent 
self-modification of memory was rejected as fundamentally opposed to 
human-curated persistence. Claude Code's unstructured auto-memory (where the 
agent writes freeform notes) was rejected because structured files with defined 
schemas produce higher-quality persistent context than unconstrained agent 
output.

**Invariant 5: Local-First / Air-Gap Capable (7 rejections)**: Sweep's 
cloud-dependent architecture was rejected as fundamentally incompatible with 
the local-first, offline-capable model. LangGraph's managed cloud deployment 
was rejected because cloud dependencies for core functionality violate 
air-gap capability.

**Invariant 6: No Default Telemetry (4 rejections)**: Continue's 
telemetry-by-default (PostHog) was rejected because it contradicts the 
local-first, privacy-respecting trust model. CrewAI's global telemetry on 
import (Scarf tracking pixel) was rejected because it violates user trust and 
breaks air-gap capability.

The remaining 9 rejections did not map to a specific invariant but were 
rejected on other architectural grounds: for example, Aider's 
full-file-content-in-context approach (which defeats token budgeting), 
AutoGen's multi-agent orchestration as core primitive (scope creep), 
and Claude Code's 30-day transcript retention limit 
(institutional knowledge should have no automatic expiration).

---

## References

[^1]: Reproducible Builds Project, *"Reproducible Builds: Increasing the Integrity of Software Supply Chains"*, 2017. <https://reproducible-builds.org/docs/definition/>
[^2]: S. McIntosh et al., *"The Impact of Build System Evolution on Software Quality"*, ICSE, 2015. <https://doi.org/10.1109/ICSE.2015.70>
[^3]: C. Manning, P. Raghavan, H. Schütze, *Introduction to Information Retrieval*, Cambridge University Press, 2008. <https://nlp.stanford.edu/IR-book/>
[^4]: M. Nygard, *"Documenting Architecture Decisions"*, Cognitect Blog, 2011. <https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions>
[^5]: L. Torvalds et al., *Git Internals - Git Objects* (content-addressed storage concepts). <https://git-scm.com/book/en/v2/Git-Internals-Git-Objects>
[^6]: Kief Morris, *Infrastructure as Code*, O'Reilly, 2016.
[^7]: J. Kreps, *"The Log: What every software engineer should know about real-time data's unifying abstraction"*, 2013. <https://engineering.linkedin.com/distributed-systems/log>
[^8]: P. Hunt et al., *"ZooKeeper: Wait-free coordination for Internet-scale systems"*, USENIX ATC, 2010. <https://www.usenix.org/legacy/event/atc10/tech/full_papers/Hunt.pdf>
