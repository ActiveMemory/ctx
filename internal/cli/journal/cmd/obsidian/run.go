//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// ObsidianMaxRelated is the maximum number of "see also" entries in the
// related sessions footer.
const ObsidianMaxRelated = 5

// runJournalObsidian generates an Obsidian vault from journal entries.
//
// Pipeline:
//  1. Scan entries (reuse core.ScanJournalEntries)
//  2. Create output dirs (entries/, topics/, files/, types/, .obsidian/)
//  3. Write .obsidian/app.json
//  4. Transform and write entries (normalize, convert links, transform
//     frontmatter, add related footer)
//  5. Build indices (reuse core.BuildTopicIndex etc.)
//  6. Generate and write MOC pages
//  7. Generate and write Home.md
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - output: Output directory for the vault
//
// Returns:
//   - error: Non-nil if generation fails
func Run(cmd *cobra.Command, output string) error {
	return BuildObsidianVault(cmd, filepath.Join(rc.ContextDir(), config.DirJournal), output)
}

// BuildObsidianVault generates an Obsidian vault from journal entries in
// journalDir and writes the output to the output directory.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - journalDir: Path to the source journal directory
//   - output: Output directory for the vault
//
// Returns:
//   - error: Non-nil if generation fails
func BuildObsidianVault(cmd *cobra.Command, journalDir, output string) error {
	if _, statErr := os.Stat(journalDir); os.IsNotExist(statErr) {
		return ctxerr.NoJournalDir(journalDir)
	}

	entries, scanErr := core.ScanJournalEntries(journalDir)
	if scanErr != nil {
		return ctxerr.ScanJournal(scanErr)
	}

	if len(entries) == 0 {
		return ctxerr.NoJournalEntries(journalDir)
	}

	// Create output directory structure
	dirs := []string{
		output,
		filepath.Join(output, config.ObsidianDirEntries),
		filepath.Join(output, config.ObsidianConfigDir),
		filepath.Join(output, config.JournalDirTopics),
		filepath.Join(output, config.JournalDirFiles),
		filepath.Join(output, config.JournalDirTypes),
	}
	for _, dir := range dirs {
		if mkErr := os.MkdirAll(dir, config.PermExec); mkErr != nil {
			return ctxerr.Mkdir(dir, mkErr)
		}
	}

	// Write .obsidian/app.json
	appConfigPath := filepath.Join(
		output, config.ObsidianConfigDir, config.ObsidianAppConfigFile,
	)
	if writeErr := os.WriteFile(
		appConfigPath, []byte(config.ObsidianAppConfig), config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(appConfigPath, writeErr)
	}

	// Write README
	readmePath := filepath.Join(output, config.FilenameReadme)
	if writeErr := os.WriteFile(
		readmePath,
		[]byte(fmt.Sprintf(config.ObsidianReadme, journalDir)),
		config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(readmePath, writeErr)
	}

	// Build indices for MOC pages and related footer
	regularEntries := core.FilterRegularEntries(entries)

	topicEntries := core.FilterEntriesWithTopics(entries)
	topics := core.BuildTopicIndex(topicEntries)

	keyFileEntries := core.FilterEntriesWithKeyFiles(entries)
	keyFiles := core.BuildKeyFileIndex(keyFileEntries)

	typeEntries := core.FilterEntriesWithType(entries)
	sessionTypes := core.BuildTypeIndex(typeEntries)

	// Build topic lookup for related footer
	topicIndex := core.BuildTopicLookup(topicEntries)

	// Transform and write entries
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(output, config.ObsidianDirEntries, entry.Filename)

		content, readErr := os.ReadFile(filepath.Clean(src))
		if readErr != nil {
			write.WarnFileErr(cmd, entry.Filename, readErr)
			continue
		}

		// Normalize content (read-only — do NOT write back to source)
		normalized := core.SoftWrapContent(
			core.MergeConsecutiveTurns(
				core.ConsolidateToolRuns(
					core.CleanToolOutputJSON(
						core.StripSystemReminders(string(content)),
					),
				),
			),
		)

		// Transform for Obsidian
		sourcePath := filepath.Join(
			config.DirContext, config.DirJournal, entry.Filename,
		)
		transformed := core.TransformFrontmatter(normalized, sourcePath)
		transformed = core.ConvertMarkdownLinks(transformed)
		transformed += core.GenerateRelatedFooter(entry, topicIndex, ObsidianMaxRelated)

		if writeErr := os.WriteFile(
			dst, []byte(transformed), config.PermFile,
		); writeErr != nil {
			write.WarnFileErr(cmd, entry.Filename, writeErr)
			continue
		}
	}

	// Write topic MOC and pages
	if len(topics) > 0 {
		topicsDir := filepath.Join(output, config.JournalDirTopics)
		mocPath := filepath.Join(output, config.ObsidianTopicsMOC)
		if writeErr := os.WriteFile(
			mocPath, []byte(core.GenerateObsidianTopicsMOC(topics)),
			config.PermFile,
		); writeErr != nil {
			return ctxerr.FileWrite(mocPath, writeErr)
		}

		for _, t := range topics {
			if !t.Popular {
				continue
			}
			pagePath := filepath.Join(topicsDir, t.Name+config.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(core.GenerateObsidianTopicPage(t)),
				config.PermFile,
			); writeErr != nil {
				write.WarnFileErr(cmd, pagePath, writeErr)
			}
		}
	}

	// Write key files MOC and pages
	if len(keyFiles) > 0 {
		filesDir := filepath.Join(output, config.JournalDirFiles)
		mocPath := filepath.Join(output, config.ObsidianFilesMOC)
		if writeErr := os.WriteFile(
			mocPath, []byte(core.GenerateObsidianFilesMOC(keyFiles)),
			config.PermFile,
		); writeErr != nil {
			return ctxerr.FileWrite(mocPath, writeErr)
		}

		for _, kf := range keyFiles {
			if !kf.Popular {
				continue
			}
			slug := core.KeyFileSlug(kf.Path)
			pagePath := filepath.Join(filesDir, slug+config.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(core.GenerateObsidianFilePage(kf)),
				config.PermFile,
			); writeErr != nil {
				write.WarnFileErr(cmd, pagePath, writeErr)
			}
		}
	}

	// Write types MOC and pages
	if len(sessionTypes) > 0 {
		typesDir := filepath.Join(output, config.JournalDirTypes)
		mocPath := filepath.Join(output, config.ObsidianTypesMOC)
		if writeErr := os.WriteFile(
			mocPath, []byte(core.GenerateObsidianTypesMOC(sessionTypes)),
			config.PermFile,
		); writeErr != nil {
			return ctxerr.FileWrite(mocPath, writeErr)
		}

		for _, st := range sessionTypes {
			pagePath := filepath.Join(typesDir, st.Name+config.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath,
				[]byte(core.GenerateObsidianTypePage(st)), config.PermFile,
			); writeErr != nil {
				write.WarnFileErr(cmd, pagePath, writeErr)
			}
		}
	}

	// Write Home.md
	homePath := filepath.Join(output, config.ObsidianHomeMOC)
	if writeErr := os.WriteFile(
		homePath,
		[]byte(core.GenerateHomeMOC(
			regularEntries,
			len(topics) > 0, len(keyFiles) > 0, len(sessionTypes) > 0,
		)),
		config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(homePath, writeErr)
	}

	write.InfoObsidianGenerated(cmd, len(entries), output)

	return nil
}
