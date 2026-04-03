//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package builder defines the dependency graph builder
// interface and registry for supported ecosystems.
//
// Each ecosystem (Go, Node, Python, Rust) implements
// the [Builder] interface. [Detect] auto-detects the
// project type; [Find] looks up a builder by name.
package builder
