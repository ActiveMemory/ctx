//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package messages

import (
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
)

// CategoryCustomizable marks messages intended for
// project-specific customization.
const CategoryCustomizable = "customizable"

// CategoryCtxSpecific marks messages specific to ctx's
// own development workflow.
const CategoryCtxSpecific = "ctx-specific"

// registryOnce, registryData, and registryErr cache the parsed hook
// message registry loaded once via sync.Once.
var (
	// registryOnce guards one-time registry loading.
	registryOnce sync.Once
	// registryData holds the parsed hook message entries.
	registryData []HookMessageInfo
	// registryErr stores any error from registry loading.
	registryErr error
)

// Registry returns the list of all hook message entries parsed from
// the embedded registry.yaml.
//
// Returns:
//   - []HookMessageInfo: All entries sorted by hook then variant
func Registry() []HookMessageInfo {
	registryOnce.Do(func() {
		raw, readErr := hook.MessageRegistry()
		if readErr != nil {
			registryErr = errFs.FileRead(asset.FileRegistryYAML, readErr)
			return
		}
		var entries []HookMessageInfo
		if unmarshalErr := yaml.Unmarshal(raw, &entries); unmarshalErr != nil {
			registryErr = errParser.ParseFile(asset.FileRegistryYAML, unmarshalErr)
			return
		}
		registryData = entries
	})
	if registryErr != nil {
		return nil
	}
	return registryData
}

// RegistryError returns any error encountered while parsing the
// embedded registry.yaml. Nil on success.
//
// Returns:
//   - error: Parse error from registry.yaml, or nil on success
func RegistryError() error {
	Registry() // ensure sync.Once has run
	return registryErr
}

// Lookup returns the HookMessageInfo for the given hook and variant,
// or nil if not found.
//
// Parameters:
//   - hookName: Hook directory name (e.g., "qa-reminder")
//   - variant: Template file stem (e.g., "gate")
//
// Returns:
//   - *HookMessageInfo: The matching entry, or nil
func Lookup(hookName, variant string) *HookMessageInfo {
	for _, info := range Registry() {
		if info.Hook == hookName && info.Variant == variant {
			return &info
		}
	}
	return nil
}

// Hooks returns a deduplicated list of hook names in the registry.
//
// Returns:
//   - []string: Hook names in alphabetical order
func Hooks() []string {
	seen := make(map[string]bool)
	var hooks []string
	for _, info := range Registry() {
		if !seen[info.Hook] {
			seen[info.Hook] = true
			hooks = append(hooks, info.Hook)
		}
	}
	return hooks
}

// Variants returns the variant names for a given hook.
//
// Parameters:
//   - hookName: Hook directory name
//
// Returns:
//   - []string: Variant names for the hook, or nil if not found
func Variants(hookName string) []string {
	var variants []string
	for _, info := range Registry() {
		if info.Hook == hookName {
			variants = append(variants, info.Variant)
		}
	}
	return variants
}
