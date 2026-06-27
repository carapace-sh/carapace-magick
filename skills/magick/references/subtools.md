# Sub-Tools — identify, mogrify, compare, composite, montage, conjure, stream

How `magick`'s sub-tools differ from the main `magick` command: which options they support, their argument structure, and their pipeline model.

> **Source of truth**: <https://imagemagick.org/script/command-line-tools.php> and each tool's `-help` output.

## Overview

In ImageMagick 7, all tools are invoked as sub-commands of `magick`:

```
magick tool [options ...] [arguments ...]
```

Legacy standalone commands (`convert`, `identify`, etc.) still work but print a deprecation warning. The `magick` command without a tool name defaults to the `convert` behavior (image pipeline processing).

| Tool | Purpose | Usage Line |
|------|---------|------------|
| (default/convert) | Image pipeline processing | `magick [options \| image ...] output_image` |
| `identify` | Describe image format and attributes | `magick identify [options ...] file [...]` |
| `mogrify` | In-place image transformation | `magick mogrify [options ...] file [...]` |
| `compare` | Assess difference between images | `magick compare [options ...] image reconstruct difference` |
| `composite` | Composite images together | `magick composite [options ...] image composite [mask] composite` |
| `montage` | Create a composite image montage | `magick montage [options ...] file [...] file` |
| `conjure` | Execute MSL scripts | `magick conjure [options ...] file [...]` |
| `stream` | Stream raw pixel data | `magick stream [options ...] input-image raw-image` |
| `animate` | Display image animation | (X11 display) |
| `display` | Interactive image viewer | (X11 display) |
| `import` | Screen capture | (X11 display) |

## Common vs. Tool-Specific Options

All tools share a **common subset** of options (settings like `-verbose`, `-debug`, `-list`), but each tool supports only a subset of the full option set. This is important for completion — the lexer must know which options are valid for each tool.

### Shared Options (all tools)

| Option | Description |
|--------|-------------|
| `-debug events` | Debug output |
| `-help` | Print options |
| `-list type` | Print supported values |
| `-log format` | Debug log format |
| `-usage` | Print usage |
| `-version` | Print version |
| `-verbose` | Detailed output |
| `-quiet` | Suppress warnings |
| `-regard-warnings` | Pay attention to warnings |
| `-seed value` | Random seed |
| `-monitor` | Progress monitoring |

## `magick` (Default / Convert)

The full image pipeline processor. Supports the complete option set: settings, operators, sequence operators, stack operators, and parentheses.

```
magick [ {option} | {image} ... ] {output_image}
```

Key characteristics:
- **Full pipeline**: supports all option categories
- **Image stack**: `-clone`, `-delete`, `-swap`, parentheses
- **Sequence operators**: `-append`, `-flatten`, `-layers`, `-composite`
- **Output required**: must end with an output filename

## `magick identify`

Read-only image inspector. Reports image format, dimensions, color depth, etc.

```
magick identify [options ...] file [ [options ...] file ...]
```

Key differences from `magick`:
- **No operators**: cannot transform images (no `-resize`, `-blur`, etc.)
- **No output argument**: reports to stdout, no output file
- **No image stack**: no parentheses, no `-clone`, no `-delete`
- **Read-only settings**: `-format`, `-verbose`, `-ping`, `-unique`, `-moments`, `-features`
- **Multiple input files**: can process several files in one invocation
- **`-format` string**: uses `%[...]` escape sequences to control output

### identify-specific options

| Option | Description |
|--------|-------------|
| `-ping` | Efficiently determine attributes (don't read full image) |
| `-unique` | Display number of unique colors |
| `-moments` | Report image moments |
| `-features distance` | Analyze image features |

## `magick mogrify`

In-place image transformer. Similar to `magick` but operates on existing files.

```
magick mogrify [options ...] file [ [options ...] file ...]
```

Key differences from `magick`:
- **In-place by default**: overwrites the input file
- **No image stack**: no parentheses, no `-clone`, no `-delete`
- **No output argument**: the input file is both source and destination
- **Full operator set**: supports most image operators (`-resize`, `-blur`, etc.)
- **Multiple files**: processes each file independently
- **`-format` for output format**: `-format png` converts to PNG while keeping the original filename base
- **`-path directory`**: write output to a different directory instead of overwriting

### mogrify-specific options

| Option | Description |
|--------|-------------|
| `-format type` | Write output in this format |
| `-path directory` | Write output files to this directory |

## `magick compare`

Compares two images and produces a visual difference image.

```
magick compare [options ...] image reconstruct difference
```

Key differences from `magick`:
- **Three positional arguments**: original image, reconstructed image, difference output
- **No image stack**: no parentheses or stack operators
- **Comparison-specific settings**: `-metric`, `-dissimilarity-threshold`, `-similarity-threshold`, `-highlight-color`, `-lowlight-color`
- **Limited operator set**: only settings, no image transformation operators

### compare-specific options

| Option | Description |
|--------|-------------|
| `-metric type` | Comparison metric (AE, MAE, MSE, PAE, PSNR, RMSE, SSIM, DSSIM, FUZZ) |
| `-dissimilarity-threshold` | Maximum dissimilarity for match |
| `-similarity-threshold` | Minimum similarity for match |
| `-highlight-color` | Color for differing pixels in visual diff |
| `-lowlight-color` | Color for similar pixels in visual diff |

## `magick composite`

Composites images together (overlay, blend, etc.).

```
magick composite [options ...] image overlay [ [options ...] mask ] [options ...] output
```

The positional arguments are: base `image`, `overlay` image to composite on top, optional `mask`, and `output` file. The official usage line uses "composite" for both the overlay and output arguments, which is confusing — here we use distinct names for clarity.

Key differences from `magick`:
- **Positional image arguments**: base image, overlay image, optional mask, output file
- **No image stack**: no parentheses or stack operators
- **No general operators**: only composite-related settings
- **`-compose` operator**: sets the compositing method (Over, Multiply, etc.)
- **`-geometry` for positioning**: places the composite image relative to the base
- **`-blend`**: blend percentage

### composite-specific options

| Option | Description |
|--------|-------------|
| `-blend geometry` | Blend percentages |
| `-displace geometry` | Shift image according to displacement map |
| `-dissolve value` | Dissolve percentage |
| `-stegano offset` | Hide watermark at offset |
| `-watermark geometry` | Watermark with given brightness/saturation |

## `magick montage`

Creates a grid of thumbnail images.

```
magick montage [options ...] file [ [options ...] file ...] file
```

Key differences from `magick`:
- **Montage-specific layout settings**: `-tile`, `-geometry`, `-frame`, `-borderwidth`
- **No general image operators**: limited to montage layout
- **Output file required**: last argument is the output

### montage-specific options

| Option | Description |
|--------|-------------|
| `-tile geometry` | Number of tiles per row/column (e.g., `4x3`) |
| `-geometry geometry` | Tile size and border |
| `-frame geometry` | Decorative frame around each tile |
| `-borderwidth value` | Border width between tiles |
| `-shadow` | Add drop shadow to tiles |
| `-texture filename` | Background texture |

## `magick conjure`

Executes MSL (Magick Scripting Language) scripts.

```
magick conjure [options ...] file [ [options ...] file ...]
```

Key differences:
- **Minimal fixed options**: only `-debug`, `-help`, `-list`, `-log`, `-verbose`, `-monitor`, `-quiet`, `-regard-warnings`, `-seed`
- **Arbitrary key-value pairs**: conjure accepts any `-key value` pair as script parameters (e.g., `-size 100x100 -color blue`)
- **Script files are the main argument**: each file is an MSL script
- **No image pipeline**: all processing is defined in the script

## `magick stream`

Streams raw pixel data from an image to a file or stdout.

```
magick stream [options ...] input-image raw-image
```

Key differences:
- **Two positional arguments**: input image, raw output
- **Minimal options**: only storage/format settings
- **No image transformation**: pure pixel extraction
- **`-storage` type**: pixel storage format (char, short, integer, float, double)

## Lexing Implications

For the completer, the tool name determines the available option set:

1. **Token 1 after `magick`**: check if it's a tool name or an option/image
2. **If tool name**: load the tool-specific option set
3. **If option/image**: default to `magick` (convert) option set
4. **Tool-specific options** only appear in their respective tools (e.g., `-metric` only in `compare`)
5. **Stack operators and parentheses** only in `magick` (convert)

## Edge Cases

- **`magick convert`**: deprecated alias that prints a warning — the completer should accept it but note it's legacy
- **`convert` as standalone**: even more deprecated, wraps `magick convert`
- **Tool name ambiguity**: tool names (`identify`, `mogrify`, etc.) could also be filenames — if a file named `identify` exists, it could be ambiguous
- **Options after tool name**: `magick identify -verbose image.png` — `-verbose` is an identify option, not a magick option
