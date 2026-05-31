//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unknown

import (
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
)

// SystemConfig is the unknown-subcommand relay for `ctx system`. Its
// copy frames the failure as a plugin/binary version skew, since
// `ctx system` is the hooks.json-wired group.
var SystemConfig = Config{
	RelayPrefixKey:  text.DescKeySystemUnknownRelayPrefix,
	BoxTitleKey:     text.DescKeySystemUnknownBoxTitle,
	BodyKey:         text.DescKeySystemUnknownBody,
	RelayMessageKey: text.DescKeySystemUnknownRelayMessage,
	HookName:        hook.System,
	Variant:         hook.VariantUnknownSubcommand,
}

// HookConfig is the unknown-subcommand relay for `ctx hook`. Its copy
// frames the failure as CLI drift between a caller (skill, loop script,
// hook) and the on-PATH binary, since `ctx hook` is consumed by name,
// not wired into hooks.json.
var HookConfig = Config{
	RelayPrefixKey:  text.DescKeyHookUnknownRelayPrefix,
	BoxTitleKey:     text.DescKeyHookUnknownBoxTitle,
	BodyKey:         text.DescKeyHookUnknownBody,
	RelayMessageKey: text.DescKeyHookUnknownRelayMessage,
	HookName:        hook.Hook,
	Variant:         hook.VariantUnknownSubcommand,
}
