//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package content defines constants for content detection
// and file-emptiness heuristics.
//
// When ctx loads context files, it needs to distinguish
// meaningful content from files that exist but contain
// only boilerplate (headers, blank lines). This package
// provides the threshold used by that detection logic.
//
// # Minimum Length
//
// [MinLen] sets the floor at 20 bytes. Files shorter than
// this are treated as "effectively empty" by the agent
// packet builder, which skips them to avoid wasting
// token budget on placeholder content.
//
// The 20-byte threshold accommodates a typical copyright
// header while still catching files that contain nothing
// but a package declaration or a single blank comment.
//
// # Why Centralized
//
// The emptiness heuristic is used by the agent budget
// allocator, the drift detector, and the bootstrap file
// lister. A shared constant ensures they all agree on
// what "empty" means.
package content
