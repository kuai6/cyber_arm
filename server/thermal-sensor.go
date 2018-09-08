package server

import (
	"log"
	"net"
)

func StreamThermalSensorData(addr *net.UDPAddr, getData func() []byte) {
	buffer := make([]byte, 2048)
	packetConn, err := net.ListenPacket(addr.Network(), addr.String())
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			_, clientAddr, err := packetConn.ReadFrom(buffer)
			if err != nil {
				log.Println(err)
				continue
			}
			go packetConn.WriteTo(getData(), clientAddr)
		}
	}()
}
