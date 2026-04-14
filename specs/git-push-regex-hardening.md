# Git Push Regex Hardening

## Problem

The `block-dangerous-command` hook's `MidGitPush` regex only matched
`git push` mid-command after `;`, `&&`, or `||`. This missed:

- Bare `git push` at start of a command
- `git -C /path push` (global flag before subcommand)
- `git --git-dir=/path push`, `--work-tree=...`, `--no-pager`, `--bare`
- Short boolean flags: `git -p push`, `git -P push`
- Env-var prefix: `GIT_DIR=/x git push`
- Command wrappers: `time git push`, `nice git push`
- Subshells and command substitution: `(git push)`, `$(git push)`, backticks
- Multi-line scripts with `\n` between commands

The agent discovered this during a session where `git -C <path> push`
slipped past the permissions deny list (which only matched `git push ...`
at prefix) AND the mid-command regex.

## Approach

Replace the narrow `MidGitPush` regex with a broader `GitPush` that:

1. **Opens the entry anchor** to cover separators, subshells, command
   substitution, and line starts: `(^|[;&|(` + "`" + `\n]\s*)`
2. **Allows arbitrary prefix tokens** before `git`: `(\S+\s+)*` catches
   env-var assignments and command wrappers
3. **Allows arbitrary tokens** between `git` and `push`: `(\s+\S+)*`
   catches any flag shape (short, long, boolean, with or without value)
4. **Tightens the tail anchor** to distinguish `push` (subcommand) from
   ref-name continuations (`push-to-remote`, `push_branch`): uses
   `([^a-zA-Z0-9._/-]|$)` — any non-ref-char or end-of-string

## Trade-offs

- **Accepted false positives**: `git log push` (push as a ref name),
  `git commit -m push` (push as commit message). Indistinguishable from
  real pushes via regex alone. Over-blocking is recoverable;
  under-blocking is not.
- **Known blind spots**: `eval "git push"`, `sh -c "git push"`, shell
  aliases. Parsing through arbitrary quoting is undecidable.

## Non-Goals

- Parsing shell semantics (aliases, `eval`, `source`)
- Making the Claude Code permission deny list obsolete — the regex
  hook is a safety net; the deny list remains the first line
