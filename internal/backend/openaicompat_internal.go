//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	cfgHTTP "github.com/ActiveMemory/ctx/internal/config/http"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// openAICompatibleFactory creates a generic OpenAI-compatible backend.
//
// Parameters:
//   - config: backend configuration
//
// Returns:
//   - Backend: configured backend
//   - error: always nil; request validation happens during calls
func openAICompatibleFactory(config Config) (Backend, error) {
	return openAICompatible{
		name:   cfgBackend.NameOpenAICompatible,
		config: config,
	}, nil
}

// models requests the model list and extracts the first model.
//
// Parameters:
//   - ctx: request context
//
// Returns:
//   - Response: model-list response with FirstModel set when present
//   - error: unreachable, upstream, or decode failure
func (backend openAICompatible) models(ctx context.Context) (Response, error) {
	raw, doErr := backend.do(
		ctx,
		cfgBackend.HTTPMethodGet,
		cfgBackend.ModelsPath,
		nil,
	)
	if doErr != nil {
		return Response{}, doErr
	}
	var decoded modelsResponse
	if decodeErr := json.Unmarshal(raw, &decoded); decodeErr != nil {
		return Response{}, errBackend.BadRequest{
			Name:  backend.name,
			Cause: decodeErr,
		}
	}
	firstModel := ""
	if len(decoded.Data) > 0 {
		firstModel = decoded.Data[0].ID
	}
	return Response{FirstModel: firstModel, Raw: raw}, nil
}

// do sends a backend HTTP request and returns the raw response body.
//
// Parameters:
//   - ctx: request context
//   - method: HTTP method constant
//   - path: endpoint path constant
//   - body: optional JSON request body
//
// Returns:
//   - []byte: raw response body
//   - error: endpoint, transport, read, or upstream status error
func (backend openAICompatible) do(
	ctx context.Context,
	method string,
	path string,
	body []byte,
) ([]byte, error) {
	endpoint, endpointErr := backend.url(path)
	if endpointErr != nil {
		return nil, endpointErr
	}
	requestBody := io.Reader(nil)
	if body != nil {
		requestBody = bytes.NewReader(body)
	}
	request, requestErr := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		requestBody,
	)
	if requestErr != nil {
		return nil, errBackend.InvalidEndpoint{
			Endpoint: backend.config.Endpoint,
			Cause:    requestErr,
		}
	}
	backend.headers(request, body != nil)
	client := http.Client{Timeout: backend.timeout()}
	response, doErr := client.Do(request)
	if doErr != nil {
		return nil, errBackend.Unreachable{
			Name:     backend.name,
			Endpoint: backend.config.Endpoint,
			Cause:    doErr,
		}
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			logWarn.Warn(cfgWarn.CloseResponse, closeErr)
		}
	}()
	raw, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, errBackend.BadRequest{Name: backend.name, Cause: readErr}
	}
	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, errBackend.Upstream{
			Name:       backend.name,
			StatusCode: response.StatusCode,
			Body:       string(raw),
		}
	}
	return raw, nil
}

// headers applies optional JSON and authorization headers.
//
// Parameters:
//   - request: outbound HTTP request
//   - hasBody: whether the request has a JSON body
func (backend openAICompatible) headers(
	request *http.Request,
	hasBody bool,
) {
	if hasBody {
		request.Header.Set(
			cfgBackend.HeaderContentType,
			cfgBackend.ContentTypeJSON,
		)
	}
	if backend.config.APIKeyEnv == "" {
		return
	}
	apiKey := os.Getenv(backend.config.APIKeyEnv)
	if apiKey == "" {
		return
	}
	request.Header.Set(
		cfgBackend.HeaderAuthorization,
		cfgBackend.AuthorizationBearerPrefix+apiKey,
	)
}

// timeout parses the configured timeout with a safe fallback.
//
// Returns:
//   - time.Duration: configured or default timeout
func (backend openAICompatible) timeout() time.Duration {
	parsed, parseErr := time.ParseDuration(backend.config.Timeout)
	if parseErr != nil || parsed <= 0 {
		return cfgBackend.DefaultTimeout
	}
	return parsed
}

// endpoint returns the fully defaulted backend endpoint.
//
// Returns:
//   - string: configured endpoint after factory defaults
func (backend openAICompatible) endpoint() string {
	return backend.config.Endpoint
}

// url combines the configured endpoint with an API path.
//
// Parameters:
//   - path: endpoint path constant
//
// Returns:
//   - string: absolute request URL
//   - error: invalid endpoint error
func (backend openAICompatible) url(path string) (string, error) {
	parsed, parseErr := url.Parse(backend.config.Endpoint)
	if parseErr != nil {
		return "", errBackend.InvalidEndpoint{
			Endpoint: backend.config.Endpoint,
			Cause:    parseErr,
		}
	}
	if parsed.Scheme != cfgHTTP.SchemeHTTP &&
		parsed.Scheme != cfgHTTP.SchemeHTTPS {
		return "", errBackend.InvalidEndpoint{
			Endpoint: backend.config.Endpoint,
		}
	}
	parsed.Path = path
	return parsed.String(), nil
}
