//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

// OpenCode config keys and provider identifiers.
const (
	// KeyProvider is the top-level provider map key.
	KeyProvider = "provider"
	// KeyBaseURL is the provider option storing the backend endpoint.
	KeyBaseURL = "baseURL"
	// KeyNPM is the optional custom provider package key.
	KeyNPM = "npm"
	// KeyOptions is the provider options object key.
	KeyOptions = "options"
	// ProviderIDAnthropic is the built-in Anthropic provider ID.
	ProviderIDAnthropic = "anthropic"
	// ProviderIDOpenAI is the built-in OpenAI provider ID.
	ProviderIDOpenAI = "openai"
	// ProviderNPMOpenAICompatible is the package for OpenAI-compatible providers.
	ProviderNPMOpenAICompatible = "@ai-sdk/openai-compatible"
)
