//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config"
)

const (
	atomNS        = "http://www.w3.org/2005/Atom"
	feedTitle     = "ctx blog"
	defaultAuthor = "Context contributors"
	xmlHeader     = `<?xml version="1.0" encoding="utf-8"?>` + "\n"
)

// blogDatePattern matches filenames like 2026-02-25-slug.md.
var blogDatePattern = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2}-.+\.md$`,
)

// blogPost holds parsed metadata from a single blog post.
type blogPost struct {
	filename string
	title    string
	date     string
	author   string
	topics   []string
	summary  string
}

// feedReport tracks what happened during feed generation.
type feedReport struct {
	included int
	skipped  []string // "filename — reason"
	warnings []string // "filename — reason"
}

// blogFrontmatter maps the YAML fields we care about.
type blogFrontmatter struct {
	Title                string   `yaml:"title"`
	Date                 string   `yaml:"date"`
	Author               string   `yaml:"author"`
	Topics               []string `yaml:"topics"`
	ReviewedAndFinalized *bool    `yaml:"reviewed_and_finalized"`
}

// feedCmd returns the "ctx site feed" subcommand.
func feedCmd() *cobra.Command {
	var (
		out     string
		baseURL string
	)

	cmd := &cobra.Command{
		Use:   "feed",
		Short: "Generate an Atom 1.0 feed from blog posts",
		Long: `Generate an Atom 1.0 feed from finalized blog posts in docs/blog/.

Parses YAML frontmatter for title, date, author, and topics. Extracts
a summary from the first paragraph after the heading. Only posts with
reviewed_and_finalized: true are included.

Examples:
  ctx site feed
  ctx site feed --out /tmp/feed.xml
  ctx site feed --base-url https://example.com`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFeed(cmd, "docs/blog", out, baseURL)
		},
	}

	cmd.Flags().StringVarP(
		&out, "out", "o", filepath.Join("site", "feed.xml"),
		"Output path for the generated feed",
	)
	cmd.Flags().StringVar(
		&baseURL, "base-url", "https://ctx.ist",
		"Base URL for entry links",
	)

	return cmd
}

// runFeed orchestrates scanning and generation.
func runFeed(
	cmd *cobra.Command, blogDir, outPath, baseURL string,
) error {
	posts, report, scanErr := scanBlogPosts(blogDir)
	if scanErr != nil {
		return scanErr
	}

	genErr := generateAtom(posts, outPath, baseURL)
	if genErr != nil {
		return genErr
	}

	report.included = len(posts)
	printReport(cmd, outPath, report)

	return nil
}

// scanBlogPosts reads blog posts from blogDir and returns parsed
// posts, a report of skipped/warned entries, and any fatal error.
func scanBlogPosts(
	blogDir string,
) ([]blogPost, feedReport, error) {
	var report feedReport

	info, statErr := os.Stat(blogDir)
	if statErr != nil || !info.IsDir() {
		return nil, report, fmt.Errorf(
			"no blog directory found at %s", blogDir,
		)
	}

	entries, readErr := os.ReadDir(blogDir)
	if readErr != nil {
		return nil, report, fmt.Errorf(
			"cannot read blog directory: %w", readErr,
		)
	}

	var posts []blogPost

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !blogDatePattern.MatchString(name) {
			continue
		}

		post, status := parsePost(
			filepath.Join(blogDir, name), name,
		)

		switch status {
		case postIncluded:
			posts = append(posts, post)
		case postSkipped:
			// reason already in report via skip helper
		case postWarn:
			posts = append(posts, post)
		}

		// Collect skip/warn messages set during parsePost.
		// parsePost returns the post and status; the caller
		// builds the report.
		if status == postSkipped {
			report.skipped = append(
				report.skipped, post.summary,
			)
		}
		if status == postWarn {
			report.warnings = append(
				report.warnings, post.summary,
			)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		if posts[i].date != posts[j].date {
			return posts[i].date > posts[j].date
		}
		return posts[i].filename > posts[j].filename
	})

	return posts, report, nil
}

type postStatus int

const (
	postIncluded postStatus = iota
	postSkipped
	postWarn
)

// parsePost reads a single blog post file and extracts metadata.
//
// Returns the post and a status indicating whether it was included,
// skipped, or included with a warning. For skipped/warned posts, the
// summary field carries the reason message.
func parsePost(path, filename string) (blogPost, postStatus) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 cannot read file",
		}, postSkipped
	}

	content := string(data)
	nl := config.NewlineLF
	sep := config.Separator

	if !strings.HasPrefix(content, sep+nl) {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 no frontmatter found",
		}, postSkipped
	}

	fmStart := len(sep + nl)
	endIdx := strings.Index(content[fmStart:], nl+sep+nl)
	if endIdx < 0 {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 malformed frontmatter",
		}, postSkipped
	}

	fmRaw := content[fmStart : fmStart+endIdx]
	body := content[fmStart+endIdx+len(nl+sep+nl):]

	var fm blogFrontmatter
	if unmarshalErr := yaml.Unmarshal(
		[]byte(fmRaw), &fm,
	); unmarshalErr != nil {
		return blogPost{
			filename: filename,
			summary: fmt.Sprintf(
				"%s \u2014 %s", filename, unmarshalErr,
			),
		}, postSkipped
	}

	// Draft gate.
	if fm.ReviewedAndFinalized == nil || !*fm.ReviewedAndFinalized {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 not finalized",
		}, postSkipped
	}

	if fm.Title == "" {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 missing title",
		}, postSkipped
	}

	if fm.Date == "" {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 missing date",
		}, postSkipped
	}

	summary := extractSummary(body)

	post := blogPost{
		filename: filename,
		title:    fm.Title,
		date:     fm.Date,
		author:   fm.Author,
		topics:   fm.Topics,
		summary:  summary,
	}

	if summary == "" {
		return blogPost{
			filename: filename,
			title:    fm.Title,
			date:     fm.Date,
			author:   fm.Author,
			topics:   fm.Topics,
			summary: filename +
				" \u2014 no summary paragraph found",
		}, postWarn
	}

	return post, postIncluded
}

// extractSummary finds the first non-empty paragraph after a
// heading line. Returns empty string if none found.
func extractSummary(body string) string {
	lines := strings.Split(body, config.NewlineLF)
	foundHeading := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, config.PrefixHeading) {
			foundHeading = true
			continue
		}

		if !foundHeading {
			continue
		}

		// Skip blank lines, images, admonitions, bylines.
		if trimmed == "" ||
			strings.HasPrefix(trimmed, "!") ||
			strings.HasPrefix(trimmed, "*") ||
			strings.HasPrefix(trimmed, config.PrefixHeading) {
			continue
		}

		return trimmed
	}

	return ""
}

// generateAtom builds the Atom XML and writes it to outPath.
func generateAtom(
	posts []blogPost, outPath, baseURL string,
) error {
	baseURL = strings.TrimRight(baseURL, "/")

	feedURL := baseURL + "/feed.xml"
	blogURL := baseURL + "/blog/"

	updated := ""
	if len(posts) > 0 {
		updated = posts[0].date + "T00:00:00Z"
	}

	feed := AtomFeed{
		NS:    atomNS,
		Title: feedTitle,
		Links: []AtomLink{
			{Href: blogURL},
			{Href: feedURL, Rel: "self"},
		},
		ID:      feedURL,
		Updated: updated,
	}

	for _, p := range posts {
		slug := strings.TrimSuffix(p.filename, ".md")
		entryURL := blogURL + slug + "/"

		entry := AtomEntry{
			Title: p.title,
			Links: []AtomLink{{Href: entryURL}},
			ID:    entryURL,
			Updated: p.date + "T00:00:00Z",
		}

		// Only set summary if it's actual content, not a
		// warning message (warn posts have the filename in
		// summary).
		if p.summary != "" &&
			!strings.Contains(p.summary, " \u2014 ") {
			entry.Summary = p.summary
		}

		author := p.author
		if author == "" {
			author = defaultAuthor
		}
		entry.Author = &AtomAuthor{Name: author}

		for _, topic := range p.topics {
			entry.Categories = append(
				entry.Categories,
				AtomCategory{Term: topic},
			)
		}

		feed.Entries = append(feed.Entries, entry)
	}

	outDir := filepath.Dir(outPath)
	if mkErr := os.MkdirAll(outDir, 0o755); mkErr != nil {
		return fmt.Errorf("cannot create output directory: %w", mkErr)
	}

	xmlData, marshalErr := xml.MarshalIndent(feed, "", "  ")
	if marshalErr != nil {
		return fmt.Errorf("cannot marshal feed: %w", marshalErr)
	}

	output := []byte(xmlHeader)
	output = append(output, xmlData...)
	output = append(output, '\n')

	if writeErr := os.WriteFile(
		outPath, output, 0o644,
	); writeErr != nil {
		return fmt.Errorf("cannot write feed: %w", writeErr)
	}

	return nil
}

// printReport outputs the generation summary.
func printReport(
	cmd *cobra.Command, outPath string, report feedReport,
) {
	cmd.Printf("\nGenerated %s (%d entries)\n",
		outPath, report.included)

	if len(report.skipped) > 0 {
		cmd.Println("\nSkipped:")
		for _, msg := range report.skipped {
			cmd.Printf("  %s\n", msg)
		}
	}

	if len(report.warnings) > 0 {
		cmd.Println("\nWarnings:")
		for _, msg := range report.warnings {
			cmd.Printf("  %s\n", msg)
		}
	}
}
