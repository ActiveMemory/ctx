//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import "encoding/xml"

// AtomFeed represents an Atom 1.0 feed document.
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
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

// AtomAuthor represents an author element in an Atom feed.
type AtomAuthor struct {
	Name string `xml:"name"`
}

// AtomCategory represents a category element in an Atom feed.
type AtomCategory struct {
	Term string `xml:"term,attr"`
}
