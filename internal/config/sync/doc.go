//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync defines constants for the ctx sync command.
//
// Glob patterns ([PatternTSConfig], [PatternMakefile], etc.) identify
// config files to check for drift. Action constants ([ActionDeps],
// [ActionConfig], [ActionNewDir]) classify sync check results.
// [ImportantDirs] and [SkipDirs] control directory scanning scope.
//
// Key exports: [ImportantDirs], [SkipDirs], [KeywordDependencies].
// Used by the sync core logic to detect project changes that need
// documentation updates.
package sync
