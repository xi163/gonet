package conn

import (
	"fmt"
	"strings"
)

// 网络地址结构
type Address struct {
	Proto string //ws/wss/tcp
	Addr  string //ip:port localhost:port
	Path  string
	Ip    string
	Port  string
}

func (s *Address) Format() (addr string) {
	switch s.Proto {
	case "ws", "wss":
		addr = fmt.Sprintf("%v://%v%v", s.Proto, s.Addr, s.Path)
	case "tcp":
		addr = fmt.Sprintf("%v://%v", s.Proto, s.Addr)
	}
	return
}

func ParseAddress(address string) *Address {
	//ws://ip:port/path wss://ip:port/path
	//ws://localhost:port/path wss://localhost:port/path
	vec := strings.Split(address, "//")
	switch len(vec) == 2 {
	case true:
		proto := strings.ToLower(strings.Trim(vec[0], ":"))
		switch proto {
		case "ws", "wss":
			v := strings.Split(vec[1], "/")
			switch len(v) {
			case 1:
				addr := v[0]
				path := "/"
				Ip := ""
				host := strings.Split(addr, ":")
				if len(host) == 2 {
					if host[0] == "localhost" {
						Ip = "127.0.0.1"
					} else {
						Ip = host[0]
					}
					port := host[1]
					return &Address{Proto: proto, Addr: addr, Path: path, Ip: Ip, Port: port}
				} else {
					panic("parse " + address + " error")
				}
			case 2:
				addr := v[0]
				path := v[1]
				if path == "" {
					path = "/"
				} else {
					path = "/" + path
				}
				Ip := ""
				host := strings.Split(addr, ":")
				if len(host) == 2 {
					if host[0] == "localhost" {
						Ip = "127.0.0.1"
					} else {
						Ip = host[0]
					}
					port := host[1]
					return &Address{Proto: proto, Addr: addr, Path: path, Ip: Ip, Port: port}
				} else {
					panic("parse " + address + " error")
				}
			default:
				panic("parse " + address + " error")
			}
		default:
			//tcp://ip:port tcp://localhost:port
			addr := vec[1]
			host := strings.Split(addr, ":")
			Ip := ""
			if len(host) == 2 {
				if host[0] == "localhost" {
					Ip = "127.0.0.1"
				} else {
					Ip = host[0]
				}
				port := host[1]
				return &Address{Proto: proto, Addr: addr, Ip: Ip, Port: port}
			}
		}
	default:
		//ip:port localhost:port
		proto := "tcp"
		addr := address
		host := strings.Split(addr, ":")
		Ip := ""
		if len(host) == 2 {
			if host[0] == "localhost" {
				Ip = "127.0.0.1"
			} else {
				Ip = host[0]
			}
			port := host[1]
			return &Address{Proto: proto, Addr: addr, Ip: Ip, Port: port}
		} else {
			panic("parse " + address + " error")
		}
	}
	return nil
}
