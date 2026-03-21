//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package zensical

// Zensical site configuration.
//
// Note that both the journal site and the main docs site share
// the same config settings.
const (
	// Toml is the zensical site configuration filename.
	Toml = "zensical.toml"
	// Bin is the zensical binary name.
	Bin = "zensical"
	// Stylesheets is the subdirectory for CSS stylesheets in site output.
	Stylesheets = "stylesheets"
	// ExtraCSS is the custom CSS filename for journal sites.
	ExtraCSS = "extra.css"
	// CmdServe is the zensical serve subcommand.
	CmdServe = "serve"
	// CmdBuild is the zensical build subcommand.
	CmdBuild = "build"
)
