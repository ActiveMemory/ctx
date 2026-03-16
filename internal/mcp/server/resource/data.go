//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
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
	{ctxCfg.Constitution, resource.Constitution, assets.TextDesc(assets.TextDescKeyMCPResConstitution)},
	{ctxCfg.Task, resource.Tasks, assets.TextDesc(assets.TextDescKeyMCPResTasks)},
	{ctxCfg.Convention, resource.Conventions, assets.TextDesc(assets.TextDescKeyMCPResConventions)},
	{ctxCfg.Architecture, resource.Architecture, assets.TextDesc(assets.TextDescKeyMCPResArchitecture)},
	{ctxCfg.Decision, resource.Decisions, assets.TextDesc(assets.TextDescKeyMCPResDecisions)},
	{ctxCfg.Learning, resource.Learnings, assets.TextDesc(assets.TextDescKeyMCPResLearnings)},
	{ctxCfg.Glossary, resource.Glossary, assets.TextDesc(assets.TextDescKeyMCPResGlossary)},
	{ctxCfg.AgentPlaybook, resource.Playbook, assets.TextDesc(assets.TextDescKeyMCPResPlaybook)},
}

// uriLookup maps full resource URIs to context file names. Populated
// by Init during server bootstrap.
var uriLookup map[string]string
