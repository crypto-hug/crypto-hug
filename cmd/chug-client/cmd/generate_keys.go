package cmd

import (
	"fmt"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/spf13/cobra"
	"gopkg.in/workanator/go-ataman.v1"
)


func init() {
	generateCmd.AddCommand(generateKeysCmd)
}

var generateKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Generate a public/private key pair in base 58 format",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		wallet := chug.NewWallet()
		pub := utils.NewBase58JsonValFromData(wallet.PubK)
		priv := utils.NewBase58JsonValFromData(wallet.PrivK)

		rndr := ataman.NewRenderer(ataman.CurlyStyle())
		fmt.Printf(rndr.MustRenderf(`{b+u}generated keys:{-}
 {-+magenta} address:      {intensive_white}%s 
 {-+magenta} public-key:   {intensive_white}%s 
 {-+magenta} private-key:  {intensive_white}%s 
`, wallet.Addr, pub.String(), priv.String()))

	},
}
