//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// baseLevel is the shallowest heading level the projector emits; deeper
// levels are indented relative to it in line output.
const baseLevel = 2

// Lines prints one heading per line, indenting deeper (L3) headings under
// their parents for a scannable table of contents.
//
// Parameters:
//   - cmd: Cobra command for the output stream.
//   - headings: Projected headings in file order.
func Lines(cmd *cobra.Command, headings []entity.Heading) {
	for _, h := range headings {
		depth := h.Level - baseLevel
		if depth < 0 {
			depth = 0
		}
		cmd.Println(strings.Repeat(token.Space+token.Space, depth) + h.Text)
	}
}

// JSON prints the headings as a JSON array of {level, text} objects.
//
// A file with no headings yields "[]", not "null", so consumers can parse
// unconditionally.
//
// Parameters:
//   - cmd: Cobra command for the output stream.
//   - headings: Projected headings in file order.
//
// Returns:
//   - error: Non-nil only if JSON marshaling fails.
func JSON(cmd *cobra.Command, headings []entity.Heading) error {
	if headings == nil {
		headings = []entity.Heading{}
	}
	b, err := json.MarshalIndent(headings, "", token.Space+token.Space)
	if err != nil {
		return err
	}
	cmd.Println(string(b))
	return nil
}
