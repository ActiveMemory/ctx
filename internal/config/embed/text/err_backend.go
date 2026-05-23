//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for AI-backend registry errors. The matching
// YAML entries live in commands/text/errors.yaml under
// the `err.backend.*` namespace; sentinels and
// constructors in internal/err/backend/ resolve them via
// desc.Text at error construction time.
const (
	// DescKeyErrBackendBackendNotFoundMsg is the sentinel
	// message lookup key for ErrBackendNotFound.
	DescKeyErrBackendBackendNotFoundMsg = "err.backend.backend-not-found-msg"
	// DescKeyErrBackendBackendNotFound is the wrapper
	// format key for the BackendNotFound constructor;
	// interpolates the sentinel and the requested name.
	DescKeyErrBackendBackendNotFound = "err.backend.backend-not-found"
	// DescKeyErrBackendNoBackends is the sentinel message
	// key for ErrNoBackends (no wrapper; the sentinel
	// itself carries the user-facing recovery hint).
	DescKeyErrBackendNoBackends = "err.backend.no-backends"
	// DescKeyErrBackendAmbiguousDefault is the sentinel
	// message key for ErrAmbiguousDefault.
	DescKeyErrBackendAmbiguousDefault = "err.backend.ambiguous-default"
	// DescKeyErrBackendDuplicateRegistrationMsg is the
	// sentinel message lookup key for
	// ErrDuplicateRegistration.
	DescKeyErrBackendDuplicateRegistrationMsg = "err.backend.duplicate-registration-msg"
	// DescKeyErrBackendDuplicateRegistration is the
	// wrapper format key for the DuplicateRegistration
	// constructor.
	DescKeyErrBackendDuplicateRegistration = "err.backend.duplicate-registration"

	// DescKeyErrBackendMissingEndpointMsg is the sentinel
	// message key for ErrMissingEndpoint.
	DescKeyErrBackendMissingEndpointMsg = "err.backend.missing-endpoint-msg"
	// DescKeyErrBackendMissingEndpoint is the wrapper
	// format key for MissingEndpoint; interpolates sentinel
	// and the offending backend name.
	DescKeyErrBackendMissingEndpoint = "err.backend.missing-endpoint"
	// DescKeyErrBackendInvalidEndpointMsg is the sentinel
	// message key for ErrInvalidEndpoint.
	DescKeyErrBackendInvalidEndpointMsg = "err.backend.invalid-endpoint-msg"
	// DescKeyErrBackendInvalidEndpoint is the wrapper
	// format key for InvalidEndpoint; interpolates name
	// and the endpoint string.
	DescKeyErrBackendInvalidEndpoint = "err.backend.invalid-endpoint"
	// DescKeyErrBackendMissingModelMsg is the sentinel
	// message key for ErrMissingModel.
	DescKeyErrBackendMissingModelMsg = "err.backend.missing-model-msg"
	// DescKeyErrBackendMissingModel is the wrapper format
	// key for MissingModel; interpolates backend name.
	DescKeyErrBackendMissingModel = "err.backend.missing-model"
	// DescKeyErrBackendUnreachableMsg is the sentinel
	// message key for ErrUnreachable.
	DescKeyErrBackendUnreachableMsg = "err.backend.unreachable-msg"
	// DescKeyErrBackendUnreachable is the wrapper format
	// key for Unreachable; interpolates name, endpoint,
	// and underlying transport error.
	DescKeyErrBackendUnreachable = "err.backend.unreachable"
	// DescKeyErrBackendUnhealthyStatusMsg is the sentinel
	// message key for ErrUnhealthyStatus.
	DescKeyErrBackendUnhealthyStatusMsg = "err.backend.unhealthy-status-msg"
	// DescKeyErrBackendUnhealthyStatus is the wrapper
	// format key for UnhealthyStatus; interpolates name,
	// HTTP status, body excerpt.
	DescKeyErrBackendUnhealthyStatus = "err.backend.unhealthy-status"
	// DescKeyErrBackendUpstreamStatusMsg is the sentinel
	// message key for ErrUpstreamStatus.
	DescKeyErrBackendUpstreamStatusMsg = "err.backend.upstream-status-msg"
	// DescKeyErrBackendUpstreamStatus is the wrapper
	// format key for UpstreamStatus; interpolates name,
	// HTTP status, body excerpt.
	DescKeyErrBackendUpstreamStatus = "err.backend.upstream-status"
	// DescKeyErrBackendMarshalRequest is the wrapper
	// format key for MarshalRequest; interpolates the
	// backend name and the underlying marshal error.
	DescKeyErrBackendMarshalRequest = "err.backend.marshal-request"
	// DescKeyErrBackendParseResponse is the wrapper format
	// key for ParseResponse; interpolates name and cause.
	DescKeyErrBackendParseResponse = "err.backend.parse-response"
	// DescKeyErrBackendReadResponse is the wrapper format
	// key for ReadResponse; interpolates name and cause.
	DescKeyErrBackendReadResponse = "err.backend.read-response"
	// DescKeyErrBackendEmptyChoices is the wrapper format
	// key for EmptyChoices; interpolates the backend name.
	DescKeyErrBackendEmptyChoices = "err.backend.empty-choices"
)
