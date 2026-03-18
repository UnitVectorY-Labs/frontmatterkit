package assertion

import (
	"testing"
)

func TestParseAssertion_Exists(t *testing.T) {
	a, err := ParseAssertion(".title exists")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".title" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".title")
	}
	if a.Operator != OpExists {
		t.Errorf("operator = %v, want OpExists", a.Operator)
	}
	if a.Value != "" {
		t.Errorf("value = %q, want empty", a.Value)
	}
}

func TestParseAssertion_NotExists(t *testing.T) {
	a, err := ParseAssertion(".title not exists")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".title" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".title")
	}
	if a.Operator != OpNotExists {
		t.Errorf("operator = %v, want OpNotExists", a.Operator)
	}
	if a.Value != "" {
		t.Errorf("value = %q, want empty", a.Value)
	}
}

func TestParseAssertion_GreaterEqual(t *testing.T) {
	a, err := ParseAssertion(".count >= 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".count" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".count")
	}
	if a.Operator != OpGreaterEqual {
		t.Errorf("operator = %v, want OpGreaterEqual", a.Operator)
	}
	if a.Value != "3" {
		t.Errorf("value = %q, want %q", a.Value, "3")
	}
}

func TestParseAssertion_Contains(t *testing.T) {
	a, err := ParseAssertion(`.tags contains "go"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".tags" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".tags")
	}
	if a.Operator != OpContains {
		t.Errorf("operator = %v, want OpContains", a.Operator)
	}
	if a.Value != `"go"` {
		t.Errorf("value = %q, want %q", a.Value, `"go"`)
	}
}

func TestParseAssertion_NotContains(t *testing.T) {
	a, err := ParseAssertion(`.tags not contains "draft"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".tags" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".tags")
	}
	if a.Operator != OpNotContains {
		t.Errorf("operator = %v, want OpNotContains", a.Operator)
	}
	if a.Value != `"draft"` {
		t.Errorf("value = %q, want %q", a.Value, `"draft"`)
	}
}

func TestParseAssertion_Equals(t *testing.T) {
	a, err := ParseAssertion(".draft == false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".draft" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".draft")
	}
	if a.Operator != OpEquals {
		t.Errorf("operator = %v, want OpEquals", a.Operator)
	}
	if a.Value != "false" {
		t.Errorf("value = %q, want %q", a.Value, "false")
	}
}

func TestParseAssertion_NotEquals(t *testing.T) {
	a, err := ParseAssertion(".draft != true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".draft" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".draft")
	}
	if a.Operator != OpNotEquals {
		t.Errorf("operator = %v, want OpNotEquals", a.Operator)
	}
	if a.Value != "true" {
		t.Errorf("value = %q, want %q", a.Value, "true")
	}
}

func TestParseAssertion_LessThan(t *testing.T) {
	a, err := ParseAssertion(".count < 50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Operator != OpLessThan {
		t.Errorf("operator = %v, want OpLessThan", a.Operator)
	}
	if a.Value != "50" {
		t.Errorf("value = %q, want %q", a.Value, "50")
	}
}

func TestParseAssertion_LessEqual(t *testing.T) {
	a, err := ParseAssertion(".count <= 42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Operator != OpLessEqual {
		t.Errorf("operator = %v, want OpLessEqual", a.Operator)
	}
	if a.Value != "42" {
		t.Errorf("value = %q, want %q", a.Value, "42")
	}
}

func TestParseAssertion_GreaterThan(t *testing.T) {
	a, err := ParseAssertion(".count > 10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Operator != OpGreaterThan {
		t.Errorf("operator = %v, want OpGreaterThan", a.Operator)
	}
	if a.Value != "10" {
		t.Errorf("value = %q, want %q", a.Value, "10")
	}
}

func TestParseAssertion_ErrorCases(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"empty expression", ""},
		{"missing operator", ".title"},
		{"unknown operator", ".title unknown"},
		{"missing value for ==", ".title =="},
		{"missing value for !=", ".title !="},
		{"missing value for contains", ".title contains"},
		{"missing value for not contains", ".title not contains"},
		{"missing value for <", ".count <"},
		{"missing value for >=", ".count >="},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseAssertion(tc.expr)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tc.expr)
			}
		})
	}
}

func TestParseAssertion_IndexedPath(t *testing.T) {
	a, err := ParseAssertion(".tags[0] == \"go\"")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Path.String() != ".tags[0]" {
		t.Errorf("path = %q, want %q", a.Path.String(), ".tags[0]")
	}
	if a.Operator != OpEquals {
		t.Errorf("operator = %v, want OpEquals", a.Operator)
	}
	if a.Value != `"go"` {
		t.Errorf("value = %q, want %q", a.Value, `"go"`)
	}
}
