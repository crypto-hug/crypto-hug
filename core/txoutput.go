package core

type TxOutput struct{
	Value int
	ScriptPubKey string
}


func (self *TxOutput) CanBeUnlocked(scriptSig string) bool{
	return self.ScriptPubKey == scriptSig
}