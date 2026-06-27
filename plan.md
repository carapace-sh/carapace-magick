# Implementation Plan — carapace-magick

Go library for parsing ImageMagick `magick` CLI argument streams into ASTs with shell completion support. Part of the [carapace-sh](https://github.com/carapace-sh) ecosystem. Module path: `github.com/carapace-sh/carapace-magick`.

## Architecture Overview

Follow the same architecture as `carapace-ffmpeg`: parser packages with AST/completion support, a completer package wiring them to carapace actions, a debug CLI, and completer CLIs per sub-tool.

```
cmd/carapace-magick/            Completer CLI for magick (default/convert)
cmd/carapace-magick-identify/   Completer CLI for magick identify
cmd/carapace-magick-mogrify/    Completer CLI for magick mogrify
cmd/carapace-magick-compare/    Completer CLI for magick compare
cmd/carapace-magick-composite/  Completer CLI for magick composite
cmd/carapace-magick-montage/    Completer CLI for magick montage
cmd/carapace-magick-debug/      Debug/diagnostic CLI (JSON output)
pkg/argstream/                  Argument stream parser (options, images, stack ops, parentheses)
pkg/definevalue/                -define format:key=value parser
pkg/completer/                  Shared completion dispatch logic
pkg/actions/tools/magick/      Carapace action functions for magick value types
man/magick/                     YAML descriptions for completion value types
skills/magick/                  AI agent reference documentation (exists already)
testdata/                       Test images for integration tests
```

### Key Differences from carapace-ffmpeg

| Aspect | ffmpeg | magick |
|--------|--------|--------|
| Pipeline model | Input → Filter → Output (positional) | Image sequence with stack operations |
| Option scope | Global / Per-File / Per-Stream | Setting / Operator / Stack / Channel / Sequence |
| Sub-argument DSLs | Stream specifier, filter graph, map value | `-define` format:key=value, `-draw` primitives, `-channel-fx` expression |
| Option prefix | `-` only | `-` and `+` (reset/inverse) |
| Parentheses | Not used | Group sub-pipelines with scope |
| Tool profiles | ffmpeg, ffplay, ffprobe | magick (convert), identify, mogrify, compare, composite, montage, conjure, stream |
| Probe tool | ffprobe | `magick identify` |

---

## Phase 1: Project Scaffolding

### 1.1 Initialize Go module and project structure

- [ ] `go mod init github.com/carapace-sh/carapace-magick`
- [ ] Add dependencies: `carapace`, `carapace-spec`, `cobra`
- [ ] Create `.gitignore`, `LICENSE`, `README.md`, `CONTRIBUTING.md`, `AGENTS.md`
- [ ] Create `.github/workflows/go.yml` (mirror carapace-ffmpeg CI: build, test, gofmt, staticcheck)
- [ ] Create `.goreleaser.yml` (binaries for all completer CLIs + debug CLI)
- [ ] Create directory structure: `cmd/`, `pkg/`, `man/`, `testdata/`

### 1.2 Test data

- [ ] `testdata/generate.go` — generate test images using `magick` itself (small PNG, GIF with frames, multi-page TIFF, etc.)
- [ ] `testdata/generate.sh` — shell script invoked by `go generate`

---

## Phase 2: Argstream Parser (`pkg/argstream/`)

This is the core parser — the equivalent of ffmpeg's argstream but adapted to magick's image pipeline model.

### 2.1 Option definitions (`options.go`)

Define `OptionDef`, `OptionCategory`, `OptionType`, `ValueType` types:

```go
type OptionCategory int  // CategorySetting, CategoryOperator, CategoryChannelOp, CategorySequenceOp, CategoryStackOp, CategoryMisc
type OptionType int       // TypeBoolean, TypeValue
type ValueType string     // "geometry", "color", "colorspace", "compose", etc.
type OptionDef struct {
    Name         string
    Description  string
    Category     OptionCategory
    Type         OptionType
    ValueType    ValueType     // only for TypeValue
    HasPlusForm  bool          // whether +name is valid
    PlusBehavior PlusBehavior  // PlusReset, PlusInverse, PlusDirectional
}
```

Populate the option index from `magick -help` output. Start with the most common options (~100), expand later. Group by the categories shown in `-help`:
- Image Settings
- Image Operators
- Image Channel Operators
- Image Sequence Operators
- Image Stack Operators
- Miscellaneous Options

### 2.2 Tool-specific option files

Like ffmpeg's `ffplay_options.go` / `ffprobe_options.go`, create per-tool option sets:

- [ ] `identify_options.go` — identify-only options
- [ ] `mogrify_options.go` — mogrify-only options (+ `-format`, `-path`)
- [ ] `compare_options.go` — compare-only options (`-metric`, `-dissimilarity-threshold`, etc.)
- [ ] `composite_options.go` — composite-only options (`-blend`, `-dissolve`, `-stegano`, etc.)
- [ ] `montage_options.go` — montage-only options (`-tile`, `-frame`, `-shadow`, etc.)
- [ ] `conjure_options.go` — minimal options
- [ ] `stream_options.go` — minimal options

### 2.3 Profile system (`profile.go`)

```go
type ToolProfile struct {
    Name             string
    HasOutputArg     bool       // identify, mogrify: false; others: true
    HasOperators     bool       // identify, conjure: false; others: true
    HasStackOps      bool       // only magick (convert)
    HasParentheses   bool       // only magick (convert)
    OptionIndex      map[string]*OptionDef
}
```

Define `DefaultMagickProfile`, `DefaultIdentifyProfile`, etc.

### 2.4 AST types (`ast.go`)

```go
type TokenKind int  // KindOption, KindImage, KindOutput, KindLParen, KindRParen, KindToolName
type Token struct { Kind TokenKind; Value string; Span Span }
type ParenGroup struct { Open, Close Span; Tokens []*Token }
type Program struct { Tool string; Tokens []*Token; Groups []*ParenGroup }
```

### 2.5 Parser (`parser.go`)

State machine for the magick argstream:

```
States: START → TOOL_SELECTION → PIPELINE → DONE

In PIPELINE state, tokens are classified:
  '('  → push scope (save settings, start sub-sequence)
  ')'  → pop scope (merge sub-sequence, restore settings)
  '-word' / '+word' → option (lookup in profile, consume value if TypeValue)
  bare word → image input (or output if last non-option token)
```

Key differences from ffmpeg's parser:
- No `-i` marker — image inputs are bare filenames (non-option tokens)
- Parentheses create scope push/pop
- `+word` is a distinct option form
- The last non-option token is the output (when `HasOutputArg` is true)
- `-read` is an explicit image read (unambiguous)
- Tool name is the first positional token (if it matches a known tool)

### 2.6 Completion parser (`completion_parser.go`, `completion.go`)

`ParseForCompletion(args, trailingSpace, profile)` → `CompletionContext`:

```go
type ExpectedToken int
const (
    ExpectedToolName      // first position: which sub-tool?
    ExpectedOptionName    // after image/setting/operator: new option
    ExpectedOptionValue   // after a value-taking option
    ExpectedPlusOptionName // after image/setting: new +option
    ExpectedImage         // image input filename
    ExpectedOutput        // output filename
    ExpectedDefineValue   // special: -define format:key=value
)

type CompletionContext struct {
    ExpectedTokens []ExpectedToken
    Tool           string
    CurrentOption  *OptionDef
    OptionForm     OptionForm  // FormDash, FormPlus
    InParentheses  bool
    // ...
}
```

### 2.7 Tests (`argstream_test.go`, `completion_test.go`)

- [ ] Table-driven tests for `Parse()` — cover tool names, options, images, parentheses, stack ops
- [ ] Table-driven tests for `ParseForCompletion()` — cover each `ExpectedToken` state
- [ ] Edge cases: filenames starting with `-`, `+` reset forms, nested parentheses

---

## Phase 3: Define Value Parser (`pkg/definevalue/`)

Parser for `-define format:key=value` strings — a simpler DSL than ffmpeg's filter graph, but still needs structured parsing for completion.

### 3.1 Structure

```
definevalue/
  parser.go          Parse("jpeg:quality=85") → *DefineValue
  completion_parser.go  ParseForCompletion("jpeg:") → *CompletionContext
  completion.go      ExpectedToken, CompletionContext types
  span.go            Span type
  definevalue_test.go
  completion_test.go
```

### 3.2 AST

```go
type DefineValue struct {
    Format string  // "jpeg", "png", etc. (empty for global)
    Key    string  // "quality", "compression-level", etc.
    Value  string  // "85", "9", etc.
}
```

### 3.3 Completion context

Three completion stages:
1. **Format prefix** → complete from `magick -list format`
2. **Key** → complete from format-specific key list
3. **Value** → complete from key-specific value list

---

## Phase 4: Action Functions (`pkg/actions/tools/magick/`)

Carapace completion actions for magick value types.

### 4.1 Dynamic actions (shell out to `magick -list`)

These query the live `magick` binary for current values:

| Action | Source Command |
|--------|---------------|
| `ActionColorspaces` | `magick -list colorspace` |
| `ActionComposes` | `magick -list compose` |
| `ActionCompressTypes` | `magick -list compress` |
| `ActionChannels` | `magick -list channel` |
| `ActionDistortMethods` | `magick -list distort` |
| `ActionFilters` | `magick -list filter` |
| `ActionGravities` | `magick -list gravity` |
| `ActionInterlaceTypes` | `magick -list interlace` |
| `ActionLayerMethods` | `magick -list layers` |
| `ActionMorphologyMethods` | `magick -list morphology` |
| `ActionTypes` | `magick -list type` |
| `ActionVirtualPixelMethods` | `magick -list virtual-pixel` |
| `ActionOrientations` | `magick -list orientation` |
| `ActionDisposes` | `magick -list dispose` |
| `ActionMetrics` | `magick -list metric` |
| `ActionEvaluateOps` | `magick -list evaluate` |
| `ActionFormats` | `magick -list format` |
| `ActionFonts` | `magick -list font` |
| `ActionColors` | `magick -list color` |
| `ActionKernels` | `magick -list kernel` |

### 4.2 Static actions (hardcoded)

| Action | Values |
|--------|--------|
| `ActionBoolean` | `true`, `false`, `1`, `0` |
| `ActionAlphaOption` | `on`, `off`, `activate`, `deactivate`, `set`, `opaque`, `copy`, `transparent`, `extract`, `background`, `shape` |
| `ActionAutoThreshold` | `Kapur`, `OTSU`, `Triangle` |
| `ActionNoiseTypes` | `Gaussian`, `Impulse`, `Laplacian`, `Multiplicative`, `Poisson`, `Uniform` |
| `ActionPreviewTypes` | `Rotate`, `Shear`, `Roll`, `Noise`, `Segment`, etc. |
| `ActionGrayscaleMethods` | (from docs) |
| `ActionIntensityMethods` | (from `magick -list intensity`) |
| `ActionStorageTypes` | `char`, `short`, `integer`, `float`, `double` |

### 4.3 Special actions

- [ ] `ActionGeometry(prefix string)` — contextual geometry completion with special characters (`^`, `!`, `>`, `<`, `@`, `%`)
- [ ] `ActionColor()` — named colors from `magick -list color` + hex forms
- [ ] `ActionDefineFormat()` — format prefixes from `magick -list format`
- [ ] `ActionDefineKeys(format string)` — format-specific define keys (maintained manually)
- [ ] `ActionDefineValues(format, key string)` — key-specific values

### 4.4 Probe action (`pkg/probe/`)

Wrap `magick identify` for image-aware completion:

```go
type ImageInfo struct {
    Width     int
    Height    int
    Format    string
    Colors    int
    Depth     int
    Colorspace string
}
func Probe(inputPath string) *ImageInfo
```

Used for dimension-aware geometry completion (e.g., suggest `50%` or `1920x1080` based on input image size).

---

## Phase 5: Completer Package (`pkg/completer/`)

Shared completion dispatch logic used by all completer CLIs.

### 5.1 Core dispatch

```
carapace.Context
  → ContextToArgs() → (args, trailingSpace)
  → argstream.ParseForCompletionWithProfile(args, trailingSpace, profile)
  → argstream.CompletionContext
  → switch on ExpectedToken:
      ExpectedToolName      → ActionToolNames
      ExpectedOptionName    → ActionOptions (with + variants)
      ExpectedOptionValue   → ActionOptionValue (switch on ValueType)
      ExpectedImage         → ActionFiles (image files)
      ExpectedOutput        → ActionFiles
      ExpectedDefineValue   → ActionDefineValue (delegate to definevalue)
```

### 5.2 Option value dispatch

`ActionOptionValue` switches on `ValueType` to pick the right action:

| ValueType | Action |
|-----------|--------|
| `geometry` | `ActionGeometry` |
| `color` | `ActionColor` |
| `colorspace` | `ActionColorspaces` |
| `compose` | `ActionComposes` |
| `compress` | `ActionCompressTypes` |
| `channel` | `ActionChannels` |
| `distort` | `ActionDistortMethods` |
| `filter` | `ActionFilters` |
| `gravity` | `ActionGravities` |
| `interlace` | `ActionInterlaceTypes` |
| `layers` | `ActionLayerMethods` |
| `morphology` | `ActionMorphologyMethods` |
| `type` | `ActionTypes` |
| `virtual_pixel` | `ActionVirtualPixelMethods` |
| `metric` | `ActionMetrics` |
| `evaluate` | `ActionEvaluateOps` |
| `define` | `ActionDefineValue` |
| `boolean` | `ActionBoolean` |
| `string` | no completion (free text) |
| `int` | no completion (numeric) |
| `float` | no completion (numeric) |
| `degrees` | no completion (numeric) |
| `font` | `ActionFonts` |
| `format` | `ActionFormats` |
| `filename` | `ActionFiles` |
| `orientation` | `ActionOrientations` |
| `dispose` | `ActionDisposes` |
| `alpha` | `ActionAlphaOption` |
| `noise` | `ActionNoiseTypes` |
| `preview` | `ActionPreviewTypes` |
| `storage` | `ActionStorageTypes` |

---

## Phase 6: Completer CLIs (`cmd/`)

### 6.1 Main completer: `cmd/carapace-magick/`

- Uses `DefaultMagickProfile` (full option set, operators, stack ops, parentheses)
- `DisableFlagParsing` + `PositionalAnyCompletion` pattern (same as ffmpeg)
- Handles tool name as first positional token — if it matches a known tool, switches to that tool's profile

### 6.2 Sub-tool completers

Separate CLIs for direct invocation (e.g., `magick identify` can also be completed as its own command):

- [ ] `cmd/carapace-magick-identify/`
- [ ] `cmd/carapace-magick-mogrify/`
- [ ] `cmd/carapace-magick-compare/`
- [ ] `cmd/carapace-magick-composite/`
- [ ] `cmd/carapace-magick-montage/`

Each uses its corresponding `ToolProfile` with restricted option set.

### 6.3 Debug CLI: `cmd/carapace-magick-debug/`

Subcommands for testing/diagnostics:
- [ ] `argstream` / `argstream-complete` — parse full arg stream
- [ ] `definevalue` / `definevalue-complete` — parse `-define` values

---

## Phase 7: Man Pages (`man/magick/`)

YAML descriptions for completion value types — following the same format as carapace-ffmpeg.

High-priority value types to document first:
- [ ] `colorspace/` — ~40 entries (sRGB, CMYK, Gray, Lab, etc.)
- [ ] `compose/` — ~60 entries (Over, Multiply, Screen, etc.)
- [ ] `filter/` — ~30 entries (Lanczos, Gaussian, etc.)
- [ ] `channel/` — ~30 entries
- [ ] `gravity/` — ~10 entries
- [ ] `type/` — ~15 entries
- [ ] `distort/` — ~16 entries
- [ ] `metric/` — ~10 entries
- [ ] `evaluate/` — ~30 entries
- [ ] `virtual-pixel/` — ~15 entries

---

## Phase 8: AGENTS.md

Document the project for AI agents:
- [ ] Project overview, module path
- [ ] Build/test/lint commands
- [ ] Architecture diagram
- [ ] Key patterns and gotchas (same style as carapace-ffmpeg's AGENTS.md)
- [ ] Completion dispatch flow diagram
- [ ] Code conventions

---

## Implementation Order

Recommended sequence to get something working end-to-end early:

1. **Phase 1** — project scaffolding (can copy most from carapace-ffmpeg)
2. **Phase 2.1–2.5** — option definitions, AST, parser (core argstream)
3. **Phase 4.1–4.2** — dynamic and static actions (enough to complete basic options)
4. **Phase 5** — completer dispatch
5. **Phase 6.1** — main `carapace-magick` CLI (first working completer!)
6. **Phase 2.6–2.7** — completion parser + tests
7. **Phase 3** — definevalue parser
8. **Phase 6.2** — sub-tool completers
9. **Phase 6.3** — debug CLI
10. **Phase 7** — man pages (iterative, can be done in parallel)
11. **Phase 8** — AGENTS.md

---

## Open Questions

- **`-draw` primitive DSL**: Should we create a dedicated parser package for `-draw` primitives (like ffmpeg's filtergraph parser)? It's a sub-DSL with its own syntax (`circle`, `rectangle`, `path`, `image`, `text`). Start without it — treat as `ValueString` and add later.
- **`-channel-fx` expression DSL**: Same question — has its own grammar (`A=>R`, `R=>A`). Treat as `ValueString` initially.
- **`-fx` expression DSL**: Mathematical expression language. Very complex grammar — defer indefinitely, treat as `ValueString`.
- **Inline image generators**: `xc:red`, `gradient:`, `plasma:` — should these be completed as image inputs? They have their own sub-syntax. Low priority.
- **Image format prefix**: `png:filename` — the `format:` prefix on filenames. Should completion handle this? Probably as a `ActionMultiParts(":")` like ffmpeg's stream specifier.
