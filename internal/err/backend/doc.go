//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend provides typed errors for AI backend registry
// failures.
//
// Callers can match these errors with errors.As while still receiving
// stable user-facing messages from the config layer. The package avoids
// sentinel strings in backend implementation code.
package backend
