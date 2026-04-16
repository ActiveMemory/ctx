//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill provides terminal output for the skill
// management commands (ctx skill install, list, remove).
//
// Skills are reusable prompt templates that extend
// agent capabilities. The output functions cover
// the full management lifecycle.
//
// # Installation
//
// [Installed] prints a confirmation that includes
// the skill name and the directory where it was
// installed.
//
// # Listing
//
// [EntryWithDesc] prints a skill entry with its
// name and description. [Entry] prints a skill
// entry with name only, used when no description
// is available. [Count] prints the total number
// of installed skills. [NoSkillsFound] handles
// the empty-list case.
//
// # Removal
//
// [Removed] prints a confirmation that the named
// skill was removed from the skills directory.
package skill
