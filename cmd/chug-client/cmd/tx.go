package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(txCmd)
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "manage transactions",
}
