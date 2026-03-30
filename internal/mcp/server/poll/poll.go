//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package poll

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/mcp/notify"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/catalog"
)

// defaultPollInterval is the default interval for resource change
// polling.
const defaultPollInterval = 5 * time.Second

// NewPoller creates a poller for the given context directory.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - notifyFn: callback to emit resource change notifications
//
// Returns:
//   - *Poller: initialized poller (not yet polling)
func NewPoller(
	contextDir string, notifyFn func(proto.Notification),
) *Poller {
	return &Poller{
		subs:       make(map[string]bool),
		mtimes:     make(map[string]time.Time),
		contextDir: contextDir,
		notifyFunc: notifyFn,
	}
}

// Subscribe adds a URI to the watch set and starts polling if needed.
//
// Goroutine lifecycle: the poller goroutine is started on the first
// subscription and stopped when the last subscription is removed or
// when Stop is called.
//
// Parameters:
//   - uri: resource URI to watch for changes
func (p *Poller) Subscribe(uri string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subs[uri] = true

	// Snapshot current mtime for the resource's file.
	if fileName := catalog.FileForURI(uri); fileName != "" {
		fPath := filepath.Join(p.contextDir, fileName)
		if info, statErr := os.Stat(fPath); statErr == nil {
			p.mtimes[fPath] = info.ModTime()
		}
	}

	// Start poller if this is the first subscription.
	if len(p.subs) == 1 && p.pollStop == nil {
		p.pollStop = make(chan struct{})
		go p.poll()
	}
}

// Unsubscribe removes a URI from the watch set and stops polling if
// empty.
//
// Parameters:
//   - uri: resource URI to stop watching
func (p *Poller) Unsubscribe(uri string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subs, uri)

	if len(p.subs) == 0 && p.pollStop != nil {
		close(p.pollStop)
		p.pollStop = nil
	}
}

// Stop shuts down the poller goroutine.
func (p *Poller) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pollStop != nil {
		close(p.pollStop)
		p.pollStop = nil
	}
}

// poll checks subscribed resources for mtime changes on a fixed
// interval.
func (p *Poller) poll() {
	ticker := time.NewTicker(defaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.pollStop:
			return
		case <-ticker.C:
			p.CheckChanges()
		}
	}
}

// SetNotifyFunc replaces the notification callback. Intended for
// tests that need to capture emitted notifications.
//
// Parameters:
//   - fn: replacement notification callback
func (p *Poller) SetNotifyFunc(fn func(proto.Notification)) {
	p.notifyFunc = fn
}

// CheckChanges compares current mtimes to snapshots and emits
// notifications. Exported for tests that need to trigger a poll
// cycle without waiting for the timer.
func (p *Poller) CheckChanges() {
	p.mu.Lock()
	uris := make([]string, 0, len(p.subs))
	for uri := range p.subs {
		uris = append(uris, uri)
	}
	p.mu.Unlock()

	for _, uri := range uris {
		fileName := catalog.FileForURI(uri)
		if fileName == "" {
			continue
		}

		fPath := filepath.Join(p.contextDir, fileName)
		info, statErr := os.Stat(fPath)
		if statErr != nil {
			continue
		}

		p.mu.Lock()
		prev, known := p.mtimes[fPath]
		if known && info.ModTime().After(prev) {
			p.mtimes[fPath] = info.ModTime()
			p.mu.Unlock()
			p.notifyFunc(proto.Notification{
				JSONRPC: server.JSONRPCVersion,
				Method:  notify.ResourcesUpdated,
				Params:  proto.ResourceUpdatedParams{URI: uri},
			})
		} else {
			if !known {
				p.mtimes[fPath] = info.ModTime()
			}
			p.mu.Unlock()
		}
	}
}
