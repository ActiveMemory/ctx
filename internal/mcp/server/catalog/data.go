//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package catalog

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/resource"
)

// mapping pairs a context file name with its MCP resource name and
// human-readable description.
type mapping struct {
	file string
	name string
	desc string
}

// table defines all individual context file resources.
var table = []mapping{
	{cfgCtx.Constitution, resource.Constitution, desc.Text(text.DescKeyMCPResConstitution)},
	{cfgCtx.Task, resource.Tasks, desc.Text(text.DescKeyMCPResTasks)},
	{cfgCtx.Convention, resource.Conventions, desc.Text(text.DescKeyMCPResConventions)},
	{cfgCtx.Architecture, resource.Architecture, desc.Text(text.DescKeyMCPResArchitecture)},
	{cfgCtx.Decision, resource.Decisions, desc.Text(text.DescKeyMCPResDecisions)},
	{cfgCtx.Learning, resource.Learnings, desc.Text(text.DescKeyMCPResLearnings)},
	{cfgCtx.Glossary, resource.Glossary, desc.Text(text.DescKeyMCPResGlossary)},
	{cfgCtx.AgentPlaybook, resource.Playbook, desc.Text(text.DescKeyMCPResPlaybook)},
}

// uriLookup maps full resource URIs to context file names. Populated
// by Init during server bootstrap.
var uriLookup map[string]string
