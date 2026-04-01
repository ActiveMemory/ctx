//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package git wraps git command execution behind typed functions.
//
// All exec.Command calls for git are centralized here. LookPath
// is checked once per call. Callers never import os/exec directly.
//
// Key exports: [Run], [Root], [RemoteURL], [LogSince],
// [LastCommitMessage], [DiffTreeHead].
// Part of the exec subsystem.
package git
