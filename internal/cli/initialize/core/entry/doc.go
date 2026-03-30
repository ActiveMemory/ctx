//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry handles template creation and insertion point
// detection during initialization.
//
// [FindInsertionPoint] locates where new content should be inserted
// in an existing file. [CreateTemplates] writes context file
// templates (TASKS.md, DECISIONS.md, etc.) to the target directory.
package entry
