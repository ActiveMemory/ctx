//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
)

// resourceMapping maps a context file name to its MCP resource URI suffix
// and human-readable description.
type resourceMapping struct {
	file string
	name string
	desc string
}

// resourceTable defines all individual context file resources.
var resourceTable = []resourceMapping{
	{ctxCfg.Constitution, "constitution", assets.TextDesc(assets.TextDescKeyMCPResConstitution)},
	{ctxCfg.Task, "tasks", assets.TextDesc(assets.TextDescKeyMCPResTasks)},
	{ctxCfg.Convention, "conventions", assets.TextDesc(assets.TextDescKeyMCPResConventions)},
	{ctxCfg.Architecture, "architecture", assets.TextDesc(assets.TextDescKeyMCPResArchitecture)},
	{ctxCfg.Decision, "decisions", assets.TextDesc(assets.TextDescKeyMCPResDecisions)},
	{ctxCfg.Learning, "learnings", assets.TextDesc(assets.TextDescKeyMCPResLearnings)},
	{ctxCfg.Glossary, "glossary", assets.TextDesc(assets.TextDescKeyMCPResGlossary)},
	{ctxCfg.AgentPlaybook, "playbook", assets.TextDesc(assets.TextDescKeyMCPResPlaybook)},
}

// resourceURI builds a resource URI from a suffix.
func resourceURI(name string) string {
	return mcp.MCPResourceURIPrefix + name
}

// handleResourcesList returns all available MCP resources.
func (s *Server) handleResourcesList(req Request) *Response {
	resources := make([]Resource, 0, len(resourceTable)+1)

	// Individual context files.
	for _, rm := range resourceTable {
		resources = append(resources, Resource{
			URI:         resourceURI(rm.name),
			Name:        rm.name,
			MimeType:    mcp.MimeMarkdown,
			Description: rm.desc,
		})
	}

	// Assembled context packet (all files in read order).
	resources = append(resources, Resource{
		URI:         resourceURI("agent"),
		Name:        "agent",
		MimeType:    mcp.MimeMarkdown,
		Description: assets.TextDesc(assets.TextDescKeyMCPResAgent),
	})

	return s.ok(req.ID, ResourceListResult{Resources: resources})
}

// handleResourcesRead returns the content of a requested resource.
func (s *Server) handleResourcesRead(req Request) *Response {
	var params ReadResourceParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, errCodeInvalidArg, assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}

	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.error(req.ID, errCodeInternal,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	// Check for individual file resources.
	for _, rm := range resourceTable {
		if params.URI == resourceURI(rm.name) {
			return s.readContextFile(req.ID, ctx, rm.file, params.URI)
		}
	}

	// Assembled agent packet.
	if params.URI == resourceURI("agent") {
		return s.readAgentPacket(req.ID, ctx)
	}

	return s.error(req.ID, errCodeInvalidArg,
		fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPUnknownResource), params.URI))
}

// readContextFile returns the content of a single context file.
func (s *Server) readContextFile(
	id json.RawMessage, ctx *context.Context, fileName, uri string,
) *Response {
	f := ctx.File(fileName)
	if f == nil {
		return s.error(id, errCodeInvalidArg,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPFileNotFound), fileName))
	}

	return s.ok(id, ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      uri,
			MimeType: mcp.MimeMarkdown,
			Text:     string(f.Content),
		}},
	})
}

// readAgentPacket assembles all context files in read order into a
// single response, respecting the configured token budget.
//
// Files are added in priority order (ReadOrder). When the token
// budget would be exceeded, remaining files are listed as "Also noted"
// summaries instead of included in full.
func (s *Server) readAgentPacket(
	id json.RawMessage, ctx *context.Context,
) *Response {
	var sb strings.Builder
	header := assets.TextDesc(assets.TextDescKeyMCPPacketHeader)
	sb.WriteString(header)

	tokensUsed := context.EstimateTokensString(header)
	budget := s.tokenBudget
	var skipped []string

	for _, fileName := range ctxCfg.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}

		section := fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPSectionFormat), fileName, string(f.Content))
		sectionTokens := context.EstimateTokensString(section)

		if budget > 0 && tokensUsed+sectionTokens > budget {
			skipped = append(skipped, fileName)
			continue
		}

		sb.WriteString(section)
		tokensUsed += sectionTokens
	}

	if len(skipped) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPAlsoNoted))
		for _, name := range skipped {
			fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPOmittedFormat), name)
		}
		sb.WriteString(token.NewlineLF)
	}

	return s.ok(id, ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      resourceURI("agent"),
			MimeType: mcp.MimeMarkdown,
			Text:     sb.String(),
		}},
	})
}
