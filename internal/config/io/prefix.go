//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

// DangerousPrefixes lists system directories where ctx
// should never read from or write to. Checked after
// filepath.Abs resolution.
var DangerousPrefixes = []string{
	// PrefixBin is the system binaries directory.
	"/bin/",
	// PrefixBoot is the boot loader directory.
	"/boot/",
	// PrefixDev is the device files directory.
	"/dev/",
	// PrefixEtc is the system configuration directory.
	"/etc/",
	// PrefixLib is the shared libraries directory.
	"/lib/",
	// PrefixLib64 is the 64-bit shared libraries directory.
	"/lib64/",
	// PrefixProc is the process information directory.
	"/proc/",
	// PrefixSbin is the system binaries directory.
	"/sbin/",
	// PrefixSys is the kernel/device tree directory.
	"/sys/",
	// PrefixUsrBin is the user binaries directory.
	"/usr/bin/",
	// PrefixUsrLib is the user libraries directory.
	"/usr/lib/",
	// PrefixUsrSbin is the user system binaries directory.
	"/usr/sbin/",
}
