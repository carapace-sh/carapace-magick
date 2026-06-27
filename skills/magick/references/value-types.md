# Value Types

The value types `magick` uses for option arguments: geometry, color, percent, and enumerated types.

> **Source of truth**: <https://imagemagick.org/script/command-line-processing.php> and <https://imagemagick.org/script/color.php>.

## Overview

`magick` options take typed arguments. The lexer must know which type each option expects to correctly consume the next token and offer appropriate completions. Many value types have their own sub-syntax.

## Geometry

The most important and complex value type. Geometry strings specify **dimensions and/or offsets** with a rich modifier syntax.

### Base Form

```
[WxH][+X+Y][%^!<>@]
```

| Component | Syntax | Meaning |
|-----------|--------|---------|
| Width | `W` or `Wx` | Width in pixels |
| Height | `xH` | Height in pixels |
| Both | `WxH` | Width and height |
| Offset | `+X+Y` or `-X-Y` | Position offset |
| Percent | `%` | Dimensions are percentage of original |
| Fit | `^` | Fill given area (may exceed dimensions) |
| Ignore aspect | `!` | Ignore aspect ratio |
| Shrink only | `>` | Only shrink larger images |
| Grow only | `<` | Only enlarge smaller images |
| Pixel count limit | `@count` | Maximum total pixel count |

### Geometry Examples

| Geometry | Meaning |
|----------|---------|
| `200x200` | 200×200 pixels |
| `200x200^` | Fill 200×200 area, keeping aspect ratio |
| `200x200!` | Exactly 200×200, ignore aspect ratio |
| `50%` | 50% of original dimensions |
| `200x200>` | Shrink to fit in 200×200, only if larger |
| `200x200<` | Enlarge to fill 200×200, only if smaller |
| `10000@` | Resize to fit within 10000 total pixels |
| `+10+20` | Offset: 10 pixels right, 20 pixels down |
| `200x200+10+20` | 200×200 with offset |

### Geometry Contexts

Different options interpret geometry differently:

| Option | Interpretation |
|--------|---------------|
| `-resize` | Resize image to fit |
| `-crop` | Extract region (offset is top-left corner) |
| `-extent` | Set canvas size (offset positions image on canvas) |
| `-border` | Add border of given width (single value = uniform) |
| `-density` | DPI: `72x72` or just `72` |
| `-page` | Canvas size and offset for multi-image formats |

## Color

`magick` supports extensive color specification syntax.

### Color Forms

| Form | Example | Description |
|------|---------|-------------|
| Named color | `red`, `DodgerBlue`, `transparent` | X11/SVG color name |
| Hex RGB | `#RGB`, `#RRGGBB`, `#RRRGGGBBB`, `#RRRRGGGGBBBB` | Hex notation |
| Hex RGBA | `#RGBA`, `#RRGGBBAA` | Hex with alpha |
| `rgb(r,g,b)` | `rgb(255,0,0)` | Functional notation |
| `rgba(r,g,b,a)` | `rgba(255,0,0,0.5)` | Functional with alpha |
| `hsl(h,s,l)` | `hsl(0,100%,50%)` | HSL notation |
| `cmyk(c,m,y,k)` | `cmyk(0,100,100,0)` | CMYK notation |
| Percent | `rgb(100%,0%,0%)` | Percentage values |
| Color from image | `xc:colorname` | Cross-reference (rare in option values) |

### Color Value Completion

Named colors come from `magick -list color` (hundreds of entries). For completion, the most useful set is the X11/SVG named colors.

## Percent Values

Many options accept percentage values:

| Pattern | Meaning |
|---------|---------|
| `50%` | 50 percent |
| `50` | Absolute value (pixels, level, etc.) — percent sign is required for percentage |

Whether a value is interpreted as percent or absolute depends on the option and whether `%` is appended. Some options (like `-resize`) accept both forms; others are always one or the other.

## Integer

Plain integer values for options like `-depth`, `-quality`, `-loop`:

```
-depth 8          # 8 bits per channel
-quality 85       # JPEG quality 85
-loop 0           # infinite animation loop
```

## Float

Floating-point values for options like `-gamma`, `-attenuate`, `-rotate`:

```
-gamma 1.5        # gamma correction
-rotate 45.5      # rotation in degrees
```

## Degrees

Used by rotation and angular operators:

```
-rotate 90        # degrees
-swirl 180        # degrees
-motion-blur 0x5+45   # geometry: radius×sigma+offset-angle
```

## Point

A point is a pair of coordinates, typically `X,Y` or `xY`:

```
-blue-primary 0.15,0.06   # chromaticity point
```

## Enumerated Types

Many options take values from a fixed enumeration, discoverable via `magick -list <type>`:

| Option | List Type | Example Values |
|--------|-----------|---------------|
| `-colorspace` | `colorspace` | `sRGB`, `RGB`, `CMYK`, `Gray`, `Lab`, `Oklab` |
| `-compose` | `compose` | `Over`, `Multiply`, `Screen`, `Blend`, `Copy` |
| `-compress` | `compress` | `None`, `JPEG`, `LZW`, `RLE`, `Zip` |
| `-channel` | `channel` | `Red`, `Green`, `Blue`, `Alpha`, `All`, `Default`, `Sync` |
| `-distort` | `distort` | `Affine`, `Perspective`, `ScaleRotateTranslate`, `Barrel`, `Polar` |
| `-filter` | `filter` | `Lanczos`, `Gaussian`, `Mitchell`, `Catrom`, `Box` |
| `-gravity` | `gravity` | `Center`, `North`, `NorthEast`, `East`, `SouthEast`, `South`, `SouthWest`, `West`, `NorthWest`, `None` |
| `-interlace` | `interlace` | `None`, `Line`, `Plane`, `Partition` |
| `-layers` | `layers` | `Optimize`, `Merge`, `Flatten`, `Composite`, `Compare` |
| `-morphology` | `morphology` | `Convolve`, `Dilate`, `Erode`, `Open`, `Close`, `TopHat`, `BottomHat` |
| `-type` | `type` | `Bilevel`, `Grayscale`, `Palette`, `TrueColor`, `ColorSeparation` |
| `-virtual-pixel` | `virtual-pixel` | `Background`, `Edge`, `Mirror`, `Tile`, `Transparent`, `Black`, `White` |
| `-orient` | `orientation` | `TopLeft`, `TopRight`, `BottomRight`, `BottomLeft`, `LeftTop`, `RightTop` |
| `-dispose` | `dispose` | `None`, `Background`, `Previous` |
| `-metric` | `metric` | `AE`, `MAE`, `MSE`, `PAE`, `PSNR`, `RMSE`, `SSIM`, `DSSIM` |
| `-evaluate` | `evaluate` | `Add`, `And`, `Abs`, `Divide`, `Log`, `Max`, `Min`, `Multiply`, `Or`, `Pow`, `Subtract`, `Xor` |
| `-preview` | `preview` | `Rotate`, `Shear`, `Roll`, `Noise`, `Segment`, `Edge`, `Solarize` |
| `-noise` | `noise` | `Gaussian`, `Impulse`, `Laplacian`, `Multiplicative`, `Poisson`, `Uniform` |

### Discovering Enumerations

The `-list` command prints all valid values for a type:

```bash
magick -list colorspace    # all valid colorspace names
magick -list compose       # all compose operators
magick -list distort       # all distort methods
```

The complete list of list categories is available via `magick -list list`.

## String

Free-form text values for options like `-label`, `-comment`, `-caption`, `-draw`:

```
-label "My Photo"          # simple string
-comment "Copyright 2024"  # annotation
-draw "circle 50,50 50,0"  # drawing primitive with its own syntax
```

The `-draw` option has a rich sub-grammar for graphic primitives (point, line, rectangle, circle, ellipse, polygon, polyline, path, image, text). This is a DSL within the string value.

## Expression

Mathematical expressions for `-fx` and `-evaluate`:

```
-fx "intensity*0.8"                    # simple expression
-evaluate-sequence add                 # per-sequence evaluation
-channel-fx "A=>R R=>A"                # channel exchange expression
```

## Edge Cases

- **Geometry without `x`**: `200` alone means width=200, height=proportional (for resize)
- **`0x0` geometry**: zero dimensions have special meaning depending on option
- **Color names with spaces**: must be quoted in shell: `-fill "Dodger Blue"`
- **Hex color shorthand**: `#F00` = `#FF0000` (3-digit shorthand)
- **Percent in geometry**: `50%x50%` or `50%` — the `%` modifies interpretation
- **Negative offsets**: `-crop 100x100-10-20` crops from bottom-right
