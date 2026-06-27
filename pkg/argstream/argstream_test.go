package argstream

import (
	"testing"
)

func TestParseEmpty(t *testing.T) {
	prog, err := Parse([]string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(prog.Tokens))
	}
}

func TestParseBooleanOption(t *testing.T) {
	prog, err := Parse([]string{"-verbose"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].Kind != KindOption {
		t.Errorf("expected KindOption, got %v", prog.Tokens[0].Kind)
	}
	if prog.Tokens[0].OptionName != "verbose" {
		t.Errorf("expected option 'verbose', got %q", prog.Tokens[0].OptionName)
	}
	if prog.Tokens[0].OptionForm != FormDash {
		t.Errorf("expected FormDash, got %v", prog.Tokens[0].OptionForm)
	}
}

func TestParseValueOption(t *testing.T) {
	prog, err := Parse([]string{"-resize", "200x200"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].OptionName != "resize" {
		t.Errorf("expected option 'resize', got %q", prog.Tokens[0].OptionName)
	}
	if prog.Tokens[0].Value != "200x200" {
		t.Errorf("expected value '200x200', got %q", prog.Tokens[0].Value)
	}
}

func TestParsePlusForm(t *testing.T) {
	prog, err := Parse([]string{"+verbose"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].OptionForm != FormPlus {
		t.Errorf("expected FormPlus, got %v", prog.Tokens[0].OptionForm)
	}
	if prog.Tokens[0].OptionName != "verbose" {
		t.Errorf("expected option 'verbose', got %q", prog.Tokens[0].OptionName)
	}
}

func TestParsePlusResetNoValue(t *testing.T) {
	prog, err := Parse([]string{"+background", "input.png"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].OptionName != "background" {
		t.Errorf("expected option 'background', got %q", prog.Tokens[0].OptionName)
	}
	if prog.Tokens[0].Value != "" {
		t.Errorf("expected no value for +background (reset), got %q", prog.Tokens[0].Value)
	}
}

func TestParseImageInput(t *testing.T) {
	prog, err := Parse([]string{"input.png"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].Kind != KindImage {
		t.Errorf("expected KindImage, got %v", prog.Tokens[0].Kind)
	}
	if prog.Tokens[0].URL != "input.png" {
		t.Errorf("expected URL 'input.png', got %q", prog.Tokens[0].URL)
	}
}

func TestParseFullCommand(t *testing.T) {
	prog, err := Parse([]string{"-resize", "200x200", "input.png", "output.png"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 3 {
		t.Fatalf("expected 3 tokens (resize+value counts as 1), got %d", len(prog.Tokens))
	}
}

func TestParseParentheses(t *testing.T) {
	prog, err := Parse([]string{"(", "-resize", "50%", ")", "input.png", "output.png"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 5 {
		t.Fatalf("expected 5 tokens (resize+value is 1), got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].Kind != KindLParen {
		t.Errorf("expected KindLParen at 0, got %v", prog.Tokens[0].Kind)
	}
	if prog.Tokens[2].Kind != KindRParen {
		t.Errorf("expected KindRParen at 2, got %v", prog.Tokens[2].Kind)
	}
}

func TestParseStackOperators(t *testing.T) {
	prog, err := Parse([]string{"input.png", "-clone", "0", "-negate", "-delete", "0", "output.png"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 5 {
		t.Fatalf("expected 5 tokens (clone+value=1, delete+value=1), got %d", len(prog.Tokens))
	}
	if prog.Tokens[1].OptionName != "clone" {
		t.Errorf("expected option 'clone', got %q", prog.Tokens[1].OptionName)
	}
	if prog.Tokens[1].Value != "0" {
		t.Errorf("expected value '0', got %q", prog.Tokens[1].Value)
	}
}

func TestParseDefine(t *testing.T) {
	prog, err := Parse([]string{"-define", "jpeg:quality=85", "input.jpg", "output.jpg"})
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(prog.Tokens))
	}
	if prog.Tokens[0].Value != "jpeg:quality=85" {
		t.Errorf("expected value 'jpeg:quality=85', got %q", prog.Tokens[0].Value)
	}
}

func TestParseIdentifyWithProfile(t *testing.T) {
	prog, err := ParseWithProfile([]string{"-verbose", "input.png"}, DefaultIdentifyProfile)
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(prog.Tokens))
	}
}

func TestLookupOption(t *testing.T) {
	opt := LookupOption("resize")
	if opt == nil {
		t.Fatal("expected option 'resize' to exist")
	}
	if opt.Type != TypeValue {
		t.Errorf("expected TypeValue for 'resize', got %v", opt.Type)
	}
	if opt.ValueType != ValueGeometry {
		t.Errorf("expected ValueGeometry for 'resize', got %v", opt.ValueType)
	}
}

func TestIsKnownToolName(t *testing.T) {
	if !IsKnownToolName("identify") {
		t.Error("expected 'identify' to be a known tool name")
	}
	if !IsKnownToolName("mogrify") {
		t.Error("expected 'mogrify' to be a known tool name")
	}
	if IsKnownToolName("resize") {
		t.Error("expected 'resize' to not be a known tool name")
	}
}
