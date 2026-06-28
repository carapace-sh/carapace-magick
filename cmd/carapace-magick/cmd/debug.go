package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/carapace-sh/carapace-magick/pkg/definevalue"
	spec "github.com/carapace-sh/carapace-spec"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Parse ImageMagick magick argument streams",
}

func init() {
	carapace.Gen(debugCmd)
	spec.Register(debugCmd)

	debugCmd.AddCommand(argstreamCmd, argstreamCompleteCmd, definevalueCmd, definevalueCompleteCmd)
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

var definevalueCmd = &cobra.Command{
	Use:   "definevalue <value>",
	Short: "Parse a -define format:key=value string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dv, err := definevalue.Parse(args[0])
		if err != nil {
			return err
		}
		m, err := json.Marshal(dv)
		if err != nil {
			return err
		}
		fmt.Println(string(m))
		return nil
	},
}

var definevalueCompleteCmd = &cobra.Command{
	Use:   "definevalue-complete <value>",
	Short: "Get completion context for a -define value string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := definevalue.ParseForCompletion(args[0])
		m, err := json.Marshal(ctx)
		if err != nil {
			return err
		}
		fmt.Println(string(m))
		return nil
	},
}

func init() {
	argstreamCompleteCmd.Flags().Bool("trailing-space", false, "cursor is at a new position after the last arg")
	argstreamCompleteCmd.Flags().String("profile", "magick", "tool profile to use (magick, identify, mogrify, compare, composite, montage)")

	carapace.Gen(argstreamCmd)
	carapace.Gen(argstreamCompleteCmd)
	carapace.Gen(definevalueCmd)
	carapace.Gen(definevalueCompleteCmd)
}
