package server

import (
	"encoding/json"
	"github.com/kuai6/cyber_arm/command"
	"log"
	"net"
)

func ListenCyberArmCommands(addr *net.UDPAddr, handler func(command *command.Command)) {
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

			commandData := buffer[:n]
			log.Printf("Got command: %s\n", string(commandData))

			c := new(command.Command)
			if err := json.Unmarshal(commandData, c); err != nil {
				log.Printf("Failed to unmarshal command: %s", err)
			}
			go handler(c)

			packetConn.WriteTo(nil, clientAddr)
		}
	}()
}
