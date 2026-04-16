//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package format provides formatting helpers for
// journal output and site generation.
//
// # Overview
//
// This package contains utility functions for
// converting raw values into human-readable or
// URL-safe representations used throughout the
// journal pipeline.
//
// # Public Surface
//
//   - [Size] -- formats a byte count as a
//     human-readable string (e.g. "512B", "1.5KB",
//     "2.3MB"). Uses IEC units (1024-based).
//   - [KeyFileSlug] -- converts a file path to a
//     URL-safe slug by replacing path separators
//     and dots with underscores, and glob stars
//     with "x".
//   - [SessionLink] -- builds a markdown list item
//     linking to a session page with a session count
//     (e.g. "- [topic](topic.md) (3 sessions)").
//
// # Usage
//
// Size is used when reporting journal entry sizes
// during import. KeyFileSlug is used by the map of
// content generator to create safe filenames from
// file paths. SessionLink is used by the MOC builder
// to generate navigation links.
package format
