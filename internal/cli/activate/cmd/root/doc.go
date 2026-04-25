//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the `ctx activate` cobra command.
//
// Activate is the shell-integration entry point under the
// explicit-context-dir model (spec: specs/explicit-context-dir.md).
// Its single job is to emit a `export CTX_DIR=...` line to stdout so
// that callers can bind the context directory for their shell via
// `eval "$(ctx activate)"`.
//
// Unlike most commands in the CLI, `activate` is in the exempt
// allowlist: it does not call rc.RequireContextDir because
// activate's reason for existing is precisely to help users declare
// CTX_DIR in the first place.
//
// Resolution:
//
//   - With an explicit path argument: the path is validated strictly
//     (exists, is a directory, contains at least one canonical
//     context file). There is no --force escape hatch in v1.
//   - Without arguments: the command scans upward from CWD using
//     rc.ScanCandidates and emits the one visible candidate when
//     there is exactly one. Zero candidates → NoCandidates error.
//     Two or more candidates → Ambiguous error listing every path;
//     activate refuses to pick automatically.
//
// This is the only command in the CLI that walks. All other
// resolution flows through rc.ContextDir / rc.RequireContextDir.
package root
