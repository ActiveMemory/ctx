This analysis will be done for the following projects:

We will keep a checklist and do this ONE BY ONE to avoid overwhelming the
agent and also keeping the context window "within reason".

This is a large deep-dive and the agent will be given a 1M-token window
for each analysis; that said, it does not have to consume the entire window.
Once a satisfactory result is established, it can end the process and
inform the user.

- https://github.com/entireio/cli
- https://github.com/entireio/claude-plugins
- https://github.com/qubicDB/qubicdb
- https://github.com/Aider-AI/aider
- https://github.com/continuedev/continue
- https://github.com/OpenDevin/OpenDevin
- https://github.com/anthropics/claude-code
- https://github.com/sweepai/sweep
- https://github.com/sourcegraph/sourcegraph
- https://github.com/letta-ai/letta
- https://github.com/run-llama/llama_index
- https://github.com/dagger/dagger
- https://github.com/dolthub/dolt
- https://github.com/pachyderm/pachyderm
- https://github.com/langchain-ai/langgraph
- https://github.com/joaomdmoura/crewai
- https://github.com/microsoft/autogen

## How this benefits our specific trajectory

Given how we are positioning **ctx as substrate / cognitive continuity infra**,
this plan:

* turns “we are different” into a **mechanical proof**
* produces **roadmap inputs**, not blog content
* gives you **competitive language grounded in primitives**

## Methodology

Study projects in **four strategic clusters**, prioritized by relevance
to ctx's problem space:

### Cluster 1 — Context Assembly (highest priority)

These systems solve "what information goes into the prompt?" — closest
to ctx's core problem.

- Continue
- Sourcegraph
- LlamaIndex
- Entire (cli + claude-plugins)

### Cluster 2 — Memory Systems

These systems solve "how does an agent remember across sessions?" —
directly comparable to ctx's persistence model.

- Letta
- Aider
- CrewAI
- AutoGen
- LangGraph

### Cluster 3 — Session Provenance

These systems solve "what happened and why?" — relevant to ctx's
journal, recall, and archaeology features.

- Claude Code
- OpenDevin
- Sweep

### Cluster 4 — Artifact / Versioned Data (lowest priority)

These systems version structured data — architecturally distant but
may reveal primitives ctx can learn from.

- Dolt
- Pachyderm
- Dagger
- QubicDB

**Execution order**: Cluster 1 → 2 → 3 → 4. Within each cluster,
order is flexible.

The goal is to extract **primitives and invariants**, not (necessarily) features.



------- PLAN MARKDOWN BELOW ----

# Comparative Analysis Plan: ctx-Adjacent Systems

## Purpose

Perform a deep, structured, and repeatable comparative analysis of multiple 
open-source projects that operate in problem spaces adjacent to `ctx`.

The goal is to extract:

1. Conceptual differences in problem framing
2. Architectural patterns
3. Workflow integration models
4. Context selection / provenance strategies
5. Opportunities for `ctx` inspiration (not imitation)

This is **not** a feature checklist exercise.  
This is a **systems and primitives** analysis.

## Repository Acquisition Protocol

All target repositories MUST be cloned into an isolated workspace.

### Workspace Layout

/home/jose/WORKSPACE/ctx-analysis/
repos/
<project-name>/
artifacts/
logs/

### Clone Rules

- Use shallow clone unless full history is required:

  git clone --depth=1 <repo>

- If commit is specified:

  git checkout <commit-sha>

- Initialize submodules if present:

  git submodule update --init --recursive

- If the target link is a GitHub organization, search relevant repositories
  in the organization, and clone all of them, and if they are more than 
  a handful, clone a representative set.

### Immutability

Once cloned:

- Do NOT pull
- Do NOT switch branches
- Treat repo as read-only

This guarantees reproducibility of analysis.

### Reuse

If the repo already exists locally:

- DO NOT re-clone
- Reuse the existing directory

---

## Inputs

### Target Repository

The agent will receive a list that can grow over time.

Each target should be treated independently and then normalized into a shared 
comparison model.

<target repositories here: note that a project may contain more than
 one repository, such as a CLI, engine, standards, specs, service, etc.>

---

## Outputs

For each project create:

```

analysis/<project-name>/
├── PROJECT.md
├── ARCHITECTURE.md
├── WORKFLOWS.md
├── CONTEXT_MODEL.md
├── EXTENSIBILITY.md
├── TRUST_PROVENANCE.md
├── CODE_QUALITY.md
├── TASK_BAKEOFF.md
└── INSIGHTS_FOR_CTX.md

```

And a global:

```

analysis/
└── COMPARISON.md

```

All outputs must be **artifact-grade**, not chat summaries.

---

## Core Comparison Axes

### 1. Problem Framing

Identify the system’s primary primitive:

Examples:
- checkpoint
- context assembly endpoint
- session graph
- memory store
- embedding index
- event log

Document:

- What is the core object?
- What problem does the project believe it solves?
- Explicit non-goals

---

### 2. Architecture

Produce a high-level architecture diagram and describe:

- data model
- storage strategy
- control plane vs data plane
- runtime dependencies
- extension points
- plugin system (if any)

Focus on **load-bearing modules**, not file counts.

---

### 3. Workflow Integration

How the system fits into real developer flow:

- CLI
- Git hooks
- CI/CD
- IDE integration
- MCP / API surface
- local-first vs cloud

Answer:

> When does a developer actually *feel* this system?

---

### 4. Context Strategy

How the system answers:

> “What information is selected and why?”

Look for:

- retrieval algorithms
- token budgeting
- ranking / scoring
- recency vs semantic vs graph expansion
- determinism vs heuristic assembly
- explainability of context

---

### 5. Provenance and Trust

Can the system answer:

- where did this come from?
- why was it chosen?
- can it be replayed?

Check:

- versioning
- audit trail
- reproducibility
- redaction / privacy model

---

### 6. Extensibility Model

- plugin mechanism
- API surface
- scripting / automation
- internal vs external extension

---

### 7. Codebase Signals

High-level quality indicators:

- test strategy
- release cadence
- issue hygiene
- documentation depth
- migration story

---

## Thin-Slice Code Analysis Strategy

Do **not** read the entire repository sequentially.

Instead:

### Step 1 — Promise Layer

- README
- docs
- examples

Extract claimed system model.

### Step 2 — Primitive Discovery

Find where the core object is:

- created
- stored
- queried
- transformed

Trace one full lifecycle.

### Step 3 — Happy Path Execution Flow

Follow one real workflow end-to-end.

### Step 4 — Failure Modes

What breaks?

- partial data
- large data
- missing services
- merge conflicts
- concurrent usage

---

## Task Bake-Off (Applied Comparison)

Each system must be evaluated against the same tasks:

### T1 — Archaeology

“Explain why this change happened.”

### T2 — Budgeted Context Assembly

“Prepare context for fixing a bug under a strict token limit.”

### T3 — Session Replay

“Reconstruct prior work deterministically.”

### T4 — Onboarding

“New contributor becomes productive in 15 minutes.”

### T5 — Offline Operation

“System works in an air-gapped environment.”

Document friction and strengths.

---

## Scoring Model

Use qualitative scoring:

- Native
- Possible with effort
- Against system design

---

## Insights for ctx

For each project produce:

### Adopt

Ideas that align directly with ctx philosophy.

### Adapt

Ideas that require transformation.

### Reject

Ideas that conflict with ctx’s core model.

### Watch

Emerging patterns worth monitoring.

Each item must include:

- value
- implementation complexity (S/M/L)
- impact (Low/Med/High)

---

## COMPARISON.md Structure

1. Primitive taxonomy
2. Architectural patterns across systems
3. Context strategy spectrum
4. Provenance capability spectrum
5. Workflow integration spectrum
6. Competitive positioning of ctx
7. Roadmap opportunities

---

## Session Cadence

Each cluster runs in its own context window. After the last project
in a cluster is analyzed:

1. Review the cluster's artifacts together for cross-cutting patterns.
2. **Stop and tell the user:**

   > "Cluster N is complete. All artifacts are written to disk.
   > Please run /clear to reset the context window before we
   > begin Cluster N+1."

   Do NOT proceed to the next cluster in the same session. Wait for
   the user to confirm they have cleared and are ready.

3. At the start of the next cluster's session, load only the prior
   clusters' `INSIGHTS_FOR_CTX.md` files for continuity — not the
   full analysis corpus.

4. The final `COMPARISON.md` synthesis happens in its own session
   after all four clusters are complete, loading only
   `INSIGHTS_FOR_CTX.md` and `ARCHITECTURE.md` from each project.

---

## Operating Rules

- Prefer code over marketing claims.
- Prefer invariants over feature lists.
- Prefer workflows over components.
- Prefer primitives over abstractions.
- Prefer reproducibility over convenience.

---

## Success Criteria

The analysis is complete when:

- Every system is reduced to its core primitive.
- All systems are mapped onto the same conceptual grid.
- ctx’s differentiation becomes mechanically obvious.
- A concrete inspiration backlog exists.

