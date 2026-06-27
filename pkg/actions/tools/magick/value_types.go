package magick

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/uid"
)

func Uid(host string, opts ...string) func(s string, uc uid.Context) (*url.URL, error) {
	return func(s string, uc uid.Context) (*url.URL, error) {
		if length := len(opts); length%2 != 0 {
			return nil, fmt.Errorf("invalid amount of arguments [magick.Uid]: %v", length)
		}

		uid := &url.URL{
			Scheme: "magick",
			Host:   host,
			Path:   s,
		}
		values := uid.Query()
		for i := 0; i < len(opts); i += 2 {
			if opts[i+1] != "" {
				values.Add(opts[i], opts[i+1])
			}
		}
		uid.RawQuery = values.Encode()

		return uid, nil
	}
}

// Dynamic actions (shell out to magick -list)

func ActionColorspaces() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "colorspace")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("colorspaces").UidF(Uid("colorspace"))
	})
}

func ActionComposes() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "compose")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("composes").UidF(Uid("compose"))
	})
}

func ActionCompressTypes() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "compress")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("compress types").UidF(Uid("compress"))
	})
}

func ActionChannels() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "channel")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("channels").UidF(Uid("channel"))
	})
}

func ActionDistortMethods() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "distort")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("distort methods").UidF(Uid("distort"))
	})
}

func ActionFilters() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "filter")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("filters").UidF(Uid("filter"))
	})
}

func ActionGravities() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "gravity")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("gravities").UidF(Uid("gravity"))
	})
}

func ActionInterlaceTypes() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "interlace")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("interlace types").UidF(Uid("interlace"))
	})
}

func ActionLayerMethods() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "layers")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("layer methods").UidF(Uid("layers"))
	})
}

func ActionMorphologyMethods() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "morphology")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("morphology methods").UidF(Uid("morphology"))
	})
}

func ActionTypes() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "type")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("types").UidF(Uid("type"))
	})
}

func ActionVirtualPixelMethods() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "virtual-pixel")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("virtual pixel methods").UidF(Uid("virtual-pixel"))
	})
}

func ActionOrientations() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "orientation")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("orientations").UidF(Uid("orientation"))
	})
}

func ActionDisposes() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "dispose")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("disposes").UidF(Uid("dispose"))
	})
}

func ActionMetrics() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "metric")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("metrics").UidF(Uid("metric"))
	})
}

func ActionEvaluateOps() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "evaluate")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("evaluate ops").UidF(Uid("evaluate"))
	})
}

func ActionFormats() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "format")(func(output []byte) carapace.Action {
		lines := strings.Split(string(output), "\n")
		r := regexp.MustCompile(`^\s+([A-Za-z0-9_]+)\s+`)
		var vals []string
		for _, line := range lines {
			if matches := r.FindStringSubmatch(line); matches != nil {
				vals = append(vals, matches[1])
			}
		}
		return carapace.ActionValues(vals...).Tag("formats").UidF(Uid("format"))
	})
}

func ActionFonts() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "font")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("fonts").UidF(Uid("font"))
	})
}

func ActionColors() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "color")(func(output []byte) carapace.Action {
		lines := strings.Split(string(output), "\n")
		r := regexp.MustCompile(`^\s+([a-zA-Z0-9_]+)\s+`)
		var vals []string
		for _, line := range lines {
			if matches := r.FindStringSubmatch(line); matches != nil {
				vals = append(vals, matches[1])
			}
		}
		return carapace.ActionValues(vals...).Tag("colors").UidF(Uid("color"))
	})
}

func ActionKernels() carapace.Action {
	return carapace.ActionExecCommand("magick", "-list", "kernel")(func(output []byte) carapace.Action {
		vals := parseListOutput(string(output))
		return carapace.ActionValues(vals...).Tag("kernels").UidF(Uid("kernel"))
	})
}

func ActionListTypes() carapace.Action {
	return carapace.ActionValuesDescribed(
		"align", "text alignment types",
		"alpha", "alpha channel types",
		"channel", "channel types",
		"color", "color names",
		"colorspace", "colorspace types",
		"compose", "compose operators",
		"compress", "compression types",
		"configure", "configure options",
		"delegate", "delegate formats",
		"density", "density options",
		"depth", "depth options",
		"dispose", "disposal methods",
		"distort", "distortion methods",
		"evaluate", "evaluate operators",
		"filter", "filter types",
		"font", "font names",
		"format", "image formats",
		"gravity", "gravity types",
		"interlace", "interlace types",
		"kernel", "kernel shapes",
		"layers", "layer methods",
		"line", "line join types",
		"list", "list types",
		"log", "log events",
		"metric", "comparison metrics",
		"method", "paint methods",
		"morphology", "morphology methods",
		"orientation", "image orientation",
		"policy", "security policy",
		"resource", "resource limits",
		"sparse-color", "sparse color methods",
		"storage", "storage types",
		"threshold", "threshold methods",
		"type", "image types",
		"virtual-pixel", "virtual pixel methods",
	).Tag("list types").Uid("magick", "list")
}

// parseListOutput parses the output of `magick -list <type>` into a list of values.
func parseListOutput(output string) []string {
	lines := strings.Split(output, "\n")
	var vals []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Path:") || strings.HasPrefix(trimmed, "Name:") || strings.HasPrefix(trimmed, "---") {
			continue
		}
		// Many list outputs have format: "  Name   Description"
		// Extract just the name (first word)
		fields := strings.Fields(trimmed)
		if len(fields) > 0 {
			name := fields[0]
			vals = append(vals, name)
		}
	}
	return vals
}

// Static actions

func ActionBoolean() carapace.Action {
	return carapace.ActionValues("true", "false", "1", "0").Tag("booleans").Uid("magick", "boolean")
}

func ActionAlphaOption() carapace.Action {
	return carapace.ActionValuesDescribed(
		"on", "enable alpha channel",
		"off", "disable alpha channel",
		"activate", "activate alpha channel",
		"deactivate", "deactivate alpha channel",
		"set", "set alpha channel",
		"opaque", "set alpha to fully opaque",
		"transparent", "set alpha to fully transparent",
		"extract", "extract alpha channel",
		"copy", "copy alpha channel",
		"background", "set alpha to background",
		"shape", "set alpha for shape",
	).Tag("alpha options").Uid("magick", "alpha")
}

func ActionAutoThreshold() carapace.Action {
	return carapace.ActionValuesDescribed(
		"Kapur", "Kapur threshold method",
		"OTSU", "OTSU threshold method",
		"Triangle", "Triangle threshold method",
	).Tag("auto threshold").Uid("magick", "auto-threshold")
}

func ActionNoiseTypes() carapace.Action {
	return carapace.ActionValuesDescribed(
		"Gaussian", "Gaussian noise",
		"Impulse", "Impulse (salt and pepper) noise",
		"Laplacian", "Laplacian noise",
		"Multiplicative", "Multiplicative noise",
		"Poisson", "Poisson noise",
		"Uniform", "Uniform noise",
	).Tag("noise types").Uid("magick", "noise")
}

func ActionPreviewTypes() carapace.Action {
	return carapace.ActionValues(
		"Rotate", "Shear", "Roll", "Noise", "Segment",
		"Blur", "Threshold", "Edge", "Spread", "Shade",
		"Solarize", "Implode", "Wave", "OilPaint", "Charcoal",
		"JPEG", "Hue", "Saturation", "Brightness", "Gamma",
		"Spiff", "Dull", "Grayscale", "Quantize", "Despeckle",
		"ReduceNoise", "AddNoise", "Sharpen", "Emboss", "EdgeDetect",
		"Raise", "Normalize", "Equalize", "Negate", "Stereo",
		"Polaroid", "Vignette", "Framed",
	).Tag("preview types").Uid("magick", "preview")
}

func ActionStorageTypes() carapace.Action {
	return carapace.ActionValuesDescribed(
		"char", "8-bit character",
		"short", "16-bit integer",
		"integer", "32-bit integer",
		"float", "32-bit float",
		"double", "64-bit double",
	).Tag("storage types").Uid("magick", "storage")
}

func ActionDirection() carapace.Action {
	return carapace.ActionValuesDescribed(
		"right-to-left", "Right to left",
		"left-to-right", "Left to right",
	).Tag("directions").Uid("magick", "direction")
}

// Tool name completion

func ActionToolNames() carapace.Action {
	return carapace.ActionValuesDescribed(
		"identify", "Describe image format and attributes",
		"mogrify", "In-place image transformation",
		"compare", "Assess difference between images",
		"composite", "Composite images together",
		"montage", "Create a composite image montage",
	).Tag("tools").Uid("magick", "tool")
}
