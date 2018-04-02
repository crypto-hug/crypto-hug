package utils

import (
	"../../core"
	"../../core/txvalidators"
	"../../persistance"
	"../../prompt"
	color "github.com/fatih/color"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func SetupPrompt(shell *ishell.Shell) {
	prompt.SetupDelegatePrinter(
		func(msg string) { // debzg
			color.Set(color.FgHiBlack, color.Italic)
			defer color.Unset()
			shell.Printf("🔘: %s\n", msg)
		},
		func(msg string) { // info
			color.Set(color.FgWhite)
			defer color.Unset()
			shell.Printf("ℹ️️: %s\n", msg)
		},
		func(msg string) { // say
			color.Set(color.FgGreen, color.Bold)
			defer color.Unset()
			shell.Printf("💬: %s\n", msg)
		},
		func(msg string) { // success
			color.Set(color.FgGreen, color.Bold)
			defer color.Unset()
			shell.Printf("✅: %s\n", msg)
		},
		func(msg string) { // warning
			color.Set(color.FgYellow, color.Bold)
			defer color.Unset()
			shell.Printf("⚠️: %s\n", msg)
		},
		func(msg string) { // panic
			color.Set(color.FgRed, color.Bold)
			defer color.Unset()
			shell.Printf("🚨️: %s\n", msg)

		})

}

func SetupBlockchain(db *hugdb.BoltDb) *core.Blockchain {
	sink := hugdb.NewBoltBlockStore(db)
	txvReg := txvalidators.SharedRegistry()

	result := core.NewBlockchain(sink, txvReg)

	return result
}

func SetupDb() *hugdb.BoltDb {
	var filePath = "./blockhain_data/c_hug.db"
	prompt.Shared().Debug("use blockchain file: %s", filePath)

	var db, err = hugdb.NewBoltDB(filePath)
	if err != nil {
		PanicExit(err)
		return nil
	}

	db.Bootstrap()

	return db
}
