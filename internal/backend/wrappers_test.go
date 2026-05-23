//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"testing"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// vendorWrapperCase exercises one per-vendor thin wrapper's
// default-substitution behavior: when Config leaves Name,
// Endpoint, or APIKeyEnv empty, the wrapper should fill
// them from cfgBackend defaults; when Config sets them,
// the wrapper must respect the override.
type vendorWrapperCase struct {
	label          string
	ctor           func(Config) (*openAICompat, error)
	wantName       string
	wantEndpoint   string
	wantAuthHeader bool // true if Authorization should be sent (API-key env set)
	expectAuthKey  string
	apiKeyEnvName  string
	apiKeyEnvValue string
}

func TestWrappers_DefaultsApplied(t *testing.T) {
	cases := []vendorWrapperCase{
		{
			label:          "openai",
			ctor:           newOpenAI,
			wantName:       cfgBackend.NameOpenAI,
			wantEndpoint:   cfgBackend.DefaultEndpointOpenAI,
			apiKeyEnvName:  cfgBackend.EnvOpenAIAPIKey,
			apiKeyEnvValue: "sk-test-openai",
			wantAuthHeader: true,
			expectAuthKey:  "sk-test-openai",
		},
		{
			label:          "anthropic",
			ctor:           newAnthropic,
			wantName:       cfgBackend.NameAnthropic,
			wantEndpoint:   cfgBackend.DefaultEndpointAnthropic,
			apiKeyEnvName:  cfgBackend.EnvAnthropicAPIKey,
			apiKeyEnvValue: "sk-test-anthropic",
			wantAuthHeader: true,
			expectAuthKey:  "sk-test-anthropic",
		},
		{
			label:          "ollama",
			ctor:           newOllama,
			wantName:       cfgBackend.NameOllama,
			wantEndpoint:   cfgBackend.DefaultEndpointOllama,
			wantAuthHeader: false,
		},
		{
			label:          "lmstudio",
			ctor:           newLMStudio,
			wantName:       cfgBackend.NameLMStudio,
			wantEndpoint:   cfgBackend.DefaultEndpointLMStudio,
			wantAuthHeader: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			if tc.apiKeyEnvName != "" {
				t.Setenv(tc.apiKeyEnvName, tc.apiKeyEnvValue)
			}
			b, err := tc.ctor(Config{})
			if err != nil {
				t.Fatalf("ctor: %v", err)
				return
			}
			if b.Name() != tc.wantName {
				t.Errorf("Name = %q, want %q", b.Name(), tc.wantName)
			}
			if b.base.String() != tc.wantEndpoint {
				t.Errorf("Endpoint = %q, want %q", b.base.String(), tc.wantEndpoint)
			}
			if tc.wantAuthHeader && b.apiKey != tc.expectAuthKey {
				t.Errorf("apiKey = %q, want %q", b.apiKey, tc.expectAuthKey)
			}
			if !tc.wantAuthHeader && b.apiKey != "" {
				t.Errorf("apiKey = %q, want empty (no-auth backend)", b.apiKey)
			}
		})
	}
}

func TestWrappers_UserOverridesWin(t *testing.T) {
	cases := []struct {
		label string
		ctor  func(Config) (*openAICompat, error)
	}{
		{"openai", newOpenAI},
		{"anthropic", newAnthropic},
		{"ollama", newOllama},
		{"lmstudio", newLMStudio},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			t.Setenv("CUSTOM_KEY", "sk-custom")
			b, err := tc.ctor(Config{
				Name:      "my-override",
				Endpoint:  "http://my-proxy.local:9000",
				APIKeyEnv: "CUSTOM_KEY",
			})
			if err != nil {
				t.Fatalf("ctor: %v", err)
				return
			}
			if b.Name() != "my-override" {
				t.Errorf("Name override lost: %q", b.Name())
			}
			if b.base.String() != "http://my-proxy.local:9000" {
				t.Errorf("Endpoint override lost: %q", b.base.String())
			}
			if b.apiKey != "sk-custom" {
				t.Errorf("APIKey override lost: %q", b.apiKey)
			}
		})
	}
}
