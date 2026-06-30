package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/argstream"
	"github.com/spf13/cobra"
)

var debug_argstreamCmd = &cobra.Command{
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

func init() {
	carapace.Gen(debug_argstreamCmd)
	debugCmd.AddCommand(debug_argstreamCmd)
}
