//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// MidSudo matches mid-command sudo after && || ;
var MidSudo = regexp.MustCompile(`(;|&&|\|\|)\s*sudo\s`)

// MidGitPush matches mid-command git push after && || ;
var MidGitPush = regexp.MustCompile(`(;|&&|\|\|)\s*git\s+push`)

// CpMvToBin matches cp/mv to bin directories.
var CpMvToBin = regexp.MustCompile(`(cp|mv)\s+\S+\s+(/usr/local/bin|/usr/bin|~/go/bin|~/.local/bin|/home/\S+/go/bin|/home/\S+/.local/bin)`)

// InstallToLocalBin matches cp/install to ~/.local/bin.
var InstallToLocalBin = regexp.MustCompile(`(cp|install)\s.*~/\.local/bin`)

// GitCommit matches git commit commands.
var GitCommit = regexp.MustCompile(`git\s+commit`)

// GitAmend matches the --amend flag.
var GitAmend = regexp.MustCompile(`--amend`)

// TaskRef matches Phase-style task references like HA.1, P-2.5, PD.3, CT.1.
var TaskRef = regexp.MustCompile(`\b[A-Z]+-?\d+\.?\d*\b`)
