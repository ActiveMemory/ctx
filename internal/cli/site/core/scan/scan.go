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
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// BlogPosts reads blog posts from blogDir, parses metadata, and
// returns them sorted by date descending.
//
// Parameters:
//   - blogDir: Path to the blog posts directory
//
// Returns:
//   - []BlogPost: Parsed blog posts sorted by date descending
//   - FeedReport: Report of skipped and warned entries
//   - error: Non-nil if directory access fails
func BlogPosts(blogDir string) ([]BlogPost, FeedReport, error) {
	var report FeedReport

	info, statErr := os.Stat(blogDir)
	if statErr != nil || !info.IsDir() {
		return nil, report, errFs.DirNotFound(blogDir)
	}

	entries, readErr := os.ReadDir(blogDir)
	if readErr != nil {
		return nil, report, errFs.ReadDir(blogDir, readErr)
	}

	var posts []BlogPost

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !regex.BlogDateFilename.MatchString(name) {
			continue
		}

		post, status := ParsePost(filepath.Join(blogDir, name), name)

		switch status {
		case PostIncluded:
			posts = append(posts, post)
		case PostSkipped:
			report.Skipped = append(report.Skipped, post.Summary)
		case PostWarn:
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
func ParsePost(path, filename string) (BlogPost, PostStatus) {
	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		return BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(
				desc.Text(text.DescKeySiteSkipCannotRead), filename),
		}, PostSkipped
	}

	content := string(data)
	nl := token.NewlineLF
	sep := token.Separator

	if !strings.HasPrefix(content, sep+nl) {
		return BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(
				desc.Text(text.DescKeySiteSkipNoFrontmatter), filename),
		}, PostSkipped
	}

	fmStart := len(sep + nl)
	endIdx := strings.Index(content[fmStart:], nl+sep+nl)
	if endIdx < 0 {
		return BlogPost{
			Filename: filename,
			Summary:  fmt.Sprintf(desc.Text(text.DescKeySiteSkipMalformed), filename),
		}, PostSkipped
	}

	fmRaw := content[fmStart : fmStart+endIdx]
	body := content[fmStart+endIdx+len(nl+sep+nl):]

	var fm BlogFrontmatter
	if unmarshalErr := yaml.Unmarshal([]byte(fmRaw), &fm); unmarshalErr != nil {
		return BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteSkipParseError),
				filename, unmarshalErr),
		}, PostSkipped
	}

	if fm.ReviewedAndFinalized == nil || !*fm.ReviewedAndFinalized {
		return BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(
				desc.Text(text.DescKeySiteSkipNotFinalized), filename),
		}, PostSkipped
	}
	if fm.Title == "" {
		return BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteSkipMissingTitle),
				filename),
		}, PostSkipped
	}
	if fm.Date == "" {
		return BlogPost{
			Filename: filename,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteSkipMissingDate),
				filename),
		}, PostSkipped
	}

	summary := ExtractSummary(body)

	if summary == "" {
		return BlogPost{
			Filename: filename, Title: fm.Title, Date: fm.Date,
			Author: fm.Author, Topics: fm.Topics,
			Summary: fmt.Sprintf(desc.Text(text.DescKeySiteWarnNoSummary), filename),
		}, PostWarn
	}

	return BlogPost{
		Filename: filename, Title: fm.Title, Date: fm.Date,
		Author: fm.Author, Topics: fm.Topics, Summary: summary,
	}, PostIncluded
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
