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

	"github.com/ActiveMemory/ctx/internal/config"
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
	{config.FileConstitution, "constitution", "Hard rules that must never be violated"},
	{config.FileTask, "tasks", "Current work items and their status"},
	{config.FileConvention, "conventions", "Code patterns and standards"},
	{config.FileArchitecture, "architecture", "System architecture documentation"},
	{config.FileDecision, "decisions", "Architectural decisions with rationale"},
	{config.FileLearning, "learnings", "Gotchas, tips, and lessons learned"},
	{config.FileGlossary, "glossary", "Project-specific terminology"},
}

// resourceURI builds a resource URI from a suffix.
func resourceURI(name string) string {
	return "ctx://context/" + name
}

// handleResourcesList returns all available MCP resources.
func (s *Server) handleResourcesList(req Request) *Response {
	resources := make([]Resource, 0, len(resourceTable)+1)

	// Individual context files.
	for _, rm := range resourceTable {
		resources = append(resources, Resource{
			URI:         resourceURI(rm.name),
			Name:        rm.name,
			MimeType:    "text/markdown",
			Description: rm.desc,
		})
	}

	// Assembled context packet (all files in read order).
	resources = append(resources, Resource{
		URI:         resourceURI("agent"),
		Name:        "agent",
		MimeType:    "text/markdown",
		Description: "All context files assembled in priority read order",
	})

	return s.ok(req.ID, ResourceListResult{Resources: resources})
}

// handleResourcesRead returns the content of a requested resource.
func (s *Server) handleResourcesRead(req Request) *Response {
	var params ReadResourceParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, errCodeInvalidArg, "invalid params")
	}

	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.error(req.ID, errCodeInternal,
			fmt.Sprintf("failed to load context: %v", err))
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
		fmt.Sprintf("unknown resource: %s", params.URI))
}

// readContextFile returns the content of a single context file.
func (s *Server) readContextFile(
	id json.RawMessage, ctx *context.Context, fileName, uri string,
) *Response {
	f := ctx.File(fileName)
	if f == nil {
		return s.error(id, errCodeInvalidArg,
			fmt.Sprintf("file not found: %s", fileName))
	}

	return s.ok(id, ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      uri,
			MimeType: "text/markdown",
			Text:     string(f.Content),
		}},
	})
}

// readAgentPacket assembles all context files in read order into a
// single response. This provides the same content that `ctx agent`
// would output, concatenated as a single markdown document.
func (s *Server) readAgentPacket(
	id json.RawMessage, ctx *context.Context,
) *Response {
	var sb strings.Builder
	sb.WriteString("# Context Packet\n\n")

	for _, fileName := range config.FileReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}
		fmt.Fprintf(&sb, "---\n## %s\n\n", fileName)
		sb.WriteString(string(f.Content))
		sb.WriteString("\n\n")
	}

	return s.ok(id, ReadResourceResult{
		Contents: []ResourceContent{{
			URI:      resourceURI("agent"),
			MimeType: "text/markdown",
			Text:     sb.String(),
		}},
	})
}
