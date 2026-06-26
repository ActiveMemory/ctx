//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	cfgHTTP "github.com/ActiveMemory/ctx/internal/config/http"
	cfgRC "github.com/ActiveMemory/ctx/internal/config/rc"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	setupErr "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
	rcPkg "github.com/ActiveMemory/ctx/internal/rc"
)

// content returns merged .ctxrc YAML for backend setup.
//
// Parameters:
//   - options: resolved backend setup options
//
// Returns:
//   - []byte: YAML content
//   - error: read or marshal failure
func content(options Options) ([]byte, error) {
	root, readErr := readRoot()
	if readErr != nil {
		return nil, readErr
	}
	backends := mapping(root, cfgRC.BackendsKey)
	if scalar(backends, cfgRC.BackendDefaultKey) == nil {
		setScalar(backends, cfgRC.BackendDefaultKey, options.Name)
	}
	entry := mapping(backends, options.Name)
	setScalar(entry, cfgRC.BackendTypeKey, options.Name)
	setScalar(entry, cfgRC.BackendEndpointKey, options.Endpoint)
	setScalar(entry, cfgRC.BackendAPIKeyEnvKey, options.APIKeyEnv)
	setScalar(entry, cfgRC.BackendTimeoutKey, options.Timeout)
	setScalar(entry, cfgRC.BackendDefaultModelKey, options.Model)
	return yaml.Marshal(root)
}

// readRoot reads .ctxrc into a YAML mapping node.
//
// Returns:
//   - *yaml.Node: root document mapping
//   - error: read or decode failure
func readRoot() (*yaml.Node, error) {
	data, readErr := ctxio.SafeReadUserFile(cfgSetup.FileCtxRC)
	if readErr != nil && !os.IsNotExist(readErr) {
		return nil, readErr
	}
	root := &yaml.Node{Kind: yaml.MappingNode}
	if len(data) == 0 {
		return root, nil
	}
	if _, validateErr := rcPkg.Validate(data); validateErr != nil {
		return nil, validateErr
	}
	var doc yaml.Node
	if decodeErr := yaml.Unmarshal(data, &doc); decodeErr != nil {
		return nil, decodeErr
	}
	if len(doc.Content) == 0 {
		return root, nil
	}
	if doc.Content[0].Kind != yaml.MappingNode {
		return nil, setupErr.BackendRCMapping()
	}
	return doc.Content[0], nil
}

// defaults applies backend-specific option defaults.
//
// Parameters:
//   - options: backend setup options
//
// Returns:
//   - Options: setup options with defaults applied
func defaults(options Options) Options {
	if options.Endpoint == "" {
		options.Endpoint = defaultEndpoint(options.Name)
	}
	if options.APIKeyEnv == "" {
		options.APIKeyEnv = defaultAPIKeyEnv(options.Name)
	}
	return options
}

// validate reports whether backend setup options name a supported backend.
//
// Parameters:
//   - options: backend setup options
//
// Returns:
//   - error: unsupported backend error, or nil when supported
func validate(options Options) error {
	if defaultEndpoint(options.Name) == "" &&
		options.Name != cfgBackend.NameOpenAICompatible {
		return setupErr.UnsupportedBackend(options.Name)
	}
	if options.Name == cfgBackend.NameOpenAICompatible && options.Endpoint == "" {
		return setupErr.BackendEndpointRequired(options.Name)
	}
	if options.Endpoint != "" {
		parsed, parseErr := url.Parse(options.Endpoint)
		if parseErr != nil ||
			(parsed.Scheme != cfgHTTP.SchemeHTTP &&
				parsed.Scheme != cfgHTTP.SchemeHTTPS) {
			return setupErr.BackendEndpointScheme(options.Name)
		}
	}
	return nil
}

// defaultEndpoint returns the default endpoint for a backend.
//
// Parameters:
//   - name: backend name
//
// Returns:
//   - string: default endpoint, or empty when unknown
func defaultEndpoint(name string) string {
	switch name {
	case cfgBackend.NameVLLM:
		return cfgBackend.DefaultEndpointVLLM
	case cfgBackend.NameOpenAI:
		return cfgBackend.DefaultEndpointOpenAI
	case cfgBackend.NameAnthropic:
		return cfgBackend.DefaultEndpointAnthropic
	case cfgBackend.NameOllama:
		return cfgBackend.DefaultEndpointOllama
	case cfgBackend.NameLMStudio:
		return cfgBackend.DefaultEndpointLMStudio
	default:
		return ""
	}
}

// defaultAPIKeyEnv returns the default credential env var for a backend.
//
// Parameters:
//   - name: backend name
//
// Returns:
//   - string: default env var, or empty when none
func defaultAPIKeyEnv(name string) string {
	switch name {
	case cfgBackend.NameOpenAI:
		return cfgBackend.DefaultAPIKeyEnvOpenAI
	case cfgBackend.NameAnthropic:
		return cfgBackend.DefaultAPIKeyEnvAnthropic
	default:
		return ""
	}
}

// downstreamEnv returns optional shell export hints for downstream tools.
//
// Parameters:
//   - name: backend name
//   - endpoint: configured backend endpoint
//
// Returns:
//   - string: shell export line, or empty when no hint applies
func downstreamEnv(name, endpoint string) string {
	if endpoint == "" {
		return ""
	}
	var envName string
	switch name {
	case cfgBackend.NameOpenAI:
		envName = cfgSetup.BackendEnvOpenAIBaseURL
	case cfgBackend.NameAnthropic:
		envName = cfgSetup.BackendEnvAnthropicBaseURL
	default:
		return ""
	}
	return fmt.Sprintf(cfgSetup.BackendShellExportLine, envName, endpoint)
}
