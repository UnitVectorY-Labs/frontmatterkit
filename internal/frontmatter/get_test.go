package frontmatter

import (
	"testing"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

func mustParsePath(t *testing.T, s string) pathsyntax.Path {
	t.Helper()
	p, err := pathsyntax.Parse(s)
	if err != nil {
		t.Fatalf("failed to parse path %q: %v", s, err)
	}
	return p
}

func TestGet_Root(t *testing.T) {
	content := "---\ntitle: Hello\nauthor: Test\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(mustParsePath(t, "."))
	if node == nil {
		t.Fatal("expected non-nil root node")
	}
}

func TestGet_ScalarField(t *testing.T) {
	content := "---\ntitle: Hello World\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(mustParsePath(t, ".title"))
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Value != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", node.Value)
	}
}

func TestGet_NestedField(t *testing.T) {
	content := "---\nauthor:\n  name: John\n  email: john@test.com\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(mustParsePath(t, ".author.name"))
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Value != "John" {
		t.Errorf("expected 'John', got %q", node.Value)
	}
}

func TestGet_ArrayElement(t *testing.T) {
	content := "---\ntags:\n  - go\n  - yaml\n  - test\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(mustParsePath(t, ".tags[0]"))
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Value != "go" {
		t.Errorf("expected 'go', got %q", node.Value)
	}

	node = doc.Get(mustParsePath(t, ".tags[2]"))
	if node == nil {
		t.Fatal("expected non-nil node for index 2")
	}
	if node.Value != "test" {
		t.Errorf("expected 'test', got %q", node.Value)
	}
}

func TestGet_MissingPath(t *testing.T) {
	content := "---\ntitle: Hello\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(mustParsePath(t, ".nonexistent"))
	if node != nil {
		t.Error("expected nil for missing path")
	}
}

func TestGet_NoFrontMatter(t *testing.T) {
	content := "# Hello\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(mustParsePath(t, ".title"))
	if node != nil {
		t.Error("expected nil when no front matter")
	}
}

func TestGetValue_Scalar(t *testing.T) {
	content := "---\ntitle: Hello World\ncount: 42\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val := doc.GetValue(mustParsePath(t, ".title"))
	if val != "Hello World" {
		t.Errorf("expected 'Hello World', got %v", val)
	}

	val = doc.GetValue(mustParsePath(t, ".count"))
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestGetValue_Missing(t *testing.T) {
	content := "---\ntitle: Hello\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val := doc.GetValue(mustParsePath(t, ".missing"))
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}
