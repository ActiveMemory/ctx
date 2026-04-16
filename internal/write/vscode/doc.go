//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package vscode provides terminal output for VS Code
// artifact generation during ctx init.
//
// When ctx initializes a project, it optionally
// creates VS Code configuration files such as
// extensions.json recommendations and workspace
// settings. The output functions report the result
// of each file operation.
//
// # File Creation
//
// [InfoCreated] confirms a VS Code configuration
// file was created at the given path.
// [InfoExistsSkipped] reports a file was skipped
// because it already exists.
//
// # Extension Recommendations
//
// [InfoRecommendationExists] reports the ctx
// extension recommendation already exists in
// extensions.json. [InfoAddManually] guides the
// user to add the extension ID manually when the
// file exists but lacks the ctx recommendation.
//
// # Error Handling
//
// [InfoWarnNonFatal] reports a non-fatal error
// during artifact creation without aborting the
// init flow. It prints a short description and
// the underlying error.
package vscode
