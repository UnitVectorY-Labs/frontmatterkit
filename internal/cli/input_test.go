package cli

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRenderJSON(t *testing.T) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte("z: '<>&'\na: [one, two]\n"), &node); err != nil {
		t.Fatal(err)
	}
	got, err := renderJSON(&node)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"a\": [\n    \"one\",\n    \"two\"\n  ],\n  \"z\": \"<>&\"\n}\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderJSONRejectsInvalidUTF8(t *testing.T) {
	node := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "\xff"}
	if _, err := renderJSON(node); err == nil {
		t.Fatal("expected invalid UTF-8 error")
	}
}
