//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package feed

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/cli/site/core"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/rss"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errSite "github.com/ActiveMemory/ctx/internal/err/site"
)

// Run orchestrates scanning and generation.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - blogDir: Path to the blog posts directory
//   - outPath: Output path for the generated feed
//   - baseURL: Base URL for entry links
//
// Returns:
//   - error: Non-nil if scanning or generation fails
func Run(
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
//
// Parameters:
//   - blogDir: Path to the blog posts directory
//
// Returns:
//   - []blogPost: Parsed blog posts sorted by date descending
//   - feedReport: Report of skipped and warned entries
//   - error: Non-nil if directory access fails
func scanBlogPosts(
	blogDir string,
) ([]blogPost, feedReport, error) {
	var report feedReport

	info, statErr := os.Stat(blogDir)
	if statErr != nil || !info.IsDir() {
		return nil, report, errFs.DirNotFound(blogDir)
	}

	entries, readErr := os.ReadDir(blogDir)
	if readErr != nil {
		return nil, report, errFs.ReadDir("blog directory", readErr)
	}

	var posts []blogPost

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !regBlogDatePattern.MatchString(name) {
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
//
// Parameters:
//   - path: Absolute path to the blog post file
//   - filename: Base filename of the blog post
//
// Returns:
//   - blogPost: Parsed blog post metadata
//   - postStatus: Whether the post was included, skipped, or warned
func parsePost(path, filename string) (blogPost, postStatus) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return blogPost{
			filename: filename,
			summary:  filename + " \u2014 cannot read file",
		}, postSkipped
	}

	content := string(data)
	nl := token.NewlineLF
	sep := token.Separator

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
//
// Parameters:
//   - body: Post body text after frontmatter
//
// Returns:
//   - string: First paragraph text, or empty if none found
func extractSummary(body string) string {
	lines := strings.Split(body, token.NewlineLF)
	foundHeading := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, token.PrefixHeading) {
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
			strings.HasPrefix(trimmed, token.PrefixHeading) {
			continue
		}

		return trimmed
	}

	return ""
}

// generateAtom builds the Atom XML and writes it to outPath.
//
// Parameters:
//   - posts: Blog posts to include in the feed
//   - outPath: Output file path for the generated XML
//   - baseURL: Base URL for entry links
//
// Returns:
//   - error: Non-nil if marshalling or writing fails
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

	feed := core.AtomFeed{
		NS:    rss.FeedAtomNS,
		Title: rss.FeedTitle,
		Links: []core.AtomLink{
			{Href: blogURL},
			{Href: feedURL, Rel: "self"},
		},
		ID:      feedURL,
		Updated: updated,
	}

	for _, p := range posts {
		slug := strings.TrimSuffix(p.filename, file.ExtMarkdown)
		entryURL := blogURL + slug + "/"

		entry := core.AtomEntry{
			Title:   p.title,
			Links:   []core.AtomLink{{Href: entryURL}},
			ID:      entryURL,
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
			author = rss.FeedDefaultAuthor
		}
		entry.Author = &core.AtomAuthor{Name: author}

		for _, topic := range p.topics {
			entry.Categories = append(
				entry.Categories,
				core.AtomCategory{Term: topic},
			)
		}

		feed.Entries = append(feed.Entries, entry)
	}

	outDir := filepath.Dir(outPath)
	if mkErr := os.MkdirAll(outDir, 0o755); mkErr != nil {
		return errFs.Mkdir("output directory", mkErr)
	}

	xmlData, marshalErr := xml.MarshalIndent(feed, "", "  ")
	if marshalErr != nil {
		return errSite.MarshalFeed(marshalErr)
	}

	output := []byte(rss.FeedXMLHeader)
	output = append(output, xmlData...)
	output = append(output, '\n')

	if writeErr := os.WriteFile(
		outPath, output, 0o644,
	); writeErr != nil {
		return errFs.FileWrite(outPath, writeErr)
	}

	return nil
}

// printReport outputs the generation summary.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - outPath: Path of the generated feed file
//   - report: Feed generation report with counts and messages
func printReport(
	cmd *cobra.Command, outPath string, report feedReport,
) {
	cmd.Println(fmt.Sprintf(
		"\nGenerated %s (%d entries)", outPath, report.included))

	if len(report.skipped) > 0 {
		cmd.Println("\nSkipped:")
		for _, msg := range report.skipped {
			cmd.Println("  " + msg)
		}
	}

	if len(report.warnings) > 0 {
		cmd.Println("\nWarnings:")
		for _, msg := range report.warnings {
			cmd.Println("  " + msg)
		}
	}
}
