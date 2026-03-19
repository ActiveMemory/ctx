//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	cfgsync "github.com/ActiveMemory/ctx/internal/config/sync"

	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ConfigPattern pairs a glob pattern with its localizable topic description.
type ConfigPattern struct {
	Pattern string
	Topic   string
}

// ConfigPatterns returns config file patterns with resolved topic descriptions.
func ConfigPatterns() []ConfigPattern {
	return []ConfigPattern{
		{cfgsync.PatternEslint, TextDesc(text.TextDescKeySyncTopicEslint)},
		{cfgsync.PatternPrettier, TextDesc(text.TextDescKeySyncTopicPrettier)},
		{cfgsync.PatternTSConfig, TextDesc(text.TextDescKeySyncTopicTSConfig)},
		{cfgsync.PatternEditorConf, TextDesc(text.TextDescKeySyncTopicEditorConfig)},
		{cfgsync.PatternMakefile, TextDesc(text.TextDescKeySyncTopicMakefile)},
		{cfgsync.PatternDockerfile, TextDesc(text.TextDescKeySyncTopicDockerfile)},
	}
}
