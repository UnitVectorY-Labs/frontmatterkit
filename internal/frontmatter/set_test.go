package frontmatter

import (
	"testing"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/pathsyntax"
)

func TestSet_NoFrontMatter(t *testing.T) {
	doc := &Document{Body: "# Hello\n"}
	path, _ := pathsyntax.Parse(".title")
	err := doc.Set(path, "Hello World", SetOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !doc.HasFrontMatter {
		t.Error("expected HasFrontMatter to be true")
	}
	node := doc.Get(path)
	if node == nil || node.Value != "Hello World" {
		t.Errorf("expected 'Hello World', got %v", node)
	}
}

func TestSet_ScalarField(t *testing.T) {
	content := "---\ntitle: Old\n---\nBody\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".title")
	err = doc.Set(path, "New Title", SetOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(path)
	if node == nil || node.Value != "New Title" {
		t.Errorf("expected 'New Title', got %v", node)
	}
}

func TestSet_NestedField(t *testing.T) {
	content := "---\ntitle: Hello\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".author.name")
	err = doc.Set(path, "John", SetOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(path)
	if node == nil || node.Value != "John" {
		t.Errorf("expected 'John', got %v", node)
	}
}

func TestSet_OverwriteExisting(t *testing.T) {
	content := "---\ntitle: Old\ncount: 1\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".count")
	err = doc.Set(path, "42", SetOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(path)
	if node == nil || node.Value != "42" {
		t.Errorf("expected '42', got %v", node)
	}
}

func TestSet_PatchMerge(t *testing.T) {
	content := "---\nauthor:\n  name: John\n  email: john@test.com\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".author")
	err = doc.Set(path, "name: Jane\nage: 30", SetPatch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// name should be updated
	namePath, _ := pathsyntax.Parse(".author.name")
	node := doc.Get(namePath)
	if node == nil || node.Value != "Jane" {
		t.Errorf("expected 'Jane', got %v", node)
	}

	// email should be preserved
	emailPath, _ := pathsyntax.Parse(".author.email")
	node = doc.Get(emailPath)
	if node == nil || node.Value != "john@test.com" {
		t.Errorf("expected 'john@test.com', got %v", node)
	}

	// age should be added
	agePath, _ := pathsyntax.Parse(".author.age")
	node = doc.Get(agePath)
	if node == nil || node.Value != "30" {
		t.Errorf("expected '30', got %v", node)
	}
}

func TestSet_CreateIntermediateMaps(t *testing.T) {
	content := "---\ntitle: Hello\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".a.b.c")
	err = doc.Set(path, "deep", SetOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := doc.Get(path)
	if node == nil || node.Value != "deep" {
		t.Errorf("expected 'deep', got %v", node)
	}
}

func TestSet_BodyPreserved(t *testing.T) {
	content := "---\ntitle: Old\n---\n# Body Content\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".title")
	err = doc.Set(path, "New", SetOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.Body != "# Body Content\n" {
		t.Errorf("body changed: %q", doc.Body)
	}
}

func TestSet_RootPatch(t *testing.T) {
	content := "---\ntitle: Hello\nexisting: true\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, _ := pathsyntax.Parse(".")
	err = doc.Set(path, "title: Updated\nnewkey: value", SetPatch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	titlePath, _ := pathsyntax.Parse(".title")
	node := doc.Get(titlePath)
	if node == nil || node.Value != "Updated" {
		t.Errorf("expected 'Updated', got %v", node)
	}

	existingPath, _ := pathsyntax.Parse(".existing")
	node = doc.Get(existingPath)
	if node == nil {
		t.Error("expected existing key to be preserved")
	}
}
