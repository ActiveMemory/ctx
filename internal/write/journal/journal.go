//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoOrphanRemoved reports a removed orphan file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Filename that was removed
func InfoOrphanRemoved(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteJournalOrphanRemoved),
		name))
}

// InfoSiteGenerated reports successful site generation with next steps.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of entries generated
//   - output: Output directory path
//   - zensicalBin: Zensical binary name
func InfoSiteGenerated(
	cmd *cobra.Command, count int, output, zensicalBin string,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteJournalSiteGeneratedBlock),
		count, output, output, zensicalBin,
	))
}

// InfoSiteStarting reports the server is starting.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSiteStarting(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteJournalSiteStarting))
}

// InfoSiteBuilding reports a build is in progress.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSiteBuilding(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteJournalSiteBuilding))
}
