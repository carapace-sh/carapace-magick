package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/carapace-sh/carapace-magick/pkg/completer"
	"github.com/spf13/cobra"
)

var mogrifyCmd = &cobra.Command{
	Use:                "mogrify",
	Short:              "ImageMagick in-place image transformer",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func init() {
	profile := argstream.DefaultMogrifyProfile

	carapace.Gen(mogrifyCmd).Standalone()
	rootCmd.AddCommand(mogrifyCmd)

	carapace.Gen(mogrifyCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			args, trailingSpace := completer.ContextToArgs(c)
			ctx := argstream.ParseForCompletionWithProfile(args, trailingSpace, profile)

			if ctx.PartialOption != "" && !trailingSpace {
				return carapace.Batch(
					completer.ActionOptions(ctx, profile),
				).ToA()
			}

			var actions []carapace.Action
			for _, token := range ctx.ExpectedTokens {
				switch token {
				case argstream.ExpectedOptionName:
					actions = append(actions, completer.ActionOptions(ctx, profile))
				case argstream.ExpectedOptionValue, argstream.ExpectedDefineValue:
					actions = append(actions, completer.ActionOptionValue(ctx))
				case argstream.ExpectedImage:
					actions = append(actions, completer.ActionImageInput())
				case argstream.ExpectedOutput:
					actions = append(actions, completer.ActionImageOutput())
				}
			}

			if len(actions) == 0 {
				return carapace.ActionValues()
			}
			return carapace.Batch(actions...).ToA()
		}),
	)
}
