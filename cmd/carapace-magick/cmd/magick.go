package cmd

import (
	"slices"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/carapace-sh/carapace-magick/pkg/completer"
	"github.com/spf13/cobra"
)

var magickCmd = &cobra.Command{
	Use:                "magick",
	Short:              "ImageMagick image pipeline processor",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func init() {
	profile := argstream.DefaultMagickProfile

	carapace.Gen(magickCmd).Standalone()
	rootCmd.AddCommand(magickCmd)

	carapace.Gen(magickCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args, trailingSpace := completer.ContextToArgs(c)
			ctx := argstream.ParseForCompletionWithProfile(args, trailingSpace, profile)

			if ctx.PartialOption != "" && !trailingSpace {
				return carapace.Batch(
					completer.ActionOptions(ctx, profile),
					actionOptionValueIfExpected(ctx),
				).ToA()
			}

			var actions []carapace.Action
			for _, token := range ctx.ExpectedTokens {
				switch token {
				case argstream.ExpectedToolName:
					actions = append(actions, completer.ActionToolNames())
				case argstream.ExpectedOptionName:
					actions = append(actions, completer.ActionOptions(ctx, profile))
				case argstream.ExpectedOptionValue, argstream.ExpectedDefineValue:
					actions = append(actions, completer.ActionOptionValue(ctx))
				case argstream.ExpectedImage:
					actions = append(actions, completer.ActionImageInput())
				case argstream.ExpectedOutput:
					actions = append(actions, completer.ActionImageOutput())
				case argstream.ExpectedLParen:
					actions = append(actions, carapace.ActionValues("("))
				case argstream.ExpectedRParen:
					actions = append(actions, carapace.ActionValues(")"))
				}
			}

			if len(actions) == 0 {
				return carapace.ActionValues()
			}
			return carapace.Batch(actions...).ToA()
		}),
	)
}

func actionOptionValueIfExpected(ctx *argstream.CompletionContext) carapace.Action {
	if ctx.CurrentOption != nil && slices.Contains(ctx.ExpectedTokens, argstream.ExpectedOptionValue) {
		return completer.ActionOptionValue(ctx)
	}
	return carapace.ActionValues()
}
