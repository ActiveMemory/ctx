//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package plugin handles Claude Code plugin detection and enablement
// during initialization.
//
// [EnableGlobally] registers the ctx plugin in Claude Code's global
// settings. [Installed] checks if the plugin binary exists.
// [EnabledGlobally] and [EnabledLocally] check registration status
// in global and project-level settings respectively.
package plugin
