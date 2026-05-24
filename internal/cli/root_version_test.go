package cli

import (
	"fmt"
	"runtime"
	"testing"
)

func TestBuildVersionOutputPrefixesSemver(t *testing.T) {
	got := buildVersionOutput("1.2.3")
	want := fmt.Sprintf("v1.2.3 (%s, %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if got != want {
		t.Fatalf("buildVersionOutput() = %q, want %q", got, want)
	}
}

func TestBuildVersionOutputKeepsNonSemver(t *testing.T) {
	got := buildVersionOutput("dev")
	want := fmt.Sprintf("dev (%s, %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if got != want {
		t.Fatalf("buildVersionOutput() = %q, want %q", got, want)
	}
}

func TestBuildVersionOutputKeepsPrefixedSemver(t *testing.T) {
	got := buildVersionOutput("v1.2.3")
	want := fmt.Sprintf("v1.2.3 (%s, %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if got != want {
		t.Fatalf("buildVersionOutput() = %q, want %q", got, want)
	}
}
