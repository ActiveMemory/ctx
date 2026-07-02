//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	cfgOpenCode "github.com/ActiveMemory/ctx/internal/config/opencode"
)

// providerDetails returns the OpenCode provider ID and optional npm package.
//
// Parameters:
//   - name: configured ctx backend name
//
// Returns:
//   - string: provider ID
//   - string: npm package for custom providers, or empty
//   - bool: true when the backend maps to an OpenCode provider
func providerDetails(name string) (string, string, bool) {
	switch name {
	case cfgBackend.NameOpenAI:
		return cfgOpenCode.ProviderIDOpenAI, "", true
	case cfgBackend.NameAnthropic:
		return cfgOpenCode.ProviderIDAnthropic, "", true
	case cfgBackend.NameOpenAICompatible,
		cfgBackend.NameVLLM,
		cfgBackend.NameOllama,
		cfgBackend.NameLMStudio:
		return name, cfgOpenCode.ProviderNPMOpenAICompatible, true
	default:
		return "", "", false
	}
}

// mergeProviderConfig overlays one provider entry into an OpenCode config map.
//
// Parameters:
//   - config: parsed OpenCode config map to mutate in place
//   - providerID: target provider identifier
//   - npmPackage: optional custom npm package for the provider
//   - endpoint: backend base URL to write into provider options
//
// Returns:
//   - none
func mergeProviderConfig(
	config map[string]any,
	providerID string,
	npmPackage string,
	endpoint string,
) {
	providers, _ := config[cfgOpenCode.KeyProvider].(map[string]any)
	if providers == nil {
		providers = map[string]any{}
	}
	entry, _ := providers[providerID].(map[string]any)
	if entry == nil {
		entry = map[string]any{}
	}
	if npmPackage != "" {
		entry[cfgOpenCode.KeyNPM] = npmPackage
	}
	options, _ := entry[cfgOpenCode.KeyOptions].(map[string]any)
	if options == nil {
		options = map[string]any{}
	}
	options[cfgOpenCode.KeyBaseURL] = endpoint
	entry[cfgOpenCode.KeyOptions] = options
	providers[providerID] = entry
	config[cfgOpenCode.KeyProvider] = providers
}
