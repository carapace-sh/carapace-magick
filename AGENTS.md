# AGENTS.md

## Project Overview

Go library for parsing ImageMagick `magick` CLI argument streams into ASTs with shell completion support. Part of the [carapace-sh](https://github.com/carapace-sh) ecosystem (shell completion framework). The module path is `github.com/carapace-sh/carapace-magick`.

## Commands

### Build & Test

```sh
go test ./...                              # run all tests
go test ./pkg/argstream/                    # argstream tests only
go test ./pkg/completer/                    # completer tests only
go build ./...                              # build all packages
```

### CI Checks (mirrors `.github/workflows/go.yml`)

```sh
go build -v ./...                                       # build
go test -v -coverprofile=profile.cov ./...               # test with coverage
[ "$(gofmt -d -s . | tee -a /dev/stderr)" = "" ]         # format check (fails if any diffs)
staticcheck ./...                                        # lint
```

Both `gofmt` and `staticcheck` are enforced in CI. Do not skip them.

### Debug CLI (`cmd/carapace-magick-debug/`)

```sh
go run ./cmd/carapace-magick-debug argstream -- -resize 200x200 input.png output.png  # parse arg stream as JSON
go run ./cmd/carapace-magick-debug argstream-complete -- -resize                        # argstream completion context as JSON
go run ./cmd/carapace-magick-debug argstream-complete --profile identify -- -verbose    # identify profile completion context
go run ./cmd/carapace-magick-debug definevalue 'jpeg:quality=85'                      # parse -define value as JSON
go run ./cmd/carapace-magick-debug definevalue-complete 'jpeg:'                         # -define completion context
```

### Completer CLIs

```sh
# magick (default/convert)
go run ./cmd/carapace-magick _carapace spec                    # generate carapace spec YAML
go run ./cmd/carapace-magick _carapace export '' ''             # complete at empty position (JSON)

# identify
go run ./cmd/carapace-magick-identify _carapace spec

# mogrify
go run ./cmd/carapace-magick-mogrify _carapace spec

# compare
go run ./cmd/carapace-magick-compare _carapace spec

# composite
go run ./cmd/carapace-magick-composite _carapace spec

# montage
go run ./cmd/carapace-magick-montage _carapace spec
```

The `_carapace spec` command generates YAML that references the `man/` directory for extended descriptions.

## Architecture

Seven CLIs, a shared completer package, and one parser package with carapace completion actions.

```
cmd/carapace-magick/            Completer CLI for magick (default/convert)
cmd/carapace-magick-identify/   Completer CLI for magick identify
cmd/carapace-magick-mogrify/    Completer CLI for magick mogrify
cmd/carapace-magick-compare/    Completer CLI for magick compare
cmd/carapace-magick-composite/  Completer CLI for magick composite
cmd/carapace-magick-montage/    Completer CLI for magick montage
cmd/carapace-magick-debug/      Debug/diagnostic CLI (JSON output)
pkg/argstream/                  Argument stream parser (options, images, stack ops, parentheses)
pkg/completer/                  Shared completion dispatch logic
pkg/actions/tools/magick/       Carapace action functions for magick value types
pkg/definevalue/                -define format:key=value parser
pkg/probe/                      magick identify wrapper for image-aware completion
man/magick/                     YAML descriptions for completion value types
skills/magick/                  AI agent reference documentation (not compiled Go)
testdata/                       Test images for integration tests (go generate)
```

### Completer CLIs

Standalone carapace completers for `magick` and its sub-tools. Each uses `DisableFlagParsing` + `PositionalAnyCompletion` with `argstream.ParseForCompletionWithProfile()` and a tool-specific `ToolProfile` to dispatch context-aware completions.

- **`root.go`** — Root cobra command with `carapace.Gen(rootCmd).Standalone()`. `PositionalAnyCompletion` callback parses args with `argstream.ParseForCompletionWithProfile()` and dispatches to shared actions from `pkg/completer/`.
- **`main.go`** — Entry point calling `cmd.Execute()`.

### Argument Stream (`pkg/argstream/`)

- **`parser.go`** — Full parser. `Parse(args)` → `*Program` AST. `ParseWithProfile(args, profile)` allows sub-tool profiles. Tokenizes argument list into options, images, parentheses, and stack operations.
- **`completion_parser.go`** — Completion parser. `ParseForCompletion(args, trailingSpace)` → `*CompletionContext`. `ParseForCompletionWithProfile(args, trailingSpace, profile)` allows sub-tool profiles.
- **`completion.go`** — Completion context types (`ExpectedToken`, `OptionContext`, `CompletionContext`). JSON-serializable with `json` tags.
- **`ast.go`** — AST node types (`Token`, `TokenKind`, `OptionForm`, `ParenGroup`, `Program`). JSON-serializable.
- **`options.go`** — Static option definitions for magick (`OptionDef`, `OptionCategory`, `OptionType`, `ValueType`, `PlusBehavior`, `OptionIndex`).
- **`profile.go`** — `ToolProfile` struct with `Name`, `HasOutputArg`, `HasOperators`, `HasStackOps`, `HasParentheses`, `OptionIndex`. Defines `DefaultMagickProfile`, `DefaultIdentifyProfile`, etc. Each sub-tool has its own option set.
- **`span.go`** — `Span` type.

### Completer (`pkg/completer/`)

- **`completer.go`** — Shared completion actions used by all completer CLIs. Key functions:
  - `ContextToArgs(c carapace.Context) (args []string, trailingSpace bool)` — converts carapace context to argstream input.
  - `ActionOptions(ctx, profile)` — option name completions (both `-` and `+` forms).
  - `ActionOptionValue(ctx)` — giant switch on `ValueType` dispatching to the correct magick action.
  - `ActionToolNames()` — sub-tool name completions for the first positional arg.
  - `ActionDefineValue(partial)` — structured completion for `-define` arguments using definevalue parser.
  - `ActionDefineKeys(format)` / `ActionDefineValues(format, key)` — format-specific define key/value completions.

### Define Value (`pkg/definevalue/`)

Parser for `-define format:key=value` argument strings.

- **`parser.go`** — Strict parser. `Parse(input)` → `*DefineValue`. Splits at `:` and `=` boundaries.
- **`completion_parser.go`** — Completion parser. `ParseForCompletion(input)` → `*CompletionContext`.
- **`completion.go`** — Completion context types (`ExpectedToken`, `CompletionContext`).
- **`ast.go`** — `DefineValue` AST type with `Format`, `Key`, `Value` fields.
- **`defines.go`** — Format-specific define key data (`FormatDefines`, `GlobalDefines`). Lists keys for JPEG, PNG, TIFF, WebP, GIF, HEIC, PSD, PDF, RAW formats.
- **`span.go`** — `Span` type.

### Probe (`pkg/probe/`)

Wraps `magick identify -verbose` for image-aware completion.

- **`probe.go`** — `Probe(inputPath)` → `*ImageInfo` with Width, Height, Format, Colors, Depth, Colorspace. Returns `nil` when magick is unavailable or file doesn't exist.

### Man Pages (`man/magick/`)

YAML descriptions for completion value types. Directory structure: `man/magick/<type>/<type>.yaml`. Each YAML file maps values to multiline descriptions.

Completed value types: colorspace, compose, filter, channel, gravity, type, distort, metric, evaluate, virtual-pixel, boolean, alpha, noise, preview, storage, orientation, dispose, auto-threshold, direction, kernel.

### Actions (`pkg/actions/tools/magick/`)

Carapace completion actions for magick value types. All actions use `.Tag()` and `.Uid()`/`.UidF()` for caching/identification.

- **`value_types.go`** — Exported action functions. **Dynamic actions** shell out to `magick -list` (colorspaces, composes, compress types, channels, distort methods, filters, gravities, interlace types, layer methods, morphology methods, types, virtual pixel methods, orientations, disposes, metrics, evaluate ops, formats, fonts, colors, kernels). **Static actions** use `carapace.ActionValues`/`ActionValuesDescribed` (booleans, alpha options, auto-threshold, noise types, preview types, storage types, directions, list types, tool names).

## Key Patterns & Gotchas

### magick's positional argument model

Unlike traditional CLI flag trees, magick arguments form a **linear stream** where settings configure future behavior, operators transform images immediately, and stack operators manipulate the image sequence. Parentheses `()` create scoped sub-pipelines.

### Option `-` vs `+` prefix

The `+` prefix is NOT just "disable" — its meaning depends on the option type:
- **Setting (value)**: `+background` resets to default (takes no value)
- **Setting (boolean)**: `+verbose` disables
- **Operator (directional)**: `-append` vertical / `+append` horizontal

The `needsValue()` function determines whether a `+form` takes a value based on `PlusBehavior`.

### No stream specifiers or short options

Unlike ffmpeg, magick has **no short single-letter options** and **no stream specifier suffix** after option names. Every option is a full word like `-resize`, `-gaussian-blur`, etc. This simplifies lexing significantly.

### Tool name as first positional token

When using `magick` (the default profile), the first positional non-option token can be a sub-tool name like `identify`, `mogrify`, etc. The completion parser detects this and sets `ctx.Tool` accordingly. For sub-tool profiles, the tool is pre-set.

### Profile isolation

Each sub-tool profile has its **own** `OptionIndex` — options not listed in a profile are unavailable. The `LookupOption` method on `ToolProfile` does NOT fall back to the main `OptionIndex` for sub-tool profiles.

### `trailingSpace` is critical for completion

`ParseForCompletion(args, trailingSpace)` behaves differently based on `trailingSpace`:
- `true` — cursor is at a new blank position after the last token
- `false` — cursor is mid-token within the last argument

### Completer uses `DisableFlagParsing` + `PositionalAnyCompletion`

The completer CLI does NOT use cobra's flag parsing. It sets `DisableFlagParsing: true` so cobra hands all arguments through as positional args.

### UIDs use `magick://` scheme

All completion actions use `magick://` UIDs for carapace's action deduplication.

## Code Conventions

- **Standard library only for parsers**: The argstream parser package uses only Go standard library. No external dependencies.
- **Carapace + Cobra for CLIs and actions**: External deps (`carapace`, `carapace-spec`, `cobra`) are only in `cmd/` and `pkg/actions/`.
- **Test style**: Table-driven tests with `testing.T` only. No testify or other assertion libraries.
- **Action test style**: Tests in `pkg/actions/tools/magick/` should use carapace's `sandbox.Action()` framework.

## Release

GoReleaser builds 7 binaries: `carapace-magick`, `carapace-magick-debug`, `carapace-magick-identify`, `carapace-magick-mogrify`, `carapace-magick-compare`, `carapace-magick-composite`, `carapace-magick-montage`. Distribution channels: Homebrew tap, Scoop bucket, AUR, nfpm (apk/deb/rpm/termux.deb), Gemfury. Releases are triggered by tag pushes in CI.
