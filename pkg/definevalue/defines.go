package definevalue

// DefineKey describes a format-specific define key.
type DefineKey struct {
	Name        string
	Description string
	ValueType   string // "int", "boolean", "string", "float"
	Values      []string
}

// FormatDefines maps format names to their supported define keys.
var FormatDefines = map[string][]DefineKey{
	"jpeg": {
		{Name: "quality", Description: "compression quality", ValueType: "int"},
		{Name: "optimize-codes", Description: "optimize Huffman coding tables", ValueType: "boolean"},
		{Name: "progressive", Description: "write progressive JPEG", ValueType: "boolean"},
		{Name: "sampling-factor", Description: "chroma subsampling", ValueType: "string"},
		{Name: "extent", Description: "maximum file size in KB", ValueType: "int"},
		{Name: "arithmetic-coding", Description: "use arithmetic coding", ValueType: "boolean"},
	},
	"png": {
		{Name: "compression-level", Description: "zlib compression level", ValueType: "int", Values: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}},
		{Name: "compression-filter", Description: "row filter method", ValueType: "int", Values: []string{"0", "1", "2", "3", "4", "5"}},
		{Name: "compression-strategy", Description: "zlib strategy", ValueType: "int", Values: []string{"0", "1", "2", "3"}},
		{Name: "color-type", Description: "force color type", ValueType: "int", Values: []string{"0", "2", "3", "4", "6"}},
		{Name: "bit-depth", Description: "force bit depth", ValueType: "int"},
		{Name: "exclude-chunk", Description: "exclude chunks from output", ValueType: "string"},
		{Name: "include-chunk", Description: "include specific chunks", ValueType: "string"},
	},
	"tiff": {
		{Name: "compress", Description: "compression method", ValueType: "string", Values: []string{"none", "lzw", "zip", "jpeg", "lzma", "piz", "pxz24", "pzip"}},
		{Name: "rows-per-strip", Description: "rows per strip", ValueType: "int"},
		{Name: "tile-geometry", Description: "tiled TIFF layout", ValueType: "string"},
		{Name: "predictor", Description: "predictor for compression", ValueType: "int"},
		{Name: "fill-order", Description: "fill order for bits", ValueType: "string"},
	},
	"webp": {
		{Name: "lossless", Description: "use lossless encoding", ValueType: "boolean"},
		{Name: "quality", Description: "lossy quality", ValueType: "float"},
		{Name: "method", Description: "compression method", ValueType: "int", Values: []string{"0", "1", "2", "3", "4", "5", "6"}},
		{Name: "alpha-quality", Description: "alpha channel quality", ValueType: "float"},
		{Name: "target-size", Description: "target file size in KB", ValueType: "int"},
		{Name: "target-psnr", Description: "target PSNR in dB", ValueType: "float"},
	},
	"gif": {
		{Name: "interlace", Description: "interlaced GIF", ValueType: "boolean"},
		{Name: "optimize", Description: "optimize animation frames", ValueType: "boolean"},
		{Name: "disposal", Description: "frame disposal method", ValueType: "string"},
	},
	"heic": {
		{Name: "quality", Description: "compression quality", ValueType: "int"},
		{Name: "lossless", Description: "lossless encoding", ValueType: "boolean"},
		{Name: "speed", Description: "encoding speed", ValueType: "int", Values: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"}},
	},
	"psd": {
		{Name: "additional-info", Description: "preserve additional info", ValueType: "boolean"},
	},
	"pdf": {
		{Name: "fit-page", Description: "fit page to image dimensions", ValueType: "boolean"},
		{Name: "use-cropbox", Description: "use crop box", ValueType: "boolean"},
		{Name: "page", Description: "page number", ValueType: "int"},
	},
	"raw": {
		{Name: "demosaic-algorithm", Description: "demosaicing algorithm", ValueType: "string"},
	},
}

// GlobalDefines lists define keys that don't require a format prefix.
var GlobalDefines = []DefineKey{
	{Name: "optimize", Description: "general optimization", ValueType: "boolean"},
	{Name: "type", Description: "hint for output image type", ValueType: "string"},
	{Name: "preserve-colorspace", Description: "preserve input colorspace", ValueType: "boolean"},
	{Name: "bypass-picture-cache", Description: "bypass the picture cache", ValueType: "boolean"},
	{Name: "dither", Description: "dithering method", ValueType: "string"},
	{Name: "exif:ignore", Description: "ignore EXIF data", ValueType: "boolean"},
	{Name: "stream:buffer-size", Description: "stream buffer size", ValueType: "int"},
}

// LookupFormatDefines returns the define keys for a given format.
func LookupFormatDefines(format string) []DefineKey {
	return FormatDefines[format]
}

// LookupGlobalDefines returns global define keys.
func LookupGlobalDefines() []DefineKey {
	return GlobalDefines
}

// LookupDefineKey finds a specific define key within a format.
func LookupDefineKey(format, key string) *DefineKey {
	keys := FormatDefines[format]
	for i := range keys {
		if keys[i].Name == key {
			return &keys[i]
		}
	}
	// Check global defines
	for i := range GlobalDefines {
		if GlobalDefines[i].Name == key {
			return &GlobalDefines[i]
		}
	}
	return nil
}
