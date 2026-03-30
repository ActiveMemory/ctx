//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook provides terminal output for the hook generation
// command (ctx hook) and hook lifecycle output.
//
// Functions cover hook deployment output ([InfoCopilotCreated],
// [InfoCopilotMerged], [InfoCopilotSkipped], [InfoCopilotSummary]),
// hook runtime output ([Nudge], [NudgeBlock], [BlockResponse],
// [Context]), and general-purpose hook helpers ([Content],
// [Separator], [InfoTool], [InfoUnknownTool]).
//
// Nudge vs NudgeBlock: [Nudge] emits a single-line relay,
// [NudgeBlock] emits a multi-line boxed message. Both are
// consumed by the agent as VERBATIM relay directives.
package hook
