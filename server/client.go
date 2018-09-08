package server

import (
	"bufio"
	"log"
	"net"
	"time"
)

func ConnectServer(addr *net.UDPAddr, message []byte) {
	conn, err := net.DialUDP(addr.Network(), nil, addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 2048)
	for i := 0; i < 5; i++ {
		_, err := conn.Write(message)
		if err != nil {
			log.Println(err)
		}
		n, err := bufio.NewReader(conn).Read(buffer)
		if err != nil {
			log.Println(err)
			continue
		}

		if n > 0 {
			log.Printf("Message from server: %s\n", buffer[:n])
		}
		time.Sleep(1 * time.Second)
	}
}
