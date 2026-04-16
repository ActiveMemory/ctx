//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering loads steering and skill content for
// inclusion in the agent context packet.
//
// # LoadBodies
//
// [LoadBodies] reads all steering files from the
// steering directory (rc.SteeringDir), filters them by
// the current tool (rc.Tool), and returns the body
// content of each matching file as a string slice.
//
// Steering files are YAML-frontmattered Markdown files
// that contain tool-specific instructions. The filtering
// uses steering.Filter with the current tool identifier,
// so only files that apply to the active AI tool are
// included in the context packet. When the steering
// directory is missing or contains no applicable files,
// LoadBodies returns nil.
//
// # LoadSkill
//
// [LoadSkill] loads a named skill from the
// .context/skills/ directory and returns its body
// content. Skills are standalone instruction files that
// can be referenced by name in agent prompts.
//
// When the skill is not found, LoadSkill returns an
// errSkill.NotFound error. Other read failures return
// errSkill.LoadQuoted with the underlying cause.
//
// # Data Flow
//
// The budget subpackage calls LoadBodies during context
// assembly to append steering instructions after the
// core context files. LoadSkill is called when a
// specific skill is requested via the --skill flag.
package steering
