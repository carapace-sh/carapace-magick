package argstream

// ParseForCompletion parses a partial magick command argument list and returns
// a CompletionContext describing what is expected at the end.
func ParseForCompletion(args []string, trailingSpace bool) *CompletionContext {
	return ParseForCompletionWithProfile(args, trailingSpace, DefaultMagickProfile)
}

// ParseForCompletionWithProfile parses a partial magick tool argument list using the given profile.
func ParseForCompletionWithProfile(args []string, trailingSpace bool, profile *ToolProfile) *CompletionContext {
	ctx := &CompletionContext{
		Tool: profile.Name,
	}

	i := 0
	var pendingOption *OptionContext
	parenDepth := 0
	toolResolved := false

	// If the profile is the default magick profile, the first positional
	// token could be a sub-tool name. For sub-tool profiles, the tool
	// is already resolved.
	if profile.Name != "magick" {
		toolResolved = true
		ctx.Tool = profile.Name
	}

	for i < len(args) {
		arg := args[i]

		// Handle pending option value
		if pendingOption != nil {
			if !isOption(arg) && arg != "(" && arg != ")" {
				if i == len(args)-1 && !trailingSpace {
					ctx.CurrentOption = pendingOption
					ctx.PartialValue = arg
					if pendingOption.ValueType == ValueDefine {
						ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedDefineValue)
					} else {
						ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOptionValue)
					}
					return ctx
				}
				pendingOption = nil
				i++
				continue
			}
			pendingOption = nil
		}

		// Parentheses
		if arg == "(" {
			if profile.HasParentheses {
				parenDepth++
				ctx.InParentheses = parenDepth > 0
				i++
				continue
			}
		}
		if arg == ")" {
			if profile.HasParentheses && parenDepth > 0 {
				parenDepth--
				ctx.InParentheses = parenDepth > 0
				i++
				continue
			}
		}

		// Check if this is an option
		if isOption(arg) {
			optName, form := parseOptionPrefix(arg)
			optDef := profile.LookupOption(optName)

			// If at the last arg and mid-token, we're completing this option
			if i == len(args)-1 && !trailingSpace {
				ctx.PartialOption = optName
				ctx.CurrentOption = buildOptionContext(optName, form, optDef)

				if optDef != nil && optDef.Type == TypeValue && needsValue(optDef, form) {
					// Complete option name is also a value-taking option —
					// expect both option name and value
					ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOptionName)
					ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOptionValue)
				} else {
					ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOptionName)
				}

				if profile.HasParentheses {
					ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedLParen)
				}

				return ctx
			}

			// Consume option
			i++

			// If value-taking option, mark as pending
			if optDef != nil && optDef.Type == TypeValue && needsValue(optDef, form) {
				pendingOption = buildOptionContext(optName, form, optDef)
			}
		} else {
			// Non-option token
			if !toolResolved && IsKnownToolName(arg) {
				// First positional arg is a sub-tool name
				ctx.Tool = arg
				toolResolved = true
				i++
				continue
			}

			// Image input
			ctx.ImageURLs = append(ctx.ImageURLs, arg)
			i++
		}
	}

	// If we have a pending option value, that's what's expected next
	if pendingOption != nil {
		ctx.CurrentOption = pendingOption
		if pendingOption.ValueType == ValueDefine {
			ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedDefineValue)
		} else {
			ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOptionValue)
		}
		return ctx
	}

	// We've consumed all args — at a new completion position
	if !toolResolved && profile.Name == "magick" {
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedToolName)
	}

	ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOptionName)
	ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedImage)

	if profile.HasOutputArg && len(ctx.ImageURLs) > 0 {
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedOutput)
	}

	if profile.HasParentheses {
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedLParen)
	}
	if ctx.InParentheses {
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedRParen)
	}

	return ctx
}

func buildOptionContext(name string, form OptionForm, optDef *OptionDef) *OptionContext {
	if optDef == nil {
		return &OptionContext{
			Name:       name,
			OptionForm: form,
		}
	}
	return &OptionContext{
		Name:       name,
		ValueType:  optDef.ValueType,
		OptionForm: form,
		Category:   optDef.Category,
		IsBoolean:  optDef.Type == TypeBoolean,
		Style:      optDef.Style(),
	}
}
