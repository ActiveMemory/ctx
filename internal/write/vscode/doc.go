//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode provides terminal output for VS Code artifact generation
// during ctx init.
//
// [InfoCreated] and [InfoExistsSkipped] report file creation results.
// [InfoRecommendationExists] and [InfoAddManually] guide users through
// manual extension setup. [InfoWarnNonFatal] reports non-fatal errors
// without aborting the init flow.
//
// Key exports: [InfoCreated], [InfoExistsSkipped],
// [InfoRecommendationExists], [InfoAddManually], [InfoWarnNonFatal].
// Used by the setup core packages when deploying VS Code integration.
package vscode
