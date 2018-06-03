package go2p

type Message struct {
	senderId   string
	receiverId string
	header     *MessageHeader
	body       []byte
}

type MessageHeader struct {
	MessageSize uint32
	HeaderSize  uint32
	BodySize    uint32
	Fields      []*HeaderField
}

type HeaderField struct {
	Name  string
	value interface{}
}
