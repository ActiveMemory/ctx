//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

// Publish budget and limits.
const (
	// DefaultPublishBudget is the default line budget for published content.
	DefaultPublishBudget = 80
	// PublishMaxTasks is the maximum number of pending tasks to publish.
	PublishMaxTasks = 10
	// PublishMaxDecisions is the maximum number of recent decisions to publish.
	PublishMaxDecisions = 5
	// PublishMaxConventions is the maximum number of convention items to publish.
	PublishMaxConventions = 10
	// PublishMaxLearnings is the maximum number of recent learnings to publish.
	PublishMaxLearnings = 5
	// PublishRecentDays is the lookback window in days for recent entries.
	PublishRecentDays = 7
)
