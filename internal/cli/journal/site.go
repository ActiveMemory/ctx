//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// multiPartPattern matches session part files like "...-p2.md", "...-p3.md", etc.
var multiPartPattern = regexp.MustCompile(`-p\d+\.md$`)

// globStarPattern matches glob-like wildcards: *.ext, */, *) etc.
var globStarPattern = regexp.MustCompile(`\*(\.\w+|[/)])`)

// toolBoldPattern matches tool-use lines like "🔧 **Glob: .context/sessions/*.md**"
var toolBoldPattern = regexp.MustCompile(`🔧\s*\*\*(.+?)\*\*`)

// turnHeaderPattern matches conversation turn headers like "### 30. Assistant (04:41:50)"
// or "### 33. Tool Output (04:41:58)"
var turnHeaderPattern = regexp.MustCompile(`^### (\d+)\. (.+?) \((\d{2}:\d{2}:\d{2})\)$`)

// fenceLinePattern matches lines that are code fence markers (3+ backticks or tildes,
// optionally followed by a language tag).
var fenceLinePattern = regexp.MustCompile("^\\s*(`{3,}|~{3,})(.*)$")

// normalizedMarkerPattern matches the metadata normalization marker (normalize.py).
var normalizedMarkerPattern = regexp.MustCompile(`<!-- normalized: \d{4}-\d{2}-\d{2} -->`)

// fencesVerifiedPattern matches the marker left after AI fence reconstruction.
// Only files with this marker skip fence stripping in the site pipeline.
var fencesVerifiedPattern = regexp.MustCompile(`<!-- fences-verified: \d{4}-\d{2}-\d{2} -->`)



// journalSiteCmd returns the journal site subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating a static site from journal entries
func journalSiteCmd() *cobra.Command {
	var (
		output string
		serve  bool
		build  bool
	)

	cmd := &cobra.Command{
		Use:   "site",
		Short: "Generate a static site from journal entries",
		Long: `Generate a zensical-compatible static site from .context/journal/ entries.

Creates a site structure with:
  - Index page with all sessions listed by date
  - Individual pages for each journal entry
  - Navigation and search support

Requires zensical to be installed for building/serving:
  pip install zensical

Examples:
  ctx journal site                    # Generate in .context/journal-site/
  ctx journal site --output ~/public  # Custom output directory
  ctx journal site --build            # Generate and build HTML
  ctx journal site --serve            # Generate and serve locally`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalSite(cmd, output, build, serve)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), "journal-site")
	cmd.Flags().StringVarP(&output, "output", "o", defaultOutput, "Output directory for site")
	cmd.Flags().BoolVar(&build, "build", false, "Run zensical build after generating")
	cmd.Flags().BoolVar(&serve, "serve", false, "Run zensical serve after generating")

	return cmd
}

// journalFrontmatter represents YAML frontmatter in enriched journal entries.
type journalFrontmatter struct {
	Title    string   `yaml:"title"`
	Date     string   `yaml:"date"`
	Type     string   `yaml:"type"`
	Outcome  string   `yaml:"outcome"`
	Topics   []string `yaml:"topics"`
	KeyFiles []string `yaml:"key_files"`
}

// journalEntry represents a parsed journal file.
type journalEntry struct {
	Filename     string
	Title        string
	Date         string
	Time         string
	Project      string
	Path         string
	Size         int64
	IsSuggestion bool
	Topics       []string
	Type         string
	Outcome      string
	KeyFiles     []string
}

// runJournalSite handles the journal site command.
//
// Scans .context/journal/ for markdown files, generates a zensical project
// structure, and optionally builds or serves the site.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - output: Output directory for the generated site
//   - build: If true, run zensical build after generating
//   - serve: If true, run zensical serve after generating
//
// Returns:
//   - error: Non-nil if generation fails
func runJournalSite(cmd *cobra.Command, output string, build, serve bool) error {
	journalDir := filepath.Join(rc.ContextDir(), "journal")

	// Check if journal directory exists
	if _, err := os.Stat(journalDir); os.IsNotExist(err) {
		return fmt.Errorf("no journal directory found at %s\nRun 'ctx recall export --all' first", journalDir)
	}

	// Scan journal files
	entries, err := scanJournalEntries(journalDir)
	if err != nil {
		return fmt.Errorf("failed to scan journal: %w", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("no journal entries found in %s\nRun 'ctx recall export --all' first", journalDir)
	}

	green := color.New(color.FgGreen).SprintFunc()

	// Create output directory structure
	docsDir := filepath.Join(output, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write README
	readmePath := filepath.Join(output, "README.md")
	if err := os.WriteFile(readmePath, []byte(generateSiteReadme(journalDir)), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Soft-wrap source journal files in-place, then copy to docs/
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(docsDir, entry.Filename)

		content, err := os.ReadFile(src)
		if err != nil {
			cmd.PrintErrf("  ! failed to read %s: %v\n", entry.Filename, err)
			continue
		}

		// Normalize source file for readability
		normalized := softWrapContent(mergeConsecutiveTurns(consolidateToolRuns(cleanToolOutputJSON(stripSystemReminders(string(content))))))
		if normalized != string(content) {
			if err := os.WriteFile(src, []byte(normalized), 0644); err != nil {
				cmd.PrintErrf("  ! failed to normalize %s: %v\n", entry.Filename, err)
			}
		}

		// Generate site copy with markdown fixes
		siteContent := normalizeContent(injectSourceLink(normalized, src))
		if err := os.WriteFile(dst, []byte(siteContent), 0644); err != nil {
			cmd.PrintErrf("  ! failed to write %s: %v\n", entry.Filename, err)
			continue
		}
	}

	// Generate index.md
	indexContent := generateIndex(entries)
	indexPath := filepath.Join(docsDir, "index.md")
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to write index.md: %w", err)
	}

	// Generate topic pages
	var topicEntries []journalEntry
	for _, e := range entries {
		if e.IsSuggestion || isMultiPartContinuation(e.Filename) || len(e.Topics) == 0 {
			continue
		}
		topicEntries = append(topicEntries, e)
	}

	topics := buildTopicIndex(topicEntries)

	if len(topics) > 0 {
		topicsDir := filepath.Join(docsDir, "topics")
		if err := os.MkdirAll(topicsDir, 0755); err != nil {
			return fmt.Errorf("failed to create topics directory: %w", err)
		}

		// Write topics index
		topicsIndexContent := generateTopicsIndex(topics)
		if err := os.WriteFile(filepath.Join(topicsDir, "index.md"), []byte(topicsIndexContent), 0644); err != nil {
			return fmt.Errorf("failed to write topics/index.md: %w", err)
		}

		// Write individual topic pages for popular topics
		for _, t := range topics {
			if !t.Popular {
				continue
			}
			pageContent := generateTopicPage(t)
			if err := os.WriteFile(filepath.Join(topicsDir, t.Name+".md"), []byte(pageContent), 0644); err != nil {
				cmd.PrintErrf("  ! failed to write topics/%s.md: %v\n", t.Name, err)
			}
		}
	}

	// Generate key files pages
	var keyFileEntries []journalEntry
	for _, e := range entries {
		if e.IsSuggestion || isMultiPartContinuation(e.Filename) || len(e.KeyFiles) == 0 {
			continue
		}
		keyFileEntries = append(keyFileEntries, e)
	}

	keyFiles := buildKeyFileIndex(keyFileEntries)

	if len(keyFiles) > 0 {
		filesDir := filepath.Join(docsDir, "files")
		if err := os.MkdirAll(filesDir, 0755); err != nil {
			return fmt.Errorf("failed to create files directory: %w", err)
		}

		filesIndexContent := generateKeyFilesIndex(keyFiles)
		if err := os.WriteFile(filepath.Join(filesDir, "index.md"), []byte(filesIndexContent), 0644); err != nil {
			return fmt.Errorf("failed to write files/index.md: %w", err)
		}

		for _, kf := range keyFiles {
			if !kf.Popular {
				continue
			}
			pageContent := generateKeyFilePage(kf)
			slug := keyFileSlug(kf.Path)
			if err := os.WriteFile(filepath.Join(filesDir, slug+".md"), []byte(pageContent), 0644); err != nil {
				cmd.PrintErrf("  ! failed to write files/%s.md: %v\n", slug, err)
			}
		}
	}

	// Generate session type pages
	var typeEntries []journalEntry
	for _, e := range entries {
		if e.IsSuggestion || isMultiPartContinuation(e.Filename) || e.Type == "" {
			continue
		}
		typeEntries = append(typeEntries, e)
	}

	sessionTypes := buildTypeIndex(typeEntries)

	if len(sessionTypes) > 0 {
		typesDir := filepath.Join(docsDir, "types")
		if err := os.MkdirAll(typesDir, 0755); err != nil {
			return fmt.Errorf("failed to create types directory: %w", err)
		}

		typesIndexContent := generateTypesIndex(sessionTypes)
		if err := os.WriteFile(filepath.Join(typesDir, "index.md"), []byte(typesIndexContent), 0644); err != nil {
			return fmt.Errorf("failed to write types/index.md: %w", err)
		}

		for _, st := range sessionTypes {
			pageContent := generateTypePage(st)
			if err := os.WriteFile(filepath.Join(typesDir, st.Name+".md"), []byte(pageContent), 0644); err != nil {
				cmd.PrintErrf("  ! failed to write types/%s.md: %v\n", st.Name, err)
			}
		}
	}

	// Generate zensical.toml
	tomlContent := generateZensicalToml(entries, topics, keyFiles, sessionTypes)
	tomlPath := filepath.Join(output, "zensical.toml")
	if err := os.WriteFile(tomlPath, []byte(tomlContent), 0644); err != nil {
		return fmt.Errorf("failed to write zensical.toml: %w", err)
	}

	cmd.Printf("%s Generated site with %d entries in %s\n", green("✓"), len(entries), output)

	// Build or serve if requested
	if serve {
		cmd.Println("\nStarting local server...")
		return runZensical(output, "serve")
	} else if build {
		cmd.Println("\nBuilding site...")
		return runZensical(output, "build")
	}

	cmd.Println("\nNext steps:")
	cmd.Printf("  cd %s && zensical serve\n", output)
	cmd.Printf("  or \n")
	cmd.Printf("  ctx journal site --serve\n")

	return nil
}

// scanJournalEntries reads all journal markdown files and extracts metadata.
//
// Parameters:
//   - journalDir: Path to .context/journal/
//
// Returns:
//   - []journalEntry: Parsed entries sorted by date (newest first)
//   - error: Non-nil if directory scanning fails
func scanJournalEntries(journalDir string) ([]journalEntry, error) {
	files, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, err
	}

	var entries []journalEntry
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		path := filepath.Join(journalDir, f.Name())
		entry := parseJournalEntry(path, f.Name())
		entries = append(entries, entry)
	}

	// Sort by datetime (newest first) - combine Date and Time
	sort.Slice(entries, func(i, j int) bool {
		// Compare Date+Time strings (YYYY-MM-DD + HH:MM:SS)
		di := entries[i].Date + " " + entries[i].Time
		dj := entries[j].Date + " " + entries[j].Time
		return di > dj
	})

	return entries, nil
}

// parseJournalEntry extracts metadata from a journal file.
//
// Parameters:
//   - path: Full path to the journal file
//   - filename: Filename (e.g., "2026-01-21-async-roaming-allen-af7cba21.md")
//
// Returns:
//   - journalEntry: Parsed entry with title, date, project extracted
func parseJournalEntry(path, filename string) journalEntry {
	entry := journalEntry{
		Filename: filename,
		Path:     path,
	}

	// Extract date from filename (YYYY-MM-DD-slug-id.md)
	if len(filename) >= 10 {
		entry.Date = filename[:10]
	}

	// Read file to extract metadata
	content, err := os.ReadFile(path)
	if err != nil {
		entry.Title = strings.TrimSuffix(filename, ".md")
		return entry
	}

	// File size
	entry.Size = int64(len(content))

	contentStr := string(content)

	// Parse YAML frontmatter if present
	if strings.HasPrefix(contentStr, "---\n") {
		if end := strings.Index(contentStr[4:], "\n---\n"); end >= 0 {
			fmRaw := contentStr[4 : 4+end]
			var fm journalFrontmatter
			if yaml.Unmarshal([]byte(fmRaw), &fm) == nil {
				if fm.Title != "" {
					entry.Title = fm.Title
				}
				entry.Topics = fm.Topics
				entry.Type = fm.Type
				entry.Outcome = fm.Outcome
				entry.KeyFiles = fm.KeyFiles
			}
		}
	}

	// Check for suggestion mode sessions
	if strings.Contains(contentStr, "[SUGGESTION MODE:") ||
		strings.Contains(contentStr, "SUGGESTION MODE:") {
		entry.IsSuggestion = true
	}

	// Line-by-line parsing as fallback for fields not in frontmatter
	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Title from first H1 (only if frontmatter didn't set it)
		if strings.HasPrefix(line, "# ") && entry.Title == "" {
			entry.Title = strings.TrimPrefix(line, "# ")
		}

		// Time from metadata
		if strings.HasPrefix(line, "**Time**:") {
			entry.Time = strings.TrimSpace(strings.TrimPrefix(line, "**Time**:"))
		}

		// Project from metadata
		if strings.HasPrefix(line, "**Project**:") {
			entry.Project = strings.TrimSpace(strings.TrimPrefix(line, "**Project**:"))
		}

		// Stop after we have all three
		if entry.Title != "" && entry.Time != "" && entry.Project != "" {
			break
		}
	}

	if entry.Title == "" {
		entry.Title = strings.TrimSuffix(filename, ".md")
	}

	return entry
}

// generateSiteReadme creates a README for the journal-site directory.
func generateSiteReadme(journalDir string) string {
	nl := config.NewlineLF
	return "# journal-site (generated)" + nl + nl +
		"This directory is generated by `ctx journal site` and is read-only." + nl +
		"Do not edit files here — changes will be overwritten on the next run." + nl + nl +
		"## To update" + nl + nl +
		"1. Edit source entries in `" + journalDir + "/`" + nl +
		"2. Regenerate:" + nl + nl +
		"```" + nl +
		"ctx journal site          # generate" + nl +
		"ctx journal site --serve  # generate and preview" + nl +
		"```" + nl
}

// generateIndex creates the index.md content for the journal site.
//
// Parameters:
//   - entries: All journal entries to include
//
// Returns:
//   - string: Markdown content for index.md
func generateIndex(entries []journalEntry) string {
	var sb strings.Builder
	nl := config.NewlineLF

	// Separate regular sessions from suggestions and multi-part continuations
	var regular, suggestions []journalEntry
	for _, e := range entries {
		if e.IsSuggestion {
			suggestions = append(suggestions, e)
		} else if isMultiPartContinuation(e.Filename) {
			// Skip part 2+ of split sessions - they're navigable from part 1
			continue
		} else {
			regular = append(regular, e)
		}
	}

	sb.WriteString("# Session Journal" + nl + nl)
	sb.WriteString("Browse your AI session history." + nl + nl)
	sb.WriteString(fmt.Sprintf("**Sessions**: %d | **Suggestions**: %d"+nl+nl, len(regular), len(suggestions)))

	// Group regular sessions by month
	months := make(map[string][]journalEntry)
	var monthOrder []string

	for _, e := range regular {
		if len(e.Date) >= 7 {
			month := e.Date[:7] // YYYY-MM
			if _, exists := months[month]; !exists {
				monthOrder = append(monthOrder, month)
			}
			months[month] = append(months[month], e)
		}
	}

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf("## %s"+nl+nl, month))

		for _, e := range months[month] {
			sb.WriteString(formatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	// Suggestions section (collapsed by default via details tag)
	if len(suggestions) > 0 {
		sb.WriteString("---" + nl + nl)
		sb.WriteString("## Suggestions" + nl + nl)
		sb.WriteString("*Auto-generated suggestion prompts from Claude Code.*" + nl + nl)

		for _, e := range suggestions {
			sb.WriteString(formatIndexEntry(e, nl))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// formatIndexEntry formats a single entry for the index.
//
// Format: - HH:MM [title](link.md) (project) [size]
func formatIndexEntry(e journalEntry, nl string) string {
	link := strings.TrimSuffix(e.Filename, ".md")

	timeStr := ""
	if e.Time != "" && len(e.Time) >= 5 {
		timeStr = e.Time[:5] + " "
	}

	project := ""
	if e.Project != "" {
		project = fmt.Sprintf(" (%s)", e.Project)
	}

	size := formatSize(e.Size)

	return fmt.Sprintf("- %s[%s](%s.md)%s `%s`"+nl, timeStr, e.Title, link, project, size)
}

// injectSourceLink inserts an "Edit source" link into a journal entry's content.
// The link is placed after YAML frontmatter if present, otherwise at the top.
func injectSourceLink(content, sourcePath string) string {
	nl := config.NewlineLF
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		absPath = sourcePath
	}
	relPath := ".context/journal/" + filepath.Base(absPath)
	link := fmt.Sprintf(`*[View source](file://%s) · <code>%s</code>*`+
		` <button onclick="navigator.clipboard.writeText('%s')" title="Copy path"`+
		` style="cursor:pointer;border:none;background:none;font-size:0.8em;vertical-align:middle">`+
		`&#x2398;</button>`+nl+nl, absPath, relPath, relPath)

	if strings.HasPrefix(content, "---\n") {
		if end := strings.Index(content[4:], "\n---\n"); end >= 0 {
			insertAt := 4 + end + 5 // after closing "---\n"
			return content[:insertAt] + nl + link + content[insertAt:]
		}
	}

	return link + content
}

// stripFences removes all code fence markers from content, leaving the inner
// text as-is. This eliminates fence nesting conflicts entirely. Files with
// <!-- fences-verified: YYYY-MM-DD --> are skipped (fences already correct).
//
// The result is plain text with structural markers preserved (turn headers,
// tool calls, section breaks). Serves as readable baseline without AI
// reconstruction, or as input for the ctx-journal-normalize skill.
func stripFences(content string) string {
	// Skip files whose fences have been verified by the AI skill
	if fencesVerifiedPattern.MatchString(content) {
		return content
	}

	lines := strings.Split(content, "\n")
	var out []string
	inFrontmatter := false

	for i, line := range lines {
		// Preserve frontmatter
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == "---" {
				inFrontmatter = false
			}
			continue
		}

		// Remove fence markers
		if fenceLinePattern.MatchString(line) {
			continue
		}

		out = append(out, line)
	}

	return strings.Join(out, "\n")
}

// normalizeContent sanitizes journal markdown for static site rendering:
//   - Strips code fence markers (eliminates nesting conflicts)
//   - Strips bold markers from tool-use lines (🔧 **Glob: *.md** -> 🔧 Glob: *.md)
//   - Escapes glob-like * characters outside code blocks
//
// Heavy formatting (metadata tables, proper fence reconstruction) is left to
// the ctx-journal-normalize skill which uses AI for context-aware cleanup.
func normalizeContent(content string) string {
	// Strip fences first — eliminates all nesting conflicts
	content = stripFences(content)

	lines := strings.Split(content, "\n")
	var out []string
	inFrontmatter := false

	for i, line := range lines {
		// Skip frontmatter
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == "---" {
				inFrontmatter = false
			}
			continue
		}

		// Strip bold from tool-use lines
		line = toolBoldPattern.ReplaceAllString(line, `🔧 $1`)

		// Escape glob stars
		if !strings.HasPrefix(line, "    ") {
			line = globStarPattern.ReplaceAllString(line, `\*$1`)
		}

		out = append(out, line)
	}

	return strings.Join(out, "\n")
}

// cleanToolOutputJSON extracts plain text from Tool Output turns whose body is
// raw JSON from the Claude API (e.g. [{"type":"text","text":"..."}]).
// The JSON text field's \n escapes become real newlines.
func cleanToolOutputJSON(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	i := 0

	for i < len(lines) {
		matches := turnHeaderPattern.FindStringSubmatch(strings.TrimSpace(lines[i]))
		if matches == nil || matches[2] != "Tool Output" {
			out = append(out, lines[i])
			i++
			continue
		}

		// Tool Output header
		out = append(out, lines[i])
		i++

		// Collect body until next header
		bodyStart := i
		for i < len(lines) {
			if turnHeaderPattern.MatchString(strings.TrimSpace(lines[i])) {
				break
			}
			i++
		}
		bodyLines := lines[bodyStart:i]

		// Strip code fences wrapping the body, then rejoin and try JSON parse
		var nonEmpty []string
		for _, l := range bodyLines {
			t := strings.TrimSpace(l)
			if t == "" || strings.HasPrefix(t, "```") {
				continue
			}
			nonEmpty = append(nonEmpty, t)
		}
		body := strings.Join(nonEmpty, " ")

		if strings.HasPrefix(body, "[{") {
			var items []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}
			if json.Unmarshal([]byte(body), &items) == nil && len(items) > 0 {
				out = append(out, "")
				for _, item := range items {
					out = append(out, item.Text)
				}
				out = append(out, "")
				continue
			}
		}

		// Not JSON or parse failed — keep original
		out = append(out, bodyLines...)
	}

	return strings.Join(out, "\n")
}

// consolidateToolRuns collapses consecutive turns with identical body content
// into a single turn with a count. Handles both tool-call turns
// ("🔧 **TaskCreate**") and tool-output turns ("The file ... has been updated").
func consolidateToolRuns(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	i := 0

	for i < len(lines) {
		// Check if this line is a turn header
		if !turnHeaderPattern.MatchString(strings.TrimSpace(lines[i])) {
			out = append(out, lines[i])
			i++
			continue
		}

		// Extract this turn: header + body (until next header or EOF)
		header := lines[i]
		bodyStart := i + 1
		// Skip blank line after header
		if bodyStart < len(lines) && strings.TrimSpace(lines[bodyStart]) == "" {
			bodyStart++
		}
		// Collect body until next turn header
		bodyEnd := bodyStart
		for bodyEnd < len(lines) {
			if turnHeaderPattern.MatchString(strings.TrimSpace(lines[bodyEnd])) {
				break
			}
			bodyEnd++
		}
		// Trim trailing blank lines for comparison
		body := strings.TrimSpace(strings.Join(lines[bodyStart:bodyEnd], "\n"))

		// Count consecutive turns with identical body
		count := 1
		j := bodyEnd
		for j < len(lines) {
			if !turnHeaderPattern.MatchString(strings.TrimSpace(lines[j])) {
				break
			}
			nextBodyStart := j + 1
			if nextBodyStart < len(lines) && strings.TrimSpace(lines[nextBodyStart]) == "" {
				nextBodyStart++
			}
			nextBodyEnd := nextBodyStart
			for nextBodyEnd < len(lines) {
				if turnHeaderPattern.MatchString(strings.TrimSpace(lines[nextBodyEnd])) {
					break
				}
				nextBodyEnd++
			}
			nextBody := strings.TrimSpace(strings.Join(lines[nextBodyStart:nextBodyEnd], "\n"))

			if nextBody != body {
				break
			}
			count++
			j = nextBodyEnd
		}

		if count > 1 {
			out = append(out, header, "", body, "", fmt.Sprintf("(\u00d7%d)", count), "")
		} else {
			// Keep original lines (preserves blank lines as-is)
			for k := i; k < bodyEnd; k++ {
				out = append(out, lines[k])
			}
		}
		i = j
	}

	return strings.Join(out, "\n")
}

// mergeConsecutiveTurns merges back-to-back turns from the same role into a
// single turn. Keeps the first header and concatenates all bodies. This reduces
// noise from sequences like 4 consecutive "Assistant" turns each with a single
// tool call.
func mergeConsecutiveTurns(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	i := 0

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		matches := turnHeaderPattern.FindStringSubmatch(trimmed)
		if matches == nil {
			out = append(out, lines[i])
			i++
			continue
		}

		role := matches[2]
		header := lines[i]

		// Collect body from this and all consecutive same-role turns,
		// explicitly skipping intermediate headers.
		var body []string
		j := i + 1
		for {
			// Collect body lines until next header or EOF
			for j < len(lines) {
				if turnHeaderPattern.MatchString(strings.TrimSpace(lines[j])) {
					break
				}
				body = append(body, lines[j])
				j++
			}
			// Check if next turn has same role
			if j >= len(lines) {
				break
			}
			nextMatches := turnHeaderPattern.FindStringSubmatch(strings.TrimSpace(lines[j]))
			if nextMatches == nil || nextMatches[2] != role {
				break
			}
			// Same role — skip the header, continue collecting body
			j++
		}

		out = append(out, header)
		out = append(out, body...)
		i = j
	}

	return strings.Join(out, "\n")
}

// stripSystemReminders removes system reminder blocks from journal content.
// Handles two formats:
//   - XML-style: <system-reminder>...</system-reminder>
//   - Bold-style: **System Reminder**: ... (paragraph until blank line)
//
// The authoritative JSONL transcripts retain them; the exported markdown
// doesn't need them.
func stripSystemReminders(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	inTagReminder := false
	inBoldReminder := false

	for _, line := range lines {
		// XML-style: <system-reminder>...</system-reminder>
		if strings.TrimSpace(line) == "<system-reminder>" {
			inTagReminder = true
			continue
		}
		if inTagReminder {
			if strings.TrimSpace(line) == "</system-reminder>" {
				inTagReminder = false
			}
			continue
		}

		// Bold-style: **System Reminder**: ... (runs until blank line)
		if strings.HasPrefix(strings.TrimSpace(line), "**System Reminder**:") {
			inBoldReminder = true
			continue
		}
		if inBoldReminder {
			if strings.TrimSpace(line) == "" {
				inBoldReminder = false
			}
			continue
		}

		out = append(out, line)
	}

	return strings.Join(out, "\n")
}

// softWrapContent wraps long lines in source journal files to ~80 characters.
// Skips only frontmatter and table rows. Wraps everything else including
// content inside code fences — journal files are reference material, not
// executable code.
func softWrapContent(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	inFrontmatter := false

	for i, line := range lines {
		// Skip frontmatter
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == "---" {
				inFrontmatter = false
			}
			continue
		}

		// Wrap long lines (skip tables)
		if len(line) > 80 && !strings.HasPrefix(strings.TrimSpace(line), "|") {
			out = append(out, softWrap(line, 80)...)
		} else {
			out = append(out, line)
		}
	}

	return strings.Join(out, "\n")
}

// softWrap breaks a long line at word boundaries, preserving leading indent.
func softWrap(line string, width int) []string {
	trimmed := strings.TrimLeft(line, " \t")
	indent := line[:len(line)-len(trimmed)]

	words := strings.Fields(trimmed)
	if len(words) == 0 {
		return []string{line}
	}

	var result []string
	current := indent + words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) > width && len(current) > len(indent) {
			result = append(result, current)
			current = indent + word
		} else {
			current += " " + word
		}
	}
	result = append(result, current)
	return result
}

// formatSize formats a file size in human-readable form.
func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	kb := float64(bytes) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%.1fKB", kb)
	}
	mb := kb / 1024
	return fmt.Sprintf("%.1fMB", mb)
}

// isMultiPartContinuation returns true if filename is a continuation part (p2, p3, etc.)
func isMultiPartContinuation(filename string) bool {
	return multiPartPattern.MatchString(filename)
}

// topicData holds aggregated data for a single topic.
type topicData struct {
	Name    string
	Entries []journalEntry
	Popular bool
}

// buildTopicIndex aggregates entries by topic and returns sorted topic data.
// Topics with 2+ sessions are marked popular. Sorted by count desc, then alpha.
func buildTopicIndex(entries []journalEntry) []topicData {
	byTopic := make(map[string][]journalEntry)
	for _, e := range entries {
		for _, t := range e.Topics {
			byTopic[t] = append(byTopic[t], e)
		}
	}

	topics := make([]topicData, 0, len(byTopic))
	for name, ents := range byTopic {
		topics = append(topics, topicData{
			Name:    name,
			Entries: ents,
			Popular: len(ents) >= 2,
		})
	}

	sort.Slice(topics, func(i, j int) bool {
		if len(topics[i].Entries) != len(topics[j].Entries) {
			return len(topics[i].Entries) > len(topics[j].Entries)
		}
		return topics[i].Name < topics[j].Name
	})

	return topics
}

// countUniqueSessions counts unique sessions across all topics.
func countUniqueSessions(topics []topicData) int {
	seen := make(map[string]bool)
	for _, t := range topics {
		for _, e := range t.Entries {
			seen[e.Filename] = true
		}
	}
	return len(seen)
}

// generateTopicsIndex creates the topics/index.md page.
// Popular topics link to dedicated pages; long-tail topics list entries inline.
func generateTopicsIndex(topics []topicData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	var popular, longtail []topicData
	for _, t := range topics {
		if t.Popular {
			popular = append(popular, t)
		} else {
			longtail = append(longtail, t)
		}
	}

	sb.WriteString("# Topics" + nl + nl)
	sb.WriteString(fmt.Sprintf("**%d topics** across **%d sessions** — **%d popular**, **%d long-tail**"+nl+nl,
		len(topics), countUniqueSessions(topics), len(popular), len(longtail)))

	// Popular topics
	if len(popular) > 0 {
		sb.WriteString("## Popular Topics" + nl + nl)
		for _, t := range popular {
			sb.WriteString(fmt.Sprintf("- [%s](%s.md) (%d sessions)"+nl, t.Name, t.Name, len(t.Entries)))
		}
		sb.WriteString(nl)
	}

	// Long-tail topics
	if len(longtail) > 0 {
		sb.WriteString("## Long-tail Topics" + nl + nl)
		for _, t := range longtail {
			e := t.Entries[0]
			link := strings.TrimSuffix(e.Filename, ".md")
			sb.WriteString(fmt.Sprintf("- **%s** — [%s](../%s.md)"+nl, t.Name, e.Title, link))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// generateTopicPage creates an individual topic page with sessions grouped by month.
func generateTopicPage(topic topicData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(fmt.Sprintf("# %s"+nl+nl, topic.Name))
	sb.WriteString(fmt.Sprintf("**%d sessions** with this topic."+nl+nl, len(topic.Entries)))

	// Group by month
	months := make(map[string][]journalEntry)
	var monthOrder []string

	for _, e := range topic.Entries {
		if len(e.Date) >= 7 {
			month := e.Date[:7]
			if _, exists := months[month]; !exists {
				monthOrder = append(monthOrder, month)
			}
			months[month] = append(months[month], e)
		}
	}

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf("## %s"+nl+nl, month))
		for _, e := range months[month] {
			link := strings.TrimSuffix(e.Filename, ".md")
			timeStr := ""
			if e.Time != "" && len(e.Time) >= 5 {
				timeStr = e.Time[:5] + " "
			}
			sb.WriteString(fmt.Sprintf("- %s[%s](../%s.md)"+nl, timeStr, e.Title, link))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// keyFileData holds aggregated data for a single file path.
type keyFileData struct {
	Path    string
	Entries []journalEntry
	Popular bool
}

// keyFileSlug converts a file path to a safe slug for use as a filename.
func keyFileSlug(path string) string {
	slug := strings.ReplaceAll(path, "/", "_")
	slug = strings.ReplaceAll(slug, ".", "_")
	slug = strings.ReplaceAll(slug, "*", "x")
	return slug
}

// buildKeyFileIndex aggregates entries by key file path.
// Files with 2+ sessions are marked popular.
func buildKeyFileIndex(entries []journalEntry) []keyFileData {
	byFile := make(map[string][]journalEntry)
	for _, e := range entries {
		for _, f := range e.KeyFiles {
			byFile[f] = append(byFile[f], e)
		}
	}

	files := make([]keyFileData, 0, len(byFile))
	for path, ents := range byFile {
		files = append(files, keyFileData{
			Path:    path,
			Entries: ents,
			Popular: len(ents) >= 2,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if len(files[i].Entries) != len(files[j].Entries) {
			return len(files[i].Entries) > len(files[j].Entries)
		}
		return files[i].Path < files[j].Path
	})

	return files
}

// generateKeyFilesIndex creates the files/index.md page.
func generateKeyFilesIndex(keyFiles []keyFileData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	var popular, longtail []keyFileData
	for _, kf := range keyFiles {
		if kf.Popular {
			popular = append(popular, kf)
		} else {
			longtail = append(longtail, kf)
		}
	}

	totalSessions := 0
	seen := make(map[string]bool)
	for _, kf := range keyFiles {
		for _, e := range kf.Entries {
			if !seen[e.Filename] {
				seen[e.Filename] = true
				totalSessions++
			}
		}
	}

	sb.WriteString("# Key Files" + nl + nl)
	sb.WriteString(fmt.Sprintf("**%d files** across **%d sessions** — **%d popular**, **%d long-tail**"+nl+nl,
		len(keyFiles), totalSessions, len(popular), len(longtail)))

	if len(popular) > 0 {
		sb.WriteString("## Frequently Touched" + nl + nl)
		for _, kf := range popular {
			slug := keyFileSlug(kf.Path)
			sb.WriteString(fmt.Sprintf("- [`%s`](%s.md) (%d sessions)"+nl, kf.Path, slug, len(kf.Entries)))
		}
		sb.WriteString(nl)
	}

	if len(longtail) > 0 {
		sb.WriteString("## Single Session" + nl + nl)
		for _, kf := range longtail {
			e := kf.Entries[0]
			link := strings.TrimSuffix(e.Filename, ".md")
			sb.WriteString(fmt.Sprintf("- `%s` — [%s](../%s.md)"+nl, kf.Path, e.Title, link))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// generateKeyFilePage creates an individual key file page with sessions grouped by month.
func generateKeyFilePage(kf keyFileData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(fmt.Sprintf("# `%s`"+nl+nl, kf.Path))
	sb.WriteString(fmt.Sprintf("**%d sessions** touching this file."+nl+nl, len(kf.Entries)))

	months := make(map[string][]journalEntry)
	var monthOrder []string

	for _, e := range kf.Entries {
		if len(e.Date) >= 7 {
			month := e.Date[:7]
			if _, exists := months[month]; !exists {
				monthOrder = append(monthOrder, month)
			}
			months[month] = append(months[month], e)
		}
	}

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf("## %s"+nl+nl, month))
		for _, e := range months[month] {
			link := strings.TrimSuffix(e.Filename, ".md")
			timeStr := ""
			if e.Time != "" && len(e.Time) >= 5 {
				timeStr = e.Time[:5] + " "
			}
			sb.WriteString(fmt.Sprintf("- %s[%s](../%s.md)"+nl, timeStr, e.Title, link))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// typeData holds aggregated data for a session type.
type typeData struct {
	Name    string
	Entries []journalEntry
}

// buildTypeIndex aggregates entries by session type.
func buildTypeIndex(entries []journalEntry) []typeData {
	byType := make(map[string][]journalEntry)
	for _, e := range entries {
		byType[e.Type] = append(byType[e.Type], e)
	}

	types := make([]typeData, 0, len(byType))
	for name, ents := range byType {
		types = append(types, typeData{
			Name:    name,
			Entries: ents,
		})
	}

	sort.Slice(types, func(i, j int) bool {
		if len(types[i].Entries) != len(types[j].Entries) {
			return len(types[i].Entries) > len(types[j].Entries)
		}
		return types[i].Name < types[j].Name
	})

	return types
}

// generateTypesIndex creates the types/index.md page.
func generateTypesIndex(sessionTypes []typeData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	totalSessions := 0
	for _, st := range sessionTypes {
		totalSessions += len(st.Entries)
	}

	sb.WriteString("# Session Types" + nl + nl)
	sb.WriteString(fmt.Sprintf("**%d types** across **%d sessions**"+nl+nl, len(sessionTypes), totalSessions))

	for _, st := range sessionTypes {
		sb.WriteString(fmt.Sprintf("- [%s](%s.md) (%d sessions)"+nl, st.Name, st.Name, len(st.Entries)))
	}
	sb.WriteString(nl)

	return sb.String()
}

// generateTypePage creates an individual session type page with sessions grouped by month.
func generateTypePage(st typeData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(fmt.Sprintf("# %s"+nl+nl, st.Name))
	sb.WriteString(fmt.Sprintf("**%d sessions** of type *%s*."+nl+nl, len(st.Entries), st.Name))

	months := make(map[string][]journalEntry)
	var monthOrder []string

	for _, e := range st.Entries {
		if len(e.Date) >= 7 {
			month := e.Date[:7]
			if _, exists := months[month]; !exists {
				monthOrder = append(monthOrder, month)
			}
			months[month] = append(months[month], e)
		}
	}

	for _, month := range monthOrder {
		sb.WriteString(fmt.Sprintf("## %s"+nl+nl, month))
		for _, e := range months[month] {
			link := strings.TrimSuffix(e.Filename, ".md")
			timeStr := ""
			if e.Time != "" && len(e.Time) >= 5 {
				timeStr = e.Time[:5] + " "
			}
			sb.WriteString(fmt.Sprintf("- %s[%s](../%s.md)"+nl, timeStr, e.Title, link))
		}
		sb.WriteString(nl)
	}

	return sb.String()
}

// generateZensicalToml creates the zensical.toml configuration.
func generateZensicalToml(entries []journalEntry, topics []topicData, keyFiles []keyFileData, sessionTypes []typeData) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(`[project]
site_name = "ctx: Session Journal"
site_description = "AI session history and notes"
site_author = "Jose Alekhinne <alekhinejose@gmail.com>"
site_url = "https://ctx.ist/"
repo_url = "https://github.com/ActiveMemory/ctx"
repo_name = "ActiveMemory/ctx"
copyright = """
Copyright &copy; 2026&ndash;present <a href="https://github.com/ActiveMemory/ctx/">Context contributors</a>.<br>
Context's code is distributed under
<a href="https://github.com/ActiveMemory/ctx/blob/main/LICENSE"><strong>Apache (v2.0)</strong></a>.<br>
"""

` + nl)

	// Build navigation
	sb.WriteString("nav = [" + nl)
	sb.WriteString(`  { "Home" = "index.md" },` + nl)
	if len(topics) > 0 {
		sb.WriteString(`  { "Topics" = "topics/index.md" },` + nl)
	}
	if len(keyFiles) > 0 {
		sb.WriteString(`  { "Files" = "files/index.md" },` + nl)
	}
	if len(sessionTypes) > 0 {
		sb.WriteString(`  { "Types" = "types/index.md" },` + nl)
	}

	// Filter out suggestion sessions and multi-part continuations from navigation
	var regular []journalEntry
	for _, e := range entries {
		if e.IsSuggestion {
			continue
		}
		// Skip part 2+ of split sessions (e.g., "...-p2.md", "...-p3.md")
		if isMultiPartContinuation(e.Filename) {
			continue
		}
		regular = append(regular, e)
	}

	// Group recent entries (last 20, excluding suggestions)
	recent := regular
	if len(recent) > 20 {
		recent = recent[:20]
	}

	sb.WriteString(`  { "Recent Sessions" = [` + nl)
	for _, e := range recent {
		title := e.Title
		if len(title) > 40 {
			title = title[:40] + "..."
		}
		// Escape quotes in title
		title = strings.ReplaceAll(title, `"`, `\"`)
		sb.WriteString(fmt.Sprintf(`    { "%s" = "%s" },`+nl, title, e.Filename))
	}
	sb.WriteString("  ]}" + nl)
	sb.WriteString("]" + nl + nl)

	sb.WriteString(`

[project.theme]
language = "en"
features = [
    "content.code.copy",
    "navigation.instant",
    "navigation.top",
    "search.highlight",
]

[[project.theme.palette]]
scheme = "default"
toggle.icon = "lucide/sun"
toggle.name = "Switch to dark mode"

[[project.theme.palette]]
scheme = "slate"
toggle.icon = "lucide/moon"
toggle.name = "Switch to light mode"

[[project.theme.palette]]
scheme = "slate"
toggle.icon = "lucide/moon"
toggle.name = "Switch to light mode"

[[project.extra.social]]
icon = "fontawesome/brands/github"
link = "https://github.com/ActiveMemory/ctx"

[[project.extra.social]]
icon = "fontawesome/brands/discord"
link = "https://discord.gg/kampus"

[project.extra]
generator = false
`)

	return sb.String()
}

// runZensical executes zensical build or serve in the output directory.
//
// Parameters:
//   - dir: Directory containing the generated site
//   - command: "build" or "serve"
//
// Returns:
//   - error: Non-nil if zensical is not found or fails
func runZensical(dir, command string) error {
	// Check if zensical is available
	_, err := exec.LookPath("zensical")
	if err != nil {
		return fmt.Errorf("zensical not found. Install with: pip install zensical")
	}

	cmd := exec.Command("zensical", command)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
