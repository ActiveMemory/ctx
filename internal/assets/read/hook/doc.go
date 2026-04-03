//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook provides access to hook message templates and the
// hook registry from embedded assets.
//
// [Message] reads a specific template file by hook name and
// filename. [MessageRegistry] returns the raw registry.yaml.
// [TraceScript] reads an embedded trace git hook script.
package hook
