//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compact implements the "ctx compact" command for cleaning up
// and consolidating context files.
//
// The compact command performs maintenance on .context/ files including
// moving completed tasks to a dedicated section, optionally archiving
// old content, and removing empty sections.
package compact
