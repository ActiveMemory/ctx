//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"

// ResolveEndpoint returns the explicit endpoint or the backend default.
//
// Parameters:
//   - name: configured backend name
//   - endpoint: explicit endpoint flag value, if any
//
// Returns:
//   - string: explicit endpoint when set, otherwise the backend default
func ResolveEndpoint(name string, endpoint string) string {
	if endpoint != "" {
		return endpoint
	}
	switch name {
	case cfgBackend.NameVLLM:
		return cfgBackend.DefaultEndpointVLLM
	case cfgBackend.NameOpenAI:
		return cfgBackend.DefaultEndpointOpenAI
	case cfgBackend.NameAnthropic:
		return cfgBackend.DefaultEndpointAnthropic
	case cfgBackend.NameOllama:
		return cfgBackend.DefaultEndpointOllama
	case cfgBackend.NameLMStudio:
		return cfgBackend.DefaultEndpointLMStudio
	default:
		return ""
	}
}
