//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

// Topic-template substitution tokens. Replaced by the topic
// scaffolder with the human-readable name and the kebab-case
// slug, respectively. Structural literals; not localizable.
const (
	// TokenTopicName is the long-form token for the topic name.
	TokenTopicName = "<TOPIC_NAME>"
	// TokenTopicSlug is the long-form token for the slug.
	//
	//nolint:gosec // G101: angle-bracket placeholder, not a credential
	TokenTopicSlug = "<TOPIC_SLUG>"
	// TokenName is the short-form token for the topic name.
	TokenName = "<NAME>"
	// TokenSlug is the short-form token for the slug.
	TokenSlug = "<SLUG>"
)

// ManagedBlock start/end markers. The reindex command rewrites
// the contents between these markers in the kb landing page.
// Structural literals parsed by the reindex regex; not
// localizable.
const (
	// ManagedKBTopicsStart opens the CTX:KB:TOPICS managed block.
	ManagedKBTopicsStart = "<!-- CTX:KB:TOPICS START -->"
	// ManagedKBTopicsEnd closes the CTX:KB:TOPICS managed block.
	ManagedKBTopicsEnd = "<!-- CTX:KB:TOPICS END -->"
	// ManagedKBTopicsEmpty is the placeholder line written into
	// the managed block when no topics exist yet.
	ManagedKBTopicsEmpty = "- _no topics yet; create one with " +
		"`ctx kb topic new \"<name>\"`_\n"
	// TopicEntryPrefix opens each topic list item.
	TopicEntryPrefix = "- [`"
	// TopicEntryMiddle separates the slug from the
	// link target in a topic list item.
	TopicEntryMiddle = "`](topics/"
	// TopicEntrySuffix closes each topic list item.
	TopicEntrySuffix = "/)\n"
)
