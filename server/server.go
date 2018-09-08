package server

import (
	"github.com/kuai6/cyber_arm/config"
	"fmt"
	"log"
	"net"
)

type Server struct {
	MessageBuffer  []byte
	Connection     net.PacketConn
	MessageChannel chan *Message
}

func NewServer(config *config.ServerConfiguration) (*Server, chan *Message, error) {
	connection, err := net.ListenPacket("udp", config.Address())
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to start UDP server: %s", err)
	}

	log.Printf("Started UDP server: %s\n", connection.LocalAddr().String())

	buffer := make([]byte, config.MessageBufferSize)
	messageChannel := make(chan *Message, config.MessageQueueSize)

	return &Server{
		MessageBuffer:  buffer,
		Connection:     connection,
		MessageChannel: messageChannel,
	}, messageChannel, nil
}

func (s *Server) ReadMessage() {
	n, addr, err := s.Connection.ReadFrom(s.MessageBuffer)
	if err != nil {
		log.Printf("Failed to receive UDP packet: %s\n", err)
	}

	s.MessageChannel <- &Message{s.MessageBuffer[0:n], addr}
}

func (s *Server) SendMessage(message *Message) error {
	_, err := s.Connection.WriteTo(message.Data, message.Address)
	if err != nil {
		return fmt.Errorf("Failed to send packet: %s\n", err)
	}

	return nil
}
