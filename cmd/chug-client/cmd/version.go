package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of this client and the connected server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		res := cli.MustGetVersion()

		my := "1.0.0"
		fmt.Printf(`
cli version:      %s
blockchain:       %s
  `, my, res.Blockchain)
	},
}
