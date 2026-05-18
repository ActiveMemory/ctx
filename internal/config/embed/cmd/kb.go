//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Cobra Use strings for the kb command family.
const (
	// UseKB is the cobra use string for the kb parent command.
	UseKB = "kb"
	// UseKBAsk is the cobra use string for `ctx kb ask`.
	UseKBAsk = "ask \"<question>\""
	// UseKBGround is the cobra use string for `ctx kb ground`.
	UseKBGround = "ground"
	// UseKBIngest is the cobra use string for `ctx kb ingest`.
	UseKBIngest = "ingest <folder|paths...>"
	// UseKBNote is the cobra use string for `ctx kb note`.
	UseKBNote = "note \"<text>\""
	// UseKBReindex is the cobra use string for `ctx kb reindex`.
	UseKBReindex = "reindex"
	// UseKBSiteReview is the cobra use string for
	// `ctx kb site-review`.
	UseKBSiteReview = "site-review"
	// UseKBTopic is the cobra use string for the topic parent
	// command.
	UseKBTopic = "topic"
	// UseKBTopicNew is the cobra use string for
	// `ctx kb topic new`.
	UseKBTopicNew = "new <name>"
)

// DescKeys for the kb command family. Values map to entries in
// the commands.yaml asset.
const (
	// DescKeyKB is the description key for the kb parent.
	DescKeyKB = "kb"
	// DescKeyKBAsk is the description key for `ctx kb ask`.
	DescKeyKBAsk = "kb.ask"
	// DescKeyKBGround is the description key for `ctx kb ground`.
	DescKeyKBGround = "kb.ground"
	// DescKeyKBIngest is the description key for `ctx kb ingest`.
	DescKeyKBIngest = "kb.ingest"
	// DescKeyKBNote is the description key for `ctx kb note`.
	DescKeyKBNote = "kb.note"
	// DescKeyKBReindex is the description key for
	// `ctx kb reindex`.
	DescKeyKBReindex = "kb.reindex"
	// DescKeyKBSiteReview is the description key for
	// `ctx kb site-review`.
	DescKeyKBSiteReview = "kb.site-review"
	// DescKeyKBTopic is the description key for the topic
	// parent command.
	DescKeyKBTopic = "kb.topic"
	// DescKeyKBTopicNew is the description key for
	// `ctx kb topic new`.
	DescKeyKBTopicNew = "kb.topic.new"
)
