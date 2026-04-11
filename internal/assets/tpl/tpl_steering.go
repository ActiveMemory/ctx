//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

import "github.com/ActiveMemory/ctx/internal/config/token"

// Foundation steering file names.
const (
	// SteeringNameProduct is the name for the product context file.
	SteeringNameProduct = "product"
	// SteeringNameTech is the name for the technology stack file.
	SteeringNameTech = "tech"
	// SteeringNameStructure is the name for the project structure file.
	SteeringNameStructure = "structure"
	// SteeringNameWorkflow is the name for the development workflow file.
	SteeringNameWorkflow = "workflow"
)

// Foundation steering file descriptions.
const (
	// SteeringDescProduct describes the product context file.
	SteeringDescProduct = "Product context, goals, and target users"
	// SteeringDescTech describes the technology stack file.
	SteeringDescTech = "Technology stack, constraints, " +
		"and dependencies"
	// SteeringDescStructure describes the project structure file.
	SteeringDescStructure = "Project structure and " +
		"directory conventions"
	// SteeringDescWorkflow describes the development workflow file.
	SteeringDescWorkflow = "Development workflow and process rules"
)

// nl is a shorthand for the LF newline used when building
// the steering body templates below, keeping line joins
// readable without drowning the file in `token.NewlineLF`
// references.
const nl = token.NewlineLF

// steeringGuidanceComment is the HTML comment prepended to
// every scaffolded steering body. It explains the three
// inclusion modes, priority semantics, and tool scope so
// the file is self-documenting on first view, and is
// invisible in rendered markdown output.
//
// The closing instruction asks the user to delete the
// comment once they've customized the file — the block is
// scaffolding, not forever guidance.
const steeringGuidanceComment = "<!--" + nl +
	"  This is a ctx steering file — persistent behavioral" + nl +
	"  rules prepended to prompts based on the frontmatter" + nl +
	"  above." + nl +
	nl +
	"  inclusion (when the rule fires):" + nl +
	"    always  → injected on EVERY tool call. On Claude" + nl +
	"              Code this is the only mode that fires" + nl +
	"              AUTOMATICALLY (the PreToolUse hook" + nl +
	"              passes an empty prompt to ctx agent)." + nl +
	"              Use for invariants and for any rule" + nl +
	"              that MUST fire reliably." + nl +
	"    auto    → injected when the prompt matches the" + nl +
	"              `description` field above." + nl +
	"                - Cursor / Cline / Kiro: native." + nl +
	"                - Claude Code: Claude calls the" + nl +
	"                  ctx_steering_get MCP tool on its" + nl +
	"                  own when it decides the rule is" + nl +
	"                  relevant. The ctx plugin ships" + nl +
	"                  the MCP auto-registration; verify" + nl +
	"                  with `claude mcp list`." + nl +
	"    manual  → only when the file is explicitly named" + nl +
	"              (e.g. via the MCP tool or a skill)." + nl +
	nl +
	"  priority (ordering within a tier):" + nl +
	"    Lower numbers inject first. 10 is a reasonable" + nl +
	"    default for invariants; use 50 for normal rules." + nl +
	nl +
	"  tools (scope the rule to specific AI tools):" + nl +
	"    Empty list = applies to all tools (default)." + nl +
	"    Example:  tools: [claude, cursor]" + nl +
	nl +
	"  Edit the body below, then delete this comment." + nl +
	"  See docs/cli/steering.md for the full reference." + nl +
	"-->" + nl + nl

// Foundation steering file body templates. Each template
// starts with the shared guidance comment so every
// scaffolded file explains the three inclusion modes to
// the user on first open.
const (
	// SteeringBodyProduct is the body for the product context file.
	SteeringBodyProduct = steeringGuidanceComment +
		"# Product Context" + nl + nl +
		"Describe the product, its goals, and target users." + nl + nl +
		"- **What is this project?**" + nl +
		"- **Who uses it?**" + nl +
		"- **What problem does it solve?**" + nl +
		"- **What is explicitly out of scope?**" + nl
	// SteeringBodyTech is the body for the technology stack file.
	SteeringBodyTech = steeringGuidanceComment +
		"# Technology Stack" + nl + nl +
		"Describe the technology stack, constraints, and " +
		"key dependencies." + nl + nl +
		"- **Languages and versions**" + nl +
		"- **Frameworks and key libraries**" + nl +
		"- **Runtime / deployment target**" + nl +
		"- **Hard constraints** (e.g. no CGO, no network " +
		"at test time)" + nl
	// SteeringBodyStructure is the body for the project structure
	// file.
	SteeringBodyStructure = steeringGuidanceComment +
		"# Project Structure" + nl + nl +
		"Describe the project layout and directory " +
		"conventions." + nl + nl +
		"- **Top-level directories and their purpose**" + nl +
		"- **Where new files should go** (and where they " +
		"should not)" + nl +
		"- **Naming conventions** for files, packages, " +
		"modules" + nl
	// SteeringBodyWorkflow is the body for the development workflow
	// file.
	SteeringBodyWorkflow = steeringGuidanceComment +
		"# Development Workflow" + nl + nl +
		"Describe the development workflow, branching " +
		"strategy, and process rules." + nl + nl +
		"- **Branch strategy** (main-only, trunk-based, " +
		"feature branches)" + nl +
		"- **Commit conventions** (message format, " +
		"signed-off-by)" + nl +
		"- **Pre-commit / pre-push checks**" + nl +
		"- **Review expectations**" + nl
)
