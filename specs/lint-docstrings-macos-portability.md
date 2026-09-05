# lint-docstrings macOS Portability

## Problem

`make audit` fails on macOS before checking anything:
`hack/lint-docstrings.sh` aborts with ``line 164: unexpected EOF
while looking for matching `'` `` (exit 2). Two portability bugs,
both invisible on Linux CI:

1. The script's shebang is `#!/bin/bash`, which on macOS is bash
   3.2.57. Bash 3.2's command-substitution re-parser treats an
   apostrophe inside a **comment** within `$( … )` as an open
   quote; the comment `# Guard: if sed didn't match …` sits inside
   the big `violations="$({ … })"` capture, so parsing swallows the
   rest of the file and dies at EOF.
2. The struct-field counters use `grep -cP` (PCRE). BSD grep has no
   `-P`, fails silently (stderr discarded), leaves `fieldcount`
   empty, and every 2+-field struct is reported as `MISSING_FIELDS
   … ( fields)` — 59 false positives once bug 1 is fixed.

## Fix

- Reword the comment (`didn't` → `did not`). No apostrophes inside
  `$( … )` comments; bash 3.2 must parse the script.
- Replace both `grep -cP` uses with `grep -cE`, a literal tab via
  `TAB=$(printf '\t')`, and `[[:space:]]` for `\s`.

## Verification

Reproduced the bash 3.2 failure with a minimal
`/bin/bash -c 'v="$({ # comment with an apostrophe … })"'` case;
after the fix `./hack/lint-docstrings.sh` runs to completion on
macOS with 0 findings (rc 0), matching Linux behavior.

## Non-Goals

Rewriting the scanner for speed; changing any docstring rule.
