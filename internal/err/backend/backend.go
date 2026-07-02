//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import "fmt"

import cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"

// DuplicateRegistration reports a duplicate backend registration.
//
// Fields:
//   - Name: duplicate backend name
type DuplicateRegistration struct {
	Name string
}

// Error formats the duplicate registration error.
//
// Returns:
//   - string: formatted duplicate registration message
func (err DuplicateRegistration) Error() string {
	return cfgBackend.ErrDuplicateRegistration + err.Name
}

// MissingBackend reports a requested backend that is not registered.
//
// Fields:
//   - Name: missing backend name
type MissingBackend struct {
	Name string
}

// Error formats the missing backend error.
//
// Returns:
//   - string: formatted missing backend message
func (err MissingBackend) Error() string {
	return cfgBackend.ErrMissingBackend + err.Name
}

// MultipleBackends reports an ambiguous default backend selection.
type MultipleBackends struct{}

// Error formats the multiple backends error.
//
// Returns:
//   - string: formatted multiple backends message
func (err MultipleBackends) Error() string {
	return cfgBackend.ErrMultipleBackends
}

// NoBackendConfigured reports that the registry is empty.
type NoBackendConfigured struct{}

// Error formats the empty registry error.
//
// Returns:
//   - string: formatted empty registry message
func (err NoBackendConfigured) Error() string {
	return cfgBackend.ErrNoBackendConfigured
}

// Factory reports a backend factory failure.
//
// Fields:
//   - Name: backend name whose factory failed
//   - Cause: underlying factory error
type Factory struct {
	Name  string
	Cause error
}

// Error formats the factory failure.
//
// Returns:
//   - string: formatted factory failure message
func (err Factory) Error() string {
	return cfgBackend.ErrFactory + err.Name
}

// Unwrap returns the underlying factory error.
//
// Returns:
//   - error: underlying factory error
func (err Factory) Unwrap() error {
	return err.Cause
}

// InvalidEndpoint reports an endpoint URL that cannot be used.
//
// Fields:
//   - Endpoint: configured endpoint value
//   - Cause: underlying parse or request construction failure
type InvalidEndpoint struct {
	Endpoint string
	Cause    error
}

// Error formats the invalid endpoint error.
//
// Returns:
//   - string: formatted invalid endpoint message
func (err InvalidEndpoint) Error() string {
	return cfgBackend.ErrInvalidEndpoint + err.Endpoint
}

// Unwrap returns the underlying endpoint error.
//
// Returns:
//   - error: underlying endpoint error
func (err InvalidEndpoint) Unwrap() error {
	return err.Cause
}

// Unreachable reports a failed backend HTTP request.
//
// Fields:
//   - Name: backend name
//   - Endpoint: configured endpoint value
//   - Cause: underlying HTTP error
type Unreachable struct {
	Name     string
	Endpoint string
	Cause    error
}

// Error formats the unreachable backend error.
//
// Returns:
//   - string: formatted unreachable backend message
func (err Unreachable) Error() string {
	return cfgBackend.ErrUnreachable + err.Name
}

// Unwrap returns the underlying HTTP error.
//
// Returns:
//   - error: underlying HTTP error
func (err Unreachable) Unwrap() error {
	return err.Cause
}

// Upstream reports a non-success status from the backend.
//
// Fields:
//   - Name: backend name
//   - StatusCode: upstream HTTP status code
//   - Body: upstream response body
type Upstream struct {
	Name       string
	StatusCode int
	Body       string
}

// Error formats the upstream status error.
//
// Returns:
//   - string: formatted upstream response message
func (err Upstream) Error() string {
	return cfgBackend.ErrUpstream + fmt.Sprintf(
		cfgBackend.FmtUpstreamStatusBody,
		err.StatusCode,
		err.Body,
	)
}

// BadRequest wraps a local request encoding or decode failure.
//
// Fields:
//   - Name: backend name
//   - Cause: underlying request error
type BadRequest struct {
	Name  string
	Cause error
}

// Error formats the backend request failure.
//
// Returns:
//   - string: formatted request failure message
func (err BadRequest) Error() string {
	return cfgBackend.ErrBadRequest + err.Name
}

// Unwrap returns the underlying request failure.
//
// Returns:
//   - error: underlying request failure
func (err BadRequest) Unwrap() error {
	return err.Cause
}

// InvalidResponseShape reports a structurally incomplete provider response.
//
// Returns:
//   - error: "invalid structured response shape"
func InvalidResponseShape() error {
	return fmt.Errorf(cfgBackend.ErrInvalidResponseShape)
}
