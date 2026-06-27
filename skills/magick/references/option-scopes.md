# Option Scopes — Settings, Operators, and Stack Ops

How `magick` classifies options into categories that determine when and how they take effect. This classification is essential for the lexer — it determines whether an option consumes a value argument and how it interacts with the image pipeline.

> **Source of truth**: <https://imagemagick.org/script/command-line-processing.php> and ImageMagick source `MagickWand/mogrify.c`.

## Overview

`magick` organizes its options into **five categories**, each with distinct semantics:

```
Settings  → configure future behavior (no immediate effect)
Operators → transform current image(s) immediately
Channel Operators → operate on specific channels
Sequence Operators → combine/rearrange the image sequence
Stack Operators   → manipulate the image list (clone, delete, swap)
```

The `-help` output groups options under these headings. The lexer must classify each option into one of these categories to determine its argument consumption and pipeline behavior.

## Image Settings

Settings configure the environment for subsequent operations. They **do not modify images** — they change how future operators behave.

### Characteristics

- **Accumulative**: settings persist until changed or reset
- **Scoped by parentheses**: with `-respect-parentheses`, settings revert when a `)` closes
- **`+form` resets to default**: `+background` resets to the built-in default

### Dual-Nature Options

Some options appear in multiple `-help` sections because they act as both settings and operators depending on their argument value:

| Option | Setting Behavior | Operator Behavior |
|--------|-----------------|-------------------|
| `-alpha` | Configures alpha channel mode | `on`, `off`, `activate`, `deactivate` immediately change the alpha channel |
| `-channel` | Sets the active channel mask for subsequent operators | `set`, `on`, `off` immediately modify the channel mask |

### Common Settings

| Setting | Value Type | Description |
|---------|-----------|-------------|
| `-background` | color | Background color |
| `-colorspace` | type | Alternate colorspace |
| `-compress` | type | Pixel compression type |
| `-density` | geometry | Horizontal/vertical resolution |
| `-depth` | value | Image depth |
| `-fill` | color | Fill color for drawing |
| `-font` | name | Font for text rendering |
| `-format` | string | Output format string (percent-escape sequence, e.g. `%w %h`) |
| `-gravity` | type | Text placement direction |
| `-interlace` | type | Interlacing scheme |
| `-quality` | value | Compression level |
| `-stroke` | color | Stroke color for drawing |
| `-type` | type | Image type |
| `-verbose` | (boolean) | Print detailed information |
| `-define` | format:key=value | Format-specific option (see [define-syntax.md](define-syntax.md)) |

## Image Operators

Operators **immediately transform** the current image(s) in the sequence. They consume the image, modify it, and replace it.

### Characteristics

- **Immediate effect**: applied as soon as encountered in the argstream
- **Per-image**: most operators apply to each image in the sequence individually
- **Some take geometry arguments**: `-resize WxH`, `-blur radiusxsigma`
- **`+form` is directional or inverse**: `-append` (vertical) / `+append` (horizontal)

### Operator Sub-Categories by Value Type

| Value Pattern | Examples |
|---------------|---------|
| No value (boolean) | `-strip`, `-flip`, `-flop`, `-negate`, `-normalize`, `-trim`, `-despeckle` |
| Geometry | `-resize`, `-scale`, `-sample`, `-thumbnail`, `-extent`, `-crop`, `-shave`, `-border`, `-chop` |
| Radius/sigma | `-blur`, `-gaussian-blur`, `-sharpen`, `-unsharp`, `-adaptive-blur`, `-adaptive-sharpen`, `-motion-blur` |
| Degrees | `-rotate`, `-swirl`, `-motion-blur`, `-polaroid` |
| Color | `-opaque`, `-transparent`, `-floodfill`, `-tint` |
| Method + args | `-distort`, `-morphology`, `-function`, `-sparse-color` |
| Expression | `-evaluate`, `-fx`, `-channel-fx` |
| Value | `-level`, `-gamma`, `-threshold`, `-posterize`, `-solarize`, `-modulate` |

## Image Channel Operators

Channel operators work on specific channels of the image. The `-channel` option (listed under "Image Operators" in `-help`) sets the active channel mask, controlling which channels subsequent operators affect.

| Operator | Value | Description |
|----------|-------|-------------|
| `-channel` | mask | Set the active channel mask (dual-nature: acts as both setting and operator) |
| `-channel-fx` | expression | Exchange, extract, or transfer channels |
| `-separate` | (none) | Separate channels into grayscale images |

The `-channel` setting controls which channels these operators affect:

```bash
magick input.png -channel RGB -separate channel_%d.png
```

## Image Sequence Operators

Sequence operators work on the **entire image sequence** (or multiple images within it), not just the current image.

| Operator | Value | Description |
|----------|-------|-------------|
| `-append` / `+append` | (none) | Append vertically/horizontally |
| `-flatten` | (none) | Flatten sequence into single image |
| `-mosaic` | (none) | Create mosaic from sequence |
| `-coalesce` | (none) | Merge sequence |
| `-combine` | (none) | Combine into color channels |
| `-compare` | (none) | Compare images |
| `-composite` | (none) | Composite images |
| `-layers` | method | Layer method (optimize, merge, etc.) |
| `-morph` | value | Morph between images |
| `-fx` | expression | Apply math expression |
| `-write` | filename | Write current sequence to file |

## Image Stack Operators

Stack operators rearrange the image sequence without transforming pixel data. They are critical for complex pipelines.

| Operator | Value | Description |
|----------|-------|-------------|
| `-clone` | indexes | Clone image(s) by index |
| `-delete` | indexes | Delete image(s) by index |
| `-duplicate` | count,indexes | Duplicate image(s) |
| `-insert` | index | Insert last image at position |
| `-reverse` | (none) | Reverse image sequence |
| `-swap` | indexes | Swap two images |

### Index Syntax

Stack operators use **index expressions** to refer to images in the sequence:

| Syntax | Meaning |
|--------|---------|
| `0` | First image |
| `1` | Second image |
| `0-3` | Images 0 through 3 (range) |
| `0,2,4` | Images 0, 2, and 4 (list) |
| `-1` | Last image |
| `-2` | Second-to-last image |
| `0--1` | All images (first through last) |

### Clone Pattern

A common pattern is to clone an image, transform the clone, then merge:

```bash
magick input.png \
  \( -clone 0 -resize 50% \) \
  \( -clone 0 -resize 25% \) \
  -delete 0 -append out.png
```

The `-delete 0` removes the original; the two resized clones remain.

## Interactions Between Categories

The categories interact in important ways:

1. **Settings before operators**: a setting must appear *before* the operator it configures
   ```bash
   magick input.png -gravity center -annotate 0x0 "text" out.png
   #  -gravity must come before -annotate
   ```

2. **Channel setting + channel operators**: `-channel` controls which channels are affected by subsequent channel operators
   ```bash
   magick input.png -channel R -separate red.png
   ```

3. **Stack operators in parentheses**: cloning inside parentheses creates images in the sub-scope
   ```bash
   magick input.png \( -clone 0 -negate \) -append out.png
   #  The cloned+negated image merges back and is appended
   ```

4. **`-write` inside pipeline**: writes current image sequence without ending the pipeline
   ```bash
   magick input.png -resize 50% -write small.png -resize 50% -write tiny.png null:
   ```

## Edge Cases

- **`-region geometry`**: restricts subsequent operators to a region of the image — a scoped operator
- **`-set property value`**: sets a metadata property — takes two arguments (name, value)
- **`-define format:key=value`**: format-specific setting — compound argument, see [define-syntax.md](define-syntax.md)
- **`-respect-parentheses`**: changes how settings behave across parentheses boundaries
- **`-list type`**: prints a list of valid values and exits — not a pipeline option
- **Options appearing after output filename**: silently ignored by some sub-tools, error in others
