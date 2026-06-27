package argstream

import (
	"slices"

	"github.com/carapace-sh/carapace/pkg/style"
)

// ToolProfile defines the option set and behavior for a magick sub-tool.
type ToolProfile struct {
	Name           string
	HasOutputArg   bool
	HasOperators   bool
	HasStackOps    bool
	HasParentheses bool
	OptionIndex    map[string]*OptionDef
}

// LookupOption looks up an option by name in the profile's option index.
// For the default magick profile (which shares the main OptionIndex), it looks up there.
// For sub-tool profiles that have their own index, only their own options are available.
func (p *ToolProfile) LookupOption(name string) *OptionDef {
	if p != nil && p.OptionIndex != nil {
		return p.OptionIndex[name]
	}
	return OptionIndex[name]
}

// DefaultMagickProfile is the profile for the magick command (default/convert).
var DefaultMagickProfile = &ToolProfile{
	Name:           "magick",
	HasOutputArg:   true,
	HasOperators:   true,
	HasStackOps:    true,
	HasParentheses: true,
	OptionIndex:    nil,
}

// DefaultIdentifyProfile is the profile for magick identify.
var DefaultIdentifyProfile = &ToolProfile{
	Name:           "identify",
	HasOutputArg:   false,
	HasOperators:   false,
	HasStackOps:    false,
	HasParentheses: false,
	OptionIndex:    nil,
}

// DefaultMogrifyProfile is the profile for magick mogrify.
var DefaultMogrifyProfile = &ToolProfile{
	Name:           "mogrify",
	HasOutputArg:   false,
	HasOperators:   true,
	HasStackOps:    false,
	HasParentheses: false,
	OptionIndex:    nil,
}

// DefaultCompareProfile is the profile for magick compare.
var DefaultCompareProfile = &ToolProfile{
	Name:           "compare",
	HasOutputArg:   true,
	HasOperators:   false,
	HasStackOps:    false,
	HasParentheses: false,
	OptionIndex:    nil,
}

// DefaultCompositeProfile is the profile for magick composite.
var DefaultCompositeProfile = &ToolProfile{
	Name:           "composite",
	HasOutputArg:   true,
	HasOperators:   false,
	HasStackOps:    false,
	HasParentheses: false,
	OptionIndex:    nil,
}

// DefaultMontageProfile is the profile for magick montage.
var DefaultMontageProfile = &ToolProfile{
	Name:           "montage",
	HasOutputArg:   true,
	HasOperators:   false,
	HasStackOps:    false,
	HasParentheses: false,
	OptionIndex:    nil,
}

// knownToolNames lists the valid sub-tool names.
var knownToolNames = []string{
	"identify",
	"mogrify",
	"compare",
	"composite",
	"montage",
}

// IsKnownToolName checks if a name is a known sub-tool.
func IsKnownToolName(name string) bool {
	return slices.Contains(knownToolNames, name)
}

// ProfileForTool returns the appropriate profile for a given tool name.
func ProfileForTool(name string) *ToolProfile {
	switch name {
	case "identify":
		return DefaultIdentifyProfile
	case "mogrify":
		return DefaultMogrifyProfile
	case "compare":
		return DefaultCompareProfile
	case "composite":
		return DefaultCompositeProfile
	case "montage":
		return DefaultMontageProfile
	default:
		return DefaultMagickProfile
	}
}

// Style returns the carapace style string for this option.
func (o *OptionDef) Style() string {
	switch o.Type {
	case TypeValue:
		return style.Carapace.FlagArg
	default:
		return style.Carapace.FlagNoArg
	}
}

func init() {
	DefaultMagickProfile.OptionIndex = OptionIndex
	DefaultIdentifyProfile.OptionIndex = buildIdentifyOptionIndex()
	DefaultMogrifyProfile.OptionIndex = buildMogrifyOptionIndex()
	DefaultCompareProfile.OptionIndex = buildCompareOptionIndex()
	DefaultCompositeProfile.OptionIndex = buildCompositeOptionIndex()
	DefaultMontageProfile.OptionIndex = buildMontageOptionIndex()
}

func buildIdentifyOptionIndex() map[string]*OptionDef {
	options := []*OptionDef{
		{Name: "debug", Description: "debug output", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "help", Description: "print program options", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "list", Description: "print supported values", Category: CategoryMisc, Type: TypeValue, ValueType: ValueList},
		{Name: "log", Description: "debug log format", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString},
		{Name: "verbose", Description: "print detailed information", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "version", Description: "print version", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "monitor", Description: "progress monitor", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "quiet", Description: "suppress warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "regard-warnings", Description: "pay attention to warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "seed", Description: "random seed", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "define", Description: "format-specific option", Category: CategorySetting, Type: TypeValue, ValueType: ValueDefine, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "density", Description: "horizontal/vertical resolution", Category: CategorySetting, Type: TypeValue, ValueType: ValueDensity},
		{Name: "format", Description: "output format string", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "ping", Description: "efficiently determine attributes", Category: CategorySetting, Type: TypeBoolean},
		{Name: "unique", Description: "display number of unique colors", Category: CategorySetting, Type: TypeBoolean},
		{Name: "moments", Description: "report image moments", Category: CategorySetting, Type: TypeBoolean},
		{Name: "features", Description: "analyze image features", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "channels", Description: "channel information", Category: CategorySetting, Type: TypeValue, ValueType: ValueChannel},
		{Name: "precision", Description: "print precision digits", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
	}
	return buildIndexFromOptions(options)
}

func buildMogrifyOptionIndex() map[string]*OptionDef {
	options := []*OptionDef{
		// Shared settings and operators (mogrify has operators but no stack/parentheses)
		{Name: "alpha", Description: "control alpha/matte channel", Category: CategorySetting, Type: TypeValue, ValueType: ValueAlpha, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "antialias", Description: "remove pixel aliasing", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "background", Description: "background color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "black-threshold", Description: "black threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "bordercolor", Description: "border color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "channel", Description: "set active channel mask", Category: CategoryChannelOp, Type: TypeValue, ValueType: ValueChannel, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "colors", Description: "preferred number of colors", Category: CategoryOperator, Type: TypeValue, ValueType: ValueInt},
		{Name: "colorspace", Description: "alternate colorspace", Category: CategorySetting, Type: TypeValue, ValueType: ValueColorspace, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "comment", Description: "annotate image with comment", Category: CategorySetting, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "compress", Description: "pixel compression type", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompress, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "crop", Description: "crop image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "debug", Description: "debug output", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "define", Description: "format-specific option", Category: CategorySetting, Type: TypeValue, ValueType: ValueDefine, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "density", Description: "horizontal/vertical resolution", Category: CategorySetting, Type: TypeValue, ValueType: ValueDensity, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "depth", Description: "image depth", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "deskew", Description: "deskew image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "despeckle", Description: "despeckle image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "dither", Description: "dithering method", Category: CategorySetting, Type: TypeValue, ValueType: ValueMethod, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "draw", Description: "annotate image with primitive", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "equalize", Description: "histogram equalization", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "evaluate", Description: "evaluate arithmetic expression", Category: CategoryOperator, Type: TypeValue, ValueType: ValueEvaluate},
		{Name: "extent", Description: "set image extent", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "fill", Description: "fill color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "flip", Description: "flip vertically", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "flop", Description: "flop horizontally", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "font", Description: "text font", Category: CategorySetting, Type: TypeValue, ValueType: ValueFont, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "format", Description: "write output in this format", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "fuzz", Description: "color distance threshold", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "gamma", Description: "gamma correction", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFloat},
		{Name: "gaussian-blur", Description: "Gaussian blur", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "geometry", Description: "preferred tile/border size", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "gravity", Description: "text placement direction", Category: CategorySetting, Type: TypeValue, ValueType: ValueGravity, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "help", Description: "print program options", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "interlace", Description: "interlacing scheme", Category: CategorySetting, Type: TypeValue, ValueType: ValueInterlace, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "label", Description: "image label", Category: CategorySetting, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "level", Description: "level adjustment", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "limit", Description: "resource limit", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "list", Description: "print supported values", Category: CategoryMisc, Type: TypeValue, ValueType: ValueList},
		{Name: "log", Description: "debug log format", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString},
		{Name: "monitor", Description: "progress monitor", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "negate", Description: "negate image", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "normalize", Description: "normalize image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "opaque", Description: "change color to fill color", Category: CategoryOperator, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "orient", Description: "image orientation", Category: CategorySetting, Type: TypeValue, ValueType: ValueOrientation},
		{Name: "path", Description: "write output files to this directory", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilename},
		{Name: "pointsize", Description: "font point size", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "profile", Description: "ICC/IOM profile", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "quality", Description: "compression level", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "quiet", Description: "suppress warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "regard-warnings", Description: "pay attention to warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "resize", Description: "resize image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "rotate", Description: "rotate image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueDegrees, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "sample", Description: "scale with pixel sampling", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "scale", Description: "scale image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "seed", Description: "random seed", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "separate", Description: "separate channels", Category: CategoryChannelOp, Type: TypeBoolean},
		{Name: "set", Description: "set image property", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "sharpen", Description: "sharpen image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "shave", Description: "shave pixels from edges", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "strip", Description: "strip image profiles/comments", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "stroke", Description: "stroke color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "strokewidth", Description: "stroke width", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "threshold", Description: "threshold image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "thumbnail", Description: "create thumbnail", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "trim", Description: "trim image edges", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "type", Description: "image type", Category: CategorySetting, Type: TypeValue, ValueType: ValueImageType, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "verbose", Description: "print detailed information", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "version", Description: "print version", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "virtual-pixel", Description: "virtual pixel method", Category: CategorySetting, Type: TypeValue, ValueType: ValueVirtualPixel},
		{Name: "white-threshold", Description: "white threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
	}
	return buildIndexFromOptions(options)
}

func buildCompareOptionIndex() map[string]*OptionDef {
	options := []*OptionDef{
		{Name: "alpha", Description: "control alpha/matte channel", Category: CategorySetting, Type: TypeValue, ValueType: ValueAlpha, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "channel", Description: "set active channel mask", Category: CategoryChannelOp, Type: TypeValue, ValueType: ValueChannel, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "colorspace", Description: "alternate colorspace", Category: CategorySetting, Type: TypeValue, ValueType: ValueColorspace},
		{Name: "compose", Description: "composite operator", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompose},
		{Name: "compress", Description: "pixel compression type", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompress},
		{Name: "debug", Description: "debug output", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "define", Description: "format-specific option", Category: CategorySetting, Type: TypeValue, ValueType: ValueDefine, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "density", Description: "horizontal/vertical resolution", Category: CategorySetting, Type: TypeValue, ValueType: ValueDensity},
		{Name: "depth", Description: "image depth", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "dissimilarity-threshold", Description: "maximum dissimilarity for match", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "highlight-color", Description: "color for differing pixels", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor},
		{Name: "help", Description: "print program options", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "interlace", Description: "interlacing scheme", Category: CategorySetting, Type: TypeValue, ValueType: ValueInterlace},
		{Name: "limit", Description: "resource limit", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "list", Description: "print supported values", Category: CategoryMisc, Type: TypeValue, ValueType: ValueList},
		{Name: "log", Description: "debug log format", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString},
		{Name: "metric", Description: "comparison metric", Category: CategorySetting, Type: TypeValue, ValueType: ValueMetric},
		{Name: "monitor", Description: "progress monitor", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "profile", Description: "ICC/IOM profile", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "quality", Description: "compression level", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "quiet", Description: "suppress warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "regard-warnings", Description: "pay attention to warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "seed", Description: "random seed", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "similarity-threshold", Description: "minimum similarity for match", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "size", Description: "image dimensions", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "subimage-search", Description: "search for subimage", Category: CategorySetting, Type: TypeBoolean},
		{Name: "transparent", Description: "make color transparent", Category: CategoryOperator, Type: TypeValue, ValueType: ValueColor},
		{Name: "type", Description: "image type", Category: CategorySetting, Type: TypeValue, ValueType: ValueImageType},
		{Name: "verbose", Description: "print detailed information", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "version", Description: "print version", Category: CategoryMisc, Type: TypeBoolean},
	}
	return buildIndexFromOptions(options)
}

func buildCompositeOptionIndex() map[string]*OptionDef {
	options := []*OptionDef{
		{Name: "alpha", Description: "control alpha/matte channel", Category: CategorySetting, Type: TypeValue, ValueType: ValueAlpha, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "background", Description: "background color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "blend", Description: "blend percentages", Category: CategorySetting, Type: TypeValue, ValueType: ValueBlend},
		{Name: "border", Description: "surround image with border", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "channel", Description: "set active channel mask", Category: CategoryChannelOp, Type: TypeValue, ValueType: ValueChannel, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "colorspace", Description: "alternate colorspace", Category: CategorySetting, Type: TypeValue, ValueType: ValueColorspace},
		{Name: "compose", Description: "composite operator", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompose},
		{Name: "compress", Description: "pixel compression type", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompress},
		{Name: "debug", Description: "debug output", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "define", Description: "format-specific option", Category: CategorySetting, Type: TypeValue, ValueType: ValueDefine, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "density", Description: "horizontal/vertical resolution", Category: CategorySetting, Type: TypeValue, ValueType: ValueDensity},
		{Name: "depth", Description: "image depth", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "displace", Description: "shift image according to displacement map", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "dissolve", Description: "dissolve percentage", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "geometry", Description: "preferred tile/border size", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "gravity", Description: "text placement direction", Category: CategorySetting, Type: TypeValue, ValueType: ValueGravity, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "help", Description: "print program options", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "identify", Description: "identify image format and attributes", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "interlace", Description: "interlacing scheme", Category: CategorySetting, Type: TypeValue, ValueType: ValueInterlace},
		{Name: "limit", Description: "resource limit", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "list", Description: "print supported values", Category: CategoryMisc, Type: TypeValue, ValueType: ValueList},
		{Name: "log", Description: "debug log format", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString},
		{Name: "monitor", Description: "progress monitor", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "negate", Description: "negate image", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "page", Description: "page geometry", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "profile", Description: "ICC/IOM profile", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "quality", Description: "compression level", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "quiet", Description: "suppress warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "regard-warnings", Description: "pay attention to warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "rotate", Description: "rotate image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueDegrees},
		{Name: "seed", Description: "random seed", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "size", Description: "image dimensions", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "stegano", Description: "hide watermark at offset", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "strip", Description: "strip image profiles/comments", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "type", Description: "image type", Category: CategorySetting, Type: TypeValue, ValueType: ValueImageType},
		{Name: "verbose", Description: "print detailed information", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "version", Description: "print version", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "virtual-pixel", Description: "virtual pixel method", Category: CategorySetting, Type: TypeValue, ValueType: ValueVirtualPixel},
		{Name: "watermark", Description: "watermark", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "write", Description: "write image to file", Category: CategorySequenceOp, Type: TypeValue, ValueType: ValueFilename},
	}
	return buildIndexFromOptions(options)
}

func buildMontageOptionIndex() map[string]*OptionDef {
	options := []*OptionDef{
		{Name: "alpha", Description: "control alpha/matte channel", Category: CategorySetting, Type: TypeValue, ValueType: ValueAlpha, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "background", Description: "background color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "bordercolor", Description: "border color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "borderwidth", Description: "border width", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "channel", Description: "set active channel mask", Category: CategoryChannelOp, Type: TypeValue, ValueType: ValueChannel, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "colorspace", Description: "alternate colorspace", Category: CategorySetting, Type: TypeValue, ValueType: ValueColorspace},
		{Name: "comment", Description: "annotate image with comment", Category: CategorySetting, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "compose", Description: "composite operator", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompose},
		{Name: "compress", Description: "pixel compression type", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompress},
		{Name: "crop", Description: "crop image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "debug", Description: "debug output", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "define", Description: "format-specific option", Category: CategorySetting, Type: TypeValue, ValueType: ValueDefine, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "density", Description: "horizontal/vertical resolution", Category: CategorySetting, Type: TypeValue, ValueType: ValueDensity},
		{Name: "depth", Description: "image depth", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "fill", Description: "fill color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "font", Description: "text font", Category: CategorySetting, Type: TypeValue, ValueType: ValueFont, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "frame", Description: "surround image with ornamental border", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "geometry", Description: "preferred tile/border size", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "gravity", Description: "text placement direction", Category: CategorySetting, Type: TypeValue, ValueType: ValueGravity, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "help", Description: "print program options", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "interlace", Description: "interlacing scheme", Category: CategorySetting, Type: TypeValue, ValueType: ValueInterlace},
		{Name: "label", Description: "image label", Category: CategorySetting, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "limit", Description: "resource limit", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "list", Description: "print supported values", Category: CategoryMisc, Type: TypeValue, ValueType: ValueList},
		{Name: "log", Description: "debug log format", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString},
		{Name: "mattecolor", Description: "frame color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor},
		{Name: "mode", Description: "mode of operation", Category: CategoryOperator, Type: TypeValue, ValueType: ValueMethod},
		{Name: "monitor", Description: "progress monitor", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "page", Description: "page geometry", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "pointsize", Description: "font point size", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "profile", Description: "ICC/IOM profile", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "quality", Description: "compression level", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "quiet", Description: "suppress warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "regard-warnings", Description: "pay attention to warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "seed", Description: "random seed", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "shadow", Description: "add drop shadow to tiles", Category: CategorySetting, Type: TypeBoolean},
		{Name: "size", Description: "image dimensions", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "strip", Description: "strip image profiles/comments", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "stroke", Description: "stroke color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "texture", Description: "background texture", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilename},
		{Name: "tile", Description: "tile geometry", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "title", Description: "montage title", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "transparent", Description: "make color transparent", Category: CategoryOperator, Type: TypeValue, ValueType: ValueColor},
		{Name: "trim", Description: "trim image edges", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "type", Description: "image type", Category: CategorySetting, Type: TypeValue, ValueType: ValueImageType},
		{Name: "verbose", Description: "print detailed information", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "version", Description: "print version", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "virtual-pixel", Description: "virtual pixel method", Category: CategorySetting, Type: TypeValue, ValueType: ValueVirtualPixel},
	}
	return buildIndexFromOptions(options)
}
