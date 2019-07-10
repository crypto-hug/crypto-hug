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
	Fields       HeaderFields
}

type HeaderField struct {
	Name  string
	value interface{}
}

type HeaderFields []*HeaderField

func (self HeaderFields) Has(name string) bool {
	for _, h := range self {
		if h.Name == name {
			return true
		}
	}

	return false
}
func (self HeaderFields) GetString(name string) string {
	val := self.Get(name)
	if val == nil {
		return ""
	}

	return val.(string)
}

func (self HeaderFields) Get(name string) interface{} {
	for _, h := range self {
		if h.Name == name {
			return h.value.(string)
		}
	}

	return ""
}

// func (self *Message) Path() string{
// }

func (self *MessageHeader) FieldPath() string {
	p := self.Fields.GetString("PATH")
	return p
}
