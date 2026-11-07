//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"bytes"
	"encoding/json"

	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
)

// MergeHooks combines an existing hooks.json with the embedded ctx
// Antigravity hooks manifest.
//
// Parameters:
//   - existing: current hooks.json file content (may be empty)
//   - embedded: embedded ctx hooks manifest
//
// Returns:
//   - []byte: content to write
//   - Outcome: Created, Merged, or Skipped
//   - error: non-nil if parsing or encoding fails
func MergeHooks(existing, embedded []byte) ([]byte, Outcome, error) {
	if len(bytes.TrimSpace(existing)) == 0 {
		return embedded, OutcomeCreated, nil
	}

	var current map[string]json.RawMessage
	if parseErr := json.Unmarshal(existing, &current); parseErr != nil {
		return nil, OutcomeSkipped, errParser.Unmarshal(parseErr)
	}

	var shipped map[string]json.RawMessage
	if parseErr := json.Unmarshal(embedded, &shipped); parseErr != nil {
		return nil, OutcomeSkipped, errParser.Unmarshal(parseErr)
	}

	shippedCtx, ok := shipped[cfgAgy.KeyCtx]
	if !ok {
		return nil, OutcomeSkipped, errSetup.MissingEmbeddedAsset(
			cfgAgy.KeyCtx,
		)
	}

	if existingCtx, hasCtx := current[cfgAgy.KeyCtx]; hasCtx {
		same, sameErr := equivalent(existingCtx, shippedCtx)
		if sameErr == nil && same {
			return existing, OutcomeSkipped, nil
		}
	}

	outcome := OutcomeMerged
	if len(current) == 0 {
		outcome = OutcomeCreated
	}

	current[cfgAgy.KeyCtx] = shippedCtx

	out, encodeErr := encodeJSON(current, token.Indent2)
	if encodeErr != nil {
		return nil, OutcomeSkipped, encodeErr
	}
	return out, outcome, nil
}

// EnsureMCPConfig registers the ctx MCP server in the Antigravity
// mcp_config.json document.
//
// Parameters:
//   - existing: current mcp_config.json content (may be empty)
//
// Returns:
//   - []byte: content to write
//   - Outcome: Created, Merged, or Skipped
//   - error: non-nil if parsing or encoding fails
func EnsureMCPConfig(existing []byte) ([]byte, Outcome, error) {
	if len(bytes.TrimSpace(existing)) == 0 {
		fresh := map[string]any{
			cfgHook.KeyMCPServers: map[string]any{
				mcpServer.Name: map[string]any{
					cfgHook.KeyCommand: mcpServer.Command,
					cfgHook.KeyArgs:    mcpServer.Args(),
				},
			},
		}
		out, encodeErr := encodeJSON(fresh, token.Indent2)
		if encodeErr != nil {
			return nil, OutcomeSkipped, encodeErr
		}
		return out, OutcomeCreated, nil
	}

	var root map[string]any
	if parseErr := json.Unmarshal(existing, &root); parseErr != nil {
		return nil, OutcomeSkipped, errParser.Unmarshal(parseErr)
	}

	servers, ok := root[cfgHook.KeyMCPServers].(map[string]any)
	if !ok || servers == nil {
		servers = make(map[string]any)
	}

	if _, present := servers[mcpServer.Name]; present {
		return existing, OutcomeSkipped, nil
	}

	servers[mcpServer.Name] = map[string]any{
		cfgHook.KeyCommand: mcpServer.Command,
		cfgHook.KeyArgs:    mcpServer.Args(),
	}
	root[cfgHook.KeyMCPServers] = servers

	out, encodeErr := encodeJSON(root, token.Indent2)
	if encodeErr != nil {
		return nil, OutcomeSkipped, encodeErr
	}
	return out, OutcomeMerged, nil
}
