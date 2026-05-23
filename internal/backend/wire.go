//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

// chatRequest mirrors OpenAI's /v1/chat/completions
// request payload. Field tags drive json.Marshal; the
// `omitempty` set is deliberate — Temperature uses a
// pointer so the wire payload can omit it when callers
// don't set a value (Request's "negative = unset"
// convention).
//
// Fields:
//   - Model: target model id.
//   - Messages: ordered chat history.
//   - MaxTokens: upper bound on output tokens; omitted
//     when zero.
//   - Temperature: pointer so 0 vs unset is
//     distinguishable on the wire.
//   - ResponseFormat: structured-output spec; nil omits
//     the field.
type chatRequest struct {
	Model          string                     `json:"model"`
	Messages       []chatRequestMessage       `json:"messages"`
	MaxTokens      int                        `json:"max_tokens,omitempty"`
	Temperature    *float64                   `json:"temperature,omitempty"`
	ResponseFormat *chatRequestResponseFormat `json:"response_format,omitempty"`
}

// chatRequestMessage is the wire shape of a single
// chat message.
//
// Fields:
//   - Role: OpenAI chat role.
//   - Content: message body.
type chatRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequestResponseFormat is the wire shape of the
// structured-output spec.
//
// Fields:
//   - Type: "json_object" or "json_schema".
//   - JSONSchema: schema payload when Type is
//     "json_schema"; omitted otherwise.
type chatRequestResponseFormat struct {
	Type       string         `json:"type"`
	JSONSchema map[string]any `json:"json_schema,omitempty"`
}

// chatResponse mirrors OpenAI's /v1/chat/completions
// response payload. Fields not consumed by ctx (id,
// object, created, usage) are intentionally absent from
// this Go shape so the decoder ignores them and a wire
// version bump doesn't break us.
//
// Fields:
//   - Model: the model id the server reports.
//   - Choices: the completion choices array.
type chatResponse struct {
	Model   string               `json:"model"`
	Choices []chatResponseChoice `json:"choices"`
}

// chatResponseChoice is one entry in the choices array.
//
// Fields:
//   - Message: the assistant reply.
//   - FinishReason: "stop", "length", etc.; preserved
//     in case future callers want to branch on it.
type chatResponseChoice struct {
	Message      chatResponseMessage `json:"message"`
	FinishReason string              `json:"finish_reason"`
}

// chatResponseMessage is the message payload inside a
// choice.
//
// Fields:
//   - Role: typically "assistant".
//   - Content: the message body.
type chatResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
