package core

import (
	"../common/crypt"
	"../common/formatters"
)

const Reward = 1

type Transaction struct{
	Hash 	[]byte
	Inputs 	[]TxInput
	Outputs	[]TxOutput
}

type Transactions []*Transaction


func (self *Transaction) getInputsHash() []byte{
	var all [][]byte
	for _, input := range self.Inputs{
		var hash = crypt.AllBytesHash(
			input.PrevTxHash,
			[]byte(formatters.HexString(int64(input.PrevTxOutputIdx))),
			[]byte(input.ScriptSig),
		)

		all = append(all, hash)
	}


	var result = crypt.AllBytesHash(all[:]...)
	return result
}
func (self *Transaction) getOutputsHash() []byte{
	var all [][]byte

	for _, output := range self.Outputs{
		var hash = crypt.AllBytesHash(
			[]byte(formatters.HexString(int64(output.Value))),
			[]byte(output.ScriptPubKey),
		)

		all = append(all, hash)
	}

	var result = crypt.AllBytesHash(all[:]...)
	return result
}

func (self Transactions) getHash() []byte{
	var all [][]byte
	for _, tx := range self{
		var hash = tx.getHash()

		all = append(all, hash)
	}

	var result = crypt.AllBytesHash(all[:]...)
	return result
}


func (self *Transaction) getHash() []byte {
	var inputsHash = self.getInputsHash()
	var outputsHash = self.getOutputsHash()
	var result = crypt.AllBytesHash(inputsHash, outputsHash)

	return result
}

func (self *Transaction) IsCoinbase() bool{
	if len(self.Inputs) != 1{
		return false
	}
	if len(self.Inputs[0].PrevTxHash) != 0{
		return false
	}
	if self.Inputs[0].PrevTxOutputIdx != -1{
		return false
	}

	return true
}

func (self *Transaction) getOutValue() int{
	var result = 0
	for _, output := range self.Outputs{
		result += output.Value
	}

	return result
}

func NewCoinbaseTransaction(to string, data string) *Transaction  {
	var in = TxInput{PrevTxHash:[]byte{}, PrevTxOutputIdx:-1, ScriptSig:data}
	var out = TxOutput{Value:Reward, ScriptPubKey: to}
	var result = Transaction{Hash:nil, Inputs:[]TxInput{in}, Outputs:[]TxOutput{out}}
	result.Hash = result.getHash()

	return &result
}