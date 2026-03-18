package pathsyntax

import (
	"fmt"
	"strconv"
	"strings"
)

// SegmentType represents the kind of path segment.
type SegmentType int

const (
	SegmentField SegmentType = iota
	SegmentIndex
)

// Segment is one step in a path.
type Segment struct {
	Type  SegmentType
	Field string // for SegmentField
	Index int    // for SegmentIndex
}

// Path is a parsed jq-like path.
type Path struct {
	Segments []Segment
	IsRoot   bool // true when path is just "."
}

// Parse parses a jq-like path string into a Path.
// Returns error for invalid syntax.
func Parse(s string) (Path, error) {
	if s == "" {
		return Path{}, fmt.Errorf("empty path")
	}
	if s[0] != '.' {
		return Path{}, fmt.Errorf("path must start with '.'")
	}

	// Root selection
	if s == "." {
		return Path{IsRoot: true}, nil
	}

	rest := s[1:] // consume leading dot
	var segments []Segment

	for len(rest) > 0 {
		// Expect a field name at this point
		if rest[0] == '.' || rest[0] == '[' {
			return Path{}, fmt.Errorf("unexpected %q in path %q", string(rest[0]), s)
		}

		name, tail := consumeFieldName(rest)
		if name == "" {
			return Path{}, fmt.Errorf("empty field name in path %q", s)
		}
		if !isValidFieldName(name) {
			return Path{}, fmt.Errorf("invalid field name %q in path %q", name, s)
		}
		segments = append(segments, Segment{Type: SegmentField, Field: name})
		rest = tail

		// Parse any trailing array indices (e.g. [0][1])
		for len(rest) > 0 && rest[0] == '[' {
			idx, tail2, err := consumeIndex(rest, s)
			if err != nil {
				return Path{}, err
			}
			segments = append(segments, Segment{Type: SegmentIndex, Index: idx})
			rest = tail2
		}

		// Consume dot separator before next field, if present
		if len(rest) > 0 {
			if rest[0] != '.' {
				return Path{}, fmt.Errorf("expected '.' or '[' but got %q in path %q", string(rest[0]), s)
			}
			rest = rest[1:] // consume the dot
			if len(rest) == 0 {
				return Path{}, fmt.Errorf("trailing '.' in path %q", s)
			}
		}
	}

	return Path{Segments: segments}, nil
}

// consumeFieldName reads a field name from the front of s, stopping at '.', '[', or end.
func consumeFieldName(s string) (name, rest string) {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' || s[i] == '[' {
			return s[:i], s[i:]
		}
	}
	return s, ""
}

// consumeIndex parses "[<int>]" from the front of s.
func consumeIndex(s, fullPath string) (int, string, error) {
	// s starts with '['
	end := strings.IndexByte(s, ']')
	if end == -1 {
		return 0, "", fmt.Errorf("unclosed '[' in path %q", fullPath)
	}
	inner := s[1:end]
	if inner == "" {
		return 0, "", fmt.Errorf("empty index in path %q", fullPath)
	}
	// Reject leading zeros (except "0" itself), signs, etc.
	if len(inner) > 1 && inner[0] == '0' {
		return 0, "", fmt.Errorf("invalid index %q in path %q", inner, fullPath)
	}
	if inner[0] == '-' {
		return 0, "", fmt.Errorf("negative index in path %q", fullPath)
	}
	idx, err := strconv.Atoi(inner)
	if err != nil {
		return 0, "", fmt.Errorf("invalid index %q in path %q", inner, fullPath)
	}
	return idx, s[end+1:], nil
}

// isValidFieldName checks that a field name only contains letters, digits, underscores, and hyphens.
func isValidFieldName(name string) bool {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return len(name) > 0
}

// String returns the canonical string representation of the path.
func (p Path) String() string {
	if p.IsRoot {
		return "."
	}
	var b strings.Builder
	for _, seg := range p.Segments {
		switch seg.Type {
		case SegmentField:
			b.WriteByte('.')
			b.WriteString(seg.Field)
		case SegmentIndex:
			b.WriteByte('[')
			b.WriteString(strconv.Itoa(seg.Index))
			b.WriteByte(']')
		}
	}
	return b.String()
}
