//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend defines the optional AI backend contract and an
// internal registry for resolving configured backends.
//
// The package contains no concrete HTTP backend implementations. It is
// the dispatch seam later rc-backed configuration and ctx ai commands
// use to select a backend, ping it, and request completions.
package backend
