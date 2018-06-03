package go2p

type Message struct {
	header *MessageHeader
	body   []byte
}

type MessageHeader struct {
	senderAddr   string
	receiverAddr string
	MessageSize  uint32
	HeaderSize   uint32
	BodySize     uint32
	Fields       []*HeaderField
}

type HeaderField struct {
	Name  string
	value interface{}
}
