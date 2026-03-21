//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rss

import "github.com/ActiveMemory/ctx/internal/config/token"

// Site feed defaults.
const (
	// DefaultFeedInputDir is the default blog source directory.
	DefaultFeedInputDir = "docs/blog"
	// DefaultFeedOutPath is the default output path for the Atom feed.
	DefaultFeedOutPath = "site/feed.xml"
	// DefaultFeedBaseURL is the default base URL for feed entry links.
	DefaultFeedBaseURL = "https://ctx.ist"
	// FeedAtomNS is the Atom XML namespace URI.
	FeedAtomNS = "http://www.w3.org/2005/Atom"
	// FeedTitle is the default feed title.
	FeedTitle = "ctx blog"
	// FeedDefaultAuthor is the default author for feed entries.
	FeedDefaultAuthor = "Jose Alekhinne"
	// FeedXMLHeader is the XML declaration prepended to feed output.
	FeedXMLHeader = `<?xml version="1.0" encoding="utf-8"?>` + token.NewlineLF
)

// Feed URL path constants.
const (
	FeedPath     = "/feed.xml"
	BlogPath     = "/blog/"
	LinkRelSelf  = "self"
	TimeSuffixZ  = "T00:00:00Z"
	SkipSentinel = " \u2014 "
	URLSlash     = "/"
)
