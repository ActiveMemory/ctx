//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package topic

import (
	"io/fs"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	kbPath "github.com/ActiveMemory/ctx/internal/cli/kb/core/path"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/slug"
)

// Scaffold creates the topic folder and writes the index.md
// rendered from the embedded template.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - name: free-text topic name (slugified to kebab-case).
//
// Returns:
//   - error: scaffolding failure or refusal when topic exists.
func Scaffold(cobraCmd *cobra.Command, name string) error {
	topicSlug := slug.Path(name)
	if topicSlug == "" {
		cobraCmd.SilenceUsage = true
		return errKbCli.ErrTopicEmptyName
	}
	topicDir, dirErr := kbPath.KBTopicDir(topicSlug)
	if dirErr != nil {
		return dirErr
	}
	indexPath := filepath.Join(topicDir, cfgKB.TopicIndex)
	if _, statErr := io.SafeStat(indexPath); statErr == nil {
		cobraCmd.SilenceUsage = true
		return errKbCli.TopicExists(topicSlug, indexPath)
	}

	mkErr := io.SafeMkdirAll(topicDir, cfgFs.PermExec)
	if mkErr != nil {
		return errKbCli.MkdirTopic(mkErr)
	}
	raw, readErr := fs.ReadFile(
		assets.FS, cfgKB.AssetTemplateTopicIndex,
	)
	if readErr != nil {
		return errKbCli.ReadTopicTemplate(readErr)
	}
	rendered := Substitute(string(raw), name, topicSlug)
	writeErr := io.SafeWriteFile(
		indexPath,
		[]byte(rendered),
		cfgFs.PermSecret,
	)
	if writeErr != nil {
		return errKbCli.WriteTopicIndex(writeErr)
	}
	io.SafeFprintf(
		cobraCmd.OutOrStdout(),
		desc.Text(text.DescKeyWriteKbCliScaffolded),
		indexPath,
	)
	return nil
}
