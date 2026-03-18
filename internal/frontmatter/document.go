package frontmatter

import "gopkg.in/yaml.v3"

// Document represents a parsed Markdown file with optional YAML front matter.
type Document struct {
	// HasFrontMatter indicates whether front matter was found.
	HasFrontMatter bool
	// RawFrontMatter is the raw YAML text between the --- delimiters.
	RawFrontMatter string
	// Body is the Markdown content after the front matter.
	Body string
	// Node is the parsed yaml.Node tree (for minimal-change editing).
	Node *yaml.Node
	// OriginalContent is the original full file content.
	OriginalContent string
}
