package packet

import (
	"encoding/binary"
	"errors"
)

const (
	SESSIONSZ = 32
	AESKEYSZ  = 16
	SERVIDSZ  = 50
)

const (
	ENC_JSON_NONE         uint8 = 0x01
	ENC_PROTOBUF_NONE     uint8 = 0x02
	ENC_JSON_BIT_MASK     uint8 = 0x11
	ENC_PROTOBUF_BIT_MASK uint8 = 0x12
	ENC_JSON_RSA          uint8 = 0x21
	ENC_PROTOBUF_RSA      uint8 = 0x22
	ENC_JSON_AES          uint8 = 0x31
	ENC_PROTOBUF_AES      uint8 = 0x32
)

// unsafe.Sizeof/binary.Size
// type Header struct {
// 	Len      [2]byte //包总长度
// 	Crc      [2]byte //CRC校验位
// 	Ver      [2]byte //版本号
// 	Sign     [2]byte //签名
// 	MainID   byte    //主消息mainID
// 	SubID    byte    //子消息subID
// 	EncType  byte    //加密类型
// 	Reserved byte    //预留
// 	ReqID    [4]byte
// 	RealSize [2]byte //用户数据长度
// }

// type InternalPrevHeader struct {
// 	Len      [2]byte
// 	Kicking  [2]int8
// 	Ok       [4]int8
// 	UserID   [8]int8
// 	Ipaddr   uint32           //来自真实IP
// 	Session  [SESSIONSZ]uint8 //用户会话
// 	Aeskey   [AESKEYSZ]uint8  //AES_KEY
// 	// ServID   [SERVIDSZ]uint8  //来自节点ID
// 	Checksum uint16           //校验和CHKSUM
// }

func Pack(msg *Msg, order binary.ByteOrder) ([]byte, error) {
	//len，2字节
	length := 18 + len(msg.Data)
	b := make([]byte, length)
	order.PutUint16(b[0:], uint16(length))
	//版本0x0001
	order.PutUint16(b[4:], uint16(msg.ver))
	//标记0x5F5F
	order.PutUint16(b[6:], uint16(msg.sign))
	//主命令ID
	b[8] = byte(msg.mainID)
	//子命令ID
	b[9] = byte(msg.subID)
	//加密类型
	b[10] = byte(msg.encType)
	//预留字段
	b[11] = byte(0x01)
	//请求ID
	order.PutUint32(b[12:], uint32(0))
	//实际大小(json/protobuf)
	order.PutUint16(b[16:], uint16(len(msg.Data)))
	//实际数据(json/protobuf)
	copy(b[18:], msg.Data)
	//CRC，2字节
	crc := GetChecksum(b[4:])
	order.PutUint16(b[2:], crc)
	return b, nil
}

func Unpack(b []byte, order binary.ByteOrder) (uint32, []byte, error) {
	//len，2字节
	length := order.Uint16(b[:2])
	if length != uint16(len(b)) {
		return 0, nil, errors.New("parse error")
	}
	//CRC，2字节
	chsum := order.Uint16(b[2:])
	//CRC校验
	crc := GetChecksum(b[4:])
	if crc != chsum {
		return 0, nil, errors.New("checksum error")
	}
	// //版本0x0001
	// ver := order.Uint16(b[4:])
	// //标记0x5F5F
	// sign := order.Uint16(b[6:])
	//主命令ID
	mainID := uint8(b[8])
	//子命令ID
	subID := uint8(b[9])
	// //加密类型
	// encType := uint8(b[10])
	// //预留字段
	// reserved := uint8(b[11])
	// //请求ID
	// reqID := order.Uint32(b[12:16])
	// //实际大小(json/protobuf)
	// realSize := order.Uint16(b[16:18])
	// logs.Debugf("ver:%#x\nsign:%#x\nmainID:%d\nsubID:%d\nencTy:%#x\nreserv:%d\nreqID:%d\nrealSize:%d",
	// 	ver, sign, mainID, subID, encType, reserved, reqID, realSize)
	// 实际数据(json/protobuf)
	data := b[18:]
	cmd := uint32(Enword(int(mainID), int(subID)))
	return cmd, data, nil
}
