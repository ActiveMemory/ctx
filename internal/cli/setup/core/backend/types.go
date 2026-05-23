//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

// Result captures the outcome of an [Apply] call so the
// caller can render an accurate confirmation to the user.
//
// Fields:
//   - Path: the absolute `.ctxrc` path that was written.
//   - Created: true when `.ctxrc` did not exist before.
//   - Updated: true when an existing entry with the same
//     name was replaced in place; false when a new entry
//     was appended.
type Result struct {
	Path    string
	Created bool
	Updated bool
}

// backendEntry is the on-wire shape of a single
// `backends:` list item. Field tags drive yaml.Marshal;
// `omitempty` keeps optional fields out of the written
// document when unset, so re-runs don't introduce empty
// scalars the user didn't choose.
//
// Fields:
//   - Name: backend type label, e.g. "vllm".
//   - Endpoint: base URL with scheme.
//   - APIKeyEnv: env-var name the backend reads for auth.
//   - Timeout: duration string (Go's time.Duration.String
//     format, e.g. "30s"); empty when caller's Timeout
//     is zero.
//   - DefaultModel: model id used when a request omits
//     one.
type backendEntry struct {
	Name         string `yaml:"name"`
	Endpoint     string `yaml:"endpoint"`
	APIKeyEnv    string `yaml:"api_key_env,omitempty"`
	Timeout      string `yaml:"timeout,omitempty"`
	DefaultModel string `yaml:"default_model,omitempty"`
}
