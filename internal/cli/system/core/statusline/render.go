//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/env"
	cfgStatusline "github.com/ActiveMemory/ctx/internal/config/statusline"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// LocationSegment renders "user@host dir" from the process
// environment and the payload's working directory. Any part that
// cannot be determined is omitted; an empty result drops the whole
// segment.
//
// Parameters:
//   - p: parsed status line payload
//
// Returns:
//   - string: location segment, or "" when nothing is resolvable
func LocationSegment(p *Payload) string {
	userHost := UserAtHost()
	dir := ShortDir(WorkingDir(p))
	switch {
	case userHost != "" && dir != "":
		return userHost + token.Space + dir
	case userHost != "":
		return userHost
	default:
		return dir
	}
}

// UserAtHost renders the "user@host" prefix from the process
// environment.
//
// Returns:
//   - string: "user@host", a lone part when only one resolves, or ""
func UserAtHost() string {
	name := os.Getenv(env.User)
	if name == "" {
		if u, userErr := user.Current(); userErr == nil {
			name = u.Username
		}
	}
	host, hostErr := os.Hostname()
	if hostErr == nil {
		// Short hostname: the domain suffix is noise at status-line
		// widths.
		host, _, _ = strings.Cut(host, cfgStatusline.HostDomainSeparator)
	}
	name = Sanitize(name)
	host = Sanitize(host)
	switch {
	case name != "" && host != "":
		return name + cfgStatusline.UserHostSeparator + host
	case name != "":
		return name
	default:
		return host
	}
}

// WorkingDir picks the payload's directory, preferring
// workspace.current_dir over the legacy cwd duplicate.
//
// Parameters:
//   - p: parsed status line payload
//
// Returns:
//   - string: working directory, or "" when the payload has none
func WorkingDir(p *Payload) string {
	if p.Workspace.CurrentDir != "" {
		return p.Workspace.CurrentDir
	}
	return p.Cwd
}

// ShortDir home-abbreviates a path and collapses overlong results so
// deep trees do not crowd out the other segments.
//
// Parameters:
//   - path: absolute directory path from the payload
//
// Returns:
//   - string: display form of the path, or "" for empty input
func ShortDir(path string) string {
	if path == "" {
		return ""
	}
	if home, homeErr := os.UserHomeDir(); homeErr == nil && home != "" {
		if rel, relErr := filepath.Rel(home, path); relErr == nil &&
			!strings.HasPrefix(rel, cfgStatusline.RelParentPrefix) {
			if rel == cfgStatusline.RelSelf {
				path = cfgStatusline.HomeAbbrev
			} else {
				path = cfgStatusline.HomeAbbrevPrefix + rel
			}
		}
	}
	path = Sanitize(path)
	if len(path) > cfgStatusline.MaxDirLen {
		path = cfgStatusline.TruncatedDirPrefix + filepath.Base(path)
	}
	return path
}

// Sanitize reduces a payload-derived string to bounded printable
// ASCII. Control bytes, ANSI escapes, and multi-byte runes are
// stripped rather than replaced; the result is trimmed and capped at
// the segment bound.
//
// Parameters:
//   - s: raw string from the payload or environment
//
// Returns:
//   - string: printable-ASCII form, at most MaxSegmentLen bytes
func Sanitize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s) && b.Len() < cfgStatusline.MaxSegmentLen; i++ {
		c := s[i]
		if c >= cfgStatusline.ASCIIPrintableMin &&
			c <= cfgStatusline.ASCIIPrintableMax {
			b.WriteByte(c)
		}
	}
	return strings.TrimSpace(b.String())
}
