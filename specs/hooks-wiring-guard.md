# Shipped-Hooks Wiring Guard

The shipped `hooks.json` and the on-PATH `ctx` binary must agree
on the command surface. Every `ctx <…>` invocation wired into the
shipped hooks must resolve to a registered subcommand on the same
binary. A half-migrated package — hooks calling a command the
binary no longer ships — must fail at `go test` time, not in a
user's session.

## Problem

`ctx system <unknown>` prints the `system` group's ~51-line Long
help and exits **0** (cobra's `legacyArgs` raises "unknown
command" only for the root command, never a non-root group). In a
`UserPromptSubmit` hook, exit 0 is read by the harness as
"hook success", so the entire help blob is injected into the
agent's context **on every prompt**.

This is not hypothetical. It shipped:

- The cwd-anchored migration (`fc7db228`, spec
  `specs/cwd-anchored-context.md`) deleted the
  `check-anchor-drift` hook command — under cwd-anchoring the
  drift it detected cannot occur.
- That migration is unreleased (no `0.8.x` tag past `v0.8.0`),
  but a plugin package was published whose bundled `hooks.json`
  still wired `ctx system check-anchor-drift` first in the
  `UserPromptSubmit` stack.
- Result on a machine running the post-migration binary against
  the pre-migration plugin: ~51 lines of `system` help injected
  every prompt, labelled "hook success". Confirmed live during
  session 0066d49b.

The working-tree asset (`internal/assets/claude/hooks/hooks.json`)
is already correct (`cd`-based, no `check-anchor-drift`). The
defect was a packaging skew, and nothing in the build caught it
because no check cross-references the shipped wiring against the
binary's actual command tree.

## Invariant

For every `ctx <token…>` invocation appearing in a `command`
string in `internal/assets/claude/hooks/hooks.json`, walking the
leading subcommand tokens from the root of the assembled cobra
command tree (`bootstrap.Initialize(bootstrap.RootCmd())`) must
land on a registered command at each step. Hidden commands count
as registered — cobra traverses them; only help display hides
them.

Direction is one-way: **wired ⊆ registered**. Registered
commands that no hook wires are fine. The guard only fails when a
hook names a verb the binary does not have.

## Guard

A Go test in `internal/compliance` (sibling to
`TestShippedHooksExcludeCheckAudit`, which already reads the same
asset):

1. Read and JSON-decode the shipped `hooks.json`.
2. Walk every event → group → hook `command` string.
3. From each command string, extract every `ctx` invocation: the
   token `ctx` followed by its run of leading subcommand-shaped
   tokens (`^[a-z][a-z0-9-]*$`), stopping at the first flag,
   redirection, or shell operator. `ctx agent --budget 8000`
   yields the path `[agent]`; `ctx system check-version` yields
   `[system check-version]`.
4. Build the full command tree and, for each extracted path,
   descend token by token matching `Command.Name()` (or an
   alias). A token with no matching child at its level fails the
   test, naming the hook command and the unresolved token.

A Go test, not a `hack/release.sh` step: it is cross-platform,
runs in CI with the rest of the suite, needs no bash, and lives
next to the binary it validates.

## Live fix (out of band)

The guard prevents recurrence; it does not unstick the already
published skewed package. That requires cutting/republishing a
release whose bundled `hooks.json` and binary come from the same
post-`fc7db228` commit, then reinstalling for users on the skewed
build. That is a release action (tag + publish), owned by the
maintainer, not part of this commit.

## Non-goals

- Hardening `parent.Cmd` so `ctx system <bogus>` exits non-zero
  instead of dumping help (spec's Bug #2). That is a separate
  task with its own spec and branch; this guard catches the
  packaging skew regardless of how the binary behaves on an
  unknown subcommand.
- Validating flags or flag values wired in hooks — only the
  command path is checked.
