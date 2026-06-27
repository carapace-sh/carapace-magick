package argstream

type ParseError struct {
	Message string
	Span    Span
}

func (e *ParseError) Error() string {
	return e.Message
}

// Parse parses a magick command argument list into a Program AST.
func Parse(args []string) (*Program, error) {
	return ParseWithProfile(args, DefaultMagickProfile)
}

// ParseWithProfile parses a magick tool argument list into a Program AST using the given profile.
func ParseWithProfile(args []string, profile *ToolProfile) (*Program, error) {
	p := &parser{args: args, profile: profile}
	return p.parseProgram()
}

type parser struct {
	args    []string
	pos     int
	profile *ToolProfile
	depth   int
}

func (p *parser) atEnd() bool {
	return p.pos >= len(p.args)
}

func (p *parser) peek() string {
	if p.pos >= len(p.args) {
		return ""
	}
	return p.args[p.pos]
}

func (p *parser) advance() string {
	if p.pos >= len(p.args) {
		return ""
	}
	arg := p.args[p.pos]
	p.pos++
	return arg
}

func (p *parser) parseProgram() (*Program, error) {
	start := p.pos
	prog := &Program{}

	for !p.atEnd() {
		arg := p.peek()

		if arg == "(" {
			p.advance()
			prog.Tokens = append(prog.Tokens, &Token{
				Kind: KindLParen,
				Span: Span{Start: p.pos - 1, End: p.pos},
			})
			p.depth++
			continue
		}

		if arg == ")" {
			p.advance()
			prog.Tokens = append(prog.Tokens, &Token{
				Kind: KindRParen,
				Span: Span{Start: p.pos - 1, End: p.pos},
			})
			p.depth--
			continue
		}

		if isOption(arg) {
			optName, form := parseOptionPrefix(arg)
			optDef := p.profile.LookupOption(optName)

			// Filter options by profile capabilities
			if optDef != nil && !p.profileOptionAllowed(optDef) {
				// Option exists but not in this profile's tool
				p.advance()
				spanEnd := p.pos
				if optDef.Type == TypeValue && needsValue(optDef, form) && !p.atEnd() && !isOption(p.peek()) && p.peek() != "(" && p.peek() != ")" {
					p.advance()
					spanEnd = p.pos
				}
				prog.Tokens = append(prog.Tokens, &Token{
					Kind:       KindOption,
					OptionName: optName,
					OptionForm: form,
					Span:       Span{Start: spanEnd - 1, End: spanEnd},
				})
				continue
			}

			p.advance()
			tok := &Token{
				Kind:       KindOption,
				OptionName: optName,
				OptionForm: form,
				Span:       Span{Start: p.pos - 1, End: p.pos},
			}

			if optDef != nil && optDef.Type == TypeValue && needsValue(optDef, form) {
				if !p.atEnd() && !isOption(p.peek()) && p.peek() != "(" && p.peek() != ")" {
					tok.Value = p.advance()
					tok.Span.End = p.pos
				}
			}
			prog.Tokens = append(prog.Tokens, tok)
		} else {
			// Non-option: image input (or output if last non-option token)
			p.advance()
			prog.Tokens = append(prog.Tokens, &Token{
				Kind: KindImage,
				URL:  arg,
				Span: Span{Start: p.pos - 1, End: p.pos},
			})
		}
	}

	prog.Span = Span{Start: start, End: p.pos}
	return prog, nil
}

// profileOptionAllowed checks if an option category is allowed by the profile.
func (p *parser) profileOptionAllowed(optDef *OptionDef) bool {
	switch optDef.Category {
	case CategoryOperator:
		return p.profile.HasOperators
	case CategoryStackOp:
		return p.profile.HasStackOps
	case CategoryChannelOp, CategorySequenceOp, CategorySetting, CategoryMisc:
		return true
	}
	return true
}

// needsValue returns true if the option in the given form needs a value argument.
// Plus-form resets (PlusReset) take no value. Directional plus-forms may take values.
func needsValue(optDef *OptionDef, form OptionForm) bool {
	if form == FormPlus {
		switch optDef.PlusBehavior {
		case PlusReset:
			return false
		case PlusInverse:
			return false
		case PlusDirectional:
			return optDef.Type == TypeValue
		}
	}
	return optDef.Type == TypeValue
}

// isOption checks if an argument looks like a magick option.
func isOption(arg string) bool {
	if len(arg) == 0 {
		return false
	}
	if arg == "--" {
		return false // end-of-options marker
	}
	if arg[0] == '-' {
		if len(arg) == 1 {
			return false // stdin/stdout
		}
		return true
	}
	if arg[0] == '+' && len(arg) > 1 {
		return true
	}
	return false
}

// parseOptionPrefix splits an option token into its name and form.
func parseOptionPrefix(arg string) (name string, form OptionForm) {
	if len(arg) > 1 && arg[0] == '+' {
		return arg[1:], FormPlus
	}
	if len(arg) > 1 && arg[0] == '-' {
		return arg[1:], FormDash
	}
	return arg, FormDash
}
