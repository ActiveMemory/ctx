//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package data

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgWhy "github.com/ActiveMemory/ctx/internal/config/why"
)

// DocAliases maps user-facing names to embedded asset names.
var DocAliases = map[string]string{
	cfgWhy.DocManifesto:  cfgWhy.DocAliasManifesto,
	cfgWhy.DocAbout:      cfgWhy.DocAliasAbout,
	cfgWhy.DocInvariants: cfgWhy.DocAliasInvariants,
}

// DocOrder defines the display order for the interactive
// menu.
var DocOrder = []DocEntry{
	{
		cfgWhy.DocManifesto,
		desc.Text(text.DescKeyWriteWhyLabelManifesto),
	},
	{
		cfgWhy.DocAbout,
		desc.Text(text.DescKeyWriteWhyLabelAbout),
	},
	{
		cfgWhy.DocInvariants,
		desc.Text(text.DescKeyWriteWhyLabelInvariants),
	},
}
