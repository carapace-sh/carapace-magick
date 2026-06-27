package definevalue

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

func TestCompletionEmpty(t *testing.T) {
	ctx := ParseForCompletion("")
	assertHasExpected(t, ctx, ExpectedFormatOrKey)
}

func TestCompletionFormatPrefix(t *testing.T) {
	ctx := ParseForCompletion("jpeg:")
	assertHasExpected(t, ctx, ExpectedKey)
	if ctx.Format != "jpeg" {
		t.Errorf("expected format 'jpeg', got %q", ctx.Format)
	}
}

func TestCompletionFormatKeyPartial(t *testing.T) {
	ctx := ParseForCompletion("jpeg:qual")
	assertHasExpected(t, ctx, ExpectedKey)
	assertHasExpected(t, ctx, ExpectedValue)
	if ctx.Format != "jpeg" {
		t.Errorf("expected format 'jpeg', got %q", ctx.Format)
	}
}

func TestCompletionFormatKeyValuePartial(t *testing.T) {
	ctx := ParseForCompletion("jpeg:quality=8")
	assertHasExpected(t, ctx, ExpectedValue)
	if ctx.Format != "jpeg" {
		t.Errorf("expected format 'jpeg', got %q", ctx.Format)
	}
	if ctx.Key != "quality" {
		t.Errorf("expected key 'quality', got %q", ctx.Key)
	}
	if ctx.Partial != "8" {
		t.Errorf("expected partial '8', got %q", ctx.Partial)
	}
}

func TestCompletionGlobalKeyPartial(t *testing.T) {
	ctx := ParseForCompletion("optimize=tr")
	assertHasExpected(t, ctx, ExpectedValue)
	if ctx.Key != "optimize" {
		t.Errorf("expected key 'optimize', got %q", ctx.Key)
	}
}

func TestCompletionNoFormatNoEquals(t *testing.T) {
	ctx := ParseForCompletion("jpeg")
	assertHasExpected(t, ctx, ExpectedFormatOrKey)
	assertHasExpected(t, ctx, ExpectedKey)
}

func TestCompletionFormatKeyEqualsNoValue(t *testing.T) {
	ctx := ParseForCompletion("png:compression-level=")
	assertHasExpected(t, ctx, ExpectedValue)
	if ctx.Format != "png" {
		t.Errorf("expected format 'png', got %q", ctx.Format)
	}
	if ctx.Key != "compression-level" {
		t.Errorf("expected key 'compression-level', got %q", ctx.Key)
	}
}
