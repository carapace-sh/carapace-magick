package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-magick/pkg/definevalue"
	"github.com/spf13/cobra"
)

var debug_definevalueCmd = &cobra.Command{
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

func init() {
	carapace.Gen(debug_definevalueCmd)
	debugCmd.AddCommand(debug_definevalueCmd)
}
