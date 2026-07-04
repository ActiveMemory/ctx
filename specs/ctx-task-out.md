# ctx-task-out — Close the Spec→Implement Gap

**Status:** Implemented (session 1306853f, 2026-07-03) — skill at
`internal/assets/claude/skills/ctx-task-out/SKILL.md`, amended after
two review rounds (role-authority split, execution ledger, amendment
mode, files column, matrix task refs, hard-gate refinements), plus a
third round from the first real consumer (zhc/os m0a run, 2026-07-03):
the template's task-table column list omitted the `st` state cell its
own ledger paragraph required (template beats prose — the generated
table was unauditable), and step 7's epic anchors lacked partition and
completion semantics (overlapping id ranges, no rule for when an epic
checks off). Both fixed: `st` column (`[ ]`/`[x]`/`[o]`) in template,
step 4, ledger rule, and checklist; step 7 now requires a disjoint
id partition with a sum check and a stated completion rule.
**Draft skill:** retired after landing (was
`specs/ctx-task-out-SKILL-draft.md`; the installed skill superseded it).
**Origin:** zhc/os session c887a6d4, 2026-07-03 — the gap surfaced in live use.

## Problem

The canonical chain has an unowned artifact. `/ctx-plan` explicitly
disclaims implementation planning ("the deliverable is a debated brief,
not a task list"). `/ctx-spec` produces the committed what/why and its
own guidance keeps specs concise ("a page is usually enough").
`/ctx-implement` opens with "use when you have a plan document" — and
**nothing in the chain produces that document.**

For single-session features the spec doubles as the plan and nothing is
missing. For multi-milestone specs the conciseness rule and the
implementable-alone rule (`specs/README.md`) become contradictory, and
the missing mass — data model, contracts, test matrix, granular tasks
with acceptance criteria — has no home. In live use this collapsed into
ad-hoc tasking: ten coarse TASKS.md entries with no acceptance criteria
and a blocking open question (language choice) sitting under task #1.
Every artifact with a ceremony came out well; the unceremonied step is
where quality collapsed.

SpecKit comparison, for orientation: ctx has an adversarial step SpecKit
lacks (`/ctx-plan`); SpecKit has the `/plan` + `/tasks` ladder ctx lacks.
This skill adds the ladder without importing a second framework (which
would create dual sources of truth for specs and tasks).

## Naming

**`ctx-task-out`.** Phrasal verb ("task this out"), precedent in
`ctx-wrap-up`, family resonance with `ctx-task-add`. Rejected:
`ctx-tasks` (noun; breaks the verb-suffix convention), `ctx-task`
(reads as noun, collides conceptually with `ctx-task-add`),
`ctx-decompose` (verb but loses the task-family naming),
`ctx-breakdown` (noun). If the implementing agent finds a strong reason
to rename, keep the phrasal-verb constraint.

## Deliverables

1. **The skill.** Install `specs/ctx-task-out-SKILL-draft.md` as
   `internal/assets/claude/skills/ctx-task-out/SKILL.md`. Mirror to the
   copilot-cli integration (see how `ctx-spec` is carried under
   `internal/assets/integrations/copilot-cli/skills/`). Register
   wherever skills are enumerated (plugin manifest, tests —
   `internal/assets/templates_test.go` may assert counts). Follow
   `CONTRIBUTING-SKILLS.md`.
2. **Chain diagrams.** The 4-step chain appears in multiple skills
   (`ctx-plan`, `ctx-spec`, `ctx-implement`, possibly `ctx-brainstorm`,
   plus copilot mirrors, docs/site). Find all:
   `grep -rn "ctx-spec  →  /ctx-implement" internal/ docs/ site/`
   Replace with the 5-step chain:
   `/ctx-brainstorm → /ctx-plan → /ctx-spec → /ctx-task-out → /ctx-implement`
   with stage labels `(vague) (contested) (committed) (decomposed) (execution)`.
3. **`ctx-spec` amendment — the delegation hook.** Add a final process
   step after "6. Cross-Reference":

   > ### 7. Hand off to tasking
   > If the spec spans multiple milestones or more than ~one session of
   > implementation, do not stop at coarse task creation: recommend
   > `/ctx-task-out --spec specs/<name>.md --milestone <first>` and say
   > why (specs stay concise; the plan carries decomposition). For
   > small specs, suggest `/ctx-implement` directly.

   Also add one line to the `--brief` flow (step 5) so the handoff
   fires in non-interactive mode too.
4. **`ctx-implement` amendment.** Its "use when you have a plan
   document" wording should name `specs/plans/<milestone>.md` (the
   `/ctx-task-out` output) as the canonical input while still accepting
   hand-written plans. Add a redirect: if invoked with only a
   multi-milestone spec, suggest `/ctx-task-out` first.
5. **Project template.** `internal/assets/project/specs-README.md`
   lifecycle gains a step between Draft/Reference and Implement:
   "Task out: for multi-milestone specs, `/ctx-task-out` writes
   `specs/plans/<milestone>.md`; TASKS.md carries epic anchors only."

## Design constraints (do not weaken)

- **Hard gates, not warnings.** The blocking-TBD gate and the
  rolling-wave gate refuse; they do not proceed-with-caveats. The
  chain's signature move is refusing to run ahead of its inputs
  (`/ctx-spec` refuses vague ideas; this refuses unresolved specs).
- **One source of truth for tasks.** The plan owns fine-grained tasks;
  TASKS.md gets epic-level anchors with a `Plan:` reference, one-way
  sync. This is the argument against adopting SpecKit wholesale — do
  not reintroduce dual truth here.
- **Decomposer, not designer.** The skill must not relitigate the spec;
  disagreements route back to `/ctx-plan`.
- **Falsifiable acceptance criteria** on every task — a command, a
  test, or an observable. "Implement X" alone is a checklist failure.

## Acceptance criteria for this work order

- [x] Skill installed + copilot parity + registered; the asset tests
      pass (full `go test ./...` green)
- [x] All chain diagrams show 5 steps (grep returns no 4-step remnants)
- [x] `/ctx-spec` ends with the tasking handoff in both interactive and
      `--brief` flows
- [x] `/ctx-implement` names the plan artifact and redirects when handed
      a bare multi-milestone spec (and refuses a `Status: Blocked` plan)
- [x] A dry run against a real spec (e.g. zhc/os `specs/v1-substrate.md`,
      milestone m0a) produces a plan passing the skill's own quality
      checklist — including a correct refusal while D-001 (language
      choice) is unresolved
      *(refusal half executed live 2026-07-03 by session c887a6d4
      hand-running the landed SKILL.md against zhc/os m0a: the gate
      surfaced TWO blockers — D-001 and D-003 (artifact schema
      language), one more than this work order predicted; the
      "schema format" example in the blocking definition earned its
      place. Both blockers resolved same day via zhc/os DECISIONS.md
      (D-003 → JSON Schema, steering; D-001 → Go, Board). The
      plan-producing run happens post-reinstall.)*
      *(plan-producing half executed 2026-07-03 by zhc/os session
      a63353a3 with the registered skill: gates behaved — one
      genuinely blocking TBD (config mechanism) was resolved via
      DECISIONS.md before decomposition, two measurement-gated TBDs
      were carried into the plan; output: 32 tasks, 15-row matrix at
      zhc/os specs/plans/m0a.md. Board audit then found the two
      skill gaps described in Status (stateless task table;
      non-partitioned epic anchors) — plan amended downstream,
      skill fixed here. The skill's own checklist passed while the
      output was unauditable, which is exactly why the checklist
      gained the two new boxes.)*

## Review

After landing, the originating session's agent (zhc/os, session
c887a6d4) reviews: gate behavior first (TBD refusal, rolling-wave
refusal), then plan-document quality on the m0a dry run, then chain-text
consistency. First real consumer: zhc/os `specs/plans/m0a.md`.
