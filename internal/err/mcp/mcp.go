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
		desc.Text(text.DescKeyMCPErrTypeContentRequired),
	)
}

// QueryRequired returns an error when query is missing from a search
// tool call.
//
// Returns:
//   - error: "query is required"
func QueryRequired() error {
	return errors.New(
		desc.Text(text.DescKeyMCPErrQueryRequired),
	)
}

// SearchRead wraps a failure to read the context directory during
// search.
//
// Parameters:
//   - dir: the directory path
//   - cause: the underlying read error
//
// Returns:
//   - error: "search: read <dir>: <cause>"
func SearchRead(dir string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyMCPErrSearchRead), dir, cause,
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
		desc.Text(text.DescKeyMCPUnknownEventType),
		eventType,
	)
}

// InputTooLong returns an error when input exceeds the allowed
// length.
//
// Parameters:
//   - field: the field name that is too long
//   - maxLen: the maximum allowed length
//
// Returns:
//   - error: "<field> exceeds maximum length of <maxLen>"
func InputTooLong(field string, maxLen int) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyMCPErrInputTooLong),
		field, maxLen,
	)
}
