//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package governance defines threshold constants for MCP hook governance
// checks.
//
// [DriftCheckInterval] controls the minimum time between drift reminders.
// [PersistNudgeAfter] and [PersistNudgeRepeat] control when and how
// often persist reminders fire based on tool call counts.
// [DriftCheckMinCalls] sets the floor before first drift check.
//
// Key exports: [DriftCheckInterval], [PersistNudgeAfter],
// [PersistNudgeRepeat], [DriftCheckMinCalls].
// Used by MCP hook handlers to throttle governance nudges.
package governance
