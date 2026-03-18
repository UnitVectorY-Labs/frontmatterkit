package frontmatter

import (
	"testing"
)

func TestParse_NoFrontMatter(t *testing.T) {
	content := "# Hello World\n\nThis is a test."
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.HasFrontMatter {
		t.Error("expected HasFrontMatter to be false")
	}
	if doc.Body != content {
		t.Errorf("expected body to be full content, got %q", doc.Body)
	}
}

func TestParse_ValidFrontMatter(t *testing.T) {
	content := "---\ntitle: Hello\nauthor: Test\n---\n# Body\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !doc.HasFrontMatter {
		t.Error("expected HasFrontMatter to be true")
	}
	if doc.RawFrontMatter != "title: Hello\nauthor: Test\n" {
		t.Errorf("unexpected raw front matter: %q", doc.RawFrontMatter)
	}
	if doc.Body != "# Body\n" {
		t.Errorf("unexpected body: %q", doc.Body)
	}
	if doc.Node == nil {
		t.Error("expected Node to be non-nil")
	}
}

func TestParse_EmptyFrontMatter(t *testing.T) {
	content := "---\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !doc.HasFrontMatter {
		t.Error("expected HasFrontMatter to be true")
	}
	if doc.Body != "" {
		t.Errorf("expected empty body, got %q", doc.Body)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	content := "---\n[invalid yaml: :\n---\n"
	_, err := Parse(content)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParse_UnclosedFrontMatter(t *testing.T) {
	content := "---\ntitle: Hello\n"
	_, err := Parse(content)
	if err == nil {
		t.Error("expected error for unclosed front matter")
	}
}

func TestParse_BodyPreservation(t *testing.T) {
	body := "# Title\n\nParagraph 1\n\nParagraph 2\n"
	content := "---\ntitle: test\n---\n" + body
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.Body != body {
		t.Errorf("body not preserved: got %q, want %q", doc.Body, body)
	}
}

func TestValidate_NoFrontMatter(t *testing.T) {
	doc := &Document{HasFrontMatter: false}
	if err := Validate(doc); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_DuplicateKeys(t *testing.T) {
	content := "---\ntitle: Hello\ntitle: World\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if err := Validate(doc); err == nil {
		t.Error("expected validation error for duplicate keys")
	}
}

func TestValidate_Anchors(t *testing.T) {
	content := "---\nbase: &base\n  key: value\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if err := Validate(doc); err == nil {
		t.Error("expected validation error for anchors")
	}
}

func TestValidate_Aliases(t *testing.T) {
	content := "---\nbase: &base\n  key: value\nref: *base\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if err := Validate(doc); err == nil {
		t.Error("expected validation error for aliases")
	}
}

func TestValidate_CustomTag(t *testing.T) {
	content := "---\nval: !custom test\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if err := Validate(doc); err == nil {
		t.Error("expected validation error for custom tag")
	}
}

func TestValidate_ValidDocument(t *testing.T) {
	content := "---\ntitle: Hello\ntags:\n  - go\n  - yaml\nnested:\n  key: value\n---\nBody text\n"
	doc, err := ParseAndValidate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !doc.HasFrontMatter {
		t.Error("expected HasFrontMatter to be true")
	}
}

func TestParseAndValidate_InvalidYAML(t *testing.T) {
	content := "---\n[bad yaml\n---\n"
	_, err := ParseAndValidate(content)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseAndValidate_DuplicateKeys(t *testing.T) {
	content := "---\na: 1\na: 2\n---\n"
	_, err := ParseAndValidate(content)
	if err == nil {
		t.Error("expected error for duplicate keys")
	}
}

func TestRender_NoFrontMatter(t *testing.T) {
	doc := &Document{Body: "# Hello\n"}
	result := doc.Render()
	if result != "# Hello\n" {
		t.Errorf("unexpected render: %q", result)
	}
}

func TestRender_WithFrontMatter(t *testing.T) {
	content := "---\ntitle: Hello\n---\n# Body\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rendered := doc.Render()
	// Re-parse the rendered content to verify round-trip
	doc2, err := Parse(rendered)
	if err != nil {
		t.Fatalf("render produced unparseable content: %v", err)
	}
	if !doc2.HasFrontMatter {
		t.Error("rendered content lost front matter")
	}
	if doc2.Body != doc.Body {
		t.Errorf("body changed after render: got %q, want %q", doc2.Body, doc.Body)
	}
}

func TestRender_EmptyFrontMatter(t *testing.T) {
	content := "---\n---\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rendered := doc.Render()
	if rendered != "---\n---\n" {
		t.Errorf("unexpected render for empty front matter: %q", rendered)
	}
}

func TestRender_NilNode(t *testing.T) {
	doc := &Document{HasFrontMatter: true, Node: nil, Body: "body\n"}
	rendered := doc.Render()
	if rendered != "---\n---\nbody\n" {
		t.Errorf("unexpected render with nil Node: %q", rendered)
	}
}
