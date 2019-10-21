package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	hugCmd.AddCommand(hugEtagCmd)
}

var hugEtagCmd = &cobra.Command{
	Use:   "etag",
	Short: "get the current etag of the given crypto hug",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			panic("⚠️  no address provided")
		}

		etag, err := cli.GetHugEtag(args[0])
		if err != nil {
			panic(fmt.Sprintf("🚨 %s\n", err.Error()))
		}

		fmt.Println(fmt.Sprintf("✴️  etag: %s", etag))
	},
}
