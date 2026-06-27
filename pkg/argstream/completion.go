package argstream

type ExpectedToken int

const (
	ExpectedToolName ExpectedToken = iota
	ExpectedOptionName
	ExpectedOptionValue
	ExpectedImage
	ExpectedOutput
	ExpectedDefineValue
	ExpectedLParen
	ExpectedRParen
)

func (t ExpectedToken) String() string {
	switch t {
	case ExpectedToolName:
		return "ToolName"
	case ExpectedOptionName:
		return "OptionName"
	case ExpectedOptionValue:
		return "OptionValue"
	case ExpectedImage:
		return "Image"
	case ExpectedOutput:
		return "Output"
	case ExpectedDefineValue:
		return "DefineValue"
	case ExpectedLParen:
		return "LParen"
	case ExpectedRParen:
		return "RParen"
	}
	return "Unknown"
}

func (t ExpectedToken) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// OptionContext provides details about the option being completed.
type OptionContext struct {
	Name       string         `json:"name"`
	ValueType  ValueType      `json:"valueType"`
	OptionForm OptionForm     `json:"optionForm"`
	Category   OptionCategory `json:"category"`
	IsBoolean  bool           `json:"isBoolean"`
	Style      string         `json:"style"`
}

// CompletionContext describes what is expected at the completion position.
type CompletionContext struct {
	ExpectedTokens []ExpectedToken `json:"expectedTokens"`
	Tool           string          `json:"tool"`
	InParentheses  bool            `json:"inParentheses"`

	CurrentOption *OptionContext `json:"currentOption,omitempty"`

	PartialOption string `json:"partialOption,omitempty"`
	PartialValue  string `json:"partialValue,omitempty"`

	// Collected image URLs from the arg stream
	ImageURLs []string `json:"imageURLs,omitempty"`
}
