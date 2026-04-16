//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ctx defines the canonical context file names,
// their priority-ordered read sequence, and the minimum
// required file set for a valid .context/ directory.
//
// Every ctx operation (bootstrap, agent packet building,
// drift detection, initialization) needs to know which
// files live in .context/ and in what order to read them.
// This package is the single source of truth for those
// names and that ordering.
//
// # Context File Names
//
// Each markdown file in .context/ has a named constant:
//
//   - [Constitution]: inviolable rules for agents
//   - [Task]: current work items and status
//   - [Convention]: code patterns and standards
//   - [Architecture]: system structure documentation
//   - [Decision]: architectural decisions with rationale
//   - [Learning]: gotchas, tips, lessons learned
//   - [Glossary]: domain terms and definitions
//   - [AgentPlaybook]: meta-instructions for using the
//     context system
//   - [AgentPlaybookGate]: distilled directives for the
//     session start hook
//   - [Dependency]: project dependency documentation
//
// # Read Order
//
// [ReadOrder] defines the priority sequence for loading
// context files into an agent's working memory:
//
//  1. Constitution (rules first)
//  2. Tasks (current focus)
//  3. Conventions (how to write code)
//  4. Architecture (system structure)
//  5. Decisions (historical rationale)
//  6. Learnings (practical tips)
//  7. Glossary (reference lookup)
//  8. Agent Playbook (operating manual last)
//
// This ordering ensures agents internalize constraints
// and current work before loading reference material.
//
// # Required Files
//
// [FilesRequired] lists the minimum files that must exist
// for a valid context directory: Constitution, Tasks, and
// Decisions. The drift detector flags missing required
// files as errors.
//
// # Stderr Prefix
//
// [StderrPrefix] ("ctx: ") is prepended to all warning
// messages written to stderr, giving users a clear
// origin for diagnostic output.
//
// # Unknown File Priority
//
// [UnknownFilePriority] (100) is assigned to context
// files not found in [ReadOrder], pushing them to the
// end of the loading sequence.
//
// # Why Centralized
//
// File names are referenced everywhere: init, bootstrap,
// agent, drift, hooks, skills. A single package prevents
// typos and ensures the read order is defined once.
package ctx
