# cli-command-design

## [2026-07-04-212434] Probing whether a ctx subcommand exists needs bare invocation, not --help

**Context**: During the PR #128 review I checked ctx command existence with 'ctx <cmd> --help' and got false positives — cobra short-circuits --help to exit 0 even for an unknown subcommand, so 'ctx system stats --help' appeared to succeed for a command that does not exist. Separately, probing with a real mutating verb ('ctx task complete 3') actually completed a live task in TASKS.md, which had to be reverted.

**Lesson**: To detect CLI command drift, invoke the BARE subcommand and check the real exit code (unknown => exit 1); --help lies because it is a persistent flag cobra handles before dispatch. Never probe with a mutating verb against the live project — ctx subcommands write to .context/. Also: 'ctx system <unknown>' emits the 'relay VERBATIM' version-skew NudgeBox and exits 1, so a wrapper that treats non-empty stdout as success will render that notice as a normal result.

**Application**: When auditing whether an editor/wrapper's ctx invocations are valid, probe each bare form for its exit code (against a throwaway repo, or use read-only forms only); treat any wrapper whose 'non-empty output => success' logic ignores exit codes as buggy.

---

## [2026-05-28-201300] cobra's legacyArgs lets unknown subcommands silently succeed on non-root groups

**Context**: Every prompt of this session injected 52 lines of `ctx system` help
text into agent context, labeled "hook success." Investigation traced it to the
0.8.1 plugin's `hooks.json` wiring `ctx system check-anchor-drift` as the first
UserPromptSubmit hook — a command the 0.8.1 binary no longer has (the command
was deleted by the cwd-anchored migration in `fc7db228`, but the plugin's hook
config wasn't updated). The harness reported "hook success" because cobra exits
0 on the unknown subcommand.

**Lesson**: cobra's `legacyArgs` only raises "unknown command" for the **root**
command (`!cmd.HasParent()`); any non-root group (built with `parent.Cmd`)
treats an unknown subcommand as non-error: it falls through to `Help()` and
returns nil → exit 0. In a UserPromptSubmit hook this is **invisible** — the
harness logs "hook success" and injects the whole help text into agent context
every prompt. The 0.8.1 plugin's stale wiring of the retired
`check-anchor-drift` caused exactly this for the entire session.

**Application**: Non-root cobra groups must have an explicit unknown-subcommand
guard. Two routes: (a) `Args: cobra.NoArgs` so unknown subcommands error loud
(non-zero exit + "unknown command" stderr); (b) a `RunE` that emits a **verbatim
relay** — which is what actually reaches the user in a UserPromptSubmit hook
context where a non-zero exit alone is invisible. Tracked under Phase CLI-FIX as
the verbatim-relay guard on `ctx system`.

---

## [2026-05-20-214821] /ctx-plan is named after its input, not its output

**Context**: Agent (and apparently other agents in prior sessions per user
observation) repeatedly inverted the canonical chain, treating /ctx-spec as the
entry point and /ctx-plan as a post-spec step. The skill description starts
'stress-test a plan' (implying user brings a plan IN) while line 44 of the body
says 'the deliverable is a debated brief, not a task list' (the OUTPUT is a
brief, not a plan).

**Lesson**: Skill names that reference their INPUT bias the agent toward the
wrong canonical position. The /ctx-plan skill takes a plan and produces a brief;
the natural mental model when scanning the name is 'plan = output', which makes
the agent place it AFTER spec instead of before. Also: /ctx-spec's 'When to Use'
section listed /ctx-brainstorm as a predecessor but never /ctx-plan, so an agent
skimming the top of the skill never learned the full chain.

**Application**: Made the canonical chain explicit at the top of both /ctx-plan
and /ctx-spec skills (Canonical Chain block with the brainstorm → plan →
spec → implement diagram) and in AGENT_PLAYBOOK_GATE Planning Work section.
/ctx-spec When-to-Use now lists /ctx-plan as a predecessor; When-NOT-to-Use says
'when the bet is contested but not yet stress-tested, use /ctx-plan first'.
/ctx-plan description now ends with '; produces a debated brief at
.context/briefs/<TS>-<slug>.md that /ctx-spec --brief consumes'.

---

