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
		{cfgsync.PatternEslint, TextDesc(text.DescKeySyncTopicEslint)},
		{cfgsync.PatternPrettier, TextDesc(text.DescKeySyncTopicPrettier)},
		{cfgsync.PatternTSConfig, TextDesc(text.DescKeySyncTopicTSConfig)},
		{cfgsync.PatternEditorConf, TextDesc(text.DescKeySyncTopicEditorConfig)},
		{cfgsync.PatternMakefile, TextDesc(text.DescKeySyncTopicMakefile)},
		{cfgsync.PatternDockerfile, TextDesc(text.DescKeySyncTopicDockerfile)},
	}
}
