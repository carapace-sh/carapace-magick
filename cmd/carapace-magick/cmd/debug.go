package cmd

import (
	"github.com/carapace-sh/carapace"
	spec "github.com/carapace-sh/carapace-spec"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Parse ImageMagick magick argument streams",
}

func init() {
	carapace.Gen(debugCmd)
	rootCmd.AddCommand(debugCmd)

	spec.Register(debugCmd)
}
