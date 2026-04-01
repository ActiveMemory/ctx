//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package why provides constants for the ctx why command's embedded
// philosophy documents.
//
// Doc* constants are user-facing alias keys (CLI args, menu items).
// DocAlias* constants are embedded asset names (file stems in assets/why/).
// Import as config/why.
package why

// User-facing document alias keys (CLI arguments and menu items).
const (
	// DocManifesto is the user-facing alias for the manifesto.
	DocManifesto = "manifesto"
	// DocAbout is the user-facing alias for the about page.
	DocAbout = "about"
	// DocInvariants is the user-facing alias for design invariants.
	DocInvariants = "invariants"
)

// Embedded asset names (file stems in internal/assets/why/).
const (
	// DocAliasManifesto is the asset name for the manifesto document.
	DocAliasManifesto = "manifesto"
	// DocAliasAbout is the asset name for the about document.
	DocAliasAbout = "about"
	// DocAliasInvariants is the asset name for the design invariants document.
	DocAliasInvariants = "design-invariants"
)
