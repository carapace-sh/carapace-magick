---
name: magick
description: >
  Use when working with ImageMagick magick CLI argument lexing or completion — the image pipeline
  command structure, option syntax, +/- flag forms, parentheses grouping, image stack operators,
  value types, option categories, and sub-tool differences. Triggers on: "magick", "magick cli",
  "magick arguments", "magick options", "magick flags", "magick pipeline", "magick completion",
  "magick lexer", "magick syntax", "magick value types", "imagemagick", "magick convert",
  "magick identify", "magick mogrify", "magick compare", "magick composite", "magick montage",
  "magick conjure", "magick stream", "-define", "image stack", "parentheses grouping",
  "geometry syntax", "magick quoting", "channel mask".
user-invocable: true
---

# ImageMagick `magick` CLI In-Depth Reference

ImageMagick's `magick` command-line is an **image pipeline processor**, not a traditional flag tree. Arguments form an interleaved sequence of image inputs, settings, operators, and image stack operations — where options apply to the *current* image in the processing pipeline and parentheses create nested image scopes. This skill provides the reference material needed to write a lexer for `magick` CLI arguments.

## Data Flow

```
global_settings
  → [settings | image_input | operators | stack_ops | parentheses_group] ...
  → output_image
```

Arguments are processed left-to-right. Settings configure how subsequent operators behave. Operators transform the current image(s). Stack operators manipulate the image sequence. Parentheses `(…)` create scoped sub-pipelines whose images merge back into the outer sequence.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| argstream, command structure, image pipeline, image sequence, input, output, positional, lexing model, token classification, image stack | [references/argstream.md](references/argstream.md) |
| option syntax, flag, dash, plus prefix, boolean flag, +/- forms, setting, operator, value-taking option, quoting, escaping | [references/option-syntax.md](references/option-syntax.md) |
| option scope, settings, operators, channel operators, sequence operators, stack operators, option categories, when options apply, positional scope | [references/option-scopes.md](references/option-scopes.md) |
| value types, geometry, color, percent, threshold, integer, float, string, expression, point, degrees, list types, -list | [references/value-types.md](references/value-types.md) |
| identify, mogrify, compare, composite, montage, conjure, stream, animate, convert, sub-tool, tool, legacy command, magick subcommand | [references/subtools.md](references/subtools.md) |
| define, -define, format:option, coder option, delegate option, key:value, define syntax, format-specific option | [references/define-syntax.md](references/define-syntax.md) |

## Quick Guide

- **How does the magick argstream work end-to-end?** → [references/argstream.md](references/argstream.md)
- **How do I parse an option name and determine +/- forms?** → [references/option-syntax.md](references/option-syntax.md)
- **Which options are settings vs operators vs stack ops?** → [references/option-scopes.md](references/option-scopes.md)
- **What value types does magick use for option arguments?** → [references/value-types.md](references/value-types.md)
- **How does `magick identify` differ from `magick convert`?** → [references/subtools.md](references/subtools.md)
- **How do I lex a `-define format:key=value` argument?** → [references/define-syntax.md](references/define-syntax.md)
- **How do parentheses and the image stack work?** → [references/argstream.md](references/argstream.md) and [references/option-scopes.md](references/option-scopes.md)
- **How does quoting/escaping work in magick option values?** → [references/option-syntax.md](references/option-syntax.md)

## Cross-Project References

- For **shell quoting rules** (bash/zsh/fish escaping that wraps magick arguments), see the **bash**, **zsh**, or **fish** skills.
- For **carapace completion framework** internals, see the **carapace** skill.
- For **ffmpeg CLI argstream patterns** (similar positional pipeline model), see the **ffmpeg** skill.
