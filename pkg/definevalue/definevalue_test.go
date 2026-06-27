package definevalue

import "testing"

func TestParseFormatKeyValue(t *testing.T) {
	dv, err := Parse("jpeg:quality=85")
	if err != nil {
		t.Fatal(err)
	}
	if dv.Format != "jpeg" {
		t.Errorf("expected format 'jpeg', got %q", dv.Format)
	}
	if dv.Key != "quality" {
		t.Errorf("expected key 'quality', got %q", dv.Key)
	}
	if dv.Value != "85" {
		t.Errorf("expected value '85', got %q", dv.Value)
	}
}

func TestParseKeyValueNoFormat(t *testing.T) {
	dv, err := Parse("optimize=true")
	if err != nil {
		t.Fatal(err)
	}
	if dv.Format != "" {
		t.Errorf("expected empty format, got %q", dv.Format)
	}
	if dv.Key != "optimize" {
		t.Errorf("expected key 'optimize', got %q", dv.Key)
	}
	if dv.Value != "true" {
		t.Errorf("expected value 'true', got %q", dv.Value)
	}
}

func TestParseFormatKeyNoValue(t *testing.T) {
	dv, err := Parse("png:compression-level")
	if err != nil {
		t.Fatal(err)
	}
	if dv.Format != "png" {
		t.Errorf("expected format 'png', got %q", dv.Format)
	}
	if dv.Key != "compression-level" {
		t.Errorf("expected key 'compression-level', got %q", dv.Key)
	}
	if dv.Value != "" {
		t.Errorf("expected empty value, got %q", dv.Value)
	}
}

func TestParseKeyNoFormatNoValue(t *testing.T) {
	dv, err := Parse("type")
	if err != nil {
		t.Fatal(err)
	}
	if dv.Format != "" {
		t.Errorf("expected empty format, got %q", dv.Format)
	}
	if dv.Key != "type" {
		t.Errorf("expected key 'type', got %q", dv.Key)
	}
	if dv.Value != "" {
		t.Errorf("expected empty value, got %q", dv.Value)
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParseValueWithCommas(t *testing.T) {
	dv, err := Parse("png:exclude-chunk=tEXt,zTXt")
	if err != nil {
		t.Fatal(err)
	}
	if dv.Format != "png" {
		t.Errorf("expected format 'png', got %q", dv.Format)
	}
	if dv.Key != "exclude-chunk" {
		t.Errorf("expected key 'exclude-chunk', got %q", dv.Key)
	}
	if dv.Value != "tEXt,zTXt" {
		t.Errorf("expected value 'tEXt,zTXt', got %q", dv.Value)
	}
}

func TestDefineValueHasFormat(t *testing.T) {
	dv := &DefineValue{Format: "jpeg", Key: "quality"}
	if !dv.HasFormat() {
		t.Error("expected HasFormat=true")
	}
	dv2 := &DefineValue{Key: "optimize"}
	if dv2.HasFormat() {
		t.Error("expected HasFormat=false")
	}
}

func TestDefineValueHasValue(t *testing.T) {
	dv := &DefineValue{Format: "jpeg", Key: "quality", Value: "85"}
	if !dv.HasValue() {
		t.Error("expected HasValue=true")
	}
	dv2 := &DefineValue{Format: "jpeg", Key: "quality"}
	if dv2.HasValue() {
		t.Error("expected HasValue=false")
	}
}
