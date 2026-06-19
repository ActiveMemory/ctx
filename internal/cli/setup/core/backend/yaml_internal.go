//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import "gopkg.in/yaml.v3"

// mapping returns a child mapping node, creating it when absent.
//
// Parameters:
//   - parent: mapping node to search
//   - key: child key
//
// Returns:
//   - *yaml.Node: child mapping node
func mapping(parent *yaml.Node, key string) *yaml.Node {
	for idx := 0; idx < len(parent.Content); idx += 2 {
		if parent.Content[idx].Value == key {
			parent.Content[idx+1].Kind = yaml.MappingNode
			return parent.Content[idx+1]
		}
	}
	child := &yaml.Node{Kind: yaml.MappingNode}
	parent.Content = append(parent.Content, scalarNode(key), child)
	return child
}

// scalar returns a child scalar node by key.
//
// Parameters:
//   - parent: mapping node to search
//   - key: child key
//
// Returns:
//   - *yaml.Node: child scalar node, or nil
func scalar(parent *yaml.Node, key string) *yaml.Node {
	for idx := 0; idx < len(parent.Content); idx += 2 {
		if parent.Content[idx].Value == key {
			return parent.Content[idx+1]
		}
	}
	return nil
}

// setScalar sets a child scalar node by key.
//
// Parameters:
//   - parent: mapping node to mutate
//   - key: child key
//   - value: scalar value
func setScalar(parent *yaml.Node, key string, value string) {
	if existing := scalar(parent, key); existing != nil {
		existing.Kind = yaml.ScalarNode
		existing.Value = value
		return
	}
	parent.Content = append(parent.Content, scalarNode(key), scalarNode(value))
}

// scalarNode creates a YAML scalar node.
//
// Parameters:
//   - value: scalar value
//
// Returns:
//   - *yaml.Node: scalar node
func scalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value}
}
