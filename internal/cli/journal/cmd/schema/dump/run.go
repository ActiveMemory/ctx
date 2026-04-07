//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dump

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	coreSchema "github.com/ActiveMemory/ctx/internal/cli/journal/core/schema"
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/journal/schema"
	writeSchema "github.com/ActiveMemory/ctx/internal/write/schema"
)

// Run executes the schema dump command.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: always nil
func Run(cmd *cobra.Command) error {
	s := schema.Default()

	writeSchema.DumpLine(cmd, fmt.Sprintf(
		cfgSchema.FmtDumpVersion, s.Version))
	writeSchema.DumpLine(cmd, fmt.Sprintf(
		cfgSchema.FmtDumpCCRange, s.CCVersionRange))
	writeSchema.DumpBlank(cmd)

	types := coreSchema.SortedRecordTypes(s.RecordTypes)

	writeSchema.DumpLine(cmd, cfgSchema.HeadingRecordTypes)
	writeSchema.DumpBlank(cmd)
	for _, t := range types {
		rs := s.RecordTypes[t]
		if len(rs.Required) == 0 && len(rs.Optional) == 0 {
			writeSchema.DumpLine(cmd, fmt.Sprintf(
				cfgSchema.FmtDumpMetadata, t))
			continue
		}
		writeSchema.DumpLine(cmd, fmt.Sprintf(
			cfgSchema.FmtDumpRecordType, t))
		writeSchema.DumpLine(cmd, fmt.Sprintf(
			cfgSchema.FmtDumpRequired,
			strings.Join(
				rs.Required, token.CommaSpace)))
		writeSchema.DumpLine(cmd, fmt.Sprintf(
			cfgSchema.FmtDumpOptional,
			strings.Join(
				rs.Optional, token.CommaSpace)))
	}

	writeSchema.DumpBlank(cmd)
	writeSchema.DumpLine(cmd, cfgSchema.HeadingBlockTypes)
	writeSchema.DumpBlank(cmd)

	blockTypes := coreSchema.SortedBlockTypes(s.BlockTypes)
	for _, bt := range blockTypes {
		kind := s.BlockTypes[bt]
		label := cfgSchema.LabelParsed
		if kind == schema.BlockKnown {
			label = cfgSchema.LabelKnown
		}
		writeSchema.DumpLine(cmd, fmt.Sprintf(
			cfgSchema.FmtDumpBlock, bt, label))
	}

	return nil
}
