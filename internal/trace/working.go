//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/env"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
)

// WorkingRefs detects context refs from the current working state.
//
// It combines in-progress task refs from TASKS.md with an active AI session
// ref (if CTX_SESSION_ID is set).
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - []string: refs like "task:1", "session:<id>"
func WorkingRefs(contextDir string) []string {
	var refs []string

	refs = append(refs, inProgressTaskRefs(contextDir)...)

	if id := os.Getenv(env.SessionID); id != "" {
		refs = append(refs, fmt.Sprintf(cfgTrace.SessionRefFormat, id))
	}

	return refs
}
