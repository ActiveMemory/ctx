//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// maxContentLen is the maximum entry content size (1MB).
// Entries are text-only (decisions, learnings, conventions).
const maxContentLen = 1 << 20

// allowedTypes is the set of valid entry types.
var allowedTypes = map[string]bool{
	"decision":   true,
	"learning":   true,
	"convention": true,
	"task":       true,
}

// validateEntry checks a PublishEntry for required fields
// and enforces size limits.
func validateEntry(pe PublishEntry) error {
	if pe.ID == "" {
		return status.Error(
			codes.InvalidArgument, "entry ID required",
		)
	}
	if !allowedTypes[pe.Type] {
		return status.Errorf(
			codes.InvalidArgument,
			"invalid entry type %q", pe.Type,
		)
	}
	if pe.Origin == "" {
		return status.Error(
			codes.InvalidArgument,
			"entry origin required",
		)
	}
	if len(pe.Content) > maxContentLen {
		return status.Error(
			codes.InvalidArgument,
			"entry content exceeds 1MB limit",
		)
	}
	return nil
}
