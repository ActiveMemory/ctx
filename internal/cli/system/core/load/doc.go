//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package load provides the injection load gate for
// context files. It detects when the total injected
// token count exceeds a configured threshold and writes
// a flag file that downstream hooks read to emit an
// oversize warning.
//
// # Oversize Detection
//
// [WriteOversizeFlag] compares the total injected token
// count against the threshold from rc.InjectionTokenWarn.
// When the threshold is exceeded, it writes a diagnostic
// flag file to .context/state/ containing:
//
//   - A timestamp of when the oversize was detected
//   - The injected token count versus the threshold
//   - A per-file breakdown showing each file's token
//     contribution
//   - A recommended action for the user
//
// The flag file is consumed by the check-context-size
// hook, which formats and relays the warning to the
// agent session.
//
// # Data Flow
//
//  1. Hook calls WriteOversizeFlag with token totals
//  2. Function checks threshold from rc config
//  3. If exceeded, writes flag to .context/state/
//  4. check-context-size hook reads the flag
//  5. Warning is relayed to the agent session
package load
