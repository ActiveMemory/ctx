//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides error constructors for schema
// validation.
//
// The sentinel ErrDrift is returned by the schema check command
// and the import integration when JSONL drift is detected. It
// signals a non-zero exit code without halting operation — drift
// warnings are informational, never blocking.
package schema
