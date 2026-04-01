//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"github.com/ActiveMemory/ctx/internal/mcp/session"
)

// Handler contains domain logic for MCP operations.
//
// It holds the context directory, token budget, and session state
// needed by tool handlers. The Server package delegates to Handler
// for all domain work and handles only protocol translation.
//
// Fields:
//   - ContextDir: Path to the .context/ directory
//   - TokenBudget: Maximum token budget for context assembly
//   - Session: Per-session advisory state
type Handler struct {
	ContextDir  string
	TokenBudget int
	Session     *session.State
}

// New creates a Handler for the given context directory.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - tokenBudget: maximum token budget for context assembly
//
// Returns:
//   - *Handler: initialized handler with fresh session state
func New(contextDir string, tokenBudget int) *Handler {
	return &Handler{
		ContextDir:  contextDir,
		TokenBudget: tokenBudget,
		Session:     session.NewState(contextDir),
	}
}
