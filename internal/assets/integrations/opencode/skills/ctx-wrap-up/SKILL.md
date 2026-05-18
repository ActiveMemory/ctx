---
name: ctx-wrap-up
description: "End-of-session context persistence ceremony. Use when wrapping up a session to capture learnings, decisions, conventions, and tasks."
---

Run the end-of-session context persistence ceremony.

## When to Use

- When ending a work session
- When switching to a different project or task area
- When context window is getting large
- Before any long break from the project

## Process

1. Review work done in this session
2. Capture any new decisions to `.context/DECISIONS.md`
3. Capture any new learnings to `.context/LEARNINGS.md`
4. Capture any new conventions to `.context/CONVENTIONS.md`
5. Update task status in `.context/TASKS.md`
6. Save a session summary to `.context/sessions/`
7. If `.context/kb/` exists, list pending closeouts under
   `.context/ingest/closeouts/` and count `open` rows in
   `.context/kb/outstanding-questions.md`; surface both
   counts so the operator sees what editorial residue is
   pending. The handover step's fold pass consumes the
   closeouts. Skip this step entirely when `.context/kb/`
   does not exist.
8. **Mandatory final step:** delegate to `/ctx-handover` so
   the next session has something to read. Draft:
   - **Title**: short noun phrase naming the session arc
     (becomes the slug in `<TS>-<slug>.md`).
   - **`--summary`** (past tense, one paragraph): what was
     done this session. Concrete, not vague.
   - **`--next`** (future tense, one paragraph): the
     specific first action the next agent should take.
   - **`--highlights`**: bullet list of notable artifacts.
     Always draft; pass empty only after the user
     explicitly says there is nothing to highlight.
   - **`--open-questions`**: bullet list of things that
     remain undecided. Always draft; pass empty only after
     explicit user confirmation.

   Confirm the draft with the user, then invoke:

   ```text
   /ctx-handover "<title>" --summary "<...>" --next "<...>" \
     [--highlights "<...>"] [--open-questions "<...>"]
   ```

   Do not declare the wrap-up complete until
   `.context/handovers/<TS>-<slug>.md` has been written.

## Self-Check

Ask: "If this session ended right now, would the next session
know what happened?" If no, persist more context before
ending. The handover step is the floor: without it, recall
degenerates to probabilistic reconstruction from canonical
files plus journal.
