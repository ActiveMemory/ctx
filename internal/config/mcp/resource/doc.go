//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resource defines the MCP resource name
// constants that identify context files exposed by the
// ctx MCP server.
//
// Each constant is the final path segment of a resource
// URI. The full URI is formed by prepending the server's
// URI prefix (see [server.ResourceURIPrefix]):
//
//	ctx://context/tasks
//	ctx://context/decisions
//	ctx://context/learnings
//
// When a client calls resources/list, the server returns
// these names. When a client calls resources/read with
// a URI, the server strips the prefix and matches
// against these constants to locate the file on disk.
//
// # Key Constants
//
//   - [Constitution]  -- CONSTITUTION.md, the hard
//     rules the agent must never violate.
//   - [Tasks]         -- TASKS.md, the current work
//     items and their status.
//   - [Conventions]   -- CONVENTIONS.md, coding
//     patterns and standards.
//   - [Architecture]  -- ARCHITECTURE.md, the system
//     architecture overview.
//   - [Decisions]     -- DECISIONS.md, architectural
//     decisions with rationale.
//   - [Learnings]     -- LEARNINGS.md, gotchas and
//     lessons learned.
//   - [Glossary]      -- GLOSSARY.md, project-specific
//     terminology definitions.
//   - [Playbook]      -- AGENT_PLAYBOOK.md, the agent
//     operating manual.
//   - [Agent]         -- the assembled context packet
//     (output of ctx agent).
//
// # Why These Are Centralized
//
// Resource registration, URI parsing, subscription
// management, and change-notification dispatch all
// reference these names. A constant ensures the
// registration list and the read handler always agree.
package resource
