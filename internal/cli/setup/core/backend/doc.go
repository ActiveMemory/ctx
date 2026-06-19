//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend writes .ctxrc backend setup entries.
//
// It owns both dry-run snippet generation and safe file updates for the
// setup command's backend mode. The package deliberately does not create
// runtime backend instances or wire command dispatch.
package backend
