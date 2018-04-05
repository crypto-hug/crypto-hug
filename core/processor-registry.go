package core

import ()

type TransactionProcessorRegistry struct {
	processors TransactionProcessors
}

func NewTransactionProcessorRegistry(processors TransactionProcessors) *TransactionProcessorRegistry {
	result := TransactionProcessorRegistry{processors: processors}
	return &result
}

func (self *TransactionProcessorRegistry) Get(tx *Transaction) TransactionProcessors {
	result := TransactionProcessors{}
	for _, processor := range self.processors {
		if processor.ShouldProcess(tx) {
			result = append(result, processor)
		}
	}

	return result
}
