package cmd_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var testBinary string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "fmk-test-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	testBinary = filepath.Join(tmpDir, "frontmatterkit")
	if runtime.GOOS == "windows" {
		testBinary += ".exe"
	}

	absRoot, err := filepath.Abs("..")
	if err != nil {
		panic("failed to resolve project root: " + err.Error())
	}

	build := exec.Command("go", "build", "-o", testBinary, ".")
	build.Dir = absRoot
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + string(out) + ": " + err.Error())
	}

	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

func testdataPath(name string) string {
	absRoot, _ := filepath.Abs("..")
	return filepath.Join(absRoot, "testdata", name)
}

func readTestdata(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(testdataPath(name))
	if err != nil {
		t.Fatalf("failed to read testdata/%s: %v", name, err)
	}
	return string(data)
}

func runCLI(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(testBinary, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run command: %v", err)
		}
	}

	return outBuf.String(), errBuf.String(), exitCode
}

// ---------------------------------------------------------------------------
// Validate tests
// ---------------------------------------------------------------------------

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdin    string
		wantExit int
	}{
		{
			name:     "valid simple file",
			args:     []string{"validate", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "valid empty front matter",
			args:     []string{"validate", testdataPath("valid-empty-frontmatter.md")},
			wantExit: 0,
		},
		{
			name:     "no front matter is valid",
			args:     []string{"validate", testdataPath("no-frontmatter.md")},
			wantExit: 0,
		},
		{
			name:     "invalid yaml",
			args:     []string{"validate", testdataPath("invalid-yaml.md")},
			wantExit: 1,
		},
		{
			name:     "duplicate keys",
			args:     []string{"validate", testdataPath("duplicate-keys.md")},
			wantExit: 1,
		},
		{
			name:     "yaml anchors rejected",
			args:     []string{"validate", testdataPath("with-anchor.md")},
			wantExit: 1,
		},
		{
			name:     "no closing delimiter",
			args:     []string{"validate", testdataPath("no-closing-delimiter.md")},
			wantExit: 1,
		},
		{
			name:     "complex valid file",
			args:     []string{"validate", testdataPath("complex.md")},
			wantExit: 0,
		},
		{
			name:     "multiline body valid",
			args:     []string{"validate", testdataPath("multiline-body.md")},
			wantExit: 0,
		},
		{
			name:     "valid via stdin",
			args:     []string{"validate"},
			stdin:    readFixture("valid-simple.md"),
			wantExit: 0,
		},
		{
			name:     "invalid via stdin",
			args:     []string{"validate"},
			stdin:    readFixture("invalid-yaml.md"),
			wantExit: 1,
		},
		{
			name:     "empty input via stdin",
			args:     []string{"validate"},
			stdin:    "",
			wantExit: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, exit := runCLI(t, tc.stdin, tc.args...)
			if exit != tc.wantExit {
				t.Errorf("exit code = %d, want %d", exit, tc.wantExit)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Get tests
// ---------------------------------------------------------------------------

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		stdin      string
		wantExit   int
		wantStdout string   // exact match (empty means skip)
		wantContains []string // substring checks
	}{
		{
			name:     "full yaml output",
			args:     []string{"get", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantContains: []string{
				"title: Hello World",
				"draft: false",
				"count: 42",
			},
		},
		{
			name:       "get scalar .title",
			args:       []string{"get", "--path", ".title", testdataPath("valid-simple.md")},
			wantExit:   0,
			wantStdout: "Hello World\n",
		},
		{
			name:       "get nested .author.name",
			args:       []string{"get", "--path", ".author.name", testdataPath("valid-simple.md")},
			wantExit:   0,
			wantStdout: "Test User\n",
		},
		{
			name:       "get array element .tags[0]",
			args:       []string{"get", "--path", ".tags[0]", testdataPath("valid-simple.md")},
			wantExit:   0,
			wantStdout: "go\n",
		},
		{
			name:       "get array element .tags[1]",
			args:       []string{"get", "--path", ".tags[1]", testdataPath("valid-simple.md")},
			wantExit:   0,
			wantStdout: "yaml\n",
		},
		{
			name:       "get nonexistent path returns null",
			args:       []string{"get", "--path", ".nonexistent", testdataPath("valid-simple.md")},
			wantExit:   0,
			wantStdout: "null\n",
		},
		{
			name:     "get json tags",
			args:     []string{"get", "--format", "json", "--path", ".tags", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantContains: []string{`"go"`, `"yaml"`},
		},
		{
			name:     "get full json",
			args:     []string{"get", "--format", "json", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantContains: []string{
				`"title"`,
				`"Hello World"`,
				`"count"`,
				`"draft"`,
			},
		},
		{
			name:       "get no-frontmatter returns null",
			args:       []string{"get", testdataPath("no-frontmatter.md")},
			wantExit:   0,
			wantStdout: "null\n",
		},
		{
			name:       "get path from no-frontmatter returns null",
			args:       []string{"get", "--path", ".title", testdataPath("no-frontmatter.md")},
			wantExit:   0,
			wantStdout: "null\n",
		},
		{
			name:       "get nested deep value from complex",
			args:       []string{"get", "--path", ".nested.deep.value", testdataPath("complex.md")},
			wantExit:   0,
			wantStdout: "42\n",
		},
		{
			name:     "get json categories from complex",
			args:     []string{"get", "--format", "json", "--path", ".categories", testdataPath("complex.md")},
			wantExit: 0,
			wantContains: []string{`"tech"`, `"programming"`, `"go"`},
		},
		{
			name:       "get via stdin",
			args:       []string{"get", "--path", ".title"},
			stdin:      readFixture("valid-simple.md"),
			wantExit:   0,
			wantStdout: "Hello World\n",
		},
		{
			name:       "get empty front matter",
			args:       []string{"get", testdataPath("valid-empty-frontmatter.md")},
			wantExit:   0,
			wantStdout: "{}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, exit := runCLI(t, tc.stdin, tc.args...)
			if exit != tc.wantExit {
				t.Errorf("exit code = %d, want %d", exit, tc.wantExit)
			}
			if tc.wantStdout != "" && stdout != tc.wantStdout {
				t.Errorf("stdout = %q, want %q", stdout, tc.wantStdout)
			}
			for _, want := range tc.wantContains {
				if !strings.Contains(stdout, want) {
					t.Errorf("stdout missing %q\ngot: %s", want, stdout)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Set tests
// ---------------------------------------------------------------------------

func TestSet(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		stdin        string
		wantExit     int
		wantFile     string   // expected output testdata file name
		wantContains []string // substring checks
	}{
		{
			name:     "set title",
			args:     []string{"set", "--set", ".title=New Title", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-set-title.md",
		},
		{
			name:     "set new key",
			args:     []string{"set", "--set", ".newkey=newvalue", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-set-newkey.md",
		},
		{
			name:     "set creates front matter on no-frontmatter file",
			args:     []string{"set", "--set", ".title=New Title", testdataPath("no-frontmatter.md")},
			wantExit: 0,
			wantFile: "expected-set-no-frontmatter.md",
		},
		{
			name:     "set numeric value",
			args:     []string{"set", "--set", ".count=100", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-set-count.md",
		},
		{
			name:     "set boolean value",
			args:     []string{"set", "--set", ".draft=true", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-set-draft-true.md",
		},
		{
			name:     "set array value",
			args:     []string{"set", "--set", ".tags=[a, b, c]", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-set-tags-array.md",
		},
		{
			name:     "set multiple values",
			args:     []string{"set", "--set", ".title=Updated", "--set", ".count=99", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-set-multiple.md",
		},
		{
			name:         "set via stdin",
			args:         []string{"set", "--set", ".title=Stdin Title"},
			stdin:        readFixture("valid-simple.md"),
			wantExit:     0,
			wantContains: []string{"title: Stdin Title", "This is the body content."},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, exit := runCLI(t, tc.stdin, tc.args...)
			if exit != tc.wantExit {
				t.Errorf("exit code = %d, want %d", exit, tc.wantExit)
			}
			if tc.wantFile != "" {
				want := readTestdata(t, tc.wantFile)
				if stdout != want {
					t.Errorf("stdout mismatch\ngot:\n%s\nwant:\n%s", stdout, want)
				}
			}
			for _, want := range tc.wantContains {
				if !strings.Contains(stdout, want) {
					t.Errorf("stdout missing %q\ngot: %s", want, stdout)
				}
			}
		})
	}
}

func TestSetInPlace(t *testing.T) {
	// Copy input to a temp file, run set --in-place, verify file content.
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "input.md")

	data, err := os.ReadFile(testdataPath("valid-simple.md"))
	if err != nil {
		t.Fatalf("read source: %v", err)
	}
	if err := os.WriteFile(src, data, 0644); err != nil {
		t.Fatalf("write temp: %v", err)
	}

	_, _, exit := runCLI(t, "", "set", "--set", ".title=InPlace", "--in-place", src)
	if exit != 0 {
		t.Fatalf("exit code = %d, want 0", exit)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if !strings.Contains(string(got), "title: InPlace") {
		t.Errorf("in-place file missing updated title\ngot:\n%s", got)
	}
	if !strings.Contains(string(got), "This is the body content.") {
		t.Errorf("in-place file missing body\ngot:\n%s", got)
	}
}

func TestSetOutput(t *testing.T) {
	tmpDir := t.TempDir()
	out := filepath.Join(tmpDir, "output.md")

	_, _, exit := runCLI(t, "", "set", "--set", ".title=OutputTest", "--output", out, testdataPath("valid-simple.md"))
	if exit != 0 {
		t.Fatalf("exit code = %d, want 0", exit)
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(got), "title: OutputTest") {
		t.Errorf("output file missing updated title\ngot:\n%s", got)
	}
	if !strings.Contains(string(got), "This is the body content.") {
		t.Errorf("output file missing body\ngot:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Unset tests
// ---------------------------------------------------------------------------

func TestUnset(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdin    string
		wantExit int
		wantFile string
	}{
		{
			name:     "unset draft",
			args:     []string{"unset", "--path", ".draft", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-unset-draft.md",
		},
		{
			name:     "unset nonexistent is no-op",
			args:     []string{"unset", "--path", ".nonexistent", testdataPath("valid-simple.md")},
			wantExit: 0,
			// Output should be same as input (re-rendered)
		},
		{
			name:     "unset nested .author.name",
			args:     []string{"unset", "--path", ".author.name", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-unset-author-name.md",
		},
		{
			name:     "unset array element .tags[0]",
			args:     []string{"unset", "--path", ".tags[0]", testdataPath("valid-simple.md")},
			wantExit: 0,
			wantFile: "expected-unset-tags0.md",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, exit := runCLI(t, tc.stdin, tc.args...)
			if exit != tc.wantExit {
				t.Errorf("exit code = %d, want %d", exit, tc.wantExit)
			}
			if tc.wantFile != "" {
				want := readTestdata(t, tc.wantFile)
				if stdout != want {
					t.Errorf("stdout mismatch\ngot:\n%s\nwant:\n%s", stdout, want)
				}
			}
		})
	}
}

func TestUnsetInPlace(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "input.md")

	data, err := os.ReadFile(testdataPath("valid-simple.md"))
	if err != nil {
		t.Fatalf("read source: %v", err)
	}
	if err := os.WriteFile(src, data, 0644); err != nil {
		t.Fatalf("write temp: %v", err)
	}

	_, _, exit := runCLI(t, "", "unset", "--path", ".draft", "--in-place", src)
	if exit != 0 {
		t.Fatalf("exit code = %d, want 0", exit)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if strings.Contains(string(got), "draft") {
		t.Errorf("in-place file still contains draft\ngot:\n%s", got)
	}
}

func TestUnsetOutput(t *testing.T) {
	tmpDir := t.TempDir()
	out := filepath.Join(tmpDir, "output.md")

	_, _, exit := runCLI(t, "", "unset", "--path", ".draft", "--output", out, testdataPath("valid-simple.md"))
	if exit != 0 {
		t.Fatalf("exit code = %d, want 0", exit)
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if strings.Contains(string(got), "draft") {
		t.Errorf("output file still contains draft\ngot:\n%s", got)
	}
}

// ---------------------------------------------------------------------------
// Assert tests
// ---------------------------------------------------------------------------

func TestAssert(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdin    string
		wantExit int
	}{
		{
			name:     "exists passes",
			args:     []string{"assert", "--assert", ".title exists", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "exists fails for missing",
			args:     []string{"assert", "--assert", ".missing exists", testdataPath("valid-simple.md")},
			wantExit: 1,
		},
		{
			name:     "not exists passes for missing",
			args:     []string{"assert", "--assert", ".missing not exists", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "equality boolean",
			args:     []string{"assert", "--assert", ".draft == false", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "greater or equal",
			args:     []string{"assert", "--assert", ".count >= 42", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "greater than fails",
			args:     []string{"assert", "--assert", ".count > 100", testdataPath("valid-simple.md")},
			wantExit: 1,
		},
		{
			name:     "less than passes",
			args:     []string{"assert", "--assert", ".count < 50", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "less or equal",
			args:     []string{"assert", "--assert", ".count <= 42", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "contains in array",
			args:     []string{"assert", "--assert", `.tags contains "go"`, testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "not contains in array",
			args:     []string{"assert", "--assert", `.tags not contains "rust"`, testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "equality string",
			args:     []string{"assert", "--assert", `.title == "Hello World"`, testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "not equals string",
			args:     []string{"assert", "--assert", `.title != "Other"`, testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "multiple asserts all pass",
			args:     []string{"assert", "--assert", ".title exists", "--assert", ".draft == false", testdataPath("valid-simple.md")},
			wantExit: 0,
		},
		{
			name:     "multiple asserts one fails",
			args:     []string{"assert", "--assert", ".title exists", "--assert", ".count > 100", testdataPath("valid-simple.md")},
			wantExit: 1,
		},
		{
			name:     "no-frontmatter exists fails",
			args:     []string{"assert", "--assert", ".title exists", testdataPath("no-frontmatter.md")},
			wantExit: 1,
		},
		{
			name:     "no-frontmatter not exists passes",
			args:     []string{"assert", "--assert", ".title not exists", testdataPath("no-frontmatter.md")},
			wantExit: 0,
		},
		{
			name:     "assert via stdin",
			args:     []string{"assert", "--assert", ".title exists"},
			stdin:    readFixture("valid-simple.md"),
			wantExit: 0,
		},
		{
			name:     "assert complex nested value",
			args:     []string{"assert", "--assert", ".nested.deep.value == 42", testdataPath("complex.md")},
			wantExit: 0,
		},
		{
			name:     "assert complex nested greater",
			args:     []string{"assert", "--assert", ".nested.deep.value > 100", testdataPath("complex.md")},
			wantExit: 1,
		},
		{
			name:     "assert on empty front matter not exists",
			args:     []string{"assert", "--assert", ".title not exists", testdataPath("valid-empty-frontmatter.md")},
			wantExit: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, exit := runCLI(t, tc.stdin, tc.args...)
			if exit != tc.wantExit {
				t.Errorf("exit code = %d, want %d", exit, tc.wantExit)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Edge case / error tests
// ---------------------------------------------------------------------------

func TestEdgeCases(t *testing.T) {
	t.Run("nonexistent file returns exit 3", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "validate", "/nonexistent/file.md")
		if exit != 3 {
			t.Errorf("exit code = %d, want 3", exit)
		}
	})

	t.Run("get with invalid path returns exit 2", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "get", "--path", "invalid", testdataPath("valid-simple.md"))
		if exit != 2 {
			t.Errorf("exit code = %d, want 2", exit)
		}
	})

	t.Run("set with bad --set format returns exit 2", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "set", "--set", "no-equals", testdataPath("valid-simple.md"))
		if exit != 2 {
			t.Errorf("exit code = %d, want 2", exit)
		}
	})

	t.Run("unset requires --path flag", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "unset", testdataPath("valid-simple.md"))
		if exit == 0 {
			t.Error("expected non-zero exit code for missing --path")
		}
	})

	t.Run("assert requires --assert flag", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "assert", testdataPath("valid-simple.md"))
		if exit != 2 {
			t.Errorf("exit code = %d, want 2", exit)
		}
	})

	t.Run("set preserves multiline body", func(t *testing.T) {
		stdout, _, exit := runCLI(t, "", "set", "--set", ".title=Changed", testdataPath("multiline-body.md"))
		if exit != 0 {
			t.Fatalf("exit code = %d, want 0", exit)
		}
		if !strings.Contains(stdout, "Line 1.") || !strings.Contains(stdout, "Line 3.") {
			t.Errorf("multiline body not preserved\ngot:\n%s", stdout)
		}
	})

	t.Run("get on invalid yaml returns exit 1", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "get", testdataPath("invalid-yaml.md"))
		if exit != 1 {
			t.Errorf("exit code = %d, want 1", exit)
		}
	})

	t.Run("assert on invalid yaml returns exit 1", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "assert", "--assert", ".title exists", testdataPath("invalid-yaml.md"))
		if exit != 1 {
			t.Errorf("exit code = %d, want 1", exit)
		}
	})

	t.Run("set on invalid yaml returns exit 1", func(t *testing.T) {
		_, _, exit := runCLI(t, "", "set", "--set", ".title=foo", testdataPath("invalid-yaml.md"))
		if exit != 1 {
			t.Errorf("exit code = %d, want 1", exit)
		}
	})
}

// readFixture reads a testdata file for use as stdin input.
// Must be called at test definition time (not in TestMain).
func readFixture(name string) string {
	absRoot, _ := filepath.Abs("..")
	data, err := os.ReadFile(filepath.Join(absRoot, "testdata", name))
	if err != nil {
		panic("readFixture: " + err.Error())
	}
	return string(data)
}
