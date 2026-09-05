//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"strconv"
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// MCPTable renders the `[mcp_servers.ctx]` table that registers
// the ctx MCP server with Codex.
//
// Returns:
//   - string: the table text, newline-terminated
func MCPTable() string {
	args := mcpServer.Args()
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, strconv.Quote(arg))
	}
	var sb strings.Builder
	sb.WriteString(cfgCodex.TOMLHeaderMCPCtx)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(cfgCodex.TOMLKeyCommand)
	sb.WriteString(tomlAssign())
	sb.WriteString(strconv.Quote(mcpServer.Command))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(cfgCodex.TOMLKeyArgs)
	sb.WriteString(tomlAssign())
	sb.WriteString(cfgCodex.TOMLBracketOpen)
	sb.WriteString(strings.Join(quoted, token.CommaSpace))
	sb.WriteString(cfgCodex.TOMLBracketClose)
	sb.WriteString(token.NewlineLF)
	return sb.String()
}

// EnsureMCPTable appends the ctx MCP table to a Codex config.toml
// unless its header is already present. Existing bytes are never
// rewritten: the table is appended after a blank line, so the
// result is valid TOML whatever the prior content.
//
// Parameters:
//   - existing: current config.toml content (may be empty)
//
// Returns:
//   - []byte: content to write (existing bytes when skipped)
//   - Outcome: Created (empty input), Merged (appended), or
//     Skipped (header already present)
func EnsureMCPTable(existing []byte) ([]byte, Outcome) {
	table := []byte(MCPTable())
	if len(bytes.TrimSpace(existing)) == 0 {
		return table, OutcomeCreated
	}
	if tableHeaderPresent(existing) {
		return existing, OutcomeSkipped
	}

	out := make([]byte, 0, len(existing)+len(table)+2)
	out = append(out, existing...)
	if !bytes.HasSuffix(out, []byte(token.NewlineLF)) {
		out = append(out, token.NewlineLF...)
	}
	out = append(out, token.NewlineLF...)
	out = append(out, table...)
	return out, OutcomeMerged
}
