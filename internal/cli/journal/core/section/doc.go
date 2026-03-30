//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package section builds topic index pages and section-based site
// output for the journal.
//
// [BuildTopicIndex] aggregates entries by topic with popularity
// thresholds. [GenerateTopicsIndex] renders the topics index page.
// [GenerateTopicPage] renders a single topic's entry list.
// [WriteFormatted] and [WriteMonths] render section content into
// string builders.
package section
