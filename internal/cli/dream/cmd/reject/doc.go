//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reject wires the "ctx dream reject <id>" subcommand.
//
// It builds the cobra command and delegates to the dispose core
// logic, which records a rejection in the ledger with no mutation. A
// rejected proposal is not re-surfaced unless its source changes.
package reject
