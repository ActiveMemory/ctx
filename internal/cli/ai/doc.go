//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ai exposes optional AI backend commands.
//
// The command surface is additive and fail-closed: it only runs when a
// backend is explicitly configured in .ctxrc. Deterministic context
// assembly commands do not depend on this package.
package ai
