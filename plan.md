# Multi-Completer Plan for carapace-magick

## Goal

Consolidate the 6 separate completer binaries (`carapace-magick`, `carapace-magick-identify`, `carapace-magick-mogrify`, `carapace-magick-compare`, `carapace-magick-composite`, `carapace-magick-montage`) into a single `carapace-magick` binary with subcommands for each completer, plus the debug subcommand.

The resulting binary would be invoked as:
- `carapace-magick magick _carapace bash ...` — completions for `magick`
- `carapace-magick identify _carapace bash ...` — completions for `identify`
- `carapace-magick mogrify _carapace bash ...` — completions for `mogrify`
- etc.

A single shell snippet from `carapace-magick _carapace bash` would register ALL commands at once.

## The Core Problem

When carapace's standard `_carapace` command generates a snippet, it produces a **single-command** snippet. For bash, that's:

```bash
_carapace-magick_completion() { ... carapace-magick _carapace bash ... }
complete -o noquote -F _carapace-magick_completion carapace-magick
```

But we need it to register `magick`, `identify`, `mogrify`, `compare`, `composite`, `montage` — not `carapace-magick`. The snippet callback must dispatch to the right subcommand:

```bash
_carapace_magick_completer() { ... carapace-magick "${command}" _carapace bash ... }
complete -o noquote -F _carapace_magick_completer magick identify mogrify compare composite montage
```

## Approaches

### Approach A: Patch the snippet (carapace-bin style)

Override the `_carapace` snippet generation in `carapace-magick`'s root command to produce a custom multi-completer snippet. Write per-shell snippet templates directly in `carapace-magick` (like carapace-bin's `snippet/bash.go`, `snippet/zsh.go`, etc.).

**How carapace-bin does it**: The `carapace` binary's `Execute()` intercepts `_carapace` with < 4 args, calls `snippet.Snippet(shell)` which generates a shared completer function calling `carapace "${command}" <shell>` and registers all completer names against it. For single-completer invocation (`carapace git bash`), it runs the completer's `Execute()` with `_carapace` trick, captures the default snippet, and patches `_carapace` → `git` via `strings.ReplaceAll`.

**How carapace-magick would do it**: Simpler — we intercept `_carapace` at the root level. When called with 0-1 args (snippet request), we generate a custom multi-completer snippet. When called with 2+ args, the subcommand's own `_carapace` handles completion normally.

**Pros**: Works today with no carapace library changes; proven pattern
**Cons**: Duplicates ~8 shell snippet templates; must maintain them alongside carapace's; no `pathSnippet`/`envSnippet` needed but still boilerplate

### Approach B: Add multi-completer support to the carapace library

Add a new API like `carapace.Gen(rootCmd).MultiCompleter(commandNames ...string)` or expose a `SnippetMulti(shell, names)` function that generates a multi-completer snippet.

The library already has `internal/shell/*/snippet.go` — the per-shell snippet generators are in `internal/` (not importable). We'd need to:
1. Either expose the snippet generators publicly (e.g., move to `pkg/shell/`)
2. Or add a `SnippetMulti` method that the internal generators handle

**Pros**: Reusable by any project; no duplicate templates; single source of truth for shell snippets
**Cons**: Requires carapace library changes; the snippet generators are fairly coupled to the single-command model; `pathSnippet`/`envSnippet` would also need to be exposed or skipped

### Approach C: Hybrid — generate individual snippets and merge

Call `Gen(subcmd).Snippet(shell)` for each subcommand to get individual snippets, then merge them. Doesn't work well because each snippet is self-contained with its own function name and registration line. Merging bash/zsh/fish snippets is fragile string manipulation.

**Verdict**: Not viable.

### Approach D: Use `strings.ReplaceAll` patching on the root snippet

Generate the root command's snippet (which references `carapace-magick` as the command name), then patch it:
1. Replace `_carapace` in the callback with `<subcommand> _carapace` 
2. Replace the single command registration with multiple registrations

This is what carapace-bin does for single-completer dispatch but adapted. The challenge is that the snippet templates embed `_carapace` in specific positions and the patching would need to be precise per shell.

**Pros**: Reuses carapace's snippet templates; less code to write
**Cons**: Fragile string replacement; different patching rules per shell; easy to break when carapace updates its snippet format

## Recommendation: Approach A (custom snippets, minimal duplication)

carapace-magick's snippet templates are **much simpler** than carapace-bin's because:
- No `pathSnippet()` (no `~/.config/carapace/bin/` PATH manipulation)
- No `envSnippet()` (no `get-env`/`set-env` helpers)
- No `CARAPACE_SHELL_*` env vars (no bridge support needed)
- Fixed set of 6 command names (no dynamic discovery)

The templates are just: one shared completer function + registration lines. ~20-30 lines per shell.

**Future improvement**: Once this is proven, we can extract the pattern into the carapace library as Approach B — a `SnippetMulti` API that other projects can use.

## Implementation Plan

### Step 1: Restructure `cmd/carapace-magick/` to have subcommands

Create a root command `carapace-magick` with subcommands:

```
carapace-magick (root)
  ├── magick      (DisableFlagParsing, PositionalAnyCompletion with DefaultMagickProfile)
  ├── identify    (DisableFlagParsing, PositionalAnyCompletion with DefaultIdentifyProfile)
  ├── mogrify     (DisableFlagParsing, PositionalAnyCompletion with DefaultMogrifyProfile)
  ├── compare     (DisableFlagParsing, PositionalAnyCompletion with DefaultCompareProfile)
  ├── composite   (DisableFlagParsing, PositionalAnyCompletion with DefaultCompositeProfile)
  ├── montage     (DisableFlagParsing, PositionalAnyCompletion with DefaultMontageProfile)
  └── debug       (existing debug CLI with argstream/argstream-complete/definevalue subcommands)
```

Each subcommand's init function contains the same `PositionalAnyCompletion` callback that currently lives in the separate binaries. The `debug` subcommand keeps its existing structure with `spec.Register`.

Key detail: The root command itself has **no** `carapace.Gen()` — we handle snippet generation ourselves. Each subcommand has `carapace.Gen(subcmd).Standalone()` so `_carapace` works normally for actual completion dispatch.

### Step 2: Custom multi-completer snippet generation

In `cmd/carapace-magick/snippet/`:

- `snippet.go` — `Snippet(shell string) string` dispatcher
- `bash.go` — bash multi-completer snippet template
- `zsh.go` — zsh multi-completer snippet template
- `fish.go` — fish multi-completer snippet template
- `elvish.go` — elvish multi-completer snippet template
- `nushell.go` — nushell multi-completer snippet template
- `powershell.go` — powershell multi-completer snippet template
- `xonsh.go` — xonsh multi-completer snippet template

The command names are hardcoded: `["magick", "identify", "mogrify", "compare", "composite", "montage"]`.

### Step 3: Intercept `_carapace` on the root command

Override the root command's `Run` to handle snippet generation:

```go
func Execute() error {
    if len(os.Args) > 1 && os.Args[1] == "_carapace" && len(os.Args) < 4 {
        // snippet request: carapace-magick _carapace [shell]
        shell := "bash"
        if len(os.Args) > 2 {
            shell = os.Args[2]
        }
        fmt.Println(snippet.Snippet(shell))
        return nil
    }
    return rootCmd.Execute()
}
```

This intercepts `_carapace` at the root level before cobra dispatches. For actual completion, `carapace-magick identify _carapace bash ...` flows through cobra to the `identify` subcommand's `_carapace` handler normally.

### Step 4: Remove separate cmd directories

Delete:
- `cmd/carapace-magick-identify/`
- `cmd/carapace-magick-mogrify/`
- `cmd/carapace-magick-compare/`
- `cmd/carapace-magick-composite/`
- `cmd/carapace-magick-montage/`
- `cmd/carapace-magick-debug/`

All their logic moves into subcommands of `cmd/carapace-magick/`.

### Step 5: Update GoReleaser

Single binary build:

```yaml
builds:
  - id: carapace-magick
    main: ./cmd/carapace-magick
    binary: carapace-magick
```

### Step 6: Update distribution (AUR, Homebrew, Scoop, nfpm)

Single package instead of 6. The `carapace-magick` binary is the only artifact.

## Open Questions

1. **Should we add a library feature to carapace?** A `SnippetMulti(names ...string)` API would avoid duplicating snippet templates across projects. This is worth doing as a follow-up after proving the approach works in carapace-magick. The library's `internal/shell/` snippet generators would need to accept a list of names and generate a shared completer function.

2. **Should the `magick` subcommand be the default?** When users run `carapace-magick` without a subcommand, should it default to `magick`? Or print help? I'd suggest printing help (cobra's default behavior).

3. **Backward compatibility**: Users who have `carapace-magick-identify` etc. in their PATH would need to remove those and use the single binary. The old separate binaries would no longer be released.

4. **`magick` vs `convert`**: The `magick` subcommand handles both `magick` and `magick convert`. Should we also register `convert` as a completion target? In practice, `magick convert` is just `magick` with the tool name as the first positional arg (already handled by `ExpectedToolName`).

5. **Shell-specific edge cases**: bash-ble, oil, tcsh, cmd-clink — should we support all shells that carapace supports, or just the main ones (bash, zsh, fish, elvish, nushell, powershell, xonsh)? Supporting all means more template code. Could start with the main ones and add others later.
