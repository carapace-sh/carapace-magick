package definevalue

// ExpectedToken indicates what kind of token is expected at the completion position.
type ExpectedToken int

const (
	ExpectedFormat ExpectedToken = iota
	ExpectedKey
	ExpectedValue
	ExpectedFormatOrKey
)

func (t ExpectedToken) String() string {
	switch t {
	case ExpectedFormat:
		return "Format"
	case ExpectedKey:
		return "Key"
	case ExpectedValue:
		return "Value"
	case ExpectedFormatOrKey:
		return "FormatOrKey"
	}
	return "Unknown"
}

func (t ExpectedToken) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// CompletionContext describes what is expected at the completion position
// within a -define argument string.
type CompletionContext struct {
	ExpectedTokens []ExpectedToken `json:"expectedTokens"`
	Format         string          `json:"format,omitempty"`
	Key            string          `json:"key,omitempty"`
	Partial        string          `json:"partial,omitempty"`
}
