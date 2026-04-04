//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Discover finds all hook scripts in the hooks directory, grouped
// by type. It iterates over each valid hook type subdirectory,
// validates each file via [ValidatePath], and skips invalid
// entries with a logged warning. Hooks within each type are sorted
// alphabetically by filename.
//
// Returns an empty map without error if hooksDir does not exist.
//
// Parameters:
//   - hooksDir: root hooks directory (e.g. .context/hooks)
//
// Returns:
//   - map[HookType][]HookInfo: discovered hooks grouped by type
//   - error: non-nil only on unexpected I/O failures
func Discover(hooksDir string) (map[HookType][]HookInfo, error) {
	result := make(map[HookType][]HookInfo)

	if _, statErr := ctxIo.SafeStat(hooksDir); os.IsNotExist(statErr) {
		return result, nil
	}

	for _, ht := range ValidTypes() {
		typeDir := filepath.Join(hooksDir, ht)

		entries, readErr := os.ReadDir(typeDir)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				continue
			}
			return nil, readErr
		}

		var hooks []HookInfo
		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			path := filepath.Join(typeDir, e.Name())

			validateErr := ValidatePath(hooksDir, path)
			if validateErr != nil {
				ctxLog.Warn(
					desc.Text(text.DescKeyTriggerSkipWarn),
					path, validateErr)
				continue
			}

			info, infoErr := e.Info()
			if infoErr != nil {
				ctxLog.Warn(warn.Readdir, typeDir, infoErr)
				continue
			}

			hooks = append(hooks, HookInfo{
				Name:    stripExt(e.Name()),
				Type:    ht,
				Path:    path,
				Enabled: info.Mode().Perm()&fs.ExecBitMask != 0,
			})
		}

		sort.Slice(hooks, func(i, j int) bool {
			return hooks[i].Name < hooks[j].Name
		})

		if len(hooks) > 0 {
			result[ht] = hooks
		}
	}

	return result, nil
}

// FindByName searches all hook type directories for a hook whose
// filename (without extension) matches name. Returns nil without
// error if no match is found.
//
// Unlike Discover, this function includes hooks regardless of their
// executable permission bit, so it can locate disabled hooks for
// enable/disable operations.
//
// Parameters:
//   - hooksDir: root hooks directory
//   - name: hook name to search for (without extension)
//
// Returns:
//   - *HookInfo: matched hook, or nil if not found
//   - error: non-nil only on unexpected I/O failures
func FindByName(hooksDir, name string) (*HookInfo, error) {
	if _, statErr := ctxIo.SafeStat(hooksDir); os.IsNotExist(statErr) {
		return nil, nil
	}

	for _, ht := range ValidTypes() {
		typeDir := filepath.Join(hooksDir, ht)

		entries, readErr := os.ReadDir(typeDir)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				continue
			}
			return nil, readErr
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if stripExt(e.Name()) == name {
				path := filepath.Join(typeDir, e.Name())
				fi, lstatErr := os.Lstat(path)
				if lstatErr != nil {
					continue
				}
				// Skip symlinks for security.
				if fi.Mode()&os.ModeSymlink != 0 {
					continue
				}
				return &HookInfo{
					Name:    name,
					Type:    ht,
					Path:    path,
					Enabled: fi.Mode().Perm()&fs.ExecBitMask != 0,
				}, nil
			}
		}
	}

	return nil, nil
}
