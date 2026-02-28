---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Invariants
icon: lucide/anchor
---

![ctx](../images/ctx-banner.png)

# The System Explains Itself

These are the properties that **must hold** for any valid `ctx` implementation.

* These are **not features**.
* These are **constraints**.

A change that violates an invariant is a **category error**, 
*not* an improvement.

---

## Cognitive State Tiers

`ctx` distinguishes between three forms of state:

* **Authoritative state**: Versioned, inspectable artifacts that define intent 
  and survive time.
* **Delivery views**: Deterministic assemblies of the authoritative state for a 
  specific budget or workflow.
* **Ephemeral working state**: Local, transient, or sensitive data that 
  assists interaction but does not define system truth.

The invariants below apply primarily to the **authoritative cognitive state**.

---

## 1. Cognitive State Is Explicit

All authoritative context lives in artifacts that can be **inspected**,
**reviewed**, and **versioned**.

If something is important, it **must** exist as a file: Not only in a prompt, 
a chat, or a model's hidden memory.

---

## 2. Assembly Is Reproducible

Given the same:

* repository state,
* configuration,
* and inputs,

context assembly produces the same result.

Heuristics may *rank* or *filter* for delivery under constraints.

They **do not** alter the authoritative state.

---

## 3. The Authoritative State Is Human-Readable

The authoritative cognitive state **must** be stored in formats that a human can:

- **read**,
- **diff**,
- **review**,
- and **edit** directly.

Sensitive working memory **may** be encrypted at rest.
However, encryption **must not** become the only representation of 
authoritative knowledge.

---

## 4. Artifacts Outlive Sessions

Sessions are *transient*.

**Knowledge persists**.

Reasoning, decisions, and outcomes **must remain available** after the 
interaction that produced them has ended.

---

## 5. Authority Is User-Defined

What enters the authoritative context is an **explicit human decision**.

Models may suggest.

Automation may assist.

**Selection is never implicit**.

---

## 6. Operation Is Local-First

The core system must function without requiring network access or a 
remote service.

External systems **may** extend `ctx`.

They **must not** be required for its operation.

---

## 7. Versioning Is the Memory Model

The evolution of the authoritative cognitive state must be:

* **preserved**,
* **inspectable**,
* and **branchable**.

Ephemeral and sensitive working state may use different retention and diff 
strategies by design.

**Understanding includes understanding how we arrived here**.

---

## 8. Structure Enables Scale

Unstructured accumulation is **not** memory.

Authoritative cognitive state must have a defined layout that:

* **communicates** intent,
* **supports** navigation,
* and **prevents** drift.

---

## 9. Verification Is the Scoreboard

Claims without recorded outcomes are noise.

**Reality** (*observed and captured*) is the **only** signal that **compounds**.

This invariant defines a required direction:

**The authoritative state must be able to record expectation and result**.

---

## 10. Capture Once, Reuse Indefinitely

Work that has already produced understanding **must not** be re-derived from 
scratch.

Explored paths, rejected options, and validated conclusions **are** 
permanent assets.

---

## 11. Policies Are Encoded, not Remembered

Alignment **must not** depend on recall or goodwill.

Constraints that matter **must** exist in machine-readable form and
**participate** in context assembly.

---

## 12. The System Explains Itself

From the repository state alone it must be possible to determine:

- **what** was authoritative,
- **what** constraints applied.

Delivery views may be optimized.

They must not become the only explanation.

---

# Non-Goals

To avoid category errors, `ctx` does **not** attempt to be:

* a skill,
* a prompt management tool,
* a chat history viewer,
* an autonomous agent runtime,
* a vector database,
* a hosted memory service.

Such systems **may** integrate with `ctx`.

They **do not** define it.

---

# Implications for Contributions

Valid contributions:

* **strengthen** an invariant,
* **reduce** the cost of maintaining an invariant,
* or **extend** the system without violating invariants.

Invalid contributions:

* introduce hidden authoritative state,
* replace reproducible assembly with non-reproducible behavior,
* make core operation depend on external services,
* reduce human inspectability of authoritative state,
* or bypass explicit user authority over what becomes authoritative.

---

# The Contract

Everything else (*commands, skills, layouts, integrations, optimizations*) 
is an implementation detail.

**These invariants are the system**.

