//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook implements git hook installation and removal for context
// tracing.
//
// [Enable] installs both prepare-commit-msg and post-commit hooks.
// [Disable] removes them if they were installed by ctx. [Install]
// writes a hook script to disk, refusing to overwrite non-ctx hooks.
// [FilePath] resolves the absolute path to a named git hook.
//
// Key exports: [Enable], [Disable], [Install], [Remove], [FilePath].
// Called by the trace hook CLI subcommand.
package hook
