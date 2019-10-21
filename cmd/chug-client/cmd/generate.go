package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(generateCmd)
}


var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generation utilities",
	Long:  ``,
}
