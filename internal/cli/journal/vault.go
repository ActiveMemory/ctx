//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// obsidianMaxRelated is the maximum number of "see also" entries in the
// related sessions footer.
const obsidianMaxRelated = 5

// runJournalObsidian generates an Obsidian vault from journal entries.
//
// Pipeline:
//  1. Scan entries (reuse scanJournalEntries)
//  2. Create output dirs (entries/, topics/, files/, types/, .obsidian/)
//  3. Write .obsidian/app.json
//  4. Transform and write entries (normalize, convert links, transform
//     frontmatter, add related footer)
//  5. Build indices (reuse buildTopicIndex etc.)
//  6. Generate and write MOC pages
//  7. Generate and write Home.md
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - output: Output directory for the vault
//
// Returns:
//   - error: Non-nil if generation fails
func runJournalObsidian(cmd *cobra.Command, output string) error {
	return buildObsidianVault(cmd, filepath.Join(rc.ContextDir(), config.DirJournal), output)
}

// buildObsidianVault generates an Obsidian vault from journal entries in
// journalDir and writes the output to the output directory.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - journalDir: Path to the source journal directory
//   - output: Output directory for the vault
//
// Returns:
//   - error: Non-nil if generation fails
func buildObsidianVault(cmd *cobra.Command, journalDir, output string) error {
	if _, err := os.Stat(journalDir); os.IsNotExist(err) {
		return errNoJournalDir(journalDir)
	}

	entries, scanErr := scanJournalEntries(journalDir)
	if scanErr != nil {
		return errScanJournal(scanErr)
	}

	if len(entries) == 0 {
		return errNoEntries(journalDir)
	}

	green := color.New(color.FgGreen).SprintFunc()

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
			return errMkdir(dir, mkErr)
		}
	}

	// Write .obsidian/app.json
	appConfigPath := filepath.Join(
		output, config.ObsidianConfigDir, config.ObsidianAppConfigFile,
	)
	if writeErr := os.WriteFile(
		appConfigPath, []byte(config.ObsidianAppConfig), config.PermFile,
	); writeErr != nil {
		return errFileWrite(appConfigPath, writeErr)
	}

	// Write README
	readmePath := filepath.Join(output, config.FilenameReadme)
	if writeErr := os.WriteFile(
		readmePath,
		[]byte(fmt.Sprintf(config.ObsidianReadme, journalDir)),
		config.PermFile,
	); writeErr != nil {
		return errFileWrite(readmePath, writeErr)
	}

	// Build indices for MOC pages and related footer
	regularEntries := filterRegularEntries(entries)

	topicEntries := filterEntriesWithTopics(entries)
	topics := buildTopicIndex(topicEntries)

	keyFileEntries := filterEntriesWithKeyFiles(entries)
	keyFiles := buildKeyFileIndex(keyFileEntries)

	typeEntries := filterEntriesWithType(entries)
	sessionTypes := buildTypeIndex(typeEntries)

	// Build topic lookup for related footer
	topicIndex := buildTopicLookup(topicEntries)

	// Transform and write entries
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(output, config.ObsidianDirEntries, entry.Filename)

		content, readErr := os.ReadFile(filepath.Clean(src))
		if readErr != nil {
			warnFileErr(cmd, entry.Filename, readErr)
			continue
		}

		// Normalize content (read-only — do NOT write back to source)
		normalized := softWrapContent(
			mergeConsecutiveTurns(
				consolidateToolRuns(
					cleanToolOutputJSON(
						stripSystemReminders(string(content)),
					),
				),
			),
		)

		// Transform for Obsidian
		sourcePath := filepath.Join(
			config.DirContext, config.DirJournal, entry.Filename,
		)
		transformed := transformFrontmatter(normalized, sourcePath)
		transformed = convertMarkdownLinks(transformed)
		transformed += generateRelatedFooter(entry, topicIndex, obsidianMaxRelated)

		if writeErr := os.WriteFile(
			dst, []byte(transformed), config.PermFile,
		); writeErr != nil {
			warnFileErr(cmd, entry.Filename, writeErr)
			continue
		}
	}

	// Write topic MOC and pages
	if len(topics) > 0 {
		topicsDir := filepath.Join(output, config.JournalDirTopics)
		mocPath := filepath.Join(output, config.ObsidianTopicsMOC)
		if writeErr := os.WriteFile(
			mocPath, []byte(generateObsidianTopicsMOC(topics)),
			config.PermFile,
		); writeErr != nil {
			return errFileWrite(mocPath, writeErr)
		}

		for _, t := range topics {
			if !t.Popular {
				continue
			}
			pagePath := filepath.Join(topicsDir, t.Name+config.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(generateObsidianTopicPage(t)),
				config.PermFile,
			); writeErr != nil {
				warnFileErr(cmd, pagePath, writeErr)
			}
		}
	}

	// Write key files MOC and pages
	if len(keyFiles) > 0 {
		filesDir := filepath.Join(output, config.JournalDirFiles)
		mocPath := filepath.Join(output, config.ObsidianFilesMOC)
		if writeErr := os.WriteFile(
			mocPath, []byte(generateObsidianFilesMOC(keyFiles)),
			config.PermFile,
		); writeErr != nil {
			return errFileWrite(mocPath, writeErr)
		}

		for _, kf := range keyFiles {
			if !kf.Popular {
				continue
			}
			slug := keyFileSlug(kf.Path)
			pagePath := filepath.Join(filesDir, slug+config.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(generateObsidianFilePage(kf)),
				config.PermFile,
			); writeErr != nil {
				warnFileErr(cmd, pagePath, writeErr)
			}
		}
	}

	// Write types MOC and pages
	if len(sessionTypes) > 0 {
		typesDir := filepath.Join(output, config.JournalDirTypes)
		mocPath := filepath.Join(output, config.ObsidianTypesMOC)
		if writeErr := os.WriteFile(
			mocPath, []byte(generateObsidianTypesMOC(sessionTypes)),
			config.PermFile,
		); writeErr != nil {
			return errFileWrite(mocPath, writeErr)
		}

		for _, st := range sessionTypes {
			pagePath := filepath.Join(typesDir, st.Name+config.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(generateObsidianTypePage(st)),
				config.PermFile,
			); writeErr != nil {
				warnFileErr(cmd, pagePath, writeErr)
			}
		}
	}

	// Write Home.md
	homePath := filepath.Join(output, config.ObsidianHomeMOC)
	if writeErr := os.WriteFile(
		homePath,
		[]byte(generateHomeMOC(
			regularEntries,
			len(topics) > 0, len(keyFiles) > 0, len(sessionTypes) > 0,
		)),
		config.PermFile,
	); writeErr != nil {
		return errFileWrite(homePath, writeErr)
	}

	cmd.Println(fmt.Sprintf(
		"%s Generated Obsidian vault with %d entries in %s",
		green("✓"), len(entries), output,
	))
	cmd.Println()
	cmd.Println("Next steps:")
	cmd.Println("  Open Obsidian → Open folder as vault → Select " + output)

	return nil
}

// filterRegularEntries returns entries excluding suggestions and multipart
// continuations.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []journalEntry: Filtered entries
func filterRegularEntries(entries []journalEntry) []journalEntry {
	var result []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// filterEntriesWithTopics returns non-suggestive, non-multipart entries
// that have topics.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []journalEntry: Entries with topics
func filterEntriesWithTopics(entries []journalEntry) []journalEntry {
	var result []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) || len(e.Topics) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

// filterEntriesWithKeyFiles returns non-suggestive, non-multipart entries
// that have key files.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []journalEntry: Entries with key files
func filterEntriesWithKeyFiles(entries []journalEntry) []journalEntry {
	var result []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) || len(e.KeyFiles) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

// filterEntriesWithType returns non-suggestive, non-multipart entries
// that have a type.
//
// Parameters:
//   - entries: All journal entries
//
// Returns:
//   - []journalEntry: Entries with type
func filterEntriesWithType(entries []journalEntry) []journalEntry {
	var result []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) || e.Type == "" {
			continue
		}
		result = append(result, e)
	}
	return result
}

// buildTopicLookup creates a map from topic name to all entries with
// that topic, for efficient related-entry lookups.
//
// Parameters:
//   - entries: Entries with topics
//
// Returns:
//   - map[string][]journalEntry: Topic name → entries
func buildTopicLookup(entries []journalEntry) map[string][]journalEntry {
	lookup := make(map[string][]journalEntry)
	for _, e := range entries {
		for _, topic := range e.Topics {
			lookup[topic] = append(lookup[topic], e)
		}
	}
	return lookup
}
