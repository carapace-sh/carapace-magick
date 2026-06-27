package argstream

// OptionForm indicates whether the option uses the dash or plus prefix.
type OptionForm int

const (
	FormDash OptionForm = iota
	FormPlus
)

func (f OptionForm) String() string {
	switch f {
	case FormDash:
		return "Dash"
	case FormPlus:
		return "Plus"
	}
	return "Unknown"
}

func (f OptionForm) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// TokenKind represents the type of argument token.
type TokenKind int

const (
	KindOption TokenKind = iota
	KindImage
	KindOutput
	KindLParen
	KindRParen
	KindToolName
)

func (k TokenKind) String() string {
	switch k {
	case KindOption:
		return "Option"
	case KindImage:
		return "Image"
	case KindOutput:
		return "Output"
	case KindLParen:
		return "LParen"
	case KindRParen:
		return "RParen"
	case KindToolName:
		return "ToolName"
	}
	return "Unknown"
}

func (k TokenKind) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// Token represents a single parsed argument token.
type Token struct {
	Kind       TokenKind
	OptionName string
	OptionForm OptionForm
	Value      string
	URL        string
	Span       Span
}

// ParenGroup represents a parenthesized sub-pipeline.
type ParenGroup struct {
	Open   Span
	Close  Span
	Tokens []*Token
}

// Program is the top-level AST for a parsed magick command line.
type Program struct {
	Tool   string
	Tokens []*Token
	Groups []*ParenGroup
	Span   Span
}
