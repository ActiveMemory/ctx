//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/catalog"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// readContextFile returns the content of a single context file.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - ctx: loaded context
//   - fileName: context file name to read
//   - uri: resource URI for the response
//
// Returns:
//   - *proto.Response: resource content or not-found error
func readContextFile(
	id json.RawMessage, ctx *entity.Context, fileName, uri string,
) *proto.Response {
	f := ctx.File(fileName)
	if f == nil {
		return out.ErrResponse(id, cfgSchema.ErrCodeInvalidArg,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPErrFileNotFound),
				fileName,
			))
	}

	return out.OkResponse(id, proto.ReadResourceResult{
		Contents: []proto.ResourceContent{{
			URI:      uri,
			MimeType: mime.Markdown,
			Text:     string(f.Content),
		}},
	})
}

// readAgentPacket assembles all context files in read order into a
// single response, respecting the configured token budget.
//
// Files are added in priority order (ReadOrder). When the token
// budget would be exceeded, remaining files are listed as "Also
// noted" summaries instead of included in full.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - ctx: loaded context
//   - budget: token budget for assembly
//
// Returns:
//   - *proto.Response: assembled context packet
func readAgentPacket(
	id json.RawMessage, ctx *entity.Context, budget int,
) *proto.Response {
	var sb strings.Builder
	header := desc.Text(text.DescKeyMCPPacketHeader)
	sb.WriteString(header)

	tokensUsed := ctxToken.EstimateString(header)
	var skipped []string

	for _, fileName := range cfgCtx.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}

		section := fmt.Sprintf(
			desc.Text(text.DescKeyMCPFormatSection),
			fileName, string(f.Content),
		)
		sectionTokens := ctxToken.EstimateString(section)

		if budget > 0 && tokensUsed+sectionTokens > budget {
			skipped = append(skipped, fileName)
			continue
		}

		sb.WriteString(section)
		tokensUsed += sectionTokens
	}

	if len(skipped) > 0 {
		sb.WriteString(
			desc.Text(text.DescKeyMCPAlsoNoted),
		)
		for _, name := range skipped {
			_, _ = fmt.Fprintf(
				&sb,
				desc.Text(text.DescKeyMCPOmittedFormat),
				name,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	return out.OkResponse(id, proto.ReadResourceResult{
		Contents: []proto.ResourceContent{{
			URI:      catalog.AgentURI(),
			MimeType: mime.Markdown,
			Text:     sb.String(),
		}},
	})
}
