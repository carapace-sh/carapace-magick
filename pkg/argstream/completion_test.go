package argstream

import (
	"slices"
	"testing"
)

func assertHasExpected(t *testing.T, ctx *CompletionContext, expected ExpectedToken) {
	t.Helper()
	if !slices.Contains(ctx.ExpectedTokens, expected) {
		t.Errorf("expected token %v not found in %v", expected, ctx.ExpectedTokens)
	}
}

func assertNotHasExpected(t *testing.T, ctx *CompletionContext, expected ExpectedToken) {
	t.Helper()
	if slices.Contains(ctx.ExpectedTokens, expected) {
		t.Errorf("did not expect token %v in %v", expected, ctx.ExpectedTokens)
	}
}

func TestCompletionEmpty(t *testing.T) {
	ctx := ParseForCompletion([]string{}, true)
	assertHasExpected(t, ctx, ExpectedToolName)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
}

func TestCompletionAfterOption(t *testing.T) {
	ctx := ParseForCompletion([]string{"-verbose"}, true)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
}

func TestCompletionOptionValueTrailingSpace(t *testing.T) {
	ctx := ParseForCompletion([]string{"-resize", "200x200"}, true)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
}

func TestCompletionOptionValueMidToken(t *testing.T) {
	ctx := ParseForCompletion([]string{"-resize", "200x"}, false)
	assertHasExpected(t, ctx, ExpectedOptionValue)
	if ctx.CurrentOption == nil {
		t.Fatal("expected CurrentOption")
	}
	if ctx.CurrentOption.Name != "resize" {
		t.Errorf("expected option name 'resize', got %q", ctx.CurrentOption.Name)
	}
	if ctx.PartialValue != "200x" {
		t.Errorf("expected partial value '200x', got %q", ctx.PartialValue)
	}
}

func TestCompletionPendingOptionValue(t *testing.T) {
	ctx := ParseForCompletion([]string{"-resize"}, true)
	assertHasExpected(t, ctx, ExpectedOptionValue)
	if ctx.CurrentOption == nil {
		t.Fatal("expected CurrentOption")
	}
	if ctx.CurrentOption.ValueType != ValueGeometry {
		t.Errorf("expected ValueGeometry, got %v", ctx.CurrentOption.ValueType)
	}
}

func TestCompletionPartialOption(t *testing.T) {
	ctx := ParseForCompletion([]string{"-resi"}, false)
	assertHasExpected(t, ctx, ExpectedOptionName)
	if ctx.PartialOption != "resi" {
		t.Errorf("expected partialOption 'resi', got %q", ctx.PartialOption)
	}
}

func TestCompletionPlusOption(t *testing.T) {
	ctx := ParseForCompletion([]string{"+verb"}, false)
	assertHasExpected(t, ctx, ExpectedPlusOptionName)
}

func TestCompletionToolName(t *testing.T) {
	ctx := ParseForCompletion([]string{"iden"}, false)
	assertHasExpected(t, ctx, ExpectedToolName)
}

func TestCompletionToolResolved(t *testing.T) {
	ctx := ParseForCompletion([]string{"identify"}, true)
	assertNotHasExpected(t, ctx, ExpectedToolName)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
	if ctx.Tool != "identify" {
		t.Errorf("expected tool 'identify', got %q", ctx.Tool)
	}
}

func TestCompletionImageInput(t *testing.T) {
	ctx := ParseForCompletion([]string{"input.png"}, true)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
	if len(ctx.ImageURLs) != 1 || ctx.ImageURLs[0] != "input.png" {
		t.Errorf("expected ImageURLs=['input.png'], got %v", ctx.ImageURLs)
	}
}

func TestCompletionDefineValue(t *testing.T) {
	ctx := ParseForCompletion([]string{"-define"}, true)
	assertHasExpected(t, ctx, ExpectedDefineValue)
}

func TestCompletionDefineValueMidToken(t *testing.T) {
	ctx := ParseForCompletion([]string{"-define", "jpeg:"}, false)
	assertHasExpected(t, ctx, ExpectedDefineValue)
	if ctx.CurrentOption == nil {
		t.Fatal("expected CurrentOption")
	}
}

func TestCompletionIdentifyProfile(t *testing.T) {
	ctx := ParseForCompletionWithProfile([]string{}, true, DefaultIdentifyProfile)
	assertNotHasExpected(t, ctx, ExpectedToolName)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
}

func TestCompletionMogrifyProfile(t *testing.T) {
	ctx := ParseForCompletionWithProfile([]string{}, true, DefaultMogrifyProfile)
	assertNotHasExpected(t, ctx, ExpectedToolName)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedImage)
}

func TestCompletionIdentifyNoOperators(t *testing.T) {
	// identify profile should not have -resize
	opt := DefaultIdentifyProfile.LookupOption("resize")
	if opt != nil {
		t.Error("identify profile should not have -resize option")
	}
}

func TestCompletionParentheses(t *testing.T) {
	ctx := ParseForCompletion([]string{}, true)
	assertHasExpected(t, ctx, ExpectedLParen)
}

func TestCompletionInParentheses(t *testing.T) {
	ctx := ParseForCompletion([]string{"("}, true)
	if !ctx.InParentheses {
		t.Error("expected InParentheses to be true")
	}
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertHasExpected(t, ctx, ExpectedRParen)
}

func TestCompletionPlusResetNoValueExpected(t *testing.T) {
	// +background should not expect a value (it's a reset)
	ctx := ParseForCompletion([]string{"+background"}, true)
	assertHasExpected(t, ctx, ExpectedOptionName)
	assertNotHasExpected(t, ctx, ExpectedOptionValue)
}
