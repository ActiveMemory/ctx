//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/mcp/session"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// Handler contains domain logic for MCP operations.
//
// It holds the context directory, token budget, and session state
// needed by tool handlers. The Server package delegates to Handler
// for all domain work and handles only protocol translation.
type Handler struct {
	ContextDir  string
	TokenBudget int
	Session     *session.State
}

// New creates a Handler for the given context directory.
func New(contextDir string, tokenBudget int) *Handler {
	return &Handler{
		ContextDir:  contextDir,
		TokenBudget: tokenBudget,
		Session:     session.NewState(contextDir),
	}
}

// checkBoundary validates the context directory boundary.
func (h *Handler) checkBoundary() error {
	return validation.ValidateBoundary(h.ContextDir)
}

// loadContext loads the context directory.
func (h *Handler) loadContext() (*context.Context, error) {
	return context.Load(h.ContextDir)
}
