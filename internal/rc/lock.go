//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import "sync"

// rc, rcOnce, and rcMu hold the singleton runtime configuration
// loaded once from .ctxrc via sync.Once.
var (
	// rc holds the singleton runtime configuration.
	rc *CtxRC
	// rcOnce guards one-time configuration loading.
	rcOnce sync.Once
	// rcMu protects concurrent access to rc.
	rcMu sync.RWMutex
)
