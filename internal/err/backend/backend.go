//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrBackendNotFound signals a Resolve call naming a
	// backend that has no Config or no Factory registered.
	// Constructor [NotFound] wraps it with the requested
	// name.
	ErrBackendNotFound = entity.Sentinel(
		text.DescKeyErrBackendBackendNotFoundMsg,
	)
	// ErrNoBackends signals a Resolve or Default call
	// against a Registry with no backends configured at
	// all. Recovery: `ctx setup --backend <name>`.
	ErrNoBackends = entity.Sentinel(
		text.DescKeyErrBackendNoBackends,
	)
	// ErrAmbiguousDefault signals a Default call against a
	// Registry where more than one backend is configured
	// and no default has been set via SetDefault.
	// Recovery: pass `--backend <name>` or set
	// `[backends].default` in `.ctxrc`.
	ErrAmbiguousDefault = entity.Sentinel(
		text.DescKeyErrBackendAmbiguousDefault,
	)
	// ErrDuplicateRegistration signals a Register call for
	// a type name already bound to a Factory. Constructor
	// [DuplicateRegistration] wraps it with the offending
	// name.
	ErrDuplicateRegistration = entity.Sentinel(
		text.DescKeyErrBackendDuplicateRegistrationMsg,
	)
	// ErrMissingEndpoint signals a Factory call against a
	// Config whose Endpoint is empty. Wrapped by
	// [MissingEndpoint] with the backend name.
	ErrMissingEndpoint = entity.Sentinel(
		text.DescKeyErrBackendMissingEndpointMsg,
	)
	// ErrInvalidEndpoint signals an Endpoint that does not
	// parse as a URL or uses a scheme other than http/https.
	// Wrapped by [InvalidEndpoint] with name and the raw
	// endpoint string.
	ErrInvalidEndpoint = entity.Sentinel(
		text.DescKeyErrBackendInvalidEndpointMsg,
	)
	// ErrMissingModel signals a Complete call against a
	// Request whose Model is empty and whose backend Config
	// has no DefaultModel. Wrapped by [MissingModel].
	ErrMissingModel = entity.Sentinel(
		text.DescKeyErrBackendMissingModelMsg,
	)
	// ErrUnreachable signals a transport-layer failure
	// (DNS, dial, TCP) talking to the backend endpoint.
	// Distinct from ErrUnhealthyStatus / ErrUpstreamStatus
	// because recovery differs: server not running, wrong
	// endpoint, network partition. Wrapped by
	// [Unreachable] with name, endpoint, cause.
	ErrUnreachable = entity.Sentinel(
		text.DescKeyErrBackendUnreachableMsg,
	)
	// ErrUnhealthyStatus signals a non-200 response from
	// the Ping endpoint (/v1/models). Wrapped by
	// [UnhealthyStatus] with name, status, body excerpt.
	ErrUnhealthyStatus = entity.Sentinel(
		text.DescKeyErrBackendUnhealthyStatusMsg,
	)
	// ErrUpstreamStatus signals a non-200 response from
	// the chat-completions endpoint. Auth failures (401,
	// 403), rate-limits (429), and model-not-found (404)
	// all reach callers through this sentinel so the CLI
	// can branch on body.
	ErrUpstreamStatus = entity.Sentinel(
		text.DescKeyErrBackendUpstreamStatusMsg,
	)
	// ErrEmptyModels signals that `/v1/models` returned
	// HTTP 200 with an empty `data` array. The backend is
	// reachable but has no models loaded — distinct from
	// transport failure so `ctx ai ping` can report a
	// usable recovery hint. Wrapped by [EmptyModels].
	ErrEmptyModels = entity.Sentinel(
		text.DescKeyErrBackendEmptyModelsMsg,
	)
)

// MissingEndpoint wraps [ErrMissingEndpoint] with the
// backend name.
//
// Parameters:
//   - name: the backend type label.
//
// Returns:
//   - error: wraps [ErrMissingEndpoint] for errors.Is.
func MissingEndpoint(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendMissingEndpoint),
		ErrMissingEndpoint, name,
	)
}

// InvalidEndpoint wraps [ErrInvalidEndpoint] with name
// and the raw endpoint string.
//
// Parameters:
//   - name: the backend type label.
//   - endpoint: the offending endpoint string.
//
// Returns:
//   - error: wraps [ErrInvalidEndpoint] for errors.Is.
func InvalidEndpoint(name, endpoint string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendInvalidEndpoint),
		ErrInvalidEndpoint, name, endpoint,
	)
}

// MissingModel wraps [ErrMissingModel] with the backend
// name so the user knows which backend lacks a default.
//
// Parameters:
//   - name: the backend type label.
//
// Returns:
//   - error: wraps [ErrMissingModel] for errors.Is.
func MissingModel(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendMissingModel),
		ErrMissingModel, name,
	)
}

// Unreachable wraps [ErrUnreachable] with name, endpoint,
// and the underlying transport error.
//
// Parameters:
//   - name: the backend type label.
//   - endpoint: the endpoint the client tried to reach.
//   - cause: the underlying net/http error.
//
// Returns:
//   - error: wraps [ErrUnreachable] for errors.Is.
func Unreachable(name, endpoint string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendUnreachable),
		ErrUnreachable, name, endpoint, cause,
	)
}

// UnhealthyStatus wraps [ErrUnhealthyStatus] with name,
// HTTP status, and a truncated body excerpt.
//
// Parameters:
//   - name: the backend type label.
//   - status: the HTTP status code returned by Ping.
//   - body: a truncated excerpt of the response body.
//
// Returns:
//   - error: wraps [ErrUnhealthyStatus] for errors.Is.
func UnhealthyStatus(name string, status int, body string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendUnhealthyStatus),
		ErrUnhealthyStatus, name, status, body,
	)
}

// UpstreamStatus wraps [ErrUpstreamStatus] with name,
// HTTP status, and a truncated body excerpt.
//
// Parameters:
//   - name: the backend type label.
//   - status: the HTTP status code from chat-completions.
//   - body: a truncated excerpt of the upstream body.
//
// Returns:
//   - error: wraps [ErrUpstreamStatus] for errors.Is.
func UpstreamStatus(name string, status int, body string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendUpstreamStatus),
		ErrUpstreamStatus, name, status, body,
	)
}

// MarshalRequest wraps a json.Marshal failure on the
// outgoing chat-completion request payload.
//
// Parameters:
//   - name: the backend type label.
//   - cause: the underlying json encode error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func MarshalRequest(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendMarshalRequest),
		name, cause,
	)
}

// ParseResponse wraps a json.Unmarshal failure on the
// upstream chat-completion response body.
//
// Parameters:
//   - name: the backend type label.
//   - cause: the underlying json decode error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ParseResponse(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendParseResponse),
		name, cause,
	)
}

// ReadResponse wraps an io.ReadAll failure draining the
// upstream response body.
//
// Parameters:
//   - name: the backend type label.
//   - cause: the underlying read error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ReadResponse(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendReadResponse),
		name, cause,
	)
}

// EmptyChoices wraps the case where a chat-completion
// response parsed cleanly but had no entries in the
// choices array.
//
// Parameters:
//   - name: the backend type label.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func EmptyChoices(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendEmptyChoices),
		name,
	)
}

// EmptyModels wraps [ErrEmptyModels] with the backend
// name so users can tell which backend has no models
// loaded.
//
// Parameters:
//   - name: the backend type label.
//
// Returns:
//   - error: wraps [ErrEmptyModels] for errors.Is.
func EmptyModels(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendEmptyModels),
		ErrEmptyModels, name,
	)
}

// NotFound wraps [ErrBackendNotFound] with the requested
// backend name so error output names the missing key.
//
// Parameters:
//   - name: the backend type label the caller asked for.
//
// Returns:
//   - error: wraps [ErrBackendNotFound] for errors.Is.
func NotFound(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendBackendNotFound),
		ErrBackendNotFound, name,
	)
}

// DuplicateRegistration wraps [ErrDuplicateRegistration]
// with the offending type name so error output identifies
// the conflict.
//
// Parameters:
//   - name: the type label that was already bound.
//
// Returns:
//   - error: wraps [ErrDuplicateRegistration] for
//     errors.Is.
func DuplicateRegistration(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendDuplicateRegistration),
		ErrDuplicateRegistration, name,
	)
}
