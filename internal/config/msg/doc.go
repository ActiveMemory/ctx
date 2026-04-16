//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package msg defines formatting constants for the
// hook message table displayed by "ctx message list".
//
// # What This Package Controls
//
// When a user runs "ctx message list", ctx prints a
// columnar table showing every registered hook message
// with its hook name, variant, category, and override
// status. This package provides the column widths and
// separator widths that keep the table aligned across
// terminal sizes.
//
// # Column Layout
//
// Four columns are defined with matching separator
// widths for the underline row:
//
//   - Hook (MessageColHook=24, sep=22)
//   - Variant (MessageColVariant=20, sep=18)
//   - Category (MessageColCategory=16, sep=14)
//   - Override (sep=8)
//
// MessageListFormat is a pre-computed printf format
// string assembled at init time from the column
// widths. Consumers call fmt.Sprintf with this
// format to produce each table row.
//
// # Why Centralize
//
// Table formatting constants are easy to scatter
// across command handlers. Keeping them here ensures
// that column widths stay consistent if the table
// gains new columns and that the separator widths
// always match.
//
// # Public Surface
//
//   - MessageListFormat: pre-computed printf format
//     for table rows.
//   - MessageColHook, MessageColVariant,
//     MessageColCategory: column widths.
//   - MessageSepHook, MessageSepVariant,
//     MessageSepCategory, MessageSepOverride --
//     separator widths for underline row.
package msg
