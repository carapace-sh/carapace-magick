package completer

import (
	"net/url"

	"github.com/carapace-sh/carapace"
	magick "github.com/carapace-sh/carapace-magick/pkg/actions/tools/magick"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/carapace-sh/carapace-magick/pkg/definevalue"
	"github.com/carapace-sh/carapace/pkg/uid"
)

// ContextToArgs converts carapace.Context to the args and trailingSpace
// expected by argstream.ParseForCompletion.
func ContextToArgs(c carapace.Context) (args []string, trailingSpace bool) {
	n := len(c.Args)
	if n > 0 && c.Args[n-1] == "" {
		n--
	}
	args = c.Args[:n]
	if c.Value != "" {
		args = append(args, c.Value)
	}
	trailingSpace = c.Value == ""
	return
}

// ActionOptions returns completions for option names appropriate to the current context.
func ActionOptions(ctx *argstream.CompletionContext, profile *argstream.ToolProfile) carapace.Action {
	var vals []string
	for name, def := range profile.OptionIndex {
		vals = append(vals, "-"+name, def.Description, def.Style())
		if def.HasPlusForm {
			vals = append(vals, "+"+name, plusFormDescription(def), styleForPlusForm(def))
		}
	}
	return carapace.ActionStyledValuesDescribed(vals...).UidF(optionUid(profile))
}

// ActionToolNames returns completions for sub-tool names.
func ActionToolNames() carapace.Action {
	return magick.ActionToolNames()
}

// ActionOptionValue returns completions for the value of the current option.
func ActionOptionValue(ctx *argstream.CompletionContext) carapace.Action {
	if ctx.CurrentOption == nil {
		return carapace.ActionValues()
	}
	switch ctx.CurrentOption.ValueType {
	case argstream.ValueGeometry:
		return carapace.ActionValues() // geometry is free-form
	case argstream.ValueColor:
		return magick.ActionColors()
	case argstream.ValueColorspace:
		return magick.ActionColorspaces()
	case argstream.ValueCompose:
		return magick.ActionComposes()
	case argstream.ValueCompress:
		return magick.ActionCompressTypes()
	case argstream.ValueChannel:
		return magick.ActionChannels()
	case argstream.ValueDistort:
		return magick.ActionDistortMethods()
	case argstream.ValueFilter:
		return magick.ActionFilters()
	case argstream.ValueGravity:
		return magick.ActionGravities()
	case argstream.ValueInterlace:
		return magick.ActionInterlaceTypes()
	case argstream.ValueLayers:
		return magick.ActionLayerMethods()
	case argstream.ValueMorphology:
		return magick.ActionMorphologyMethods()
	case argstream.ValueImageType:
		return magick.ActionTypes()
	case argstream.ValueVirtualPixel:
		return magick.ActionVirtualPixelMethods()
	case argstream.ValueMetric:
		return magick.ActionMetrics()
	case argstream.ValueEvaluate:
		return magick.ActionEvaluateOps()
	case argstream.ValueDefine:
		return ActionDefineValue(ctx.PartialValue)
	case argstream.ValueBoolean:
		return magick.ActionBoolean()
	case argstream.ValueFont:
		return magick.ActionFonts()
	case argstream.ValueFormat:
		return magick.ActionFormats()
	case argstream.ValueFilename:
		return carapace.ActionFiles()
	case argstream.ValueOrientation:
		return magick.ActionOrientations()
	case argstream.ValueDispose:
		return magick.ActionDisposes()
	case argstream.ValueAlpha:
		return magick.ActionAlphaOption()
	case argstream.ValueNoise:
		return magick.ActionNoiseTypes()
	case argstream.ValuePreview:
		return magick.ActionPreviewTypes()
	case argstream.ValueStorage:
		return magick.ActionStorageTypes()
	case argstream.ValueKernel:
		return magick.ActionKernels()
	case argstream.ValueList:
		return magick.ActionListTypes()
	case argstream.ValueAutoThreshold:
		return magick.ActionAutoThreshold()
	case argstream.ValueDirection:
		return magick.ActionDirection()
	case argstream.ValueDegrees, argstream.ValueInt, argstream.ValueFloat, argstream.ValueRatio, argstream.ValueThreshold, argstream.ValueExpression:
		return carapace.ActionValues()
	case argstream.ValueString, argstream.ValueMethod, argstream.ValueDensity, argstream.ValuePoint, argstream.ValueBlend, argstream.ValueCount, argstream.ValueIndex:
		return carapace.ActionValues()
	default:
		return carapace.ActionValues()
	}
}

func plusFormDescription(def *argstream.OptionDef) string {
	switch def.PlusBehavior {
	case argstream.PlusReset:
		return "reset " + def.Name + " to default"
	case argstream.PlusInverse:
		return "inverse of " + def.Name
	case argstream.PlusDirectional:
		return "alternate direction for " + def.Name
	}
	return def.Description
}

func styleForPlusForm(_ *argstream.OptionDef) string {
	return "dim" // dim style for plus forms
}

func optionUid(_ *argstream.ToolProfile) func(s string, uc uid.Context) (*url.URL, error) {
	return func(s string, uc uid.Context) (*url.URL, error) {
		prefix := ""
		if len(s) > 1 && s[0] == '+' {
			prefix = "+"
		}
		name := s
		if len(name) > 1 && (name[0] == '-' || name[0] == '+') {
			name = name[1:]
		}
		return &url.URL{
			Scheme: "magick",
			Host:   "option",
			Path:   prefix + name,
		}, nil
	}
}

// ActionDefineValue returns completions for a -define argument value.
// It parses the partial input to determine whether to complete format prefixes,
// define keys, or define values.
func ActionDefineValue(partial string) carapace.Action {
	ctx := definevalue.ParseForCompletion(partial)

	for _, expected := range ctx.ExpectedTokens {
		switch expected {
		case definevalue.ExpectedFormatOrKey:
			return carapace.Batch(
				ActionDefineFormats(),
				ActionDefineGlobalKeys(),
				ActionDefineFormatPrefixes(),
			).ToA()

		case definevalue.ExpectedFormat:
			return ActionDefineFormatPrefixes()

		case definevalue.ExpectedKey:
			if ctx.Format != "" {
				return ActionDefineKeys(ctx.Format)
			}
			return carapace.Batch(
				ActionDefineGlobalKeys(),
				ActionDefineFormatPrefixes(),
			).ToA()

		case definevalue.ExpectedValue:
			if ctx.Format != "" && ctx.Key != "" {
				return ActionDefineValues(ctx.Format, ctx.Key)
			}
			return carapace.ActionValues()
		}
	}

	return carapace.ActionValues()
}

// ActionDefineFormatPrefixes returns completions for format: prefixes in -define values.
func ActionDefineFormatPrefixes() carapace.Action {
	return magick.ActionFormats().Suffix(":")
}

// ActionDefineFormats returns format names for -define completion (without colon).
func ActionDefineFormats() carapace.Action {
	return magick.ActionFormats()
}

// ActionDefineGlobalKeys returns completions for global define keys.
func ActionDefineGlobalKeys() carapace.Action {
	keys := definevalue.LookupGlobalDefines()
	vals := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		vals = append(vals, k.Name, k.Description)
	}
	return carapace.ActionValuesDescribed(vals...).Tag("define keys").Uid("magick", "define-keys", "global")
}

// ActionDefineKeys returns completions for format-specific define keys.
func ActionDefineKeys(format string) carapace.Action {
	keys := definevalue.LookupFormatDefines(format)
	if len(keys) == 0 {
		return carapace.ActionValues()
	}
	vals := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		vals = append(vals, k.Name, k.Description)
	}
	return carapace.ActionValuesDescribed(vals...).Tag("define keys").Uid("magick", "define-keys", format)
}

// ActionDefineValues returns completions for a specific define key's values.
func ActionDefineValues(format, key string) carapace.Action {
	defKey := definevalue.LookupDefineKey(format, key)
	if defKey == nil || len(defKey.Values) == 0 {
		switch defKey.ValueType {
		case "boolean":
			return magick.ActionBoolean()
		case "int", "float":
			return carapace.ActionValues()
		default:
			return carapace.ActionValues()
		}
	}
	return carapace.ActionValues(defKey.Values...).Tag("define values").Uid("magick", "define-values", format, key)
}
