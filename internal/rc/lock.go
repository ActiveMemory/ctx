//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import "sync"

var (
	rc            *CtxRC
	rcOnce        sync.Once
	rcOverrideDir string
	rcMu          sync.RWMutex
)
