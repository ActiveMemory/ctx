//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// TestDisclosureInvariants runs the progressive-disclosure guards and
// invariants against this repo's own .context/ knowledge files. The
// tree is not yet migrated, so every check passes vacuously — this is
// the guard that no canonical file drifts into a shape the guards would
// reject before migration even begins.
//
// It is proven-both-ways below (TestDisclosureInvariants_PlantedViolation):
// a check that only ever passes proves nothing.
func TestDisclosureInvariants(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	ctxDir := filepath.Join(root, ".context")

	cases := []struct {
		file     string
		themeDir string
		kind     disclosure.Kind
	}{
		{cfgCtx.Learning, cfgDisc.ThemeDirLearning, disclosure.KindLearning},
		{cfgCtx.Decision, cfgDisc.ThemeDirDecision, disclosure.KindDecision},
		{cfgCtx.Convention, cfgDisc.ThemeDirConvention, disclosure.KindConvention},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			data, readErr := os.ReadFile(filepath.Join(ctxDir, tc.file))
			if readErr != nil {
				t.Fatalf("read %s: %v", tc.file, readErr)
			}
			r := disclosure.Parse(string(data), tc.kind)
			themeDir := filepath.Join(ctxDir, tc.themeDir)

			if vErr := disclosure.Validate(r); vErr != nil {
				t.Errorf("Validate(%s) = %v, want nil", tc.file, vErr)
			}
			if pErr := disclosure.CheckPairing(r, themeDir); pErr != nil {
				t.Errorf("CheckPairing(%s) = %v, want nil", tc.file, pErr)
			}
			if uErr := disclosure.CheckUniqueness(r, themeDir); uErr != nil {
				t.Errorf("CheckUniqueness(%s) = %v, want nil", tc.file, uErr)
			}
			if lErr := disclosure.CheckLinks(r, ctxDir); lErr != nil {
				t.Errorf("CheckLinks(%s) = %v, want nil", tc.file, lErr)
			}
		})
	}
}

// TestDisclosureInvariants_PlantedViolation proves the guard actually
// fires: a root with two ## Themes sections must be refused.
func TestDisclosureInvariants_PlantedViolation(t *testing.T) {
	planted := "# Learnings\n\n## Themes\n\n- a — g → [a](learnings/a.md)\n\n" +
		"## Themes\n\n- b — g → [b](learnings/b.md)\n"

	err := disclosure.Validate(disclosure.Parse(planted, disclosure.KindLearning))
	if !errors.Is(err, errDisc.ErrMultipleThemes) {
		t.Errorf("Validate(planted double-Themes) = %v, want ErrMultipleThemes", err)
	}
}
