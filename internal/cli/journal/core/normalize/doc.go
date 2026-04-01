//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package normalize sanitizes journal markdown for site rendering.
//
// [MatchTurnHeader] parses conversation turn headers. [FindTurnBoundary]
// locates turn boundaries in content. [TrimBlankLines] removes
// leading and trailing blank lines from a slice. The main Content
// function (not exported here) handles fence stripping, heading
// demotion, and HTML escaping.
package normalize
