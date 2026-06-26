//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// Format strings for sentinel-wrapping in err/context constructors.
// Centralized here so the magic-string audit (which exempts
// internal/config) does not flag them at the call site.
const (
	// FmtWrapColon wraps a sentinel and a tailored message:
	//   fmt.Errorf(FmtWrapColon, ErrFoo, "tailored detail")
	//   ↦ "<ErrFoo.Error()>: tailored detail".
	FmtWrapColon = "%w: %s"

	// ErrBackendsMapping reports a non-mapping backends value.
	ErrBackendsMapping = "backends must be a mapping"
	// ErrBackendsDefaultScalar reports a non-scalar default backend value.
	ErrBackendsDefaultScalar = "backends.default must be a scalar"
	// ErrBackendsDefaultMissing reports a default backend with no definition.
	ErrBackendsDefaultMissing = "backends.default references missing backend: "
	// ErrBackendsEndpointRequired reports a backend missing its endpoint.
	ErrBackendsEndpointRequired = "backends.%s.endpoint is required"
	// ErrBackendsEndpointScheme reports a backend endpoint with invalid scheme.
	ErrBackendsEndpointScheme = "backends.%s.endpoint must be http or https"
	// ErrBackendsUnknownField reports an unknown key under a named backend.
	ErrBackendsUnknownField = "backends.%s.%s"

	// BackendDefaultKey is the reserved key for default backend selection.
	BackendDefaultKey = "default"
	// BackendsKey is the top-level backend configuration key.
	BackendsKey = "backends"
	// BackendTypeKey is the backend implementation type key.
	BackendTypeKey = "type"
	// BackendEndpointKey is the endpoint URL key.
	BackendEndpointKey = "endpoint"
	// BackendAPIKeyEnvKey is the credential environment variable key.
	BackendAPIKeyEnvKey = "api_key_env"
	// BackendTimeoutKey is the request timeout key.
	BackendTimeoutKey = "timeout"
	// BackendDefaultModelKey is the default model key.
	BackendDefaultModelKey = "default_model"
)
