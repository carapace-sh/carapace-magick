package definevalue

import "strings"

// ParseForCompletion parses a partial -define argument and returns a
// CompletionContext describing what is expected at the cursor position.
func ParseForCompletion(input string) *CompletionContext {
	ctx := &CompletionContext{}

	if input == "" {
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedFormatOrKey)
		return ctx
	}

	// Check if there's a colon (format prefix present)
	colonIdx := strings.Index(input, ":")
	eqIdx := strings.Index(input, "=")

	// If '=' comes before ':', treat the whole prefix as key (no format)
	if eqIdx >= 0 && (colonIdx < 0 || eqIdx < colonIdx) {
		keyPart := input[:eqIdx]
		ctx.Key = keyPart
		ctx.Partial = input[eqIdx+1:]
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedValue)
		return ctx
	}

	// If there's a colon and no equals (or colon before equals)
	if colonIdx >= 0 {
		ctx.Format = input[:colonIdx]

		afterColon := input[colonIdx+1:]
		if afterColon == "" {
			// Cursor is right after "format:" — need key
			ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedKey)
			return ctx
		}

		// Check for '=' in the key part
		beforeEq, afterEq, hasEq := strings.Cut(afterColon, "=")
		if hasEq {
			ctx.Key = beforeEq
			ctx.Partial = afterEq
			ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedValue)
			return ctx
		}

		// Partial key — could be completing key or this is the key with no value yet
		ctx.Key = afterColon
		ctx.Partial = afterColon
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedKey)
		ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedValue)
		return ctx
	}

	// No colon, no equals — could be a global key or start of format prefix
	ctx.Partial = input
	ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedFormatOrKey)
	ctx.ExpectedTokens = append(ctx.ExpectedTokens, ExpectedKey)
	return ctx
}
