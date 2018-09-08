package server

import "net"

type Message struct {
	Data    []byte
	Address net.Addr
}
