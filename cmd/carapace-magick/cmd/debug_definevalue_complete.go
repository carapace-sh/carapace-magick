package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/definevalue"
	"github.com/spf13/cobra"
)

var debug_definevalueCompleteCmd = &cobra.Command{
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
	carapace.Gen(debug_definevalueCompleteCmd)
	debugCmd.AddCommand(debug_definevalueCompleteCmd)
}
