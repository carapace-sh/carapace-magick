package definevalue

type ParseError struct {
	Message string
	Span    Span
}

func (e *ParseError) Error() string {
	return e.Message
}

// Parse parses a -define argument string into a DefineValue.
// Format: [format:]key[=value]
func Parse(input string) (*DefineValue, error) {
	if input == "" {
		return nil, &ParseError{Message: "empty define value", Span: Span{Start: 0, End: 0}}
	}

	dv := &DefineValue{
		Span: Span{Start: 0, End: len(input)},
	}

	// Split at the first '=' to separate key part from value.
	eqIdx := indexOfEquals(input)
	if eqIdx >= 0 {
		dv.Value = input[eqIdx+1:]
		keyPart := input[:eqIdx]
		dv.Span.End = len(input)
		splitKeyPart(dv, keyPart)
	} else {
		// No '=' — just [format:]key
		splitKeyPart(dv, input)
	}

	return dv, nil
}

// indexOfEquals finds the first '=' that is not inside the format prefix.
// The format prefix ends at the first ':', so we skip '=' before the first ':'.
func indexOfEquals(input string) int {
	colonIdx := -1
	for i, c := range input {
		if c == ':' && colonIdx == -1 {
			colonIdx = i
		}
		if c == '=' {
			// If we haven't seen a colon yet, this '=' is the value separator
			if colonIdx == -1 || i > colonIdx {
				return i
			}
		}
	}
	return -1
}

// splitKeyPart splits [format:]key into dv.Format and dv.Key.
func splitKeyPart(dv *DefineValue, keyPart string) {
	colonIdx := -1
	for i, c := range keyPart {
		if c == ':' {
			colonIdx = i
			break
		}
	}
	if colonIdx >= 0 {
		dv.Format = keyPart[:colonIdx]
		dv.Key = keyPart[colonIdx+1:]
	} else {
		dv.Key = keyPart
	}
}
