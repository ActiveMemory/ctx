//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains the individual health checks for the
// doctor command.
//
// Each Check* function adds findings to a shared Report:
// [CheckContextInitialized] verifies .context/ exists,
// [CheckRequiredFiles] verifies mandatory files are present,
// [CheckCtxrcValidation] validates .ctxrc syntax, and
// [CheckDrift] runs the drift detector.
package core
