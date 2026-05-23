//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"fmt"
	"strings"

	cfgExtract "github.com/ActiveMemory/ctx/internal/config/extract"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
)

// Compose wraps the backend's JSON response in a minimal
// markdown envelope so the proposal file is
// self-describing when a human (or future ratification
// skill) opens it. The template is centralised in
// cfgExtract.ProposalTemplate.
//
// Parameters:
//   - backendName: the backend type label.
//   - model: model id reported by the backend.
//   - jsonBody: the chat completion's content (expected
//     to be valid JSON because the request used
//     response_format: json_object).
//
// Returns:
//   - string: the markdown body to write.
func Compose(backendName, model, jsonBody string) string {
	if !strings.HasSuffix(jsonBody, cfgToken.NewlineLF) {
		jsonBody += cfgToken.NewlineLF
	}
	return fmt.Sprintf(
		cfgExtract.ProposalTemplate, backendName, model, jsonBody,
	)
}
