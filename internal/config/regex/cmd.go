//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// MidSudo matches mid-command sudo after && || ;
var MidSudo = regexp.MustCompile(`(;|&&|\|\|)\s*sudo\s`)

// GitPush matches `git push` invocations across common shell shapes.
//
// Covered entry points (prefix anchor `[^|(`+"`"+`\n]`):
//   - Bare `git push` at start of the command
//   - After statement separators: `;`, `&&`, `||`, `|`, `&`
//   - Subshells and command substitution: `(git push)`, `$(git push)`,
//     backtick-wrapped `git push`
//   - New lines in multi-line command input
//
// Covered prefixes (the `(\S+\s+)*` group before `git`):
//   - Environment variable assignments: `GIT_DIR=/foo git push`
//   - Command wrappers: `time git push`, `nice git push`, `nohup git push`
//
// Covered flag shapes between `git` and `push` (the `(\s+\S+)*` group):
//   - Short flags with values: `-C /path`, `-c key=value`
//   - Short boolean flags: `-p`, `-P`, `-h`, `-v`
//   - Long flags (boolean or `=value`): `--git-dir=PATH`, `--no-pager`,
//     `--bare`, `--work-tree=PATH`, etc.
//
// Trailing anchor `([^a-zA-Z0-9._/-]|$)`: matches any shell terminator
// (whitespace, `)`, backtick, `;`, `|`, `&`, `>`, `<`, quote, newline,
// end-of-string) but rejects ref-name continuations like `push-to-remote`
// or `push_branch` so `git push-to-remote` (an imagined alias) does not
// false-positive as a push subcommand.
//
// Known blind spots:
//   - False-positives on literal `push` as an argument in other
//     subcommands, e.g. `git log push` when `push` is a branch name.
//     Accepted as a safer-than-sorry trade-off for a push guard:
//     over-blocking is recoverable, under-blocking is not.
//   - Does not match through `eval` or `sh -c` quoting, e.g.
//     `eval "git push"` or `sh -c "git push"`. Parsing through arbitrary
//     shell quoting is undecidable in the general case.
//   - Shell aliases (`alias p=push; git p`) are invisible to static
//     regex matching.
//
// Uses Go's RE2 engine, so `(\S+\s+)*` is linear-time despite its
// nested-quantifier appearance. Do not port this regex to a PCRE
// engine without reviewing backtracking behavior.
var GitPush = regexp.MustCompile(
	`(^|[;&|(` + "`" + `\n]\s*)(\S+\s+)*git(\s+\S+)*\s+push([^a-zA-Z0-9._/-]|$)`,
)

// CpMvToBin matches cp/mv to bin directories.
var CpMvToBin = regexp.MustCompile(
	`(cp|mv)\s+\S+\s+` +
		`(/usr/local/bin|/usr/bin|~/go/bin|~/.local/bin` +
		`|/home/\S+/go/bin|/home/\S+/.local/bin)`)

// InstallToLocalBin matches cp/install to ~/.local/bin.
var InstallToLocalBin = regexp.MustCompile(`(cp|install)\s.*~/\.local/bin`)

// GitCommit matches git commit commands.
var GitCommit = regexp.MustCompile(`git\s+commit`)

// GitAmend matches the --amend flag.
var GitAmend = regexp.MustCompile(`--amend`)

// TaskRef matches Phase-style task references like HA.1, P-2.5, PD.3, CT.1.
var TaskRef = regexp.MustCompile(`\b[A-Z]+-?\d+\.?\d*\b`)
