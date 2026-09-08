//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for hub write confirmations.
const (
	// DescKeyWriteHubAddedPeer is the text key for hub added
	// peer messages.
	DescKeyWriteHubAddedPeer = "write.hub-added-peer"
	// DescKeyWriteHubRemovedPeer is the text key for hub
	// removed peer messages.
	DescKeyWriteHubRemovedPeer = "write.hub-removed-peer"
	// DescKeyWriteHubLeadershipTransferred is the text key for
	// hub leadership transferred messages.
	DescKeyWriteHubLeadershipTransferred = "write.hub-leadership-transferred"
	// DescKeyWriteHubLeader is the text key for hub leader
	// messages.
	DescKeyWriteHubLeader = "write.hub-leader"
	// DescKeyWriteHubRole is the text key for hub role
	// messages.
	DescKeyWriteHubRole = "write.hub-role"
	// DescKeyWriteHubLeaderUnknown is the text key for the
	// leader line while no leader is known: an election is in
	// progress, or quorum was lost.
	DescKeyWriteHubLeaderUnknown = "write.hub-leader-unknown"
	// DescKeyWriteHubClusterStats is the text key for hub
	// cluster statistics.
	DescKeyWriteHubClusterStats = "write.hub-cluster-stats"
	// DescKeyWriteHubEntries is the text key for the entry
	// count on a standalone hub, where the cluster-stats line
	// would report a peer count that does not exist.
	DescKeyWriteHubEntries = "write.hub-entries"
	// DescKeyWriteHubDroppedListeners is the text key for the
	// cumulative slow-listener disconnect count. Printed only
	// when the count is non-zero.
	DescKeyWriteHubDroppedListeners = "write.hub-dropped-listeners"
	// DescKeyWriteHubRevoked is the text key for the hub client
	// revocation confirmation.
	DescKeyWriteHubRevoked = "write.hub-revoked"
)
