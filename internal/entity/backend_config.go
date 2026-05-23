//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import "time"

// BackendConfig is the per-backend slice of `.ctxrc`
// settings consumed by the AI-backend registry. Lives in
// `entity/` because both `internal/backend/` (which
// constructs Backend instances from it) and
// `internal/rc/` (which loads it from YAML) reference
// it; cross-package types belong here by audit rule.
//
// Fields:
//   - Name: backend type label, e.g., "vllm". Must match
//     a registered Factory at resolve time.
//   - Endpoint: base URL, e.g., "http://localhost:8000".
//     Path segments below `/v1/...` are appended by the
//     backend implementation; this value should not end
//     in a slash.
//   - APIKeyEnv: name of the environment variable that
//     holds the bearer token. Empty means no auth header
//     is sent (the vLLM default).
//   - Timeout: per-request deadline. Zero means use the
//     backend implementation's default.
//   - DefaultModel: model ID to use when a request does
//     not specify one. Empty means require explicit
//     model selection at call time.
type BackendConfig struct {
	Name         string
	Endpoint     string
	APIKeyEnv    string
	Timeout      time.Duration
	DefaultModel string
}
