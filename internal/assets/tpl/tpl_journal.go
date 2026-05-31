//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Journal site format templates.
//
// These templates define the structure of generated journal site pages.
// Each uses fmt.Sprintf verbs for interpolation.
const (
	// JournalIndexIntro is the introductory line on the journal index.
	JournalIndexIntro = "Browse your AI session history."

	// JournalIndexStats formats the session/suggestion count line.
	// Args: regular count, suggestion count.
	JournalIndexStats = "**Sessions**: %d | **Suggestions**: %d"

	// JournalSuggestionsNote is the description under the suggestions heading.
	JournalSuggestionsNote = "*Auto-generated suggestion" +
		" prompts from Claude Code.*"

	// JournalMonthHeading formats a month section heading.
	// Args: month string (YYYY-MM).
	JournalMonthHeading = "## %s"

	// JournalIndexEntry formats a single entry in the journal index.
	// Args: timeStr, title, link, project, size.
	JournalIndexEntry = "- %s[%s](%s.md)%s `%s`"

	// JournalIndexSummary formats the summary line below an index entry.
	// Args: summary text. Indented to nest visually under the parent bullet.
	JournalIndexSummary = "    *%s*"

	// JournalSummaryAdmonition formats the summary as an abstract admonition
	// on individual journal entry pages.
	// Args: summary text.
	JournalSummaryAdmonition = "!!! abstract \"Summary\"\n    %s"

	// JournalSourceLink formats the "View source" link injected into entries.
	// Args: absPath, relPath, relPath.
	JournalSourceLink = `*[View source](file://%s) · <code>%s</code>*` +
		` <button onclick="navigator.clipboard.writeText('%s')" title="Copy path"` +
		` style="cursor:pointer;border:none;background:none;` +
		`font-size:0.8em;vertical-align:middle">` +
		`&#x2398;</button>`

	// JournalTopicStats formats the topics index summary line.
	// Args: topic count, session count, popular count, longtail count.
	JournalTopicStats = "**%d topics** across" +
		" **%d sessions** - **%d popular**," +
		" **%d long-tail**"

	// JournalTopicPageStats formats an individual topic page summary.
	// Args: session count.
	JournalTopicPageStats = "**%d sessions** with this topic."

	// JournalFileStats formats the key files index summary line.
	// Args: file count, session count, popular count, longtail count.
	JournalFileStats = "**%d files** across" +
		" **%d sessions** - **%d popular**," +
		" **%d long-tail**"

	// JournalFilePageStats formats an individual key file page summary.
	// Args: session count.
	JournalFilePageStats = "**%d sessions** touching this file."

	// JournalTypeStats formats the session types index summary line.
	// Args: type count, session count.
	JournalTypeStats = "**%d types** across **%d sessions**"

	// JournalTypePageStats formats an individual type page summary.
	// Args: session count, type name.
	JournalTypePageStats = "**%d sessions** of type *%s*."

	// JournalPageHeading formats a Markdown heading for an index page.
	// Args: name.
	JournalPageHeading = "# %s"

	// JournalCodePageHeading formats a Markdown heading with code styling.
	// Args: path.
	JournalCodePageHeading = "# `%s`"

	// JournalSubpageEntry formats a session link on a subpage (topic, file, type).
	// Args: timeStr, title, linkPrefix, link.
	JournalSubpageEntry = "- %s[%s](%s%s.md)"

	// JournalLongtailEntry formats an inline longtail topic entry.
	// Args: name, title, link.
	JournalLongtailEntry = "- **%s** - [%s](../%s.md)"

	// JournalLongtailCodeEntry formats an inline longtail key file entry.
	// Args: path, title, link.
	JournalLongtailCodeEntry = "- `%s` - [%s](../%s.md)"

	// JournalNavItem formats a navigation item in zensical.toml.
	// Args: label, path.
	JournalNavItem = `  { "%s" = "%s" },`

	// JournalNavSection formats a navigation section opening.
	// Args: label.
	JournalNavSection = `  { "%s" = [`

	// JournalNavSessionItem formats a session entry in navigation.
	// Args: title, filename.
	JournalNavSessionItem = `    { "%s" = "%s" },`

	// ZensicalExtraCSS is the extra_css line for zensical.toml.
	// Must appear under [project] (after nav, before [project.theme]).
	ZensicalExtraCSS = `extra_css = ["stylesheets/extra.css"]`
)
