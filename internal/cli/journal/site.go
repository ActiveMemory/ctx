//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

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

	defaultOutput := filepath.Join(rc.GetContextDir(), "journal-site")
	cmd.Flags().StringVarP(&output, "output", "o", defaultOutput, "Output directory for site")
	cmd.Flags().BoolVar(&build, "build", false, "Run zensical build after generating")
	cmd.Flags().BoolVar(&serve, "serve", false, "Run zensical serve after generating")

	return cmd
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
	journalDir := filepath.Join(rc.GetContextDir(), "journal")

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

	// Copy journal files to docs/
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(docsDir, entry.Filename)

		content, err := os.ReadFile(src)
		if err != nil {
			cmd.PrintErrf("  ! failed to read %s: %v\n", entry.Filename, err)
			continue
		}

		if err := os.WriteFile(dst, content, 0644); err != nil {
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

	// Generate zensical.toml
	tomlContent := generateZensicalToml(entries)
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

	// Check for suggestion mode sessions
	contentStr := string(content)
	if strings.Contains(contentStr, "[SUGGESTION MODE:") ||
		strings.Contains(contentStr, "SUGGESTION MODE:") {
		entry.IsSuggestion = true
	}

	lines := strings.Split(contentStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Title from first H1
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

	// Separate regular sessions from suggestions
	var regular, suggestions []journalEntry
	for _, e := range entries {
		if e.IsSuggestion {
			suggestions = append(suggestions, e)
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

// generateZensicalToml creates the zensical.toml configuration.
//
// Parameters:
//   - entries: All journal entries for navigation
//
// Returns:
//   - string: TOML content for zensical.toml
func generateZensicalToml(entries []journalEntry) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(`[project]
site_name = "Session Journal"
site_description = "AI session history and notes"
` + nl)

	// Build navigation
	sb.WriteString("nav = [" + nl)
	sb.WriteString(`  { "Home" = "index.md" },` + nl)

	// Group recent entries (last 20)
	recent := entries
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

	sb.WriteString(`[project.theme]
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
