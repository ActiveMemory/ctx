//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import "time"

// copilotCLIRawMessage represents a single JSONL line from a Copilot CLI
// session file. The exact format may evolve as Copilot CLI matures.
type copilotCLIRawMessage struct {
	ID        string    `json:"id,omitempty"`
	SessionID string    `json:"sessionId,omitempty"`
	Role      string    `json:"role,omitempty"`
	Type      string    `json:"type,omitempty"`
	Text      string    `json:"text,omitempty"`
	Model     string    `json:"model,omitempty"`
	CWD       string    `json:"cwd,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
