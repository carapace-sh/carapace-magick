# Multi-Completer Plan for carapace-magick

## Goal

Consolidate the 6 separate completer binaries into a single `carapace-magick` binary with subcommands, custom multi-completer snippet generation, and proper bridge integration with carapace-bin.

## Completed: Multi-Completer Restructure

Done in commit `fb525ec`. Single binary with subcommands `magick`, `identify`, `mogrify`, `compare`, `composite`, `montage`, `debug`. Custom per-shell snippet templates in `cmd/carapace-magick/cmd/snippet/`. Shell auto-detection via `ps.DetermineShell()`.

## Completed: Bridge Routing via Arg Rewriting

Done in commit `c8005a6`. Bridge routing in `Execute()` that rewrites `os.Args` to route `bridge.ActionCarapace` invocations to the correct cobra subcommand.

### How bridge.ActionCarapace Constructs Its Invocation

`bridge.ActionCarapace("carapace-magick", "identify")` in `carapace-bridge/pkg/actions/bridge/carapace.go`:

```go
args := []string{"_carapace", "export", ""}
args = append(args, command[1:]...)  // "identify" (the subcommand hint)
args = append(args, c.Args...)       // user's already-typed args
args = append(args, c.Value)         // current word being completed
// executes: carapace-magick _carapace export "" identify -verbose image.png
```

The `command[1]` value (e.g. "identify") is always spliced in at position 4 (`os.Args[4]`), before the user's args. This is the key to routing — the subcommand name is always at a known position.

### The Arg Rewriting Logic

When `Execute()` intercepts `_carapace` with 4+ args (completion/export request), it checks `os.Args[4]` against a list of known completer subcommand names. If matched, it rewrites `os.Args` to route through cobra:

```
bridge.ActionCarapace("carapace-magick", "identify") calls:
  carapace-magick _carapace export "" identify -verbose image.png

Rewriting detects "identify" at os.Args[4], produces:
  carapace-magick identify _carapace export "" -verbose image.png
```

```
bridge.ActionCarapace("carapace-magick", "magick") calls:
  carapace-magick _carapace export "" magick -resize 200

Rewriting detects "magick" at os.Args[4], produces:
  carapace-magick magick _carapace export "" -resize 200
```

When `os.Args[4]` is NOT a known subcommand (e.g. `ActionCarapace("carapace-magick")` without a subcommand hint, where `os.Args[4]` = "-resize"), it defaults to the `magick` subcommand:

```
bridge.ActionCarapace("carapace-magick") calls:
  carapace-magick _carapace export "" -resize 200

No match at os.Args[4], defaults to magick:
  carapace-magick magick _carapace export "" -resize 200
```

### The os.Args[4] Ambiguity Is Not a Problem

`os.Args[4]` can be both a known subcommand name AND a valid user arg. For example, when a user types `magick identify -verbose`, the bridge produces:

```
carapace-magick _carapace export "" magick identify -verbose
```

Here `os.Args[4]` = "magick" (the `command[1]` value, NOT the user's arg). The rewriting correctly picks subcommand = "magick" and removes it:

```
carapace-magick magick _carapace export "" identify -verbose
```

The user's "identify" is preserved at `os.Args[5]` and reaches the argstream as a tool name — which is correct behavior for the `magick` profile where `identify` can be a sub-tool name.

This works because `bridge.ActionCarapace` always inserts `command[1]` at position 4. The user's args start at position 5. So `os.Args[4]` is always the bridge subcommand hint, never a user arg.

### The Subcommand Path After Rewriting

After rewriting, the subcommand's `_carapace` handler calls `complete(parentCmd, args)`. The `args` layout:

```
magick subcommand:  args = [export, "", -resize, 200]
                     args[2:] = [-resize, 200]  (argstream: option + partial value)

identify subcommand: args = [export, "", -verbose, image.png]
                      args[2:] = [-verbose, image.png]  (argstream: option + image)
```

The empty string at `args[1]` is the cobra subcommand path (empty = root of that subcommand). `args[2:]` are the actual positional args for the argstream parser.

### Implementation in root.go

```go
func Execute() {
    if len(os.Args) > 1 && os.Args[1] == "_carapace" {
        if len(os.Args) < 4 {
            // snippet request: carapace-magick _carapace [shell]
            ...
            return
        }
        // completion/export request — route to correct subcommand
        subcommand := "magick" // default
        if len(os.Args) > 4 && isCompleterSubcommand(os.Args[4]) {
            subcommand = os.Args[4]
            os.Args = append(
                []string{os.Args[0], subcommand, "_carapace", os.Args[2], os.Args[3]},
                os.Args[5:]...,
            )
        } else {
            os.Args = append(
                []string{os.Args[0], subcommand, "_carapace"},
                os.Args[2:]...,
            )
        }
    }
    // rootCmd.Execute() routes to the subcommand's cobra command,
    // which has its own _carapace handler via carapace.Gen(subcmd).Standalone()
    rootCmd.Execute()
}

func isCompleterSubcommand(name string) bool {
    return slices.Contains([]string{"magick", "identify", "mogrify", "compare", "composite", "montage"}, name)
}
```

## Completed: carapace-bin Bridge Stubs

Done in commit `914f1769d` on branch `magick-bridge-stubs`.

### magick_completer (replaced)

The native magick completer (~300 lines of cobra flags + action map) was replaced with a bridge stub:

```go
carapace.Gen(rootCmd).PositionalAnyCompletion(
    carapace.ActionCallback(func(c carapace.Context) carapace.Action {
        if _, err := exec.LookPath("carapace-magick"); err == nil {
            return bridge.ActionCarapace("carapace-magick", "magick")
        }
        return bridge.ActionBridge("magick")
    }),
)
```

### 5 New Stubs

- `identify_completer` → `bridge.ActionCarapace("carapace-magick", "identify")`
- `mogrify_completer` → `bridge.ActionCarapace("carapace-magick", "mogrify")`
- `compare_completer` → `bridge.ActionCarapace("carapace-magick", "compare")`
- `composite_completer` → `bridge.ActionCarapace("carapace-magick", "composite")`
- `montage_completer` → `bridge.ActionCarapace("carapace-magick", "montage")`

Each falls back to `bridge.ActionBridge("<name>")` when `carapace-magick` is not in PATH.

### Completion Path in carapace-bin

Shell completion uses the `carapace <completer> <shell>` path, NOT `carapace <completer> _carapace export`. The flow:

1. Shell calls `carapace identify bash -verbose ''` (via snippet callback)
2. carapace-bin routes to `invokeCompleter("identify")`
3. `invokeCompleter` calls `completer.Execute()` with `os.Args[1] = "_carapace"`
4. The completer's cobra processes `_carapace export bash '' -verbose ''`
5. carapace's `complete()` calls `traverse()` → `PositionalAnyCompletion` callback
6. `bridge.ActionCarapace("carapace-magick", "identify")` executes as subprocess
7. The subprocess `carapace-magick _carapace export "" identify -verbose ""` hits our arg rewriting
8. Rewriting routes to `carapace-magick identify _carapace export "" -verbose ""`
9. The identify subcommand's _carapace handler produces the completion output

The `_carapace export` path is only used internally by carapace's `_carapace` subcommand for re-invocation. It does NOT work for bridge completers in carapace-bin (the re-invoked binary loses the completer name context). This is fine because shell completion uses the direct `carapace <completer> <shell>` path.

## Future: SnippetMulti in carapace Library

Once this pattern is proven, extract the multi-completer snippet generation and bridge routing into a public API in the carapace library. This would enable any carapace-based completer to become a multi-completer without duplicating shell templates or arg rewriting logic.

### Proposed API

```go
// carapace package
func GenMulti(rootCmd *cobra.Command, subcommands []*cobra.Command, opts ...MultiOption) {
    // 1. Registers all subcommands on rootCmd
    // 2. Adds _carapace interception in Execute() for snippet + routing
    // 3. Generates multi-completer snippets for all shells
    // 4. Handles arg rewriting for bridge.ActionCarapace routing
}
```

### What Needs to Be Extracted

1. **Arg rewriting logic** — The `os.Args` rewriting in `Execute()` that detects the subcommand name at `os.Args[4]` and routes to the correct cobra subcommand. This is the critical piece for bridge routing. It must know the list of completer subcommand names to detect at `os.Args[4]`.

2. **Snippet generation** — The per-shell templates in `cmd/carapace-magick/cmd/snippet/`. These are simpler than carapace-bin's templates (no `pathSnippet`/`envSnippet`/`CARAPACE_SHELL_*` env vars). A shared completer function that calls `<binary> <command> _carapace <shell>` and registers all command names.

3. **Shell auto-detection** — `ps.DetermineShell()` for when `_carapace` is called without a shell arg.

### Key Design Constraints for the Generic Solution

1. **`os.Args[4]` is always the subcommand hint** — `bridge.ActionCarapace("binary", "subcommand")` always places the subcommand name at `os.Args[4]`, before any user args. The rewriting logic can safely check this position against a known list.

2. **The root command has no `carapace.Gen()`** — We intercept `_carapace` ourselves before cobra processes it. The root is just a dispatcher. Each completer subcommand has its own `carapace.Gen(subcmd).Standalone()`.

3. **Snippet template structure** — Each shell template creates a shared completer function that calls `<binary> <command> _carapace <shell>` and registers all command names. Much simpler than carapace-bin's templates.

4. **Fallback default** — When `os.Args[4]` is not a known subcommand, route to a default subcommand (typically the primary one, like `magick`).

5. **PowerShell backtick gofmt issue** — PowerShell templates contain backtick characters that `gofmt -s` normalizes differently. Must use a `const` raw string literal (backtick-delimited) to avoid gofmt mangling the template. This is the same pattern used in the carapace library's own PowerShell snippet.

### The Generic Routing Logic

For a generic `GenMulti`, the arg rewriting would look like:

```go
func Execute() {
    if len(os.Args) > 1 && os.Args[1] == "_carapace" {
        if len(os.Args) < 4 {
            // snippet request
            fmt.Println(multiSnippet(shell, subcommandNames))
            return
        }
        // completion/export request — route to correct subcommand
        subcommand := defaultSubcommand
        if len(os.Args) > 4 && isMultiSubcommand(os.Args[4]) {
            subcommand = os.Args[4]
            os.Args = append(
                []string{os.Args[0], subcommand, "_carapace", os.Args[2], os.Args[3]},
                os.Args[5:]...,
            )
        } else {
            os.Args = append(
                []string{os.Args[0], subcommand, "_carapace"},
                os.Args[2:]...,
            )
        }
    }
    rootCmd.Execute()
}
```

This is identical to the current carapace-magick implementation — the only parameterization needed is the list of subcommand names and the default subcommand.
