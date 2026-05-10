//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	"github.com/ActiveMemory/ctx/internal/sanitize"
)

// buildEntry renders a structured entry prompt (decision or
// learning) from the given spec and returns the formatted response.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - spec: entry prompt specification (header, footer, fields)
//
// Returns:
//   - *proto.Response: formatted entry prompt
func buildEntry(
	id json.RawMessage, spec entity.PromptEntrySpec,
) *proto.Response {
	fieldFmt := desc.Text(spec.FieldFmtK)

	var sb strings.Builder
	sb.WriteString(desc.Text(spec.KeyHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	for _, f := range spec.Fields {
		// MCP-SAN.3: sanitize user-supplied content before
		// embedding in the prompt output.
		_, _ = fmt.Fprintf(
			&sb,
			fieldFmt, desc.Text(f.KeyLabel),
			sanitize.Content(f.Value),
		)
	}
	sb.WriteString(token.NewlineLF)
	sb.WriteString(desc.Text(spec.KeyFooter))

	return out.OkResponse(id, proto.GetPromptResult{
		Description: desc.Text(spec.KeyResultD),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: sb.String(),
				},
			},
		},
	})
}
