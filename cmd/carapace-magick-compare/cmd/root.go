package cmd

import (
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/carapace-sh/carapace-magick/pkg/completer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                "compare",
	Short:              "ImageMagick image comparison",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	profile := argstream.DefaultCompareProfile

	carapace.Gen(rootCmd).Standalone()

	carapace.Gen(rootCmd).PositionalAnyCompletion(
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
				case argstream.ExpectedOptionName, argstream.ExpectedPlusOptionName:
					actions = append(actions, completer.ActionOptions(ctx, profile))
				case argstream.ExpectedOptionValue, argstream.ExpectedDefineValue:
					actions = append(actions, completer.ActionOptionValue(ctx))
				case argstream.ExpectedImage, argstream.ExpectedOutput:
					actions = append(actions, carapace.ActionFiles())
				}
			}

			if len(actions) == 0 {
				return carapace.ActionValues()
			}
			return carapace.Batch(actions...).ToA()
		}),
	)
}
