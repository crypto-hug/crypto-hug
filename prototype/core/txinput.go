package core

type TxInput struct{
	PrevTxHash []byte
	PrevTxOutputIdx int
	ScriptSig string
}

func (self *TxInput) CanUnlock(scriptPubKey string) bool{
	return self.ScriptSig == scriptPubKey
}