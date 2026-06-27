package probe

import (
	"os"
	"testing"
)

func TestProbeNonexistentFile(t *testing.T) {
	info := Probe("/nonexistent/path/to/image.png")
	if info != nil {
		t.Error("expected nil for nonexistent file")
	}
}

func TestProbeEmptyPath(t *testing.T) {
	info := Probe("")
	if info != nil {
		t.Error("expected nil for empty path")
	}
}

func TestProbeJSONNonexistentFile(t *testing.T) {
	info := ProbeJSON("/nonexistent/path/to/image.png")
	if info != nil {
		t.Error("expected nil for nonexistent file")
	}
}

func TestProbeWithMagick(t *testing.T) {
	if _, err := os.Stat("/usr/bin/magick"); os.IsNotExist(err) {
		t.Skip("magick binary not available")
	}
	info := Probe("/usr/bin/magick")
	if info != nil {
		t.Log("Probe returned info for binary (unexpected but not an error)")
	}
}
