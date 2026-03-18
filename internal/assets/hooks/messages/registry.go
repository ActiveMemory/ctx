//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package messages provides metadata for hook message templates.
//
// The registry is parsed from an embedded registry.yaml that lives
// alongside the .txt template files. It is the metadata layer over
// the embedded FS and changes only when hooks are added or removed.
package messages

import (
	"sync"

	"github.com/ActiveMemory/ctx/internal/assets"
	fserr "github.com/ActiveMemory/ctx/internal/err/fs"
	errparser "github.com/ActiveMemory/ctx/internal/err/parser"
	"gopkg.in/yaml.v3"
)

// HookMessageInfo describes a single hook message template entry.
type HookMessageInfo struct {
	// Hook is the hook directory name (e.g., "qa-reminder").
	Hook string `yaml:"hook"`

	// Variant is the template file stem (e.g., "gate").
	Variant string `yaml:"variant"`

	// Category is "customizable" or "ctx-specific".
	Category string `yaml:"category"`

	// Description is a one-line human description of this message.
	Description string `yaml:"description"`

	// TemplateVars lists available Go template variables (e.g., "PromptsSinceNudge").
	TemplateVars []string `yaml:"vars,omitempty"`
}

// CategoryCustomizable marks messages intended for project-specific customization.
const CategoryCustomizable = "customizable"

// CategoryCtxSpecific marks messages specific to ctx's own development workflow.
const CategoryCtxSpecific = "ctx-specific"

var (
	registryOnce sync.Once
	registryData []HookMessageInfo
	registryErr  error
)

// Registry returns the list of all hook message entries parsed from
// the embedded registry.yaml.
//
// Returns:
//   - []HookMessageInfo: All entries sorted by hook then variant
func Registry() []HookMessageInfo {
	registryOnce.Do(func() {
		raw, readErr := assets.HookMessageRegistry()
		if readErr != nil {
			registryErr = fserr.FileRead("registry.yaml", readErr)
			return
		}
		var entries []HookMessageInfo
		if unmarshalErr := yaml.Unmarshal(raw, &entries); unmarshalErr != nil {
			registryErr = errparser.ParseFile("registry.yaml", unmarshalErr)
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
func RegistryError() error {
	Registry() // ensure sync.Once has run
	return registryErr
}

// Lookup returns the HookMessageInfo for the given hook and variant,
// or nil if not found.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//   - variant: Template file stem (e.g., "gate")
//
// Returns:
//   - *HookMessageInfo: The matching entry, or nil
func Lookup(hook, variant string) *HookMessageInfo {
	for _, info := range Registry() {
		if info.Hook == hook && info.Variant == variant {
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
//   - hook: Hook directory name
//
// Returns:
//   - []string: Variant names for the hook, or nil if hook not found
func Variants(hook string) []string {
	var variants []string
	for _, info := range Registry() {
		if info.Hook == hook {
			variants = append(variants, info.Variant)
		}
	}
	return variants
}
