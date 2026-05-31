//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unknown

// Config parameterizes the unknown-subcommand relay for one command
// group. Each group opts in by setting its RunE to [HandlerFor](cfg).
//
// Fields:
//   - RelayPrefixKey: text key for the "relay verbatim" prefix line
//   - BoxTitleKey: text key for the NudgeBox title
//   - BodyKey: text key for the box body (format string taking the verb)
//   - RelayMessageKey: text key for the event-log/webhook message
//     (format string taking the verb)
//   - HookName: relay-ref hook label (e.g. hook.System, hook.Hook)
//   - Variant: relay-ref variant (e.g. hook.VariantUnknownSubcommand)
type Config struct {
	RelayPrefixKey  string
	BoxTitleKey     string
	BodyKey         string
	RelayMessageKey string
	HookName        string
	Variant         string
}
