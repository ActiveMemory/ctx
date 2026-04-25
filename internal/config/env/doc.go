//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package env centralizes environment variable names
// and toggle values used by the ctx CLI at runtime.
//
// Rather than scattering string literals like "HOME" or
// "CTX_DIR" across the codebase, every environment
// variable is defined here once. This makes it easy to
// audit which variables ctx reads, prevents typos, and
// keeps grep-ability high.
//
// # Variable Names
//
// Core variables that control ctx behavior:
//
//   - Home: the user's home directory ($HOME)
//   - CtxDir: overrides the default .context/
//     directory location ($CTX_DIR)
//   - CtxTokenBudget: overrides the default token
//     budget for context window sizing
//     ($CTX_TOKEN_BUDGET)
//   - SessionID: active AI session identifier used
//     by ctx trace ($CTX_SESSION_ID)
//   - SkipPathCheck: skips PATH validation during
//     init; set to "1" in tests
//     ($CTX_SKIP_PATH_CHECK)
//
// # OS-Specific Variables
//
//   - OSWindows: the runtime.GOOS value for Windows,
//     used in platform-specific path resolution
//   - LocalAppData: the Windows %LOCALAPPDATA%
//     variable for finding config directories
//
// # Toggle Values
//
// The True constant ("1") is the canonical truthy value
// for environment variable toggles. It is used instead
// of comparing against multiple truthy strings.
//
// # Why Centralized
//
// Environment variable names are referenced by the init
// command, bootstrap logic, and test helpers.
// Centralizing them here prevents naming drift and makes
// it trivial to add new variables with consistent
// documentation.
package env
