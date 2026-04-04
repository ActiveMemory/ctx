//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"bytes"

	"gopkg.in/yaml.v3"

	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
)

// defaultInclusion is the default inclusion mode when omitted.
var defaultInclusion = cfgSteering.InclusionManual

// Parse reads a steering file from bytes, extracting YAML frontmatter
// and markdown body. The filePath is stored on the returned SteeringFile
// for error reporting and identification.
//
// Frontmatter must be delimited by --- lines at the top of the file.
// Missing optional fields receive defaults: inclusion → manual,
// tools → nil (all), priority → 50.
//
// Returns an error if frontmatter contains invalid YAML, identifying
// the file path and the parsing failure.
func Parse(data []byte, filePath string) (*SteeringFile, error) {
	raw, body, splitErr := splitFrontmatter(data)
	if splitErr != nil {
		return nil, errSteering.Parse(filePath, splitErr)
	}

	sf := &SteeringFile{}
	if unmarshalErr := yaml.Unmarshal(raw, sf); unmarshalErr != nil {
		return nil, errSteering.InvalidYAML(filePath, unmarshalErr)
	}

	applyDefaults(sf)
	sf.Body = body
	sf.Path = filePath

	return sf, nil
}

// Print serializes a SteeringFile back to frontmatter + markdown bytes.
//
// The output format is:
//
//	---
//	<yaml frontmatter>
//	---
//	<markdown body>
//
// Round-trip property: Parse(Print(Parse(data))) == Parse(data) for all
// valid inputs.
func Print(sf *SteeringFile) []byte {
	var buf bytes.Buffer

	raw, _ := yaml.Marshal(sf)

	buf.WriteString(token.FrontmatterDelimiter)
	buf.WriteByte(token.NewlineLF[0])
	buf.Write(raw)
	buf.WriteString(token.FrontmatterDelimiter)
	buf.WriteByte(token.NewlineLF[0])

	if sf.Body != "" {
		buf.WriteString(sf.Body)
	}

	return buf.Bytes()
}
