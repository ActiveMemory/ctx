//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory provides terminal output for the memory bridge
// commands (ctx memory status, sync, diff).
//
// Status output renders a dashboard with source path, mirror path,
// sync timestamps, line counts, drift detection, and archive counts.
// Functions are composable: the caller assembles the status display
// by calling [BridgeHeader], [Source], [Mirror], [LastSync],
// [SourceLines], [DriftDetected]/[DriftNone], and [Archives] in
// sequence, separated by [StatusSeparator].
//
// Example (status command):
//
//	write.BridgeHeader(cmd)
//	write.Source(cmd, sourcePath)
//	write.Mirror(cmd, mirrorRelPath)
//	write.LastSync(cmd, formatted, ago)
//	write.StatusSeparator(cmd)
//	write.SourceLines(cmd, count, drifted)
package memory
