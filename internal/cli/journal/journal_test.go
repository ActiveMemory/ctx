//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "journal" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "journal")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}

func TestCmd_HasSiteSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "site" {
			found = true
			if sub.Short == "" {
				t.Error("site subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("site subcommand has no RunE function")
			}

			// Check flags
			outputFlag := sub.Flags().Lookup("output")
			if outputFlag == nil {
				t.Error("site subcommand missing --output flag")
			}

			buildFlag := sub.Flags().Lookup("build")
			if buildFlag == nil {
				t.Error("site subcommand missing --build flag")
			}

			serveFlag := sub.Flags().Lookup("serve")
			if serveFlag == nil {
				t.Error("site subcommand missing --serve flag")
			}

			break
		}
	}

	if !found {
		t.Error("site subcommand not found")
	}
}

func TestCmd_HasMarkSubcommand(t *testing.T) {
	cmd := Cmd()

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Use == "mark <filename> <stage>" {
			found = true
			if sub.Short == "" {
				t.Error("mark subcommand has empty Short description")
			}
			if sub.RunE == nil {
				t.Error("mark subcommand has no RunE function")
			}

			checkFlag := sub.Flags().Lookup("check")
			if checkFlag == nil {
				t.Error("mark subcommand missing --check flag")
			}

			break
		}
	}

	if !found {
		t.Error("mark subcommand not found")
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0B"},
		{100, "100B"},
		{1023, "1023B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{10240, "10.0KB"},
		{1048576, "1.0MB"},
		{1572864, "1.5MB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatSize(tt.bytes)
			if got != tt.want {
				t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestParseJournalEntry(t *testing.T) {
	// Create a temp file with journal content
	tmpDir := t.TempDir()
	filename := "2026-01-21-test-slug-abc12345.md"
	content := `# Test Session

**Time**: 14:30:00
**Project**: my-project

Some content here.
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	if entry.Filename != filename {
		t.Errorf("Filename = %q, want %q", entry.Filename, filename)
	}

	if entry.Date != "2026-01-21" {
		t.Errorf("Date = %q, want %q", entry.Date, "2026-01-21")
	}

	if entry.Title != "Test Session" {
		t.Errorf("Title = %q, want %q", entry.Title, "Test Session")
	}

	if entry.Time != "14:30:00" {
		t.Errorf("Time = %q, want %q", entry.Time, "14:30:00")
	}

	if entry.Project != "my-project" {
		t.Errorf("Project = %q, want %q", entry.Project, "my-project")
	}

	if entry.Size != int64(len(content)) {
		t.Errorf("Size = %d, want %d", entry.Size, len(content))
	}
}

func TestParseJournalEntry_SuggestionMode(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-01-21-suggestion-abc12345.md"
	content := `# Suggestion

[SUGGESTION MODE: some suggestion]

Content here.
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	if !entry.Suggestive {
		t.Error("Suggestive should be true for suggestion mode sessions")
	}
}

func TestParseJournalEntry_MissingFile(t *testing.T) {
	entry := parseJournalEntry("/nonexistent/path.md", "2026-01-21-test.md")

	// Should use filename as title fallback
	if entry.Title != "2026-01-21-test" {
		t.Errorf("Title = %q, want %q", entry.Title, "2026-01-21-test")
	}
}

func TestGenerateIndex(t *testing.T) {
	entries := []journalEntry{
		{
			Filename: "2026-01-21-session-one-abc12345.md",
			Title:    "Session One",
			Date:     "2026-01-21",
			Time:     "14:30:00",
			Project:  "project-a",
			Size:     1024,
		},
		{
			Filename: "2026-01-20-session-two-def67890.md",
			Title:    "Session Two",
			Date:     "2026-01-20",
			Time:     "10:00:00",
			Project:  "project-b",
			Size:     2048,
		},
		{
			Filename:   "2026-01-19-suggestion-ghi11111.md",
			Title:      "Suggestion",
			Date:       "2026-01-19",
			Time:       "09:00:00",
			Suggestive: true,
			Size:       512,
		},
	}

	index := generateIndex(entries)

	// Should have header
	if !strings.Contains(index, "# Session Journal") {
		t.Error("index missing header")
	}

	// Should have session count
	if !strings.Contains(index, "**Sessions**: 2") {
		t.Error("index missing session count")
	}

	// Should have suggestions count
	if !strings.Contains(index, "**Suggestions**: 1") {
		t.Error("index missing suggestions count")
	}

	// Should have month headers
	if !strings.Contains(index, "## 2026-01") {
		t.Error("index missing month header")
	}

	// Should have entry links
	if !strings.Contains(index, "[Session One]") {
		t.Error("index missing session one link")
	}

	// Should have suggestions section
	if !strings.Contains(index, "## Suggestions") {
		t.Error("index missing suggestions section")
	}
}

func TestInjectSourceLink_WithFrontmatter(t *testing.T) {
	content := "---\ntitle: Test\n---\n\n# Heading\n"
	result := injectSourceLink(content, "/home/user/.context/journal/test.md")

	// Should have file:// link
	if !strings.Contains(result, "[View source](file:///home/user/.context/journal/test.md)") {
		t.Errorf("missing file:// link:\n%s", result)
	}
	// Should have copyable path with copy button
	if !strings.Contains(result, ".context/journal/test.md") {
		t.Errorf("missing relative path:\n%s", result)
	}
	// Original content should still be there
	if !strings.Contains(result, "# Heading") {
		t.Error("original content missing")
	}
}

func TestInjectSourceLink_NoFrontmatter(t *testing.T) {
	content := "# Heading\n\nSome text.\n"
	result := injectSourceLink(content, "/path/to/file.md")

	// Link should be at the top
	if !strings.HasPrefix(result, "*[View source](file:///path/to/file.md)") {
		t.Errorf("source link not at top:\n%s", result)
	}
	if !strings.Contains(result, ".context/journal/file.md") {
		t.Errorf("missing relative path:\n%s", result)
	}
	if !strings.Contains(result, "# Heading") {
		t.Error("original content missing")
	}
}

func TestNormalizeContent(t *testing.T) {
	tests := []struct {
		name, input    string
		fencesVerified bool
		check          func(t *testing.T, got string)
	}{
		{
			"strips tool bold",
			`🔧 **Glob: .context/journal/*.md**`,
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "**Glob") {
					t.Error("bold markers not stripped from tool line")
				}
				if !strings.Contains(got, "🔧 Glob:") {
					t.Error("tool prefix missing")
				}
			},
		},
		{
			"escapes glob stars",
			`pattern: src/*/main.go`,
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, `\*/`) {
					t.Error("glob star not escaped")
				}
			},
		},
		{
			"strips fences and escapes content",
			"```\n*.md\n```",
			false,
			func(t *testing.T, got string) {
				// Fences should be stripped
				if strings.Contains(got, "```") {
					t.Error("fence markers should be stripped")
				}
				// Content survives (after fence strip, *.md may be escaped)
			},
		},
		{
			"skips frontmatter",
			"---\ntitle: test\n---\nsome text",
			false,
			func(t *testing.T, got string) {
				if !strings.HasPrefix(got, "---\ntitle: test\n---\n") {
					t.Errorf("frontmatter mangled: %q", got)
				}
			},
		},
		{
			"does not wrap (site output is read-only)",
			"This is a very long line that exceeds eighty characters and should not be wrapped since the site output is read-only.",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "\n") {
					t.Error("normalizeContent should not wrap lines")
				}
			},
		},
		{
			"inline code with angle brackets gets quoted",
			"the link text contains `</com` which is broken",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "`</com`") {
					t.Error("backtick code with angle bracket should be replaced")
				}
				if !strings.Contains(got, `"&lt;/com"`) {
					t.Errorf("expected quoted entity, got: %s", got)
				}
			},
		},
		{
			"inline code without angles is untouched",
			"run `ctx status` to check",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "`ctx status`") {
					t.Error("safe inline code should not be modified")
				}
			},
		},
		{
			"inline code with both angles",
			"found `<div>` tag in output",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, `"&lt;div&gt;"`) {
					t.Errorf("expected quoted entities, got: %s", got)
				}
			},
		},
		{
			"H1 with Claude tags gets sanitized",
			"# <command-message>ctx:ctx-journal-normalize</command-message> more text",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "command-message") {
					t.Error("Claude tags should be stripped from H1")
				}
				if !strings.HasPrefix(got, "# ctx:ctx-journal-normalize more text") {
					t.Errorf("unexpected H1: %s", got)
				}
			},
		},
		{
			"long H1 gets truncated",
			"# " + strings.Repeat("word ", 20),
			false,
			func(t *testing.T, got string) {
				heading := strings.TrimPrefix(got, "# ")
				if len([]rune(heading)) > 75 {
					t.Errorf("H1 not truncated: %d runes", len([]rune(heading)))
				}
			},
		},
		{
			"tool output wrapped in pre/code",
			"### 5. Tool Output (10:30:00)\n\n# this is not a heading\n---\n<details>bad\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("tool output should be wrapped in <pre><code>")
				}
				// Content should be HTML-escaped inside <pre><code>
				if !strings.Contains(got, "# this is not a heading") {
					t.Error("# line should be preserved")
				}
				if !strings.Contains(got, "&lt;details&gt;bad") {
					t.Error("<details> should be HTML-escaped")
				}
				// Next turn header should survive
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("next turn header should not be consumed")
				}
			},
		},
		{
			"tool output in details gets re-wrapped as pre/code",
			"### 5. Tool Output (10:30:00)\n\n<details>\n<summary>79 lines</summary>\n<pre>\n# heading\n---\n&lt;div&gt;\n</pre>\n</details>\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("tool output should be wrapped in <pre><code>")
				}
				// <details>/<summary>/<pre> wrappers stripped
				if strings.Contains(got, "<summary>") {
					t.Error("<summary> wrapper should be stripped")
				}
				// Content unescaped from export then re-escaped for <pre><code>:
				// &lt;div&gt; -> <div> -> &lt;div&gt;
				if !strings.Contains(got, "&lt;div&gt;") {
					t.Error("HTML entities should be re-escaped in <pre><code>")
				}
				if !strings.Contains(got, "# heading") {
					t.Error("# heading should be preserved")
				}
				// Next turn should survive
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("next turn header should not be consumed")
				}
			},
		},
		{
			"boilerplate tool output stripped — empty body",
			"### 5. Tool Output (10:30:00)\n\n\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("empty tool output header should be stripped")
				}
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("next turn should survive")
				}
			},
		},
		{
			"boilerplate tool output stripped — no matches found",
			"### 5. Tool Output (10:30:00)\n\nNo matches found\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("'No matches found' tool output should be stripped")
				}
				if strings.Contains(got, "No matches found") {
					t.Error("boilerplate line should not appear")
				}
			},
		},
		{
			"boilerplate tool output stripped — edit confirmation",
			"### 5. Tool Output (10:30:00)\n\nThe file /home/jose/foo.go has been updated successfully.\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("edit confirmation tool output should be stripped")
				}
			},
		},
		{
			"boilerplate tool output stripped — hook denial",
			"### 5. Tool Output (10:30:00)\n\nHook PreToolUse:Bash denied this tool\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("hook denial tool output should be stripped")
				}
			},
		},
		{
			"non-boilerplate tool output preserved",
			"### 5. Tool Output (10:30:00)\n\nactual useful content here\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "actual useful content here") {
					t.Error("non-boilerplate content should be preserved")
				}
				if !strings.Contains(got, "<pre><code>") {
					t.Error("tool output should be wrapped in <pre><code>")
				}
			},
		},
		{
			"boilerplate stripped — multi-line edit confirmation",
			"### 5. Tool Output (10:30:00)\n\nThe file /home/jose/WORKSPACE/ctx/internal/config/limit.go has been updated\nsuccessfully.\n\n### 6. Assistant (10:30:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "Tool Output") {
					t.Error("multi-line edit confirmation should be stripped")
				}
				if strings.Contains(got, "has been updated") {
					t.Error("boilerplate content should not appear")
				}
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("next turn should survive")
				}
			},
		},
		{
			"fencesVerified wraps already-fenced tool output in pre/code",
			"### 5. Tool Output (10:30:00)\n\n```\nactual output\n```\n\n### 6. Assistant (10:30:01)\n\nhi",
			true,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("tool output should be wrapped in <pre><code>")
				}
				if !strings.Contains(got, "actual output") {
					t.Error("content should be preserved")
				}
			},
		},
		{
			"fencesVerified converts details/pre to pre/code",
			"### 5. Tool Output (10:30:00)\n\n<details>\n<summary>3 lines</summary>\n\n<pre>\nfoo\n&lt;bar&gt;\n</pre>\n</details>\n\n### 6. Assistant (10:30:01)\n\nhi",
			true,
			func(t *testing.T, got string) {
				if strings.Contains(got, "<summary>") {
					t.Error("<summary> wrapper should be stripped")
				}
				if !strings.Contains(got, "<pre><code>") {
					t.Error("should wrap in <pre><code>")
				}
				// Unescaped from export then re-escaped
				if !strings.Contains(got, "&lt;bar&gt;") {
					t.Error("HTML entities should be re-escaped")
				}
			},
		},
		{
			"inner fences become literal text in pre/code",
			"### 5. Tool Output (10:30:00)\n\n<details>\n<summary>5 lines</summary>\n\n<pre>\n## Heading\n\n```\ncode block\n```\n</pre>\n</details>\n\n### 6. Assistant (10:30:01)\n\nhi",
			true,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("should wrap in <pre><code>")
				}
				// Inner fences are just literal text (HTML-escaped has no effect on backticks)
				if !strings.Contains(got, "```") {
					t.Error("inner fence markers should be preserved as literal text")
				}
				if !strings.Contains(got, "## Heading") {
					t.Error("content should be preserved")
				}
			},
		},
		{
			"embedded turn headers inside pre are not boundaries",
			"### 5. Tool Output (10:30:00)\n\n<details>\n<summary>3 lines</summary>\n\n<pre>\n### 800. Assistant (15:00:00)\n\nembedded content\n</pre>\n</details>\n\n### 6. Assistant (10:30:01)\n\nreal next turn",
			true,
			func(t *testing.T, got string) {
				// Embedded turn header should be inside the pre/code block
				if !strings.Contains(got, "### 800. Assistant") {
					t.Error("embedded turn header should be preserved in output")
				}
				if !strings.Contains(got, "embedded content") {
					t.Error("embedded content should be in the pre/code block")
				}
				// Real next turn must survive
				if !strings.Contains(got, "### 6. Assistant") {
					t.Error("real next turn should not be swallowed")
				}
				if !strings.Contains(got, "real next turn") {
					t.Error("real next turn content should survive")
				}
			},
		},
		{
			"user turn wrapped in pre/code",
			"### 1. User (10:00:00)\n\nHello world\n\n### 2. Assistant (10:00:01)\n\nHi there",
			false,
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("user body should be wrapped in <pre><code>")
				}
				if !strings.Contains(got, "</code></pre>") {
					t.Error("user body should have closing </code></pre>")
				}
				if !strings.Contains(got, "Hello world") {
					t.Error("user content should be preserved")
				}
				if !strings.Contains(got, "### 2. Assistant") {
					t.Error("next turn should survive")
				}
			},
		},
		{
			"user turn with stray fence does not swallow subsequent turns",
			"### 1. User (10:00:00)\n\nsome text\n```\nmore text\n\n### 2. Assistant (10:00:01)\n\nresponse here",
			true, // fencesVerified — the dangerous case
			func(t *testing.T, got string) {
				// Stray fence should be HTML-escaped inside <pre><code>
				if !strings.Contains(got, "<pre><code>") {
					t.Error("user body should be wrapped in <pre><code>")
				}
				// The ``` inside <pre><code> is inert to markdown parsing.
				// Backticks aren't HTML-escaped by html.EscapeString,
				// so we accept either raw ``` or &#96;&#96;&#96; here.
				// Critical: next turn must NOT be swallowed
				if !strings.Contains(got, "### 2. Assistant") {
					t.Error("stray fence must not swallow subsequent turns")
				}
				if !strings.Contains(got, "response here") {
					t.Error("assistant content should survive")
				}
			},
		},
		{
			"user turn with HTML is escaped",
			"### 1. User (10:00:00)\n\n<script>alert('xss')</script>\n\n### 2. Assistant (10:00:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "<script>") {
					t.Error("HTML in user body should be escaped")
				}
				if !strings.Contains(got, "&lt;script&gt;") {
					t.Error("HTML should be entity-escaped")
				}
			},
		},
		{
			"empty user turn not wrapped",
			"### 1. User (10:00:00)\n\n\n\n### 2. Assistant (10:00:01)\n\nhi",
			false,
			func(t *testing.T, got string) {
				if strings.Contains(got, "<pre><code>") {
					t.Error("empty user turn should not get pre/code wrapper")
				}
				if !strings.Contains(got, "### 2. Assistant") {
					t.Error("next turn should survive")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeContent(tt.input, tt.fencesVerified)
			tt.check(t, got)
		})
	}
}

func TestWrapUserTurns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, got string)
	}{
		{
			"basic user turn",
			"### 1. User (10:00:00)\n\nHello\n\n### 2. Assistant (10:00:01)\n\nHi",
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>\nHello\n</code></pre>") {
					t.Errorf("unexpected wrapping:\n%s", got)
				}
				if !strings.Contains(got, "### 2. Assistant") {
					t.Error("next turn lost")
				}
			},
		},
		{
			"user turn with angle brackets escaped",
			"### 1. User (10:00:00)\n\ncheck <div> tag\n\n### 2. Assistant (10:00:01)\n\nok",
			func(t *testing.T, got string) {
				if !strings.Contains(got, "check &lt;div&gt; tag") {
					t.Errorf("HTML not escaped:\n%s", got)
				}
			},
		},
		{
			"user turn with ampersand escaped",
			"### 1. User (10:00:00)\n\nfoo & bar\n\n### 2. Assistant (10:00:01)\n\nok",
			func(t *testing.T, got string) {
				if !strings.Contains(got, "foo &amp; bar") {
					t.Errorf("ampersand not escaped:\n%s", got)
				}
			},
		},
		{
			"user turn at EOF",
			"### 1. User (10:00:00)\n\nlast message",
			func(t *testing.T, got string) {
				if !strings.Contains(got, "<pre><code>") {
					t.Error("user turn at EOF should still be wrapped")
				}
				if !strings.Contains(got, "last message") {
					t.Error("content should be preserved")
				}
			},
		},
		{
			"multiple user turns",
			"### 1. User (10:00:00)\n\nfirst\n\n### 2. Assistant (10:00:01)\n\nhi\n\n### 3. User (10:00:02)\n\nsecond\n\n### 4. Assistant (10:00:03)\n\nbye",
			func(t *testing.T, got string) {
				if strings.Count(got, "<pre><code>") != 2 {
					t.Errorf("expected 2 pre/code blocks, got %d",
						strings.Count(got, "<pre><code>"))
				}
				if !strings.Contains(got, "first") {
					t.Error("first user message lost")
				}
				if !strings.Contains(got, "second") {
					t.Error("second user message lost")
				}
			},
		},
		{
			"non-user turns untouched",
			"### 1. Assistant (10:00:00)\n\nI'll help\n\n### 2. Tool Output (10:00:01)\n\nresult",
			func(t *testing.T, got string) {
				if strings.Contains(got, "<pre><code>") {
					t.Error("non-user turns should not be wrapped")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapUserTurns(tt.input)
			tt.check(t, got)
		})
	}
}

func TestSoftWrap(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  int // expected number of output lines
	}{
		{"short line", "hello world", 80, 1},
		{"long line wraps", "word " + strings.Repeat("x ", 50), 80, 2},
		{"preserves indent", "    indented " + strings.Repeat("word ", 20), 80, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := softWrap(tt.input, tt.width)
			if len(lines) < tt.want {
				t.Errorf("got %d lines, want >= %d", len(lines), tt.want)
			}
			// Verify indent preserved on continuation
			if strings.HasPrefix(tt.input, "    ") && len(lines) > 1 {
				if !strings.HasPrefix(lines[1], "    ") {
					t.Errorf("continuation lost indent: %q", lines[1])
				}
			}
		})
	}
}

func TestConsolidateToolRuns(t *testing.T) {
	input := strings.Join([]string{
		"Some text before.",
		"",
		"### 30. Assistant (04:41:50)",
		"",
		"🔧 **TaskCreate**",
		"",
		"### 31. Assistant (04:41:51)",
		"",
		"🔧 **TaskCreate**",
		"",
		"### 32. Assistant (04:41:53)",
		"",
		"🔧 **TaskCreate**",
		"",
		"### 33. Tool Output (04:41:58)",
		"",
		"Some output.",
	}, "\n")

	got := consolidateToolRuns(input)

	// Should collapse 3 TaskCreate into 1
	if !strings.Contains(got, "🔧 **TaskCreate**") || !strings.Contains(got, "(×3)") {
		t.Errorf("expected consolidated count:\n%s", got)
	}

	// Should keep the first header
	if !strings.Contains(got, "### 30. Assistant (04:41:50)") {
		t.Error("first header missing")
	}

	// Should NOT keep the duplicate headers
	if strings.Contains(got, "### 31. Assistant") {
		t.Error("duplicate header not removed")
	}
	if strings.Contains(got, "### 32. Assistant") {
		t.Error("duplicate header not removed")
	}

	// Should keep non-tool content
	if !strings.Contains(got, "Some text before.") {
		t.Error("surrounding content lost")
	}
	if !strings.Contains(got, "### 33. Tool Output") {
		t.Error("following turn lost")
	}
}

func TestConsolidateToolRuns_DifferentTools(t *testing.T) {
	input := strings.Join([]string{
		"### 10. Assistant (04:00:00)",
		"",
		"🔧 **TaskCreate**",
		"",
		"### 11. Assistant (04:00:01)",
		"",
		"🔧 **TaskUpdate**",
		"",
	}, "\n")

	got := consolidateToolRuns(input)

	// Different tools should NOT be consolidated
	if strings.Contains(got, "×") {
		t.Error("different tools should not be consolidated")
	}
}

func TestConsolidateToolRuns_ToolOutput(t *testing.T) {
	input := strings.Join([]string{
		"### 140. Tool Output (04:46:21)",
		"",
		"The file DECISIONS.md has been updated successfully.",
		"",
		"### 141. Tool Output (04:46:21)",
		"",
		"The file DECISIONS.md has been updated successfully.",
		"",
		"### 142. Tool Output (04:46:22)",
		"",
		"The file DECISIONS.md has been updated successfully.",
		"",
		"### 143. Assistant (04:46:23)",
		"",
		"Done with updates.",
	}, "\n")

	got := consolidateToolRuns(input)

	// Should collapse 3 identical outputs
	if !strings.Contains(got, "(×3)") {
		t.Errorf("expected ×3 count:\n%s", got)
	}

	// Should keep first header
	if !strings.Contains(got, "### 140. Tool Output") {
		t.Error("first header missing")
	}

	// Should not keep duplicates
	if strings.Contains(got, "### 141.") || strings.Contains(got, "### 142.") {
		t.Error("duplicate headers not removed")
	}

	// Should keep the following different turn
	if !strings.Contains(got, "### 143. Assistant") {
		t.Error("following turn lost")
	}
	if !strings.Contains(got, "Done with updates.") {
		t.Error("following content lost")
	}
}

func TestConsolidateToolRuns_SingleTurn(t *testing.T) {
	input := "### 10. Assistant (04:00:00)\n\n🔧 **Read**\n\nSome text."

	got := consolidateToolRuns(input)

	// Single tool turn should be unchanged
	if strings.Contains(got, "×") {
		t.Error("single turn should not get a count")
	}
}

func TestConsolidateToolRuns_FenceSafe(t *testing.T) {
	// The (×N) annotation must NOT be appended to a closing fence line,
	// otherwise "``` (×2)" parses as an opening fence with info string.
	input := strings.Join([]string{
		"### 50. Tool Output (05:00:00)",
		"",
		"```",
		"some output",
		"```",
		"",
		"### 51. Tool Output (05:00:01)",
		"",
		"```",
		"some output",
		"```",
		"",
		"### 52. Assistant (05:00:02)",
		"",
		"Done.",
	}, "\n")

	got := consolidateToolRuns(input)

	if !strings.Contains(got, "(×2)") {
		t.Errorf("expected ×2 count:\n%s", got)
	}
	// The closing ``` must NOT have (×N) appended to it
	if strings.Contains(got, "``` (×2)") || strings.Contains(got, "```(×2)") {
		t.Errorf("(×N) must not be appended to closing fence line:\n%s", got)
	}
}

func TestSoftWrapContent(t *testing.T) {
	long := "This is a very long line that exceeds eighty characters and should be wrapped at a word boundary somewhere."
	input := "---\ntitle: test\n---\n\n" + long + "\n\n```\n" + long + "\n```\n"

	got := softWrapContent(input)

	// All lines should be wrapped (including inside code fences)
	for _, line := range strings.Split(got, "\n") {
		if len(line) > 85 { // allow slack for word boundaries
			t.Errorf("line too long (%d): %q", len(line), line)
		}
	}

	// Frontmatter should be intact
	if !strings.HasPrefix(got, "---\ntitle: test\n---\n") {
		t.Error("frontmatter mangled")
	}
}

func TestFormatIndexEntry(t *testing.T) {
	entry := journalEntry{
		Filename: "2026-01-21-test-abc12345.md",
		Title:    "Test Session",
		Date:     "2026-01-21",
		Time:     "14:30:00",
		Project:  "my-project",
		Size:     1536,
	}

	result := formatIndexEntry(entry, "\n")

	// Should have time prefix
	if !strings.Contains(result, "14:30") {
		t.Error("entry missing time prefix")
	}

	// Should have title link
	if !strings.Contains(result, "[Test Session]") {
		t.Error("entry missing title")
	}

	// Should have link to md file
	if !strings.Contains(result, "(2026-01-21-test-abc12345.md)") {
		t.Error("entry missing link")
	}

	// Should have project
	if !strings.Contains(result, "(my-project)") {
		t.Error("entry missing project")
	}

	// Should have size
	if !strings.Contains(result, "1.5KB") {
		t.Error("entry missing size")
	}
}

func TestParseJournalEntry_WithFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-02-04--80ac7de2.md"
	content := `---
title: "Skill audit: distill prompt files"
date: 2026-02-04
type: refactor
outcome: completed
topics:
  - skills
  - conventions
  - code-quality
key_files:
  - internal/cli/run.go
  - cmd/main.go
---

# 2026-02-04--80ac7de2

**Time**: 09:15:00
**Project**: ctx

Some content here.
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	// Frontmatter title should override H1
	if entry.Title != "Skill audit: distill prompt files" {
		t.Errorf("Title = %q, want frontmatter title", entry.Title)
	}

	if entry.Type != "refactor" {
		t.Errorf("Type = %q, want %q", entry.Type, "refactor")
	}

	if entry.Outcome != "completed" {
		t.Errorf("Outcome = %q, want %q", entry.Outcome, "completed")
	}

	if len(entry.Topics) != 3 {
		t.Fatalf("Topics length = %d, want 3", len(entry.Topics))
	}
	if entry.Topics[0] != "skills" {
		t.Errorf("Topics[0] = %q, want %q", entry.Topics[0], "skills")
	}

	if len(entry.KeyFiles) != 2 {
		t.Fatalf("KeyFiles length = %d, want 2", len(entry.KeyFiles))
	}
	if entry.KeyFiles[0] != "internal/cli/run.go" {
		t.Errorf("KeyFiles[0] = %q, want %q", entry.KeyFiles[0], "internal/cli/run.go")
	}

	// Line-by-line fallback should still work
	if entry.Time != "09:15:00" {
		t.Errorf("Time = %q, want %q", entry.Time, "09:15:00")
	}
	if entry.Project != "ctx" {
		t.Errorf("Project = %q, want %q", entry.Project, "ctx")
	}
}

func TestBuildTopicIndex(t *testing.T) {
	entries := []journalEntry{
		{Filename: "a.md", Date: "2026-01-21", Topics: []string{"go", "cli"}},
		{Filename: "b.md", Date: "2026-01-22", Topics: []string{"go", "testing"}},
		{Filename: "c.md", Date: "2026-01-23", Topics: []string{"cli"}},
		{Filename: "d.md", Date: "2026-01-24", Topics: []string{"docs"}},
	}

	topics := buildTopicIndex(entries)

	if len(topics) != 4 {
		t.Fatalf("got %d topics, want 4", len(topics))
	}

	// First should be most popular: "cli" and "go" both have 2, alpha order means "cli" first
	if topics[0].Name != "cli" {
		t.Errorf("topics[0].Name = %q, want %q", topics[0].Name, "cli")
	}
	if topics[1].Name != "go" {
		t.Errorf("topics[1].Name = %q, want %q", topics[1].Name, "go")
	}

	// Popular flag
	if !topics[0].Popular {
		t.Error("cli should be popular (2 sessions)")
	}
	if !topics[1].Popular {
		t.Error("go should be popular (2 sessions)")
	}

	// Long-tail topics (1 session) should not be popular
	for _, topic := range topics[2:] {
		if topic.Popular {
			t.Errorf("%q should not be popular (1 session)", topic.Name)
		}
	}
}

func TestGenerateTopicsIndex(t *testing.T) {
	topics := []topicData{
		{Name: "go", Entries: []journalEntry{
			{Filename: "a.md", Title: "Session A"},
			{Filename: "b.md", Title: "Session B"},
		}, Popular: true},
		{Name: "docs", Entries: []journalEntry{
			{Filename: "c.md", Title: "Session C"},
		}, Popular: false},
	}

	index := generateTopicsIndex(topics)

	if !strings.Contains(index, "# Topics") {
		t.Error("missing header")
	}
	if !strings.Contains(index, "## Popular Topics") {
		t.Error("missing popular section")
	}
	if !strings.Contains(index, "## Long-tail Topics") {
		t.Error("missing long-tail section")
	}
	// Popular topics link to dedicated pages
	if !strings.Contains(index, "[go](go.md)") {
		t.Error("popular topic should link to go.md")
	}
	// Long-tail topics link inline to session
	if !strings.Contains(index, "[Session C](../c.md)") {
		t.Error("long-tail topic should link to session inline")
	}
}

func TestGenerateTopicPage(t *testing.T) {
	topic := topicData{
		Name: "cli",
		Entries: []journalEntry{
			{Filename: "2026-02-01-a.md", Title: "Session A", Date: "2026-02-01", Time: "14:30:00"},
			{Filename: "2026-01-15-b.md", Title: "Session B", Date: "2026-01-15", Time: "09:00:00"},
		},
	}

	page := generateTopicPage(topic)

	if !strings.Contains(page, "# cli") {
		t.Error("missing topic title")
	}
	if !strings.Contains(page, "**2 sessions**") {
		t.Error("missing session count")
	}
	// Month grouping
	if !strings.Contains(page, "## 2026-02") {
		t.Error("missing month group 2026-02")
	}
	if !strings.Contains(page, "## 2026-01") {
		t.Error("missing month group 2026-01")
	}
	// Relative links
	if !strings.Contains(page, "(../2026-02-01-a.md)") {
		t.Error("missing relative link to session")
	}
}

func TestBuildKeyFileIndex(t *testing.T) {
	entries := []journalEntry{
		{Filename: "a.md", Date: "2026-01-21", KeyFiles: []string{"cmd/main.go", "internal/cli/run.go"}},
		{Filename: "b.md", Date: "2026-01-22", KeyFiles: []string{"cmd/main.go", "README.md"}},
		{Filename: "c.md", Date: "2026-01-23", KeyFiles: []string{"internal/cli/run.go"}},
		{Filename: "d.md", Date: "2026-01-24", KeyFiles: []string{"go.mod"}},
	}

	keyFiles := buildKeyFileIndex(entries)

	if len(keyFiles) != 4 {
		t.Fatalf("got %d key files, want 4", len(keyFiles))
	}

	// Most popular first
	if keyFiles[0].Path != "cmd/main.go" {
		t.Errorf("keyFiles[0].Path = %q, want %q", keyFiles[0].Path, "cmd/main.go")
	}
	if !keyFiles[0].Popular {
		t.Error("cmd/main.go should be popular (2 sessions)")
	}
	if !keyFiles[1].Popular {
		t.Error("internal/cli/run.go should be popular (2 sessions)")
	}

	// Long-tail
	for _, kf := range keyFiles[2:] {
		if kf.Popular {
			t.Errorf("%q should not be popular (1 session)", kf.Path)
		}
	}
}

func TestGenerateKeyFilesIndex(t *testing.T) {
	keyFiles := []keyFileData{
		{Path: "cmd/main.go", Entries: []journalEntry{
			{Filename: "a.md", Title: "Session A"},
			{Filename: "b.md", Title: "Session B"},
		}, Popular: true},
		{Path: "go.mod", Entries: []journalEntry{
			{Filename: "c.md", Title: "Session C"},
		}, Popular: false},
	}

	index := generateKeyFilesIndex(keyFiles)

	if !strings.Contains(index, "# Key Files") {
		t.Error("missing header")
	}
	if !strings.Contains(index, "## Frequently Touched") {
		t.Error("missing popular section")
	}
	if !strings.Contains(index, "## Single Session") {
		t.Error("missing long-tail section")
	}
	// Popular files link to dedicated pages
	slug := keyFileSlug("cmd/main.go")
	if !strings.Contains(index, slug+".md") {
		t.Errorf("popular file should link to %s.md", slug)
	}
	// Long-tail files link inline to session
	if !strings.Contains(index, "[Session C](../c.md)") {
		t.Error("long-tail file should link to session inline")
	}
}

func TestGenerateKeyFilePage(t *testing.T) {
	kf := keyFileData{
		Path: "internal/cli/run.go",
		Entries: []journalEntry{
			{Filename: "2026-02-01-a.md", Title: "Session A", Date: "2026-02-01", Time: "14:30:00"},
			{Filename: "2026-01-15-b.md", Title: "Session B", Date: "2026-01-15", Time: "09:00:00"},
		},
	}

	page := generateKeyFilePage(kf)

	if !strings.Contains(page, "# `internal/cli/run.go`") {
		t.Error("missing file path title")
	}
	if !strings.Contains(page, "**2 sessions**") {
		t.Error("missing session count")
	}
	if !strings.Contains(page, "## 2026-02") {
		t.Error("missing month group")
	}
	if !strings.Contains(page, "(../2026-02-01-a.md)") {
		t.Error("missing relative link to session")
	}
}

func TestKeyFileSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"cmd/main.go", "cmd_main_go"},
		{"internal/cli/run.go", "internal_cli_run_go"},
		{".context/journal/*.md", "_context_journal_x_md"},
	}
	for _, tt := range tests {
		got := keyFileSlug(tt.input)
		if got != tt.want {
			t.Errorf("keyFileSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildTypeIndex(t *testing.T) {
	entries := []journalEntry{
		{Filename: "a.md", Date: "2026-01-21", Type: "feature"},
		{Filename: "b.md", Date: "2026-01-22", Type: "bugfix"},
		{Filename: "c.md", Date: "2026-01-23", Type: "feature"},
		{Filename: "d.md", Date: "2026-01-24", Type: "feature"},
	}

	types := buildTypeIndex(entries)

	if len(types) != 2 {
		t.Fatalf("got %d types, want 2", len(types))
	}
	// Most popular first
	if types[0].Name != "feature" {
		t.Errorf("types[0].Name = %q, want %q", types[0].Name, "feature")
	}
	if len(types[0].Entries) != 3 {
		t.Errorf("feature should have 3 entries, got %d", len(types[0].Entries))
	}
}

func TestGenerateTypesIndex(t *testing.T) {
	types := []typeData{
		{Name: "feature", Entries: []journalEntry{
			{Filename: "a.md"}, {Filename: "b.md"}, {Filename: "c.md"},
		}},
		{Name: "bugfix", Entries: []journalEntry{
			{Filename: "d.md"},
		}},
	}

	index := generateTypesIndex(types)

	if !strings.Contains(index, "# Session Types") {
		t.Error("missing header")
	}
	if !strings.Contains(index, "[feature](feature.md)") {
		t.Error("missing link to feature page")
	}
	if !strings.Contains(index, "(3 sessions)") {
		t.Error("missing session count for feature")
	}
}

func TestGenerateTypePage(t *testing.T) {
	st := typeData{
		Name: "bugfix",
		Entries: []journalEntry{
			{Filename: "2026-02-01-a.md", Title: "Fix login", Date: "2026-02-01", Time: "14:30:00"},
			{Filename: "2026-01-15-b.md", Title: "Fix crash", Date: "2026-01-15", Time: "09:00:00"},
		},
	}

	page := generateTypePage(st)

	if !strings.Contains(page, "# bugfix") {
		t.Error("missing type title")
	}
	if !strings.Contains(page, "**2 sessions**") {
		t.Error("missing session count")
	}
	if !strings.Contains(page, "## 2026-02") {
		t.Error("missing month group")
	}
	if !strings.Contains(page, "(../2026-02-01-a.md)") {
		t.Error("missing relative link")
	}
}

func TestGenerateZensicalToml_WithAllNav(t *testing.T) {
	entries := []journalEntry{{Filename: "a.md", Title: "A"}}
	topics := []topicData{{Name: "go"}}
	keyFiles := []keyFileData{{Path: "main.go"}}
	types := []typeData{{Name: "feature"}}

	toml := generateZensicalToml(entries, topics, keyFiles, types)

	if !strings.Contains(toml, `"Topics" = "topics/index.md"`) {
		t.Error("missing Topics nav")
	}
	if !strings.Contains(toml, `"Files" = "files/index.md"`) {
		t.Error("missing Files nav")
	}
	if !strings.Contains(toml, `"Types" = "types/index.md"`) {
		t.Error("missing Types nav")
	}
}

func TestGenerateZensicalToml(t *testing.T) {
	entries := []journalEntry{
		{
			Filename: "2026-01-21-test.md",
			Title:    "Test Session",
		},
	}

	topics := []topicData{
		{Name: "go", Popular: true},
	}

	toml := generateZensicalToml(entries, topics, nil, nil)

	// Verify required structural elements exist (not exact content)
	requiredPatterns := []struct {
		pattern string
		desc    string
	}{
		{"[project]", "project section"},
		{"site_name = ", "site_name field"},
		{"nav = [", "nav array"},
		{"[project.theme]", "theme section"},
		{`"Topics" = "topics/index.md"`, "topics nav entry"},
	}

	for _, tc := range requiredPatterns {
		if !strings.Contains(toml, tc.pattern) {
			t.Errorf("toml missing %s (expected %q)", tc.desc, tc.pattern)
		}
	}
}

func TestStripFences(t *testing.T) {
	input := "# Heading\n\n```go\nfunc main() {}\n```\n\nMore text.\n"
	got := stripFences(input, false)

	// Fence markers removed
	if strings.Contains(got, "```") {
		t.Error("fence markers not stripped")
	}
	// Content preserved
	if !strings.Contains(got, "func main() {}") {
		t.Error("content inside fence lost")
	}
	if !strings.Contains(got, "# Heading") || !strings.Contains(got, "More text.") {
		t.Error("surrounding content lost")
	}
}

func TestStripFences_PreservesFrontmatter(t *testing.T) {
	input := "---\ntitle: test\n---\n\n```\ncode\n```\n"
	got := stripFences(input, false)

	if !strings.HasPrefix(got, "---\ntitle: test\n---\n") {
		t.Error("frontmatter damaged")
	}
	if strings.Contains(got, "```") {
		t.Error("fence not stripped after frontmatter")
	}
}

func TestStripFences_SkipsFencesVerified(t *testing.T) {
	input := "```go\ncode\n```\n"
	got := stripFences(input, true)

	// Should be unchanged — fences verified via state
	if got != input {
		t.Error("should skip files when fencesVerified is true")
	}
}

func TestStripFences_NestedFences(t *testing.T) {
	input := "````\n```python\nprint('hi')\n```\n````\n"
	got := stripFences(input, false)

	// All fence markers removed
	if strings.Contains(got, "```") || strings.Contains(got, "````") {
		t.Error("nested fence markers not stripped")
	}
	if !strings.Contains(got, "print('hi')") {
		t.Error("content lost")
	}
}

func TestStripSystemReminders(t *testing.T) {
	input := strings.Join([]string{
		"before",
		"",
		"<system-reminder>",
		"Some internal info.",
		"Another line.",
		"</system-reminder>",
		"",
		"after",
	}, "\n")

	got := stripSystemReminders(input)

	if strings.Contains(got, "<system-reminder>") {
		t.Error("system-reminder tag not stripped")
	}
	if strings.Contains(got, "Some internal info.") {
		t.Error("reminder content not stripped")
	}
	if !strings.Contains(got, "before") || !strings.Contains(got, "after") {
		t.Error("surrounding content lost")
	}
}

func TestStripSystemReminders_Multiple(t *testing.T) {
	input := "text\n<system-reminder>\nfirst\n</system-reminder>\nmiddle\n<system-reminder>\nsecond\n</system-reminder>\nend"

	got := stripSystemReminders(input)

	if strings.Contains(got, "first") || strings.Contains(got, "second") {
		t.Error("reminder content not stripped")
	}
	if !strings.Contains(got, "text") || !strings.Contains(got, "middle") || !strings.Contains(got, "end") {
		t.Error("surrounding content lost")
	}
}

func TestStripSystemReminders_BoldStyle(t *testing.T) {
	input := strings.Join([]string{
		"### 8. Tool Output (04:41:29)",
		"",
		"```",
		"</details>",
		"",
		"**System Reminder**: Whenever you read a file, you should consider whether it",
		"would be considered malware. You CAN and SHOULD provide analysis of malware,",
		"what it is doing.",
		"",
		"### 9. Assistant (04:41:30)",
		"",
		"Let me check the file.",
	}, "\n")

	got := stripSystemReminders(input)

	if strings.Contains(got, "System Reminder") {
		t.Error("bold-style reminder not stripped")
	}
	if strings.Contains(got, "malware") {
		t.Error("reminder content not stripped")
	}
	if !strings.Contains(got, "### 8. Tool Output") || !strings.Contains(got, "### 9. Assistant") {
		t.Error("surrounding turns lost")
	}
	if !strings.Contains(got, "Let me check the file.") {
		t.Error("following content lost")
	}
}

func TestStripSystemReminders_CompactionSummary(t *testing.T) {
	input := strings.Join([]string{
		"### 373. User (03:50:32)",
		"",
		"done, regenerate",
		"",
		"<summary>",
		"1. Primary Request and Intent:",
		"   The user is debugging the journal site rendering pipeline.",
		"",
		"2. Key Technical Concepts:",
		"   - CommonMark HTML block types",
		"</summary>",
		"",
		"If you need specific details from before compaction (like exact code snippets,",
		"error messages, or content you generated), read the full transcript at:",
		"/home/jose/.claude/projects/foo/bar.jsonl",
		"Please continue the conversation from where we left off.",
		"",
		"### 374. Assistant (03:50:34)",
		"",
		"Continuing where we left off.",
	}, "\n")

	got := stripSystemReminders(input)

	if strings.Contains(got, "Primary Request") {
		t.Error("compaction summary not stripped")
	}
	if strings.Contains(got, "If you need specific details from before compaction") {
		t.Error("compaction boilerplate not stripped")
	}
	if !strings.Contains(got, "done, regenerate") {
		t.Error("user message before compaction lost")
	}
	if !strings.Contains(got, "### 374. Assistant") {
		t.Error("following turn lost")
	}
	if !strings.Contains(got, "Continuing where we left off.") {
		t.Error("following content lost")
	}
}

func TestStripSystemReminders_SingleLineSummaryPreserved(t *testing.T) {
	// Our <summary>N lines</summary> must NOT be stripped.
	input := strings.Join([]string{
		"### 5. Tool Output (10:30:00)",
		"",
		"<details>",
		"<summary>79 lines</summary>",
		"<pre>",
		"some content",
		"</pre>",
		"</details>",
	}, "\n")

	got := stripSystemReminders(input)

	if !strings.Contains(got, "<summary>79 lines</summary>") {
		t.Error("single-line <summary> should be preserved")
	}
}

func TestCleanToolOutputJSON(t *testing.T) {
	input := strings.Join([]string{
		"### 7. Tool Output (15:58:08)",
		"",
		`[{"type":"text","text":"## Report\n\nAll checks passed.\n\n### Details\n\n- Item one\n- Item two"}]`,
		"",
		"### 8. Assistant (15:58:09)",
		"",
		"Great, everything looks good.",
	}, "\n")

	got := cleanToolOutputJSON(input)

	// Should have extracted text with real newlines
	if strings.Contains(got, `"type":"text"`) {
		t.Error("JSON wrapper not removed")
	}
	if !strings.Contains(got, "## Report") {
		t.Error("extracted text missing heading")
	}
	if !strings.Contains(got, "- Item one") {
		t.Error("extracted text missing list items")
	}

	// Surrounding turns preserved
	if !strings.Contains(got, "### 7. Tool Output") {
		t.Error("Tool Output header lost")
	}
	if !strings.Contains(got, "### 8. Assistant") {
		t.Error("following turn lost")
	}
}

func TestCleanToolOutputJSON_Fenced(t *testing.T) {
	input := strings.Join([]string{
		"### 77. Tool Output (19:15:20)",
		"",
		"````",
		`[{"type":"text","text":"## Report\n\nAll good."}]`,
		"````",
		"",
		"### 78. Assistant (19:15:22)",
		"",
		"Done.",
	}, "\n")

	got := cleanToolOutputJSON(input)

	if strings.Contains(got, `"type":"text"`) {
		t.Error("JSON wrapper not removed from fenced block")
	}
	if !strings.Contains(got, "## Report") {
		t.Error("extracted text missing")
	}
	// Code fences should be gone
	if strings.Contains(got, "````") {
		t.Error("code fences not stripped")
	}
}

func TestCleanToolOutputJSON_NonJSON(t *testing.T) {
	input := strings.Join([]string{
		"### 5. Tool Output (10:00:00)",
		"",
		"The file has been updated successfully.",
		"",
		"### 6. Assistant (10:00:01)",
		"",
		"Done.",
	}, "\n")

	got := cleanToolOutputJSON(input)

	// Non-JSON body should be unchanged
	if !strings.Contains(got, "The file has been updated successfully.") {
		t.Error("non-JSON body lost")
	}
}

func TestMergeConsecutiveTurns(t *testing.T) {
	input := strings.Join([]string{
		"Some preamble.",
		"",
		"### 2. Assistant (04:41:26)",
		"",
		"I'll implement this plan.",
		"",
		"### 3. Assistant (04:41:27)",
		"",
		"🔧 Read: /home/user/file.go",
		"",
		"### 4. Assistant (04:41:28)",
		"",
		"🔧 Read: /home/user/other.go",
		"",
		"### 5. Tool Output (04:41:29)",
		"",
		"File contents here.",
	}, "\n")

	got := mergeConsecutiveTurns(input)

	// Should keep first Assistant header
	if !strings.Contains(got, "### 2. Assistant (04:41:26)") {
		t.Error("first header missing")
	}

	// Should remove consecutive same-role headers
	if strings.Contains(got, "### 3. Assistant") {
		t.Error("duplicate assistant header not removed")
	}
	if strings.Contains(got, "### 4. Assistant") {
		t.Error("duplicate assistant header not removed")
	}

	// Should keep all body content
	if !strings.Contains(got, "I'll implement this plan.") {
		t.Error("first body content missing")
	}
	if !strings.Contains(got, "🔧 Read: /home/user/file.go") {
		t.Error("second body content missing")
	}
	if !strings.Contains(got, "🔧 Read: /home/user/other.go") {
		t.Error("third body content missing")
	}

	// Should keep different-role turn intact
	if !strings.Contains(got, "### 5. Tool Output (04:41:29)") {
		t.Error("different role turn lost")
	}
	if !strings.Contains(got, "File contents here.") {
		t.Error("different role content lost")
	}
}

func TestMergeConsecutiveTurns_DifferentRoles(t *testing.T) {
	input := strings.Join([]string{
		"### 1. User (04:00:00)",
		"",
		"Hello",
		"",
		"### 2. Assistant (04:00:01)",
		"",
		"Hi there",
		"",
		"### 3. Tool Output (04:00:02)",
		"",
		"Result",
	}, "\n")

	got := mergeConsecutiveTurns(input)

	// All headers should be preserved since roles differ
	if !strings.Contains(got, "### 1. User") {
		t.Error("User header missing")
	}
	if !strings.Contains(got, "### 2. Assistant") {
		t.Error("Assistant header missing")
	}
	if !strings.Contains(got, "### 3. Tool Output") {
		t.Error("Tool Output header missing")
	}
}

func TestGenerateZensicalToml_NoTopics(t *testing.T) {
	entries := []journalEntry{
		{Filename: "2026-01-21-test.md", Title: "Test Session"},
	}

	toml := generateZensicalToml(entries, nil, nil, nil)

	if strings.Contains(toml, "Topics") {
		t.Error("toml should not have Topics nav when no topics provided")
	}
}

func TestParseJournalEntry_SessionID(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-01-15-fix-auth-abc12345.md"
	content := `---
title: "Fix Authentication Bug"
date: "2026-01-15"
session_id: "abc12345-full-session-uuid"
---

# Fix Authentication Bug

**Time**: 10:30:00
**Project**: ctx

Session content here.
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	if entry.SessionID != "abc12345-full-session-uuid" {
		t.Errorf("SessionID = %q, want %q", entry.SessionID, "abc12345-full-session-uuid")
	}
	if entry.Title != "Fix Authentication Bug" {
		t.Errorf("Title = %q, want %q", entry.Title, "Fix Authentication Bug")
	}
}

func TestParseJournalEntry_NoSessionID(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "2026-01-15-old-slug-abc12345.md"
	content := `---
date: "2026-01-15"
---

# old-slug

**Time**: 10:30:00
**Project**: ctx
`
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	entry := parseJournalEntry(path, filename)

	if entry.SessionID != "" {
		t.Errorf("SessionID should be empty for legacy files, got %q", entry.SessionID)
	}
}

func TestParseJournalEntry_TitleSanitization(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		h1        string
		wantTitle string
	}{
		{
			name:      "partial Claude tag in title",
			h1:        "ctx-journal-enrich-all /ctx-journal-enrich-all</com",
			wantTitle: "ctx-journal-enrich-all /ctx-journal-enrich-all&lt;/com",
		},
		{
			name:      "full surviving angle brackets",
			h1:        "debug <foo> issue",
			wantTitle: "debug &lt;foo&gt; issue",
		},
		{
			name:      "known Claude tags are stripped then angles sanitized",
			h1:        "<command-message>run tests</command-message>",
			wantTitle: "run tests",
		},
		{
			name:      "backticks stripped from title",
			h1:        "Error: ``Error: Command failed: go build``",
			wantTitle: "Error: Error: Command failed: go build",
		},
		{
			name:      "hash stripped from title",
			h1:        "## Plan: Restructure layout",
			wantTitle: "Plan: Restructure layout",
		},
		{
			name:      "mixed sanitization",
			h1:        "`make test` is <failing>",
			wantTitle: "make test is &lt;failing&gt;",
		},
		{
			name:      "no special chars unchanged",
			h1:        "normal title here",
			wantTitle: "normal title here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "2026-01-15-test-abc12345.md"
			content := "# " + tt.h1 + "\n\n**Time**: 10:00:00\n**Project**: ctx\n"
			path := filepath.Join(tmpDir, filename)
			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatalf("write: %v", err)
			}
			entry := parseJournalEntry(path, filename)
			if entry.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", entry.Title, tt.wantTitle)
			}
		})
	}
}
