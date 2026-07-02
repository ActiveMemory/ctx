//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

// Options controls backend setup output.
//
// Fields:
//   - Name: backend name to configure
//   - Endpoint: backend HTTP endpoint
//   - APIKeyEnv: environment variable containing credentials
//   - Model: default model name
//   - Timeout: backend request timeout duration
//   - Write: whether to write .ctxrc instead of printing a snippet
type Options struct {
	Name      string
	Endpoint  string
	APIKeyEnv string
	Model     string
	Timeout   string
	Write     bool
}
