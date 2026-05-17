//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package handover defines the typed error constructors for
// Phase KB handover artifacts. Handovers are per-session
// recall artifacts written under `.context/handovers/`; this
// package owns every error surface the handover writer /
// reader / fold mechanism can return.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/handover]
//     supplies the sentinel-message and format-string
//     constants.
//   - [github.com/ActiveMemory/ctx/internal/write/handover]
//     is the primary caller.
//   - [github.com/ActiveMemory/ctx/internal/err/closeout]
//     defines the parallel error surface for the closeout
//     artifacts that the handover fold consumes.
package handover
