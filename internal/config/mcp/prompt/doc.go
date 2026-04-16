//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prompt defines MCP prompt names that mirror ctx
// CLI skills exposed as interactive prompt templates.
//
// MCP prompts are pre-authored message sequences that an
// AI client can retrieve and inject into its conversation.
// Each prompt maps to a ctx skill; for example,
// "ctx-session-start" guides the agent through the
// bootstrap and context-load ceremony that every session
// should begin with.
//
// # Key Constants
//
//   - [SessionStart] ("ctx-session-start"): the
//     session initialization prompt. Instructs the
//     agent to run ctx system bootstrap, read the
//     playbook, and load the context packet.
//   - [AddDecision] ("ctx-decision-add"): guides
//     the agent through recording an architectural
//     decision in DECISIONS.md.
//   - [AddLearning] ("ctx-learning-add"): guides
//     the agent through recording a learning or
//     gotcha in LEARNINGS.md.
//   - [Reflect] ("ctx-reflect"): prompts a mid-
//     session reflection on progress and blockers.
//   - [Checkpoint] ("ctx-checkpoint"): prompts a
//     lightweight checkpoint save.
//   - [RoleUser] ("user"): the MCP message role
//     for user-originated prompt content.
//
// # Why These Are Centralized
//
// Prompt registration, prompt-get handlers, and the
// skills that generate prompt content all reference
// these names. A constant prevents a registration name
// from silently diverging from the handler that
// resolves it.
package prompt
