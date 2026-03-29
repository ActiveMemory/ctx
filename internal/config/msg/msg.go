//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package msg

import "fmt"

// MessageListFormat is the pre-computed printf format for message list rows.
var MessageListFormat = fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%s",
	MessageColHook, MessageColVariant, MessageColCategory)

// Message table formatting.
const (
	// MessageColHook is the column width for the Hook field
	// in message list output.
	MessageColHook = 24
	// MessageColVariant is the column width for the Variant
	// field in message list output.
	MessageColVariant = 20
	// MessageColCategory is the column width for the Category
	// field in message list output.
	MessageColCategory = 16
	// MessageSepHook is the separator width for the Hook column underline.
	MessageSepHook = 22
	// MessageSepVariant is the separator width for the Variant column underline.
	MessageSepVariant = 18
	// MessageSepCategory is the separator width for the Category column underline.
	MessageSepCategory = 14
	// MessageSepOverride is the separator width for the Override column underline.
	MessageSepOverride = 8
)
