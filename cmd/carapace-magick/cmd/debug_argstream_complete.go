package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/spf13/cobra"
)

var debug_argstreamCompleteCmd = &cobra.Command{
	Use:   "argstream-complete <args...>",
	Short: "Get completion context for magick argument stream",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		trailing, _ := cmd.Flags().GetBool("trailing-space")
		profileName, _ := cmd.Flags().GetString("profile")
		profile := argstream.DefaultMagickProfile
		switch profileName {
		case "identify":
			profile = argstream.DefaultIdentifyProfile
		case "mogrify":
			profile = argstream.DefaultMogrifyProfile
		case "compare":
			profile = argstream.DefaultCompareProfile
		case "composite":
			profile = argstream.DefaultCompositeProfile
		case "montage":
			profile = argstream.DefaultMontageProfile
		}
		ctx := argstream.ParseForCompletionWithProfile(args, trailing, profile)
		m, err := json.Marshal(ctx)
		if err != nil {
			return err
		}
		fmt.Println(string(m))
		return nil
	},
}

func init() {
	debug_argstreamCompleteCmd.Flags().Bool("trailing-space", false, "cursor is at a new position after the last arg")
	debug_argstreamCompleteCmd.Flags().String("profile", "magick", "tool profile to use (magick, identify, mogrify, compare, composite, montage)")

	carapace.Gen(debug_argstreamCompleteCmd)
	debugCmd.AddCommand(debug_argstreamCompleteCmd)
}
