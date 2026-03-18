package assertion

import (
	"testing"

	"github.com/UnitVectorY-Labs/frontmatterkit/internal/frontmatter"
)

const testDocument = `---
title: Hello World
draft: false
count: 42
tags:
  - go
  - yaml
author:
  name: Test
description: "A test document about Go"
---
Body content
`

const emptyDocument = `No front matter here`

func mustParseDoc(t *testing.T, content string) *frontmatter.Document {
	t.Helper()
	doc, err := frontmatter.Parse(content)
	if err != nil {
		t.Fatalf("failed to parse document: %v", err)
	}
	return doc
}

func mustParseAssertion(t *testing.T, expr string) *Assertion {
	t.Helper()
	a, err := ParseAssertion(expr)
	if err != nil {
		t.Fatalf("failed to parse assertion %q: %v", expr, err)
	}
	return a
}

func TestEvaluate_Exists(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".title exists")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_NotExists_Pass(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".missing not exists")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_Exists_Fail(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".missing exists")
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_EqualsQuotedString(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.title == "Hello World"`)
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_EqualsUnquotedString(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".title == Hello World")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_EqualsBool(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".draft == false")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_GreaterEqual(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".count >= 3")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_GreaterThan_Fail(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".count > 100")
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_LessThan(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".count < 50")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_LessEqual(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".count <= 42")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_ContainsArray(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.tags contains "go"`)
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_NotContainsArray(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.tags not contains "rust"`)
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_ContainsString(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.description contains "Go"`)
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_NotEquals(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".draft != true")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_EmptyDoc_ExistsFail(t *testing.T) {
	doc := mustParseDoc(t, emptyDocument)
	a := mustParseAssertion(t, ".title exists")
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_EmptyDoc_NotExistsPass(t *testing.T) {
	doc := mustParseDoc(t, emptyDocument)
	a := mustParseAssertion(t, ".title not exists")
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_NestedPath(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.author.name == "Test"`)
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestEvaluate_ContainsArray_Fail(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.tags contains "rust"`)
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_NotContainsArray_Fail(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.tags not contains "go"`)
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_NumericOnNonExistent(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".missing > 5")
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_NumericOnNonNumeric(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, ".title > 5")
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_ContainsOnNonExistent(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.missing contains "x"`)
	if err := a.Evaluate(doc); err == nil {
		t.Errorf("expected failure, got pass")
	}
}

func TestEvaluate_NotContainsOnNonExistent(t *testing.T) {
	doc := mustParseDoc(t, testDocument)
	a := mustParseAssertion(t, `.missing not contains "x"`)
	if err := a.Evaluate(doc); err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}
