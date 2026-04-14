//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill manages **reusable instruction bundles** —
// the `SKILL.md` + supporting-files trees that ship under
// `.claude/skills/<name>/` and tell an AI tool how to perform
// a recurring workflow.
//
// A skill is a self-contained directory:
//
//	skills/<name>/
//	  SKILL.md                # YAML frontmatter + instructions
//	  references/...          # optional supporting docs
//	  <executable>            # optional supporting script
//
// The package's job is to **install, list, load, and remove**
// these bundles. It does not execute them — execution is the
// AI tool's responsibility (Claude Code, Copilot CLI, etc.).
//
// # The Frontmatter Schema
//
// Each `SKILL.md` declares a [Manifest] in YAML
// frontmatter:
//
//   - **name** — globally unique identifier; used as the
//     directory name and as the slash-command alias.
//   - **description** — one-line trigger phrase the AI uses
//     to decide when to invoke the skill.
//   - **tools** — Copilot-style allowed-tools list (`bash`,
//     `read`, `write`, `edit`, `glob`, `grep`).
//   - **allowed-tools** — Claude-Code-style permission
//     scopes (`Bash(ctx:*)`, `Read`, etc.).
//
// [manifest.go] parses and validates the frontmatter;
// missing required fields produce a typed error from
// [internal/err/skill] that names the file path.
//
// # Public Surface
//
//   - **[Install]** — copies a source skill directory into
//     the target `skillsDir/<name>/`. Refuses to overwrite
//     an existing skill (the user must `Remove` first); use
//     `--force` at the CLI for replacement.
//   - **[Load]** — reads one skill by name, returns its
//     full [Skill] with manifest + body + path.
//   - **[LoadAll]** — walks the skills directory, returns
//     every loadable skill. Skills that fail to parse are
//     reported in the error slice rather than aborting the
//     load.
//   - **[Remove]** — deletes a skill directory after
//     verifying it lives under the canonical skills
//     directory (boundary check guards against `..`
//     escape).
//
// # File-Copy Semantics
//
// [copy.go] does the recursive copy with three rules:
//
//  1. **Preserve mode bits** — executable scripts stay
//     executable.
//  2. **Skip dotfiles at the source root** —
//     `.DS_Store`, `.git`, etc. never end up installed.
//  3. **Validate the destination** lies within the
//     skills-dir boundary.
//
// # Concurrency
//
// All operations are filesystem-bound and stateless.
// Callers serialize through process-level execution.
//
// # Related Packages
//
//   - [internal/cli/skill]    — the `ctx skill install /
//     list / remove` CLI surface.
//   - [internal/err/skill]    — typed error constructors
//     used here and by callers for consistent messaging.
//   - [internal/entity]       — [Skill], [Manifest]
//     domain types.
//   - [internal/assets/claude/skills]                 —
//     the project's own bundled skills, deployed to
//     user environments at `ctx init` time.
package skill
