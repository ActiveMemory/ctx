//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/notify"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// resourceMapping maps a context file name to its MCP resource URI suffix
// and human-readable description.
type resourceMapping struct {
	file string
	name string
	desc string
}

// resources defines all individual context file resources.
var resources = []resourceMapping{
	{
		ctxCfg.Constitution,
		"constitution",
		assets.TextDesc(assets.TextDescKeyMCPResConstitution),
	},
	{
		ctxCfg.Task,
		"tasks",
		assets.TextDesc(assets.TextDescKeyMCPResTasks),
	},
	{
		ctxCfg.Convention,
		"conventions",
		assets.TextDesc(assets.TextDescKeyMCPResConventions),
	},
	{
		ctxCfg.Architecture,
		"architecture",
		assets.TextDesc(assets.TextDescKeyMCPResArchitecture),
	},
	{
		ctxCfg.Decision,
		"decisions",
		assets.TextDesc(assets.TextDescKeyMCPResDecisions),
	},
	{
		ctxCfg.Learning,
		"learnings",
		assets.TextDesc(assets.TextDescKeyMCPResLearnings),
	},
	{
		ctxCfg.Glossary,
		"glossary",
		assets.TextDesc(assets.TextDescKeyMCPResGlossary),
	},
	{
		ctxCfg.AgentPlaybook,
		"playbook",
		assets.TextDesc(assets.TextDescKeyMCPResPlaybook),
	},
}

// resourceURI builds a resource URI from a suffix.
func resourceURI(name string) string {
	return server.ResourceURIPrefix + name
}

// handleResourcesList returns all available MCP resources.
func (s *Server) handleResourcesList(req proto.Request) *proto.Response {
	rr := make([]proto.Resource, 0, len(resources)+1)

	// Individual context files.
	for _, rm := range resources {
		rr = append(rr, proto.Resource{
			URI:         resourceURI(rm.name),
			Name:        rm.name,
			MimeType:    mime.Markdown,
			Description: rm.desc,
		})
	}

	// Assembled context packet (all files in read order).
	rr = append(rr, proto.Resource{
		URI:         resourceURI("agent"),
		Name:        "agent",
		MimeType:    mime.Markdown,
		Description: assets.TextDesc(assets.TextDescKeyMCPResAgent),
	})

	return s.ok(req.ID, proto.ResourceListResult{Resources: rr})
}

// handleResourcesRead returns the content of a requested resource.
func (s *Server) handleResourcesRead(req proto.Request) *proto.Response {
	var params proto.ReadResourceParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, proto.ErrCodeInvalidArg, assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}

	ctx, err := context.Load(s.handler.ContextDir)
	if err != nil {
		return s.error(req.ID, proto.ErrCodeInternal,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	// Check for individual file resources.
	for _, rm := range resources {
		if params.URI == resourceURI(rm.name) {
			return s.readContextFile(req.ID, ctx, rm.file, params.URI)
		}
	}

	// Assembled agent packet.
	if params.URI == resourceURI("agent") {
		return s.readAgentPacket(req.ID, ctx)
	}

	return s.error(req.ID, proto.ErrCodeInvalidArg,
		fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPUnknownResource), params.URI))
}

// readContextFile returns the content of a single context file.
func (s *Server) readContextFile(
	id json.RawMessage, ctx *context.Context, fileName, uri string,
) *proto.Response {
	f := ctx.File(fileName)
	if f == nil {
		return s.error(id, proto.ErrCodeInvalidArg,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPFileNotFound), fileName))
	}

	return s.ok(id, proto.ReadResourceResult{
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
// budget would be exceeded, remaining files are listed as "Also noted"
// summaries instead of included in full.
func (s *Server) readAgentPacket(
	id json.RawMessage, ctx *context.Context,
) *proto.Response {
	var sb strings.Builder
	header := assets.TextDesc(assets.TextDescKeyMCPPacketHeader)
	sb.WriteString(header)

	tokensUsed := context.EstimateTokensString(header)
	budget := s.handler.TokenBudget
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

	return s.ok(id, proto.ReadResourceResult{
		Contents: []proto.ResourceContent{{
			URI:      resourceURI("agent"),
			MimeType: mime.Markdown,
			Text:     sb.String(),
		}},
	})
}

// defaultPollInterval is the default interval for resource change polling.
const defaultPollInterval = 5 * time.Second

// ResourcePoller tracks subscribed resources and polls for file changes.
type ResourcePoller struct {
	mu         sync.Mutex
	subs       map[string]bool      // URI → subscribed
	mtimes     map[string]time.Time // file path → last known mtime
	contextDir string
	pollStop   chan struct{}
	notifyFunc func(proto.Notification) // callback to emit notifications
}

// NewResourcePoller creates a poller for the given context directory.
func NewResourcePoller(contextDir string, notify func(proto.Notification)) *ResourcePoller {
	return &ResourcePoller{
		subs:       make(map[string]bool),
		mtimes:     make(map[string]time.Time),
		contextDir: contextDir,
		notifyFunc: notify,
	}
}

// subscribe adds a URI to the watch set and starts polling if needed.
//
// Goroutine lifecycle: the poller goroutine is started on the first
// subscription and stopped when the last subscription is removed or
// when Server.Serve returns (via poller.stop in the deferred cleanup).
func (p *ResourcePoller) subscribe(uri string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subs[uri] = true

	// Snapshot current mtime for the resource's file.
	if fileName := p.uriToFile(uri); fileName != "" {
		fpath := filepath.Join(p.contextDir, fileName)
		if info, err := os.Stat(fpath); err == nil {
			p.mtimes[fpath] = info.ModTime()
		}
	}

	// Start poller if this is the first subscription.
	if len(p.subs) == 1 && p.pollStop == nil {
		p.pollStop = make(chan struct{})
		go p.poll()
	}
}

// unsubscribe removes a URI from the watch set and stops polling if empty.
func (p *ResourcePoller) unsubscribe(uri string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subs, uri)

	if len(p.subs) == 0 && p.pollStop != nil {
		close(p.pollStop)
		p.pollStop = nil
	}
}

// Stop shuts down the poller goroutine.
func (p *ResourcePoller) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pollStop != nil {
		close(p.pollStop)
		p.pollStop = nil
	}
}

// uriToFile maps a resource URI to its context file name.
func (p *ResourcePoller) uriToFile(uri string) string {
	for _, rm := range resources {
		if uri == resourceURI(rm.name) {
			return rm.file
		}
	}
	return ""
}

// poll checks subscribed resources for mtime changes on a fixed interval.
func (p *ResourcePoller) poll() {
	ticker := time.NewTicker(defaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.pollStop:
			return
		case <-ticker.C:
			p.checkChanges()
		}
	}
}

// checkChanges compares current mtimes to snapshots and emits notifications.
func (p *ResourcePoller) checkChanges() {
	p.mu.Lock()
	uris := make([]string, 0, len(p.subs))
	for uri := range p.subs {
		uris = append(uris, uri)
	}
	p.mu.Unlock()

	for _, uri := range uris {
		fileName := p.uriToFile(uri)
		if fileName == "" {
			continue
		}
		fpath := filepath.Join(p.contextDir, fileName)
		info, err := os.Stat(fpath)
		if err != nil {
			continue
		}

		p.mu.Lock()
		prev, known := p.mtimes[fpath]
		if known && info.ModTime().After(prev) {
			p.mtimes[fpath] = info.ModTime()
			p.mu.Unlock()
			p.notifyFunc(proto.Notification{
				JSONRPC: server.JSONRPCVersion,
				Method:  notify.ResourcesUpdated,
				Params:  proto.ResourceUpdatedParams{URI: uri},
			})
		} else {
			if !known {
				p.mtimes[fpath] = info.ModTime()
			}
			p.mu.Unlock()
		}
	}
}

// handleResourcesSubscribe registers a resource for change notifications.
func (s *Server) handleResourcesSubscribe(req proto.Request) *proto.Response {
	return s.subscriptionAction(req, s.poller.subscribe)
}

// handleResourcesUnsubscribe removes a resource from change notifications.
func (s *Server) handleResourcesUnsubscribe(req proto.Request) *proto.Response {
	return s.subscriptionAction(req, s.poller.unsubscribe)
}

// subscriptionAction unmarshals a URI param and applies the given
// subscription function (subscribe or unsubscribe).
func (s *Server) subscriptionAction(
	req proto.Request, fn func(string),
) *proto.Response {
	var params proto.SubscribeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, proto.ErrCodeInvalidArg, assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}
	if params.URI == "" {
		return s.error(req.ID, proto.ErrCodeInvalidArg, assets.TextDesc(assets.TextDescKeyMCPURIRequired))
	}
	fn(params.URI)
	return s.ok(req.ID, struct{}{})
}
