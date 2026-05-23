//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"time"

	"gopkg.in/yaml.v3"

	cfgRc "github.com/ActiveMemory/ctx/internal/config/rc"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// parseDocument decodes raw .ctxrc bytes into a yaml.Node
// tree. Empty input yields a fresh document containing
// a single empty mapping so Apply can append the first
// `backends:` key into it.
//
// Parameters:
//   - data: raw .ctxrc bytes (may be empty).
//
// Returns:
//   - *yaml.Node: a DocumentNode with one mapping child.
//   - error: yaml.Unmarshal failure on malformed input.
func parseDocument(data []byte) (*yaml.Node, error) {
	doc := &yaml.Node{}
	if len(data) == 0 {
		doc.Kind = yaml.DocumentNode
		doc.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
		return doc, nil
	}
	if err := yaml.Unmarshal(data, doc); err != nil {
		return nil, err
	}
	if doc.Kind == 0 || len(doc.Content) == 0 {
		doc.Kind = yaml.DocumentNode
		doc.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}
	if doc.Content[0].Kind != yaml.MappingNode {
		doc.Content[0] = &yaml.Node{Kind: yaml.MappingNode}
	}
	return doc, nil
}

// upsert mutates doc's root mapping so the `backends:`
// list contains an entry matching entry.Name. Returns
// true when an existing entry was replaced; false when a
// new entry was appended (or when the list itself was
// created).
//
// Parameters:
//   - doc: DocumentNode produced by [parseDocument].
//   - entry: the backend config to add or update.
//
// Returns:
//   - bool: true if updated in place, false if appended.
//   - error: yaml.Marshal failure when constructing the
//     entry sub-document (rare).
func upsert(doc *yaml.Node, entry entity.BackendConfig) (bool, error) {
	root := doc.Content[0]
	backendsNode := findChildValue(root, cfgRc.YAMLKeyBackends)
	entryNode, entryErr := buildEntryNode(entry)
	if entryErr != nil {
		return false, entryErr
	}
	if backendsNode == nil {
		appendMappingPair(root, cfgRc.YAMLKeyBackends, &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: []*yaml.Node{entryNode},
		})
		return false, nil
	}
	if backendsNode.Kind != yaml.SequenceNode {
		// Hostile shape (e.g. user set backends: null);
		// reset to a fresh sequence with our entry.
		*backendsNode = yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: []*yaml.Node{entryNode},
		}
		return false, nil
	}
	for i, item := range backendsNode.Content {
		if mapValueString(item, cfgRc.YAMLKeyBackendName) == entry.Name {
			backendsNode.Content[i] = entryNode
			return true, nil
		}
	}
	backendsNode.Content = append(backendsNode.Content, entryNode)
	return false, nil
}

// findChildValue returns the value node for the given key
// in a MappingNode, or nil when absent.
//
// Parameters:
//   - mapping: a MappingNode (or any node — non-mappings
//     return nil).
//   - key: the scalar key to look up.
//
// Returns:
//   - *yaml.Node: the value node, or nil.
func findChildValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// mapValueString returns the scalar string under `key`
// in a MappingNode, or "" when the key is absent or the
// value is not a scalar.
//
// Parameters:
//   - mapping: candidate MappingNode.
//   - key: scalar key to look up.
//
// Returns:
//   - string: the scalar value or "".
func mapValueString(mapping *yaml.Node, key string) string {
	v := findChildValue(mapping, key)
	if v == nil || v.Kind != yaml.ScalarNode {
		return ""
	}
	return v.Value
}

// appendMappingPair appends a key/value pair to a
// MappingNode in declaration order.
//
// Parameters:
//   - mapping: the MappingNode to extend.
//   - key: the scalar key.
//   - value: the value node (any kind).
func appendMappingPair(mapping *yaml.Node, key string, value *yaml.Node) {
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		value,
	)
}

// buildEntryNode constructs a MappingNode for a single
// `backends:` list entry from an entity.BackendConfig.
// Done by marshaling a typed Go struct and re-parsing as
// a Node so the resulting node carries proper omitempty
// behavior for free.
//
// Parameters:
//   - entry: the backend config to render.
//
// Returns:
//   - *yaml.Node: the MappingNode (already unwrapped from
//     the document/sequence wrapping yaml.Unmarshal adds).
//   - error: yaml.Marshal failure (rare).
func buildEntryNode(entry entity.BackendConfig) (*yaml.Node, error) {
	be := backendEntry{
		Name:         entry.Name,
		Endpoint:     entry.Endpoint,
		APIKeyEnv:    entry.APIKeyEnv,
		Timeout:      durationString(entry.Timeout),
		DefaultModel: entry.DefaultModel,
	}
	raw, marshalErr := yaml.Marshal(be)
	if marshalErr != nil {
		return nil, marshalErr
	}
	var wrapper yaml.Node
	if unmarshalErr := yaml.Unmarshal(raw, &wrapper); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	// wrapper is a DocumentNode; its sole child is the
	// MappingNode for the entry.
	return wrapper.Content[0], nil
}

// durationString returns a Go-format duration string for
// non-zero d, or "" when d is zero (so `omitempty` keeps
// the field out of the document).
//
// Parameters:
//   - d: the duration.
//
// Returns:
//   - string: the duration in Go's stringer format, or "".
func durationString(d time.Duration) string {
	if d == 0 {
		return ""
	}
	return d.String()
}
