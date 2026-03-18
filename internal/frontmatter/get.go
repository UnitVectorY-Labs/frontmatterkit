package frontmatter

import (
	"gopkg.in/yaml.v3"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

// Get retrieves the yaml.Node at the given path.
// Returns nil if not found or if there is no front matter.
func (d *Document) Get(path pathsyntax.Path) *yaml.Node {
	if !d.HasFrontMatter || d.Node == nil {
		return nil
	}

	root := d.Node
	// Unwrap DocumentNode
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			return nil
		}
		root = root.Content[0]
	}

	if path.IsRoot {
		return root
	}

	current := root
	for _, seg := range path.Segments {
		if current == nil {
			return nil
		}
		switch seg.Type {
		case pathsyntax.SegmentField:
			current = findMapValue(current, seg.Field)
		case pathsyntax.SegmentIndex:
			current = findSeqElement(current, seg.Index)
		}
	}
	return current
}

// GetValue retrieves the Go value at the given path.
// Returns nil if path not found.
func (d *Document) GetValue(path pathsyntax.Path) interface{} {
	node := d.Get(path)
	if node == nil {
		return nil
	}
	var val interface{}
	if err := node.Decode(&val); err != nil {
		return nil
	}
	return val
}

// findMapValue finds the value node for a key in a MappingNode.
func findMapValue(node *yaml.Node, key string) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// findSeqElement returns the element at index in a SequenceNode.
func findSeqElement(node *yaml.Node, index int) *yaml.Node {
	if node.Kind != yaml.SequenceNode {
		return nil
	}
	if index < 0 || index >= len(node.Content) {
		return nil
	}
	return node.Content[index]
}
