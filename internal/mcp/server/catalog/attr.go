//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package catalog

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/resource"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Init builds the URI-to-file lookup map. Must be called once before
// FileForURI is used (called from NewServer during bootstrap).
func Init() {
	uriLookup = make(map[string]string, len(table))
	for _, m := range table {
		uriLookup[URI(m.name)] = m.file
	}
}

// URI builds a resource URI from a name suffix.
//
// Parameters:
//   - name: resource name suffix (e.g., "tasks", "agent")
//
// Returns:
//   - string: full URI (e.g., "ctx://context/tasks")
func URI(name string) string {
	return server.ResourceURIPrefix + name
}

// AgentURI returns the full URI for the assembled agent packet.
//
// Returns:
//   - string: agent resource URI
func AgentURI() string {
	return URI(resource.Agent)
}

// FileForURI returns the context file name for a resource URI or
// empty string if the URI is not a known file resource.
//
// Parameters:
//   - uri: full resource URI to look up
//
// Returns:
//   - string: context file name, or "" if unknown
func FileForURI(uri string) string {
	return uriLookup[uri]
}

// ToList constructs the immutable resource list. Called once from
// NewServer.
//
// Returns:
//   - proto.ResourceListResult: all resources including the agent packet
func ToList() proto.ResourceListResult {
	rr := make([]proto.Resource, 0, len(table)+1)

	for _, m := range table {
		rr = append(rr, proto.Resource{
			URI:         URI(m.name),
			Name:        m.name,
			MimeType:    mime.Markdown,
			Description: m.desc,
		})
	}

	rr = append(rr, proto.Resource{
		URI:         AgentURI(),
		Name:        resource.Agent,
		MimeType:    mime.Markdown,
		Description: desc.Text(text.DescKeyMCPResAgent),
	})

	return proto.ResourceListResult{Resources: rr}
}
