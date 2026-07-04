---
name: ctx-experimental-spec
description: "EXPERIMENTAL (discardable). Turn a debated brief into a LOOSE intent spec at .context/specs/intent-<slug>.md — deliberately not pre-shaped into spec-kit's template. Second step of the experimental chain: /ctx-experimental-plan → /ctx-experimental-spec → /ctx-experimental-handoff."
allowed-tools: Read, Edit, Write, Glob, Grep
---

# Write an intent spec from a debrief — experimental

> **Experimental / discardable.** ctx-native port of an external
> intent-spec skill. It deliberately writes a **loose intent spec** in
> a ctx-namespaced location (`.context/specs/`), separate from
> spec-kit's repo-root `specs/<NNN-slug>/`, so the two trees never
> collide. This differs from canonical `/ctx-spec`, which walks the full
> `specs/tpl/spec-template.md` and writes a complete spec to repo-root
> `specs/`. The point of *this* skill is to seed spec-kit, not to
> double-specify. Trial it, then decide what to fold into `/ctx-spec`.

This skill **creates or extends a Markdown intent spec** under
`.context/specs/` from material produced by `/ctx-experimental-plan` (a
debated brief) or **equivalent notes** — meaning a **file path or
user-provided paste** only, not the agent's memory of a prior
conversation. It is **not** a second adversarial interview — it is
**structuring intent for a downstream spec-kit handoff**, without
replacing frozen contracts in `docs/`.

Separation of concerns:

- **`/ctx-experimental-plan`** — stress-test the bet; output is a
  **debated brief** (no milestone task list as the primary deliverable).
- **`/ctx-experimental-spec`** — turn agreed content into a **loose
  intent spec** that `/ctx-experimental-handoff` hands to spec-kit.

## When to use this skill

- User invokes `/ctx-experimental-spec`.
- User says: "write the intent spec", "formalize the debrief", "turn the
  plan into something we can hand to spec-kit."
- Immediately after saving a brief under `.context/briefs/`.

## Input contract (pick one path)

1. **Preferred:** User gives a **path** to an existing brief, e.g.
   `.context/briefs/<TS>-<slug>.md`. **Read that file first** and treat
   it as authoritative source text.

2. **Paste:** User pastes the brief inline. Use their text verbatim
   where it matters; do not invent decisions they did not state.

3. **Missing input:** If neither path nor paste exists, **ask once:**
   "Where is the brief — path under `.context/briefs/`, or paste here?"
   If they have only run `/ctx-experimental-plan` in chat, ask them to
   **save** the brief first — do not fabricate a debrief from memory.

## Authority order (do not invert)

When `docs/`, `DECISIONS.md`, and the brief could conflict, treat
sources in this order:

1. **Frozen contracts** under `docs/` (release notes, public CLI docs)
2. **Recorded decisions** in `.context/DECISIONS.md`
3. **The supplied brief** or pasted notes (this skill's input)
4. **Agent inference** — only when explicitly labeled **`TBD`** in the
   spec body, never as silent fact

Specs are where accidental authority inversion happens; follow this
stack. If the brief contradicts a frozen contract, surface the
contradiction to the user; do not silently follow the brief.

## Procedure

1. **Confirm scope.** One sentence: what feature or problem this spec
   covers. If the brief mixes multiple unrelated bets, ask which
   **single** spec to write first (or split into multiple files with
   explicit user choice).

2. **Name the file deterministically.** Write to
   `.context/specs/intent-<slug>.md`, where `<slug>` is the **same slug**
   as the brief this spec derives from — so the brief
   `<TS>-<slug>.md` and the intent spec `intent-<slug>.md` share one
   identity (the handle `/ctx-experimental-handoff` resolves against).
   **List `.context/specs/`** first (create it if absent); if
   `intent-<slug>.md` already exists you are extending it, not creating a
   second. Only ask the user about the slug when there is no brief to
   inherit it from.

3. **Draft the spec** (adapt sections to the problem — not all sections
   apply every time). A typical shape:

   - `# Intent spec: <title>` — one-line promise of what "done" means for
     *this* doc (not the whole product).
   - **Context / problem** — why this exists (pull from brief).
   - **Goals and non-goals** — in scope / explicitly out.
   - **Proposed approach** — high-level design or behavior (no code dump
     unless the brief already had it).
   - **User-visible contract** — CLI flags, errors, file paths, if
     relevant.
   - **Risks and failure modes** — from the brief; extend from code
     reading **only** when the brief names **concrete** files, commands,
     packages, or existing behavior to inspect. Do **not** perform a
     broad architecture review as part of this skill.
   - **Open questions** — what still needs a decision or spike.
   - **References** — link to the brief path, related `TASKS.md` lines,
     `DECISIONS.md` headings, `docs/` specs if any.

   Keep it **loose**. Inspired by strong specs (problem → solution →
   concepts → edge cases) but shape is flexible; omit sections that add
   noise.

4. **Write the file** at `.context/specs/intent-<slug>.md`. End with a
   **Source** line (path to the brief or `"inline paste from user"`) and
   a single **Handoff** reference line naming the downstream pickup:
   `Handoff: /ctx-experimental-handoff → /speckit-specify`. Keep the spec
   body project-defined — **do not** pre-shape it into spec-kit's
   `spec-template.md` (User Scenarios P1/P2/P3, Functional Requirements,
   Success Criteria). `/speckit-specify` re-derives that structure and
   mints its own REQ-IDs from prose; pre-shaping duplicates its job and
   splits REQ-ID authority. That single Handoff line is the only
   downstream-aware element.

5. **Suggest follow-ups** (do not act without consent): e.g.
   `/ctx-task-add` to add a TASKS.md item pointing at this spec,
   `/ctx-decision-add` if a new architectural commitment was stated, or
   `/ctx-experimental-handoff` to seed `/speckit-specify` from this
   intent spec (optional; it warns and continues if spec-kit is absent).

## Response contract

After writing the spec file, reply with:

- **Path written** (one line).
- **2–4 bullets** summarizing what the spec commits to (not a wall of
  text).
- **Open questions**, if any remain.
- **Suggested next action** — usually `/ctx-experimental-handoff` to seed
  spec-kit, `/ctx-task-add`, or `/ctx-decision-add` when a trade-off
  needs recording.

Do **not** paste the full spec into chat unless the user asks.

## Edge cases

- **No `.context/`:** suggest `ctx init` first; `specs/` is meaningless
  without project memory.
- **Conflict with `docs/`:** if the user wants to change a **frozen**
  contract, say so plainly — `.context/specs/` is for **intent**;
  changing `docs/` may need a recorded decision via `/ctx-decision-add`.
- **User asks for implementation tasks inside the spec:** a short
  **Suggested milestones** subsection is OK if the brief already implies
  order; otherwise keep the spec intent-level and point them to
  `/ctx-task-add`. Suggested milestones must not become the primary
  artifact — if the list would exceed **five** items, stop and recommend
  `/ctx-task-add` for the rest.

## Anti-patterns

- Do not copy-paste proprietary specs from other products verbatim;
  **rephrase** for this repo.
- Do not invent `--flag` names or file paths the brief did not imply;
  mark as **TBD** instead.
- Do not pre-shape the intent spec into spec-kit's template — that is the
  whole reason this skill is separate from `/ctx-spec`.
- Do not treat this file as `DECISIONS.md`; record **approved** trade-offs
  there via `/ctx-decision-add` when the user commits.
- Do not use this skill as permission to **infer** scope, silently
  override `docs/` or `DECISIONS.md`, or **re-open** broad discovery — the
  **Authority order** and **Risks** bounds above are hard limits.
