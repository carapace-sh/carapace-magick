# `-define` Syntax

How the `-define` option works: format-specific key=value pairs with a `format:key=value` structure.

> **Source of truth**: <https://imagemagick.org/script/command-line-options.php#define> and ImageMagick source `MagickWand/mogrify.c`.

## Overview

`-define` provides a way to pass format-specific or internal settings that don't have dedicated command-line options. It's the escape hatch for fine-grained control over coders, delegates, and internal algorithms.

```
-define format:key=value
-define key=value
```

## Syntax

### Basic Form

The argument to `-define` is a single string with up to three parts:

```
[format:]key=value
```

|| Component | Required | Description |
||-----------|----------|-------------|
|| `format` | Optional | Coder/format name (e.g., `png`, `jpeg`, `tiff`, `webp`) |
|| `key` | Required | The setting name |
|| `value` | Required | The setting value |

When `format` is omitted, the define applies globally (not tied to a specific coder).

### The `+define` Form

`+define format:key` removes a previously set define:

```bash
-define jpeg:quality=85    # set JPEG quality
+define jpeg:quality       # remove the define, revert to default
```

### Multiple Defines

Each `-define` sets one key. For multiple settings, use multiple `-define` options:

```bash
magick input.png -define png:compression-level=9 -define png:compression-filter=0 -define png:compression-strategy=2 output.png
```

## Format-Specific Defines

Each image coder (format) supports its own set of define keys. These are not discoverable from the command line — they are documented in the ImageMagick documentation for each format.

### Common Format Defines

#### JPEG (`jpeg:`)

|| Key | Value | Description |
||-----|-------|-------------|
|| `jpeg:quality` | integer (1-100) | Compression quality |
|| `jpeg:optimize-codes` | boolean | Optimize Huffman coding tables |
|| `jpeg:progressive` | boolean | Write progressive JPEG |
|| `jpeg:sampling-factor` | geometry | Chroma subsampling (e.g., `2x2`, `4:2:0`) |
|| `jpeg:extent` | size limit | Maximum file size in KB |

#### PNG (`png:`)

|| Key | Value | Description |
||-----|-------|-------------|
|| `png:compression-level` | integer (0-9) | zlib compression level |
|| `png:compression-filter` | integer (0-5) | Row filter method |
|| `png:compression-strategy` | integer (0-3) | zlib strategy |
|| `png:color-type` | integer | Force color type (0=Gray, 2=RGB, 3=Indexed, 4=GrayAlpha, 6=RGBA) |
|| `png:bit-depth` | integer | Force bit depth |
|| `png:exclude-chunk` | chunk names | Exclude chunks from output |
|| `png:include-chunk` | chunk names | Include specific chunks |

#### TIFF (`tiff:`)

|| Key | Value | Description |
||-----|-------|-------------|
|| `tiff:compress` | string | Compression (none, lzw, zip, jpeg, etc.) |
|| `tiff:rows-per-strip` | integer | Rows per strip |
|| `tiff:tile-geometry` | geometry | Tiled TIFF layout |
|| `tiff:predictor` | integer | Predictor for compression |

#### WebP (`webp:`)

|| Key | Value | Description |
||-----|-------|-------------|
|| `webp:lossless` | boolean | Use lossless encoding |
|| `webp:quality` | float (0-100) | Lossy quality |
|| `webp:method` | integer (0-6) | Compression method (0=fast, 6=slowest) |
|| `webp:alpha-quality` | float | Alpha channel quality |

#### GIF (`gif:`)

|| Key | Value | Description |
||-----|-------|-------------|
|| `gif:interlace` | boolean | Interlaced GIF |
|| `gif:optimize` | boolean | Optimize animation frames |
|| `gif:disposal` | method | Frame disposal method |

#### HEIC/AVIF (`heic:`)

|| Key | Value | Description |
||-----|-------|-------------|
|| `heic:quality` | integer | Compression quality |
|| `heic:lossless` | boolean | Lossless encoding |
|| `heic:speed` | integer (0-8) | Encoding speed |

## Global Defines (No Format Prefix)

Some defines apply globally:

|| Key | Value | Description |
||-----|-------|-------------|
|| `optimize` | boolean | General optimization |
|| `type` | type | Hint for output image type |
|| `preserve-colorspace` | boolean | Preserve input colorspace |

## Coder Discovery

The available format prefixes correspond to the coders listed by `magick -list coder` and `magick -list format`. The format column in `-list format` shows the coder names that can be used as prefixes.

## Lexing Implications

For completion, the `-define` argument has a structured grammar:

1. **Format prefix**: complete from `magick -list format` coder names
2. **Colon separator**: `format:`
3. **Key**: complete from format-specific key list (must be maintained manually or extracted from source)
4. **Equals sign**: `=`
5. **Value**: complete from format+key-specific value list

The lexer should split the `-define` argument at `:` and `=` boundaries:

```
-define jpeg:quality=85
         ↑     ↑      ↑
       format  key   value
```

### Partial Completion States

|| Cursor Position | Complete |
||-----------------|---------|
|| `-define ` | Format prefixes + global keys |
|| `-define jpeg:` | JPEG-specific keys |
|| `-define jpeg:quality=` | Quality values (1-100) |
|| `-define png:color-type=` | Color type values (0, 2, 3, 4, 6) |

## Edge Cases

- **Multiple colons**: some define keys contain colons (rare) — `delegate:delegate=value`
- **Values with equals**: `png:exclude-chunk=tEXt,zTXt` — the value contains commas, not additional `=`
- **Boolean values**: some defines accept `true`, `false`, `1`, `0`, `yes`, `no`
- **`+define` without value**: `+define jpeg:quality` — removes the define (no `=value` part)
- **Unknown defines**: `magick` silently ignores unknown define keys — no validation error
