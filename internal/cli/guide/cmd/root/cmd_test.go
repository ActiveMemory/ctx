//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

func TestParseSkillFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    skillMeta
		wantErr bool
	}{
		{
			name: "valid frontmatter",
			input: "---\nname: ctx-test\n" +
				"description: \"A test skill.\"\n---\nBody text.",
			want: skillMeta{Name: "ctx-test", Description: "A test skill."},
		},
		{
			name:  "missing frontmatter",
			input: "No frontmatter here.",
			want:  skillMeta{},
		},
		{
			name:  "unterminated frontmatter",
			input: "---\nname: ctx-test\nno closing delimiter",
			want:  skillMeta{},
		},
		{
			name:    "invalid YAML",
			input:   "---\n: :\n  bad:\n    - [\n---\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSkillFrontmatter([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSkillFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseSkillFrontmatter() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGuideLineCount(t *testing.T) {
	guide := desc.Text(text.DescKeyGuideDefault)
	lines := strings.Split(guide, "\n")
	if len(lines) > 50 {
		t.Errorf("guide default has %d lines, want at most 50", len(lines))
	}
}
