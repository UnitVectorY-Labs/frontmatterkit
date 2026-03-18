package frontmatter

import (
	"testing"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

func TestUnset_ExistingKey(t *testing.T) {
	content := "---\ntitle: Hello\nauthor: Test\n---\nBody\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".title")
	removed := doc.Unset(path)
	if !removed {
		t.Error("expected Unset to return true")
	}
	node := doc.Get(path)
	if node != nil {
		t.Error("expected title to be removed")
	}

	// author should remain
	authorPath, _ := pathsyntax.Parse(".author")
	node = doc.Get(authorPath)
	if node == nil || node.Value != "Test" {
		t.Errorf("expected author to be preserved, got %v", node)
	}
}

func TestUnset_MissingKey(t *testing.T) {
	content := "---\ntitle: Hello\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".nonexistent")
	removed := doc.Unset(path)
	if removed {
		t.Error("expected Unset to return false for missing key")
	}
}

func TestUnset_NestedKey(t *testing.T) {
	content := "---\nauthor:\n  name: John\n  email: john@test.com\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".author.email")
	removed := doc.Unset(path)
	if !removed {
		t.Error("expected Unset to return true")
	}
	node := doc.Get(path)
	if node != nil {
		t.Error("expected email to be removed")
	}

	// name should remain
	namePath, _ := pathsyntax.Parse(".author.name")
	node = doc.Get(namePath)
	if node == nil || node.Value != "John" {
		t.Errorf("expected name to be preserved, got %v", node)
	}
}

func TestUnset_ArrayElement(t *testing.T) {
	content := "---\ntags:\n  - go\n  - yaml\n  - test\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".tags[1]")
	removed := doc.Unset(path)
	if !removed {
		t.Error("expected Unset to return true")
	}

	// Should now have 2 elements: go, test
	tagsPath, _ := pathsyntax.Parse(".tags")
	tagsNode := doc.Get(tagsPath)
	if tagsNode == nil {
		t.Fatal("expected tags to exist")
	}
	if len(tagsNode.Content) != 2 {
		t.Errorf("expected 2 elements, got %d", len(tagsNode.Content))
	}
	if tagsNode.Content[0].Value != "go" {
		t.Errorf("expected first element 'go', got %q", tagsNode.Content[0].Value)
	}
	if tagsNode.Content[1].Value != "test" {
		t.Errorf("expected second element 'test', got %q", tagsNode.Content[1].Value)
	}
}

func TestUnset_BodyPreserved(t *testing.T) {
	content := "---\ntitle: Hello\nauthor: Test\n---\n# My Body\nParagraph.\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".title")
	doc.Unset(path)
	if doc.Body != "# My Body\nParagraph.\n" {
		t.Errorf("body changed: %q", doc.Body)
	}
}

func TestUnset_NoFrontMatter(t *testing.T) {
	doc := &Document{Body: "# Hello\n"}
	path, _ := pathsyntax.Parse(".title")
	removed := doc.Unset(path)
	if removed {
		t.Error("expected Unset to return false with no front matter")
	}
}

func TestUnset_Root(t *testing.T) {
	content := "---\ntitle: Hello\nauthor: Test\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".")
	removed := doc.Unset(path)
	if !removed {
		t.Error("expected Unset root to return true")
	}
}
