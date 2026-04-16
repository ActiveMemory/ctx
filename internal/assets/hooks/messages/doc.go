//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package messages provides metadata for hook message
// templates.
//
// The embedded registry.yaml maps each hook+variant
// pair to a category and description. The registry
// is parsed once via sync.Once and cached for the
// lifetime of the process.
//
// # Registry Access
//
// Registry returns the full list of HookMessageInfo
// entries parsed from the embedded YAML file.
//
//	for _, info := range messages.Registry() {
//	    fmt.Println(info.Hook, info.Variant)
//	}
//
// # Lookup
//
// Lookup finds a specific entry by hook name and
// variant. Returns nil if no match is found.
//
//	info := messages.Lookup("qa-reminder", "gate")
//
// # Variants
//
// Variants returns all variant names for a given
// hook, enabling enumeration of available templates.
//
//	names := messages.Variants("qa-reminder")
//	// => ["gate", "nudge"]
//
// # Types
//
// HookMessageInfo describes a single hook message
// template with hook name, variant, category
// (customizable or ctx-specific), description, and
// optional template variable names.
package messages
