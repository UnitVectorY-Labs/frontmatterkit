package assertion

import (
	"fmt"
	"strings"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

// OpType represents the assertion operator.
type OpType int

const (
	OpExists OpType = iota
	OpNotExists
	OpEquals
	OpNotEquals
	OpContains
	OpNotContains
	OpLessThan
	OpLessEqual
	OpGreaterThan
	OpGreaterEqual
)

// Assertion represents a parsed assertion expression.
type Assertion struct {
	Path     pathsyntax.Path
	Operator OpType
	Value    string // raw YAML value string for comparison (empty for exists/not exists)
}

// ParseAssertion parses an assertion expression string.
// Format: "<path> <operator> [<value>]"
// Examples:
//
//	".title exists"
//	".count >= 3"
//	".tags contains \"go\""
//	".draft == false"
//	".title not exists"
func ParseAssertion(expr string) (*Assertion, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty assertion expression")
	}

	// Find the end of the path: first space not inside brackets
	pathEnd := findPathEnd(expr)
	if pathEnd <= 0 {
		return nil, fmt.Errorf("missing operator in assertion %q", expr)
	}

	pathStr := expr[:pathEnd]
	rest := strings.TrimSpace(expr[pathEnd:])

	path, err := pathsyntax.Parse(pathStr)
	if err != nil {
		return nil, fmt.Errorf("invalid path in assertion %q: %w", expr, err)
	}

	if rest == "" {
		return nil, fmt.Errorf("missing operator in assertion %q", expr)
	}

	// Parse operator and value
	op, value, err := parseOperatorAndValue(rest, expr)
	if err != nil {
		return nil, err
	}

	return &Assertion{
		Path:     path,
		Operator: op,
		Value:    value,
	}, nil
}

// findPathEnd returns the index of the first space not inside brackets.
func findPathEnd(expr string) int {
	depth := 0
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case '[':
			depth++
		case ']':
			depth--
		case ' ', '\t':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// parseOperatorAndValue extracts the operator and optional value from the remainder string.
func parseOperatorAndValue(rest, fullExpr string) (OpType, string, error) {
	// Check for two-word operators first: "not exists", "not contains"
	if strings.HasPrefix(rest, "not ") {
		after := strings.TrimSpace(rest[4:])
		if after == "exists" {
			return OpNotExists, "", nil
		}
		if strings.HasPrefix(after, "contains") {
			value := strings.TrimSpace(after[8:])
			if value == "" {
				return 0, "", fmt.Errorf("missing value for 'not contains' operator in assertion %q", fullExpr)
			}
			return OpNotContains, value, nil
		}
	}

	// Single-word or symbol operators
	if rest == "exists" {
		return OpExists, "", nil
	}

	// Symbol operators: ==, !=, <=, >=, <, >
	// Check two-char operators first
	for _, pair := range []struct {
		prefix string
		op     OpType
	}{
		{"==", OpEquals},
		{"!=", OpNotEquals},
		{"<=", OpLessEqual},
		{">=", OpGreaterEqual},
		{"<", OpLessThan},
		{">", OpGreaterThan},
	} {
		if strings.HasPrefix(rest, pair.prefix) {
			value := strings.TrimSpace(rest[len(pair.prefix):])
			if value == "" {
				return 0, "", fmt.Errorf("missing value for %q operator in assertion %q", pair.prefix, fullExpr)
			}
			return pair.op, value, nil
		}
	}

	if strings.HasPrefix(rest, "contains") {
		value := strings.TrimSpace(rest[8:])
		if value == "" {
			return 0, "", fmt.Errorf("missing value for 'contains' operator in assertion %q", fullExpr)
		}
		return OpContains, value, nil
	}

	// Extract first word for error message
	firstWord := rest
	if idx := strings.IndexByte(rest, ' '); idx != -1 {
		firstWord = rest[:idx]
	}
	return 0, "", fmt.Errorf("unknown operator %q in assertion %q", firstWord, fullExpr)
}
