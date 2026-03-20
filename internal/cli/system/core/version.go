//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ParseMajorMinor extracts major and minor version numbers from a semver
// string like "1.2.3". Returns ok=false for unparseable versions.
//
// Parameters:
//   - ver: version string in semver format
//
// Returns:
//   - major: major version number
//   - minor: minor version number
//   - ok: true if parsing succeeded
func ParseMajorMinor(ver string) (major, minor int, ok bool) {
	parts := strings.SplitN(ver, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	var m, n int
	if _, scanErr := fmt.Sscanf(parts[0], "%d", &m); scanErr != nil {
		return 0, 0, false
	}
	if _, scanErr := fmt.Sscanf(parts[1], "%d", &n); scanErr != nil {
		return 0, 0, false
	}
	return m, n, true
}

// CheckKeyAge emits a nudge when the encryption key is older than the
// configured rotation threshold.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: current session identifier
func CheckKeyAge(cmd *cobra.Command, sessionID string) {
	kp := rc.KeyPath()
	info, statErr := os.Stat(kp)
	if statErr != nil {
		return // no key: nothing to check
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	threshold := rc.KeyRotationDays()

	if ageDays < threshold {
		return
	}

	keyFallback := fmt.Sprintf(
		desc.TextDesc(text.DescKeyCheckVersionKeyFallback), ageDays,
	)
	keyContent := LoadMessage(hook.CheckVersion, hook.VariantKeyRotation,
		map[string]any{tpl.VarKeyAgeDays: ageDays}, keyFallback)
	if keyContent == "" {
		return
	}

	boxTitle := desc.TextDesc(text.DescKeyCheckVersionKeyBoxTitle)
	relayPrefix := desc.TextDesc(text.DescKeyCheckVersionKeyRelayPrefix)

	cmd.Println(token.NewlineLF + NudgeBox(relayPrefix, boxTitle, keyContent))

	keyRef := notify.NewTemplateRef(hook.CheckVersion, hook.VariantKeyRotation,
		map[string]any{tpl.VarKeyAgeDays: ageDays})
	keyNotifyMsg := fmt.Sprintf(
		desc.TextDesc(text.DescKeyRelayPrefixFormat),
		hook.CheckVersion,
		fmt.Sprintf(
			desc.TextDesc(text.DescKeyCheckVersionKeyRelayFormat), ageDays,
		),
	)
	NudgeAndRelay(keyNotifyMsg, sessionID, keyRef)
}
