//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for journal output.
const (
	// DescKeyJournalConsolidateCount is the text key for journal consolidate
	// count messages.
	DescKeyJournalConsolidateCount = "journal.consolidate-count"
	// DescKeyJournalProjectLabel is the text key for journal project label
	// messages.
	DescKeyJournalProjectLabel = "journal.project-label"
	// DescKeyJournalMocBrowseBy is the text key for journal moc browse by
	// messages.
	DescKeyJournalMocBrowseBy = "journal.moc.browse-by"
	// DescKeyJournalMocFilePageStats is the text key for journal moc file page
	// stats messages.
	DescKeyJournalMocFilePageStats = "journal.moc.file-page-stats"
	// DescKeyJournalMocFileStats is the text key for journal moc file stats
	// messages.
	DescKeyJournalMocFileStats = "journal.moc.file-stats"
	// DescKeyJournalMocFilesDesc is the text key for journal moc files desc
	// messages.
	DescKeyJournalMocFilesDesc = "journal.moc.files-description"
	// DescKeyJournalMocNavDescription is the text key for journal moc nav
	// description messages.
	DescKeyJournalMocNavDescription = "journal.moc.nav-description"
	// DescKeyJournalMocSessionLink is the text key for journal moc session link
	// messages.
	DescKeyJournalMocSessionLink = "journal.moc.session-link"
	// DescKeyJournalMocTopicPageStats is the text key for journal moc topic page
	// stats messages.
	DescKeyJournalMocTopicPageStats = "journal.moc.topic-page-stats"
	// DescKeyJournalMocTopicStats is the text key for journal moc topic stats
	// messages.
	DescKeyJournalMocTopicStats = "journal.moc.topic-stats"
	// DescKeyJournalMocTopicsDesc is the text key for journal moc topics desc
	// messages.
	DescKeyJournalMocTopicsDesc = "journal.moc.topics-description"
	// DescKeyJournalMocTopicsLabel is the text key for journal moc topics label
	// messages.
	DescKeyJournalMocTopicsLabel = "journal.moc.topics-label"
	// DescKeyJournalMocTypeLabel is the text key for journal moc type label
	// messages.
	DescKeyJournalMocTypeLabel = "journal.moc.type-label"
	// DescKeyJournalMocTypePageStats is the text key for journal moc type page
	// stats messages.
	DescKeyJournalMocTypePageStats = "journal.moc.type-page-stats"
	// DescKeyJournalMocTypeStats is the text key for journal moc type stats
	// messages.
	DescKeyJournalMocTypeStats = "journal.moc.type-stats"
	// DescKeyJournalMocTypesDesc is the text key for journal moc types desc
	// messages.
	DescKeyJournalMocTypesDesc = "journal.moc.types-description"
	// DescKeyJournalMocBrowseItem is the text key for journal moc browse item
	// messages.
	DescKeyJournalMocBrowseItem = "journal.moc.browse-item"
	// DescKeyJournalMocHeadingTopics is the text key for journal moc heading
	// topics messages.
	DescKeyJournalMocHeadingTopics = "journal.moc.heading-topics"
	// DescKeyJournalMocHeadingPopular is the text key for journal moc heading
	// popular messages.
	DescKeyJournalMocHeadingPopular = "journal.moc.heading-popular"
	// DescKeyJournalMocHeadingLongtail is the text key for journal moc heading
	// longtail messages.
	DescKeyJournalMocHeadingLongtail = "journal.moc.heading-longtail"
	// DescKeyJournalMocHeadingFiles is the text key for journal moc heading files
	// messages.
	DescKeyJournalMocHeadingFiles = "journal.moc.heading-files"
	// DescKeyJournalMocHeadingFreq is the text key for journal moc heading freq
	// messages.
	DescKeyJournalMocHeadingFreq = "journal.moc.heading-frequent"
	// DescKeyJournalMocHeadingSingle is the text key for journal moc heading
	// single messages.
	DescKeyJournalMocHeadingSingle = "journal.moc.heading-single"
	// DescKeyJournalMocHeadingTypes is the text key for journal moc heading types
	// messages.
	DescKeyJournalMocHeadingTypes = "journal.moc.heading-types"
	// DescKeyJournalMocHeadingMonth is the text key for journal moc heading month
	// messages.
	DescKeyJournalMocHeadingMonth = "journal.moc.heading-month"
	// DescKeyJournalMocItemSessions is the text key for journal moc item sessions
	// messages.
	DescKeyJournalMocItemSessions = "journal.moc.item-sessions"
	// DescKeyJournalMocItemNamed is the text key for journal moc item named
	// messages.
	DescKeyJournalMocItemNamed = "journal.moc.item-named"
	// DescKeyJournalMocItemFileSess is the text key for journal moc item file
	// sess messages.
	DescKeyJournalMocItemFileSess = "journal.moc.item-file-sessions"
	// DescKeyJournalMocItemFileNamed is the text key for journal moc item file
	// named messages.
	DescKeyJournalMocItemFileNamed = "journal.moc.item-file-named"
	// DescKeyJournalMocItemListed is the text key for journal moc item listed
	// messages.
	DescKeyJournalMocItemListed = "journal.moc.item-listed"
	// DescKeyJournalMocPageTitle is the text key for journal moc page title
	// messages.
	DescKeyJournalMocPageTitle = "journal.moc.page-title"
	// DescKeyJournalMocCodeTitle is the text key for journal moc code title
	// messages.
	DescKeyJournalMocCodeTitle = "journal.moc.code-title"
	// DescKeyJournalMocTopicsMocLink is the text key for journal moc topics moc
	// link messages.
	DescKeyJournalMocTopicsMocLink = "journal.moc.topics-moc-link"
	// DescKeyJournalMocTopicSep is the text key for journal moc topic sep
	// messages.
	DescKeyJournalMocTopicSep = "journal.moc.topic-separator"
)

// DescKeys for journal write output.
const (
	// DescKeyWriteJournalOrphanRemoved is the text key for write journal orphan
	// removed messages.
	DescKeyWriteJournalOrphanRemoved = "write.journal-orphan-removed"
	// DescKeyWriteJournalSiteBuilding is the text key for write journal site
	// building messages.
	DescKeyWriteJournalSiteBuilding = "write.journal-site-building"
	// DescKeyWriteJournalSiteGeneratedBlock is the text key for write journal
	// site generated block messages.
	DescKeyWriteJournalSiteGeneratedBlock = "write.journal-site-generated-block"
	// DescKeyWriteJournalSiteStarting is the text key for write journal site
	// starting messages.
	DescKeyWriteJournalSiteStarting = "write.journal-site-starting"
	// DescKeyWriteJournalSyncLocked is the text key for write journal sync locked
	// messages.
	DescKeyWriteJournalSyncLocked = "write.journal-sync-locked"
	// DescKeyWriteJournalSyncLockedCount is the text key for write journal sync
	// locked count messages.
	DescKeyWriteJournalSyncLockedCount = "write.journal-sync-locked-count"
	// DescKeyWriteJournalSyncMatch is the text key for write journal sync match
	// messages.
	DescKeyWriteJournalSyncMatch = "write.journal-sync-match"
	// DescKeyWriteJournalSyncNone is the text key for write journal sync none
	// messages.
	DescKeyWriteJournalSyncNone = "write.journal-sync-none"
	// DescKeyWriteJournalSyncUnlocked is the text key for write journal sync
	// unlocked messages.
	DescKeyWriteJournalSyncUnlocked = "write.journal-sync-unlocked"
	// DescKeyWriteJournalSyncUnlockedCount is the text key for write journal sync
	// unlocked count messages.
	DescKeyWriteJournalSyncUnlockedCount = "write.journal-sync-unlocked-count"
)
