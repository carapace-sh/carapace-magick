package definevalue

// DefineValue represents a parsed -define argument: [format:]key[=value].
type DefineValue struct {
	Format string `json:"format,omitempty"`
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	Span   Span   `json:"span"`
}

// HasFormat returns true if the define has a format prefix.
func (d *DefineValue) HasFormat() bool {
	return d.Format != ""
}

// HasValue returns true if the define has a value part.
func (d *DefineValue) HasValue() bool {
	return d.Value != ""
}
