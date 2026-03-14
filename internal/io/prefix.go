//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

// dangerousPrefixes lists system directories where ctx should never
// read from or write to. Checked after filepath.Abs resolution.
var dangerousPrefixes = []string{
	"/bin/",
	"/boot/",
	"/dev/",
	"/etc/",
	"/lib/",
	"/lib64/",
	"/proc/",
	"/sbin/",
	"/sys/",
	"/usr/bin/",
	"/usr/lib/",
	"/usr/sbin/",
}
