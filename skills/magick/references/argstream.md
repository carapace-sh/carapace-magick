# Argstream — Command-Line Structure

How `magick` processes its command-line arguments: the image pipeline model, token classification, and how options, images, and operators flow through the command line.

> **Source of truth**: <https://imagemagick.org/script/command-line-processing.php> and ImageMagick source `MagickCore/magick-cli.c`.

## Overview

`magick` is an **image pipeline processor**. Unlike traditional CLIs where flags form a tree of subcommands, `magick` arguments form a **linear sequence** where tokens are classified into settings, operators, images, and stack manipulations — all processed left-to-right, building up an image sequence that is written to the final output.

The fundamental invocation syntax is:

```
magick [tool] [ {option} | {image} ... ] {output_image}
magick [ {option} | {image} ... ] -script {filename} [ {script_args} ...]
```

## Token Classification

Every argument token falls into one of these categories:

| Category | Prefix | Effect | Example |
|----------|--------|--------|---------|
| **Image input** | none (filename) or `-` for stdin | Adds image(s) to the sequence | `input.png`, `xc:red` |
| **Image output** | none (last non-option token) | Writes final image sequence | `output.png` |
| **Setting** | `-word value` / `+word` | Configures future operators (does not modify images) | `-background blue`, `-quality 85` |
| **Operator** | `-word [args]` | Transforms the current image(s) immediately | `-resize 200x200`, `-blur 0x5` |
| **Stack operator** | `-word [args]` | Manipulates the image sequence (clone, delete, swap) | `-clone 0`, `-delete 1`, `-swap 0,1` |
| **Parenthesis** | `\(` ... `\)` | Creates a scoped sub-pipeline | `\(` `-clone 0` `-negate` `\)` |

The lexer cannot determine whether a `-word` token is a setting or operator from syntax alone — it must consult a static option definition table.

## Pipeline Flow

```
┌──────────────────────────────────────────────────────────────┐
│  CLI:  [setting] [image] [operator] [stack_op] [...] output │
│                                                              │
│  Image Sequence:  img0 → img1 → img2 → ... (in memory)      │
│      settings configure ──► operators transform              │
│      stack ops rearrange ──► parentheses scope               │
└──────────────────────────────────────────────────────────────┘
```

Arguments are processed strictly left-to-right:

1. **Settings** accumulate — they configure the environment for subsequent operators
2. **Image inputs** are read and appended to the image sequence
3. **Operators** immediately transform the current image(s) in the sequence
4. **Stack operators** rearrange images in the sequence (clone, delete, swap, insert)
5. **Parentheses** push/pop the image sequence (see below)
6. **The last non-option token** is the output filename

## Image Input Forms

An image input can be:

| Form | Example | Description |
|------|---------|-------------|
| Filename | `photo.jpg` | Read from file (format auto-detected) |
| Format prefix | `png:photo.dat` | Explicit format override |
| STDIN | `-` | Read from standard input |
| Built-in generator | `xc:red`, `gradient:`, `plasma:`, `pattern:checkerboard` | Generate an image without a file |
| `-read filename` | `-read input.png` | Explicit read (same as bare filename but unambiguous) |

The format prefix syntax is `format:filename` — the colon separates the coder name from the filename. This disambiguates when the file extension doesn't match the format or when reading from stdin: `png:-`.

## Image Output

The output is the **last non-option argument**. Like inputs, it supports format prefix:

```bash
magick input.jpg png:output.dat   # force PNG format regardless of extension
magick input.jpg -                # write to stdout
```

Multiple outputs are possible with `-write filename` within the pipeline.

## Parentheses — Scoped Sub-Pipelines

Parentheses create a nested image sequence scope. The shell requires escaping: `\(` and `\)` (or quoting).

```
magick \( -size 100x100 xc:red \) \( -size 100x100 xc:blue \) -append out.png
```

### How Parentheses Work

1. **`(` — push**: saves the current image sequence and settings onto a stack, starts a fresh sequence
2. **`)` — pop**: merges the sub-pipeline's image sequence back into the outer sequence; restores settings

### Scoped Settings

Settings inside parentheses only affect that scope. With `-respect-parentheses`, settings revert when the scope closes:

```bash
magick -background red \( -background blue -size 50x50 xc: ) -append out.png
#  -background blue only applies inside ()
#  -background red is restored after )
```

Without `-respect-parentheses`, settings **leak out** of the parentheses — the inner setting persists after the closing `)`.

## Image Sequence vs. Single Image

Many operators work on the "current image" (the last image in the sequence). Sequence operators work on all images. The distinction matters for:

- **Single-image operators**: `-resize`, `-blur`, `-rotate` — applied per-image
- **Sequence operators**: `-append`, `-flatten`, `-morph`, `-layers` — combine or rearrange the sequence
- **Stack operators**: `-clone`, `-delete`, `-swap`, `-insert` — rearrange the sequence without transforming pixels

## Script Mode

With `-script filename`, `magick` reads commands from a file. This is used for complex pipelines where shell escaping becomes unwieldy. In a script:

- No shell escaping needed for parentheses
- Each line is a separate command or continuation
- Arguments after the script filename are available as `%1`, `%2`, etc.

## Edge Cases

- **Filenames starting with `-`**: ambiguous with options. Use `./-photo.jpg` or explicit `-read`.
- **`-` as output**: writes to stdout; the format is determined by the output filename prefix or `-format` setting.
- **No output**: always an error — magick requires an output filename. Use `null:` as a dummy output when only side effects matter (e.g., `magick input.png -identify null:`).
- **Inline images**: `xc:red` is an image input, not an option — the lexer must recognize built-in coders as image sources.
- **Plus-form reset**: `+setting` (e.g., `+background`) resets the setting to its default. The lexer must handle `+` prefix as a distinct token form.
