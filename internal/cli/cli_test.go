package cli_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const (
	rootCommandDir = "root"
)

var testBinary string
var repoRoot string

type testCase struct {
	command string
	name    string
	path    string
}

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "fmk-test-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	testBinary = filepath.Join(tmpDir, "frontmatterkit")
	if runtime.GOOS == "windows" {
		testBinary += ".exe"
	}

	repoRoot, err = filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		panic("failed to resolve project root: " + err.Error())
	}

	build := exec.Command("go", "build", "-o", testBinary, ".")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + string(out) + ": " + err.Error())
	}

	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

func TestCLICases(t *testing.T) {
	cases := discoverCases(t)
	if len(cases) == 0 {
		t.Fatal("no CLI test cases found")
	}

	for _, tc := range cases {
		tc := tc
		testName := tc.name
		if tc.command != rootCommandDir {
			testName = tc.command + "/" + tc.name
		}

		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			runCLICase(t, tc)
		})
	}
}

func discoverCases(t *testing.T) []testCase {
	t.Helper()

	root := filepath.Join(repoRoot, "testdata")
	commandEntries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read CLI test root: %v", err)
	}

	var cases []testCase
	for _, commandEntry := range commandEntries {
		if !commandEntry.IsDir() {
			continue
		}

		command := commandEntry.Name()
		commandPath := filepath.Join(root, command)
		caseEntries, err := os.ReadDir(commandPath)
		if err != nil {
			t.Fatalf("read command directory %s: %v", command, err)
		}

		for _, caseEntry := range caseEntries {
			if !caseEntry.IsDir() {
				continue
			}

			caseName := caseEntry.Name()
			casePath := filepath.Join(commandPath, caseName)
			if !fileExists(filepath.Join(casePath, caseName+".args")) {
				continue
			}

			cases = append(cases, testCase{
				command: command,
				name:    caseName,
				path:    casePath,
			})
		}
	}

	slices.SortFunc(cases, func(a, b testCase) int {
		if a.command != b.command {
			return strings.Compare(a.command, b.command)
		}
		return strings.Compare(a.name, b.name)
	})

	return cases
}

func runCLICase(t *testing.T, tc testCase) {
	t.Helper()

	tempRoot := t.TempDir()
	sourceTestdata := filepath.Join(repoRoot, "testdata")
	tempTestdata := filepath.Join(tempRoot, "testdata")
	if err := copyDir(sourceTestdata, tempTestdata); err != nil {
		t.Fatalf("copy testdata: %v", err)
	}

	caseDir := filepath.Join(tempTestdata, tc.command, tc.name)
	basePath := filepath.Join(caseDir, tc.name)

	args := loadArgs(t, basePath+".args")
	if tc.command != rootCommandDir {
		args = append([]string{tc.command}, args...)
	}

	stdin := loadStdin(t, basePath, args)
	wantExit := loadExitCode(t, basePath+".exitcode")

	stdout, stderr, exitCode := runCLI(t, caseDir, stdin, args...)
	if exitCode != wantExit {
		t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", exitCode, wantExit, stdout, stderr)
	}

	assertStream(t, basePath, "stdout", stdout)
	assertStream(t, basePath, "stderr", stderr)
	assertExpectedFiles(t, caseDir, tc.name)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, 0644)
	})
}

func loadArgs(t *testing.T, path string) []string {
	t.Helper()

	lines := loadLineFile(t, path, false)
	args := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		args = append(args, line)
	}
	return args
}

func loadExitCode(t *testing.T, path string) int {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read exit code %s: %v", path, err)
	}

	value, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("parse exit code %s: %v", path, err)
	}
	return value
}

func loadStdin(t *testing.T, basePath string, args []string) string {
	t.Helper()

	exactPath := basePath + ".stdin"
	markdownPath := basePath + ".stdin.md"
	defaultMarkdownPath := basePath + ".md"

	switch {
	case fileExists(exactPath):
		data, err := os.ReadFile(exactPath)
		if err != nil {
			t.Fatalf("read %s: %v", exactPath, err)
		}
		return string(data)
	case fileExists(markdownPath):
		data, err := os.ReadFile(markdownPath)
		if err != nil {
			t.Fatalf("read %s: %v", markdownPath, err)
		}
		return string(data)
	case fileExists(defaultMarkdownPath) && !argsReferenceFile(args, filepath.Base(defaultMarkdownPath)):
		data, err := os.ReadFile(defaultMarkdownPath)
		if err != nil {
			t.Fatalf("read %s: %v", defaultMarkdownPath, err)
		}
		return string(data)
	default:
		return ""
	}
}

func argsReferenceFile(args []string, fileName string) bool {
	for _, arg := range args {
		if arg == fileName {
			return true
		}
	}
	return false
}

func assertStream(t *testing.T, basePath, streamName, got string) {
	t.Helper()

	exactPath := basePath + "." + streamName + ".expected"
	containsPath := basePath + "." + streamName + ".contains"

	var want string
	var hasExpected bool
	if fileExists(exactPath) {
		data, err := os.ReadFile(exactPath)
		if err != nil {
			t.Fatalf("read %s: %v", exactPath, err)
		}
		want = string(data)
		hasExpected = true
	}

	if hasExpected && got != want {
		t.Fatalf("%s mismatch\ngot:\n%s\nwant:\n%s", streamName, got, want)
	}

	contains := []string(nil)
	if fileExists(containsPath) {
		contains = loadLineFile(t, containsPath, false)
	}

	if !hasExpected && len(contains) == 0 {
		if got != "" {
			t.Fatalf("%s mismatch\ngot:\n%s\nwant empty output", streamName, got)
		}
		return
	}

	for _, line := range contains {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(got, line) {
			t.Fatalf("%s missing %q\ngot:\n%s", streamName, line, got)
		}
	}
}

func assertExpectedFiles(t *testing.T, caseDir, caseName string) {
	t.Helper()

	err := filepath.WalkDir(caseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".expected") {
			return nil
		}
		if strings.HasSuffix(path, ".stdout.expected") || strings.HasSuffix(path, ".stderr.expected") {
			return nil
		}
		if !strings.HasPrefix(filepath.Base(path), caseName+".") {
			return nil
		}

		actualPath := strings.TrimSuffix(path, ".expected")
		want, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		got, err := os.ReadFile(actualPath)
		if err != nil {
			return fmt.Errorf("read actual file %s: %w", actualPath, err)
		}
		if string(got) != string(want) {
			return fmt.Errorf("%s mismatch\ngot:\n%s\nwant:\n%s", actualPath, string(got), string(want))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func loadLineFile(t *testing.T, path string, requireSingle bool) []string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if requireSingle {
		var filtered []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			filtered = append(filtered, line)
		}
		if len(filtered) != 1 {
			t.Fatalf("%s must contain exactly one non-comment line", path)
		}
		return filtered
	}

	return lines
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func runCLI(t *testing.T, workDir, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(testBinary, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Dir = workDir

	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	err := cmd.Run()
	if err == nil {
		return outBuf.String(), errBuf.String(), 0
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("failed to run command: %v", err)
	}

	return outBuf.String(), errBuf.String(), exitErr.ExitCode()
}
