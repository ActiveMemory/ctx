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
	res "github.com/ActiveMemory/ctx/internal/mcp/server/resource"
)

// handleResourcesList returns all available MCP resources.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *proto.Response: resource list result
func (s *Server) handleResourcesList(req proto.Request) *proto.Response {
	return s.ok(req.ID, s.resourceList)
}

// handleResourcesRead returns the content of a requested resource.
//
// Parameters:
//   - req: the MCP request containing the resource URI
//
// Returns:
//   - *proto.Response: resource content or error
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

	// Individual file resource.
	if fileName := res.FileForURI(params.URI); fileName != "" {
		return s.readContextFile(req.ID, ctx, fileName, params.URI)
	}

	// Assembled agent packet.
	if params.URI == res.AgentURI() {
		return s.readAgentPacket(req.ID, ctx)
	}

	return s.error(req.ID, proto.ErrCodeInvalidArg,
		fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPUnknownResource), params.URI))
}

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
//
// Parameters:
//   - id: JSON-RPC request ID
//   - ctx: loaded context
//
// Returns:
//   - *proto.Response: assembled context packet
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
			URI:      res.AgentURI(),
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
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - notify: callback to emit resource change notifications
//
// Returns:
//   - *ResourcePoller: initialized poller (not yet polling)
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
	if fileName := res.FileForURI(uri); fileName != "" {
		fpath := filepath.Join(p.contextDir, fileName)
		if info, statErr := os.Stat(fpath); statErr == nil {
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
		fileName := res.FileForURI(uri)
		if fileName == "" {
			continue
		}
		fpath := filepath.Join(p.contextDir, fileName)
		info, statErr := os.Stat(fpath)
		if statErr != nil {
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
//
// Parameters:
//   - req: the MCP request containing the resource URI
//
// Returns:
//   - *proto.Response: empty success or validation error
func (s *Server) handleResourcesSubscribe(req proto.Request) *proto.Response {
	return s.subscriptionAction(req, s.poller.subscribe)
}

// handleResourcesUnsubscribe removes a resource from change notifications.
//
// Parameters:
//   - req: the MCP request containing the resource URI
//
// Returns:
//   - *proto.Response: empty success or validation error
func (s *Server) handleResourcesUnsubscribe(req proto.Request) *proto.Response {
	return s.subscriptionAction(req, s.poller.unsubscribe)
}

// subscriptionAction unmarshals a URI param and applies the given
// subscription function (subscribe or unsubscribe).
//
// Parameters:
//   - req: the MCP request containing the resource URI
//   - fn: subscription action to apply (subscribe or unsubscribe)
//
// Returns:
//   - *proto.Response: empty success or validation error
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
