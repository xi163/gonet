package packet

import "google.golang.org/protobuf/proto"

type Msg struct {
	ver     uint16
	sign    uint16
	encType uint8
	mainID  uint8
	subID   uint8
	Data    []byte
}

func New(mainID uint8, subID uint8, v proto.Message) *Msg {
	b, err := proto.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return &Msg{
		ver:     0x0001,
		sign:    0x5F5F,
		encType: 0x02,
		mainID:  mainID,
		subID:   subID,
		Data:    b,
	}
}

func Enword(mainID, subID int) int {
	return ((0xFF & mainID) << 8) | (0xFF & subID)
}

func Deword(cmd int) (mainID, subID int) {
	mainID = (0xFF & (cmd >> 8))
	subID = (0xFF & cmd)
	return
}

const (
	SubCmdID = iota + ((0xFF & 20) << 8) | (0xFF & 19)
)
