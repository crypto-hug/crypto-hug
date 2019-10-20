package cmd

import (
	"strings"

	"github.com/crypto-hug/crypto-hug/cmd/chug-node/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	root = &cobra.Command{
		Use:   "chug",
		Short: "Crypto Hug CLI",
		Long:  ``,
	}

	cli *client.Client = client.NewAPIClient()
)

// Execute executes the root command.
func Execute() error {
	return root.Execute()
}

func init() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		postInitCommands(root.Commands())
	})
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
