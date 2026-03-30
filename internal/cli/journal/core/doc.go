//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains the journal processing pipeline: import
// planning, execution, site generation, normalization, and
// Obsidian vault building.
//
// The pipeline flows: plan → confirm → execute → normalize →
// generate site. Each stage is idempotent and tracks progress
// via .state.json.
package core
