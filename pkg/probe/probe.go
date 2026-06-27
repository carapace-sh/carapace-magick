package probe

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// ImageInfo holds metadata about an image probed via magick identify.
type ImageInfo struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Format     string `json:"format"`
	Colors     int    `json:"colors"`
	Depth      int    `json:"depth"`
	Colorspace string `json:"colorspace"`
}

// Probe runs `magick identify -verbose` on the given input path and
// extracts key image attributes. Returns nil, nil if magick is
// unavailable or the file is not local.
func Probe(inputPath string) *ImageInfo {
	cmd := exec.Command("magick", "identify", "-verbose", inputPath)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	info := &ImageInfo{}
	for line := range strings.SplitSeq(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Geometry:") {
			parseGeometry(line, info)
		} else if strings.HasPrefix(line, "Format:") {
			info.Format = parseFieldValue(line, "Format:")
		} else if strings.HasPrefix(line, "Colorspace:") {
			info.Colorspace = parseFieldValue(line, "Colorspace:")
		} else if strings.HasPrefix(line, "Depth:") {
			fmt.Sscanf(parseFieldValue(line, "Depth:"), "%d", &info.Depth)
		} else if strings.HasPrefix(line, "Colors:") {
			fmt.Sscanf(parseFieldValue(line, "Colors:"), "%d", &info.Colors)
		}
	}

	return info
}

// ProbeJSON runs magick identify with JSON output if available,
// falling back to verbose parsing.
func ProbeJSON(inputPath string) *ImageInfo {
	cmd := exec.Command("magick", "identify", inputPath, "-format",
		`{"width":%w,"height":%h,"format":"%m","depth":%z,"colorspace":"%r"}`)
	output, err := cmd.Output()
	if err != nil {
		return Probe(inputPath)
	}

	var info ImageInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return Probe(inputPath)
	}
	return &info
}

func parseGeometry(line string, info *ImageInfo) {
	field := parseFieldValue(line, "Geometry:")
	parts := strings.SplitN(field, "x", 2)
	if len(parts) >= 2 {
		fmt.Sscanf(parts[0], "%d", &info.Width)
		heightPart := parts[1]
		if plusIdx := strings.Index(heightPart, "+"); plusIdx >= 0 {
			heightPart = heightPart[:plusIdx]
		}
		fmt.Sscanf(heightPart, "%d", &info.Height)
	}
}

func parseFieldValue(line, prefix string) string {
	field := strings.TrimPrefix(line, prefix)
	return strings.TrimSpace(field)
}
