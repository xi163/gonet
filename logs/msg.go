package logs

import "sync"

var (
	msg = sync.Pool{
		New: func() any {
			return &Msg{}
		},
	}
	flags = sync.Pool{
		New: func() any {
			return &Flags{}
		},
	}
	message = sync.Pool{
		New: func() any {
			return &Message{}
		},
	}
	messageT = sync.Pool{
		New: func() any {
			return &MessageT{}
		},
	}
)

type Msg struct {
	first  string
	second string
}

func NewMsg(first, second string) *Msg {
	s := msg.Get().(*Msg)
	s.first = first
	s.second = second
	return s
}

func (s *Msg) Put() {
	msg.Put(s)
}

type Flags struct {
	first  int
	second Style
}

func NewFlags(first int, second Style) *Flags {
	s := flags.Get().(*Flags)
	s.first = first
	s.second = second
	return s
}

func (s *Flags) Put() {
	flags.Put(s)
}

type Message struct {
	first  *Msg
	second string
}

func NewMessage(first *Msg, second string) *Message {
	s := message.Get().(*Message)
	s.first = first
	s.second = second
	return s
}

func (s *Message) Put() {
	message.Put(s)
}

type MessageT struct {
	first  *Message
	second *Flags
}

func NewMessageT(first *Message, second *Flags) *MessageT {
	s := messageT.Get().(*MessageT)
	s.first = first
	s.second = second
	return s
}

func (s *MessageT) Put() {
	messageT.Put(s)
}
