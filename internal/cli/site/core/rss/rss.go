//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rss

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/site/core"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/rss"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errSite "github.com/ActiveMemory/ctx/internal/err/site"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Atom builds the Atom XML and writes it to outPath.
//
// Parameters:
//   - posts: Blog posts to include in the feed
//   - outPath: Output file path for the generated XML
//   - baseURL: Base URL for entry links
//
// Returns:
//   - error: Non-nil if marshalling or writing fails
func Atom(posts []core.BlogPost, outPath, baseURL string) error {
	baseURL = strings.TrimRight(baseURL, rss.URLSlash)

	feedURL := baseURL + rss.FeedPath
	blogURL := baseURL + rss.BlogPath

	updated := ""
	if len(posts) > 0 {
		updated = posts[0].Date + rss.TimeSuffixZ
	}

	feed := AtomFeed{
		NS:    rss.FeedAtomNS,
		Title: rss.FeedTitle,
		Links: []AtomLink{
			{Href: blogURL},
			{Href: feedURL, Rel: rss.LinkRelSelf},
		},
		ID:      feedURL,
		Updated: updated,
	}

	for _, p := range posts {
		slug := strings.TrimSuffix(p.Filename, file.ExtMarkdown)
		entryURL := blogURL + slug + rss.URLSlash

		entry := AtomEntry{
			Title:   p.Title,
			Links:   []AtomLink{{Href: entryURL}},
			ID:      entryURL,
			Updated: p.Date + rss.TimeSuffixZ,
		}

		if p.Summary != "" && !strings.Contains(p.Summary, rss.SkipSentinel) {
			entry.Summary = p.Summary
		}

		author := p.Author
		if author == "" {
			author = rss.FeedDefaultAuthor
		}
		entry.Author = &AtomAuthor{Name: author}

		for _, topic := range p.Topics {
			entry.Categories = append(entry.Categories, AtomCategory{Term: topic})
		}

		feed.Entries = append(feed.Entries, entry)
	}

	outDir := filepath.Dir(outPath)
	if mkErr := ctxIo.SafeMkdirAll(outDir, 0o755); mkErr != nil {
		return errFs.Mkdir(outDir, mkErr)
	}

	xmlData, marshalErr := xml.MarshalIndent(feed, "", "  ")
	if marshalErr != nil {
		return errSite.MarshalFeed(marshalErr)
	}

	output := []byte(rss.FeedXMLHeader)
	output = append(output, xmlData...)
	output = append(output, token.NewlineLF[0])

	if writeErr := ctxIo.SafeWriteFile(outPath, output, 0o644); writeErr != nil {
		return errFs.FileWrite(outPath, writeErr)
	}

	return nil
}
