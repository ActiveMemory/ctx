//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tpl handles template deployment during initialization.
//
// [DeployTemplates] copies embedded template files to the target
// directory, creating subdirectories as needed. Existing files
// are skipped unless force mode is enabled.
package tpl
