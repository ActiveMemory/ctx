//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package poll

import (
	"sync"
	"time"

	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Poller tracks subscribed resources and polls for file
// changes.
type Poller struct {
	mu         sync.Mutex
	subs       map[string]bool      // URI → subscribed
	mtimes     map[string]time.Time // file path → last known mtime
	contextDir string
	pollStop   chan struct{}
	notifyFunc func(proto.Notification) // callback to emit notifications
}
