//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/site/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
)

var regBlogDatePattern = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2}-.+\.md$`,
)

// ScanBlogPosts reads blog posts from blogDir, parses metadata, and
// returns them sorted by date descending.
//
// Parameters:
//   - blogDir: Path to the blog posts directory
//
// Returns:
//   - []BlogPost: Parsed blog posts sorted by date descending
//   - FeedReport: Report of skipped and warned entries
//   - error: Non-nil if directory access fails
func ScanBlogPosts(blogDir string) ([]core.BlogPost, core.FeedReport, error) {
	var report core.FeedReport

	info, statErr := os.Stat(blogDir)
	if statErr != nil || !info.IsDir() {
		return nil, report, errFs.DirNotFound(blogDir)
	}

	entries, readErr := os.ReadDir(blogDir)
	if readErr != nil {
		return nil, report, errFs.ReadDir(blogDir, readErr)
	}

	var posts []core.BlogPost

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !regBlogDatePattern.MatchString(name) {
			continue
		}

		post, status := ParsePost(filepath.Join(blogDir, name), name)

		switch status {
		case core.PostIncluded:
			posts = append(posts, post)
		case core.PostSkipped:
			report.Skipped = append(report.Skipped, post.Summary)
		case core.PostWarn:
			posts = append(posts, post)
			report.Warnings = append(report.Warnings, post.Summary)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Date != posts[j].Date {
			return posts[i].Date > posts[j].Date
		}
		return posts[i].Filename > posts[j].Filename
	})

	return posts, report, nil
}

// ParsePost reads a single blog post file and extracts metadata.
//
// Parameters:
//   - path: Absolute path to the blog post file
//   - filename: Base filename of the blog post
//
// Returns:
//   - BlogPost: Parsed blog post metadata
//   - PostStatus: Whether the post was included, skipped, or warned
func ParsePost(path, filename string) (core.BlogPost, core.PostStatus) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return core.BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(
				desc.Text(text.DescKeySiteSkipCannotRead), filename),
		}, core.PostSkipped
	}

	content := string(data)
	nl := token.NewlineLF
	sep := token.Separator

	if !strings.HasPrefix(content, sep+nl) {
		return core.BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(
				desc.Text(text.DescKeySiteSkipNoFrontmatter), filename),
		}, core.PostSkipped
	}

	fmStart := len(sep + nl)
	endIdx := strings.Index(content[fmStart:], nl+sep+nl)
	if endIdx < 0 {
		return core.BlogPost{
			Filename: filename,
			Summary:  fmt.Sprintf(desc.Text(text.DescKeySiteSkipMalformed), filename),
		}, core.PostSkipped
	}

	fmRaw := content[fmStart : fmStart+endIdx]
	body := content[fmStart+endIdx+len(nl+sep+nl):]

	var fm core.BlogFrontmatter
	if unmarshalErr := yaml.Unmarshal([]byte(fmRaw), &fm); unmarshalErr != nil {
		return core.BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteSkipParseError),
				filename, unmarshalErr),
		}, core.PostSkipped
	}

	if fm.ReviewedAndFinalized == nil || !*fm.ReviewedAndFinalized {
		return core.BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(
				desc.Text(text.DescKeySiteSkipNotFinalized), filename),
		}, core.PostSkipped
	}
	if fm.Title == "" {
		return core.BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteSkipMissingTitle),
				filename),
		}, core.PostSkipped
	}
	if fm.Date == "" {
		return core.BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteSkipMissingDate),
				filename),
		}, core.PostSkipped
	}

	summary := ExtractSummary(body)

	if summary == "" {
		return core.BlogPost{
			Filename: filename, Title: fm.Title, Date: fm.Date,
			Author: fm.Author, Topics: fm.Topics,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteWarnNoSummary), filename),
		}, core.PostWarn
	}

	return core.BlogPost{
		Filename: filename, Title: fm.Title, Date: fm.Date,
		Author: fm.Author, Topics: fm.Topics, Summary: summary,
	}, core.PostIncluded
}

// ExtractSummary finds the first non-empty paragraph after a heading line.
//
// Parameters:
//   - body: Post body text after frontmatter
//
// Returns:
//   - string: First paragraph text, or empty if none found
func ExtractSummary(body string) string {
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

		if trimmed == "" ||
			strings.HasPrefix(trimmed, token.PrefixBang) ||
			strings.HasPrefix(trimmed, token.PrefixStar) ||
			strings.HasPrefix(trimmed, token.PrefixHeading) {
			continue
		}

		return trimmed
	}

	return ""
}
