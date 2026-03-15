//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ctx

// `ctx` file name constants for the .context/ directory.
const (
	// Constitution contains inviolable rules for agents.
	Constitution = "CONSTITUTION.md"
	// Task contains current work items and their status.
	Task = "TASKS.md"
	// Convention contains code patterns and standards.
	Convention = "CONVENTIONS.md"
	// Architecture contains system structure documentation.
	Architecture = "ARCHITECTURE.md"
	// Decision contains architectural decisions with rationale.
	Decision = "DECISIONS.md"
	// Learning contains gotchas, tips, and lessons learned.
	Learning = "LEARNINGS.md"
	// Glossary contains domain terms and definitions.
	Glossary = "GLOSSARY.md"
	// AgentPlaybook contains the meta-instructions for using the
	// context system.
	AgentPlaybook = "AGENT_PLAYBOOK.md"
	// Dependency contains project dependency documentation.
	Dependency = "DEPENDENCIES.md"
)

// ReadOrder defines the priority order for reading context files.
//
// The order follows a logical progression for AI agents:
//
//  1. CONSTITUTION: Inviolable rules. Must be loaded first so the agent
//     knows what it cannot do before attempting anything.
//
//  2. TASKS: Current work items. What the agent should focus on.
//
//  3. CONVENTIONS: How to write code. Patterns and standards to follow.
//
//  4. ARCHITECTURE: System structure. Understanding of components and
//     boundaries before making changes.
//
//  5. DECISIONS: Historical context. Why things are the way they are,
//     to avoid re-debating settled decisions.
//
//  6. LEARNINGS: Gotchas and tips. Lessons from past work that inform
//     current implementation.
//
//  7. GLOSSARY: Reference material. Domain terms and abbreviations for
//     lookup as needed.
//
//  8. AGENT_PLAYBOOK: Meta instructions. How to use this context system.
//     Loaded last because it's about the system itself, not the work.
//     The agent should understand the content before the operating manual.
var ReadOrder = []string{
	Constitution,
	Task,
	Convention,
	Architecture,
	Decision,
	Learning,
	Glossary,
	AgentPlaybook,
}
