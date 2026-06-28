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

### Debug CLI (`debug` subcommand)

```sh
go run ./cmd/carapace-magick debug argstream -- -resize 200x200 input.png output.png  # parse arg stream as JSON
go run ./cmd/carapace-magick debug argstream-complete -- -resize                        # argstream completion context as JSON
go run ./cmd/carapace-magick debug argstream-complete --profile identify -- -verbose    # identify profile completion context
go run ./cmd/carapace-magick debug definevalue 'jpeg:quality=85'                      # parse -define value as JSON
go run ./cmd/carapace-magick debug definevalue-complete 'jpeg:'                         # -define completion context
```

### Completer CLI

```sh
# Multi-completer snippet (registers all 6 ImageMagick commands + carapace-magick itself)
go run ./cmd/carapace-magick _carapace bash

# Single-command snippet (registers only one command)
go run ./cmd/carapace-magick carapace-magick _carapace bash   # carapace-magick itself
go run ./cmd/carapace-magick identify _carapace bash          # identify only

# Bridge completion (used by carapace-bin via bridge.ActionCarapace)
go run ./cmd/carapace-magick _carapace export '' identify -verbose image.png

# Pseudo-subcommand completion (carapace-magick as self-completer)
go run ./cmd/carapace-magick carapace-magick _carapace export bash '' '' ''           # root subcommands
go run ./cmd/carapace-magick carapace-magick _carapace export bash '' debug argstream-complete '-'  # debug flags
```

## Architecture

Single binary with subcommands for each ImageMagick completer, a pseudo-subcommand for self-completion, a shared completer package, and parser packages with carapace completion actions.

```
cmd/carapace-magick/                  Single binary with multi-completer subcommands
  cmd/root.go                        Root command, Execute() interception for bridge routing, pseudo-subcommand, and snippet dispatch
  cmd/magick.go                      magick (default/convert) completer subcommand
  cmd/identify.go                    identify completer subcommand
  cmd/mogrify.go                     mogrify completer subcommand
  cmd/compare.go                     compare completer subcommand
  cmd/composite.go                   composite completer subcommand
  cmd/montage.go                     montage completer subcommand
  cmd/debug.go                       debug subcommand (argstream, definevalue parsers)
  cmd/snippet/                       Shell snippet generators
    snippet.go                       Snippet() and SingleSnippet() dispatcher
    bash.go, zsh.go, fish.go, etc.   Per-shell snippet templates (multi-completer and single-command)
pkg/argstream/                        Argument stream parser (options, images, stack ops, parentheses)
pkg/completer/                        Shared completion dispatch logic
pkg/actions/tools/magick/             Carapace action functions for magick value types
pkg/definevalue/                      -define format:key=value parser
pkg/probe/                            magick identify wrapper for image-aware completion
man/magick/                           YAML descriptions for completion value types
skills/magick/                        AI agent reference documentation (not compiled Go)
testdata/                             Test images for integration tests (go generate)
```

### Multi-Completer Architecture

`carapace-magick` is a single binary that acts as a multi-completer with three distinct `_carapace` dispatch paths handled in `Execute()` before cobra sees the args:

1. **Bridge routing** (`carapace-magick _carapace export "" <subcommand> ...`) — Used by `carapace-bin` via `bridge.ActionCarapace("carapace-magick", "identify")`. Rewrites `os.Args` to route to the correct completer subcommand (defaults to `magick` if no subcommand specified).

2. **Pseudo-subcommand** (`carapace-magick carapace-magick _carapace ...`) — Handles self-completion without a cobra command. For root-level completion, strips the pseudo-subcommand and falls through to `rootCmd.Execute()`, where `carapace.Gen(rootCmd).Standalone()` + `PositionalAnyCompletion` returns subcommand names. For deeper navigation (e.g. `debug argstream-complete`), rewrites `os.Args` to route to the actual subcommand's cobra command. Must handle two call formats:
   - **Shell format**: `_carapace <shell> <args...>` — user args start at `os.Args[4]`
   - **Export format**: `_carapace export <shell> "" <args...>` — user args start at `os.Args[6]`

3. **Subcommand-level snippet** (`carapace-magick <subcommand> _carapace <shell>`) — Returns a single-command snippet via `snippet.SingleSnippet()`.

The pseudo-subcommand is NOT a cobra command because adding one with `carapace.Gen().Standalone()` causes a `_carapace` re-invocation loop: carapace's `PositionalAnyCompletion` re-invokes the binary via `ActionExecCommand(os.Executable(), args...)`, and for a pseudo-subcommand cobra command this chain defaults to `magick` instead of staying within the pseudo-subcommand.

- **`root.go`** — Root cobra command with `carapace.Gen(rootCmd).Standalone()` and `PositionalAnyCompletion` returning subcommand names. `Execute()` intercepts `_carapace` calls with arg rewriting before dispatching to cobra.
- **`magick.go` / `identify.go` / etc.** — Completer subcommands. Each uses `DisableFlagParsing` + `PositionalAnyCompletion` with `argstream.ParseForCompletionWithProfile()` and a tool-specific `ToolProfile`.
- **`debug.go`** — Debug subcommand with `argstream`, `argstream-complete`, `definevalue`, `definevalue-complete` sub-subcommands. Uses `carapace-spec` for spec generation.
- **`snippet/`** — Per-shell snippet templates with both multi-completer (`Snippet()`) and single-command (`SingleSnippet()`) generators. The multi-completer snippet dispatches via `${command}` (the invoked binary name), while single-command snippets use the completer name directly.

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

Each completer subcommand does NOT use cobra's flag parsing. It sets `DisableFlagParsing: true` so cobra hands all arguments through as positional args.

### Pseudo-subcommand is NOT a cobra command

The `"carapace-magick"` pseudo-subcommand is handled purely in `Execute()` by checking `os.Args[1]`. Adding it as a cobra command with `carapace.Gen().Standalone()` causes a `_carapace` re-invocation loop because `PositionalAnyCompletion` re-invokes the binary, and the chain defaults to `magick` instead of staying within the pseudo-subcommand. Root-level completion falls through to `rootCmd.Execute()` where `carapace.Gen(rootCmd).Standalone()` handles shell formatting.

### `_carapace` call formats differ by context

The pseudo-subcommand handler must distinguish between two arg layouts:
- **Shell format** (from shell snippet functions): `_carapace <shell> <args...>` — user args at `os.Args[4:]`
- **Export format** (from carapace's internal re-invocation): `_carapace export <shell> "" <args...>` — user args at `os.Args[6:]`

Using the wrong index causes off-by-one bugs where root subcommand names are returned instead of the target subcommand's completions.

### UIDs use `magick://` scheme

All completion actions use `magick://` UIDs for carapace's action deduplication.

## Code Conventions

- **Standard library only for parsers**: The argstream parser package uses only Go standard library. No external dependencies.
- **Carapace + Cobra for CLIs and actions**: External deps (`carapace`, `carapace-spec`, `cobra`) are only in `cmd/` and `pkg/actions/`.
- **Test style**: Table-driven tests with `testing.T` only. No testify or other assertion libraries.
- **Action test style**: Tests in `pkg/actions/tools/magick/` should use carapace's `sandbox.Action()` framework.

## Release

GoReleaser builds a single `carapace-magick` binary. Distribution channels: Homebrew tap, Scoop bucket, AUR, nfpm (apk/deb/rpm/termux.deb), Gemfury. Releases are triggered by tag pushes in CI.
