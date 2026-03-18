package frontmatter

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse splits a Markdown document into front matter and body.
// If the file doesn't start with ---, HasFrontMatter is false and Body is the full content.
// If it starts with --- but YAML is invalid, returns an error.
func Parse(content string) (*Document, error) {
	doc := &Document{OriginalContent: content}

	if !strings.HasPrefix(content, "---") {
		doc.Body = content
		return doc, nil
	}

	// Find the closing ---
	// The opening delimiter is "---" possibly followed by a newline.
	rest := content[3:]
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 0 && rest[0] == '\r' {
		rest = rest[1:]
		if len(rest) > 0 && rest[0] == '\n' {
			rest = rest[1:]
		}
	}

	// Search for closing "---" on its own line
	closingIdx := -1
	lines := strings.SplitAfter(rest, "\n")
	pos := 0
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "---" {
			closingIdx = pos
			break
		}
		pos += len(line)
	}

	if closingIdx == -1 {
		return nil, fmt.Errorf("unclosed front matter: file starts with '---' but no closing '---' found")
	}

	doc.HasFrontMatter = true
	doc.RawFrontMatter = rest[:closingIdx]

	// Body starts after the closing --- and its newline
	afterClose := rest[closingIdx+3:]
	if len(afterClose) > 0 && afterClose[0] == '\n' {
		afterClose = afterClose[1:]
	} else if len(afterClose) > 0 && afterClose[0] == '\r' {
		afterClose = afterClose[1:]
		if len(afterClose) > 0 && afterClose[0] == '\n' {
			afterClose = afterClose[1:]
		}
	}
	doc.Body = afterClose

	// Parse the YAML into a node tree
	if strings.TrimSpace(doc.RawFrontMatter) == "" {
		// Empty front matter is valid
		doc.Node = &yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{Kind: yaml.MappingNode, Tag: "!!map"},
			},
		}
		return doc, nil
	}

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(doc.RawFrontMatter), &node); err != nil {
		return nil, fmt.Errorf("invalid YAML in front matter: %w", err)
	}
	doc.Node = &node

	return doc, nil
}

// allowedTags contains the standard YAML tags that yaml.v3 sets via implicit type resolution.
var allowedTags = map[string]bool{
	"!!str":       true,
	"!!int":       true,
	"!!float":     true,
	"!!bool":      true,
	"!!null":      true,
	"!!seq":       true,
	"!!map":       true,
	"!!binary":    true,
	"!!timestamp": true,
	"!!merge":     true,
}

// Validate checks the document against strict YAML rules.
// Returns nil if valid, error describing the problem otherwise.
// Missing front matter is valid.
func Validate(doc *Document) error {
	if !doc.HasFrontMatter {
		return nil
	}
	if doc.Node == nil {
		return nil
	}
	return validateNode(doc.Node)
}

func validateNode(node *yaml.Node) error {
	if node.Anchor != "" {
		return fmt.Errorf("anchors are not allowed (found anchor %q)", node.Anchor)
	}
	if node.Kind == yaml.AliasNode {
		return fmt.Errorf("aliases are not allowed")
	}
	if node.Tag != "" && !allowedTags[node.Tag] {
		return fmt.Errorf("custom tag %q is not allowed", node.Tag)
	}

	// Check for duplicate keys in mapping nodes
	if node.Kind == yaml.MappingNode {
		seen := make(map[string]bool)
		for i := 0; i < len(node.Content)-1; i += 2 {
			key := node.Content[i]
			keyStr := key.Value
			if seen[keyStr] {
				return fmt.Errorf("duplicate key %q", keyStr)
			}
			seen[keyStr] = true
		}
	}

	for _, child := range node.Content {
		if err := validateNode(child); err != nil {
			return err
		}
	}
	return nil
}

// ParseAndValidate combines Parse and Validate.
func ParseAndValidate(content string) (*Document, error) {
	doc, err := Parse(content)
	if err != nil {
		return nil, err
	}
	if err := Validate(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// Render produces the complete markdown content from the document.
func (d *Document) Render() string {
	if !d.HasFrontMatter {
		return d.Body
	}

	if d.Node == nil {
		return "---\n---\n" + d.Body
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	// Encode the node. If it's a DocumentNode, encode its content.
	var toEncode *yaml.Node
	if d.Node.Kind == yaml.DocumentNode && len(d.Node.Content) > 0 {
		toEncode = d.Node
	} else {
		toEncode = &yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{d.Node},
		}
	}

	if err := enc.Encode(toEncode); err != nil {
		// Fallback: use raw front matter
		return "---\n" + d.RawFrontMatter + "---\n" + d.Body
	}
	enc.Close()

	yamlStr := buf.String()
	// yaml.v3 encoder may add a trailing "...\n" or newline; clean up
	yamlStr = strings.TrimSuffix(yamlStr, "\n")
	yamlStr = strings.TrimSuffix(yamlStr, "...")
	yamlStr = strings.TrimSuffix(yamlStr, "\n")

	// Check if the mapping is empty (yaml encodes empty map as "{}")
	if yamlStr == "{}" {
		return "---\n---\n" + d.Body
	}

	return "---\n" + yamlStr + "\n---\n" + d.Body
}
