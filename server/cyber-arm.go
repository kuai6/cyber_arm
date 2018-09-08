package server

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type Command struct {
	Name      string   `json:"name"`
	Arguments []string `json:"arguments,omitempty"`
}

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

			commandData := buffer[:n]
			log.Printf("Got command: %s\n", string(commandData))

			command := new(Command)
			if err := json.Unmarshal(commandData, command); err != nil {
				log.Printf("Failed to unmarshal command: %s", err)
			}
			go handleCommand(command)

			packetConn.WriteTo(nil, clientAddr)
		}
	}()
}

func handleCommand(command *Command) {
	switch command.Name {
	case "ROTATE":
		alpha, err := strconv.ParseFloat(command.Arguments[0], 64)
		if err != nil {
			log.Printf("Failed to parse argument: %s", err)
		}
		beta, err := strconv.ParseFloat(command.Arguments[1], 64)
		if err != nil {
			log.Printf("Failed to parse argument: %s", err)
		}
		log.Printf("Perform cyber-arm rotation to (%f,%f)\n", alpha, beta)
		//rotate(alpha, beta)
	case "FIRE":
		log.Printf("Perform fire action\n")
		//fire()
	}
}
