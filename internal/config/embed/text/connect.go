//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for connect write output.
const (
	// DescKeyWriteConnectRegistered is the label for
	// successful registration.
	DescKeyWriteConnectRegistered = "write.connect-registered"
	// DescKeyWriteConnectSubscribed is the label for
	// successful subscription.
	DescKeyWriteConnectSubscribed = "write.connect-subscribed"
	// DescKeyWriteConnectSynced is the format string for
	// sync entry count.
	DescKeyWriteConnectSynced = "write.connect-synced"
	// DescKeyWriteConnectPublished is the format string for
	// publish entry count.
	DescKeyWriteConnectPublished = "write.connect-published"
	// DescKeyWriteConnectListening is the message shown when
	// entering listen mode.
	DescKeyWriteConnectListening = "write.connect-listening"
	// DescKeyWriteConnectReceived is the label for received
	// entries.
	DescKeyWriteConnectReceived = "write.connect-received"
	// DescKeyWriteConnectPublishWarning is the format string
	// for publish failure warnings.
	DescKeyWriteConnectPublishWarning = "write.connect-publish-warning"
	// DescKeyWriteConnectHubLabel is the label for hub status
	// output.
	DescKeyWriteConnectHubLabel = "write.connect-hub-label"
	// DescKeyWriteConnectHubStats is the format string for
	// hub entry and client counts.
	DescKeyWriteConnectHubStats = "write.connect-hub-stats"
	// DescKeyWriteConnectHubSync is the format string for
	// hub sync status messages.
	DescKeyWriteConnectHubSync = "write.connect-hub-sync"
)

// DescKeys for agent section headings.
const (
	// DescKeyWriteAgentSectionHub is the markdown heading
	// for the hub section in agent output.
	DescKeyWriteAgentSectionHub = "write.agent-section-hub"
)
