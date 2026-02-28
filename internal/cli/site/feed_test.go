//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writePost creates a blog post file in dir with the given content.
func writePost(t *testing.T, dir, filename, content string) {
	t.Helper()
	if writeErr := os.WriteFile(
		filepath.Join(dir, filename),
		[]byte(content),
		0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
}

// finalizedPost returns a complete blog post with all fields.
func finalizedPost(
	title, date, author string, topics []string,
) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString("title: " + title + "\n")
	sb.WriteString("date: " + date + "\n")
	sb.WriteString("author: " + author + "\n")
	sb.WriteString("reviewed_and_finalized: true\n")
	if len(topics) > 0 {
		sb.WriteString("topics:\n")
		for _, t := range topics {
			sb.WriteString("  - " + t + "\n")
		}
	}
	sb.WriteString("---\n")
	sb.WriteString("# " + title + "\n\n")
	sb.WriteString("This is the summary paragraph.\n\n")
	sb.WriteString("More content here.\n")
	return sb.String()
}

// draftPost returns a post with reviewed_and_finalized: false.
func draftPost(title, date string) string {
	return "---\n" +
		"title: " + title + "\n" +
		"date: " + date + "\n" +
		"reviewed_and_finalized: false\n" +
		"---\n" +
		"# " + title + "\n\n" +
		"Draft content.\n"
}

func TestFeed_Basic(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-first.md",
		finalizedPost("First", "2026-01-01", "Alice", nil))
	writePost(t, dir, "2026-01-02-second.md",
		finalizedPost("Second", "2026-01-02", "Bob", nil))
	writePost(t, dir, "2026-01-03-third.md",
		finalizedPost("Third", "2026-01-03", "Carol", nil))

	posts, report, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}
	if report.included != 0 {
		// included is set by runFeed, not scanBlogPosts
	}
	if len(report.skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(report.skipped))
	}

	outPath := filepath.Join(t.TempDir(), "feed.xml")
	genErr := generateAtom(posts, outPath, "https://example.com")
	if genErr != nil {
		t.Fatal(genErr)
	}

	data, readErr := os.ReadFile(outPath)
	if readErr != nil {
		t.Fatal(readErr)
	}

	var feed AtomFeed
	if unmarshalErr := xml.Unmarshal(data, &feed); unmarshalErr != nil {
		t.Fatalf("invalid XML: %v", unmarshalErr)
	}

	if len(feed.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(feed.Entries))
	}
}

func TestFeed_SkipsDrafts(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-final.md",
		finalizedPost("Final", "2026-01-01", "Alice", nil))
	writePost(t, dir, "2026-01-02-draft.md",
		draftPost("Draft", "2026-01-02"))

	// Post with no reviewed_and_finalized field at all.
	writePost(t, dir, "2026-01-03-implicit-draft.md",
		"---\ntitle: Implicit\ndate: 2026-01-03\n---\n# Implicit\n\nContent.\n")

	posts, report, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0].title != "Final" {
		t.Errorf("expected 'Final', got %q", posts[0].title)
	}
	if len(report.skipped) != 2 {
		t.Errorf(
			"expected 2 skipped, got %d: %v",
			len(report.skipped), report.skipped,
		)
	}
}

func TestFeed_MissingTitle(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-no-title.md",
		"---\ndate: 2026-01-01\nreviewed_and_finalized: true\n---\n# Heading\n\nContent.\n")

	posts, report, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
	if len(report.skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(report.skipped))
	}
	if !strings.Contains(report.skipped[0], "missing title") {
		t.Errorf("expected 'missing title' reason, got %q",
			report.skipped[0])
	}
}

func TestFeed_MissingDate(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-no-date.md",
		"---\ntitle: No Date\nreviewed_and_finalized: true\n---\n# No Date\n\nContent.\n")

	posts, report, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
	if len(report.skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(report.skipped))
	}
	if !strings.Contains(report.skipped[0], "missing date") {
		t.Errorf("expected 'missing date' reason, got %q",
			report.skipped[0])
	}
}

func TestFeed_NoSummary(t *testing.T) {
	dir := t.TempDir()
	// Post with heading but no paragraph after it.
	writePost(t, dir, "2026-01-01-no-summary.md",
		"---\ntitle: No Summary\ndate: 2026-01-01\nreviewed_and_finalized: true\n---\n# No Summary\n")

	posts, report, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 1 {
		t.Fatalf("expected 1 post (with warning), got %d", len(posts))
	}
	if len(report.warnings) != 1 {
		t.Fatalf(
			"expected 1 warning, got %d", len(report.warnings),
		)
	}
	if !strings.Contains(report.warnings[0], "no summary") {
		t.Errorf("expected 'no summary' warning, got %q",
			report.warnings[0])
	}

	// Verify summary is omitted in XML output.
	outPath := filepath.Join(t.TempDir(), "feed.xml")
	genErr := generateAtom(posts, outPath, "https://example.com")
	if genErr != nil {
		t.Fatal(genErr)
	}

	data, _ := os.ReadFile(outPath)
	if strings.Contains(string(data), "<summary>") {
		t.Error("feed should not contain <summary> for warned post")
	}
}

func TestFeed_EmptyBlog(t *testing.T) {
	dir := t.TempDir()

	posts, _, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	outPath := filepath.Join(t.TempDir(), "feed.xml")
	genErr := generateAtom(posts, outPath, "https://example.com")
	if genErr != nil {
		t.Fatal(genErr)
	}

	data, readErr := os.ReadFile(outPath)
	if readErr != nil {
		t.Fatal(readErr)
	}

	var feed AtomFeed
	if unmarshalErr := xml.Unmarshal(data, &feed); unmarshalErr != nil {
		t.Fatalf("invalid XML for empty feed: %v", unmarshalErr)
	}
	if len(feed.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(feed.Entries))
	}
}

func TestFeed_SortOrder(t *testing.T) {
	dir := t.TempDir()
	// Write posts in non-chronological order.
	writePost(t, dir, "2026-01-15-middle.md",
		finalizedPost("Middle", "2026-01-15", "A", nil))
	writePost(t, dir, "2026-01-01-oldest.md",
		finalizedPost("Oldest", "2026-01-01", "A", nil))
	writePost(t, dir, "2026-01-30-newest.md",
		finalizedPost("Newest", "2026-01-30", "A", nil))

	posts, _, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	expected := []string{"Newest", "Middle", "Oldest"}
	for i, want := range expected {
		if posts[i].title != want {
			t.Errorf(
				"position %d: expected %q, got %q",
				i, want, posts[i].title,
			)
		}
	}
}

func TestFeed_MalformedFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-bad.md",
		"---\n: invalid: yaml: [[\n---\n# Bad\n\nContent.\n")

	posts, report, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
	if len(report.skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(report.skipped))
	}
}

func TestFeed_Idempotent(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-post.md",
		finalizedPost("Post", "2026-01-01", "Alice", []string{"go"}))

	outDir := t.TempDir()
	outPath := filepath.Join(outDir, "feed.xml")

	posts1, _, _ := scanBlogPosts(dir)
	genErr1 := generateAtom(posts1, outPath, "https://example.com")
	if genErr1 != nil {
		t.Fatal(genErr1)
	}
	data1, _ := os.ReadFile(outPath)

	posts2, _, _ := scanBlogPosts(dir)
	genErr2 := generateAtom(posts2, outPath, "https://example.com")
	if genErr2 != nil {
		t.Fatal(genErr2)
	}
	data2, _ := os.ReadFile(outPath)

	if string(data1) != string(data2) {
		t.Error("feed output is not idempotent")
	}
}

func TestFeed_Categories(t *testing.T) {
	dir := t.TempDir()
	topics := []string{"hooks", "agent behavior", "security"}
	writePost(t, dir, "2026-01-01-topics.md",
		finalizedPost("Topics", "2026-01-01", "Alice", topics))

	posts, _, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	outPath := filepath.Join(t.TempDir(), "feed.xml")
	genErr := generateAtom(posts, outPath, "https://example.com")
	if genErr != nil {
		t.Fatal(genErr)
	}

	data, _ := os.ReadFile(outPath)
	var feed AtomFeed
	xml.Unmarshal(data, &feed)

	if len(feed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(feed.Entries))
	}

	cats := feed.Entries[0].Categories
	if len(cats) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(cats))
	}

	for i, want := range topics {
		if cats[i].Term != want {
			t.Errorf(
				"category %d: expected %q, got %q",
				i, want, cats[i].Term,
			)
		}
	}
}

func TestFeed_CustomBaseURL(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "2026-01-01-post.md",
		finalizedPost("Post", "2026-01-01", "Alice", nil))

	posts, _, _ := scanBlogPosts(dir)
	outPath := filepath.Join(t.TempDir(), "feed.xml")
	genErr := generateAtom(
		posts, outPath, "https://custom.example.com",
	)
	if genErr != nil {
		t.Fatal(genErr)
	}

	data, _ := os.ReadFile(outPath)
	content := string(data)

	if !strings.Contains(content, "https://custom.example.com/blog/") {
		t.Error("feed does not use custom base URL for blog link")
	}
	if !strings.Contains(content, "https://custom.example.com/feed.xml") {
		t.Error("feed does not use custom base URL for self link")
	}
	if !strings.Contains(
		content,
		"https://custom.example.com/blog/2026-01-01-post/",
	) {
		t.Error("entry URL does not use custom base URL")
	}
}

func TestFeed_FilenameFilter(t *testing.T) {
	dir := t.TempDir()
	// Valid blog post.
	writePost(t, dir, "2026-01-01-valid.md",
		finalizedPost("Valid", "2026-01-01", "Alice", nil))
	// Non-matching filenames.
	writePost(t, dir, "index.md", "# Blog Index\n")
	writePost(t, dir, "draft-ideas.md", "# Ideas\n")
	writePost(t, dir, "README.md", "# README\n")

	posts, _, scanErr := scanBlogPosts(dir)
	if scanErr != nil {
		t.Fatal(scanErr)
	}

	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}
	if posts[0].title != "Valid" {
		t.Errorf("expected 'Valid', got %q", posts[0].title)
	}
}
