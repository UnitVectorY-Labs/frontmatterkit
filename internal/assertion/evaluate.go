package assertion

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
)

// Evaluate checks the assertion against the document.
// Returns nil if the assertion passes, error with explanation if it fails.
func (a *Assertion) Evaluate(doc *frontmatter.Document) error {
	node := doc.Get(a.Path)

	switch a.Operator {
	case OpExists:
		if node == nil {
			return fmt.Errorf("assertion failed: %s does not exist", a.Path.String())
		}
		return nil

	case OpNotExists:
		if node != nil {
			return fmt.Errorf("assertion failed: %s exists", a.Path.String())
		}
		return nil

	case OpEquals:
		return evalEquals(a, node, false)

	case OpNotEquals:
		return evalEquals(a, node, true)

	case OpContains:
		return evalContains(a, node, false)

	case OpNotContains:
		return evalContains(a, node, true)

	case OpLessThan:
		return evalNumeric(a, node, func(actual, expected float64) bool { return actual < expected })

	case OpLessEqual:
		return evalNumeric(a, node, func(actual, expected float64) bool { return actual <= expected })

	case OpGreaterThan:
		return evalNumeric(a, node, func(actual, expected float64) bool { return actual > expected })

	case OpGreaterEqual:
		return evalNumeric(a, node, func(actual, expected float64) bool { return actual >= expected })

	default:
		return fmt.Errorf("unsupported operator")
	}
}

// evalEquals handles == and != operators.
func evalEquals(a *Assertion, node *yaml.Node, negate bool) error {
	// Decode the expected value via YAML
	var expected interface{}
	if err := yaml.Unmarshal([]byte(a.Value), &expected); err != nil {
		return fmt.Errorf("assertion failed: cannot parse expected value %q as YAML: %w", a.Value, err)
	}

	// Decode the actual value
	var actual interface{}
	if node != nil {
		if err := node.Decode(&actual); err != nil {
			return fmt.Errorf("assertion failed: cannot decode node value: %w", err)
		}
	}

	equal := valuesEqual(actual, expected)

	if negate {
		if equal {
			return fmt.Errorf("assertion failed: %s == %v, expected != %s", a.Path.String(), actual, a.Value)
		}
		return nil
	}

	if !equal {
		return fmt.Errorf("assertion failed: %s is %v, expected %s", a.Path.String(), actual, a.Value)
	}
	return nil
}

// valuesEqual compares two decoded YAML values.
func valuesEqual(a, b interface{}) bool {
	// Normalize numeric types for comparison: both int and float should compare equal
	af, aIsNum := toFloat64(a)
	bf, bIsNum := toFloat64(b)
	if aIsNum && bIsNum {
		return af == bf
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// evalContains handles contains and not contains operators.
func evalContains(a *Assertion, node *yaml.Node, negate bool) error {
	if node == nil {
		if negate {
			return nil
		}
		return fmt.Errorf("assertion failed: %s does not exist", a.Path.String())
	}

	var found bool

	switch node.Kind {
	case yaml.SequenceNode:
		// Parse the expected value from YAML
		var expected interface{}
		if err := yaml.Unmarshal([]byte(a.Value), &expected); err != nil {
			return fmt.Errorf("assertion failed: cannot parse expected value %q as YAML: %w", a.Value, err)
		}

		for _, elem := range node.Content {
			var elemVal interface{}
			if err := elem.Decode(&elemVal); err != nil {
				continue
			}
			if valuesEqual(elemVal, expected) {
				found = true
				break
			}
		}

	case yaml.ScalarNode:
		// For scalar strings, check substring containment
		// Parse the expected value to get a string
		var expected interface{}
		if err := yaml.Unmarshal([]byte(a.Value), &expected); err != nil {
			return fmt.Errorf("assertion failed: cannot parse expected value %q as YAML: %w", a.Value, err)
		}
		expectedStr := fmt.Sprintf("%v", expected)
		found = strings.Contains(node.Value, expectedStr)

	default:
		return fmt.Errorf("assertion failed: %s is not a sequence or string, cannot use contains", a.Path.String())
	}

	if negate {
		if found {
			return fmt.Errorf("assertion failed: %s contains %s", a.Path.String(), a.Value)
		}
		return nil
	}

	if !found {
		return fmt.Errorf("assertion failed: %s does not contain %s", a.Path.String(), a.Value)
	}
	return nil
}

// evalNumeric handles <, <=, >, >= operators.
func evalNumeric(a *Assertion, node *yaml.Node, cmp func(actual, expected float64) bool) error {
	if node == nil {
		return fmt.Errorf("assertion failed: %s does not exist, cannot compare numerically", a.Path.String())
	}

	actualF, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return fmt.Errorf("assertion failed: %s value %q is not numeric", a.Path.String(), node.Value)
	}

	expectedF, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return fmt.Errorf("assertion failed: expected value %q is not numeric", a.Value)
	}

	if !cmp(actualF, expectedF) {
		return fmt.Errorf("assertion failed: %s is %v, expected %s %s", a.Path.String(), actualF, operatorSymbol(a.Operator), a.Value)
	}
	return nil
}

// toFloat64 tries to convert a value to float64.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	case float32:
		return float64(n), true
	default:
		return 0, false
	}
}

// operatorSymbol returns the string representation of a comparison operator.
func operatorSymbol(op OpType) string {
	switch op {
	case OpLessThan:
		return "<"
	case OpLessEqual:
		return "<="
	case OpGreaterThan:
		return ">"
	case OpGreaterEqual:
		return ">="
	default:
		return "?"
	}
}
