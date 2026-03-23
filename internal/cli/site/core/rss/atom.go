//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rss

import "encoding/xml"

// AtomFeed represents an Atom 1.0 feed document.
//
// Fields:
//   - XMLName: XML element name ("feed")
//   - NS: Atom namespace URI
//   - Title: Feed title
//   - Links: Feed links (self and alternate)
//   - ID: Unique feed identifier
//   - Updated: Last update timestamp
//   - Entries: Feed entries
type AtomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	NS      string      `xml:"xmlns,attr"`
	Title   string      `xml:"title"`
	Links   []AtomLink  `xml:"link"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Entries []AtomEntry `xml:"entry"`
}

// AtomEntry represents a single entry in an Atom feed.
//
// Fields:
//   - Title: Entry title
//   - Links: Entry links
//   - ID: Unique entry identifier
//   - Updated: Entry update timestamp
//   - Summary: Optional entry summary
//   - Author: Optional entry author
//   - Categories: Optional entry categories
type AtomEntry struct {
	Title      string         `xml:"title"`
	Links      []AtomLink     `xml:"link"`
	ID         string         `xml:"id"`
	Updated    string         `xml:"updated"`
	Summary    string         `xml:"summary,omitempty"`
	Author     *AtomAuthor    `xml:"author,omitempty"`
	Categories []AtomCategory `xml:"category,omitempty"`
}

// AtomLink represents a link element in an Atom feed.
//
// Fields:
//   - Href: Link URL
//   - Rel: Link relation type (e.g., "self", "alternate")
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

// AtomAuthor represents an author element in an Atom feed.
//
// Fields:
//   - Name: Author name
type AtomAuthor struct {
	Name string `xml:"name"`
}

// AtomCategory represents a category element in an Atom feed.
//
// Fields:
//   - Term: Category term
type AtomCategory struct {
	Term string `xml:"term,attr"`
}
