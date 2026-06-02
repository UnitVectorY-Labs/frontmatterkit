package frontmatter

import (
	"errors"
	"strings"
	"testing"
)

const testSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["title", "draft", "tags"],
  "properties": {
    "title": { "type": "string" },
    "draft": { "type": "boolean" },
    "tags": {
      "type": "array",
      "items": { "type": "string" },
      "minItems": 1
    }
  },
  "additionalProperties": true
}`

func TestValidateJSONSchemaPassesMatchingDocument(t *testing.T) {
	doc := mustParseAndValidate(t, `---
title: Hello
draft: false
tags:
  - go
---
Body
`)

	if err := ValidateJSONSchema(doc, strings.NewReader(testSchema)); err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v", err)
	}
}

func TestValidateJSONSchemaFailsNonMatchingDocument(t *testing.T) {
	doc := mustParseAndValidate(t, `---
title: Hello
draft: "no"
tags: []
---
Body
`)

	err := ValidateJSONSchema(doc, strings.NewReader(testSchema))
	if err == nil {
		t.Fatal("ValidateJSONSchema() error = nil, want schema validation error")
	}
	if errors.Is(err, ErrInvalidJSONSchema) {
		t.Fatalf("ValidateJSONSchema() error = %v, want document validation error", err)
	}
}

func TestValidateJSONSchemaFailsMalformedSchemaJSON(t *testing.T) {
	doc := mustParseAndValidate(t, `---
title: Hello
---
Body
`)

	err := ValidateJSONSchema(doc, strings.NewReader(`{`))
	if !errors.Is(err, ErrInvalidJSONSchema) {
		t.Fatalf("ValidateJSONSchema() error = %v, want ErrInvalidJSONSchema", err)
	}
}

func TestValidateJSONSchemaFailsInvalidSchemaDefinition(t *testing.T) {
	doc := mustParseAndValidate(t, `---
title: Hello
---
Body
`)

	err := ValidateJSONSchema(doc, strings.NewReader(`{ "type": 123 }`))
	if !errors.Is(err, ErrInvalidJSONSchema) {
		t.Fatalf("ValidateJSONSchema() error = %v, want ErrInvalidJSONSchema", err)
	}
}

func TestValidateJSONSchemaFailsNoFrontMatterWithRequiredFields(t *testing.T) {
	doc, err := ParseAndValidate("no front matter here")
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}

	err = ValidateJSONSchema(doc, strings.NewReader(testSchema))
	if err == nil {
		t.Fatal("ValidateJSONSchema() error = nil, want schema validation error for missing required fields")
	}
	if errors.Is(err, ErrInvalidJSONSchema) {
		t.Fatalf("ValidateJSONSchema() error = %v, want document validation error", err)
	}
}

func TestValidateJSONSchemaPassesNoFrontMatterOptionalSchema(t *testing.T) {
	doc, err := ParseAndValidate("no front matter here")
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}

	err = ValidateJSONSchema(doc, strings.NewReader(`{"type": "object"}`))
	if err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v, want nil", err)
	}
}

func TestValidateJSONSchemaPassesEmptyFrontMatterOptionalSchema(t *testing.T) {
	doc := mustParseAndValidate(t, `---
---
Body
`)

	err := ValidateJSONSchema(doc, strings.NewReader(`{"type": "object"}`))
	if err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v, want nil", err)
	}
}

func TestValidateJSONSchemaFailsEmptyFrontMatterWithRequiredFields(t *testing.T) {
	doc := mustParseAndValidate(t, `---
---
Body
`)

	err := ValidateJSONSchema(doc, strings.NewReader(testSchema))
	if err == nil {
		t.Fatal("ValidateJSONSchema() error = nil, want schema validation error for missing required fields")
	}
	if errors.Is(err, ErrInvalidJSONSchema) {
		t.Fatalf("ValidateJSONSchema() error = %v, want document validation error", err)
	}
}

func TestValidateJSONSchemaPassesNilDocument(t *testing.T) {
	err := ValidateJSONSchema(nil, strings.NewReader(`{"type": "object"}`))
	if err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v, want nil", err)
	}
}

func mustParseAndValidate(t *testing.T, content string) *Document {
	t.Helper()

	doc, err := ParseAndValidate(content)
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}
	return doc
}
