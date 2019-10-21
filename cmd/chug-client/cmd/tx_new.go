package cmd

import (
	"fmt"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/cmd/chug-client/print"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/v-braun/go-must"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/workanator/go-ataman.v1"
)

// example keys:
// genesis key
// address:      oH59mogZ1nFbWa2AuvnELcbUg4tVaYyrN
// public-key:   4AbRc2eSQetaPgzDqrknhyuyP6v1yrLHNkXwMsoGQkXqcFwArtNCxY869baaq5aYhpAFzQLxfLBfF7KWN4kW1BsK
// private-key:  3FMrspUmNTh6n2ozmjRPpT8TMsq1dNcA6Nmdzqr4mDdd
//
//
// another example key
// address:      QZLPwyodQMRFkzWeX5y4hETPrkaNgRcWS
// public-key:   5cVgdi1iJV346ke6Uae9aJbwtX5wNTWczfxQf6fShXsMf3Dz4h29ifzboWC3RKLAJuqJkpXu2HbfzjV7XgKziCTB
// private-key:  7BRMCLDE6mUvoi9iqv5Nh6YQBhJxLyRekvQVL3Q9QGLM

/*
example commands:

generate genesis tx and print it:

./bin/chug tx new \
--genesis \
--issuer-pub=4AbRc2eSQetaPgzDqrknhyuyP6v1yrLHNkXwMsoGQkXqcFwArtNCxY869baaq5aYhpAFzQLxfLBfF7KWN4kW1BsK \
--issuer-priv=3FMrspUmNTh6n2ozmjRPpT8TMsq1dNcA6Nmdzqr4mDdd \
--data="hug the planet" \
--print=true


generate spawn hug tx:
./bin/chug tx new \
--issuer-pub=4AbRc2eSQetaPgzDqrknhyuyP6v1yrLHNkXwMsoGQkXqcFwArtNCxY869baaq5aYhpAFzQLxfLBfF7KWN4kW1BsK \
--issuer-priv=3FMrspUmNTh6n2ozmjRPpT8TMsq1dNcA6Nmdzqr4mDdd \
--validator-pub=5cVgdi1iJV346ke6Uae9aJbwtX5wNTWczfxQf6fShXsMf3Dz4h29ifzboWC3RKLAJuqJkpXu2HbfzjV7XgKziCTB \
--validator-priv=7BRMCLDE6mUvoi9iqv5Nh6YQBhJxLyRekvQVL3Q9QGLM \
--data="hug the universe" \
--print=true \
--send \
--query-etag


*/

func askCreateTxQuestions(isGenesis bool) *chug.Transaction {
	rndr := ataman.NewRenderer(ataman.CurlyStyle())

	var tx *chug.Transaction
	if isGenesis {
		tx = chug.NewTransaction(chug.SpawnGenesisHugTransactionType)
	} else {
		tx = chug.NewTransaction(chug.GiveHugTransactionType)
	}

	answers := struct {
		IssuerEtag    string
		ValidatorEtag string
		Data          string

		IssuerPub     string
		IssuerPriv    string
		ValidatorPub  string
		ValidatorPriv string
	}{
		IssuerEtag:    viper.GetString("issuer-etag"),
		ValidatorEtag: viper.GetString("validator-etag"),
		Data:          viper.GetString("data"),
		IssuerPub:     viper.GetString("issuer-pub"),
		IssuerPriv:    viper.GetString("issuer-priv"),
		ValidatorPub:  viper.GetString("validator-pub"),
		ValidatorPriv: viper.GetString("validator-priv"),
	}

	questions := []*survey.Question{}
	if !isGenesis && !viper.GetBool("query-etag") && viper.GetString("issuer-etag") == "" {
		questions = append(questions, &survey.Question{
			Name:     "issuerEtag",
			Prompt:   &survey.Input{Message: "enter your etag"},
			Validate: survey.Required,
		})
	}
	if !isGenesis && !viper.GetBool("query-etag") && viper.GetString("validator-etag") == "" {
		questions = append(questions, &survey.Question{
			Name:   "validatorEtag",
			Prompt: &survey.Input{Message: "enter peer etag"},
		})
	}
	if viper.GetString("data") == "" {
		questions = append(questions, &survey.Question{
			Name:   "data",
			Prompt: &survey.Multiline{Message: "enter data to store in the transaction"},
		})
	}

	survey.Ask(questions, &answers)
	fmt.Printf(rndr.MustRenderf("✅ {-+light+cyan}hash transaction: {-+white}{-}\n"))

	tx.IssuerEtag = answers.IssuerEtag
	tx.ValidatorEtag = answers.ValidatorEtag
	if answers.Data != "" {
		tx.Data = utils.NewBase58JsonValFromData([]byte(answers.Data))
	}

	questions = []*survey.Question{}
	if viper.GetString("issuer-pub") == "" {
		questions = append(questions, &survey.Question{
			Name:   "issuerPub",
			Prompt: &survey.Input{Message: "enter your pub. key"},
		})
	}
	if viper.GetString("issuer-priv") == "" {
		questions = append(questions, &survey.Question{
			Name:   "issuerPriv",
			Prompt: &survey.Input{Message: "enter your priv. key"},
		})
	}
	survey.Ask(questions, &answers)

	questions = []*survey.Question{}
	if isGenesis {
		answers.ValidatorPub = answers.IssuerPub
		answers.ValidatorPriv = answers.IssuerPriv
	} else {
		if viper.GetString("validator-pub") == "" {
			questions = append(questions, &survey.Question{
				Name:   "validatorPub",
				Prompt: &survey.Input{Message: "enter peer pub. key"},
			})
		}
		if viper.GetString("validator-priv") == "" {
			questions = append(questions, &survey.Question{
				Name:   "validatorPriv",
				Prompt: &survey.Input{Message: "enter peer priv. key"},
			})
		}

		survey.Ask(questions, &answers)
	}

	validatorPubK := utils.Base58FromStringMust(answers.ValidatorPub)
	issuerPubK := utils.Base58FromStringMust(answers.IssuerPub)

	if viper.GetBool("query-etag") && !viper.GetBool("genesis") {
		issuerAddr, err := chug.NewAddress(issuerPubK)
		must.NoError(err, "invalid issuer address")
		issuerEtag, err := cli.GetHugEtag(issuerAddr)
		must.NoError(err, "%s", err)
		fmt.Printf(rndr.MustRenderf("✅ {-+light+cyan}issuer {-+white}[%s]{-+light+cyan} etag queried: {-+white}%s{-}\n", issuerAddr, issuerEtag))

		validatorAddr, err := chug.NewAddress(validatorPubK)
		must.NoError(err, "invalid validator address")
		validatorEtag, err := cli.GetHugEtag(validatorAddr)
		must.NoError(err, "%s", err)
		fmt.Printf(rndr.MustRenderf("✅ {-+light+cyan}validator {-+white}[%s]{-+light+cyan} etag queried: {-+white}%s{-}\n", validatorAddr, validatorEtag))

		tx.IssuerEtag = issuerEtag
		tx.ValidatorEtag = validatorEtag
	}

	tx.HashTx()
	print.LineTpl("✅ tx hash {{.hash}} generated", print.Fields{"hash": tx.Hash.String()})

	print.LineTpl("✅ lock tx with issuer keys", nil)
	tx.LockIssuer(utils.Base58FromStringMust(answers.IssuerPriv), issuerPubK)

	print.LineTpl("✅ lock tx with validator keys", nil)
	tx.LockValidator(utils.Base58FromStringMust(answers.ValidatorPriv), validatorPubK)

	if isGenesis {
		fmt.Printf(rndr.MustRenderf("✅ {-+light+cyan}{reverse}GENESIS{-+light+cyan} tx generated{-}\n"))
	} else {
		print.LineTpl("✅ tx generated", print.Fields{})
	}

	return tx
}

var txNewCmd = &cobra.Command{
	Use:   "new",
	Short: "create new genesis tx",

	Run: func(cmd *cobra.Command, args []string) {
		tx := askCreateTxQuestions(viper.GetBool("genesis"))

		if viper.GetBool("print") {
			print.LineTpl(`
Version:		 {{.Version}}
Type:			 {{.Type}}
Timestamp:		 {{.Timestamp}}
Hash:			 {{.Hash}}
IssuerPubKey:		 {{.IssuerPubKey}}
IssuerLock:		 {{.IssuerLock}}
IssuerEtag:		 {{.IssuerEtag}}
ValidatorPubKey: 	 {{.ValidatorPubKey}}
ValidatorLock: 		 {{.ValidatorLock}}
ValidatorEtag:	 	 {{.ValidatorEtag}}
Data:			 {{.Data}}			
`, tx)

		}

		if viper.GetBool("send") {
			if viper.GetBool("genesis") {
				fmt.Println("⚠️ genesis transaction cannot be send")
				return
			}

			err := cli.PostTransaction(tx)
			if err != nil {
				fmt.Printf("🚨 %s\n", err.Error())
			}
		}

	},
}

func init() {
	txNewCmd.Flags().StringP("issuer-pub", "", "", "issuer public key")
	txNewCmd.Flags().StringP("issuer-priv", "", "", "issuer private key")
	txNewCmd.Flags().StringP("issuer-etag", "", "", "issuer etag")

	txNewCmd.Flags().StringP("validator-pub", "", "", "validator public key")
	txNewCmd.Flags().StringP("validator-priv", "", "", "validator private key")
	txNewCmd.Flags().StringP("validator-etag", "", "", "validator etag")

	txNewCmd.Flags().StringP("data", "", "", "data to store in the tx")

	txNewCmd.Flags().BoolP("query-etag", "", true, "if true etags will be queried from the blockchain")
	txNewCmd.Flags().BoolP("genesis", "", false, "true if genesis tx else false")
	txNewCmd.Flags().BoolP("print", "", true, "true if the transaction should be printed")
	txNewCmd.Flags().BoolP("send", "", false, "true if the transaction should be sended")

	txCmd.AddCommand(txNewCmd)
}
