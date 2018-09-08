package server

import (
	"log"
	"net"
)

func ListenCyberArmCommands(addr *net.UDPAddr) {
	buffer := make([]byte, 2048)
	packetConn, err := net.ListenPacket(addr.Network(), addr.String())
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			n, clientAddr, err := packetConn.ReadFrom(buffer)
			if err != nil {
				log.Println(err)
				continue
			}
			packetConn.WriteTo(nil, clientAddr)
			log.Printf("Got command: %s, perfomring some action...\n", buffer[:n])
		}
	}()
}
