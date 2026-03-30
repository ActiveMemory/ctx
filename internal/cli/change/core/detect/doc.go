//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package detect resolves reference timestamps for change detection.
//
// [FromMarkers] reads the session marker file for the last known
// timestamp. [FromEvents] reads the event log. [ReferenceTime]
// combines both with the --since flag to pick the best reference.
// [ParseSinceFlag] parses user-provided duration or date strings.
package detect
