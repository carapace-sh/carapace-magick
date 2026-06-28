package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/cmd/carapace-magick/cmd/snippet"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/carapace-sh/carapace-magick/pkg/completer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "carapace-magick",
	Short:             "ImageMagick completion provider",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

var carapaceMagickCmd = &cobra.Command{
	Use:                "carapace-magick",
	Short:              "ImageMagick completion provider",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func Execute() {
	if len(os.Args) > 1 && os.Args[1] == "_carapace" {
		if len(os.Args) < 4 {
			shell := ""
			if len(os.Args) > 2 {
				shell = os.Args[2]
			}
			fmt.Println(snippet.Snippet(shell))
			return
		}
		// Route completion/export requests to the correct subcommand.
		// bridge.ActionCarapace("carapace-magick", "identify") calls:
		//   carapace-magick _carapace export "" identify -verbose image.png
		// Rewrite to:
		//   carapace-magick identify _carapace export "" -verbose image.png
		subcommand := "magick"
		if len(os.Args) > 4 && isCompleterSubcommand(os.Args[4]) {
			subcommand = os.Args[4]
			os.Args = append(
				[]string{os.Args[0], subcommand, "_carapace", os.Args[2], os.Args[3]},
				os.Args[5:]...,
			)
		} else {
			os.Args = append(
				[]string{os.Args[0], subcommand, "_carapace"},
				os.Args[2:]...,
			)
		}
	}
	if len(os.Args) > 2 && isCompleterSubcommand(os.Args[1]) && os.Args[2] == "_carapace" {
		// Subcommand-level snippet request: carapace-magick <subcommand> _carapace [shell]
		if len(os.Args) < 5 {
			shell := ""
			if len(os.Args) > 3 {
				shell = os.Args[3]
			}
			fmt.Println(snippet.SingleSnippet(shell, os.Args[1]))
			return
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isCompleterSubcommand(name string) bool {
	return slices.Contains([]string{"carapace-magick", "magick", "identify", "mogrify", "compare", "composite", "montage"}, name)
}

func init() {
	rootCmd.AddCommand(
		carapaceMagickCmd,
		magickCmd,
		identifyCmd,
		mogrifyCmd,
		compareCmd,
		compositeCmd,
		montageCmd,
		debugCmd,
	)

	carapace.Gen(carapaceMagickCmd).Standalone()

	carapace.Gen(carapaceMagickCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionValues("magick", "identify", "mogrify", "compare", "composite", "montage", "debug")
		}),
	)
}

var magickCmd = &cobra.Command{
	Use:                "magick",
	Short:              "ImageMagick image pipeline processor",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func init() {
	profile := argstream.DefaultMagickProfile

	carapace.Gen(magickCmd).Standalone()

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
