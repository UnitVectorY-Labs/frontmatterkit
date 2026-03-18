package frontmatter

import (
	"gopkg.in/yaml.v3"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

// Unset removes a value at the given path.
// Returns true if something was removed, false if path not found.
func (d *Document) Unset(path pathsyntax.Path) bool {
	if !d.HasFrontMatter || d.Node == nil {
		return false
	}
	if path.IsRoot {
		// Clear all front matter content
		root := d.Node
		if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
			root.Content[0] = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			return true
		}
		return false
	}

	root := d.Node
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			return false
		}
		root = root.Content[0]
	}

	// Navigate to the parent of the target
	segments := path.Segments
	current := root
	for i := 0; i < len(segments)-1; i++ {
		seg := segments[i]
		switch seg.Type {
		case pathsyntax.SegmentField:
			current = findMapValue(current, seg.Field)
		case pathsyntax.SegmentIndex:
			current = findSeqElement(current, seg.Index)
		}
		if current == nil {
			return false
		}
	}

	last := segments[len(segments)-1]
	switch last.Type {
	case pathsyntax.SegmentField:
		return removeMapKey(current, last.Field)
	case pathsyntax.SegmentIndex:
		return removeSeqElement(current, last.Index)
	}
	return false
}

// removeMapKey removes a key-value pair from a MappingNode.
func removeMapKey(node *yaml.Node, key string) bool {
	if node.Kind != yaml.MappingNode {
		return false
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			return true
		}
	}
	return false
}

// removeSeqElement removes an element from a SequenceNode by index.
func removeSeqElement(node *yaml.Node, index int) bool {
	if node.Kind != yaml.SequenceNode {
		return false
	}
	if index < 0 || index >= len(node.Content) {
		return false
	}
	node.Content = append(node.Content[:index], node.Content[index+1:]...)
	return true
}
