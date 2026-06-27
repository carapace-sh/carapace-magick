package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	spec "github.com/carapace-sh/carapace-spec"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-magick-debug",
	Short: "Parse ImageMagick magick argument streams",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	carapace.Gen(rootCmd)
	spec.Register(rootCmd)
}

var argstreamCmd = &cobra.Command{
	Use:   "argstream <args...>",
	Short: "Parse magick argument stream",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		prog, err := argstream.Parse(args)
		if err != nil {
			return err
		}
		m, err := json.Marshal(prog)
		if err != nil {
			return err
		}
		fmt.Println(string(m))
		return nil
	},
	DisableFlagParsing: true,
}

var argstreamCompleteCmd = &cobra.Command{
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
	rootCmd.AddCommand(argstreamCmd)
	rootCmd.AddCommand(argstreamCompleteCmd)

	argstreamCompleteCmd.Flags().Bool("trailing-space", false, "cursor is at a new position after the last arg")
	argstreamCompleteCmd.Flags().String("profile", "magick", "tool profile to use (magick, identify, mogrify, compare, composite, montage)")

	carapace.Gen(argstreamCmd)
	carapace.Gen(argstreamCompleteCmd)
}
