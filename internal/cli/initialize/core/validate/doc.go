//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate runs pre-flight checks during the
// ctx init pipeline.
//
// # Overview
//
// Before the init command creates or modifies context
// files, it runs validation checks to ensure the
// environment is properly configured. This package
// provides those checks.
//
// # Behavior
//
// [CheckCtxInPath] uses exec.LookPath to verify the ctx
// binary is reachable via PATH, warning if it is missing.
// [EssentialFilesPresent] checks for TASKS.md,
// CONSTITUTION.md, or DECISIONS.md, treating a directory
// without them as uninitialised.
//
// # Behavior
//
// CheckCtxInPath uses exec.LookPath to search for the
// ctx binary. If the binary is not found, it prints a
// warning via the initialize write layer and returns
// an error. The check can be skipped by setting the
// CTX_SKIP_PATH_CHECK environment variable to "true".
//
// EssentialFilesPresent checks for the presence of any
// file in the required files list (TASKS.md,
// CONSTITUTION.md, DECISIONS.md). A directory that
// contains only logs or other non-essential content is
// considered uninitialised, allowing init to run a
// fresh setup.
package validate
