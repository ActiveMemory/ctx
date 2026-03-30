//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package format renders context entries into Markdown with
// structured sections.
//
// Each entry type has its own formatter: [Task] adds priority
// labels, [Decision] adds Context/Rationale/Consequence sections,
// [Learning] adds Context/Lesson/Application sections, and
// [Convention] wraps content with a timestamp header.
package format
