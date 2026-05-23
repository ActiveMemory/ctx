//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

// OpenAI request-shape constants used by the dispatch
// payload (role labels + JSON-mode response format) plus
// the temperature-unset sentinel.
const (
	// RoleSystem is the OpenAI system-role label used in
	// the dispatch payload.
	RoleSystem = "system"
	// RoleUser is the OpenAI user-role label used in the
	// dispatch payload.
	RoleUser = "user"
	// ResponseFormatType is the OpenAI response_format
	// type value asking the model to emit valid JSON.
	// Honoured by vLLM, OpenAI, Anthropic (OpenAI-compat
	// mode), and others; ignored backends still return
	// raw text for the caller to validate.
	ResponseFormatType = "json_object"
	// UnsetTemperature is the sentinel value
	// (Request.Temperature negative = unset) used so the
	// dispatch sends no `temperature` field on the wire
	// and the backend uses its own default.
	UnsetTemperature = -1.0
)

// Slug + header constants for the generated proposal
// markdown file.
const (
	// ProposalSlug is the per-extract filename slug; the
	// timestamp is the discriminating component, so the
	// slug stays stable.
	ProposalSlug = "extract"
	// ProposalTemplate is the markdown-format template
	// for the generated proposal file. Interpolation
	// order: backend name, model id, JSON body. Trailing
	// newline on the JSON body is the caller's
	// responsibility.
	ProposalTemplate = "# ctx ai extract proposal\n\n" +
		"- backend: `%s`\n" +
		"- model: `%s`\n\n" +
		"```json\n%s```\n"
)

// SystemPrompt is the system message sent ahead of the
// user input. Asks the model to extract candidates for
// canonical ctx files (decisions, learnings, tasks, open
// questions) in a stable JSON shape that future
// ratification skills can parse.
const SystemPrompt = "" +
	"You extract project knowledge from text. " +
	"Respond ONLY with a single JSON object. " +
	"Schema (all keys optional, arrays may be empty):\n" +
	"  {\n" +
	"    \"decisions\":    [{" +
	"\"title\":\"...\", \"context\":\"...\"," +
	" \"rationale\":\"...\"}],\n" +
	"    \"learnings\":    [{" +
	"\"title\":\"...\", \"context\":\"...\"," +
	" \"lesson\":\"...\", \"application\":\"...\"}],\n" +
	"    \"tasks\":        [{" +
	"\"summary\":\"...\"," +
	" \"priority\":\"low|medium|high\"}],\n" +
	"    \"open_questions\":[{" +
	"\"question\":\"...\", \"context\":\"...\"}]\n" +
	"  }\n" +
	"Omit any category with no candidates. Do not " +
	"include commentary outside the JSON object."
