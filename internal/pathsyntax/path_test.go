package pathsyntax

import (
	"testing"
)

func TestParseRoot(t *testing.T) {
	p, err := Parse(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.IsRoot {
		t.Fatal("expected IsRoot to be true")
	}
	if len(p.Segments) != 0 {
		t.Fatalf("expected 0 segments, got %d", len(p.Segments))
	}
}

func TestParseSimpleField(t *testing.T) {
	p, err := Parse(".title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.IsRoot {
		t.Fatal("expected IsRoot to be false")
	}
	if len(p.Segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "title")
}

func TestParseNestedFields(t *testing.T) {
	p, err := Parse(".author.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "author")
	assertField(t, p.Segments[1], "name")
}

func TestParseArrayIndex(t *testing.T) {
	p, err := Parse(".tags[0]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "tags")
	assertIndex(t, p.Segments[1], 0)
}

func TestParseComplexPath(t *testing.T) {
	p, err := Parse(".a.b[2].c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 4 {
		t.Fatalf("expected 4 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "a")
	assertField(t, p.Segments[1], "b")
	assertIndex(t, p.Segments[2], 2)
	assertField(t, p.Segments[3], "c")
}

func TestParseFieldWithHyphenAndUnderscore(t *testing.T) {
	p, err := Parse(".my-field.another_field")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "my-field")
	assertField(t, p.Segments[1], "another_field")
}

func TestParseFieldWithDigits(t *testing.T) {
	p, err := Parse(".field1.item2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "field1")
	assertField(t, p.Segments[1], "item2")
}

func TestParseMultipleIndices(t *testing.T) {
	p, err := Parse(".matrix[0][1]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "matrix")
	assertIndex(t, p.Segments[1], 0)
	assertIndex(t, p.Segments[2], 1)
}

func TestParseLargeIndex(t *testing.T) {
	p, err := Parse(".items[42]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(p.Segments))
	}
	assertField(t, p.Segments[0], "items")
	assertIndex(t, p.Segments[1], 42)
}

func TestStringRootPath(t *testing.T) {
	p := Path{IsRoot: true}
	if got := p.String(); got != "." {
		t.Fatalf("expected %q, got %q", ".", got)
	}
}

func TestStringRoundTrip(t *testing.T) {
	cases := []string{
		".",
		".title",
		".author.name",
		".tags[0]",
		".a.b[2].c",
		".my-field.another_field",
		".matrix[0][1]",
		".foo.bar[2].baz",
	}
	for _, input := range cases {
		p, err := Parse(input)
		if err != nil {
			t.Fatalf("Parse(%q) unexpected error: %v", input, err)
		}
		got := p.String()
		if got != input {
			t.Errorf("String() round-trip: input=%q got=%q", input, got)
		}
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		input string
		desc  string
	}{
		{"", "empty string"},
		{"foo", "missing dot prefix"},
		{"foo.bar", "missing dot prefix with nested"},
		{"..", "double dot"},
		{".foo.", "trailing dot"},
		{".foo[]", "empty index"},
		{".foo[-1]", "negative index"},
		{".foo[abc]", "non-numeric index"},
		{".foo[", "unclosed bracket"},
		{".foo[0", "unclosed bracket with number"},
		{".[0]", "index without field name"},
		{".foo.bar.", "trailing dot nested"},
		{".foo..bar", "double dot nested"},
		{".foo[01]", "leading zero in index"},
		{".foo bar", "space in field name"},
		{".foo/bar", "slash in field name"},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := Parse(tc.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error for %s, got nil", tc.input, tc.desc)
			}
		})
	}
}

func assertField(t *testing.T, seg Segment, name string) {
	t.Helper()
	if seg.Type != SegmentField {
		t.Fatalf("expected SegmentField, got %v", seg.Type)
	}
	if seg.Field != name {
		t.Fatalf("expected field %q, got %q", name, seg.Field)
	}
}

func assertIndex(t *testing.T, seg Segment, idx int) {
	t.Helper()
	if seg.Type != SegmentIndex {
		t.Fatalf("expected SegmentIndex, got %v", seg.Type)
	}
	if seg.Index != idx {
		t.Fatalf("expected index %d, got %d", idx, seg.Index)
	}
}
