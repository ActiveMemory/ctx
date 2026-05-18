---
sha: 60543e46
branch: feat/phase-kb
mode: handover
pass-mode: n/a
life-stage: maintenance
generated-at: 2026-05-18T00:41:18Z
title: phase-kb sentinel-strings + post-amend punch list
---

## Summary

Phase KB branch (`feat/phase-kb` @ `60543e46`) is landed and
linted clean across ~313 files folded into one signed commit.
This session's final iterations did, in order: package
relocation (`initkb` → `initialize/kb`, kb-prefix flat dirs
→ `kb/<sub>` nested, plural→singular kb subdir names,
`topicnew/` flat → `topic/cmd/newcmd/` nested,
`internal/cli/kb/cmd.go` → `kb.go`); function and file
renames (`RunNew` → `Scaffold`, `LatestCursor` → `Latest`,
`CopyKBLanding` → `CopyLanding` for the stutter audit,
`nextid.go` → `next_id.go`, `ctxio` alias → `ctxIo`,
`topics.go` → `topic.go`, `Slugify` removed in favor of
`internal/slug.Path`); shared-package extractions
(`internal/slug/` from `internal/cli/journal/core/slug/`,
`internal/write/kb/row/` from triplicated
contradiction/decision/question append flows,
`internal/cli/setup/core/copilot_cli/github_asset.go` from
the `deployAgent`+`deployInstructions` duplication);
`cmd.go`/`run.go` split across 8 kb+handover subcommands;
Phase KB skill parity ported to copilot-cli and opencode
trees; the
`.context/{TASKS,DECISIONS,LEARNINGS}.md` Phase KB additions
em-dash-swept; markdown→Markdown, en-US, Title-Case,
NDA-name (your-project / your-domain placeholders) and
rot-prone repo-spec-link sweeps run across docs/specs/skills;
`commands.yaml` gained Examples blocks; localizable strings
moved from cfg `messages.go` files into
`commands/text/{errors,write}.yaml` + DescKey constants in
`internal/config/embed/text/`; YAML hierarchy aligned to
`err.kb.<sub>.<verb>` with dot separators; finally, all 13
remaining `messages.go` files were renamed to package-singular
names — no `messages.go` exists anywhere in the repo now.

## Next Session

Fix the `ErrMsg`-as-string-sentinel anti-pattern. Currently
in `internal/config/{handover,closeout,git_meta,kb/cli,
kb/evidence,kb/sourcecoverage,rc,initialize}/<pkg>.go` we
still have things like

```go
ErrMsgMissingGitTree = "git working tree required"
```

These strings back package-level
`var ErrX = errors.New(cfgPkg.ErrMsgX)`. The `errors.Is`
contract uses identity, not text — but the embedded English
string still leaks into `.Error()` output and breaks
localization. The correct shape: (a) the sentinel value
carries identity, not text — use `errors.New("")` or a tiny
typed sentinel with `Is(target error) bool`; (b) the
user-facing text moves into `commands/text/errors.yaml` as
`err.<pkg>.<name>` and is rendered at error-display time by
the `err/<pkg>/<pkg>.go` wrapping constructor (which already
uses `desc.Text` for the format wrapping).

Sweep every `ErrMsg*` in `internal/config/**/*.go`. Verify
with:

```bash
grep -rn 'ErrMsg.*= "' internal/config/
```

After the sweep, also kill any "sentinel mirrors YAML key"
doc comments I left in the cfg files explaining the
duplication — that justification was wrong.

## Highlights

- 8 amends this session culminating at `60543e46`.
- Phase KB skill parity across 3 host trees
  (`claude` / `copilot-cli` / `opencode`).
- Phase RG (`git`-required) and Phase KB (editorial
  pipeline) shipped paired.
- `internal/cli/handover/core/path/` extracted from
  `kb/core/path/` (handover is session-glue, not a KB
  feature; `.context/handovers/<TS>-<slug>.md` files are
  timestamped so concurrent agent runs never overwrite).
- `docs/cli/kb-handover.md` split into `kb.md` + `handover.md`
  (the combined page was a category error).
- `CLAUDE.md` restructured: `## Session Handovers` is now an
  h2 sibling of `## KB Editorial Workflow`, not a sub-step.
- `Phase KB-followup` task filed in `.context/TASKS.md` for
  the adversarial design review of the 3-tree skill drift
  problem (parallel claude/copilot-cli/opencode trees).
- New `internal/slug.Path` for slash-preserving slugs;
  duplicate `Slugify` in `kb/core/topic` removed.
- New `internal/write/kb/row.Append` + `entity.KBRowHooks`
  collapsed three triplicated append flows.
- Doc structure audit clean; cross-package types audit
  clean; mixed-visibility audit clean; gocritic clean;
  stutter audit clean.

## Open Questions

- Should we keep `errors.Is`-style sentinels at all for
  these wrappers, or switch to typed-error structs with
  `Is(target error) bool` methods? Either works; pick one
  and apply uniformly.
- My "`ErrMsg` consts must stay Go-const for package-init
  timing" framing in several recent commit messages and
  doc comments is partially wrong — only the
  `var ErrX = errors.New()` values need a non-empty
  backing string for the `.Error()` fallback, but if every
  caller uses the wrapping constructor (which goes through
  `desc.Text`), the backing string is dead text. Decide
  whether to delete the `ErrMsg` consts entirely (typed
  sentinel with empty `.Error()`), or keep them as a
  degraded-mode fallback when the embedded YAML lookup
  table hasn't been populated yet.
- The `Phase KB-followup` adversarial-review task assumes a
  body-extract builder is the right shape. Confirm before
  building anything.

## Folded Closeouts

none (this is a manual session-end handover; no editorial
closeouts pending).
