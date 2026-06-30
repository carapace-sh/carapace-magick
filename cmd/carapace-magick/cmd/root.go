package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "carapace-magick",
	Short:             "ImageMagick completion provider",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	carapace.Gen(rootCmd).Execute()
}

func init() {
	carapace.Gen(rootCmd,
		carapace.WithSubcommands(magickCmd, identifyCmd, mogrifyCmd, compareCmd, compositeCmd, montageCmd),
	).Standalone()
}
