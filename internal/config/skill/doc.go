//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill defines the manifest filename and
// parsing constants for the ctx skill system.
//
// Skills are reusable prompt-and-instruction bundles
// stored as directories under .claude/skills/. Each
// skill directory must contain a manifest file whose
// name is declared here as [SkillManifest]
// ("SKILL.md"). The skill loader, the skill
// scaffolder, and the skill audit command all
// reference this constant to locate manifests.
//
// # Key Constants
//
//   - [SkillManifest]: the expected filename inside
//     every skill directory ("SKILL.md"). The loader
//     rejects directories that lack this file.
//
// # Why Centralized
//
// The skill loader, the CLI skill create/list/audit
// commands, and the scheduled-task creator all need
// the manifest filename. A single constant prevents
// one path from drifting while others stay fixed.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package skill
