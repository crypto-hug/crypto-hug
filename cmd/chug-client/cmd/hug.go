package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(hugCmd)
}

var hugCmd = &cobra.Command{
	Use:   "hug",
	Short: "manage crypto hugs",
}
