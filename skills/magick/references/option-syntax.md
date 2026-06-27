# Option Syntax

How `magick` option names are structured: the `-`/`+` prefix convention, boolean flags, value-taking options, and quoting rules.

> **Source of truth**: <https://imagemagick.org/script/command-line-options.php> and ImageMagick source `MagickWand/mogrify.c`.

## Option Name Patterns

`magick` uses a **single-dash prefix** for all options — there are no double-dash `--` long options. The `+` prefix is the **reset/inverse** form.

### Core Patterns

| Pattern | Example | Description |
|---------|---------|-------------|
| `-word` | `-verbose`, `-strip` | Boolean flag (enable) |
| `-word value` | `-resize 200x200`, `-quality 85` | Option with value |
| `+word` | `+verbose`, `+append` | Reset/inverse form |
| `+word` | `+background` | Reset to default (no value) |
| `+word value` | `+clone 0` | Inverse operation (rare, most +forms take no value) |

### The `-` vs `+` Convention

This is a **core pattern** unique to ImageMagick. The `+` prefix is not just "disable" — its meaning depends on the option type:

| Option Type | `-form` | `+form` | Example |
|------------|---------|---------|---------|
| **Setting** (boolean) | Enable | Disable | `-verbose` / `+verbose` |
| **Setting** (value) | Set value | Reset to default | `-background blue` / `+background` |
| **Operator** (directional) | One direction | Other direction | `-append` (top-to-bottom) / `+append` (left-to-right) |
| **Operator** (on/off) | Apply | Remove/undo | `-clip` / `+clip` |
| **Stack operator** | Normal | Inverse | `-clone 0` / `+clone 0` (rarely different) |

### Lexing Implication

The lexer must:
1. Detect the `+` prefix and classify it as the reset/inverse form
2. Know that `+word` for a value-taking setting takes **no argument** (reset has no value)
3. Know that `+word` for a directional operator may still take arguments

## Value-less vs Value-taking Options

Options are classified by whether they consume a following token:

| Class | Prefix | Takes Value? | Example |
|-------|--------|-------------|---------|
| Boolean setting | `-` | No | `-verbose`, `-antialias`, `-monitor` |
| Boolean setting | `+` | No | `+verbose`, `+antialias` |
| Value setting | `-` | Yes | `-background color`, `-quality value` |
| Value setting (reset) | `+` | No | `+background`, `+quality` |
| Value operator | `-` | Yes | `-resize geometry`, `-blur geometry` |
| No-value operator | `-` | No | `-strip`, `-flip`, `-negate` |
| Directional operator | `-`/`+` | No or Yes | `-append` / `+append` (no value) |
| Stack operator | `-`/`+` | Yes | `-clone indexes`, `-delete indexes` |

The classification must be known from a **static option definition table** — it cannot be inferred from syntax alone.

## Option Name Structure

All option names are **multi-letter words** prefixed by `-` or `+`. There are no single-letter options (unlike ffmpeg's `-y`, `-n`, etc.). Names use lowercase with hyphens:

```
-adaptive-blur     -black-threshold    -gaussian-blur
-channel-fx        -sigmoidal-contrast  -write-mask
```

### No Short Forms

Unlike many CLIs, `magick` has **no single-letter aliases**. Every option is a full word. This simplifies lexing — there is no ambiguity between `-c` (short flag) and `-clip` (word flag).

## Quoting and Escaping

`magick` uses its own quoting rules for option values, **separate from shell quoting**. The shell strips its quoting first, then `magick` applies its own.

### magick Quoting Rules

1. **Backslash escaping** `\x` — escapes any special character
2. **Percent escaping** `%x` — format escape sequences in strings like `-format` and `-label`
3. **No single-quote delimiters** — unlike ffmpeg, `magick` does not have its own single-quote quoting layer

### Interaction with Shell Quoting

There are **two layers**: the shell layer, then the `magick` layer:

```bash
# Shell strips outer quotes, magick sees the inner string
magick input.jpg -label "Hello World" output.jpg
magick input.jpg -annotate 0x0 "My text" output.jpg
```

### Special Characters in Values

| Context | Special Characters | Escape |
|---------|-------------------|--------|
| General option values | `\`, `%`, `@` | `\` prefix |
| Geometry values | `x`, `%`, `^`, `!`, `>`, `<`, `@` | Contextual (not escaped, parsed as geometry) |
| Color values | `#`, `rgb()`, `hsl()` | Not escaped (parsed as color) |
| `-format` strings | `%[...]`, `%w`, `%h`, `%m` | `%%` for literal `%` (e.g. `-format "%%w"`) |
| `-draw` primitives | `"`, `'`, `,`, `(`, `)` | `\` or quoting |
| `-define` keys | `:` | `\:` (rarely needed) |

### Parentheses in Shell

Parentheses **must be escaped** in the shell, since `(` and `)` are shell operators:

```bash
# Bash/zsh: escape with backslash
magick \( -size 50x50 xc:red \) -append out.png

# Or quote them
magick '(' -size 50x50 xc:red ')' -append out.png
```

The lexer receives the already-escaped `\(` or `'('` as a bare `(` token after the shell processes it.

## Edge Cases

- **`-list type`** consumes a list category name (e.g., `-list colorspace`) — not an image option
- **`-set property value`** has two arguments — the property name and the value
- **`-define format:key=value`** has a compound value — see [define-syntax.md](define-syntax.md)
- **`-channel-fx expression`** has a complex expression value with its own grammar
- **`+clone` / `+delete`** — some stack operators have `+` forms with subtly different behavior
- **`-read filename` vs bare filename** — both add an image to the sequence, but `-read` is unambiguous
