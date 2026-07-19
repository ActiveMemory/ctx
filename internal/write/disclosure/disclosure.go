//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/disclosure"
)

// Human prints an Inspection as a scannable summary: the kind, the
// staged entries awaiting digestion, and the current themes.
//
// Parameters:
//   - cmd: Cobra command for the output stream
//   - insp: the inspection to render
func Human(cmd *cobra.Command, insp disclosure.Inspection) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteDisclosureKind), insp.Kind))

	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteDisclosureStagedHeader), len(insp.Staging),
	))
	if len(insp.Staging) == 0 {
		cmd.Println(desc.Text(text.DescKeyWriteDisclosureNone))
	}
	for _, e := range insp.Staging {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDisclosureStagedLine),
			e.Timestamp, e.Title,
		))
	}

	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteDisclosureThemesHeader), len(insp.Themes),
	))
	if len(insp.Themes) == 0 {
		cmd.Println(desc.Text(text.DescKeyWriteDisclosureNone))
	}
	for _, t := range insp.Themes {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteDisclosureThemeLine), t.Name, t.Link,
		))
	}
}

// JSON prints an Inspection as indented JSON for machine consumption.
//
// Parameters:
//   - cmd: Cobra command for the output stream
//   - insp: the inspection to render
//
// Returns:
//   - error: non-nil only if JSON marshaling fails
func JSON(cmd *cobra.Command, insp disclosure.Inspection) error {
	b, err := json.MarshalIndent(insp, "", token.Space+token.Space)
	if err != nil {
		return err
	}
	cmd.Println(string(b))
	return nil
}

// ApplyHuman prints an ApplyResult as a one-line summary: how many
// entries moved, into how many themes, and which theme slugs.
//
// Parameters:
//   - cmd: Cobra command for the output stream
//   - res: the apply result to render
func ApplyHuman(cmd *cobra.Command, res disclosure.ApplyResult) {
	slugs := strings.Join(res.Themes, token.CommaSpace)
	if slugs == "" {
		slugs = token.Dash
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteDisclosureApplied),
		res.Moved, len(res.Themes), slugs,
	))
}

// ApplyJSON prints an ApplyResult as indented JSON.
//
// Parameters:
//   - cmd: Cobra command for the output stream
//   - res: the apply result to render
//
// Returns:
//   - error: non-nil only if JSON marshaling fails
func ApplyJSON(cmd *cobra.Command, res disclosure.ApplyResult) error {
	b, err := json.MarshalIndent(res, "", token.Space+token.Space)
	if err != nil {
		return err
	}
	cmd.Println(string(b))
	return nil
}
