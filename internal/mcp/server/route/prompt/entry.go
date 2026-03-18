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

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	promptdef "github.com/ActiveMemory/ctx/internal/mcp/server/def/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
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
	id json.RawMessage, spec promptdef.EntryPromptSpec,
) *proto.Response {
	fieldFmt := assets.TextDesc(spec.FieldFmtK)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(spec.KeyHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	for _, f := range spec.Fields {
		_, _ = fmt.Fprintf(
			&sb,
			fieldFmt, assets.TextDesc(f.KeyLabel), f.Value,
		)
	}
	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(spec.KeyFooter))

	return out.OkResponse(id, proto.GetPromptResult{
		Description: assets.TextDesc(spec.KeyResultD),
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
