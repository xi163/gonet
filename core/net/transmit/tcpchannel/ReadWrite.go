package tcpchannel

import (
	"net"

	logs "github.com/cwloo/gonet/logs"
)

// 读指定长度
func Read(conn net.Conn, buf []byte) (int, error) {
	size := len(buf)
	n, err := conn.Read(buf[0:size])
	if err != nil {
		// logs.Debugf("%v", err)
		return n, err
	}
	// buf = buf[0:n]
	return n, nil
}

// 读指定长度
func ReadFull(conn net.Conn, buf []byte) error {
	length := 0
	size := len(buf)
	for {
		n, err := conn.Read(buf[length:size])
		if err != nil {
			// logs.Debugf("%v", err)
			return err
		}
		length += n
		if length == size {
			return nil
		}
	}
}

// 写指定长度
func WriteFull(conn net.Conn, buf []byte) error {
	length := 0
	size := len(buf)
	for {
		n, err := conn.Write(buf[length:size])
		if err != nil {
			logs.Debugf("%v", err)
			return err
		}
		length += n
		if length == size {
			return nil
		}
	}
}
