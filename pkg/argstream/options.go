package argstream

// OptionCategory defines where an option applies in the magick pipeline.
type OptionCategory int

const (
	CategorySetting OptionCategory = iota
	CategoryOperator
	CategoryChannelOp
	CategorySequenceOp
	CategoryStackOp
	CategoryMisc
)

func (c OptionCategory) String() string {
	switch c {
	case CategorySetting:
		return "Setting"
	case CategoryOperator:
		return "Operator"
	case CategoryChannelOp:
		return "ChannelOp"
	case CategorySequenceOp:
		return "SequenceOp"
	case CategoryStackOp:
		return "StackOp"
	case CategoryMisc:
		return "Misc"
	}
	return "Unknown"
}

// OptionType defines whether an option takes a value.
type OptionType int

const (
	TypeBoolean OptionType = iota
	TypeValue
)

// PlusBehavior defines what the +form of an option does.
type PlusBehavior int

const (
	PlusReset       PlusBehavior = iota // +option resets to default (no value)
	PlusInverse                         // +option is the inverse/opposite
	PlusDirectional                     // +option is alternate direction (may take value)
)

// ValueType defines the type of value an option expects.
type ValueType string

const (
	ValueString        ValueType = "string"
	ValueInt           ValueType = "int"
	ValueFloat         ValueType = "float"
	ValueDegrees       ValueType = "degrees"
	ValueGeometry      ValueType = "geometry"
	ValueColor         ValueType = "color"
	ValueColorspace    ValueType = "colorspace"
	ValueCompose       ValueType = "compose"
	ValueCompress      ValueType = "compress"
	ValueChannel       ValueType = "channel"
	ValueDistort       ValueType = "distort"
	ValueFilter        ValueType = "filter"
	ValueGravity       ValueType = "gravity"
	ValueInterlace     ValueType = "interlace"
	ValueLayers        ValueType = "layers"
	ValueMorphology    ValueType = "morphology"
	ValueImageType     ValueType = "type"
	ValueVirtualPixel  ValueType = "virtual_pixel"
	ValueMetric        ValueType = "metric"
	ValueEvaluate      ValueType = "evaluate"
	ValueDefine        ValueType = "define"
	ValueBoolean       ValueType = "boolean"
	ValueFont          ValueType = "font"
	ValueFormat        ValueType = "format"
	ValueFilename      ValueType = "filename"
	ValueOrientation   ValueType = "orientation"
	ValueDispose       ValueType = "dispose"
	ValueAlpha         ValueType = "alpha"
	ValueNoise         ValueType = "noise"
	ValuePreview       ValueType = "preview"
	ValueStorage       ValueType = "storage"
	ValueKernel        ValueType = "kernel"
	ValueList          ValueType = "list"
	ValueDensity       ValueType = "density"
	ValuePoint         ValueType = "point"
	ValueRatio         ValueType = "ratio"
	ValueThreshold     ValueType = "threshold"
	ValueExpression    ValueType = "expression"
	ValueIndex         ValueType = "index"
	ValueCount         ValueType = "count"
	ValueMethod        ValueType = "method"
	ValueAutoThreshold ValueType = "auto_threshold"
	ValueBlend         ValueType = "blend"
	ValueDirection     ValueType = "direction"
)

// OptionDef defines a single magick option.
type OptionDef struct {
	Name         string
	Description  string
	Category     OptionCategory
	Type         OptionType
	ValueType    ValueType
	HasPlusForm  bool
	PlusBehavior PlusBehavior
}

// OptionIndex maps option names to their definitions.
var OptionIndex map[string]*OptionDef

func init() {
	OptionIndex = buildOptionIndex()
}

func buildOptionIndex() map[string]*OptionDef {
	options := []*OptionDef{
		// Image Settings
		{Name: "adaptive-resize", Description: "resize adaptively", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry, HasPlusForm: false},
		{Name: "adjoin", Description: "join images into a single multi-image file", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "alpha", Description: "control alpha/matte channel", Category: CategorySetting, Type: TypeValue, ValueType: ValueAlpha, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "antialias", Description: "remove pixel aliasing", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "authenticate", Description: "password for encrypted images", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "auto-orient", Description: "automatically orient image", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "auto-threshold", Description: "auto threshold method", Category: CategoryOperator, Type: TypeValue, ValueType: ValueAutoThreshold},
		{Name: "background", Description: "background color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "bias", Description: "add bias to convolve", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "black-point-compensation", Description: "black point compensation", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "black-threshold", Description: "black threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "blend", Description: "blend percentages", Category: CategorySetting, Type: TypeValue, ValueType: ValueBlend},
		{Name: "blue-primary", Description: "blue primary point", Category: CategorySetting, Type: TypeValue, ValueType: ValuePoint},
		{Name: "bordercolor", Description: "border color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "caption", Description: "caption text", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "channel", Description: "set active channel mask", Category: CategoryChannelOp, Type: TypeValue, ValueType: ValueChannel, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "colors", Description: "preferred number of colors", Category: CategoryOperator, Type: TypeValue, ValueType: ValueInt},
		{Name: "colorspace", Description: "alternate colorspace", Category: CategorySetting, Type: TypeValue, ValueType: ValueColorspace, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "comment", Description: "annotate image with comment", Category: CategorySetting, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "compose", Description: "composite operator", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompose},
		{Name: "compress", Description: "pixel compression type", Category: CategorySetting, Type: TypeValue, ValueType: ValueCompress, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "convolve", Description: "apply convolution kernel", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "crop", Description: "crop image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "cycle", Description: "cycle image colormap", Category: CategoryOperator, Type: TypeValue, ValueType: ValueInt},
		{Name: "debug", Description: "debug output", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "define", Description: "format-specific option", Category: CategorySetting, Type: TypeValue, ValueType: ValueDefine, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "delay", Description: "display delay (1/100 second)", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "density", Description: "horizontal/vertical resolution", Category: CategorySetting, Type: TypeValue, ValueType: ValueDensity, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "depth", Description: "image depth", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "direction", Description: "text rendering direction", Category: CategorySetting, Type: TypeValue, ValueType: ValueDirection},
		{Name: "display", Description: "X11 display", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "dispose", Description: "disposal method", Category: CategorySetting, Type: TypeValue, ValueType: ValueDispose},
		{Name: "dissimilarity-threshold", Description: "maximum dissimilarity for match", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "dither", Description: "dithering method", Category: CategorySetting, Type: TypeValue, ValueType: ValueMethod, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "draw", Description: "annotate image with primitive", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "encoding", Description: "text encoding type", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "endian", Description: "endianness", Category: CategorySetting, Type: TypeValue, ValueType: ValueMethod},
		{Name: "equalize", Description: "histogram equalization", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "evaluate", Description: "evaluate arithmetic expression", Category: CategoryOperator, Type: TypeValue, ValueType: ValueEvaluate},
		{Name: "extent", Description: "set image extent", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "fill", Description: "fill color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "filter", Description: "resampling filter", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilter},
		{Name: "flatten", Description: "flatten sequence", Category: CategorySequenceOp, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "flip", Description: "flip vertically", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "flop", Description: "flop horizontally", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "font", Description: "text font", Category: CategorySetting, Type: TypeValue, ValueType: ValueFont, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "format", Description: "output format string", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "frame", Description: "surround image with ornamental border", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "fuzz", Description: "color distance threshold", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "gamma", Description: "gamma correction", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFloat},
		{Name: "gaussian-blur", Description: "Gaussian blur", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "gravity", Description: "text placement direction", Category: CategorySetting, Type: TypeValue, ValueType: ValueGravity, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "green-primary", Description: "green primary point", Category: CategorySetting, Type: TypeValue, ValueType: ValuePoint},
		{Name: "help", Description: "print program options", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "identify", Description: "identify image format and attributes", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "immutable", Description: "make image immutable", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "interlace", Description: "interlacing scheme", Category: CategorySetting, Type: TypeValue, ValueType: ValueInterlace, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "interpolate", Description: "interpolation method", Category: CategorySetting, Type: TypeValue, ValueType: ValueMethod},
		{Name: "kernel", Description: "convolution kernel", Category: CategoryOperator, Type: TypeValue, ValueType: ValueKernel},
		{Name: "label", Description: "image label", Category: CategorySetting, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "lat", Description: "local adaptive threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "layers", Description: "layer method", Category: CategorySequenceOp, Type: TypeValue, ValueType: ValueLayers},
		{Name: "level", Description: "level adjustment", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "limit", Description: "resource limit", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "linear-stretch", Description: "linear stretch", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "list", Description: "print supported values", Category: CategoryMisc, Type: TypeValue, ValueType: ValueList},
		{Name: "log", Description: "debug log format", Category: CategoryMisc, Type: TypeValue, ValueType: ValueString},
		{Name: "loop", Description: "animation loop count", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "matte", Description: "store matte channel", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "metric", Description: "comparison metric", Category: CategorySetting, Type: TypeValue, ValueType: ValueMetric},
		{Name: "mode", Description: "mode of operation", Category: CategoryOperator, Type: TypeValue, ValueType: ValueMethod},
		{Name: "monitor", Description: "progress monitor", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "morph", Description: "morph image sequence", Category: CategorySequenceOp, Type: TypeValue, ValueType: ValueInt},
		{Name: "morphology", Description: "morphology method", Category: CategoryOperator, Type: TypeValue, ValueType: ValueMorphology},
		{Name: "mosaic", Description: "mosaic image sequence", Category: CategorySequenceOp, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "negate", Description: "negate image", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "noise", Description: "add/reduce noise", Category: CategoryOperator, Type: TypeValue, ValueType: ValueNoise, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "normalize", Description: "normalize image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "opaque", Description: "change color to fill color", Category: CategoryOperator, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "ordered-dither", Description: "ordered dither", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "orient", Description: "image orientation", Category: CategorySetting, Type: TypeValue, ValueType: ValueOrientation},
		{Name: "page", Description: "page geometry", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "path", Description: "write output files to this directory", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilename},
		{Name: "ping", Description: "efficiently determine attributes", Category: CategorySetting, Type: TypeBoolean},
		{Name: "pointsize", Description: "font point size", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "posterize", Description: "posterize image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueInt},
		{Name: "precision", Description: "print precision digits", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "preview", Description: "preview type", Category: CategoryOperator, Type: TypeValue, ValueType: ValuePreview},
		{Name: "process", Description: "custom filter", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "profile", Description: "ICC/IOM profile", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "quality", Description: "compression level", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "quiet", Description: "suppress warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "read", Description: "explicitly read image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "regard-warnings", Description: "pay attention to warnings", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "region", Description: "restrict operations to region", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "remap", Description: "remap image colors", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "render", Description: "render vector operations", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "repage", Description: "page geometry", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "resample", Description: "resample to resolution", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "resize", Description: "resize image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "respect-parentheses", Description: "respect parentheses settings", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "reverse", Description: "reverse image sequence", Category: CategoryStackOp, Type: TypeBoolean},
		{Name: "rotate", Description: "rotate image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueDegrees, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "sample", Description: "scale with pixel sampling", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "sampling-factor", Description: "sampling factor", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "scale", Description: "scale image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "scene", Description: "scene number", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "screen", Description: "screen capture", Category: CategorySetting, Type: TypeBoolean},
		{Name: "seed", Description: "random seed", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "separate", Description: "separate channels", Category: CategoryChannelOp, Type: TypeBoolean},
		{Name: "sepia-tone", Description: "simulate sepia-toned photo", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "set", Description: "set image property", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "shade", Description: "shade image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "shadow", Description: "add drop shadow", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "sharpen", Description: "sharpen image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "shave", Description: "shave pixels from edges", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "shear", Description: "shear image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueDegrees},
		{Name: "sigmoidal-contrast", Description: "sigmoidal contrast adjustment", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "similarity-threshold", Description: "minimum similarity for match", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "size", Description: "image dimensions", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "sketch", Description: "sketch image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "solarize", Description: "solarize image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "sparse-color", Description: "sparse color method", Category: CategoryOperator, Type: TypeValue, ValueType: ValueMethod},
		{Name: "stegano", Description: "hide watermark at offset", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "stereo", Description: "stereo image composite", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "strip", Description: "strip image profiles/comments", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "stroke", Description: "stroke color", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "strokewidth", Description: "stroke width", Category: CategorySetting, Type: TypeValue, ValueType: ValueFloat},
		{Name: "style", Description: "font style", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "swap", Description: "swap two images", Category: CategoryStackOp, Type: TypeValue, ValueType: ValueIndex},
		{Name: "swirl", Description: "swirl image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueDegrees},
		{Name: "synchronize", Description: "synchronize image to storage", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "taint", Description: "declare image as modified", Category: CategorySetting, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "texture", Description: "background texture", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilename},
		{Name: "threshold", Description: "threshold image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "thumbnail", Description: "create thumbnail", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "tile", Description: "tile geometry", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "tile-offset", Description: "tile offset", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "tint", Description: "tint image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "transform", Description: "affine transform", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "transparent", Description: "make color transparent", Category: CategoryOperator, Type: TypeValue, ValueType: ValueColor, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "transpose", Description: "transpose image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "transverse", Description: "transverse image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "treedepth", Description: "color tree depth", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "trim", Description: "trim image edges", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "type", Description: "image type", Category: CategorySetting, Type: TypeValue, ValueType: ValueImageType, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "unique-colors", Description: "discard all but one of each pixel color", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "units", Description: "resolution units", Category: CategorySetting, Type: TypeValue, ValueType: ValueMethod},
		{Name: "unsharp", Description: "unsharp mask", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "update", Description: "detect when image file is modified", Category: CategoryMisc, Type: TypeValue, ValueType: ValueInt},
		{Name: "usage", Description: "print usage", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "verbose", Description: "print detailed information", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "version", Description: "print version", Category: CategoryMisc, Type: TypeBoolean},
		{Name: "virtual-pixel", Description: "virtual pixel method", Category: CategorySetting, Type: TypeValue, ValueType: ValueVirtualPixel},
		{Name: "watermark", Description: "watermark", Category: CategorySetting, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "wave", Description: "wave distortion", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "wavelet-denoise", Description: "wavelet denoise", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "weight", Description: "font weight", Category: CategorySetting, Type: TypeValue, ValueType: ValueString},
		{Name: "white-balance", Description: "auto white balance", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "white-point", Description: "chromaticity white point", Category: CategorySetting, Type: TypeValue, ValueType: ValuePoint},
		{Name: "white-threshold", Description: "white threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "write", Description: "write image to file", Category: CategorySequenceOp, Type: TypeValue, ValueType: ValueFilename, HasPlusForm: true, PlusBehavior: PlusDirectional},

		// Stack operators
		{Name: "clone", Description: "clone image(s)", Category: CategoryStackOp, Type: TypeValue, ValueType: ValueIndex},
		{Name: "delete", Description: "delete image(s)", Category: CategoryStackOp, Type: TypeValue, ValueType: ValueIndex},
		{Name: "duplicate", Description: "duplicate image(s)", Category: CategoryStackOp, Type: TypeValue, ValueType: ValueCount},
		{Name: "insert", Description: "insert last image at position", Category: CategoryStackOp, Type: TypeValue, ValueType: ValueIndex},

		// Special operators
		{Name: "append", Description: "append image sequence", Category: CategorySequenceOp, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "coalesce", Description: "merge image sequence", Category: CategorySequenceOp, Type: TypeBoolean},
		{Name: "combine", Description: "combine into color channels", Category: CategorySequenceOp, Type: TypeBoolean},
		{Name: "composite", Description: "composite images", Category: CategorySequenceOp, Type: TypeBoolean},
		{Name: "distort", Description: "distort image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueDistort},
		{Name: "channel-fx", Description: "channel operations", Category: CategoryChannelOp, Type: TypeValue, ValueType: ValueExpression},
		{Name: "fx", Description: "apply math expression", Category: CategorySequenceOp, Type: TypeValue, ValueType: ValueExpression},
		{Name: "function", Description: "apply math function", Category: CategoryOperator, Type: TypeValue, ValueType: ValueMethod},
		{Name: "annotate", Description: "annotate image with text", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "adaptive-blur", Description: "adaptive blur", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "adaptive-sharpen", Description: "adaptive sharpen", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "blur", Description: "blur image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "border", Description: "surround image with border", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "charcoal", Description: "charcoal drawing simulation", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "chop", Description: "chop image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "colorize", Description: "colorize image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "color-threshold", Description: "color threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "contrast", Description: "enhance image contrast", Category: CategoryOperator, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "contrast-stretch", Description: "contrast stretch", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "decipher", Description: "decipher image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "deskew", Description: "deskew image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueThreshold},
		{Name: "despeckle", Description: "despeckle image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "encipher", Description: "encipher image", Category: CategoryOperator, Type: TypeValue, ValueType: ValueFilename},
		{Name: "enhance", Description: "enhance image", Category: CategoryOperator, Type: TypeBoolean},
		{Name: "highlight-color", Description: "color for differing pixels", Category: CategorySetting, Type: TypeValue, ValueType: ValueColor},
		{Name: "hough-lines", Description: "hough line detection", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "mean-shift", Description: "mean shift segmentation", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "motion-blur", Description: "motion blur", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "polaroid", Description: "polaroid effect", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "raise", Description: "raise/lower image edges", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry, HasPlusForm: true, PlusBehavior: PlusDirectional},
		{Name: "random-threshold", Description: "random threshold", Category: CategoryOperator, Type: TypeValue, ValueType: ValueString},
		{Name: "read-mask", Description: "read mask", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilename, HasPlusForm: true, PlusBehavior: PlusReset},
		{Name: "red-primary", Description: "red primary point", Category: CategorySetting, Type: TypeValue, ValueType: ValuePoint},
		{Name: "selective-blur", Description: "selective blur", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "spread", Description: "displace image pixels", Category: CategoryOperator, Type: TypeValue, ValueType: ValueInt},
		{Name: "statistic", Description: "replace pixels with statistic", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "vignette", Description: "vignette effect", Category: CategoryOperator, Type: TypeValue, ValueType: ValueGeometry},
		{Name: "write-mask", Description: "write mask", Category: CategorySetting, Type: TypeValue, ValueType: ValueFilename, HasPlusForm: true, PlusBehavior: PlusReset},

		// Special settings
		{Name: "concurrent", Description: "concurrent processing", Category: CategoryMisc, Type: TypeBoolean, HasPlusForm: true, PlusBehavior: PlusInverse},
		{Name: "dissolve", Description: "dissolve percentage", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "features", Description: "analyze image features", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
		{Name: "moments", Description: "report image moments", Category: CategorySetting, Type: TypeBoolean},
		{Name: "unique", Description: "display number of unique colors", Category: CategorySetting, Type: TypeBoolean},
		{Name: "borderwidth", Description: "border width", Category: CategorySetting, Type: TypeValue, ValueType: ValueInt},
	}

	return buildIndexFromOptions(options)
}

func buildIndexFromOptions(options []*OptionDef) map[string]*OptionDef {
	index := make(map[string]*OptionDef)
	for _, opt := range options {
		index[opt.Name] = opt
	}
	return index
}

// LookupOption looks up an option by name (without leading - or +).
func LookupOption(name string) *OptionDef {
	if opt, ok := OptionIndex[name]; ok {
		return opt
	}
	return nil
}
