//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compat contains backward-compatibility integration tests
// that verify the hooks-and-steering extensions do not break existing
// ctx workflows when the new directories are absent.
//
// Tests exercise init, status, agent, and drift commands in a
// clean project to confirm graceful degradation when .context/hooks/
// and .context/steering/ directories do not exist.
package compat
