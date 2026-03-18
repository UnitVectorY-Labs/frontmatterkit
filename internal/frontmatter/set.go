package frontmatter

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

// SetMode controls how values are merged.
type SetMode int

const (
	// SetOverwrite replaces the value at the path entirely.
	SetOverwrite SetMode = iota
	// SetPatch merges maps recursively; arrays and scalars are replaced.
	SetPatch
)

// Set sets a value at the given path in the document's front matter.
// If no front matter exists, it creates one.
// The yamlValue parameter is a YAML string that will be parsed.
func (d *Document) Set(path pathsyntax.Path, yamlValue string, mode SetMode) error {
	// Parse the value
	var valDoc yaml.Node
	if err := yaml.Unmarshal([]byte(yamlValue), &valDoc); err != nil {
		return fmt.Errorf("invalid YAML value: %w", err)
	}
	if len(valDoc.Content) == 0 {
		return fmt.Errorf("empty YAML value")
	}
	newVal := valDoc.Content[0]

	// Ensure front matter exists
	if !d.HasFrontMatter || d.Node == nil {
		d.HasFrontMatter = true
		d.Node = &yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: "!!map"},
			},
		}
	}

	root := d.Node
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			root.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		}
		root = root.Content[0]
	}

	if path.IsRoot {
		// Replace the entire front matter
		if mode == SetPatch && root.Kind == yaml.MappingNode && newVal.Kind == yaml.MappingNode {
			mergeMappingNodes(root, newVal)
		} else {
			d.Node.Content[0] = newVal
		}
		return nil
	}

	// Navigate/create intermediate nodes
	current := root
	segments := path.Segments
	for i, seg := range segments {
		isLast := i == len(segments)-1

		switch seg.Type {
		case pathsyntax.SegmentField:
			if current.Kind != yaml.MappingNode {
				return fmt.Errorf("cannot set field %q on non-map node", seg.Field)
			}

			existing := findMapValueIndex(current, seg.Field)
			if isLast {
				if existing >= 0 {
					valNode := current.Content[existing+1]
					if mode == SetPatch && valNode.Kind == yaml.MappingNode && newVal.Kind == yaml.MappingNode {
						mergeMappingNodes(valNode, newVal)
					} else {
						current.Content[existing+1] = newVal
					}
				} else {
					// Add new key-value pair
					keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: seg.Field, Tag: "!!str"}
					current.Content = append(current.Content, keyNode, newVal)
				}
			} else {
				if existing >= 0 {
					current = current.Content[existing+1]
				} else {
					// Create intermediate map
					newMap := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
					keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: seg.Field, Tag: "!!str"}
					current.Content = append(current.Content, keyNode, newMap)
					current = newMap
				}
			}

		case pathsyntax.SegmentIndex:
			if current.Kind != yaml.SequenceNode {
				return fmt.Errorf("cannot index into non-sequence node")
			}
			if seg.Index < 0 || seg.Index >= len(current.Content) {
				return fmt.Errorf("index %d out of range (length %d)", seg.Index, len(current.Content))
			}
			if isLast {
				if mode == SetPatch && current.Content[seg.Index].Kind == yaml.MappingNode && newVal.Kind == yaml.MappingNode {
					mergeMappingNodes(current.Content[seg.Index], newVal)
				} else {
					current.Content[seg.Index] = newVal
				}
			} else {
				current = current.Content[seg.Index]
			}
		}
	}

	return nil
}

// findMapValueIndex returns the index of the key node in a MappingNode, or -1 if not found.
func findMapValueIndex(node *yaml.Node, key string) int {
	if node.Kind != yaml.MappingNode {
		return -1
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return i
		}
	}
	return -1
}

// mergeMappingNodes merges src into dst recursively.
// New keys from src are added to dst; existing keys are updated.
func mergeMappingNodes(dst, src *yaml.Node) {
	for i := 0; i < len(src.Content)-1; i += 2 {
		srcKey := src.Content[i]
		srcVal := src.Content[i+1]

		existingIdx := findMapValueIndex(dst, srcKey.Value)
		if existingIdx >= 0 {
			dstVal := dst.Content[existingIdx+1]
			if dstVal.Kind == yaml.MappingNode && srcVal.Kind == yaml.MappingNode {
				mergeMappingNodes(dstVal, srcVal)
			} else {
				dst.Content[existingIdx+1] = srcVal
			}
		} else {
			dst.Content = append(dst.Content, srcKey, srcVal)
		}
	}
}
