//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx doctor command for context
// health diagnostics.
//
// [Cmd] builds the cobra.Command with --json flag. [Run] executes
// all health checks (initialization, required files, ctxrc
// validation, drift, token budget) and renders the results as
// a checklist or JSON.
package root
