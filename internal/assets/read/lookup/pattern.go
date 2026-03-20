//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/sync"
)

// ConfigPattern pairs a glob pattern with its localizable topic description.
type ConfigPattern struct {
	Pattern string
	Topic   string
}

// ConfigPatterns returns config file patterns with resolved topic descriptions.
func ConfigPatterns() []ConfigPattern {
	return []ConfigPattern{
		{sync.PatternEslint, TextDesc(text.DescKeySyncTopicEslint)},
		{sync.PatternPrettier, TextDesc(text.DescKeySyncTopicPrettier)},
		{sync.PatternTSConfig, TextDesc(text.DescKeySyncTopicTSConfig)},
		{sync.PatternEditorConf, TextDesc(text.DescKeySyncTopicEditorConfig)},
		{sync.PatternMakefile, TextDesc(text.DescKeySyncTopicMakefile)},
		{sync.PatternDockerfile, TextDesc(text.DescKeySyncTopicDockerfile)},
	}
}
