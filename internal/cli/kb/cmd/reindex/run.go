//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	kbPath "github.com/ActiveMemory/ctx/internal/cli/kb/core/path"
	reindexCore "github.com/ActiveMemory/ctx/internal/cli/kb/core/reindex"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run rebuilds the managed block from the current topic
// folders.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//
// Returns:
//   - error: I/O or parse failure.
func Run(cobraCmd *cobra.Command) error {
	indexPath, pathErr := kbPath.KBIndexFile()
	if pathErr != nil {
		return pathErr
	}
	topicsDir, topicsErr := kbPath.KBTopicsDir()
	if topicsErr != nil {
		return topicsErr
	}

	slugs, listErr := reindexCore.ListTopics(topicsDir)
	if listErr != nil {
		return listErr
	}

	raw, readErr := io.SafeReadUserFile(indexPath)
	if readErr != nil {
		return errKbCli.ReadKBIndex(readErr)
	}
	if !regex.ManagedKBTopics.MatchString(string(raw)) {
		cobraCmd.SilenceUsage = true
		return errKbCli.ErrReindexMissingBlock
	}
	block := reindexCore.RenderBlock(slugs)
	updated := regex.ManagedKBTopics.ReplaceAllString(
		string(raw), block,
	)
	writeErr := io.SafeWriteFile(
		indexPath,
		[]byte(updated),
		cfgFs.PermSecret,
	)
	if writeErr != nil {
		return errKbCli.WriteKBIndex(writeErr)
	}
	io.SafeFprintf(
		cobraCmd.OutOrStdout(),
		desc.Text(text.DescKeyWriteKbReindexed),
		len(slugs),
		indexPath,
	)
	return nil
}
