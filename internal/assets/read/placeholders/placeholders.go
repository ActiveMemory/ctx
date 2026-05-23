//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package placeholders

import (
	"path"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	errPlaceholders "github.com/ActiveMemory/ctx/internal/err/placeholders"
	"github.com/ActiveMemory/ctx/internal/i18n"
)

// loadedMu guards `loaded` for safe concurrent reads
// (validator hot path) and rare writes (first call per
// locale + test-only Reset).
var loadedMu sync.RWMutex

// loaded caches the per-locale folded set. Each entry is
// the set of i18n.Fold(placeholder) values for fast O(1)
// membership testing in the validator hot path. Populated
// lazily on first Load(locale).
var loaded = make(map[string]map[string]struct{})

// Load returns the case-folded placeholder set for locale,
// memoized on first call per locale. Callers compare
// against `i18n.Fold(strings.TrimSpace(input))`.
//
// Parameters:
//   - locale: locale identifier matching a file under
//     `<DirI18nPlaceholders>/<locale>.yaml`. Use the
//     constants from `internal/config/asset` (e.g.
//     `asset.LocaleEN`).
//
// Returns:
//   - map[string]struct{}: the folded placeholder set.
//     Keys are already case-folded; callers must fold
//     their input before lookup.
//   - error: non-nil if the locale file is missing or the
//     YAML is malformed.
func Load(locale string) (map[string]struct{}, error) {
	loadedMu.RLock()
	if cached, ok := loaded[locale]; ok {
		loadedMu.RUnlock()
		return cached, nil
	}
	loadedMu.RUnlock()

	p := path.Join(asset.DirI18nPlaceholders, locale+cfgFile.ExtYAML)
	data, readErr := assets.FS.ReadFile(p)
	if readErr != nil {
		return nil, errPlaceholders.ReadLocale(locale, readErr)
	}
	var parsed file
	if parseErr := yaml.Unmarshal(data, &parsed); parseErr != nil {
		return nil, errPlaceholders.ParseLocale(locale, parseErr)
	}

	set := make(map[string]struct{}, len(parsed.Placeholders))
	for _, raw := range parsed.Placeholders {
		set[i18n.Fold(raw)] = struct{}{}
	}

	loadedMu.Lock()
	loaded[locale] = set
	loadedMu.Unlock()
	return set, nil
}

// Reset clears the in-process cache. Test-only; production
// callers should treat the cache as permanent for the
// lifetime of the binary.
func Reset() {
	loadedMu.Lock()
	loaded = make(map[string]map[string]struct{})
	loadedMu.Unlock()
}
