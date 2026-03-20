//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// TypeContentRequired returns an error when type or content is missing
// from an MCP tool call.
//
// Returns:
//   - error: "type and content are required"
func TypeContentRequired() error {
	return errors.New(
		desc.TextDesc(text.DescKeyMCPErrTypeContentRequired),
	)
}

// UnknownEventType returns an error for an unrecognized session event
// type.
//
// Parameters:
//   - eventType: the unrecognized event type string
//
// Returns:
//   - error: "unknown event type: <eventType>"
func UnknownEventType(eventType string) error {
	return fmt.Errorf(
		desc.TextDesc(text.TextDescKeyMCPUnknownEventType),
		eventType,
	)
}
